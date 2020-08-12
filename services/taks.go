package services

import (
	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/model"
)

// GetTasksByUser get the tasks by the id of the user
func GetTasksByUser(user *model.User) []model.Task {
	var todos []model.Task
	config.GetDB().Where("user_id = ?", user.ID).Order("created_at desc").Find(&todos)

	return todos
}

// UpdateTask updates title, description, completed of o from n
func UpdateTask(o, n *model.Task) {
	config.GetDB().Model(o).Update("title", n.Title)
	config.GetDB().Model(o).Update("description", n.Description)
	config.GetDB().Model(o).Update("completed", n.Completed)
}
