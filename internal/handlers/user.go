package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isdiemer/crossword-backend/internal/service"
	"github.com/isdiemer/crossword-backend/internal/storage"
)

func RegisterUser(c *gin.Context) {
	var input RegisterInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid input"})
		return
	}
	user, err := service.RegisterNewUser(input.Username, input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Username Taken!": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"created":  user.CreatedAt,
	})
}
func PingHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

func MeHandler(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID := userIDVal.(uint)

	user, err := storage.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}
