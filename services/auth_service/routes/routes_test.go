package routes

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/piyushsharma67/movie_booking/services/auth_service/databases"
	"github.com/piyushsharma67/movie_booking/services/auth_service/repository"
	"github.com/piyushsharma67/movie_booking/services/auth_service/service"
	"github.com/stretchr/testify/assert"
)

func setupSharedTestServer() *gin.Engine {
	sqlDB := databases.InitSharedSqliteTestDB()

	db := databases.NewSqliteDB(sqlDB)
	repo := repository.NewUserRepository(db)
	svc := service.NewAuthService(repo)
	gin.SetMode(gin.TestMode) //setting so that we don't get debug logs during testing

	return InitRoutes(svc)
}

func TestSignupAPI_SQLite(t *testing.T) {
	router := setupSharedTestServer()

	body := []byte(`{
		"name": "Piyush",
		"email": "piyush@test.com",
		"password": "password123"
	}`)

	req, _ := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// assert.Contains(t, w.Body.String(), `"role":"user"`)
	assert.NotContains(t, w.Body.String(), `"password":"password123"`)
}
func TestLoginAPI_SQLite(t *testing.T) {
	router := setupSharedTestServer()

	// First signup
	signup := []byte(`{
		"name": "Piyush",
		"email": "login@test.com",
		"password": "secret"
	}`)
	req1, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(signup))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	// Then login
	login := []byte(`{
		"email": "login@test.com",
		"password": "secret"
	}`)
	req2, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(login))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
}
