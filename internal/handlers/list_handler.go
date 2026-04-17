package handlers

import (
	"strconv"

	"github.com/GydeonZ/task-api/internal/models"
	"github.com/GydeonZ/task-api/internal/service"
	"github.com/gofiber/fiber/v2"
)

// ListHandler handles list-related HTTP requests
type ListHandler struct {
	service *service.ListService
}

// NewListHandler creates a new ListHandler
func NewListHandler(service *service.ListService) *ListHandler {
	return &ListHandler{service: service}
}

// CreateList handles POST /boards/:boardId/lists
func (h *ListHandler) CreateList(c *fiber.Ctx) error {
	boardID, err := strconv.Atoi(c.Params("boardId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid board ID",
		})
	}

	var req models.CreateListRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	list, err := h.service.CreateList(c.Context(), boardID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "List created successfully",
		Data:    list,
	})
}

// GetList handles GET /lists/:id
func (h *ListHandler) GetList(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid list ID",
		})
	}

	list, err := h.service.GetList(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    list,
	})
}

// GetBoardLists handles GET /boards/:boardId/lists
func (h *ListHandler) GetBoardLists(c *fiber.Ctx) error {
	boardID, err := strconv.Atoi(c.Params("boardId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid board ID",
		})
	}

	lists, err := h.service.GetBoardLists(c.Context(), boardID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    lists,
	})
}

// UpdateList handles PUT /lists/:id
func (h *ListHandler) UpdateList(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid list ID",
		})
	}

	var req models.UpdateListRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	list, err := h.service.UpdateList(c.Context(), id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "List updated successfully",
		Data:    list,
	})
}

// DeleteList handles DELETE /lists/:id
func (h *ListHandler) DeleteList(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid list ID",
		})
	}

	err = h.service.DeleteList(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "List deleted successfully",
	})
}
