package handlers

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", PingHandler)
	r.POST("/register", RegisterUser)
	r.POST("/login", LoginHandler)
	r.POST("/puzzles/:id/validate", ValidatePuzzle)
}
