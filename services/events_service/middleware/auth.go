package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/piyushsharma67/events_booking/services/events_service/utils"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateJWT(token, os.Getenv("JWT_SECRET"))
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func RoleAuthMiddleware(allowedRoles string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleAny, exists := c.Get("role")

		if !exists || roleAny == "" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		role := strings.TrimSpace(roleAny.(string))
		if role == allowedRoles {

			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "access denied",
		})
	}
}

func SetRoleAndIdFromHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("X-User-Role")

		if role != "" {
			c.Set("role", role)
		}
		organiserId := c.GetHeader("X-User-ID")

		if organiserId != "" {
			c.Set("user_id", organiserId)
		}
		c.Next()
	}
}
