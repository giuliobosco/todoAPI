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
	authMiddleware, err := auth.SetupAuth(config.CredsFilePath)

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	router.GET("/", baseRoute)

	v1 := router.Group("/v1")
	{

		todo := v1.Group("todo")
		{
			todo.POST("/create", authMiddleware.MiddlewareFunc(), controller.CreateTask)
			todo.GET("/all", authMiddleware.MiddlewareFunc(), controller.FetchAllTask)
			todo.GET("/get/:id", authMiddleware.MiddlewareFunc(), controller.FetchSingleTask)
			todo.PUT("/update/:id", authMiddleware.MiddlewareFunc(), controller.UpdateTask)
			todo.DELETE("/delete/:id", authMiddleware.MiddlewareFunc(), controller.DeleteTask)
		}

		user := v1.Group("user")
		{
			user.GET("/get", authMiddleware.MiddlewareFunc(), controller.FetchUser)
			user.PUT("/update", authMiddleware.MiddlewareFunc(), controller.UpdateUser)
			user.DELETE("/delete", authMiddleware.MiddlewareFunc(), controller.DeleteUser)
		}
	}

	authGroup := router.Group("auth")
	{
		authGroup.GET("/refresh_token", authMiddleware.RefreshHandler)
		authGroup.POST("/login", authMiddleware.LoginHandler)
		authGroup.POST("/logout", authMiddleware.LogoutHandler)

		authGroup.GET("/oauth_url", auth.OAuthURL)
		authGroup.GET("/oauth", authMiddleware.LoginHandler)

		authGroup.POST("/register", controller.RegisterEndPoint)

		authGroup.GET("/confirm", controller.ConfirmUser)
		authGroup.GET("/sendAgainConfirm", controller.SendUserConfirmAgain)

		authGroup.GET("/requestPasswordRecovery", controller.RequestPasswordRecovery)

		authGroup.POST("/executePasswordRecovery", controller.ExecutePasswordRecovery)

		authGroup.POST("/updatePassword", authMiddleware.MiddlewareFunc(), controller.UpdatePassword)
	}

	return router
}

// baseRoute handle the GET /
func baseRoute(c *gin.Context) {
	c.String(http.StatusOK, config.SWelcome)
}

// SetupTestRouter prepare a test router
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.Default()
}
