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
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "Cookie", "Set-Cookie"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Set-Cookie"},
		MaxAge:           12 * time.Hour,
	}

	// Exact-match origins from env (comma-separated)
	exact := map[string]struct{}{}
	if v := os.Getenv("ALLOWED_ORIGINS"); v != "" {
		for _, o := range strings.Split(v, ",") {
			o = strings.TrimSpace(o)
			exact[o] = struct{}{}
			log.Printf("Added allowed origin: %s", o)
		}
	} else {
		log.Printf("Warning: No ALLOWED_ORIGINS environment variable set")
	}

	cfg.AllowOriginFunc = func(origin string) bool {
		// Log the incoming origin
		log.Printf("Checking origin: %s", origin)

		// 1️⃣ explicit list
		if _, ok := exact[origin]; ok {
			log.Printf("Origin %s allowed by exact match", origin)
			return true
		}
		// 2️⃣ any Vercel preview, e.g. https://my-branch--app.vercel.app
		isVercel := strings.HasPrefix(origin, "https://") && strings.HasSuffix(origin, ".vercel.app")
		if isVercel {
			log.Printf("Origin %s allowed as Vercel preview", origin)
		}
		return isVercel
	}

	return cfg
}

func main() {
	// Gin in release mode unless GIN_MODE=debug
	if gin.Mode() != gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()

	// Load .env in local dev; ignore if it doesn't exist
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Log environment variables
	log.Printf("Environment Mode: %s", gin.Mode())
	log.Printf("Cookie Domain: %s", os.Getenv("COOKIE_DOMAIN"))
	log.Printf("Allowed Origins: %s", os.Getenv("ALLOWED_ORIGINS"))

	corsConfig := buildCORS()
	r.Use(cors.New(corsConfig))

	// init DB & routes
	storage.InitDatabase()
	routes.RegisterRoutes(r)

	// Port from env, fallback 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
