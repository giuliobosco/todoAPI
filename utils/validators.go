package utils

import (
	"errors"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/model"
	"github.com/giuliobosco/todoAPI/services"

	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
)

// EmailValidator validate email address, by his format and the host
func EmailValidator(email string) (bool, error) {
	if err := checkmail.ValidateFormat(email); err != nil {
		return false, err
	}

	return true, nil
}

// UserValidator validate user parameters
func UserValidator(c *gin.Context, usePassword bool) (*model.User, error) {
	var u model.User
	if err := c.ShouldBindJSON(&u); err != nil {
		return nil, err
	}

	var missing []string

	if len(u.Email) == 0 {
		missing = append(missing, "email")
	}
	if len(u.Password) == 0 && usePassword {
		missing = append(missing, "password")
	}
	if len(u.Firstname) == 0 {
		missing = append(missing, "Firstname")
	}
	if len(u.Lastname) == 0 {
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

		return nil, errors.New(errorString)
	}

	if ok, err := EmailValidator(u.Email); !ok {
		return nil, errors.New("Email: " + err.Error())
	}

	return &u, nil
}

// ConfirmUserValidator confirmation of the user (email link)
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

	return services.VerifyUserEmailToken(e, t)
}

// PasswordRecovery used for context checks
type PasswordRecovery struct {
	Email       string `json:"email"`
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// PasswordRecoveryValidator checks the password recovey validation, returns user with new password
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

	user, err := services.VerifyUserEmailToken(pr.Email, pr.Token)

	if err != nil {
		return nil, errors.New(config.SUserPasswordRecoveryError)
	}

	user.VerifyToken = ""
	user.Password = pr.NewPassword

	return user, nil
}

// TaskValidator validate task title
func TaskValidator(c *gin.Context) (*model.Task, error) {
	var t model.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		return nil, err
	}

	if len(t.Title) == 0 {
		return nil, errors.New("Missing: title")
	}

	return &t, nil
}
