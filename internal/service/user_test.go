package service_test

import (
	"testing"

	"github.com/isdiemer/crossword-backend/internal/model"
	"github.com/isdiemer/crossword-backend/internal/service"
	"github.com/isdiemer/crossword-backend/internal/storage"
	"golang.org/x/crypto/bcrypt"
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

func TestRegisterNewUser(t *testing.T) {
	setupTestDB(t)

	user, err := service.RegisterNewUser("alice", "alice@example.com", "password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.ID == 0 {
		t.Errorf("expected user ID to be set")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("password123")) != nil {
		t.Errorf("stored password hash does not match original")
	}

	fetched, err := storage.GetUserByUsername("alice")
	if err != nil {
		t.Fatalf("failed to fetch user: %v", err)
	}
	if fetched.ID != user.ID {
		t.Errorf("expected fetched ID %d got %d", user.ID, fetched.ID)
	}

	if _, err := service.RegisterNewUser("alice", "alice2@example.com", "pw"); err == nil {
		t.Errorf("expected error for duplicate username")
	}
}
