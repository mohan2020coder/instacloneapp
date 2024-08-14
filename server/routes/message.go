package routes

import (
	"instacloneapp/server/controller"
	"instacloneapp/server/middleware"
	"instacloneapp/server/pkg/db"

	"github.com/cloudinary/cloudinary-go"
	"github.com/gin-gonic/gin"
)

// SetupMessageRoutes sets up the routes for message-related endpoints

func SetupMessageRoutes(router *gin.Engine, database db.Database, cloudinaryClient *cloudinary.Cloudinary) {

	controller.InitUser(database, cloudinaryClient)
	messageRoutes := router.Group("/api/v1/message")
	{
		// Route to send a message
		messageRoutes.POST("/send/:id", middleware.IsAuthenticated(), controller.SendMessage())

		// Route to get messages
		messageRoutes.GET("/all/:id", middleware.IsAuthenticated(), controller.GetMessages())
	}
}
