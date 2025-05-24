package user

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"socialnetwork/dto/request"
	"socialnetwork/internal/otp"
	"socialnetwork/models"
	"socialnetwork/pkg/auth"
	"socialnetwork/pkg/email"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Service interface {
	Register(ctx context.Context, user *models.User) error
	LoginEmail(ctx context.Context, email, password string) (string, error)
	LoginPhone(ctx context.Context, phone, password string) (string, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateProfile(ctx context.Context, id string, req *request.UpdateProfileRequest) error
	ChangePassword(ctx context.Context, userID string, req *request.ChangePasswordRequest) error
	SendForgotPasswordOTP(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, req *request.ResetPasswordRequest) error
	ChangeEmailRequest(ctx context.Context, userID string, req *request.ChangeEmailRequest) error
	VerifyEmailRequest(ctx context.Context, userID string, req *request.VerifyEmailRequest) error
	SendFriendRequest(ctx context.Context, fromID, toID primitive.ObjectID) error
	AcceptFriendRequest(ctx context.Context, userID, requesterID primitive.ObjectID) error
	BlockUser(ctx context.Context, userID, targetID primitive.ObjectID) error
	ToggleHideProfile(ctx context.Context, userID primitive.ObjectID, hide bool) error
	CancelFriendRequest(ctx context.Context, fromID, toID primitive.ObjectID) error
	FriendRequestExists(ctx context.Context, fromID, toID primitive.ObjectID) (bool, error)
}

type service struct {
	repo        Repository
	otpService  otp.Service // interface quản lý OTP
	emailSender email.EmailSender
}

func NewService(repo Repository, otpService otp.Service, emailSender email.EmailSender) Service {
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

func (s *service) LoginEmail(ctx context.Context, email, password string) (string, error) {
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

func (s *service) LoginPhone(ctx context.Context, phone, password string) (string, error) {
	user, err := s.repo.FindByPhone(ctx, phone)
	if err != nil {
		return "", errors.New("tài khoản không tồn tại")
	}
	if user == nil {
		return "", errors.New("tài khoản không tồn tại")
	}

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

	return s.repo.UpdateByID(ctx, id, bson.M{"$set": update})
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
		"$set": bson.M{
			"password": hashedPassword,
		},
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

	update := bson.M{
		"$set": bson.M{
			"password": hashedPassword,
		},
	}

	if err := s.repo.UpdateByID(ctx, user.ID.Hex(), update); err != nil {
		return err
	}

	// Xoá OTP sau khi dùng
	s.otpService.DeleteOTP(ctx, "forgot_password:"+req.Email)

	return nil
}

func (s *service) ChangeEmailRequest(ctx context.Context, userID string, req *request.ChangeEmailRequest) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("người dùng không tồn tại")
	}

	if user.Email != req.OldEmail {
		return errors.New("email hiện tại không đúng")
	}

	existingUser, err := s.repo.FindByEmail(ctx, req.NewEmail)
	if err == nil && existingUser != nil {
		return errors.New("email mới đã tồn tại")
	}

	// 🔐 Tạo OTP
	otp := otp.GenerateOTP(6)

	// 🔑 Lưu vào Redis: key = change_email:<userID>, value = <newEmail>:<otp>
	key := fmt.Sprintf("change_email:%s", user.ID.Hex())
	value := fmt.Sprintf("%s:%s", req.NewEmail, otp)

	err = s.otpService.SaveOTP(ctx, key, value, 5*time.Minute)
	if err != nil {
		return errors.New("không thể lưu mã OTP")
	}

	// ✉️ Gửi OTP qua email mới
	message := fmt.Sprintf("Mã xác thực thay đổi email của bạn là %s. Có hiệu lực trong 5 phút.", otp)
	err = s.otpService.SendRawEmail(ctx, req.NewEmail, "Xác thực thay đổi email", message)
	if err != nil {
		return errors.New("không thể gửi email xác thực")
	}

	return nil
}

func (s *service) VerifyEmailRequest(ctx context.Context, userID string, req *request.VerifyEmailRequest) error {
	key := fmt.Sprintf("change_email:%s", userID)

	val, err := s.otpService.GetRawOTP(ctx, key)
	if err != nil {
		return errors.New("mã OTP không hợp lệ hoặc đã hết hạn")
	}

	parts := strings.Split(val, ":")
	if len(parts) != 2 {
		return errors.New("dữ liệu OTP không hợp lệ")
	}

	newEmail := parts[0]
	storedOTP := parts[1]

	if req.OTP != storedOTP {
		return errors.New("mã OTP không chính xác")
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("người dùng không tồn tại")
	}

	update := bson.M{
		"$set": bson.M{
			"email":      newEmail,
			"isVerified": true,
			"updatedAt":  time.Now(),
		},
	}

	if err := s.repo.UpdateByID(ctx, user.ID.Hex(), update); err != nil {
		return err
	}

	// Xóa OTP sau khi dùng
	s.otpService.DeleteOTP(ctx, key)

	return nil
}

func (s *service) SendFriendRequest(ctx context.Context, fromID, toID primitive.ObjectID) error {
	if fromID == toID {
		return errors.New("không thể gửi yêu cầu cho chính mình")
	}
	return s.repo.SendFriendRequest(ctx, fromID, toID)
}

func (s *service) AcceptFriendRequest(ctx context.Context, receiverID, senderID primitive.ObjectID) error {
	// 1. Lấy thông tin người nhận (receiver)
	receiver, err := s.repo.FindByID(ctx, receiverID.Hex())
	if err != nil {
		return fmt.Errorf("Không tìm thấy người nhận lời mời")
	}

	// 2. Kiểm tra xem sender có nằm trong friendRequests của receiver không
	found := false
	for _, id := range receiver.FriendRequests {
		if id == senderID {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("Không có lời mời kết bạn từ người này")
	}

	// 3. Cập nhật cả 2 user: thêm vào friends, xóa khỏi friendRequests
	// Sử dụng repository để cập nhật
	updateReceiver := bson.M{
		"$pull":     bson.M{"friendRequests": senderID},
		"$addToSet": bson.M{"friends": senderID},
	}
	if err := s.repo.UpdateByID(ctx, receiverID.Hex(), updateReceiver); err != nil {
		return fmt.Errorf("Lỗi khi cập nhật người nhận")
	}

	updateSender := bson.M{
		"$addToSet": bson.M{"friends": receiverID},
	}
	if err := s.repo.UpdateByID(ctx, senderID.Hex(), updateSender); err != nil {
		return fmt.Errorf("Lỗi khi cập nhật người gửi")
	}

	return nil
}

func (s *service) BlockUser(ctx context.Context, userID, targetID primitive.ObjectID) error {
	if userID == targetID {
		return errors.New("không thể chặn chính mình")
	}
	return s.repo.BlockUser(ctx, userID, targetID)
}

func (s *service) ToggleHideProfile(ctx context.Context, userID primitive.ObjectID, hide bool) error {
	return s.repo.ToggleHideProfile(ctx, userID, hide)
}

func (s *service) CancelFriendRequest(ctx context.Context, fromID, toID primitive.ObjectID) error {
	return s.repo.CancelFriendRequest(ctx, fromID, toID)
}

func (s *service) FriendRequestExists(ctx context.Context, fromID, toID primitive.ObjectID) (bool, error) {
	return s.repo.FriendRequestExists(ctx, fromID, toID)
}
