package config

import "github.cojm/jinzhu/gorm"

var DB *gorm.DB

func Init() *gorm.DB {
	db, err := gorm.Open("postgres",
		"host=postgrestodo port=5432 user=admin dbname=tododb password=123 sslmode=disable")

	if err != nil {
		panic(err.Error())
	}

	Db = db
	return DB
}

func GetDb() *gorm.DB {
	return DB
}