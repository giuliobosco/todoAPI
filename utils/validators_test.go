package utils

import (
	"strings"
	"testing"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/mock"
	"github.com/giuliobosco/todoAPI/model"
	"github.com/giuliobosco/todoAPI/tu"

	"github.com/gin-gonic/gin"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
)

// TestEmailValidator
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

// TestUserValidator full user with password
func TestUserValidatorPassword(t *testing.T) {
	c := tu.GetContext()

	expectedUser := mock.GetMockUser(true)

	req, err := tu.GetRequestPost(expectedUser, "/")

	c.Request = req

	actualUser, err := UserValidator(c, true)

	assert.Equal(t, &expectedUser, actualUser)
	assert.Nil(t, err)
}

// TestUserValidatorNoPassword, user without password
func TestUserValidatorNoPassword(t *testing.T) {
	c := tu.GetContext()

	expectedUser := mock.GetMockUser(false)

	req, err := tu.GetRequestPost(expectedUser, "/")

	c.Request = req

	userActual, err := UserValidator(c, false)

	assert.Equal(t, &expectedUser, userActual)
	assert.Nil(t, err)
}

// testUserValidatorErrors, automatically test error for validator (automated)
func testUserValidatorErrors(t *testing.T, u model.User, missing []string, usePassword bool) {
	c := tu.GetContext()

	req, err := tu.GetRequestPost(u, "/")
	assert.Nil(t, err)

	c.Request = req

	actualUser, err := UserValidator(c, usePassword)

	assert.Nil(t, actualUser)

	for _, v := range missing {
		assert.True(t, strings.Index(err.Error(), v) >= 0, "Error message should contains:"+v)
	}
}

// userValidatorTestUnit automation unit for user validator test
type userValidatorTestUnit struct {
	Userr       interface{}
	ErrorStrigs []string
	UsePassword bool
}

// TestUserValidatorPasswordErrors auto test all possible errors for userValidator
func TestUserValidatorPasswordErrors(t *testing.T) {
	m := []userValidatorTestUnit{
		{model.User{Email: "a@b.ch", Firstname: "first", Lastname: "last"}, []string{"password"}, true},
		{model.User{Email: "a@b.ch", Firstname: "first"}, []string{"password", "Lastname"}, true},
		{model.User{Email: "a@b.ch"}, []string{"password", "Lastname", "Firstname"}, true},
		{model.User{}, []string{"password", "Lastname", "Firstname", "email"}, true},

		{model.User{Email: "a@b.ch", Firstname: "first"}, []string{"Lastname"}, false},
		{model.User{Email: "a@b.ch"}, []string{"Lastname", "Firstname"}, false},
		{model.User{}, []string{"Lastname", "Firstname", "email"}, false},
	}

	for _, v := range m {
		testUserValidatorErrors(t, v.Userr.(model.User), v.ErrorStrigs, v.UsePassword)
	}
}

// TestConfirmUserValidator
func TestConfirmUserValidatorMissingEmail(t *testing.T) {
	m := make(map[string][]string)
	m["token"] = []string{"helo"}

	u, err := ConfirmUserValidator(m)

	assert.Nil(t, u)
	assert.True(t, strings.Index(err.Error(), "email") >= 0, "Error message should contains: email")
}

func TestConfirmUserValidatorMissingToken(t *testing.T) {
	m := make(map[string][]string)
	m["email"] = []string{"helo"}

	u, err := ConfirmUserValidator(m)

	assert.Nil(t, u)
	assert.True(t, strings.Index(err.Error(), "token") >= 0, "Error message should conains: token")
}

func TestConfirmUserValidatorMissingTokenEmail(t *testing.T) {
	m := make(map[string][]string)

	u, err := ConfirmUserValidator(m)

	assert.Nil(t, u)
	assert.True(t, strings.Index(err.Error(), "email") >= 0, "Error message should contains: email")
	assert.True(t, strings.Index(err.Error(), "token") >= 0, "Error message should conains: token")
}

func TestConfirmUserValidatorDBerror(t *testing.T) {
	config.TestInit()
	m := make(map[string][]string)
	m["email"] = []string{"a@b.ch"}
	m["token"] = []string{"123qwe"}

	dbResponse := []map[string]interface{}{{"id": 0}}
	mocket.Catcher.Reset().NewMock().WithArgs("a@b.ch", "123qwe").WithReply(dbResponse)

	u, err := ConfirmUserValidator(m)

	assert.Nil(t, u)
	assert.True(t, strings.Index(err.Error(), "Not valid request") >= 0, "Error message should contain: Not valid request")
}

