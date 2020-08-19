package controller

import (
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
// RegisterEndPoint()

// TestRegisterEndPointMissingParam test parameter validation
func TestRegisterEndPointMissingParam(t *testing.T) {
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)

	req, err := tu.GetRequestPost(u, "/")
	assert.Nil(t, err)

	c.Request = req

	RegisterEndPoint(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestRegisterEndPointEmailConflict test email check on db
func TestRegisterEndPoinEmailConflict(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(true)

	req, err := tu.GetRequestPost(u, "/")
	assert.Nil(t, err)

	c.Request = req

	dbResponse := mock.GetMapArrayByUser(u)
	mocket.Catcher.Reset().NewMock().WithArgs(u.Email).WithReply(dbResponse)

	RegisterEndPoint(c)

	assert.Equal(t, http.StatusConflict, w.Code)
}

// TestRegisterEndPoint test user registration
func TestRegisterEndPoint(t *testing.T) {
	config.TestInit()
	utils.SetTesting(true)
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUserID0(true)

	req, err := tu.GetRequestPost(u, "/")
	assert.Nil(t, err)

	c.Request = req

	dbResponse := mock.GetMapArrayByUser(mock.GetMockUserID0(false))
	mocket.Catcher.Reset().NewMock().WithQuery("SELECT").WithArgs(u.Email).WithReply(dbResponse)
	mocket.Catcher.NewMock().WithQuery("INSERT")

	RegisterEndPoint(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.True(t, mocket.Catcher.Mocks[0].Triggered)
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
}

// TestConfirmUserNoParams test user confirmation without params
func TestConfirmUserNoParams(t *testing.T) {
	w, c := tu.GetRecorderContext()

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)

	c.Request = req

	ConfirmUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestConfirmUser(t *testing.T) {
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	u.Active = false
	u.VerifyToken = "123qwe"
	p := "/a?email=" + u.Email + "&token=" + u.VerifyToken

	req, err := tu.GetRequestPost(nil, p)
	assert.Nil(t, err)

	c.Request = req

	dbResponse := mock.GetMapArrayByUser(u)
	mocket.Catcher.Reset().NewMock().WithQuery("SELECT").WithArgs(u.Email, u.VerifyToken).WithReply(dbResponse)
	mocket.Catcher.NewMock().WithQuery("UPDATE")

	ConfirmUser(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, mocket.Catcher.Mocks[0].Triggered)
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
	assert.False(t, mocket.Catcher.Mocks[1].Once)
}
