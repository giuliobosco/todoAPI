package controller

import (
	"net/http"

	"github.com/giuliobosco/todoAPI/services"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/model"
	"github.com/giuliobosco/todoAPI/utils"

	"github.com/gin-gonic/gin"
)

// FetchUser is the funciton for load the actual user data
func FetchUser(c *gin.Context) {
	user, err := getUserByContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser update user
func UpdateUser(c *gin.Context) {
	dbUser, err := getUserWithPasswordByContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}

	var user *model.User
	user, err = utils.UserValidator(c, false)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}

	if dbUser.Email != user.Email {
		if len(dbUser.Password) <= 0 {
			c.JSON(http.StatusConflict, gin.H{sError: config.SUserFailUpdate})
			return
		}

		if services.EmailExists(user.Email) {
			c.JSON(http.StatusConflict, gin.H{sMessage: config.SUserEmailAlreadyExists})
			return
		}

		user.Active = false
		var err error
		user.VerifyToken, err = utils.GenerateRandomStringURLSafe(config.TokenLength)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{sError: config.SUserFailUpdate})
			return
		}

		services.UpdateUserEmail(dbUser, user)

		if !utils.IsTesting() {
			utils.UserConfirmationSendMail(user)
		}
	}

	services.UpdateUserAnagraphic(dbUser, user)

	c.JSON(http.StatusOK, gin.H{sMessage: config.SUserUpdated})
}

// DeleteUser is the function for delete the user
func DeleteUser(c *gin.Context) {
	dbUser, err := getUserWithPasswordByContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}

	if !utils.ComparePasswordHash(dbUser.Password, user.Password) && len(dbUser.Password) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{sError: config.SWrongPassword})
		return
	}

	config.GetDB().Delete(dbUser)

	c.JSON(http.StatusOK, gin.H{sMessage: config.SUserDeleted, config.SUser: dbUser})
}
