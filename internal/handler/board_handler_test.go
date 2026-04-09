package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/GydeonZ/task-api/internal/domain"
	"github.com/GydeonZ/task-api/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockBoardService struct {
	mock.Mock
}

func (m *MockBoardService) CreateBoard(userID uint, req *domain.CreateBoardRequest) (*domain.Board, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Board), args.Error(1)
}

func (m *MockBoardService) GetBoards(userID uint) ([]domain.Board, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.Board), args.Error(1)
}

func (m *MockBoardService) GetBoard(userID, boardID uint) (*domain.Board, error) {
	args := m.Called(userID, boardID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Board), args.Error(1)
}

func (m *MockBoardService) UpdateBoard(userID, boardID uint, req *domain.UpdateBoardRequest) (*domain.Board, error) {
	args := m.Called(userID, boardID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Board), args.Error(1)
}

func (m *MockBoardService) DeleteBoard(userID, boardID uint) error {
	args := m.Called(userID, boardID)
	return args.Error(0)
}

func setupBoardRouter(svc *MockBoardService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handler.NewBoardHandler(svc)

	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	})

	r.GET("/boards", h.GetBoards)
	r.POST("/boards", h.CreateBoard)
	r.GET("/boards/:id", h.GetBoard)
	r.DELETE("/boards/:id", h.DeleteBoard)
	return r
}

func newBoardWithID(id, userID uint, title string) *domain.Board {
	b := &domain.Board{Title: title, UserID: userID}
	b.Model = gorm.Model{ID: id}
	return b
}

func TestGetBoards_Handler(t *testing.T) {
	mockSvc := new(MockBoardService)
	r := setupBoardRouter(mockSvc)

	boards := []domain.Board{*newBoardWithID(1, 1, "Board 1")}
	mockSvc.On("GetBoards", uint(1)).Return(boards, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/boards", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	mockSvc.AssertExpectations(t)
}

func TestCreateBoard_Handler(t *testing.T) {
	mockSvc := new(MockBoardService)
	r := setupBoardRouter(mockSvc)

	board := newBoardWithID(1, 1, "New Board")
	mockSvc.On("CreateBoard", uint(1), mock.AnythingOfType("*domain.CreateBoardRequest")).Return(board, nil)

	body, _ := json.Marshal(domain.CreateBoardRequest{Title: "New Board"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/boards", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	mockSvc.AssertExpectations(t)
}

func TestGetBoard_Handler(t *testing.T) {
	mockSvc := new(MockBoardService)
	r := setupBoardRouter(mockSvc)

	board := newBoardWithID(1, 1, "My Board")
	mockSvc.On("GetBoard", uint(1), uint(1)).Return(board, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/boards/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	mockSvc.AssertExpectations(t)
}

func TestDeleteBoard_Handler(t *testing.T) {
	mockSvc := new(MockBoardService)
	r := setupBoardRouter(mockSvc)

	mockSvc.On("DeleteBoard", uint(1), uint(1)).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/boards/1", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.True(t, resp["success"].(bool))
	mockSvc.AssertExpectations(t)
}
