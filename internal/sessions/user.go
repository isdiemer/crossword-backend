package sessions

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/isdiemer/crossword-backend/internal/model"
	"github.com/isdiemer/crossword-backend/internal/storage"
)

// Creates a sessionID for a particular user and inserts into DB
func Create(userID uint) (token string, err error) {
	b := make([]byte, 32)
	rand.Read(b)
	token = hex.EncodeToString(b)

	sesh := model.Session{
		Token:     token,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	if err := storage.CreateSession(sesh); err != nil {
		return "", err
	}

	return token, nil
}

func GetSessionByToken(token string) (*model.Session, error) {
	return storage.GetSessionByToken(token)
}
func GetSessionByUsername(name string) (*model.Session, error) {
	return storage.GetSessionByUsername(name)
}
func DropSessionByToken(token string) error {
	return storage.DropSessionByToken(token)
}
func RemoveUserByID(ID uint) error {
	return storage.RemoveUserByID(ID)
}
func RemoveAllSessionsByID(ID uint) error {
	return storage.RemoveAllSessionsByID(ID)
}
