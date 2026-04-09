package handler

import (
	"github.com/GydeonZ/task-api/internal/middleware"
	"github.com/GydeonZ/task-api/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authSvc service.AuthService,
	boardSvc service.BoardService,
	listSvc service.ListService,
	cardSvc service.CardService,
	jwtSecret string,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			authHandler := NewAuthHandler(authSvc)
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtSecret))
		{
			boardHandler := NewBoardHandler(boardSvc)
			boards := protected.Group("/boards")
			{
				boards.GET("", boardHandler.GetBoards)
				boards.POST("", boardHandler.CreateBoard)
				boards.GET("/:id", boardHandler.GetBoard)
				boards.PUT("/:id", boardHandler.UpdateBoard)
				boards.DELETE("/:id", boardHandler.DeleteBoard)

				listHandler := NewListHandler(listSvc)
				boardLists := boards.Group("/:boardID/lists")
				{
					boardLists.GET("", listHandler.GetLists)
					boardLists.POST("", listHandler.CreateList)
					boardLists.PUT("/:id", listHandler.UpdateList)
					boardLists.DELETE("/:id", listHandler.DeleteList)
				}

				cardHandler := NewCardHandler(cardSvc)
				boardCards := boards.Group("/:boardID/lists/:listID/cards")
				{
					boardCards.GET("", cardHandler.GetCards)
					boardCards.POST("", cardHandler.CreateCard)
					boardCards.GET("/:id", cardHandler.GetCard)
					boardCards.PUT("/:id", cardHandler.UpdateCard)
					boardCards.DELETE("/:id", cardHandler.DeleteCard)
				}
			}

			cardHandler := NewCardHandler(cardSvc)
			protected.PATCH("/cards/:id/move", cardHandler.MoveCard)
		}
	}

	return r
}
