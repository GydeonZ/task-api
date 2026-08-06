package routes

import (
	"github.com/GydeonZ/task-api/internal/handlers"
	"github.com/GydeonZ/task-api/internal/middleware"
	"github.com/GydeonZ/task-api/internal/repository"
	"github.com/GydeonZ/task-api/internal/service"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes initializes all routes for the application
func SetupRoutes(app *fiber.App) {
	// Initialize repositories
	userRepo := repository.NewUserRepository()
	boardRepo := repository.NewBoardRepository()
	listRepo := repository.NewListRepository()
	taskRepo := repository.NewTaskRepository()

	// Initialize services
	userService := service.NewUserService(userRepo)
	boardService := service.NewBoardService(boardRepo)
	listService := service.NewListService(listRepo)
	taskService := service.NewTaskService(taskRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	boardHandler := handlers.NewBoardHandler(boardService)
	listHandler := handlers.NewListHandler(listService)
	taskHandler := handlers.NewTaskHandler(taskService)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Serve static OpenAPI spec
	app.Static("/openapi.yaml", "./docs/openapi.yaml")

	// Serve Scalar docs UI
	app.Get("/docs", func(c *fiber.Ctx) error {
		html := `<!doctype html>
<html>
  <head>
    <title>Task API Docs</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <script id="api-reference" data-url="/openapi.yaml"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(html)
	})

	// API v1 routes (versioning for future-proofing & for public API)
	api := app.Group("/api/v1")

	// =========== Public Routes ===========

	// Register routes
	register := api.Group("/register")
	register.Post("", userHandler.Register)

	// Login routes
	login := api.Group("/login")
	login.Post("", userHandler.Login)

	// =========== Protected Routes ===========

	// Protected routes with authentication middleware
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware)

	// Board routes
	// boards := api.Group("/boards")C
	boards := protected.Group("/boards")
	boards.Post("", boardHandler.CreateBoard)
	boards.Get("", boardHandler.GetUserBoards)
	boards.Get("/:id", boardHandler.GetBoard)
	boards.Put("/:id", boardHandler.UpdateBoard)
	boards.Delete("/:id", boardHandler.DeleteBoard)

	// List routes
	lists := protected.Group("/lists")
	lists.Get("/:id", listHandler.GetList)
	lists.Put("/:id", listHandler.UpdateList)
	lists.Delete("/:id", listHandler.DeleteList)

	boardLists := protected.Group("/boards/:boardId/lists")
	boardLists.Post("", listHandler.CreateList)
	boardLists.Get("", listHandler.GetBoardLists)

	// Task routes
	tasks := protected.Group("/tasks")
	tasks.Get("/:id", taskHandler.GetTask)
	tasks.Put("/:id", taskHandler.UpdateTask)
	tasks.Delete("/:id", taskHandler.DeleteTask)
	tasks.Patch("/:id/complete", taskHandler.CompleteTask)

	listTasks := protected.Group("/lists/:listId/tasks")
	listTasks.Post("", taskHandler.CreateTask)
	listTasks.Get("", taskHandler.GetListTasks)

	// 404 handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Route not found",
		})
	})
}
