package config

import (
	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"

	// importing postgres for start open gorm connection
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// DB is the gorm connection to the postgress database
var DB *gorm.DB

// Init initialize the connection to the database.
func Init() *gorm.DB {
	db, err := gorm.Open("postgres", "host=postgrestodo port=5432 user=admin dbname=tododb password=123  sslmode=disable")

	if err != nil {
		panic(err.Error())
	}

	DB = db
	return DB
}

// GetDB returns the connection to the database.
func GetDB() *gorm.DB {
	return DB
}

// TestInit initialize the connection to the mock database driver for tests
func TestInit() *gorm.DB {
	mocket.Catcher.Register()
	mocket.Catcher.Logging = true

	db, _ := gorm.Open(mocket.DriverName, "connection_string")
	DB = db

	return DB
}
