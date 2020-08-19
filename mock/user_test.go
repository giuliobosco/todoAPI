package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetMockUserID0NoPassowrd checks the random generation of the user
// with ID 0 and password
func TestGetMockUserID0NoPassword(t *testing.T) {
	for i := 0; i < runTests; i++ {
		u := GetMockUserID0(false)

		assert.True(t, uint(0) == u.ID, "ID of the user should be 0")
		assert.Empty(t, u.Password)
		assert.NotEmpty(t, u.Email)
		assert.NotEmpty(t, u.Firstname)
		assert.NotEmpty(t, u.Lastname)
	}
}

// TestGetMockUserID0 checks the random generation of the user with ID 0
func TestGetMockUserID0(t *testing.T) {
	for i := 0; i < runTests; i++ {
		u := GetMockUserID0(true)

		assert.True(t, uint(0) == u.ID, "ID of the user should be 0")
		assert.NotEmpty(t, u.Password)
		assert.NotEmpty(t, u.Email)
		assert.NotEmpty(t, u.Firstname)
		assert.NotEmpty(t, u.Lastname)
	}
}

// TestGetMockUserNoPassword checks the random generation of the user
func TestGetMockUserNoPassword(t *testing.T) {
	for i := 0; i < runTests; i++ {
		u := GetMockUser(false)

		assert.True(t, uint(0) != u.ID, "ID of the user should not be 0")
		assert.Empty(t, u.Password)
		assert.NotEmpty(t, u.Email)
		assert.NotEmpty(t, u.Firstname)
		assert.NotEmpty(t, u.Lastname)
	}
}

// TestGetMockUser checks the random generation of the user
func TestGetMockUser(t *testing.T) {
	for i := 0; i < runTests; i++ {
		u := GetMockUser(true)

		assert.True(t, uint(0) != u.ID, "ID of the user should not be 0")
		assert.NotEmpty(t, u.Password)
		assert.NotEmpty(t, u.Email)
		assert.NotEmpty(t, u.Firstname)
		assert.NotEmpty(t, u.Lastname)
	}
}
