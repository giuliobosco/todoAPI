package services

import (
	"testing"

	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/model"
)

// TestGetTaskByUser test the function with mock db
func TestGetTaskByUser(t *testing.T) {
	config.TestInit()

	user := model.User{Base: model.Base{ID: 10}}
	expectedTasks := []model.Task{
		{Title: "Tit1"},
		{Title: "Tit2"},
	}

	dbResponse := []map[string]interface{}{
		{"title": expectedTasks[0].Title},
		{"title": expectedTasks[1].Title},
	}
	mocket.Catcher.Reset().NewMock().WithQuery("SELECT").WithReply(dbResponse)
	actualTasks := GetTasksByUser(&user)

	assert.Equal(t, expectedTasks, actualTasks)
}

// TestUpdateTask checks all query are executed
func TestUpdateTask(t *testing.T) {
	config.TestInit()

	o := model.Task{Base: model.Base{ID: 10}, Title: "title", Description: "desciption"}
	n := model.Task{Title: "new_title", Description: "new_description", Done: true}

	mocket.Catcher.Reset().NewMock().WithQuery(`UPDATE "tasks" SET "title"`)
	mocket.Catcher.NewMock().WithQuery(`UPDATE "tasks" SET "description"`)
	mocket.Catcher.NewMock().WithQuery(`UPDATE "tasks" SET "done"`)

	UpdateTask(&o, &n)

	for _, v := range mocket.Catcher.Mocks {
		assert.True(t, v.Triggered)
	}
}
