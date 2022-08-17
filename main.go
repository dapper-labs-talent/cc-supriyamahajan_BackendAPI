package main

import (
	"cc-supriyamahajan_BackendAPI/controllers"
	"cc-supriyamahajan_BackendAPI/db"
	"cc-supriyamahajan_BackendAPI/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database
	db.Connect()

	// Initialize Router
	router := initializeRouter()
	router.Run(":8080")
}

func initializeRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	api := router.Group("/v1")
	{
		api.POST("/user/signup", controllers.SignUp)
		api.POST("/user/login", controllers.Login)
		api.Use(middleware.Auth())
		{
			api.GET("/users", controllers.GetUsers)
			api.PUT("/users", controllers.UpdateUser)
		}
	}
	return router
}
