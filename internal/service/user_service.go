package service

import (
"context"

"github.com/google/uuid"
"golang.org/x/crypto/bcrypt"

"github.com/GydeonZ/task-api/internal/config"
"github.com/GydeonZ/task-api/internal/domain"
"github.com/GydeonZ/task-api/internal/repository"
"github.com/GydeonZ/task-api/pkg/apperror"
jwtutil "github.com/GydeonZ/task-api/pkg/jwt"
)

// AuthTokens holds the access and refresh tokens returned after authentication.
type AuthTokens struct {
AccessToken  string `json:"access_token"`
RefreshToken string `json:"refresh_token"`
}

// UserService defines business-logic operations for users.
type UserService interface {
Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error)
Login(ctx context.Context, req *domain.LoginRequest) (*AuthTokens, error)
GetProfile(ctx context.Context, id uuid.UUID) (*domain.User, error)
UpdateProfile(ctx context.Context, id uuid.UUID, req *domain.UpdateProfileRequest) (*domain.User, error)
}

type userService struct {
userRepo   repository.UserRepository
jwtManager *jwtutil.Manager
cfg        *config.Config
}

// NewUserService creates a new UserService.
func NewUserService(userRepo repository.UserRepository, jwtManager *jwtutil.Manager, cfg *config.Config) UserService {
return &userService{
userRepo:   userRepo,
jwtManager: jwtManager,
cfg:        cfg,
}
}

func (s *userService) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
if existing != nil {
return nil, apperror.Conflict("email already registered")
}

hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
if err != nil {
return nil, apperror.Internal(err)
}

user := &domain.User{
Name:     req.Name,
Email:    req.Email,
Password: string(hashed),
}

if err := s.userRepo.Create(ctx, user); err != nil {
return nil, err
}
return user, nil
}

func (s *userService) Login(ctx context.Context, req *domain.LoginRequest) (*AuthTokens, error) {
user, err := s.userRepo.FindByEmail(ctx, req.Email)
if err != nil {
return nil, apperror.Unauthorized("invalid credentials")
}

if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
return nil, apperror.Unauthorized("invalid credentials")
}

accessToken, err := s.jwtManager.Generate(user.ID, jwtutil.AccessToken, s.cfg.JWT.AccessTokenTTL)
if err != nil {
return nil, apperror.Internal(err)
}

refreshToken, err := s.jwtManager.Generate(user.ID, jwtutil.RefreshToken, s.cfg.JWT.RefreshTokenTTL)
if err != nil {
return nil, apperror.Internal(err)
}

return &AuthTokens{
AccessToken:  accessToken,
RefreshToken: refreshToken,
}, nil
}

func (s *userService) GetProfile(ctx context.Context, id uuid.UUID) (*domain.User, error) {
return s.userRepo.FindByID(ctx, id)
}

func (s *userService) UpdateProfile(ctx context.Context, id uuid.UUID, req *domain.UpdateProfileRequest) (*domain.User, error) {
user, err := s.userRepo.FindByID(ctx, id)
if err != nil {
return nil, err
}

if req.Name != "" {
user.Name = req.Name
}
if req.Avatar != "" {
user.Avatar = req.Avatar
}

if err := s.userRepo.Update(ctx, user); err != nil {
return nil, err
}
return user, nil
}
