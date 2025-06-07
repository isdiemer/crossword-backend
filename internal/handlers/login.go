package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isdiemer/crossword-backend/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}

	user, err := storage.GetUserByUsername(input.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or password incorrect"})
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or password incorrect"})
		return
	}

	token := GenerateToken()
	sessions[token] = fmt.Sprint(user.ID)

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func GenerateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
