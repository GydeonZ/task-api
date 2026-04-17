package service

import (
	"context"

	"github.com/GydeonZ/task-api/internal/models"
	"github.com/GydeonZ/task-api/internal/repository"
	"github.com/GydeonZ/task-api/pkg/errors"
)

// ListService handles list business logic
type ListService struct {
	repo *repository.ListRepository
}

// NewListService creates a new ListService
func NewListService(repo *repository.ListRepository) *ListService {
	return &ListService{repo: repo}
}

// CreateList creates a new list in a board
func (s *ListService) CreateList(ctx context.Context, boardID int, req *models.CreateListRequest) (*models.List, error) {
	if req.Title == "" {
		return nil, errors.ErrInvalidInput
	}

	list := &models.List{
		BoardID:  boardID,
		Title:    req.Title,
		Position: 1,
	}

	return s.repo.Create(ctx, list)
}

// GetList retrieves a list by ID
func (s *ListService) GetList(ctx context.Context, id int) (*models.List, error) {
	if id <= 0 {
		return nil, errors.ErrInvalidInput
	}

	return s.repo.GetByID(ctx, id)
}

// GetBoardLists retrieves all lists in a board
func (s *ListService) GetBoardLists(ctx context.Context, boardID int) ([]models.List, error) {
	if boardID <= 0 {
		return nil, errors.ErrInvalidInput
	}

	return s.repo.GetByBoardID(ctx, boardID)
}

// UpdateList updates a list
func (s *ListService) UpdateList(ctx context.Context, id int, req *models.UpdateListRequest) (*models.List, error) {
	if id <= 0 {
		return nil, errors.ErrInvalidInput
	}

	list, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		list.Title = req.Title
	}
	if req.Position > 0 {
		list.Position = req.Position
	}

	return s.repo.Update(ctx, list)
}

// DeleteList deletes a list
func (s *ListService) DeleteList(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.ErrInvalidInput
	}

	return s.repo.Delete(ctx, id)
}
