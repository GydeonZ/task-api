package service

import (
	"context"
	"os"
	"time"

	"github.com/GydeonZ/task-api/internal/models"
	"github.com/GydeonZ/task-api/internal/repository"
	"github.com/GydeonZ/task-api/pkg/errors"
	"github.com/GydeonZ/task-api/pkg/validator"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Register - registers a new user
func (s *UserService) RegisterUser(ctx context.Context, req *models.RegisterUserRequest) (*models.User, error) {
	// Validasi Username
	valid, msg := validator.ValidateUsername(req.Username)
	if !valid {
		return nil, &errors.AppError{
			Code:    "INVALID_USERNAME",
			Message: msg,
			Status:  400,
		}
	}

	// Validasi Email format
	if !validator.ValidateEmail(req.Email) {
		return nil, &errors.AppError{
			Code:    "INVALID_EMAIL",
			Message: "Email format is invalid. Please provide a valid email address.",
			Status:  400,
		}
	}

	// Validasi Password
	valid, msg = validator.ValidatePassword(req.Password)
	if !valid {
		return nil, &errors.AppError{
			Code:    "INVALID_PASSWORD",
			Message: msg,
			Status:  400,
		}
	}

	// Check username that has been registered
	_, err := s.repo.FindUserByUsername(ctx, req.Username)
	if err == nil {
		// username already exist
		return nil, &errors.AppError{
			Code:    "DUPLICATE_USERNAME",
			Message: "Username '" + req.Username + "' is already taken. Please choose a different one.",
			Status:  409,
		}
	}

	// Check email that has been registered
	_, err = s.repo.FindUserByEmail(ctx, req.Email)
	if err == nil {
		// Email Already Exists
		return nil, &errors.AppError{
			Code:    "DUPLICATE_EMAIL",
			Message: "Email '" + req.Email + "' is already registered. Please use a different email.",
			Status:  409,
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &errors.AppError{
			Code:    "HASHING_ERROR",
			Message: "Failed to hash password: " + err.Error(),
			Status:  500,
			Err:     err,
		}
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	createdUser, err := s.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

// LoginUser - authenticates user and returns JWT token
func (s *UserService) LoginUser(ctx context.Context, req *models.LoginRequest) (string, error) {
	// Cek email ada di database
	user, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return "", &errors.AppError{
			Code:    "USER_NOT_FOUND",
			Message: "Email '" + req.Email + "' not found. Please check your email or register for a new account.",
			Status:  404,
			Err:     err,
		}
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	)

	if err != nil {
		return "", &errors.AppError{
			Code:    "INVALID_CREDENTIALS",
			Message: "Invalid email or password",
			Status:  401,
			Err:     err,
		}
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	t, err := token.SignedString(
		[]byte(os.Getenv("JWT_SECRET")),
	)

	if err != nil {
		return "", &errors.AppError{
			Code:    "JWT_ERROR",
			Message: "Failed to generate token",
			Status:  500,
			Err:     err,
		}
	}

	return t, nil
}

// TODO: Implement CreateUser after repository method is added
