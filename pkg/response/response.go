package response

import "github.com/gofiber/fiber/v2"

// Meta holds pagination metadata.
type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Response is the standard JSON envelope returned by every endpoint.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// OK sends a 200 response with data.
func OK(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{Success: true, Data: data})
}

// Created sends a 201 response with data.
func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{Success: true, Data: data})
}

// OKWithMeta sends a 200 response with data and pagination metadata.
func OKWithMeta(c *fiber.Ctx, data interface{}, meta Meta) error {
	return c.Status(fiber.StatusOK).JSON(Response{Success: true, Data: data, Meta: &meta})
}

// NoContent sends a 204 response.
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// Error sends a failure response with the given HTTP status code.
func Error(c *fiber.Ctx, code int, message string) error {
	return c.Status(code).JSON(Response{Success: false, Message: message})
}

// ValidationError sends a 422 response with field-level validation errors.
func ValidationError(c *fiber.Ctx, errs interface{}) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{
		Success: false,
		Message: "validation failed",
		Errors:  errs,
	})
}
