package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/isdiemer/crossword-backend/internal/routes"
	"github.com/isdiemer/crossword-backend/internal/storage"
)

// buildCORS returns a cors.Config that
//   - echoes exact origins from ALLOWED_ORIGINS
//   - and green-lights any https://*.vercel.app preview URL.
func buildCORS() cors.Config {
	cfg := cors.Config{
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"Accept",
			"Cache-Control",
			"X-Requested-With",
			"Cookie",
			"Set-Cookie",
		},
		AllowCredentials: true,
		ExposeHeaders: []string{
			"Set-Cookie",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Credentials",
		},
		MaxAge: 12 * time.Hour,
	}

	// Exact-match origins from env (comma-separated)
	exact := map[string]struct{}{
		"https://crossword-frontend-one.vercel.app": struct{}{},
	}
	if v := os.Getenv("ALLOWED_ORIGINS"); v != "" {
		for _, o := range strings.Split(v, ",") {
			o = strings.TrimSpace(o)
			exact[o] = struct{}{}
			log.Printf("Added allowed origin: %s", o)
		}
	} else {
		log.Printf("Warning: No ALLOWED_ORIGINS environment variable set, using default Vercel domain")
	}

	cfg.AllowOriginFunc = func(origin string) bool {
		// Log the incoming origin
		log.Printf("CORS: Checking origin: %s", origin)

		// Always allow the main Vercel domain
		if origin == "https://crossword-frontend-one.vercel.app" {
			log.Printf("CORS: Origin %s allowed as main Vercel domain", origin)
			return true
		}

		// Check explicit list
		if _, ok := exact[origin]; ok {
			log.Printf("CORS: Origin %s allowed by exact match", origin)
			return true
		}

		// Allow Vercel preview deployments
		if strings.HasPrefix(origin, "https://") && strings.HasSuffix(origin, ".vercel.app") {
			log.Printf("CORS: Origin %s allowed as Vercel preview", origin)
			return true
		}

		log.Printf("CORS: Origin %s denied", origin)
		return false
	}

	return cfg
}

func main() {
	// Configure logging
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Gin in release mode unless GIN_MODE=debug
	if gin.Mode() != gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Load .env in local dev; ignore if it doesn't exist
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Log all environment variables
	log.Printf("Environment Configuration:")
	log.Printf("- GIN_MODE: %s", gin.Mode())
	log.Printf("- ALLOWED_ORIGINS: %s", os.Getenv("ALLOWED_ORIGINS"))
	log.Printf("- PORT: %s", os.Getenv("PORT"))

	// Configure CORS
	corsConfig := buildCORS()
	r.Use(cors.New(corsConfig))

	// Add middleware to log all requests
	r.Use(func(c *gin.Context) {
		log.Printf("Incoming request: %s %s", c.Request.Method, c.Request.URL.Path)
		log.Printf("Request headers:")
		for k, v := range c.Request.Header {
			log.Printf("  %s: %v", k, v)
		}
		c.Next()
	})

	// init DB & routes
	storage.InitDatabase()
	routes.RegisterRoutes(r)

	// Port from env, fallback 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
