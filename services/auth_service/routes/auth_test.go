package routes

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/piyushsharma67/movie_booking/services/auth_service/utils"
)

func TestSignup_Success(t *testing.T) {
	r, _ := setupTestRouter(t)

	body := `{
		"name": "Piyush",
		"email": "piyush@test.com",
		"password": "secret123"
	}`

	req := httptest.NewRequest("POST", "/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestSignup_DuplicateEmail(t *testing.T) {
	r, db := setupTestRouter(t)

	seedUser(t, db, "dup@test.com", "secret", "user")

	body := `{
		"name": "Test",
		"email": "dup@test.com",
		"password": "secret"
	}`

	req := httptest.NewRequest("POST", "/signup", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	r, db := setupTestRouter(t)

	seedUser(t, db, "login@test.com", "password123", "user")

	body := `{
		"email": "login@test.com",
		"password": "password123"
	}`

	req := httptest.NewRequest("GET", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestLogin_InvalidPassword(t *testing.T) {
	r, db := setupTestRouter(t)

	seedUser(t, db, "badpass@test.com", "correct", "user")

	body := `{
		"email": "badpass@test.com",
		"password": "wrong"
	}`

	req := httptest.NewRequest("GET", "/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestValidate_MissingToken(t *testing.T) {
	r, _ := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/validate", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestValidate_ValidToken(t *testing.T) {
	r, _ := setupTestRouter(t)

	token, err := utils.GenerateJWT("1", "test@test.com", "user", os.Getenv("JWT_SECRET"))
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("GET", "/validate", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if w.Header().Get("X-User-Id") == "" {
		t.Fatal("expected X-User-Id header")
	}
}
