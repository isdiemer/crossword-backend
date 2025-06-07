package storage

import "github.com/isdiemer/crossword-backend/internal/model"

func CreateUser(user *model.User) error {
	result := DB.Create(user)
	return result.Error
}

func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByID(ID uint) (*model.User, error) {
	var user model.User
	result := DB.Where("ID = ?", ID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	result := DB.Where("Username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
