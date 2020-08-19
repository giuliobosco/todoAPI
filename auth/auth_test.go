package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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
// payload()

// TestPayloadEmpty payload function empty data
func TestPayloadEmpty(t *testing.T) {
	d := ""

	expected := jwtapple2.MapClaims{}
	actual := payload(&d)

	assert.Equal(t, expected, actual)
}

// TestPayload payload function with user
func TestPayload(t *testing.T) {
	d := mock.GetMockUser(false)

	expected := jwtapple2.MapClaims{"id": d.ID}
	actual := payload(&d)

	assert.Equal(t, expected, actual)
}

// ################# TESTS
// identityHandler()

// TestIdentityHandlerNoId test identityHandler function without id
func TestIdentityHandlerNoId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	i := identityHandler(c)

	assert.Nil(t, i)
}

// TestIdentityHandler test identityHandler function with id
func TestIdentityHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	config.TestInit()
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	expectedUser := mock.GetMockUser(false)
	claims := jwtapple2.MapClaims{"id": 1.0}
	c.Set("JWT_PAYLOAD", claims)

	dbResponse := mock.GetMapArrayByUser(expectedUser)
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "users"`).WithReply(dbResponse)
	actualUser := identityHandler(c)

	assert.Equal(t, expectedUser, actualUser)
}

// ################# TESTS
// authenticator()

// TestAuthenticatorNoType tests authenticator without the auth type
func TestAuthenticatorNoType(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	req, err := http.NewRequest("POST", "/", nil)
	assert.Nil(t, err)

	c.Request = req

	i, err := authenticator(c)

	assert.Nil(t, i)
	assert.Equal(t, config.SMissingAuthType, err.Error())
}

// TestAuthenticatorEmail tests authenticator via email
func TestAuthenticatorEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	req, err := http.NewRequest("POST", "/?type=email", nil)
	assert.Nil(t, err)

	c.Request = req

	i, err := authenticator(c)

	assert.Empty(t, i)
	assert.Equal(t, jwtapple2.ErrMissingLoginValues.Error(), err.Error())
}

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

// ################# TESTS
// emailAuthenticator()

// TestAuthorizatorNoData test authorizator() without data
func TestAuthorizatorNoData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	d := ""

	b := authorizator(d, c)

	assert.False(t, b)
}

// TestAuthorizatorID0 test authorizator() with user id 0
func TestAuthorizatorID0(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	d := mock.GetMockUserID0(false)
	b := authorizator(d, c)

	assert.False(t, b)
}

// TestAuthorizatorNotActive test authorizator() with not active user
func TestAuthorizatorNotActive(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	d := mock.GetMockUser(false)
	d.Active = false
	b := authorizator(d, c)

	assert.False(t, b)
}

// TestAuthorizator test authorizator() with working user
func TestAuthorizator(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	d := mock.GetMockUser(false)
	d.Active = true
	b := authorizator(d, c)

	assert.True(t, b)
}

// ################# TESTS
// unauthorized()

// TestUnauthorized tests unauthorized() with differentes codes
func TestUnauthorized(t *testing.T) {
	codes := []int{400, 401}
	for _, v := range codes {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		code := v
		message := testutils.RandomString12()

		unauthorized(c, code, message)

		assert.Equal(t, code, w.Code)
		assert.True(t, strings.Index(w.Body.String(), message) > 0, "Response should contains:"+message)
	}
}

// ################# TESTS
// loginResponse()

// TestLoginResponse tests unauthorized() with differentes codes
func TestLoginResponse(t *testing.T) {
	codes := []int{400, 401}
	for _, v := range codes {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		code := v
		token := testutils.RandomString12()
		expire := time.Now()

		loginResponse(c, code, token, expire)

		assert.Equal(t, code, w.Code)
		assert.True(t, strings.Index(w.Body.String(), token) > 0, "Response should contains:"+token)
	}
}
