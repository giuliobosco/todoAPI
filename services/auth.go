package services

import (
	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/model"
)

// GetUserByID Gets the user by id
func GetUserByID(id int) model.User {
	var user model.User
	config.GetDB().Where("id = ?", id).First(&user)

	return user
}
