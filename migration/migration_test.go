package migration

import (
	"testing"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/stretchr/testify/assert"

	mocket "github.com/selvatico/go-mocket"
)

// TestMigrate tests the migration of all the tables
func TestMigrate(t *testing.T) {
	config.TestInit()
	mocket.Catcher.Reset().NewMock().WithQuery(`CREATE TABLE "tasks"`)
	mocket.Catcher.NewMock().WithQuery(`CREATE TABLE "users"`)

	Migrate(config.GetDB())

	assert.True(t, mocket.Catcher.Mocks[0].Triggered)
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
}
