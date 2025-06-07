package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/isdiemer/crossword-backend/internal/handlers"
	"github.com/isdiemer/crossword-backend/internal/storage"
)

func main() {
	r := gin.Default()
	storage.InitDatabase()

	handlers.RegisterRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
