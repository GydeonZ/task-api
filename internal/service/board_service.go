package service

import (
	"context"
	"strings"

	"github.com/GydeonZ/task-api/internal/models"
	"github.com/GydeonZ/task-api/internal/repository"
	"github.com/GydeonZ/task-api/pkg/errors"
)

// BoardService handles board business logic
type BoardService struct {
	repo *repository.BoardRepository
}

// NewBoardService creates a new BoardService
func NewBoardService(repo *repository.BoardRepository) *BoardService {
	return &BoardService{repo: repo}
}

// CreateBoard creates a new board
func (s *BoardService) CreateBoard(ctx context.Context, userID int, req *models.CreateBoardRequest) (*models.Board, error) {
	if req.Title == "" {
		return nil, errors.ErrInvalidInput
	}

	board := &models.Board{
		UserID: userID,
		Title:  req.Title,
		Slug:   slugify(req.Title),
	}

	return s.repo.Create(ctx, board)
}

// GetBoard retrieves a board by ID
func (s *BoardService) GetBoard(ctx context.Context, id int) (*models.Board, error) {
	if id <= 0 {
		return nil, errors.ErrInvalidInput
	}

	return s.repo.GetByID(ctx, id)
}

// GetUserBoards retrieves all boards for a user
func (s *BoardService) GetUserBoards(ctx context.Context, userID int) ([]models.Board, error) {
	if userID <= 0 {
		return nil, errors.ErrInvalidInput
	}

	return s.repo.GetByUserID(ctx, userID)
}

// UpdateBoard updates a board
func (s *BoardService) UpdateBoard(ctx context.Context, id int, req *models.UpdateBoardRequest) (*models.Board, error) {
	if id <= 0 || req.Title == "" {
		return nil, errors.ErrInvalidInput
	}

	board, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	board.Title = req.Title
	board.Slug = slugify(req.Title)

	return s.repo.Update(ctx, board)
}

// DeleteBoard deletes a board
func (s *BoardService) DeleteBoard(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.ErrInvalidInput
	}

	return s.repo.Delete(ctx, id)
}

// slugify converts a string to a URL-friendly slug
func slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, s)
	return s
}
