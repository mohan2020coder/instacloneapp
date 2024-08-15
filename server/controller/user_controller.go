// package controller

// import (
// 	"fmt"
// 	"instacloneapp/server/pkg/db"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"golang.org/x/crypto/bcrypt"
// )

// // UserController handles user-related API requests
// type UserController struct {
// 	DB db.Database
// }

// func (uc *UserController) CreateUser(c *gin.Context) {
// 	// Create an instance of the User struct
// 	var user db.User

// 	// Bind the incoming JSON data to the User struct
// 	if err := c.BindJSON(&user); err != nil {
// 		// Log the error for debugging
// 		fmt.Println("Error binding JSON:", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
// 		return
// 	}

// 	// Validate required fields (you can add more validation if needed)
// 	if user.Username == "" || user.Email == "" || user.Password == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, email, and password are required"})
// 		return
// 	}

// 	// Hash the password before saving it (ensure you have a hashing function)
// 	user.Password = hashPassword(user.Password)

// 	// Set the ID and timestamps manually
// 	user.ID = primitive.NewObjectID() // Generate a new ObjectID for the user
// 	user.CreatedAt = time.Now()
// 	user.UpdatedAt = time.Now()

// 	// Log the user data to check what is being saved
// 	fmt.Println("Creating user:", user)

// 	// Save the user to the database
// 	createdUser, err := uc.DB.CreateUser(user)
// 	if err != nil {
// 		// Log the error for debugging
// 		fmt.Println("Error saving user to database:", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
// 		return
// 	}

//		// Respond with the created user data
//		c.JSON(http.StatusCreated, createdUser)
//	}
//
//	func hashPassword(password string) string {
//		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
//		if err != nil {
//			log.Fatalf("Failed to hash password: %v", err)
//		}
//		return string(hashedPassword)
//	}
package controller

import (
	"bytes"
	"context"
	"fmt"
	"instacloneapp/server/pkg/db"
	"instacloneapp/server/utils"
	"log"
	"net/http"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	//"golang.org/x/crypto/bcrypt"
)

// Controller holds dependencies
var (
	dbInstance       db.Database
	cloudinaryClient *cloudinary.Cloudinary
)

// NewController creates a new instance of Controller
func InitUser(database db.Database, cloudinary *cloudinary.Cloudinary) {
	dbInstance = database
	cloudinaryClient = cloudinary
}

// func hashPassword(password string) string {
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		log.Fatalf("Failed to hash password: %v", err)
// 	}
// 	return string(hashedPassword)
// }

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		filter := bson.M{}

		cursor, err := dbInstance.GetUsers(filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cursor.Close(context.Background())

		var users []db.User
		for cursor.Next(context.Background()) {
			var user db.User
			if err := cursor.Decode(&user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			users = append(users, user)
		}

		c.JSON(http.StatusOK, users)
	}
}

// Register handles user registration
func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
			return
		}

		if req.Username == "" || req.Email == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "All fields are required"})
			return
		}

		// Retrieve user by email
		user, err := dbInstance.GetUserByEmail(req.Email)
		if err != nil {
			// Log the error for further investigation
			log.Printf("Error checking email: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error checking email"})
			return
		}

		// Check if the user was found
		// Assuming a zero value of user indicates "not found"
		if user.Email != "" { // or any other field to check if user exists
			c.JSON(http.StatusBadRequest, gin.H{"message": "Email already exists"})
			return
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		newUser := db.User{
			Username: req.Username,
			Email:    req.Email,
			Password: string(hashedPassword),
		}
		_, err = dbInstance.CreateUser(newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating account"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Account created successfully"})
	}
}

// Login handles user login
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
			return
		}

		// Retrieve the user by email
		user, err := dbInstance.GetUserByEmail(req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect email or password"})
			return
		}

		// Check if the user exists and validate password
		if isEmptyUser(user) {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect email or password"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect email or password"})
			return
		}

		// Generate token
		token, err := utils.GenerateToken(user.ID.Hex())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error generating token"})
			return
		}

		fmt.Println("Token from cookie:", token)

		// Set the token in the cookie
		c.SetCookie("token", token, 24*60*60, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	}
}

// isEmptyUser checks if a user object is empty
func isEmptyUser(user db.User) bool {
	return user.ID == primitive.NilObjectID && user.Username == "" && user.Email == ""
}

