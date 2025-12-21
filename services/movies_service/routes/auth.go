package routes

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/piyushsharma67/movie_booking/services/movies_service/databases"
	"github.com/piyushsharma67/movie_booking/services/movies_service/models"
	"github.com/piyushsharma67/movie_booking/services/movies_service/utils"
	"golang.org/x/crypto/bcrypt"
)

func SignUpHandler(db databases.Database) gin.HandlerFunc {
	return func(g *gin.Context) {
		var request models.UserSignUp

		if err := g.ShouldBindJSON(&request); err != nil {
			g.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
			return
		}

		user := &databases.User{
			Name:         request.Username,
			Email:        request.Email,
			PasswordHash: string(hashed),
			Role:         "user",
		}

		if err := db.InsertUser(g.Request.Context(), user); err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		jwt, err := utils.GenerateJWT(user.ID, user.Email, user.Role, os.Getenv("JWT_SECRET"))
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
			return
		}

		g.JSON(http.StatusCreated, gin.H{"id": user.ID, "email": user.Email, "token": jwt})

	}
}

func LoginHandler(db databases.Database) gin.HandlerFunc {
	return func(g *gin.Context) {
		var request models.UserLogin

		if err := g.ShouldBindJSON(&request); err != nil {
			g.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
			return
		}
		user, err := db.GetUserByEmail(g.Request.Context(), request.Email)

		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
			g.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		jwt, err := utils.GenerateJWT(user.ID, user.Email, user.Role, os.Getenv("JWT_SECRET"))
		if err != nil {
			g.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
			return
		}
		g.JSON(http.StatusOK, gin.H{"id": user.ID, "email": user.Email, "token": jwt})

	}
}

func ValidateHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Status(http.StatusUnauthorized)
			return
		}

		// Expecting "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Status(http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		claims, err := utils.ValidateJWT(tokenStr, os.Getenv("JWT_SECRET"))
		if err != nil {
			c.Status(http.StatusUnauthorized)
			return
		}

		// Set headers for NGINX to pass to upstream
		c.Header("X-User-Id", claims.UserID)
		c.Header("X-User-Role", claims.Role)

		// Return 200 OK for NGINX
		c.Status(http.StatusOK)
	}
}
