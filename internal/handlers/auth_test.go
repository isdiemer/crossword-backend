package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/isdiemer/crossword-backend/internal/handlers"
	"github.com/isdiemer/crossword-backend/internal/model"
	"github.com/isdiemer/crossword-backend/internal/service"
	"github.com/isdiemer/crossword-backend/internal/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	storage.DB = db
	if err := db.AutoMigrate(&model.User{}, &model.Session{}, &model.Puzzle{}, &model.Guess{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	setupTestDB(t)

	_, err := service.RegisterNewUser("carol", "carol@example.com", "mypassword")
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	body := bytes.NewBufferString(`{"username":"carol","password":"mypassword"}`)
	req, _ := http.NewRequest(http.MethodPost, "/login", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handlers.LoginHandler(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 got %d", w.Code)
	}
	cookieFound := false
	for _, ck := range w.Result().Cookies() {
		if ck.Name == "session_token" {
			cookieFound = true
			break
		}
	}
	if !cookieFound {
		t.Errorf("session cookie not set")
	}
}
