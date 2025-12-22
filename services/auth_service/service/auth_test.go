package service

import (
	"context"
	"testing"

	"github.com/piyushsharma67/movie_booking/services/auth_service/databases"
	"github.com/piyushsharma67/movie_booking/services/auth_service/models"
	"github.com/piyushsharma67/movie_booking/services/auth_service/repository"
	"github.com/stretchr/testify/assert"
)

func TestAuthService(t *testing.T) {
	sqliteDB, err := databases.InitSqliteTestDB()
	assert.NoError(t, err)
	defer sqliteDB.Close()

	db := databases.NewSqliteDB(sqliteDB)
	repo := repository.NewUserRepository(db)
	svc := NewAuthService(repo)

	user := models.User{
		Name:     "Piyush",
		Email:    "piyush@test.com",
		Password: "password123",
	}

	created, err := svc.SignUp(context.Background(), user)

	assert.NoError(t, err)
	// assert.NotEmpty(t, created.ID)
	assert.Equal(t, "user", created.Role)
	assert.NotEqual(t, "password123", created.PasswordHash)
}
