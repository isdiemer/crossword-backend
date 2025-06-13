package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isdiemer/crossword-backend/internal/sessions"
	"github.com/isdiemer/crossword-backend/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("input received: %+v", input)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := storage.GetUserByUsername(input.Username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or password incorrect"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username or password incorrect"})
		return
	}

	token, err := sessions.Create(user.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error Inserting Token"})
		return
	}

	c.SetCookie("session_token", // name
		token, // value
		3600,  // expiration
		"/",   // path
		"",    // domain
		false, // secure
		true,  // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func LogoutHandler(c *gin.Context) {
	token, err := c.Cookie("session_token")
	if err == nil {
		_ = sessions.DropSessionByToken(token)
	}

	c.SetCookie("session_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

func DeleteHandler(c *gin.Context) {
	userIDRaw, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDRaw.(uint)

	// Remove all sessions
	_ = sessions.RemoveAllSessionsByID(userID)

	// Delete the user account
	err := storage.RemoveUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete user"})
		return
	}

	// Clear the cookie
	c.SetCookie("session_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "account deleted"})
}
