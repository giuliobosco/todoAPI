package main

import (
	"github.com/calo001/todoAPI/config"
	"github.com/calo001/todoAPI/migration"
	"github.com/calo001/todoAPI/route"
	"github.com/gin-gonic/gin"
	"log"
	"os"
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