package main

import (
	"github.com/gin-gonic/gin"

	"github.com/rony-tampubolon/sample-rest-api/controllers"
	"github.com/rony-tampubolon/sample-rest-api/middlewares"
	"github.com/rony-tampubolon/sample-rest-api/models"
)

func main() {
	router := gin.Default()

	models.ConnectToDatabase() // new!

	// Authentication
	router.POST("/register", controllers.SignUp)
	router.POST("/login", controllers.SignIn)
	router.POST("/refresh/token", middlewares.AuthenticatedUser(), controllers.RefreshToken)

	// Create Post Group
	postGroup := router.Group("/post").Use(middlewares.AuthenticatedUser())
	// Element in Router  POST
	postGroup.POST("/", controllers.CreatePost)
	postGroup.GET("/", controllers.GetAllPost)
	postGroup.GET("/:id", controllers.GetById)
	postGroup.PUT("/:id", controllers.UpdatePost)
	postGroup.DELETE("/:id", controllers.DeletePost)

	router.Run("localhost:8080")
}
