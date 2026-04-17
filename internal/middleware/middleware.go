package middleware

import (
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang-jwt/jwt/v5"
)

// SetupMiddleware initializes all middleware for the application
func SetupMiddleware(app *fiber.App) {
	// Recover from panics
	app.Use(recover.New())

	// Request logging
	app.Use(logger.New(logger.Config{
		Format: "[${ip}] ${status} - ${method} ${path}\n",
	}))

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
}

// ErrorHandler is a custom error handling middleware
func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := err.Error()

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	log.Printf("Error: %v", err)

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"error":   message,
	})
}

// AuthMiddleware validates JWT Token
	func AuthMiddleware(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": "Missing authorization header",
			})
		}

		// Format: "Bearer <token>"
		parts := strings.Split(authHeader, " ") 
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": "Invalid authorization header format",
			})
		}
		
		tokenString := parts[1]

		// Verify token
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": "Invalid or expired token",
			})
		}

		// Save user information in context for later use
		userID := claims["user_id"]
		c.Locals("user_id", userID)
		c.Locals("email", claims["email"])
		return c.Next()
	}
