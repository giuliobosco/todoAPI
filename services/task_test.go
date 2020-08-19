package services

import (
	"testing"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/mock"
	"github.com/giuliobosco/todoAPI/model"

	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
)

// TestGetTaskByUser test the function with mock db
func TestGetTaskByUser(t *testing.T) {
	config.TestInit()

	user := model.User{Base: model.Base{ID: 10}}
	expectedTasks := mock.GetMockTasks()

	dbResponse := mock.GetMapArrayByTasks(expectedTasks)
	mocket.Catcher.Reset().NewMock().WithQuery("SELECT").WithReply(dbResponse)
	actualTasks := GetTasksByUser(&user)

	assert.Equal(t, expectedTasks, actualTasks)
}

// TestUpdateTask checks all query are executed
func TestUpdateTask(t *testing.T) {
	config.TestInit()

	o := mock.GetMockTask()
	n := mock.GetMockTask()
	n.Title = "new_title"
	n.Description = "new_description"
	n.Done = !n.Done

	mocket.Catcher.Reset().NewMock().WithQuery(`UPDATE "tasks" SET "title"`)
	mocket.Catcher.NewMock().WithQuery(`UPDATE "tasks" SET "description"`)
	mocket.Catcher.NewMock().WithQuery(`UPDATE "tasks" SET "done"`)

	UpdateTask(&o, &n)

	for _, v := range mocket.Catcher.Mocks {
		assert.True(t, v.Triggered)
	}
}
