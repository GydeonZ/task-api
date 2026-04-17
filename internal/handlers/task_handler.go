package handlers

import (
	"strconv"

	"github.com/GydeonZ/task-api/internal/models"
	"github.com/GydeonZ/task-api/internal/service"
	"github.com/gofiber/fiber/v2"
)

// TaskHandler handles task-related HTTP requests
type TaskHandler struct {
	service *service.TaskService
}

// NewTaskHandler creates a new TaskHandler
func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// CreateTask handles POST /lists/:listId/tasks
func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	listID, err := strconv.Atoi(c.Params("listId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid list ID",
		})
	}

	var req models.CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	task, err := h.service.CreateTask(c.Context(), listID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "Task created successfully",
		Data:    task,
	})
}

// GetTask handles GET /tasks/:id
func (h *TaskHandler) GetTask(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid task ID",
		})
	}

	task, err := h.service.GetTask(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    task,
	})
}

// GetListTasks handles GET /lists/:listId/tasks
func (h *TaskHandler) GetListTasks(c *fiber.Ctx) error {
	listID, err := strconv.Atoi(c.Params("listId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid list ID",
		})
	}

	tasks, err := h.service.GetListTasks(c.Context(), listID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    tasks,
	})
}

// UpdateTask handles PUT /tasks/:id
func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid task ID",
		})
	}

	var req models.UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	task, err := h.service.UpdateTask(c.Context(), id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Task updated successfully",
		Data:    task,
	})
}

// DeleteTask handles DELETE /tasks/:id
func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid task ID",
		})
	}

	err = h.service.DeleteTask(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Task deleted successfully",
	})
}

// CompleteTask handles PATCH /tasks/:id/complete
func (h *TaskHandler) CompleteTask(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid task ID",
		})
	}

	task, err := h.service.CompleteTask(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Task marked as completed",
		Data:    task,
	})
}
