package handlers

import (
	"github.com/GydeonZ/task-api/internal/models"
	"github.com/GydeonZ/task-api/internal/service"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service *service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// User LoginHandler
func (h *UserHandler) Login(c *fiber.Ctx) error {

	var req models.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	token, err := h.service.LoginUser(
		c.Context(),
		&req,
	)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

// Add RegistrationHandler
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req models.RegisterUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error:   "Invalid request body",
		})
	}

	register, err := h.service.RegisterUser(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Success: true,
		Message: "User Created Succesfully",
		Data: fiber.Map{
			"id":         register.ID,
			"username":   register.Username,
			"email":      register.Email,
			"created_at": register.CreatedAt,
			"updated_at": register.UpdatedAt,
		},
	})
}
