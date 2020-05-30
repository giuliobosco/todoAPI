package migration

import (
	"github.com/giuliobosco/todoAPI/model"
	"github.com/jinzhu/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&model.Task{})
	db.AutoMigrate(&model.User{})
}