// Logout handles user logout
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("token", "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}

// GetProfile retrieves user profile
func GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
			return
		}

		user, err := dbInstance.GetUserByID(objectID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

// EditProfile handles updating a user's profile
func EditProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		objectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
			return
		}

		var req struct {
			Bio    string `json:"bio"`
			Gender string `json:"gender"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
			return
		}

		var cloudResponse *uploader.UploadResult

		// Handle profile picture upload if provided
		file, _, err := c.Request.FormFile("profile_picture")
		if err == nil { // If there's no file, this will be skipped
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(file); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error reading file"})
				return
			}

			resp, err := cloudinaryClient.Upload.Upload(context.Background(), buf.Bytes(), uploader.UploadParams{
				Folder: "profile_pictures",
			})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading file"})
				return
			}
			cloudResponse = resp
		}

		// Find the user
		user, err := dbInstance.GetUserByID(objectID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		// Update user details
		updateFields := bson.M{}
		if req.Bio != "" {
			updateFields["bio"] = req.Bio
		}
		if req.Gender != "" {
			updateFields["gender"] = req.Gender
		}
		if cloudResponse != nil {
			updateFields["profilePicture"] = cloudResponse.SecureURL
		}

		if len(updateFields) > 0 {
			err := dbInstance.UpdateUser(objectID, bson.M{"$set": updateFields})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating profile"})
				return
			}
		}

		// Return response
		c.JSON(http.StatusOK, gin.H{
			"message": "Profile updated successfully",
			"success": true,
			"user":    user,
		})
	}
}

// // UpdateProfile handles profile update
// func UpdateProfile() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userID := c.Param("id")
// 		objectID, err := primitive.ObjectIDFromHex(userID)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
// 			return
// 		}

// 		var updateFields map[string]interface{}
// 		if err := c.BindJSON(&updateFields); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
// 			return
// 		}

// 		err = dbInstance.UpdateUser(objectID, bson.M{"$set": updateFields})
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating profile"})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
// 	}
// }

// // UploadProfilePicture handles profile picture upload
// func UploadProfilePicture() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		file, _, err := c.Request.FormFile("profile_picture")
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"message": "No file uploaded"})
// 			return
// 		}

// 		// Convert the file to a bytes.Buffer
// 		var buf bytes.Buffer
// 		if _, err := buf.ReadFrom(file); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error reading file"})
// 			return
// 		}

// 		// Upload to Cloudinary
// 		resp, err := cloudinaryClient.Upload.Upload(context.Background(), buf.Bytes(), uploader.UploadParams{
// 			Folder: "profile_pictures",
// 		})
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading file"})
// 			return
// 		}

// 		userID := c.Param("id")
// 		objectID, err := primitive.ObjectIDFromHex(userID)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
// 			return
// 		}

// 		err = dbInstance.UpdateUser(objectID, bson.M{"$set": bson.M{"profilePicture": resp.SecureURL}})
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating profile picture"})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"message": "Profile picture uploaded successfully"})
// 	}
// }

// GetSuggestedUsers retrieves suggested users
func GetSuggestedUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		// userID := c.Param("id")
		userID, exists := c.Get("userID") // The ID of the user initiating the follow/unfollow action
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}
		objectID, err := primitive.ObjectIDFromHex(userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
			return
		}

		// Retrieve the user to check their following list
		user, err := dbInstance.GetUserByID(objectID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		// Find users who are not followed by the current user
		filter := bson.M{
			"$and": []bson.M{
				{"_id": bson.M{"$ne": objectID}},        // Exclude the current user
				{"_id": bson.M{"$nin": user.Following}}, // Exclude users already followed
			},
		}

		cursor, err := dbInstance.GetUsers(filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching suggested users"})
			return
		}
		defer cursor.Close(context.Background())

		var suggestions []db.User
		for cursor.Next(context.Background()) {
			var suggestedUser db.User
			if err := cursor.Decode(&suggestedUser); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error decoding user data"})
				return
			}
			suggestions = append(suggestions, suggestedUser)
		}

		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Cursor error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"suggested_users": suggestions})
	}
}
func FollowOrUnfollowUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		followKrneWala, exists := c.Get("userID") // The ID of the user initiating the follow/unfollow action
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		jiskoFollowKrunga := c.Param("id") // The ID of the target user
		if followKrneWala == jiskoFollowKrunga {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "You cannot follow/unfollow yourself",
				"success": false,
			})
			return
		}

		followingUserID, err := primitive.ObjectIDFromHex(followKrneWala.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
			return
		}

		targetUserID, err := primitive.ObjectIDFromHex(jiskoFollowKrunga)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid target user ID"})
			return
		}

		// Determine the action (follow/unfollow)
		user, err := dbInstance.GetUserByID(followingUserID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		isFollowing := contains(user.Following, targetUserID)
		action := "follow"
		if isFollowing {
			action = "unfollow"
		}

		// Use the FollowOrUnfollowUser function from the db interface
		result, err := dbInstance.FollowOrUnfollowUser(followingUserID, targetUserID, action)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error performing action"})
			return
		}

		message := "Followed successfully"
		if action == "unfollow" {
			message = "Unfollowed successfully"
		}

		c.JSON(http.StatusOK, gin.H{
			"message": message,
			"success": true,
			"result":  result,
		})
	}
}

// Helper function to check if a user is already following another user
func contains(following []primitive.ObjectID, targetUserID primitive.ObjectID) bool {
	for _, id := range following {
		if id == targetUserID {
			return true
		}
	}
	return false
}
