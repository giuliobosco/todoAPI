package controller

import (
	"net/http"

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
	user.Password = ""

	c.JSON(http.StatusOK, user)
}

// UpdateUser update user
func UpdateUser(c *gin.Context) {
	dbUser, err := getUserByContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}

	var user *model.User
	user, err = utils.UserValidator(c, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
	}

	if dbUser.Email != user.Email {
		if err = emailCheck(user.Email); err != nil {
			c.JSON(http.StatusConflict, gin.H{sMessage: err.Error()})
			return
		}

		user.Active = false
		var err error
		user.VerifyToken, err = utils.GenerateRandomStringURLSafe(config.TokenLength)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{sError: config.SUserFailUpdate})
			return
		}
		config.GetDB().Model(&dbUser).Update("email", user.Email)
		config.GetDB().Model(&dbUser).Update("active", user.Active)
		config.GetDB().Model(&dbUser).Update("verify_token", user.VerifyToken)
		utils.UserConfirmationSendMail(user)
	}

	config.GetDB().Model(&dbUser).Update("firstname", user.Firstname)
	config.GetDB().Model(&dbUser).Update("lastname", user.Lastname)

	c.JSON(http.StatusCreated, gin.H{sMessage: config.SUserUpdated})
}

// DeleteUser is the function for delete the user
func DeleteUser(c *gin.Context) {
	dbUser, err := getUserByContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}

	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}

	if !utils.ComparePasswordHash(dbUser.Password, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{sError: config.SWrongPassword})
		return
	}

	config.GetDB().Delete(dbUser)

	c.JSON(http.StatusOK, gin.H{sMessage: config.SUserDeleted, config.SUser: dbUser})
}