func TestConfirmUserValidator(t *testing.T) {
	config.TestInit()
	m := make(map[string][]string)
	m["email"] = []string{"a@b.ch"}
	m["token"] = []string{"123qwe"}

	expectedUser := model.User{Base: model.Base{ID: 1}, Email: "a@b.c"}
	dbResponse := []map[string]interface{}{{"id": expectedUser.ID, "email": expectedUser.Email}}
	mocket.Catcher.Reset().NewMock().WithArgs("a@b.ch", "123qwe").WithReply(dbResponse)
	actualUser, err := ConfirmUserValidator(m)

	assert.Nil(t, err)
	assert.Equal(t, &expectedUser, actualUser)
}

// testUserValidatorErrors, automatically test error for validator (automated)
func testPasswordRecoveryValidator(t *testing.T, unit userValidatorTestUnit) {
	c := tu.GetContext()

	req, err := tu.GetRequestPost(unit.Userr, "/")
	assert.Nil(t, err)

	c.Request = req

	actualUser, err := PasswordRecoveryValidator(c)

	assert.Nil(t, actualUser)

	for _, v := range unit.ErrorStrigs {
		assert.True(t, strings.Index(err.Error(), v) >= 0, "Error message should contains:"+v)
	}
}

// TestUserValidatorPasswordErrors auto test all possible errors for userValidator
func TestPasswordRecoveryValidatorErrors(t *testing.T) {
	m := []userValidatorTestUnit{
		{PasswordRecovery{Email: "email", Token: "token"}, []string{"new_password"}, true},
		{PasswordRecovery{Email: "email", NewPassword: "new_password"}, []string{"token"}, true},
		{PasswordRecovery{NewPassword: "new_password", Token: "token"}, []string{"email"}, true},
		{PasswordRecovery{NewPassword: "new_password"}, []string{"email", "token"}, true},
		{PasswordRecovery{Email: "email"}, []string{"new_password", "token"}, true},
		{PasswordRecovery{Token: "token"}, []string{"new_password", "email"}, true},
		{PasswordRecovery{}, []string{"new_password", "email", "token"}, true},
	}

	for _, v := range m {
		testPasswordRecoveryValidator(t, v)
	}
}

// TestPasswordRecoveryValidatorDBerrors test the function with query error
func TestPasswordRecoveryValidatorDBerrors(t *testing.T) {
	config.TestInit()

	c := tu.GetContext()

	pr := PasswordRecovery{Email: "a@b.ch", Token: "Token", NewPassword: "new_password"}

	req, err := tu.GetRequestPost(pr, "/")
	assert.Nil(t, err)

	c.Request = req

	expectedUser := model.User{Base: model.Base{ID: 0}}
	dbResponse := []map[string]interface{}{{"id": expectedUser.ID}}
	mocket.Catcher.Reset().NewMock().WithArgs(pr.Email, pr.Token).WithReply(dbResponse)

	actualUser, err := PasswordRecoveryValidator(c)
	expectedUser.Password = pr.NewPassword

	assert.Nil(t, actualUser)
	assert.True(t, strings.Index(err.Error(), config.SUserPasswordRecoveryError) >= 0, "Error message should cotains: "+config.SUserPasswordRecoveryError)
}

// TestPasswordRecoveryValidatorDB test the function with query, no errors
func TestPasswordRecoveryValidatorDB(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.TestInit()

	c := tu.GetContext()

	pr := PasswordRecovery{Email: "a@b.ch", Token: "Token", NewPassword: "new_password"}

	req, err := tu.GetRequestPost(pr, "/")
	assert.Nil(t, err)

	c.Request = req

	expectedUser := model.User{Base: model.Base{ID: 1}, Email: "a@b.c"}
	dbResponse := []map[string]interface{}{{"id": expectedUser.ID, "email": expectedUser.Email}}
	mocket.Catcher.Reset().NewMock().WithArgs(pr.Email, pr.Token).WithReply(dbResponse)

	actualUser, err := PasswordRecoveryValidator(c)
	expectedUser.Password = pr.NewPassword

	assert.Nil(t, err)
	assert.Equal(t, &expectedUser, actualUser)
}
