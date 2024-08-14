package controller

import (
	"context"
	"instacloneapp/server/socket"
	"net/http"

	"instacloneapp/server/pkg/db"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// getUserIDFromContext retrieves the user ID from the Gin context
// func getUserIDFromContext(c *gin.Context) string {
// 	if userID, exists := c.Get("userId"); exists {
// 		return userID.(string)
// 	}
// 	return ""
// }

// AddNewPost handles adding a new post
func AddNewPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorID := c.Param("author_id")
		authorIDObjectID, err := primitive.ObjectIDFromHex(authorID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid author ID"})
			return
		}

		// Parse the form to handle file upload and JSON data
		if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to parse form"})
			return
		}

		// Assuming the image is uploaded as a file
		image, imageHeader, err := c.Request.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Image required"})
			return
		}

		// Upload image to Cloudinary
		uploadResult, err := cloudinaryClient.Upload.Upload(context.Background(), image, uploader.UploadParams{
			PublicID: imageHeader.Filename,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Image upload failed"})
			return
		}

		// Get the image URL from the upload result
		imageURL := uploadResult.SecureURL

		// Bind JSON data for caption
		var req struct {
			Caption string `json:"caption"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
			return
		}

		// Create a new post
		post := db.Post{
			Author:  authorIDObjectID,
			Caption: req.Caption,
			Image:   imageURL,
			// Set other necessary fields like CreatedAt, etc.
		}

		// Insert the post into the database and get a pointer to the created post
		createdPost, err := dbInstance.CreatePost(post)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating post"})
			return
		}

		// Update user with the new post
		err = dbInstance.AddPostToUser(authorIDObjectID, createdPost.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating user"})
			return
		}

		// Return the created post in the response
		c.JSON(http.StatusCreated, gin.H{
			"message": "New post added",
			"post":    createdPost,
			"success": true,
		})
	}
}

// GetAllPosts retrieves all posts
func GetAllPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		posts, err := dbInstance.GetAllPosts()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving posts"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"posts":   posts,
			"success": true,
		})
	}
}

// GetUserPosts retrieves all posts by a specific user
func GetUserPosts() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorID := c.Param("author_id")
		authorIDObjectID, err := primitive.ObjectIDFromHex(authorID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid author ID"})
			return
		}

		posts, err := dbInstance.GetPostsByUserID(authorIDObjectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving posts"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"posts":   posts,
			"success": true,
		})
	}
}

// LikePost handles liking a post
func LikePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("user_id")
		postID := c.Param("post_id")

		userIDObjectID, err := primitive.ObjectIDFromHex(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
			return
		}

		postIDObjectID, err := primitive.ObjectIDFromHex(postID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid post ID"})
			return
		}

		// Like the post
		err = dbInstance.AddLikeToPost(postIDObjectID, userIDObjectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error liking post"})
			return
		}

		// Notify the post owner
		post, err := dbInstance.GetPostByID(postIDObjectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving post"})
			return
		}

		postOwnerID := post.Author.Hex()
		if postOwnerID != userID {
			notification := bson.M{
				"type":    "like",
				"userId":  userID,
				"postId":  postID,
				"message": "Your post was liked",
			}
			socketID := socket.GetReceiverSocketID(postOwnerID)
			if socketID != "" {
				socket.BroadcastMessageToUser(socketID, "notification", notification)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Post liked",
			"success": true,
		})
	}
}

// DislikePost handles the logic for disliking a post
func DislikePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Post ID"})
			return
		}

		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User ID not found in context"})
		}

		userUID, err := primitive.ObjectIDFromHex(userID.(string)) // Extract user ID from the context or request

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User ID Hexx failed"})
			return
		}

		post, err := dbInstance.GetPostByID(postID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
			return
		}

		err = dbInstance.RemoveLikeFromPost(postID, userUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error disliking post"})
			return
		}

		// Implement socket.io for real-time notification
		user, err := dbInstance.GetUserByID(userUID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		postOwnerID := post.Author.Hex()
		if postOwnerID != userID {
			notification := bson.M{
				"type":    "dislike",
				"userId":  userID,
				"user":    user,
				"postId":  postID,
				"message": "Your post was disliked",
			}
			socketID := socket.GetReceiverSocketID(postOwnerID)
			if socketID != "" {
				socket.BroadcastMessageToUser(socketID, "notification", notification)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Post disliked",
			"success": true,
		})
	}
}

// AddComment handles adding a new comment to a post
func AddComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Post ID"})
			return
		}

		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User ID not found in context"})
		}

		userUID, err := primitive.ObjectIDFromHex(userID.(string))

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User ID Hexx failed"})
			return
		}
		// userID := getUserIDFromContext(c.Request) // Extract user ID from the context or request

		var req struct {
			Text string `json:"text"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
			return
		}

		if req.Text == "" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Text is required"})
			return
		}

		comment, err := dbInstance.CreateComment(userUID, postID, req.Text)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating comment"})
			return
		}

		err = dbInstance.AddCommentToPost(postID, comment.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error adding comment to post"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Comment added",
			"comment": comment,
			"success": true,
		})
	}
}

// GetCommentsOfPost handles fetching comments for a post
func GetCommentsOfPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Post ID"})
			return
		}

		comments, err := dbInstance.GetCommentsByPostID(postID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "No comments found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"comments": comments,
		})
	}
}

// DeletePost handles deleting a post and its comments
func DeletePost() gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Post ID"})
			return
		}

		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User ID not found in context"})
		}

		userUID, err := primitive.ObjectIDFromHex(userID.(string)) // Extract user ID from the context or request

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User ID Hexx failed"})
			return
		}

		post, err := dbInstance.GetPostByID(postID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
			return
		}

		if post.Author.Hex() != userID {
			c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
			return
		}

		err = dbInstance.DeletePost(postID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting post"})
			return
		}

		err = dbInstance.DeleteCommentsByPostID(postID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error deleting comments"})
			return
		}

		err = dbInstance.RemovePostFromUser(userUID, postID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating user posts"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Post deleted",
			"success": true,
		})
	}
}

// BookmarkPost handles bookmarking or removing a bookmark for a post
func BookmarkPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve post ID from URL parameters
		postIDStr := c.Param("id")
		postID, err := primitive.ObjectIDFromHex(postIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Post ID"})
			return
		}

		// Extract user ID from context
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User ID not found in context"})
			return
		}

		userUID, err := primitive.ObjectIDFromHex(userID.(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid User ID"})
			return
		}

		// Retrieve the user by their ID
		user, err := dbInstance.GetUserByID(userUID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		// Check if the post is already bookmarked by the user
		if contains(user.Bookmarks, postID) {
			// Remove bookmark if it already exists
			err = dbInstance.RemoveBookmarkFromUser(userUID, postID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error removing bookmark"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"type":    "unsaved",
				"message": "Post removed from bookmark",
				"success": true,
			})
		} else {
			// Add bookmark if it doesn't exist
			err = dbInstance.AddBookmarkToUser(userUID, postID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error adding bookmark"})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"type":    "saved",
				"message": "Post bookmarked",
				"success": true,
			})
		}
	}
}
