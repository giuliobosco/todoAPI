package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/mock"
	"github.com/giuliobosco/todoAPI/model"
	"github.com/giuliobosco/todoAPI/testutils"
	"github.com/giuliobosco/todoAPI/utils"

	jwtapple2 "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
)

// ################# TESTS
// emailAuthenticator()

// TestEmailAuthenticatorNoData test withoud request data
func TestEmailAuthenticatorNoData(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.TestInit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	//loginVals := model.Base{} //mock.GetLoginVals(mock.GetMockUser(true))
	//jsonStr, err := json.Marshal(loginVals)
	//assert.Nil(t, err)

	//jsonBytes := []byte(jsonStr)

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)

	c.Request = req

	u, err := emailAuthenticator(c)

	assert.Equal(t, "", u)
	assert.Equal(t, jwtapple2.ErrMissingLoginValues.Error(), string(err.Error()))
}

// TestEmailAuthenticatorWrongEmail test with wrong email
func TestEmailAuthenticatorWrongEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.TestInit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	loginVals := mock.GetLoginVals(mock.GetMockUser(true))
	jsonStr, err := json.Marshal(loginVals)
	assert.Nil(t, err)

	jsonBytes := []byte(jsonStr)

	req, err := http.NewRequest("POST", "/", bytes.NewReader(jsonBytes))
	assert.Nil(t, err)

	c.Request = req

	dbResponse := mock.GetMapArrayByUser(mock.GetMockUserID0(false))
	mocket.Catcher.Reset().NewMock().WithArgs(loginVals.Email).WithReply(dbResponse)

	u, err := emailAuthenticator(c)

	assert.Nil(t, u)
	assert.Equal(t, jwtapple2.ErrFailedAuthentication.Error(), err.Error())
}

// TestEmailAuthenticatorNotActive with not active user
func TestEmailAuthenticatorNotActive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.TestInit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	expectedUser := mock.GetMockUser(true)
	expectedUser.Active = false
	loginVals := mock.GetLoginVals(expectedUser)
	jsonStr, err := json.Marshal(loginVals)
	assert.Nil(t, err)

	jsonBytes := []byte(jsonStr)

	req, err := http.NewRequest("POST", "/", bytes.NewReader(jsonBytes))
	assert.Nil(t, err)

	c.Request = req
	dbResponse := mock.GetMapArrayByUser(expectedUser)
	mocket.Catcher.Reset().NewMock().WithArgs(loginVals.Email).WithReply(dbResponse)

	u, err := emailAuthenticator(c)

	assert.Nil(t, u)
	assert.Equal(t, config.SUserNotConfirmed, err.Error())
}

// TestEmailAuthenticatorWrongPassword with wrong password
func TestEmailAuthenticatorWrongPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.TestInit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	expectedUser := mock.GetMockUser(true)
	expectedUser.Active = true
	loginVals := mock.GetLoginVals(expectedUser)
	jsonStr, err := json.Marshal(loginVals)
	assert.Nil(t, err)

	jsonBytes := []byte(jsonStr)

	req, err := http.NewRequest("POST", "/", bytes.NewReader(jsonBytes))
	assert.Nil(t, err)

	c.Request = req
	dbResponse := mock.GetMapArrayByUser(expectedUser)
	mocket.Catcher.Reset().NewMock().WithArgs(loginVals.Email).WithReply(dbResponse)

	u, err := emailAuthenticator(c)

	assert.Nil(t, u)
	assert.Equal(t, jwtapple2.ErrFailedAuthentication.Error(), err.Error())
}

// TestEmailAuthenticatorNoToken without token
func TestEmailAuthenticatorNoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.TestInit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	expectedUser := mock.GetMockUser(true)
	expectedUser.Active = true
	expectedUser.VerifyToken = ""
	loginVals := mock.GetLoginVals(expectedUser)
	var e error
	expectedUser.Password, e = utils.PasswordHash(expectedUser.Password)
	assert.Nil(t, e)

	jsonStr, err := json.Marshal(loginVals)
	assert.Nil(t, err)

	jsonBytes := []byte(jsonStr)

	req, err := http.NewRequest("POST", "/", bytes.NewReader(jsonBytes))
	assert.Nil(t, err)

	c.Request = req
	dbResponse := mock.GetMapArrayByUser(expectedUser)
	mocket.Catcher.Reset().NewMock().WithArgs(loginVals.Email).WithReply(dbResponse)
	mocket.Catcher.NewMock().WithArgs(`UPDATE "users" SET "verify_token"`)

	u, err := emailAuthenticator(c)

	assert.Nil(t, err)
	assert.Equal(t, &expectedUser, u)
	assert.False(t, mocket.Catcher.Mocks[1].Triggered)
}

// TestEmailAuthenticatorWithToken with token
func TestEmailAuthenticatorWithToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.TestInit()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	expectedUser := mock.GetMockUser(true)
	expectedUser.Active = true
	expectedUser.VerifyToken = testutils.RandomString12()
	loginVals := mock.GetLoginVals(expectedUser)
	var e error
	expectedUser.Password, e = utils.PasswordHash(expectedUser.Password)
	assert.Nil(t, e)

	jsonStr, err := json.Marshal(loginVals)
	assert.Nil(t, err)

	jsonBytes := []byte(jsonStr)

	req, err := http.NewRequest("POST", "/", bytes.NewReader(jsonBytes))
	assert.Nil(t, err)

	c.Request = req
	dbResponse := mock.GetMapArrayByUser(expectedUser)
	mocket.Catcher.Reset().NewMock().WithArgs(loginVals.Email).WithReply(dbResponse)
	mocket.Catcher.NewMock().WithQuery(`UPDATE "users" SET "updated_at" = ?, "verify_token" = ?`)

	expectedUser.VerifyToken = ""
	u, err := emailAuthenticator(c)

	user := u.(*model.User)

	assert.Nil(t, err)
	assert.Equal(t, expectedUser.Email, user.Email)
	assert.Equal(t, expectedUser.VerifyToken, user.VerifyToken)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
}
