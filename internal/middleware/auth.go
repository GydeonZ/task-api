package middleware

import (
	"strings"

	"github.com/GydeonZ/task-api/pkg/apperrors"
	"github.com/GydeonZ/task-api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Error(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apperrors.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			response.Error(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.Error(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			response.Error(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set("userID", uint(userIDFloat))
		c.Next()
	}
}
