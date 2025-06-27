package main

import (
	"log"

	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/isdiemer/crossword-backend/internal/routes"
	"github.com/isdiemer/crossword-backend/internal/storage"
)

func main() {
	r := gin.Default()

	// Configure CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200", "http://127.0.0.1:4200"}
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Cookie"}

	r.Use(cors.New(config))

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
