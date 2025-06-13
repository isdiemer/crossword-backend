package storage

import (
	"log"
	"os"

	"github.com/isdiemer/crossword-backend/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error

	dsn := os.Getenv("DATABASE_URL")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to Supabase:", err)
	}

	err = DB.AutoMigrate(&model.User{}, &model.Session{})
	if err != nil {
		log.Fatal("failed to migrate:", err)
	}
}
