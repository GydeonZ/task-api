package main

import (
	"log"

	"github.com/GydeonZ/task-api/internal/config"
	"github.com/GydeonZ/task-api/internal/database"
	"github.com/GydeonZ/task-api/internal/middleware"
	"github.com/GydeonZ/task-api/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	err = database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Task API v1.0",
	})

	// Setup middleware
	middleware.SetupMiddleware(app)

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Server starting on %s", addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
