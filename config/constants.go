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
