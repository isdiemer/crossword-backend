package storage

import (
	"github.com/isdiemer/crossword-backend/internal/model"
)

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
func CreateSession(session model.Session) error {
	return DB.Create(&session).Error
}
func GetSessionByToken(token string) (*model.Session, error) {
	var session model.Session
	result := DB.Where("token = ?", token).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func GetSessionByUsername(name string) (*model.Session, error) {
	var session model.Session
	result := DB.Where("Username = ?", name).First(&session)
	if result.Error != nil {
		return nil, result.Error
	}
	return &session, nil
}

func DropSessionByToken(token string) error {
	return DB.Delete("token = ?", token).Delete(&model.Session{}).Error
}
func RemoveUserByID(ID uint) error {
	return DB.Where("id = ?", ID).Delete(&model.User{}).Error
}

func RemoveAllSessionsByID(ID uint) error {
	return DB.Where("user_id = ?", ID).Delete(&model.Session{}).Error
}
