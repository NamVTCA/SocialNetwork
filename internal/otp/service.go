package otp

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"socialnetwork/models"
	"socialnetwork/pkg/email"
	"socialnetwork/pkg/sms"
	"socialnetwork/pkg/utils"
)

type Service interface {
	SaveOTP(ctx context.Context, key string, code string, duration time.Duration) error
	VerifyOTP(ctx context.Context, req *models.VerifyOTPRequest) error
	DeleteOTP(ctx context.Context, key string) error
	SendOTP(ctx context.Context, req *models.SendOTPRequest) error
	SendOTPWithCustomKey(ctx context.Context, req *models.SendOTPRequest) error // ← thêm
	PeekIdentifierByCustomKey(ctx context.Context, key string) (string, error)
	SendForgotPasswordOTP(ctx context.Context, email string) error
}

type OTPService interface {
	SaveOTP(ctx context.Context, key string, code string, duration time.Duration) error
	VerifyOTP(ctx context.Context, req *models.VerifyOTPRequest) error
	DeleteOTP(ctx context.Context, key string) error
	SendOTP(ctx context.Context, req *models.SendOTPRequest) error
	SendForgotPasswordOTP(ctx context.Context, email string) error
	SendRawEmail(ctx context.Context, to, subject, body string) error
	GetRawOTP(ctx context.Context, key string) (string, error)
}

type service struct {
	redisClient *redis.Client
	emailSender email.EmailSender
	smsSender   sms.SMSSender
}

func normalizePhone(phone string) string {
	// Nếu bắt đầu bằng 0 → đổi thành +84
	if strings.HasPrefix(phone, "0") {
		return "+84" + phone[1:]
	}
	return phone
}

func NewService(redisClient *redis.Client, emailSender email.EmailSender, smsSender sms.SMSSender) Service {
	return &service{
		redisClient: redisClient,
		emailSender: emailSender,
		smsSender:   smsSender,
	}
}

func GenerateOTP(length int) string {
	rand.Seed(time.Now().UnixNano())
	format := "%06d"
	if length == 4 {
		format = "%04d"
	}
	max := 1000000
	if length == 4 {
		max = 10000
	}
	return fmt.Sprintf(format, rand.Intn(max))
}

// ✅ Chuẩn hóa identifier để dùng làm key
func normalizeIdentifier(channel, identifier string) (string, error) {
	if channel == "phone" {
		formatted := utils.FormatPhoneToE164(identifier)
		if !utils.IsValidPhoneE164(formatted) {
			return "", errors.New("invalid phone number format")
		}
		return formatted, nil
	}
	return identifier, nil
}

func (s *service) SendOTP(ctx context.Context, req *models.SendOTPRequest) error {
	normalizedID := req.Identifier
	if req.Channel == "phone" {
		normalizedID = utils.FormatPhoneToE164(req.Identifier)
		if !utils.IsValidPhoneE164(normalizedID) {
			return errors.New("invalid phone number format")
		}
	}

	otp := GenerateOTP(6)
	key := fmt.Sprintf("otp:%s:%s", normalizedID, req.Purpose)
	expire := 5 * time.Minute

	// fmt.Printf("[DEBUG][SendOTP] key = %s, otp = %s\n", key, otp)

	if err := s.redisClient.Set(ctx, key, otp, expire).Err(); err != nil {
		return err
	}

	message := fmt.Sprintf("Your OTP code is %s. It is valid for 5 minutes.", otp)

	switch req.Channel {
	case "email":
		// fmt.Println("[DEBUG][SendOTP] Sending OTP via email")
		return s.emailSender.Send(req.Identifier, "OTP Verification", message)
	case "phone":
		// fmt.Println("[DEBUG][SendOTP] Sending OTP via SMS")
		return s.smsSender.Send(normalizedID, message)
	default:
		return errors.New("invalid channel")
	}
}

func (s *service) SendOTPWithCustomKey(ctx context.Context, req *models.SendOTPRequest) error {
	normalizedID := req.Identifier
	if req.Channel == "phone" {
		normalizedID = utils.FormatPhoneToE164(req.Identifier)
		if !utils.IsValidPhoneE164(normalizedID) {
			return errors.New("invalid phone number format")
		}
	}

	otp := GenerateOTP(6)
	key := req.CustomKey
	if key == "" {
		// fallback nếu không có custom key
		key = fmt.Sprintf("otp:%s:%s", normalizedID, req.Purpose)
	}
	expire := 5 * time.Minute

	if err := s.redisClient.Set(ctx, key, normalizedID+":"+otp, expire).Err(); err != nil {
		return err
	}

	message := fmt.Sprintf("Your OTP code is %s. It is valid for 5 minutes.", otp)

	switch req.Channel {
	case "email":
		return s.emailSender.Send(req.Identifier, "OTP Verification", message)
	case "phone":
		return s.smsSender.Send(normalizedID, message)
	default:
		return errors.New("invalid channel")
	}
}

func (s *service) PeekIdentifierByCustomKey(ctx context.Context, key string) (string, error) {
	val, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", errors.New("OTP không tồn tại hoặc đã hết hạn")
	}

	parts := strings.Split(val, ":")
	if len(parts) != 2 {
		return "", errors.New("dữ liệu OTP không hợp lệ")
	}

	return parts[0], nil // Trả về identifier (email/số điện thoại)
}

func (s *service) SaveOTP(ctx context.Context, key, code string, duration time.Duration) error {
	return s.redisClient.Set(ctx, key, code, duration).Err()
}

func (s *service) DeleteOTP(ctx context.Context, key string) error {
	return s.redisClient.Del(ctx, key).Err()
}

func (s *service) VerifyOTP(ctx context.Context, req *models.VerifyOTPRequest) error {
	// fmt.Printf("[DEBUG][VerifyOTP] identifier = %s, purpose = %s, input OTP = %s\n", req.Identifier, req.Purpose, req.OTP)

	// ✅ Chuẩn hóa lại identifier theo đúng channel
	normalizedID, err := normalizeIdentifier(req.Channel, req.Identifier)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("otp:%s:%s", normalizedID, req.Purpose)

	storedOTP, err := s.redisClient.Get(ctx, key).Result()
	if err != nil {
		// fmt.Printf("[DEBUG][VerifyOTP] Không tìm thấy OTP trong Redis với key = %s\n", key)
		return errors.New("mã OTP không hợp lệ hoặc đã hết hạn")
	}

	// fmt.Printf("[DEBUG][VerifyOTP] stored OTP = %s\n", storedOTP)

	if storedOTP != req.OTP {
		// fmt.Println("[DEBUG][VerifyOTP] OTP KHÔNG khớp!")
		return errors.New("mã OTP không chính xác")
	}

	// fmt.Println("[DEBUG][VerifyOTP] OTP hợp lệ.")
	return nil
}

func (s *service) SendForgotPasswordOTP(ctx context.Context, email string) error {
	req := &models.SendOTPRequest{
		Identifier: email,
		Purpose:    "forgot_password",
		Channel:    "email",
	}
	return s.SendOTP(ctx, req)
}

func (s *service) SendRawEmail(ctx context.Context, to, subject, body string) error {
	return s.emailSender.Send(to, subject, body)
}

func (s *service) GetRawOTP(ctx context.Context, key string) (string, error) {
	return s.redisClient.Get(ctx, key).Result()
}

