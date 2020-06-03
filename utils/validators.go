package utils

import (
	"errors"

	"github.com/giuliobosco/todoAPI/model"

	"github.com/badoux/checkmail"
)

// EmailValidator validate email address, by his format and the host
func EmailValidator(email string) (bool, error) {
	if err := checkmail.ValidateFormat(email); err != nil {
		return false, err
	}

	if err := checkmail.ValidateHost(email); err != nil {
		if smtpErr, ok := err.(checkmail.SmtpError); ok && err != nil {
			return false, smtpErr
		}
		return false, err
	}

	return true, nil
}

// UserValidator validate user parameters
func UserValidator(user model.User) (bool, error) {
	var missing []string

	if len(user.Email) == 0 {
		missing = append(missing, "email")
	}
	if len(user.Password) == 0 {
		missing = append(missing, "password")
	}
	if len(user.Firstname) == 0 {
		missing = append(missing, "Firstname")
	}
	if len(user.Firstname) == 0 {
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
