package services

import (
	"strings"
	"testing"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/model"

	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
)

// TestVerifyUserEmailTokenError tests the query for check email and token with errors
func TestVerifyUserEmailTokenError(t *testing.T) {
	config.TestInit()
	expectedUser := model.User{Base: model.Base{ID: 0}, Email: "a@b.c"}
	dbResponse := []map[string]interface{}{{"id": expectedUser.ID, "email": expectedUser.Email}}
	mocket.Catcher.Reset().NewMock().WithArgs("a@b.ch", "123qwe").WithReply(dbResponse)

	actualUser, err := VerifyUserEmailToken("a@b.ch", "123qwe")

	assert.Nil(t, actualUser)
	assert.True(t, strings.Index(err.Error(), "Not valid request") >= 0, "Error message should contain: Not valid request")
}

// TestVerifyUserEmailTokenError tests the query for check email and token
func TestVerifyUserEmailToken(t *testing.T) {
	config.TestInit()
	expectedUser := model.User{Base: model.Base{ID: 1}, Email: "a@b.c"}
	dbResponse := []map[string]interface{}{{"id": expectedUser.ID, "email": expectedUser.Email}}
	mocket.Catcher.Reset().NewMock().WithArgs("a@b.ch", "123qwe").WithReply(dbResponse)

	actualUser, err := VerifyUserEmailToken("a@b.ch", "123qwe")

	assert.Nil(t, err)
	assert.Equal(t, &expectedUser, actualUser)
}
