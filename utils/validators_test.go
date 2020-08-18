package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/model"

	"github.com/gin-gonic/gin"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
)

func getTestUser() model.User {
	return model.User{Email: "a@b.ch", Password: "123qwe", Lastname: "lastname", Firstname: "firstname"}
}

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
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	expectedUser := getTestUser()

	jsonStr, _ := json.Marshal(expectedUser)
	var jsonBytes = []byte(jsonStr)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonBytes))

	c.Request = req

	actualUser, err := UserValidator(c, true)

	assert.Equal(t, &expectedUser, actualUser)
	assert.Nil(t, err)
}

// TestUserValidatorNoPassword, user without password
func TestUserValidatorNoPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	expectedUser := getTestUser()
	expectedUser.Password = ""

	jsonStr, _ := json.Marshal(expectedUser)
	var jsonBytes = []byte(jsonStr)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonBytes))

	c.Request = req

	userActual, err := UserValidator(c, false)

	assert.Equal(t, &expectedUser, userActual)
	assert.Nil(t, err)
}

// testUserValidatorErrors, automatically test error for validator (automated)
func testUserValidatorErrors(t *testing.T, u model.User, missing []string, usePassword bool) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	jsonStr, err := json.Marshal(u)
	assert.Nil(t, err)

	jsonBytes := []byte(jsonStr)

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonBytes))
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
	Userr       model.User
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
		testUserValidatorErrors(t, v.Userr, v.ErrorStrigs, v.UsePassword)
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
