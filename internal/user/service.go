package user

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	// "fmt"
	"socialnetwork/dto/request"
	"socialnetwork/pkg/auth"
	"socialnetwork/models"
)

type Service interface {
	Register(ctx context.Context, user *models.User) error
	Login(ctx context.Context, email, password string) (string, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateProfile(ctx context.Context, id string, req *request.UpdateProfileRequest) error
	ChangePassword(ctx context.Context, userID string, req *request.ChangePasswordRequest) error

}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
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


