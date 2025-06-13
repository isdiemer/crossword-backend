package service

import (
	"log"

	"github.com/isdiemer/crossword-backend/internal/model"
	"github.com/isdiemer/crossword-backend/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

func RegisterNewUser(username, email, pw string) (model.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	newUser := model.User{
		ID:           0,
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
	}
	userErr := storage.CreateUser(&newUser)
	if userErr != nil {
		return model.User{}, userErr
	}

	return newUser, nil
}
