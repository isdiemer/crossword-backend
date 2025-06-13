package storage

import (
	"log"

	"github.com/isdiemer/crossword-backend/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("crossword.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	err = DB.AutoMigrate(&model.User{}, &model.Puzzle{})
	if err != nil {
		log.Fatal("failed to migrate database:", err)
	}
}
