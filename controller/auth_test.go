package controller

import (
	"net/http"
	"testing"

	"github.com/giuliobosco/todoAPI/model"
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

// ################# TESTS
// ConfirmUser()

// TestConfirmUserNoParams test user confirmation without params
func TestConfirmUserNoParams(t *testing.T) {
	w, c := tu.GetRecorderContext()

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)

	c.Request = req

	ConfirmUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestConfirmUser test user confirmation correctly
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
	assert.Contains(t, w.Body.String(), config.SUserConfirmed)
	assert.True(t, mocket.Catcher.Mocks[0].Triggered)
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
	assert.False(t, mocket.Catcher.Mocks[1].Once)
}

// ################# TESTS
// getUSerByEmailParam()

// TestGetUserByEmailParamNoEmail check the function without email
func TestGetUserByEmailParamNoEmail(t *testing.T) {
	c := tu.GetContext()

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)

	c.Request = req

	u, err := getUserByEmailParam(c)

	assert.Nil(t, u)
	assert.Equal(t, config.SMissingEmail, err.Error())
}

// TestGetUserByEmailParamNoUser check the function with no user in db
func TestGetUserByEmailParamNoUser(t *testing.T) {
	c := tu.GetContext()

	req, err := tu.GetRequestPost(nil, "/?email=a@b.ch")
	assert.Nil(t, err)

	c.Request = req

	mocket.Catcher.Reset().NewMock().WithArgs("a@b.ch")

	u, err := getUserByEmailParam(c)

	assert.Nil(t, u)
	assert.Equal(t, config.SUserInvalid, err.Error())
}

// TestGetUserByEmailParamNoUser check the function with user
func TestGetUserByEmailParam(t *testing.T) {
	c := tu.GetContext()

	expectedUser := mock.GetMockUser(false)
	req, err := tu.GetRequestPost(nil, "/?email="+expectedUser.Email)
	assert.Nil(t, err)

	c.Request = req

	dbResponse := mock.GetMapArrayByUser(expectedUser)
	mocket.Catcher.Reset().NewMock().WithArgs(expectedUser.Email).WithReply(dbResponse)

	u, err := getUserByEmailParam(c)

	assert.Equal(t, &expectedUser, u)
	assert.Nil(t, err)
}

// ################# TESTS
// SendUserConfirmAgain()

// TestSendUserConfirmAgainNoEmail test func with no email
func TestSendUserConfirmAgainNoEmail(t *testing.T) {
	w, c := tu.GetRecorderContext()

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)

	c.Request = req

	SendUserConfirmAgain(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), config.SMissingEmail)
}

// TestSendUserConfirmAgain test func
func TestSendUserConfirmAgain(t *testing.T) {
	w, c := tu.GetRecorderContext()

	expectedUser := mock.GetMockUser(false)
	req, err := tu.GetRequestPost(nil, "/?email="+expectedUser.Email)
	assert.Nil(t, err)

	c.Request = req

	dbResponse := mock.GetMapArrayByUser(expectedUser)
	mocket.Catcher.Reset().NewMock().WithArgs(expectedUser.Email).WithReply(dbResponse)
	mocket.Catcher.NewMock().WithQuery("UPDATE")

	SendUserConfirmAgain(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), config.SUserSentConfirmationMailAgain)
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
}

// ################# TESTS
// RequestPasswordRecovery()

// TestRequestPasswordRecoveryNoEmail tests the func without the email parameter
func TestRequestPasswordRecoveryNoEmail(t *testing.T) {
	w, c := tu.GetRecorderContext()

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)

	c.Request = req

	RequestPasswordRecovery(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), config.SMissingEmail)
}

// TestRequestPasswordRecoveryNotActive tests the func with user not active
func TestRequestPasswordRecoveryNotActive(t *testing.T) {
	w, c := tu.GetRecorderContext()

	expectedUser := mock.GetMockUser(false)
	expectedUser.Active = false
	req, err := tu.GetRequestPost(nil, "/?email="+expectedUser.Email)
	assert.Nil(t, err)

	c.Request = req

	dbResponse := mock.GetMapArrayByUser(expectedUser)
	mocket.Catcher.Reset().NewMock().WithArgs(expectedUser.Email).WithReply(dbResponse)

	RequestPasswordRecovery(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), config.SUserNotConfirmed)
}

// TestRequestPasswordRecovery tests the func with user ok
func TestRequestPasswordRecovery(t *testing.T) {
	w, c := tu.GetRecorderContext()

	expectedUser := mock.GetMockUser(false)
	expectedUser.Active = true
	req, err := tu.GetRequestPost(nil, "/?email="+expectedUser.Email)
	assert.Nil(t, err)

	c.Request = req

	dbResponse := mock.GetMapArrayByUser(expectedUser)
	mocket.Catcher.Reset().NewMock().WithArgs(expectedUser.Email).WithReply(dbResponse)
	mocket.Catcher.NewMock().WithQuery("UPDATE")

	RequestPasswordRecovery(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), config.SUserPasswordRecoveryMailSent)
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
}

