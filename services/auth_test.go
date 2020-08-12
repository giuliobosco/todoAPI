package services

import (
	"testing"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/mock"
	"github.com/giuliobosco/todoAPI/model"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	mocket.Catcher.Logging = true
	config.TestInit()

	expectedUser := mock.GetMockUser()
	commonReply := mock.GetMapByUser(expectedUser)
	mocket.Catcher.Reset().NewMock().WithQuery("SELECT").WithReply(commonReply)

	var actualUser model.User
	actualUser = GetUserByID(expectedUser.ID)

	assert.Equal(t, expectedUser, actualUser)
}

func TestGetUserByIDNotFound(t *testing.T) {
	config.TestInit()

	commonReply := []map[string]interface{}{{}}
	mocket.Catcher.Reset().NewMock().WithQuery("SELECT").WithReply(commonReply)

	var actualUser model.User
	actualUser = GetUserByID(1)

	assert.Equal(t, 0, int(actualUser.ID))
}

func TestGetUserByEmail(t *testing.T) {
	config.TestInit()

	expectedUser := mock.GetMockUser()
	commonReply := mock.GetMapByUser(expectedUser)
	mocket.Catcher.Reset().NewMock().WithQuery("SELECT").WithReply(commonReply)

	var actualUser model.User
	actualUser = GetUserByEmail(expectedUser.Email)

	assert.Equal(t, expectedUser, actualUser)
}

func TestGetUserByEmailNotFound(t *testing.T) {
	config.TestInit()

	commonReply := []map[string]interface{}{{}}
	mocket.Catcher.Reset().NewMock().WithQuery("SELECT").WithReply(commonReply)

	var actualUser model.User
	actualUser = GetUserByEmail(mock.GetMockUser().Email)

	assert.Equal(t, 0, int(actualUser.ID))
}
