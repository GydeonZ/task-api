package handlers

import (
	"strconv"

	"github.com/GydeonZ/task-api/internal/models"
	"github.com/GydeonZ/task-api/internal/service"
	"github.com/gofiber/fiber/v2"
)

// BoardHandler handles board-related HTTP requests
type BoardHandler struct {
	service *service.BoardService
}

// NewBoardHandler creates a new BoardHandler
func NewBoardHandler(service *service.BoardService) *BoardHandler {
	return &BoardHandler{service: service}
}

// CreateBoard handles POST /boards
func (h *BoardHandler) CreateBoard(c *fiber.Ctx) error {
	var req models.CreateBoardRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	//Extract userID from JWT token
	userID := c.Locals("user_id").(float64)

	board, err := h.service.CreateBoard(c.Context(), int(userID), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "Board created successfully",
		Data:    board,
	})
}

// GetBoard handles GET /boards/:id
func (h *BoardHandler) GetBoard(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid board ID",
		})
	}

	board, err := h.service.GetBoard(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    board,
	})
}

// GetUserBoards handles GET /boards
func (h *BoardHandler) GetUserBoards(c *fiber.Ctx) error {
	// Extract userID from JWT token
	userID := c.Locals("user_id").(float64)

	boards, err := h.service.GetUserBoards(c.Context(), int(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Data:    boards,
	})
}

// UpdateBoard handles PUT /boards/:id
func (h *BoardHandler) UpdateBoard(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid board ID",
		})
	}

	var req models.UpdateBoardRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	board, err := h.service.UpdateBoard(c.Context(), id, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Board updated successfully",
		Data:    board,
	})
}

// DeleteBoard handles DELETE /boards/:id
func (h *BoardHandler) DeleteBoard(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid board ID",
		})
	}

	err = h.service.DeleteBoard(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(models.APIResponse{
		Success: true,
		Message: "Board deleted successfully",
	})
}
