package mock

import "github.com/giuliobosco/todoAPI/model"

// GetMockUser build a mock user for tests
func GetMockUser() model.User {
	return model.User{
		Email:     "first.last@email.com",
		Firstname: "firstname",
		Lastname:  "lastname",
	}
}
