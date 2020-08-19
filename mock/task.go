package mock

import (
	"github.com/giuliobosco/todoAPI/model"
	"github.com/giuliobosco/todoAPI/testutils"
)

// GetMockTaskID0 build a mock task with ID 0
func GetMockTaskID0() model.Task {
	return model.Task{
		Title:       testutils.RandomString12(),
		Description: testutils.RandomString12(),
		Done:        false,
	}
}

// GetMockTask build a mock task
func GetMockTask() model.Task {
	t := GetMockTaskID0()
	t.ID = testutils.RandomUintNo0()

	return t
}

// GetMockTasksID0 build an array of mock tasks with ID 0
func GetMockTasksID0() []model.Task {
	return []model.Task{
		GetMockTaskID0(),
		GetMockTaskID0(),
	}
}

// GetMockTasks build an array of mock tasks
func GetMockTasks() []model.Task {
	a := GetMockTasksID0()
	a[0].ID = testutils.RandomUintNo0()
	a[1].ID = testutils.RandomUintNo0()

	return a
}

// GetMapByTask convert task to task map
func GetMapByTask(t model.Task) map[string]interface{} {
	return map[string]interface{}{
		"id":          t.ID,
		"title":       t.Title,
		"description": t.Description,
		"done":        t.Done,
	}
}

// GetMapArrayByTask convert task to array of task map
func GetMapArrayByTask(t model.Task) []map[string]interface{} {
	return []map[string]interface{}{GetMapByTask(t)}
}

// GetMapArrayByTasks convert an array of tasks in
func GetMapArrayByTasks(tasks []model.Task) []map[string]interface{} {
	var m []map[string]interface{}

	for _, v := range tasks {
		m = append(m, GetMapByTask(v))
	}

	return m
}
