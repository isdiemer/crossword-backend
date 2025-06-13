package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isdiemer/crossword-backend/internal/model"
	"github.com/isdiemer/crossword-backend/internal/storage"
)

const sessionKey = "session"

func AuthMiddleware(c *gin.Context) {

	token, err := c.Cookie("session_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing session"})
		return
	}

	session, err := storage.GetSessionByToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	c.Set(sessionKey, session)
	c.Set("userID", session.UserID)
	c.Next()
}

func GetSessionFromContext(c *gin.Context) (*model.Session, error) {
	val, ok := c.Get(sessionKey)
	if !ok {
		return nil, errors.New("session not found in context")
	}
	session, ok := val.(*model.Session)
	if !ok {
		return nil, errors.New("session has wrong type")
	}
	return session, nil

}
