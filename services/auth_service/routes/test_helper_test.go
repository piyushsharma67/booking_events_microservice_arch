package routes

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/piyushsharma67/movie_booking/services/auth_service/databases"
	"github.com/piyushsharma67/movie_booking/services/auth_service/utils"
)

func setupTestRouter(t *testing.T) (*gin.Engine, databases.Database) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	db := databases.GetTestDB()

	db, err := databases.InitSqliteTestDB()
	if err != nil {
		t.Fatalf("failed to init sqlite db: %v", err)
	}

	r := InitialiseRoutes(db)
	return r, db
}

func seedUser(t *testing.T, db databases.Database, email, password, role string) {
	t.Helper()

	hashed, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	err = db.InsertUser(context.Background(), &databases.User{
		Name:         "Test User",
		Email:        email,
		PasswordHash: hashed,
		Role:         role,
	})
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}
}
