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

func GetMapByUser(u model.User) []map[string]interface{} {
	return []map[string]interface{}{{
		"id":        u.ID,
		"email":     u.Email,
		"firstname": u.Firstname,
		"lastname":  u.Lastname,
		"active":    u.Active,
	}}
}
