package services

import (
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
