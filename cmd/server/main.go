package main

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/isdiemer/crossword-backend/internal/routes"
	"github.com/isdiemer/crossword-backend/internal/storage"
)

func main() {
	r := gin.Default()

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Configure CORS middleware
	config := cors.DefaultConfig()

	// Get allowed origins from environment variable, fallback to localhost if not set
	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		// Default to localhost and Vercel preview URLs
		allowedOrigins = "http://localhost:4200,http://localhost:3000,https://*.vercel.app"
	}

	// Handle wildcard domains for Vercel
	origins := strings.Split(allowedOrigins, ",")
	if containsWildcard(origins) {
		config.AllowOriginFunc = func(origin string) bool {
			return isAllowedOrigin(origin, origins)
		}
	} else {
		config.AllowOrigins = origins
	}

	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Cookie"}

	r.Use(cors.New(config))

	storage.InitDatabase()
	routes.RegisterRoutes(r)

	// Get port from environment variable, fallback to 8080 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// Helper function to check if origins contain wildcard
func containsWildcard(origins []string) bool {
	for _, origin := range origins {
		if strings.Contains(origin, "*") {
			return true
		}
	}
	return false
}

// Helper function to check if an origin matches the allowed patterns
func isAllowedOrigin(origin string, allowedPatterns []string) bool {
	for _, pattern := range allowedPatterns {
		if strings.Contains(pattern, "*") {
			// Convert wildcard pattern to regex
			regexPattern := strings.Replace(pattern, "*", ".*", -1)
			if matched, _ := regexp.MatchString(regexPattern, origin); matched {
				return true
			}
		} else if origin == pattern {
			return true
		}
	}
	return false
}
