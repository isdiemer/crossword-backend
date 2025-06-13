package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"github.com/isdiemer/crossword-backend/internal/routes"
	"github.com/isdiemer/crossword-backend/internal/storage"
)

func main() {
	r := gin.Default()
	err := godotenv.Load()
	storage.InitDatabase()
	if err != nil {
		log.Println("No .env file found")
	}

	routes.RegisterRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
