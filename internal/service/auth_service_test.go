package service_test

import (
	"testing"

	"github.com/GydeonZ/task-api/internal/domain"
	"github.com/GydeonZ/task-api/internal/service"
	"github.com/GydeonZ/task-api/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewAuthService(mockRepo, "testsecret")

	req := &domain.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	resp, err := svc.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Test User", resp.Name)
	assert.Equal(t, "test@example.com", resp.Email)
	mockRepo.AssertExpectations(t)
}

func TestRegister_EmailConflict(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewAuthService(mockRepo, "testsecret")

	req := &domain.RegisterRequest{
		Name:     "Test User",
		Email:    "existing@example.com",
		Password: "password123",
	}

	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(apperrors.ErrConflict)

	resp, err := svc.Register(req)

	assert.Nil(t, resp)
	assert.Equal(t, apperrors.ErrConflict, err)
	mockRepo.AssertExpectations(t)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewAuthService(mockRepo, "testsecret")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &domain.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}
	user.ID = 1

	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)

	token, err := svc.Login(req)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewAuthService(mockRepo, "testsecret")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	user := &domain.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: string(hashedPassword),
	}
	user.ID = 1

	req := &domain.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	mockRepo.On("FindByEmail", "test@example.com").Return(user, nil)

	token, err := svc.Login(req)

	assert.Empty(t, token)
	assert.Equal(t, apperrors.ErrUnauthorized, err)
	mockRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := service.NewAuthService(mockRepo, "testsecret")

	req := &domain.LoginRequest{
		Email:    "notfound@example.com",
		Password: "password123",
	}

	mockRepo.On("FindByEmail", "notfound@example.com").Return(nil, apperrors.ErrNotFound)

	token, err := svc.Login(req)

	assert.Empty(t, token)
	assert.Equal(t, apperrors.ErrUnauthorized, err)
	mockRepo.AssertExpectations(t)
}
