package mock

import (
	"github.com/giuliobosco/todoAPI/model"
	"github.com/giuliobosco/todoAPI/tu"
)

// GetMockUserID0 build a user for tests with ID 0
func GetMockUserID0(password bool) model.User {
	u := model.User{
		Email:     tu.RandomEmail(),
		Firstname: tu.RandomString12(),
		Lastname:  tu.RandomString12(),
	}

	if password {
		u.Password = tu.RandomString12()
	}

	return u
}

// GetMockUser build a user for tests
func GetMockUser(password bool) model.User {
	u := GetMockUserID0(password)
	u.ID = tu.RandomUintNo0()

	return u
}

// GetMapByUser Gets the user as map
func GetMapByUser(u model.User) map[string]interface{} {
	return map[string]interface{}{
		"id":           u.ID,
		"email":        u.Email,
		"firstname":    u.Firstname,
		"lastname":     u.Lastname,
		"active":       u.Active,
		"verify_token": u.VerifyToken,
		"password":     u.Password,
	}
}

// GetMapArrayByUser gets an array of users as map
func GetMapArrayByUser(u model.User) []map[string]interface{} {
	return []map[string]interface{}{GetMapByUser(u)}
}

// GetLoginVals conver user to only login vals
func GetLoginVals(u model.User) model.User {
	return model.User{Email: u.Email, Password: u.Password}
}
