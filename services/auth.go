package services

import (
	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/model"
)

// GetUserByID Gets the user by id
func GetUserByID(id uint) model.User {
	var user model.User
	config.GetDB().Where("id = ?", id).First(&user)

	return user
}

// GetUserByEmail Get the user by email
func GetUserByEmail(email string) model.User {
	var user model.User
	config.GetDB().Where("email = ?", email).First(&user)

	return user
}

// EmptyUserVerifyToken empty the verify token of the user
func EmptyUserVerifyToken(user *model.User) {
	config.GetDB().Model(&user).Update("verify_token", "")
}
