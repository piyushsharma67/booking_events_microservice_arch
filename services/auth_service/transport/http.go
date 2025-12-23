package transport

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/endpoint"
	"github.com/piyushsharma67/movie_booking/services/auth_service/service"
	"github.com/piyushsharma67/movie_booking/services/auth_service/utils"
)

func GinHandler(e endpoint.Endpoint, newRequest func() interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request interface{}
		request = newRequest()

		if err := c.ShouldBindJSON(request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Call the Go Kit endpoint
		resp, err := e(c, request)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				c.JSON(http.StatusRequestTimeout, gin.H{"error": "request timed out"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}

func ValidateGinHandler(svc service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Status(http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ValidateJWT(token, os.Getenv("JWT_SECRET"))
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		// 🔥 NGINX reads these headers
		c.Header("X-User-Id", claims.UserID)
		c.Header("X-User-Role", claims.Role)

		c.Status(http.StatusOK)
	}
}