// ################# TESTS
// ExecutePasswordRecovery()

// TestExecutePasswordRecoveryNoUser test func without user
func TestExecutePasswordRecoveryNoUser(t *testing.T) {
	w, c := tu.GetRecorderContext()

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)

	c.Request = req

	ExecutePasswordRecovery(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestExecutePasswordRecovery test func with ok params
func TestExecutePasswordRecovery(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	pr := utils.PasswordRecovery{
		Email:       tu.RandomEmail(),
		Token:       tu.RandomString12(),
		NewPassword: tu.RandomString12(),
	}
	expectedUser := model.User{
		Email:       pr.Email,
		VerifyToken: pr.Token,
		Password:    tu.RandomString12(),
	}
	expectedUser.ID = tu.RandomUintNo0()

	req, err := tu.GetRequestPost(pr, "/")
	assert.Nil(t, err)

	c.Request = req

	dbResponse := mock.GetMapArrayByUser(expectedUser)
	mocket.Catcher.Reset().NewMock().WithArgs(pr.Email, pr.Token).WithReply(dbResponse)
	mocket.Catcher.NewMock().WithQuery("UPDATE")

	ExecutePasswordRecovery(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), config.SUserPasswordUpdated)
}

// ################# TESTS
// getUserByContext()

// TestGetUserByContextNoUser test func without db user
func TestGetUserByContextNoUser(t *testing.T) {
	c := tu.GetContext()
	expectedUser := mock.GetMockUserID0(false)
	mock.ConfigClaims(c, expectedUser)

	actualUser, err := getUserByContext(c)

	assert.Nil(t, actualUser)
	assert.Contains(t, config.SUserInvalid, err.Error())
}

// TestGetUserByContext test func with db user
func TestGetUserByContext(t *testing.T) {
	c := tu.GetContext()
	expectedUser := mock.GetMockUser(true)
	mock.ConfigClaims(c, expectedUser)

	actualUser, err := getUserByContext(c)

	assert.Equal(t, expectedUser.Email, actualUser.Email)
	assert.Equal(t, expectedUser.Firstname, actualUser.Firstname)
	assert.Equal(t, expectedUser.Lastname, actualUser.Lastname)
	assert.Equal(t, expectedUser.VerifyToken, actualUser.VerifyToken)
	assert.Equal(t, expectedUser.Active, actualUser.Active)
	assert.Empty(t, actualUser.Password)
	assert.Nil(t, err)
}

// ################# TESTS
// getUserByContext()

// TestGetUserWithPasswordByContextNoUser test func without db user
func TestGetUserWithPasswordByContextNoUser(t *testing.T) {
	c := tu.GetContext()
	expectedUser := mock.GetMockUserID0(false)
	mock.ConfigClaims(c, expectedUser)

	actualUser, err := getUserWithPasswordByContext(c)

	assert.Nil(t, actualUser)
	assert.Contains(t, config.SUserInvalid, err.Error())
}

// TestGetUserWithPasswordByContext test func with db user
func TestGetUserWithPasswordByContext(t *testing.T) {
	c := tu.GetContext()
	expectedUser := mock.GetMockUser(true)
	mock.ConfigClaims(c, expectedUser)

	actualUser, err := getUserWithPasswordByContext(c)

	assert.Equal(t, &expectedUser, actualUser)
	assert.Nil(t, err)
}

// ################# TESTS
// UpdatePassword()

// TestUpdatePasswordNoUser test func with no user in db
func TestUpdatePasswordNoUser(t *testing.T) {
	w, c := tu.GetRecorderContext()
	expectedUser := mock.GetMockUserID0(true)
	mock.ConfigClaims(c, expectedUser)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)

	c.Request = req

	UpdatePassword(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), config.SUserInvalid)
}

// TestUpdatePasswordMissingData test func with missing request data
func TestUpdatePasswordMissingData(t *testing.T) {
	w, c := tu.GetRecorderContext()
	expectedUser := mock.GetMockUser(true)
	mock.ConfigClaims(c, expectedUser)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)

	c.Request = req

	UpdatePassword(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), config.SMissingOldNewPassword)
}

// TestUpdatePasswordFailCompare test func with wrong password
func TestUpdatePasswordFailCompare(t *testing.T) {
	w, c := tu.GetRecorderContext()
	expectedUser := mock.GetMockUser(true)

	up := UpdatePasswordObj{
		OldPassword: expectedUser.Password,
		NewPassword: tu.RandomString12(),
	}

	hash, err := utils.PasswordHash(expectedUser.Password)
	assert.Nil(t, err)
	expectedUser.Password = hash
	mock.ConfigClaims(c, expectedUser)
	req, err := tu.GetRequestPost(up, "/")
	assert.Nil(t, err)

	c.Request = req

	mocket.Catcher.NewMock().WithQuery("UPDATE")

	UpdatePassword(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), config.SUserPasswordUpdated)
	assert.True(t, mocket.Catcher.Mocks[0].Triggered)
}
