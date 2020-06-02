// todoAPI Engine
package main

import (
	"log"
	"os"

	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/migration"
	"github.com/giuliobosco/todoAPI/route"

	"github.com/gin-gonic/gin"
)

// Init initialize the application
func init() {
	db := config.Init()
	migration.Migrate(db)
}

// main starts the app.
func main() {
	gin.SetMode(gin.ReleaseMode)

	router := route.SetupRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
