package controller

// import (
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/nfnt/resize"
// 	"github.com/cloudinary/cloudinary-go/v2"
// 	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
// 	"github.com/sirupsen/logrus"
// 	"gopkg.in/mgo.v2/bson"

// 	"your_project/models"
// 	"your_project/socket"
// )

// var cld *cloudinary.Cloudinary

// func init() {
// 	cld, _ = cloudinary.NewFromParams("cloud_name", "api_key", "api_secret")
// }

// func AddNewPost(c *gin.Context) {
// 	authorID := c.GetString("id")
// 	caption := c.PostForm("caption")

// 	file, _, err := c.Request.FormFile("image")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Image required"})
// 		return
// 	}

// 	// Resize image using the "resize" package
// 	img, _, err := image.Decode(file)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to decode image"})
// 		return
// 	}

// 	resizedImg := resize.Resize(800, 800, img, resize.Lanczos3)
// 	buffer := new(bytes.Buffer)
// 	if err := jpeg.Encode(buffer, resizedImg, nil); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to encode image"})
// 		return
// 	}

// 	// Upload to Cloudinary
// 	fileUri := fmt.Sprintf("data:image/jpeg;base64,%s", base64.StdEncoding.EncodeToString(buffer.Bytes()))
// 	resp, err := cld.Upload.Upload(context.Background(), fileUri, uploader.UploadParams{Folder: "posts"})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to upload image"})
// 		return
// 	}

// 	post := models.Post{
// 		Caption: caption,
// 		Image:   resp.SecureURL,
// 		Author:  authorID,
// 	}

// 	if err := post.Save(); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save post"})
// 		return
// 	}

// 	user := models.User{}
// 	if err := user.FindByID(authorID); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to find user"})
// 		return
// 	}
// 	user.Posts = append(user.Posts, post.ID)
// 	if err := user.Save(); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user"})
// 		return
// 	}

// 	post.Author = user
// 	c.JSON(http.StatusCreated, gin.H{
// 		"message": "New post added",
// 		"post":    post,
// 		"success": true,
// 	})
// }

// func GetAllPosts(c *gin.Context) {
// 	posts, err := models.GetPosts()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve posts"})
// 		return
// 	}

// 	for i := range posts {
// 		posts[i].PopulateAuthor()
// 		posts[i].PopulateComments()
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"posts":  posts,
// 		"success": true,
// 	})
// }

// func GetUserPosts(c *gin.Context) {
// 	authorID := c.GetString("id")
// 	posts, err := models.GetPostsByAuthor(authorID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve posts"})
// 		return
// 	}

// 	for i := range posts {
// 		posts[i].PopulateAuthor()
// 		posts[i].PopulateComments()
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"posts":  posts,
// 		"success": true,
// 	})
// }

// func LikePost(c *gin.Context) {
// 	userID := c.GetString("id")
// 	postID := c.Param("id")

// 	post := models.Post{}
// 	if err := post.FindByID(postID); err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found", "success": false})
// 		return
// 	}

// 	if err := post.AddLike(userID); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to like post", "success": false})
// 		return
// 	}

// 	postOwnerID := post.Author
// 	if postOwnerID != userID {
// 		user := models.User{}
// 		if err := user.FindByID(userID); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to find user"})
// 			return
// 		}

// 		notification := models.Notification{
// 			Type:        "like",
// 			UserID:      userID,
// 			UserDetails: user,
// 			PostID:      postID,
// 			Message:     "Your post was liked",
// 		}
// 		postOwnerSocketID := socket.GetReceiverSocketId(postOwnerID)
// 		socket.IO.To(postOwnerSocketID).Emit("notification", notification)
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Post liked", "success": true})
// }

// func DislikePost(c *gin.Context) {
// 	userID := c.GetString("id")
// 	postID := c.Param("id")

// 	post := models.Post{}
// 	if err := post.FindByID(postID); err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found", "success": false})
// 		return
// 	}

// 	if err := post.RemoveLike(userID); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to dislike post", "success": false})
// 		return
// 	}

