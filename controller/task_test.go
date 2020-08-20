package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/mock"
	"github.com/giuliobosco/todoAPI/tu"

	"github.com/gin-gonic/gin"
	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
)

// ################# TESTS
// CreateTask()

// TestCreateTaskNotAuthenticated func with no authenticated user
func TestCreateTaskNotAuthenticated(t *testing.T) {
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUserID0(false)
	dbResponse := mock.GetMapArrayByUser(u)
	mocket.Catcher.Reset().NewMock().WithQuery("SELECT").WithReply(dbResponse)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Request = req

	CreateTask(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), config.SUserInvalid)
}

// TestCreateTaskNoData with no data in request
func TestCreateTaskNoData(t *testing.T) {
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Request = req

	CreateTask(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing: title")
}

// TestCreateTask with no data in request
func TestCreateTask(t *testing.T) {
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	expectedTask := mock.GetMockTaskID0()
	req, err := tu.GetRequestPost(expectedTask, "/")
	assert.Nil(t, err)
	c.Request = req

	mocket.Catcher.NewMock().WithQuery("INSERT")

	CreateTask(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
}

// ################# TESTS
// FetchAllTask()

// TestFetchAllTaskNotAuthenticated func with no authenticated user
func TestFetchAllTaskNotAuthenticated(t *testing.T) {
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUserID0(false)
	dbResponse := mock.GetMapArrayByUser(u)
	mocket.Catcher.Reset().NewMock().WithQuery("SELECT").WithReply(dbResponse)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Request = req

	FetchAllTask(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), config.SUserInvalid)
}

// TestFetchAllTaskNoData with no data in db
func TestFetchAllTaskNoData(t *testing.T) {
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Request = req

	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "tasks"`)

	FetchAllTask(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), config.STaskNotFound)
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
}

// TestFetchAllTaskNoData test func with data
func TestFetchAllTask(t *testing.T) {
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Request = req

	tasks := mock.GetMockTasks()
	dbResponse := mock.GetMapArrayByTasks(tasks)
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "tasks"`).WithReply(dbResponse)

	FetchAllTask(c)

	tasksJSON, err := json.Marshal(tasks)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), string(tasksJSON))
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
}

// ################# TESTS
// FetchSingleTask()

// TestFetchSingleTaskNotAuthenticated func with no authenticated user
func TestFetchSingleTaskNotAuthenticated(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUserID0(false)
	dbResponse := mock.GetMapArrayByUser(u)
	mocket.Catcher.Reset().NewMock().WithQuery("SELECT").WithReply(dbResponse)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Request = req

	FetchSingleTask(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), config.SUserInvalid)
}

// TestFetchSingleTaskNoData with no data in request
func TestFetchSingleTaskNoData(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Request = req

	FetchSingleTask(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), config.STaskInvalid)
}

// TestFetchSingleTaskNotFound test not found
func TestFetchSingleTaskNotFound(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Params = append(c.Params, gin.Param{"id", strconv.FormatUint(uint64(u.ID), 10)})
	c.Request = req

	task := mock.GetMockTaskID0()
	dbResponse := mock.GetMapArrayByTask(task)
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "tasks"`).WithReply(dbResponse)

	FetchSingleTask(c)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), config.STaskNotFound)
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
}

// TestFetchSingleTaskUnauthorized test unauthorized
func TestFetchSingleTaskUnauthorized(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Params = append(c.Params, gin.Param{"id", strconv.FormatUint(uint64(u.ID), 10)})
	c.Request = req

	task := mock.GetMockTask()
	dbResponse := mock.GetMapArrayByTask(task)
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "tasks"`).WithReply(dbResponse)

	FetchSingleTask(c)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), config.STaskUnauthorized)
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
}

// TestFetchSingleTask test return of todo
func TestFetchSingleTask(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Params = append(c.Params, gin.Param{"id", strconv.FormatUint(uint64(u.ID), 10)})
	c.Request = req

	task := mock.GetMockTask()
	task.UserID = u.ID
	dbResponse := mock.GetMapArrayByTask(task)
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "tasks"`).WithReply(dbResponse)

	FetchSingleTask(c)

	taskJSON, err := json.Marshal(task)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), string(taskJSON))
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
}

// ################# TESTS
// UpdateTask()

// TestUpdateTaskNotAuthenticated func with no authenticated user
func TestUpdateTaskNotAuthenticated(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUserID0(false)
	dbResponse := mock.GetMapArrayByUser(u)
	mocket.Catcher.Reset().NewMock().WithQuery("SELECT").WithReply(dbResponse)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Request = req

	UpdateTask(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), config.SUserInvalid)
}

// TestUpdateTaskNoData with no data in request
func TestUpdateTaskNoData(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Request = req

	UpdateTask(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), config.STaskInvalid)
}

// TestUpdateTaskNotFound test not found
func TestUpdateTaskNotFound(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Params = append(c.Params, gin.Param{"id", strconv.FormatUint(uint64(u.ID), 10)})
	c.Request = req

	task := mock.GetMockTaskID0()
	dbResponse := mock.GetMapArrayByTask(task)
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "tasks"`).WithReply(dbResponse)

	UpdateTask(c)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), config.STaskNotFound)
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
}

