// todoAPI Engine
package main

import (
	"log"
	"os"
	"strings"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/migration"
	"github.com/giuliobosco/todoAPI/route"
	"github.com/giuliobosco/todoAPI/utils"

	"github.com/gin-gonic/gin"
)

// Init initialize the application
func init() {
	utils.SetTesting(strings.Index(os.Args[0], ".test") >= 0)

	if utils.IsTesting() {
		config.TestInit()
	} else {
		db := config.Init()
		migration.Migrate(db)
	}
}

// main starts the app.
func main() {
	if utils.IsTesting() {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := route.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
