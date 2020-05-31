// Package model contains the models of the data of the API Engine.
package model

import (
	"time"
)

// User is the rapresentation of the user
type User struct {
	Base            // use base object as parent
	Username string `json:"username"` // username of the user
	Password string `json:"password"` // password of the user
	Todos    []Task `json:"todos"`    // list of the todos of the user
}

// Task is the rappresentation of a task
type Task struct {
	Base               // user base object as parent
	Title       string `json:"title"`       // title of the task
	Description string `json:"description"` // description of the task
	UserID      uint   `json:"userid"`      // id of the user owner of the task
}

// Base is the basic object with basic components
type Base struct {
	ID        uint       `gorm:"primary_key" json:"id"` // id of the object
	CreatedAt time.Time  `json:"created_at"`            // object creation time
	UpdatedAt time.Time  `json:"updated_at"`            // object updating time
	DeletedAt *time.Time `json:"deleted_at"`            // object deleting time
}
