package repository

import (
	"context"
	"time"

	"github.com/GydeonZ/task-api/internal/database"
	"github.com/GydeonZ/task-api/internal/models"
	"github.com/GydeonZ/task-api/pkg/errors"
)

// TaskRepository handles task data operations
type TaskRepository struct{}

// NewTaskRepository creates a new TaskRepository
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

// Create creates a new task
func (r *TaskRepository) Create(ctx context.Context, task *models.Task) (*models.Task, error) {
	query := `
		INSERT INTO tasks (list_id, title, description, position, is_completed, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, list_id, title, description, position, is_completed, due_date, created_at, updated_at
	`

	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	if task.Position == 0 {
		task.Position = 1
	}

	err := database.DB.QueryRowContext(ctx, query,
		task.ListID,
		task.Title,
		task.Description,
		task.Position,
		task.IsCompleted,
		task.DueDate,
		task.CreatedAt,
		task.UpdatedAt,
	).Scan(&task.ID, &task.ListID, &task.Title, &task.Description, &task.Position,
		&task.IsCompleted, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return nil, errors.NewWithErr("DB_ERROR", "Failed to create task", 500, err)
	}

	return task, nil
}

// GetByID retrieves a task by ID
func (r *TaskRepository) GetByID(ctx context.Context, id int) (*models.Task, error) {
	query := `
		SELECT id, list_id, title, description, position, is_completed, due_date, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	task := &models.Task{}
	err := database.DB.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.ListID,
		&task.Title,
		&task.Description,
		&task.Position,
		&task.IsCompleted,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err != nil {
		return nil, errors.ErrTaskNotFound
	}

	return task, nil
}

// GetByListID retrieves all tasks in a list
func (r *TaskRepository) GetByListID(ctx context.Context, listID int) ([]models.Task, error) {
	query := `
		SELECT id, list_id, title, description, position, is_completed, due_date, created_at, updated_at
		FROM tasks
		WHERE list_id = $1
		ORDER BY position ASC, created_at DESC
	`

	rows, err := database.DB.QueryContext(ctx, query, listID)
	if err != nil {
		return nil, errors.NewWithErr("DB_ERROR", "Failed to retrieve tasks", 500, err)
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		task := models.Task{}
		err := rows.Scan(&task.ID, &task.ListID, &task.Title, &task.Description, &task.Position,
			&task.IsCompleted, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, errors.NewWithErr("DB_ERROR", "Failed to scan task", 500, err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Update updates a task
func (r *TaskRepository) Update(ctx context.Context, task *models.Task) (*models.Task, error) {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, position = $3, is_completed = $4, due_date = $5, updated_at = $6
		WHERE id = $7
		RETURNING id, list_id, title, description, position, is_completed, due_date, created_at, updated_at
	`

	task.UpdatedAt = time.Now()

	err := database.DB.QueryRowContext(ctx, query,
		task.Title,
		task.Description,
		task.Position,
		task.IsCompleted,
		task.DueDate,
		task.UpdatedAt,
		task.ID,
	).Scan(&task.ID, &task.ListID, &task.Title, &task.Description, &task.Position,
		&task.IsCompleted, &task.DueDate, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return nil, errors.NewWithErr("DB_ERROR", "Failed to update task", 500, err)
	}

	return task, nil
}

// Delete deletes a task
func (r *TaskRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id = $1`

	result, err := database.DB.ExecContext(ctx, query, id)
	if err != nil {
		return errors.NewWithErr("DB_ERROR", "Failed to delete task", 500, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewWithErr("DB_ERROR", "Failed to check rows affected", 500, err)
	}

	if rowsAffected == 0 {
		return errors.ErrTaskNotFound
	}

	return nil
}
