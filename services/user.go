package services

import (
	"errors"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/model"
)

// UpdateUserEmail updates user email, active flag and verify token
func UpdateUserEmail(o, n *model.User) {
	config.GetDB().Model(&o).Update("email", n.Email)
	config.GetDB().Model(&o).Update("active", n.Active)
	config.GetDB().Model(&o).Update("verify_token", n.VerifyToken)
}

// UpdateUserAnagraphic updates user firstname and lastname
func UpdateUserAnagraphic(o, n *model.User) {
	config.GetDB().Model(&o).Update("firstname", n.Firstname)
	config.GetDB().Model(&o).Update("lastname", n.Lastname)
}

// VerifyUserEmailToken checks the user email and tokens
func VerifyUserEmailToken(e, t string) (*model.User, error) {
	var userCheck model.User
	config.GetDB().Where("email = ? AND verify_token = ?", e, t).First(&userCheck)

	if userCheck.ID == 0 {
		return nil, errors.New("Not valid request")
	}

	return &userCheck, nil
}

// EmailExists checks if the email address exists in the db
func EmailExists(email string) bool {
	var u model.User
	config.GetDB().First(&u, "email = ?", email)

	return u.ID > 0
}
