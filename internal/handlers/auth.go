package handlers

import (
	"log"
	"net/http"
	"os"
	"strings"

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

	// Log all request headers for debugging
	log.Printf("All request headers:")
	for k, v := range c.Request.Header {
		log.Printf("%s: %v", k, v)
	}

	// Determine cookie domain
	origin := c.GetHeader("Origin")
	log.Printf("Request Origin: %s", origin)
	host := c.Request.Host
	domain := os.Getenv("COOKIE_DOMAIN")
	if domain == "" {
		domain = strings.Split(host, ":")[0]
	}
	if strings.Contains(host, "vercel.app") {
		// allow all Vercel preview subdomains to share the cookie
		domain = ".vercel.app"
	}

	// Set cookie attributes
	cookieName := "session_token"
	cookieValue := token
	maxAge := 3600
	path := "/"
	secure := true
	httpOnly := true

	// Create a new cookie instance
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    cookieValue,
		Path:     path,
		Domain:   domain,
		MaxAge:   maxAge,
		Secure:   secure,
		HttpOnly: httpOnly,
		SameSite: http.SameSiteNoneMode,
	}

	// Set the cookie using http.SetCookie
	http.SetCookie(c.Writer, cookie)

	// Log response headers
	log.Printf("Response headers being set:")
	for k, v := range c.Writer.Header() {
		log.Printf("%s: %v", k, v)
	}

	// Try to read the cookie back
	if cookie, err := c.Cookie("session_token"); err != nil {
		log.Printf("Warning: Could not read back cookie: %v", err)
	} else {
		log.Printf("Cookie successfully set and readable: %v", cookie != "")
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"debug_info": gin.H{
			"origin":      origin,
			"domain":      domain,
			"headers_set": c.Writer.Header(),
		},
	})
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
