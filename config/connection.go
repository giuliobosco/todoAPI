package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func Init() *gorm.DB {
	db, err := gorm.Open("postgres", "host=postgrestodo port=5432 user=admin dbname=tododb password=123  sslmode=disable")

	if err != nil {
		panic(err.Error())
	}

	DB = db
	return DB
}

func GetDB() *gorm.DB {
	return DB
}
