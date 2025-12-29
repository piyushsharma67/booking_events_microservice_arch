package middleware

import (
	"fmt"
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
		fmt.Println("roleAny:", roleAny, "exists:", exists)
		fmt.Println("allowedRoles:", allowedRoles)

		if !exists || roleAny == "" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		role := strings.TrimSpace(roleAny.(string))
		if role == allowedRoles {
			fmt.Println("Role matches! proceeding...")
			c.Next()
			return
		}

		fmt.Println("Role mismatch! access denied")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "access denied",
		})
	}
}

func SetRoleAndIdFromHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetHeader("X-User-Role")
		fmt.Println("role parsed in event service is", role)
		if role != "" {
			c.Set("role", role)
		}
		organiserId := c.GetHeader("X-User-ID")
		fmt.Println("useriD parsed in event service is", organiserId)
		if organiserId != "" {
			c.Set("user_id", organiserId)
		}
		c.Next()
	}
}
