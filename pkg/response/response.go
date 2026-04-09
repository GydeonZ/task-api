package response

import (
	"github.com/GydeonZ/task-api/pkg/apperrors"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, err error) {
	status := apperrors.HTTPStatus(err)
	c.JSON(status, Response{
		Success: false,
		Message: err.Error(),
	})
}
