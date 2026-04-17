package repository

import (
	"context"
	"time"

	"github.com/GydeonZ/task-api/internal/database"
	"github.com/GydeonZ/task-api/internal/models"
	"github.com/GydeonZ/task-api/pkg/errors"
)

// ListRepository handles list data operations
type ListRepository struct{}

// NewListRepository creates a new ListRepository
func NewListRepository() *ListRepository {
	return &ListRepository{}
}

// Create creates a new list
func (r *ListRepository) Create(ctx context.Context, list *models.List) (*models.List, error) {
	query := `
		INSERT INTO lists (board_id, title, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, board_id, title, position, created_at, updated_at
	`

	list.CreatedAt = time.Now()
	list.UpdatedAt = time.Now()

	err := database.DB.QueryRowContext(ctx, query,
		list.BoardID,
		list.Title,
		list.Position,
		list.CreatedAt,
		list.UpdatedAt,
	).Scan(&list.ID, &list.BoardID, &list.Title, &list.Position, &list.CreatedAt, &list.UpdatedAt)

	if err != nil {
		return nil, errors.NewWithErr("DB_ERROR", "Failed to create list", 500, err)
	}

	return list, nil
}

// GetByID retrieves a list by ID
func (r *ListRepository) GetByID(ctx context.Context, id int) (*models.List, error) {
	query := `
		SELECT id, board_id, title, position, created_at, updated_at
		FROM lists
		WHERE id = $1
	`

	list := &models.List{}
	err := database.DB.QueryRowContext(ctx, query, id).Scan(
		&list.ID,
		&list.BoardID,
		&list.Title,
		&list.Position,
		&list.CreatedAt,
		&list.UpdatedAt,
	)

	if err != nil {
		return nil, errors.ErrListNotFound
	}

	return list, nil
}

// GetByBoardID retrieves all lists in a board
func (r *ListRepository) GetByBoardID(ctx context.Context, boardID int) ([]models.List, error) {
	query := `
		SELECT id, board_id, title, position, created_at, updated_at
		FROM lists
		WHERE board_id = $1
		ORDER BY position ASC
	`

	rows, err := database.DB.QueryContext(ctx, query, boardID)
	if err != nil {
		return nil, errors.NewWithErr("DB_ERROR", "Failed to retrieve lists", 500, err)
	}
	defer rows.Close()

	lists := []models.List{}
	for rows.Next() {
		list := models.List{}
		err := rows.Scan(&list.ID, &list.BoardID, &list.Title, &list.Position, &list.CreatedAt, &list.UpdatedAt)
		if err != nil {
			return nil, errors.NewWithErr("DB_ERROR", "Failed to scan list", 500, err)
		}
		lists = append(lists, list)
	}

	return lists, nil
}

// Update updates a list
func (r *ListRepository) Update(ctx context.Context, list *models.List) (*models.List, error) {
	query := `
		UPDATE lists
		SET title = $1, position = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, board_id, title, position, created_at, updated_at
	`

	list.UpdatedAt = time.Now()

	err := database.DB.QueryRowContext(ctx, query,
		list.Title,
		list.Position,
		list.UpdatedAt,
		list.ID,
	).Scan(&list.ID, &list.BoardID, &list.Title, &list.Position, &list.CreatedAt, &list.UpdatedAt)

	if err != nil {
		return nil, errors.NewWithErr("DB_ERROR", "Failed to update list", 500, err)
	}

	return list, nil
}

// Delete deletes a list
func (r *ListRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM lists WHERE id = $1`

	result, err := database.DB.ExecContext(ctx, query, id)
	if err != nil {
		return errors.NewWithErr("DB_ERROR", "Failed to delete list", 500, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewWithErr("DB_ERROR", "Failed to check rows affected", 500, err)
	}

	if rowsAffected == 0 {
		return errors.ErrListNotFound
	}

	return nil
}
