package user

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	// "fmt"
	"socialnetwork/dto/request"
	"socialnetwork/models"
	"socialnetwork/pkg/auth"
	"socialnetwork/pkg/email"

	"time"
)

type Service interface {
	Register(ctx context.Context, user *models.User) error
	Login(ctx context.Context, email, password string) (string, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateProfile(ctx context.Context, id string, req *request.UpdateProfileRequest) error
	ChangePassword(ctx context.Context, userID string, req *request.ChangePasswordRequest) error
	SendForgotPasswordOTP(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, req *request.ResetPasswordRequest) error
	ChangeEmailRequest(ctx context.Context, req *request.ChangeEmailRequest) error
	VerifyEmailRequest(ctx context.Context, req *request.VerifyEmailRequest) error
}

type OTPService interface {
	SaveOTP(ctx context.Context, key string, code string, duration time.Duration) error
	VerifyOTP(ctx context.Context, req *models.VerifyOTPRequest) error
	DeleteOTP(ctx context.Context, key string) error
	SendOTP(ctx context.Context, req *models.SendOTPRequest) error
	SendForgotPasswordOTP(ctx context.Context, email string) error
}

type service struct {
	repo        Repository
	otpService  OTPService // interface quản lý OTP
	emailSender email.EmailSender
}

func NewService(repo Repository, otpService OTPService, emailSender email.EmailSender) Service {
	return &service{
		repo:        repo,
		otpService:  otpService,
		emailSender: emailSender,
	}
}

func (s *service) Register(ctx context.Context, user *models.User) error {
	hashedPassword, err := auth.HashPassword(user.Password)
	// fmt.Println("Password trước khi lưu:", hashedPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	return s.repo.Create(ctx, user)
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("tài khoản không tồn tại")
	}

	// fmt.Println("DEBUG LOGIN")
	// fmt.Println("Mật khẩu người dùng nhập:", password)
	// fmt.Println("Hash lưu trong DB:", user.Password)
	// fmt.Println("Check result:", auth.CheckPasswordHash(password, user.Password))

	if !auth.CheckPasswordHash(password, user.Password) {
		return "", errors.New("mật khẩu không đúng")
	}

	token, err := auth.GenerateJWT(user.ID.Hex())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) GetByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *service) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) UpdateProfile(ctx context.Context, id string, req *request.UpdateProfileRequest) error {
	update := bson.M{}

	if req.DisplayName != "" {
		update["displayName"] = req.DisplayName
	}
	if req.Bio != "" {
		update["bio"] = req.Bio
	}
	if req.Gender != "" {
		update["gender"] = req.Gender
	}
	if req.BirthDate != nil {
		update["birthDate"] = req.BirthDate
	}
	if req.AvatarURL != "" {
		update["avatarUrl"] = req.AvatarURL
	}
	if req.CoverURL != "" {
		update["coverUrl"] = req.CoverURL
	}
	if req.Location != "" {
		update["location"] = req.Location
	}
	if req.Website != "" {
		update["website"] = req.Website
	}
	if req.Phone != "" {
		update["phone"] = req.Phone
	}

	if len(update) == 0 {
		return nil // không có gì để cập nhật
	}

	return s.repo.UpdateByID(ctx, id, update)
}

func (s *service) ChangePassword(ctx context.Context, userID string, req *request.ChangePasswordRequest) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if !auth.CheckPasswordHash(req.OldPassword, user.Password) {
		return errors.New("mật khẩu cũ không đúng")
	}

	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	update := bson.M{
		"password": hashedPassword,
	}

	return s.repo.UpdateByID(ctx, userID, update)
}

// Gửi mã OTP quên mật khẩu
func (s *service) SendForgotPasswordOTP(ctx context.Context, email string) error {
	_, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("email không tồn tại")
	}

	return s.otpService.SendForgotPasswordOTP(ctx, email)
}

// Reset mật khẩu bằng OTP
func (s *service) ResetPassword(ctx context.Context, req *request.ResetPasswordRequest) error {
	// Tạo VerifyOTPRequest từ req
	verifyReq := &models.VerifyOTPRequest{
		Identifier: req.Email,
		Purpose:    "forgot_password",
		OTP:        req.OTP,
		Channel:    "email",
	}

	// Kiểm tra OTP hợp lệ
	if err := s.otpService.VerifyOTP(ctx, verifyReq); err != nil {
		return errors.New("mã OTP không hợp lệ hoặc đã hết hạn")
	}

	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("Email không tồn tại")
	}

	hashedPassword, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	update := bson.M{"password": hashedPassword}
	if err := s.repo.UpdateByID(ctx, user.ID.Hex(), update); err != nil {
		return err
	}

	// Xoá OTP sau khi dùng
	s.otpService.DeleteOTP(ctx, "forgot_password:"+req.Email)

	return nil
}

func (s *service) ChangeEmailRequest(ctx context.Context, req *request.ChangeEmailRequest) error {
	// Kiểm tra xem email mới đã tồn tại chưa
	existingUser, err := s.repo.FindByEmail(ctx, req.NewEmail)
	if err == nil && existingUser != nil {
		return errors.New("email mới đã tồn tại")
	}

	// Gửi mã OTP đến email mới
	err = s.otpService.SendOTP(ctx, &models.SendOTPRequest{
		Identifier: req.NewEmail,
		Purpose:    "change_email",
		Channel:    "email",
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *service) VerifyEmailRequest(ctx context.Context, req *request.VerifyEmailRequest) error {
	// Kiểm tra OTP
	verifyReq := &models.VerifyOTPRequest{
		Identifier: req.NewEmail,
		Purpose:    "change_email",
		OTP:        req.OTP,
		Channel:    "email",
	}

	if err := s.otpService.VerifyOTP(ctx, verifyReq); err != nil {
		return errors.New("mã OTP không hợp lệ hoặc đã hết hạn")
	}

	user, err := s.repo.FindByID(ctx, req.UserID)
	if err != nil {
		return errors.New("người dùng không tồn tại")
	}

	update := bson.M{"email": req.NewEmail}
	if err := s.repo.UpdateByID(ctx, user.ID.Hex(), update); err != nil {
		return err
	}

	s.otpService.DeleteOTP(ctx, "change_email:"+req.NewEmail)

	return nil
}