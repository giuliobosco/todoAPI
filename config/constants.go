package config

const (
	// URL application url
	URL = "http://localhost:8080"
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
