package controller

import (
	"fmt"
	"instacloneapp/server/pkg/db"
	"instacloneapp/server/socket"
	"instacloneapp/server/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// getUserIDFromContext retrieves the user ID from the Gin context
func getUserIDFromContext(c *gin.Context) string {
	if userID, exists := c.Get("userId"); exists {
		return userID.(string)
	}
	return ""
}

// AddNewPost handles adding a new post
func AddNewPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract author ID from URL parameters
		authorID := c.Param("author_id")
		authorIDObjectID, err := primitive.ObjectIDFromHex(authorID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid author ID"})
			return
		}

		// Parse caption from the request body
		var req struct {
			Caption string `json:"caption"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
			return
		}

		// Get the image from form data
		image, _, err := c.Request.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Image required"})
			return
		}
		defer image.Close()

		// Upload the image using the existing Cloudinary client
		imageURL, err := utils.UploadImageToCloudinary(cloudinaryClient, image, "posts")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error uploading image"})
			return
		}
		// Create a Post object
		postInput := db.Post{
			Caption:   req.Caption,
			Image:     imageURL,
			Author:    authorIDObjectID,
			CreatedAt: time.Now(),
		}
		// Create a new post in the database
		post, err := dbInstance.CreatePost(postInput)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating post"})
			return
		}

		// Update user by adding the post ID to the user's posts array
		err = dbInstance.AddPostToUser(authorIDObjectID, post.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating user with post"})
			return
		}

		// Return success response
		c.JSON(http.StatusCreated, gin.H{
			"message": "New post added",
			"post":    post,
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

		userID, err := primitive.ObjectIDFromHex(getUserIDFromContext(c)) // Extract user ID from the context or request

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid User ID"})

			return
		}
		post, err := dbInstance.GetPostByID(postID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
			return
		}

		err = dbInstance.RemoveLikeFromPost(postID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error disliking post"})
			return
		}

		// Implement socket.io for real-time notification
		user, err := dbInstance.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		postOwnerID := post.Author.Hex()
		if postOwnerID != userID.Hex() {
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

		userID, err := primitive.ObjectIDFromHex(getUserIDFromContext(c)) // Extract user ID from the context or request
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid User ID"})
			return
		}

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

		comment, err := dbInstance.CreateComment(userID, postID, req.Text)
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

		userID, err := primitive.ObjectIDFromHex(getUserIDFromContext(c)) // Extract user ID from the context or request
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid User ID"})
			return
		}

		post, err := dbInstance.GetPostByID(postID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
			return
		}

		if post.Author.Hex() != userID.Hex() {
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

		err = dbInstance.RemovePostFromUser(userID, postID)
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
		postID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Post ID"})
			return
		}

		userID, err := primitive.ObjectIDFromHex(getUserIDFromContext(c)) // Extract user ID from the context or request
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid User ID"})
			return
		}

		post, err := dbInstance.GetPostByID(postID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Post not found"})
			return
		}

		fmt.Println(post)

		user, err := dbInstance.GetUserByID(userID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}

		if contains(user.Bookmarks, postID) {
			err = dbInstance.RemoveBookmarkFromUser(userID, postID)
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
			err = dbInstance.AddBookmarkToUser(userID, postID)
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
