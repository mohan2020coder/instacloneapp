// package routes

// import (
// 	"instacloneapp/server/controller"
// 	"instacloneapp/server/pkg/db"

// 	"github.com/gin-gonic/gin"
// )

// func SetupRoutes(router *gin.Engine, database db.Database) {
// 	userController := &controller.UserController{DB: database}

//		router.GET("/api/v1/user", userController.GetUsers)
//		router.POST("/api/v1/user", userController.CreateUser)
//	}
package routes

import (
	"instacloneapp/server/controller"
	"instacloneapp/server/pkg/db"

	"github.com/cloudinary/cloudinary-go"
	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up the routes for the application
func SetupRoutes(router *gin.Engine, database db.Database, cloudinaryClient *cloudinary.Cloudinary) {

	controller.InitUser(database, cloudinaryClient)
	// Create a new group for user-related routes
	userRoutes := router.Group("/api/v1/user")
	{
		// Route for user registration
		userRoutes.POST("/register", controller.Register())

		// Route for user login
		userRoutes.POST("/login", controller.Login())

		// Route for user logout
		userRoutes.POST("/logout", controller.Logout())

		// Route to get a user's profile by ID
		userRoutes.GET("/:id/profile", controller.GetProfile())

		// Route to edit a user's profile (e.g., username, bio, etc.)
		userRoutes.PUT("/profile/edit", controller.EditProfile())

		// Route to get suggested users (for follow suggestions, etc.)
		userRoutes.GET("/suggested", controller.GetSuggestedUsers())

		// Route to  unfollow a user based on their ID
		userRoutes.POST("/followorunfollow/:id", controller.FollowOrUnfollowUser())

	}

}
