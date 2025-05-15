package user

import (
	"context"
	"errors"
	// "fmt"
	"socialnetwork/pkg/auth"
	"socialnetwork/models"
)

type Service interface {
	Register(ctx context.Context, user *models.User) error
	Login(ctx context.Context, email, password string) (string, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetAllUsers() ([]*models.User, error)
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

func (s *service) GetAllUsers() ([]*models.User, error) {
	return s.repo.GetAllUsers()
}


