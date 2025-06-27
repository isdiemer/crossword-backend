package service

import (
	"log"

	"github.com/isdiemer/crossword-backend/internal/model"
	"github.com/isdiemer/crossword-backend/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

func RegisterNewUser(username, email, pw string) (model.User, error) {
	log.Printf("Service: Starting registration for username: %s", username)

	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Service: Password hashing error: %v", err)
		return model.User{}, err
	}

	log.Printf("Service: Password hashed successfully")

	newUser := model.User{
		ID:           0,
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
	}

	log.Printf("Service: Attempting to create user in database")
	userErr := storage.CreateUser(&newUser)
	if userErr != nil {
		log.Printf("Service: Database creation error: %v", userErr)
		return model.User{}, userErr
	}

	log.Printf("Service: User created successfully with ID: %d", newUser.ID)
	return newUser, nil
}
