package controller

import (
	"net/http"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/model"
	"github.com/giuliobosco/todoAPI/utils"

	jwtapple2 "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const sMessage string = config.SMessage
const sError string = config.SError
const sData string = config.SData
const sTask string = config.STask

// RegisterEndPoint registration API End Point
func RegisterEndPoint(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}
	if ok, err := utils.UserValidator(user); !ok {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}

	var userCheck model.User
	config.GetDB().First(&userCheck, "email = ?", user.Email)

	if userCheck.ID > 0 {
		c.JSON(http.StatusConflict, gin.H{sMessage: config.SUserExists})
		return
	}

	user.Active = false
	var err error
	user.VerifyToken, err = utils.GenerateRandomStringURLSafe(config.TokenLength)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{sError: config.SUserFailCreation})
	}

	config.GetDB().Save(&user)
	utils.UserConfirmation(user)

	c.JSON(http.StatusCreated, gin.H{sMessage: config.SUserCreated})
}

// CreateTask is the function for create a task
func CreateTask(c *gin.Context) {
	claims := jwtapple2.ExtractClaims(c)

	var user model.User
	config.GetDB().Where("id = ?", claims[config.IdentityKey]).First(&user)

	if user.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{sError: config.SUserInvalid})
		return
	}

	var todo model.Task
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}

	todo.UserID = user.ID
	config.GetDB().Save(&todo)
	c.JSON(http.StatusCreated, gin.H{sMessage: config.STaskCreated, sTask: todo})
}

// FetchAllTask is the function for fetch all tasks
func FetchAllTask(c *gin.Context) {
	claims := jwtapple2.ExtractClaims(c)

	var user model.User
	config.GetDB().Where("id = ?", claims[config.IdentityKey]).First(&user)

	if user.ID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{sError: config.SUserInvalid})
		return
	}

	var todos []model.Task
	config.GetDB().Where("user_id = ?", user.ID).Order("created_at desc").Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{sMessage: config.STaskNotFound, sData: todos})
		return
	}

	c.JSON(http.StatusOK, gin.H{sData: todos})
}

// FetchSingleTask is the function for fetch a single task by id
func FetchSingleTask(c *gin.Context) {
	todoID := c.Param("id")

	if len(todoID) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{sError: config.STaskInvalid})
		return
	}

	var todo model.Task
	config.GetDB().First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{sMessage: config.STaskNotFound})
		return
	}

	c.JSON(http.StatusOK, todo)
}

// UpdateTask is the function for update a task by id
func UpdateTask(c *gin.Context) {
	todoID := c.Param("id")

	if len(todoID) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{sError: config.SUserInvalid})
		return
	}

	var newTodo model.Task
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{sError: err.Error()})
		return
	}

	var todo model.Task
	config.GetDB().First(&todo, todoID)

	if todo.ID <= 0 {
		c.JSON(http.StatusNotFound, gin.H{sMessage: config.STaskNotFound})
		return
	}

	config.GetDB().Model(&todo).Update("title", newTodo.Title)
	config.GetDB().Model(&todo).Update("description", newTodo.Description)
	config.GetDB().Model(&todo).Update("completed", newTodo.Completed)

	config.GetDB().First(&todo, todoID)

	c.JSON(http.StatusOK, gin.H{sMessage: config.STaskUpdated, sTask: todo})
}

// DeleteTask is the function for delete a task by id
func DeleteTask(c *gin.Context) {
	var todo model.Task
	todoID := c.Param("id")

	if len(todoID) <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{sError: config.STaskNotFound})
		return
	}

	config.GetDB().First(&todo, todoID)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{sMessage: config.STaskNotFound})
		return
	}

	config.GetDB().Delete(&todo)
	c.JSON(http.StatusOK, gin.H{sMessage: config.STaskDeleted, sTask: todo})
}
