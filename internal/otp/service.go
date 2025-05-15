package otp

import (
    "context"
    "errors"
    "fmt"
    "math/rand"
    "time"

    "github.com/go-redis/redis/v8"
	"socialnetwork/models"
    "socialnetwork/pkg/email"
    "socialnetwork/pkg/sms"
)

type Service interface {
    SendOTP(ctx context.Context, req *models.SendOTPRequest) error
    VerifyOTP(ctx context.Context, req *models.VerifyOTPRequest) error
}

type service struct {
    redisClient *redis.Client
    emailSender email.EmailSender
    smsSender   sms.SMSSender
}

func NewService(redisClient *redis.Client, emailSender email.EmailSender, smsSender sms.SMSSender) Service {
    return &service{
        redisClient: redisClient,
        emailSender: emailSender,
        smsSender:   smsSender,
    }
}

func generateOTP() string {
    rand.Seed(time.Now().UnixNano())
    return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (s *service) SendOTP(ctx context.Context, req *models.SendOTPRequest) error {
    otp := generateOTP()
    key := fmt.Sprintf("otp:%s:%s", req.Identifier, req.Purpose)
    expire := time.Minute * 5

    err := s.redisClient.Set(ctx, key, otp, expire).Err()
    if err != nil {
        return err
    }

    message := fmt.Sprintf("Your OTP code is %s. It is valid for 5 minutes.", otp)
    switch req.Channel {
    case "email":
        return s.emailSender.Send(req.Identifier, "OTP Verification", message)
    case "phone":
        return s.smsSender.Send(req.Identifier, message)
    default:
        return errors.New("invalid channel")
    }
}

func (s *service) VerifyOTP(ctx context.Context, req *models.VerifyOTPRequest) error {
    key := fmt.Sprintf("otp:%s:%s", req.Identifier, req.Purpose)
    otp, err := s.redisClient.Get(ctx, key).Result()
    if err != nil {
        return errors.New("OTP expired or invalid")
    }

    if otp != req.OTP {
        return errors.New("OTP incorrect")
    }

    // Xóa OTP sau khi dùng thành công
    s.redisClient.Del(ctx, key)
    return nil
}
