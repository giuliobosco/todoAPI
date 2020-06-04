package controller

import (
	"errors"
	"net/http"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/model"
	"github.com/giuliobosco/todoAPI/utils"

	jwtapple2 "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const sMessage string = config.SMessage
const sError string = config.SError
const sData string = config.SData

// RegisterEndPoint registration API End Point
func RegisterEndPoint(c *gin.Context) {
	user, err := utils.UserValidator(c, true)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
	}

	if err = emailCheck(user.Email); err != nil {
		c.JSON(http.StatusConflict, gin.H{sMessage: err.Error()})
		return
	}

	user.Active = false
	user.VerifyToken, err = utils.GenerateRandomStringURLSafe(config.TokenLength)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{sError: config.SUserFailCreation})
		return
	}

	user.Password, err = utils.PasswordHash(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{sError: err.Error()})
		return
	}

	config.GetDB().Save(&user)
	utils.UserConfirmationSendMail(user)

	c.JSON(http.StatusCreated, gin.H{sMessage: config.SUserCreated})
}

// ConfirmUser is the function for confirm a user
func ConfirmUser(c *gin.Context) {
	p := c.Request.URL.Query()
	user, err := utils.ConfirmUserValidator(p)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}

	config.GetDB().Model(&user).Update("active", true)
	config.GetDB().Model(&user).Update("verify_token", "")

	c.JSON(http.StatusOK, gin.H{sMessage: config.SUserConfirmed})
}

func getUserByEmailParam(c *gin.Context) (*model.User, error) {
	p := c.Request.URL.Query()

	if p["email"] == nil || len(p["email"]) == 0 {
		return nil, errors.New(config.SMissingEmail)
	}

	var user model.User
	config.GetDB().Where("email = ?", p["email"][0]).First(&user)

	if user.ID <= 0 {
		return nil, errors.New(config.SUserInvalid)
	}

	return &user, nil
}

func SendUserConfirmAgain(c *gin.Context) {
	user, err := getUserByEmailParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}

	user.VerifyToken, err = utils.GenerateRandomStringURLSafe(config.TokenLength)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{sError: config.SUserFailCreation})
		return
	}

	config.GetDB().Model(user).Update("verify_token", user.VerifyToken)
	utils.UserConfirmationSendMail(user)

	c.JSON(http.StatusOK, gin.H{sMessage: config.SUserSentConfirmationMailAgain})
}

// RequestPasswordRecovery is the function for request the password recovery
func RequestPasswordRecovery(c *gin.Context) {
	user, err := getUserByEmailParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}

	if !user.Active {
		c.JSON(http.StatusBadRequest, gin.H{sError: config.SUserNotConfirmed})
		return
	}

	user.VerifyToken, err = utils.GenerateRandomStringURLSafe(config.TokenLength)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{sError: config.SUserPasswordRecoveryError})
		return
	}

	config.GetDB().Model(user).Update(user)
	utils.UserPasswordRecoverySendMail(user)

	c.JSON(http.StatusOK, gin.H{sMessage: config.SUserPasswordRecoveryMailSent})
}

// ExecutePasswordRecovery is the function for execute the password recovery
func ExecutePasswordRecovery(c *gin.Context) {
	user, err := utils.PasswordRecoveryValidator(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}

	user.Password, err = utils.PasswordHash(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{sError: err.Error()})
		return
	}

	config.GetDB().Model(&user).Update("verify_token", user.VerifyToken)
	config.GetDB().Model(&user).Update("password", user.Password)

	c.JSON(http.StatusOK, gin.H{sMessage: config.SUserPasswordUpdated})
}

func getUserByContext(c *gin.Context) (*model.User, error) {
	cl := jwtapple2.ExtractClaims(c)

	var u model.User
	config.GetDB().Where("id = ?", cl[config.IdentityKey]).First(&u)

	if u.ID <= 0 {
		return nil, errors.New(config.SUserInvalid)
	}

	return &u, nil
}

func emailCheck(email string) error {
	var u model.User
	config.GetDB().First(&u, "email = ?", email)

	if u.ID > 0 {
		return errors.New(config.SUserEmailAlreadyExists)
	}

	return nil
}

// UpdatePasswordObj is the object containing the update password data
type UpdatePasswordObj struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// UpdatePassword function
func UpdatePassword(c *gin.Context) {
	user, err := getUserByContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}

	var pr UpdatePasswordObj
	if err = c.ShouldBindJSON(&pr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}
	if len(pr.OldPassword) == 0 || len(pr.NewPassword) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{sError: config.SMissingOldNewPassword})
		return
	}

	if !utils.ComparePasswordHash(user.Password, pr.OldPassword) {
		c.JSON(http.StatusBadRequest, gin.H{sError: config.SWrongPassword})
		return
	}

	pr.NewPassword, err = utils.PasswordHash(pr.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{sError: err.Error()})
		return
	}

	config.GetDB().Model(&user).Update("password", pr.NewPassword)

	c.JSON(http.StatusOK, gin.H{sMessage: config.SUserPasswordUpdated})
}

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
