package migration

import (
	"todoAPI/model"

	"github.com/jinzhu/gorm"
)

// Migrate is the function for migrate the database.
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&model.Task{})
	db.AutoMigrate(&model.User{})
}
