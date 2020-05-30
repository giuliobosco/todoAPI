package model

import (
	"time"
)

type User struct {
	Base
	Username string `json:"username"`
	Password string `json:"password"`
	Todos    []Task `json:"todos"`
}

type Task struct {
	Base
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      uint   `json:"userid"`
}

type Base struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
