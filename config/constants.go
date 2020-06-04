package config

import (
	"os"

	"github.com/giuliobosco/todoAPI/model"
)

var (
	// URL application url
	URL = os.Getenv("URL")
)

const (
	// TokenLength is the length of the user activation token
	TokenLength = 64
	// IdentityKey represent the parameter used as connection key.
	IdentityKey = "id"
	// Key is the internal secret key of the API Engine.
	Key = "my_secret_key_8F6E2P"
	// SWelcome is the welcome string
	SWelcome = "Welcome to my Todo App"
	// SUserExists is the user already exists string
	SUserExists = "User already exists"
	// SUserMissingParams is the user missing parameters string
	SUserMissingParams = "Missing Username or password"
	// SUserCreated is the user created string
	SUserCreated = "User created successfully!"
	// SUserInvalid is the invalid user id string
	SUserInvalid = "Invalid user id"
	// SUserFailCreation is the user creation internal error string
	SUserFailCreation = "Error while creating user"
	// SUserConfirmed is the user confirmed string
	SUserConfirmed = "User confirmed!"
	// SUserNotConfirmed is the user not confirmed string
	SUserNotConfirmed = "User not confirmed!"
	// SUserPasswordRecoveryError is the user passwor recovery error string
	SUserPasswordRecoveryError = "Error while recovery user password."
	// SUserPasswordRecoveryMailSent is the user password recovery mail sent
	SUserPasswordRecoveryMailSent = "User password recovery mail sent."
	// SUserPasswordUpdated is the user password updated string
	SUserPasswordUpdated = "User password updated"
	// SMissingOldNewPassword is the missing old or new password string
	SMissingOldNewPassword = "Missing old or new password"
	// SMissingEmail is the missing email string
	SMissingEmail = "Missing: email"
	// SUserUpdated is the user updated string
	SUserUpdated = "User updated."
	// SUserNotFound is the user not found string
	SUserNotFound = "User not found"
	// SUserEmailAlreadyExists is the user email already exists string
	SUserEmailAlreadyExists = "The email address is already used."
	// SUserFailUpdate is the user fail update string
	SUserFailUpdate = "Error while updating user"
	// SWrongPassword is the wrong password string
	SWrongPassword = "Wrong password"
	// SUserDeleted is the user deleted string
	SUserDeleted = "User deleted"
	// SUser user string
	SUser = "user"
	// STaskCreated is the task created string
	STaskCreated = "Task created successfully!"
	// STaskNotFound is the task not found string
	STaskNotFound = "No todo found!"
	// STaskInvalid is the invalid todo id string
	STaskInvalid = "Invalid todo id"
	// STaskUpdated is the task updated string
	STaskUpdated = "Task updated successfully!"
	// STaskDeleted is the task delted string
	STaskDeleted = "Task deleted successfully!"
	// SMessage is the message string
	SMessage = "message"
	// SError is the error string
	SError = "error"
	// SData is the data string
	SData = "data"
	// STask is the task string
	STask = "task"
	// SExpire is the expire string
	SExpire = "expire"
	// SToken is the token string
	SToken = "token"
)

func BuildConfirmEmail(user model.User, smtpUsername string) []byte {
	var link string = URL + "v1/confirm?email=" + user.Email + "&token=" + user.VerifyToken

	return []byte("To: " + user.Email + "\r\n" +
		"From: " + smtpUsername + "\r\n" +
		"Subject: TodoAPI: confirm you email address!\r\n" +
		"\r\n" +
		"Hi " + user.Firstname + " " + user.Lastname + ",\r\n\r\n" +
		"Confirm your email address for todoAPI with the following link\r\n\r\n" +
		link + "\r\n\r\n" +
		"Thanks for using todoAPI\r\n" +
		"The todoAPI team\r\n")
}

func BuildPasswordRecovery(user model.User, smtpUsername string) []byte {
	var link string = URL + "v1/executePasswordRecovery?email=" + user.Email + "&token=" + user.VerifyToken

	return []byte("To: " + user.Email + "\r\n" +
		"From: " + smtpUsername + "\r\n" +
		"Subject: TodoAPI: Password recovery link!\r\n" +
		"\r\n" +
		"Hi " + user.Firstname + " " + user.Lastname + ",\r\n\r\n" +
		"Use the following link for recovery your password\r\n\r\n" +
		link + "\r\n\r\n" +
		"Thanks for using todoAPI\r\n" +
		"The todoAPI team\r\n")
}
