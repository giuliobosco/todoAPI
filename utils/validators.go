package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/model"

	"github.com/badoux/checkmail"
)

// EmailValidator validate email address, by his format and the host
func EmailValidator(email string) (bool, error) {
	if err := checkmail.ValidateFormat(email); err != nil {
		return false, err
	}

	return true, nil
}

// UserValidator validate user parameters
func UserValidator(user model.User, usePassword bool) (bool, error) {
	var missing []string

	if len(user.Email) == 0 {
		missing = append(missing, "email")
	}
	if len(user.Password) == 0 && usePassword {
		missing = append(missing, "password")
	}
	if len(user.Firstname) == 0 {
		missing = append(missing, "Firstname")
	}
	if len(user.Lastname) == 0 {
		missing = append(missing, "Lastname")
	}

	if len(missing) > 0 {
		var errorString string = "Missing: "

		for i, m := range missing {
			if i > 0 {
				errorString += ","
			}
			errorString += " " + m
		}

		return false, errors.New(errorString)
	}

	if ok, err := EmailValidator(user.Email); !ok {
		return false, errors.New("Email: " + err.Error())
	}

	return true, nil
}

func ConfirmUserValidator(m map[string][]string) (*model.User, error) {
	var missing []string

	if m["email"] == nil || len(m["email"]) == 0 {
		missing = append(missing, "email")
	}
	if m["token"] == nil || len(m["token"]) == 0 {
		missing = append(missing, "token")
	}

	if len(missing) > 0 {
		var errorString string = "Missing: "

		for i, m := range missing {
			if i > 0 {
				errorString += ","
			}
			errorString += " " + m
		}

		return nil, errors.New(errorString)
	}

	e := m["email"][0]
	t := m["token"][0]

	var userCheck model.User
	config.GetDB().Where("email = ? AND verify_token = ?", e, t).First(&userCheck)

	if userCheck.ID == 0 {
		return nil, errors.New("Not valid request")
	}

	return &userCheck, nil
}

type PasswordRecovery struct {
	Email       string `json:"email"`
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func PasswordRecoveryValidator(c *gin.Context) (*model.User, error) {
	var missing []string

	var pr PasswordRecovery
	if err := c.ShouldBindJSON(&pr); err != nil {
		return nil, err
	}

	if len(pr.Email) == 0 {
		missing = append(missing, "email")
	}
	if len(pr.Token) == 0 {
		missing = append(missing, "token")
	}
	if len(pr.NewPassword) == 0 {
		missing = append(missing, "new_password")
	}

	if len(missing) > 0 {
		var errorString = "Missing: "

		for i, m := range missing {
			if i > 0 {
				errorString += ","
			}
			errorString += " " + m
		}

		return nil, errors.New(errorString)
	}
	var user model.User
	config.GetDB().Where("email = ? AND verify_token = ?", pr.Email, pr.Token).First(&user)

	if user.ID == 0 {
		return nil, errors.New(config.SUserPasswordRecoveryError)
	}

	user.VerifyToken = ""
	user.Password = pr.NewPassword

	return &user, nil
}
