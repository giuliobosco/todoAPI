package services

import (
	"testing"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/model"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	mocket.Catcher.Logging = true
	config.TestInit()

	commonReply := []map[string]interface{}{{"id": 1}}
	mocket.Catcher.Reset().NewMock().WithQuery("SELECT").WithReply(commonReply)

	var user model.User
	user = GetUserByID(1)

	assert.Equal(t, 1, int(user.ID))
}

func TestGetUserByIDNotFound(t *testing.T) {
	config.TestInit()

	commonReply := []map[string]interface{}{{}}
	mocket.Catcher.Reset().NewMock().WithQuery("SELECT").WithReply(commonReply)

	var user model.User
	user = GetUserByID(1)

	assert.Equal(t, 0, int(user.ID))
}
