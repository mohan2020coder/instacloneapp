package routes

import (
	"instacloneapp/server/controller"
	"instacloneapp/server/pkg/db"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, database db.Database) {
	userController := &controller.UserController{DB: database}

	router.GET("/api/v1/user", userController.GetUsers)
	router.POST("/api/v1/user", userController.CreateUser)
}
