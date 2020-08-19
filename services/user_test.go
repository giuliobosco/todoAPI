package services

import (
	"strings"
	"testing"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/mock"

	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
)

// TestVerifyUserEmailTokenError tests the query for check email and token with errors
func TestVerifyUserEmailTokenError(t *testing.T) {
	config.TestInit()
	expectedUser := mock.GetMockUserID0(false)
	dbResponse := mock.GetMapArrayByUser(expectedUser)
	mocket.Catcher.Reset().NewMock().WithArgs("a@b.ch", "123qwe").WithReply(dbResponse)

	actualUser, err := VerifyUserEmailToken("a@b.ch", "123qwe")

	assert.Nil(t, actualUser)
	assert.True(t, strings.Index(err.Error(), "Not valid request") >= 0, "Error message should contain: Not valid request")
}

// TestVerifyUserEmailTokenError tests the query for check email and token
func TestVerifyUserEmailToken(t *testing.T) {
	config.TestInit()
	expectedUser := mock.GetMockUser(false)
	dbResponse := mock.GetMapArrayByUser(expectedUser)
	mocket.Catcher.Reset().NewMock().WithArgs("a@b.ch", "123qwe").WithReply(dbResponse)

	actualUser, err := VerifyUserEmailToken("a@b.ch", "123qwe")

	assert.Nil(t, err)
	assert.Equal(t, &expectedUser, actualUser)
}