// TestUpdateTaskEmpty test empty data request
func TestUpdateTaskEmpty(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Params = append(c.Params, gin.Param{"id", strconv.FormatUint(uint64(u.ID), 10)})
	c.Request = req

	task := mock.GetMockTask()
	dbResponse := mock.GetMapArrayByTask(task)
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "tasks"`).WithReply(dbResponse)

	UpdateTask(c)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Missing: title")
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
}

// TestUpdateTask test update of todo
func TestUpdateTask(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	newTask := mock.GetMockTask()

	req, err := tu.GetRequestPost(newTask, "/")
	assert.Nil(t, err)
	c.Params = append(c.Params, gin.Param{"id", strconv.FormatUint(uint64(u.ID), 10)})
	c.Request = req

	task := mock.GetMockTask()
	task.ID = newTask.ID
	task.UserID = u.ID
	dbResponse := mock.GetMapArrayByTask(task)
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "tasks"`).WithReply(dbResponse)
	mocket.Catcher.NewMock().WithQuery(`UPDATE "tasks" SET "title"`)
	mocket.Catcher.NewMock().WithQuery(`UPDATE "tasks" SET "description"`)
	mocket.Catcher.NewMock().WithQuery(`UPDATE "tasks" SET "done"`)

	UpdateTask(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, mocket.Catcher.Mocks[1].Triggered && !mocket.Catcher.Mocks[1].Once)
	assert.True(t, mocket.Catcher.Mocks[2].Triggered)
	assert.True(t, mocket.Catcher.Mocks[3].Triggered)
	assert.True(t, mocket.Catcher.Mocks[4].Triggered)
}

// ################# TESTS
// DeleteTask()

// TestDeleteTaskNotAuthenticated func with no authenticated user
func TestDeleteTaskNotAuthenticated(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUserID0(false)
	dbResponse := mock.GetMapArrayByUser(u)
	mocket.Catcher.Reset().NewMock().WithQuery("SELECT").WithReply(dbResponse)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Request = req

	DeleteTask(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), config.SUserInvalid)
}

// TestDeleteTaskNoData with no data in request
func TestDeleteTaskNoData(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Request = req

	DeleteTask(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), config.STaskInvalid)
}

// TestDeleteTaskNotFound test not found
func TestDeleteTaskNotFound(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Params = append(c.Params, gin.Param{"id", strconv.FormatUint(uint64(u.ID), 10)})
	c.Request = req

	task := mock.GetMockTaskID0()
	dbResponse := mock.GetMapArrayByTask(task)
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "tasks"`).WithReply(dbResponse)

	DeleteTask(c)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), config.STaskNotFound)
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
}

// TestDeleteTaskUnauthorized test unauthorized
func TestDeleteTaskUnauthorized(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Params = append(c.Params, gin.Param{"id", strconv.FormatUint(uint64(u.ID), 10)})
	c.Request = req

	task := mock.GetMockTask()
	dbResponse := mock.GetMapArrayByTask(task)
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "tasks"`).WithReply(dbResponse)

	DeleteTask(c)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), config.STaskUnauthorized)
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
}

// TestDeleteTask test return of todo
func TestDeleteTask(t *testing.T) {
	config.TestInit()
	w, c := tu.GetRecorderContext()

	u := mock.GetMockUser(false)
	mock.ConfigClaims(c, u)

	req, err := tu.GetRequestPost(nil, "/")
	assert.Nil(t, err)
	c.Params = append(c.Params, gin.Param{"id", strconv.FormatUint(uint64(u.ID), 10)})
	c.Request = req

	task := mock.GetMockTask()
	task.UserID = u.ID
	dbResponse := mock.GetMapArrayByTask(task)
	mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "tasks"`).WithReply(dbResponse)

	DeleteTask(c)

	taskJSON, err := json.Marshal(task)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), string(taskJSON))
	assert.True(t, mocket.Catcher.Mocks[1].Triggered)
}