// 	postOwnerID := post.Author
// 	if postOwnerID != userID {
// 		user := models.User{}
// 		if err := user.FindByID(userID); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to find user"})
// 			return
// 		}

// 		notification := models.Notification{
// 			Type:        "dislike",
// 			UserID:      userID,
// 			UserDetails: user,
// 			PostID:      postID,
// 			Message:     "Your post was disliked",
// 		}
// 		postOwnerSocketID := socket.GetReceiverSocketId(postOwnerID)
// 		socket.IO.To(postOwnerSocketID).Emit("notification", notification)
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Post disliked", "success": true})
// }

// func AddComment(c *gin.Context) {
// 	postID := c.Param("id")
// 	userID := c.GetString("id")
// 	text := c.PostForm("text")

// 	if text == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Text is required", "success": false})
// 		return
// 	}

// 	post := models.Post{}
// 	if err := post.FindByID(postID); err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found", "success": false})
// 		return
// 	}

// 	comment := models.Comment{
// 		Text:   text,
// 		Author: userID,
// 		PostID: postID,
// 	}
// 	if err := comment.Save(); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to save comment"})
// 		return
// 	}

// 	post.Comments = append(post.Comments, comment.ID)
// 	if err := post.Save(); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update post"})
// 		return
// 	}

// 	comment.PopulateAuthor()
// 	c.JSON(http.StatusCreated, gin.H{
// 		"message": "Comment Added",
// 		"comment": comment,
// 		"success": true,
// 	})
// }

// func GetCommentsOfPost(c *gin.Context) {
// 	postID := c.Param("id")

// 	comments, err := models.GetCommentsByPostID(postID)
// 	if err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "No comments found for this post", "success": false})
// 		return
// 	}

// 	for i := range comments {
// 		comments[i].PopulateAuthor()
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true, "comments": comments})
// }

// func DeletePost(c *gin.Context) {
// 	postID := c.Param("id")
// 	authorID := c.GetString("id")

// 	post := models.Post{}
// 	if err := post.FindByID(postID); err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found", "success": false})
// 		return
// 	}

// 	if post.Author != authorID {
// 		c.JSON(http.StatusForbidden, gin.H{"message": "Unauthorized"})
// 		return
// 	}

// 	if err := post.Delete(); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete post"})
// 		return
// 	}

// 	user := models.User{}
// 	if err := user.FindByID(authorID); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to find user"})
// 		return
// 	}
// 	user.Posts = remove(user.Posts, postID)
// 	if err := user.Save(); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update user"})
// 		return
// 	}

// 	if err := models.DeleteCommentsByPostID(postID); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete comments"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Post deleted"})
// }

// func BookmarkPost(c *gin.Context) {
// 	postID := c.Param("id")
// 	userID := c.GetString("id")

// 	post := models.Post{}
// 	if err := post.FindByID(postID); err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found", "success": false})
// 		return
// 	}

// 	user := models.User{}
// 	if err := user.FindByID(userID); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to find user"})
// 		return
// 	}

// 	if contains(user.Bookmarks, post.ID) {
// 		user.Bookmarks = remove(user.Bookmarks, post.ID)
// 	} else {
// 		user.Bookmarks = append(user.Bookmarks, post.ID)
// 	}

// 	if err := user.Save(); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update bookmarks"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"type":    ifContains(user.Bookmarks, post.ID) ? "saved" : "unsaved",
// 		"message": ifContains(user.Bookmarks, post.ID) ? "Post bookmarked" : "Post removed from bookmark",
// 		"success": true,
// 	})
// }

// func contains(slice []bson.ObjectId, item bson.ObjectId) bool {
// 	for _, v := range slice {
// 		if v == item {
// 			return true
// 		}
// 	}
// 	return false
// }

// func remove(slice []bson.ObjectId, item bson.ObjectId) []bson.ObjectId {
// 	for i, v := range slice {
// 		if v == item {
// 			return append(slice[:i], slice[i+1:]...)
// 		}
// 	}
// 	return slice
// }
