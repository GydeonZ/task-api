package service

import (
	"context"

	"github.com/GydeonZ/task-api/internal/models"
	"github.com/GydeonZ/task-api/internal/repository"
	"github.com/GydeonZ/task-api/pkg/errors"
)

// TaskService handles task business logic
type TaskService struct {
	repo *repository.TaskRepository
}

// NewTaskService creates a new TaskService
func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// CreateTask creates a new task in a list
func (s *TaskService) CreateTask(ctx context.Context, listID int, req *models.CreateTaskRequest) (*models.Task, error) {
	if req.Title == "" {
		return nil, errors.ErrInvalidInput
	}

	task := &models.Task{
		ListID:      listID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		IsCompleted: false,
		Position:    1,
	}

	return s.repo.Create(ctx, task)
}

// GetTask retrieves a task by ID
func (s *TaskService) GetTask(ctx context.Context, id int) (*models.Task, error) {
	if id <= 0 {
		return nil, errors.ErrInvalidInput
	}

	return s.repo.GetByID(ctx, id)
}

// GetListTasks retrieves all tasks in a list
func (s *TaskService) GetListTasks(ctx context.Context, listID int) ([]models.Task, error) {
	if listID <= 0 {
		return nil, errors.ErrInvalidInput
	}

	return s.repo.GetByListID(ctx, listID)
}

// UpdateTask updates a task
func (s *TaskService) UpdateTask(ctx context.Context, id int, req *models.UpdateTaskRequest) (*models.Task, error) {
	if id <= 0 {
		return nil, errors.ErrInvalidInput
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	task.IsCompleted = req.IsCompleted
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if req.Position > 0 {
		task.Position = req.Position
	}

	return s.repo.Update(ctx, task)
}

// DeleteTask deletes a task
func (s *TaskService) DeleteTask(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.ErrInvalidInput
	}

	return s.repo.Delete(ctx, id)
}

// CompleteTask marks a task as completed
func (s *TaskService) CompleteTask(ctx context.Context, id int) (*models.Task, error) {
	if id <= 0 {
		return nil, errors.ErrInvalidInput
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	task.IsCompleted = true
	return s.repo.Update(ctx, task)
}
