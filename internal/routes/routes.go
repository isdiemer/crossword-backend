package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/isdiemer/crossword-backend/internal/handlers"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", handlers.PingHandler)
	r.GET("/me", handlers.AuthMiddleware, handlers.MeHandler)
	r.POST("/register", handlers.RegisterUser)
	r.POST("/login", handlers.LoginHandler)
	r.POST("/logout", handlers.LogoutHandler)
	r.POST("/delete", handlers.AuthMiddleware, handlers.DeleteHandler)
	r.GET("/my-puzzles", handlers.AuthMiddleware, handlers.GetMyPuzzlesHandler)
	r.GET("/puzzles/:id", handlers.GetPuzzleByIDHandler)
	r.POST("/puzzles", handlers.AuthMiddleware, handlers.CreatePuzzleHandler)
	r.PUT("/puzzles/:id", handlers.AuthMiddleware, handlers.UpdatePuzzleHandler)
	r.POST("/guess", handlers.AuthMiddleware, handlers.SubmitGuessHandler)
	r.POST("/puzzles/:id/validate-guess", handlers.AuthMiddleware, handlers.ValidateGuessHandler)
	r.DELETE("/puzzles/:id", handlers.AuthMiddleware, handlers.DeletePuzzleHandler)

}
