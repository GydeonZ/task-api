package service

import (
	"time"

	"github.com/GydeonZ/task-api/internal/domain"
	"github.com/GydeonZ/task-api/internal/repository"
	"github.com/GydeonZ/task-api/pkg/apperrors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req *domain.RegisterRequest) (*domain.UserResponse, error)
	Login(req *domain.LoginRequest) (string, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *authService) Register(req *domain.RegisterRequest) (*domain.UserResponse, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *authService) Login(req *domain.LoginRequest) (string, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return "", apperrors.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", apperrors.ErrUnauthorized
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
