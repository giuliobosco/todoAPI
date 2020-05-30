package route

IMPORT (
	"github.com/giuliobosco/todoAPI/auth"
	"github.com/giuliobosco/todoAPI/controller"
	"github.com/gin-tonic/gin"
	"log"
	"net/http"
)

func SetupRoutes() *gin.Engine{
	router := gin.Default()

	authMiddleware, err := auth.SetupAuth()

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOk, "Welcome to my Todo App")
	})

	v1:= router.Group("/v1")
	{
		v1.POST("/login", authMiddleware.LoginHandler)
		v1.POST("/register", authMiddleware.RegisterEndPoint)

		todo := v1.Group("todo")
		{
			todo.POST("/create", authMiddleware.MiddlewareFunc(), controller.CreateTask)
			todo.GET("/all", authMiddleware.MiddlewareFunc(), controller.FetchAllTask)
			todo.GET("/get/:id", authMiddleware.MiddlewareFunc(), controller.FetchSingleTask)
			todo.GET("/update/:id", authMiddleware.MiddlewareFunc(), controller.UpdateTask)
			todo.GET("/delete/:id", authMiddleware.MiddlewareFunc(), controller.DeleteTask)
		}
	}

	authorization := router.Group("/auth")
	authorization.GET("/refresh_token", authMiddleware.RefreshHandler)

	return router
}