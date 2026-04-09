package service_test

import (
	"testing"

	"github.com/GydeonZ/task-api/internal/domain"
	"github.com/GydeonZ/task-api/internal/service"
	"github.com/GydeonZ/task-api/pkg/apperrors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockBoardRepository struct {
	mock.Mock
}

func (m *MockBoardRepository) Create(board *domain.Board) error {
	args := m.Called(board)
	return args.Error(0)
}

func (m *MockBoardRepository) FindByUserID(userID uint) ([]domain.Board, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.Board), args.Error(1)
}

func (m *MockBoardRepository) FindByID(id uint) (*domain.Board, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Board), args.Error(1)
}

func (m *MockBoardRepository) Update(board *domain.Board) error {
	args := m.Called(board)
	return args.Error(0)
}

func (m *MockBoardRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func newTestBoard(id, userID uint, title string) *domain.Board {
	b := &domain.Board{
		Title:  title,
		UserID: userID,
	}
	b.Model = gorm.Model{ID: id}
	return b
}

func TestCreateBoard_Success(t *testing.T) {
	mockRepo := new(MockBoardRepository)
	svc := service.NewBoardService(mockRepo)

	req := &domain.CreateBoardRequest{Title: "My Board", Description: "Test"}
	mockRepo.On("Create", mock.AnythingOfType("*domain.Board")).Return(nil)

	board, err := svc.CreateBoard(1, req)

	assert.NoError(t, err)
	assert.NotNil(t, board)
	assert.Equal(t, "My Board", board.Title)
	assert.Equal(t, uint(1), board.UserID)
	mockRepo.AssertExpectations(t)
}

func TestGetBoards_Success(t *testing.T) {
	mockRepo := new(MockBoardRepository)
	svc := service.NewBoardService(mockRepo)

	boards := []domain.Board{
		*newTestBoard(1, 1, "Board 1"),
		*newTestBoard(2, 1, "Board 2"),
	}
	mockRepo.On("FindByUserID", uint(1)).Return(boards, nil)

	result, err := svc.GetBoards(1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

func TestGetBoard_Success(t *testing.T) {
	mockRepo := new(MockBoardRepository)
	svc := service.NewBoardService(mockRepo)

	board := newTestBoard(1, 1, "My Board")
	mockRepo.On("FindByID", uint(1)).Return(board, nil)

	result, err := svc.GetBoard(1, 1)

	assert.NoError(t, err)
	assert.Equal(t, "My Board", result.Title)
	mockRepo.AssertExpectations(t)
}

func TestGetBoard_Unauthorized(t *testing.T) {
	mockRepo := new(MockBoardRepository)
	svc := service.NewBoardService(mockRepo)

	board := newTestBoard(1, 2, "Other User's Board")
	mockRepo.On("FindByID", uint(1)).Return(board, nil)

	result, err := svc.GetBoard(1, 1)

	assert.Nil(t, result)
	assert.Equal(t, apperrors.ErrUnauthorized, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateBoard_Success(t *testing.T) {
	mockRepo := new(MockBoardRepository)
	svc := service.NewBoardService(mockRepo)

	board := newTestBoard(1, 1, "Old Title")
	req := &domain.UpdateBoardRequest{Title: "New Title"}

	mockRepo.On("FindByID", uint(1)).Return(board, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Board")).Return(nil)

	result, err := svc.UpdateBoard(1, 1, req)

	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
	mockRepo.AssertExpectations(t)
}

func TestDeleteBoard_Success(t *testing.T) {
	mockRepo := new(MockBoardRepository)
	svc := service.NewBoardService(mockRepo)

	board := newTestBoard(1, 1, "My Board")
	mockRepo.On("FindByID", uint(1)).Return(board, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := svc.DeleteBoard(1, 1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
