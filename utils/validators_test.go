package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailValidator(t *testing.T) {
	// working email
	email := "giuliobva@gmail.com"
	ok, _ := EmailValidator(email)
	assert.True(t, ok)

	// not working email
	emails := []string{"giuliobva", "giuliobva@", "@gmail.com", "giuliobva@."}
	for _, e := range emails {
		ok, _ = EmailValidator(e)
		assert.False(t, ok)
	}
}

/*
func TestUserValidator(t *testing.T) {
	users := []model.User{
		{},
		{Email: "a"},
		{Email: "a", Password: "b"},
		{Email: "a", Password: "b", Firstname: "c"},
		{Email: "a", Password: "b", Firstname: "c", Lastname: "d"},
	}

	for _, user := range users {
		ok, err := UserValidator(user)
		assert.False(t, ok, err)
	}

	user := model.User{Email: "giuliobva@gmail.com", Password: "b", Firstname: "c", Lastname: "d"}
	ok, err := UserValidator(user)
	assert.True(t, ok, err)
}*/
