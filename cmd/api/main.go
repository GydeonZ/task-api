package main

import (
	"log"

	"github.com/GydeonZ/task-api/internal/config"
	"github.com/GydeonZ/task-api/internal/domain"
	"github.com/GydeonZ/task-api/internal/handler"
	"github.com/GydeonZ/task-api/internal/repository"
	"github.com/GydeonZ/task-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()
	gin.SetMode(cfg.GinMode)

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&domain.User{}, &domain.Board{}, &domain.List{}, &domain.Card{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	boardRepo := repository.NewBoardRepository(db)
	listRepo := repository.NewListRepository(db)
	cardRepo := repository.NewCardRepository(db)

	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	boardSvc := service.NewBoardService(boardRepo)
	listSvc := service.NewListService(listRepo, boardRepo)
	cardSvc := service.NewCardService(cardRepo, listRepo, boardRepo)

	r := handler.SetupRouter(authSvc, boardSvc, listSvc, cardSvc, cfg.JWTSecret)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
