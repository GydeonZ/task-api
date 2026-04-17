package repository

import (
	"context"
	"time"

	"github.com/GydeonZ/task-api/internal/database"
	"github.com/GydeonZ/task-api/internal/models"
	"github.com/GydeonZ/task-api/pkg/errors"
)

// BoardRepository handles board data operations
type BoardRepository struct{}

// NewBoardRepository creates a new BoardRepository
func NewBoardRepository() *BoardRepository {
	return &BoardRepository{}
}

// Create creates a new board
func (r *BoardRepository) Create(ctx context.Context, board *models.Board) (*models.Board, error) {
	query := `
		INSERT INTO boards (user_id, title, slug, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, title, slug, created_at, updated_at
	`

	board.CreatedAt = time.Now()
	board.UpdatedAt = time.Now()

	err := database.DB.QueryRowContext(ctx, query,
		board.UserID,
		board.Title,
		board.Slug,
		board.CreatedAt,
		board.UpdatedAt,
	).Scan(&board.ID, &board.UserID, &board.Title, &board.Slug, &board.CreatedAt, &board.UpdatedAt)

	if err != nil {
		return nil, errors.NewWithErr("DB_ERROR", "Failed to create board", 500, err)
	}

	return board, nil
}

// GetByID retrieves a board by ID
func (r *BoardRepository) GetByID(ctx context.Context, id int) (*models.Board, error) {
	query := `
		SELECT id, user_id, title, slug, created_at, updated_at
		FROM boards
		WHERE id = $1
	`

	board := &models.Board{}
	err := database.DB.QueryRowContext(ctx, query, id).Scan(
		&board.ID,
		&board.UserID,
		&board.Title,
		&board.Slug,
		&board.CreatedAt,
		&board.UpdatedAt,
	)

	if err != nil {
		return nil, errors.NewWithErr("DB_ERROR", "Board not found", 404, err)
	}

	return board, nil
}

// GetByUserID retrieves all boards for a user
func (r *BoardRepository) GetByUserID(ctx context.Context, userID int) ([]models.Board, error) {
	query := `
		SELECT id, user_id, title, slug, created_at, updated_at
		FROM boards
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := database.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, errors.NewWithErr("DB_ERROR", "Failed to retrieve boards", 500, err)
	}
	defer rows.Close()

	boards := []models.Board{}
	for rows.Next() {
		board := models.Board{}
		err := rows.Scan(&board.ID, &board.UserID, &board.Title, &board.Slug, &board.CreatedAt, &board.UpdatedAt)
		if err != nil {
			return nil, errors.NewWithErr("DB_ERROR", "Failed to scan board", 500, err)
		}
		boards = append(boards, board)
	}

	return boards, nil
}

// Update updates a board
func (r *BoardRepository) Update(ctx context.Context, board *models.Board) (*models.Board, error) {
	query := `
		UPDATE boards
		SET title = $1, slug = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, user_id, title, slug, created_at, updated_at
	`

	board.UpdatedAt = time.Now()

	err := database.DB.QueryRowContext(ctx, query,
		board.Title,
		board.Slug,
		board.UpdatedAt,
		board.ID,
	).Scan(&board.ID, &board.UserID, &board.Title, &board.Slug, &board.CreatedAt, &board.UpdatedAt)

	if err != nil {
		return nil, errors.NewWithErr("DB_ERROR", "Failed to update board", 500, err)
	}

	return board, nil
}

// Delete deletes a board
func (r *BoardRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM boards WHERE id = $1`

	result, err := database.DB.ExecContext(ctx, query, id)
	if err != nil {
		return errors.NewWithErr("DB_ERROR", "Failed to delete board", 500, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewWithErr("DB_ERROR", "Failed to check rows affected", 500, err)
	}

	if rowsAffected == 0 {
		return errors.ErrBoardNotFound
	}

	return nil
}
