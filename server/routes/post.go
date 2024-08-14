package routes

import (
	"instacloneapp/server/controller"
	"instacloneapp/server/middleware"
	"instacloneapp/server/pkg/db"

	"github.com/cloudinary/cloudinary-go"
	"github.com/gin-gonic/gin"
)

// SetupPostRoutes sets up the routes for post-related endpoints
func SetupPostRoutes(router *gin.Engine, database db.Database, cloudinaryClient *cloudinary.Cloudinary) {
	controller.InitUser(database, cloudinaryClient)

	postRoutes := router.Group("/api/v1/post")
	{
		// Route to add a new post
		postRoutes.POST("/addpost", middleware.IsAuthenticated(), controller.AddNewPost())

		// Route to get all posts
		postRoutes.GET("/all", middleware.IsAuthenticated(), controller.GetAllPosts())

		// Route to get posts by a user
		postRoutes.GET("/userpost/all", middleware.IsAuthenticated(), controller.GetUserPosts())

		// Route to like a post
		postRoutes.GET("/:id/like", middleware.IsAuthenticated(), controller.LikePost())

		// Route to dislike a post
		postRoutes.GET("/:id/dislike", middleware.IsAuthenticated(), controller.DislikePost())

		// Route to add a comment to a post
		postRoutes.POST("/:id/comment", middleware.IsAuthenticated(), controller.AddComment())

		// Route to get all comments for a post
		postRoutes.POST("/:id/comment/all", middleware.IsAuthenticated(), controller.GetCommentsOfPost())

		// Route to delete a post
		postRoutes.DELETE("/delete/:id", middleware.IsAuthenticated(), controller.DeletePost())

		// Route to bookmark or unbookmark a post
		postRoutes.GET("/:id/bookmark", middleware.IsAuthenticated(), controller.BookmarkPost())
	}
}
