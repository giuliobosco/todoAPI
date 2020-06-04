/*
Package route create the routes of the API Server.
*/
package route

import (
	"log"
	"net/http"

	"github.com/giuliobosco/todoAPI/auth"
	"github.com/giuliobosco/todoAPI/config"
	"github.com/giuliobosco/todoAPI/controller"

	"github.com/gin-gonic/gin"
)

// SetupRoutes create the router of the API Engine.
func SetupRoutes() *gin.Engine {
	router := gin.Default()
	authMiddleware, err := auth.SetupAuth()

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, config.SWelcome)
	})

	v1 := router.Group("/v1")
	{
		v1.POST("/login", authMiddleware.LoginHandler)
		v1.POST("/logout", authMiddleware.LogoutHandler)

		v1.POST("/register", controller.RegisterEndPoint)

		v1.GET("/confirm", controller.ConfirmUser)

		v1.GET("/requestPasswordRecovery", controller.RequestPasswordRecovery)

		v1.POST("/executePasswordRecovery", controller.ExecutePasswordRecovery)

		v1.POST("/updatePassword", authMiddleware.MiddlewareFunc(), controller.UpdatePassword)

		v1.PUT("/updateUser", authMiddleware.MiddlewareFunc(), controller.UpdateUser)

		v1.GET("/user", authMiddleware.MiddlewareFunc(), controller.FetchUser)

		todo := v1.Group("todo")
		{
			todo.POST("/create", authMiddleware.MiddlewareFunc(), controller.CreateTask)
			todo.GET("/all", authMiddleware.MiddlewareFunc(), controller.FetchAllTask)
			todo.GET("/get/:id", authMiddleware.MiddlewareFunc(), controller.FetchSingleTask)
			todo.PUT("/update/:id", authMiddleware.MiddlewareFunc(), controller.UpdateTask)
			todo.DELETE("/delete/:id", authMiddleware.MiddlewareFunc(), controller.DeleteTask)
		}
	}

	authorization := router.Group("/auth")
	authorization.GET("/refresh_token", authMiddleware.RefreshHandler)

	return router
}
