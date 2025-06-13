package sessions_test

import (
	"testing"

	"github.com/isdiemer/crossword-backend/internal/model"
	"github.com/isdiemer/crossword-backend/internal/service"
	"github.com/isdiemer/crossword-backend/internal/sessions"
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

func TestCreateSession(t *testing.T) {
	setupTestDB(t)

	user, err := service.RegisterNewUser("bob", "bob@example.com", "secret123")
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}

	token, err := sessions.Create(user.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(token) != 64 {
		t.Errorf("expected token length 64 got %d", len(token))
	}

	sess, err := storage.GetSessionByToken(token)
	if err != nil {
		t.Fatalf("session not persisted: %v", err)
	}
	if sess.UserID != user.ID {
		t.Errorf("expected userID %d got %d", user.ID, sess.UserID)
	}
}
