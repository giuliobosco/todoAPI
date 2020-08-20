package controller

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/giuliobosco/todoAPI/utils"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/mock"
	"github.com/giuliobosco/todoAPI/tu"

	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
)

// ################# TESTS
// FetchAllTask()

// TestFetchUserAuthenticated test with not authenticated
func TestFetchUserUnauthenticated(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()
	mocket.Catcher.Reset()

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)

	c.Request = req

	FetchUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), config.SUserInvalid)
}

// TestFetchUser test with user
func TestFetchUser(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)

	c.Request = req

	FetchUser(c)

	jsonUser, err := json.Marshal(u)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), string(jsonUser))
}

// ################# TESTS
// UpdateUser()

// TestUpdateUserUnauthorized without authentication
func TestUpdateUserUnauhtorized(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()
	mocket.Catcher.Reset()

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)

	c.Request = req

	UpdateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), config.SUserInvalid)
}

// TestUserUpdateEmpty with empty data
func TestUserUpdateEmpty(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)

	c.Request = req

	UpdateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestUserUpdateNoPassword test user if in db has no password
// means is oauth user, so cannot change email address
func TestUserUpdateNoPassword(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	newUser := mock.GetMockUser(false)

	req, err := tu.GetRequestPost(newUser, "/")
	assert.Nil(t, err)

	c.Request = req

	UpdateUser(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), config.SUserFailUpdate)
}

// TestUserUpdateEmailExists test if user tries an email who already exists
func TestUserUpdateEmailExists(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(true)
	mock.ConfigClaims(c, u)

	newUser := mock.GetMockUser(true)

	req, err := tu.GetRequestPost(newUser, "/")
	assert.Nil(t, err)

	dbResponse := mock.GetMapArrayByUser(newUser)
	mocket.Catcher.NewMock().WithArgs(newUser.Email).WithReply(dbResponse)

	c.Request = req

	UpdateUser(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), config.SUserEmailAlreadyExists)
}

// TestUserUpdateEmail test update with email
func TestUserUpdateEmail(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(true)
	mock.ConfigClaims(c, u)

	newUser := mock.GetMockUser(true)

	req, err := tu.GetRequestPost(newUser, "/")
	assert.Nil(t, err)

	newUser.ID = 0
	dbResponse := mock.GetMapArrayByUser(newUser)
	mocket.Catcher.Reset().NewMock().WithArgs(newUser.Email).WithReply(dbResponse)
	dbResponse = mock.GetMapArrayByUser(u)
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND ((id =`).WithReply(dbResponse)

	c.Request = req

	UpdateUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestUserUpdate without email change
func TestUserUpdate(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(true)
	mock.ConfigClaims(c, u)

	newUser := u
	newUser.Firstname = tu.RandomString12()

	req, err := tu.GetRequestPost(newUser, "/")
	assert.Nil(t, err)

	c.Request = req

	UpdateUser(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

// ################# TESTS
// DeleteUser()

// TestDeleteUserUnauhtenticated tests func without authentication
func TestDeleteUserUnauthenticated(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()
	mocket.Catcher.Reset()

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)

	c.Request = req

	DeleteUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), config.SUserInvalid)
}

// TestDeleteUserNoData tests func without req data
func TestDeleteUserNoData(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(true)
	mock.ConfigClaims(c, u)

	req, err := tu.GetRequestPost(u, "/")
	assert.Nil(t, err)

	c.Request = req

	DeleteUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), config.SWrongPassword)
}

// TestDeleteUser test delete user
func TestDeleteUser(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(true)

	req, err := tu.GetRequestPost(u, "/")
	assert.Nil(t, err)

	u.Password, err = utils.PasswordHash(u.Password)
	assert.Nil(t, err)

	mock.ConfigClaims(c, u)

	c.Request = req

	DeleteUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
