package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/migration"
	"github.com/giuliobosco/todoAPI/route"
)

func init() {
	db := config.Init()
	migration.Migrate(db)
}

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
