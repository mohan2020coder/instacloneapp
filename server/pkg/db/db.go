package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database interface {
	GetUsers(filter interface{}) (*mongo.Cursor, error) // Get users based on filter
	GetUserByID(id primitive.ObjectID) (User, error)    // Retrieve a single user by ID
	CreateUser(User) (User, error)
	GetUserByEmail(email string) (User, error)                                                                         // Create a new user
	UpdateUser(id primitive.ObjectID, update interface{}) error                                                        // Update a user's information
	DeleteUser(id primitive.ObjectID) (*mongo.DeleteResult, error)                                                     // Delete a user by ID
	FollowOrUnfollowUser(followingUserID, targetUserID primitive.ObjectID, action string) (*mongo.UpdateResult, error) // Follow or unfollow a user

	// Conversation operations
	GetConversation(senderID, receiverID primitive.ObjectID) (*Conversation, error)
	UpdateConversation(id primitive.ObjectID, update interface{}) error
	CreateConversation(participant1, participant2 primitive.ObjectID) (*Conversation, error)

	// Message operations
	GetMessagesByIDs(ids []primitive.ObjectID) ([]Message, error)
	CreateMessage(senderID, receiverID primitive.ObjectID, messageText string) (*Message, error)
	RemoveBookmarkFromUser(userID, postID primitive.ObjectID) error
	AddBookmarkToUser(userID, postID primitive.ObjectID) error
	RemovePostFromUser(userID, postID primitive.ObjectID) error
	CreateComment(authorID, postID primitive.ObjectID, text string) (*Comment, error)
	DeleteCommentsByPostID(postID primitive.ObjectID) error
	GetPostByID(postID primitive.ObjectID) (*Post, error)
	RemoveLikeFromPost(postID, userID primitive.ObjectID) error
	AddCommentToPost(postID, commentID primitive.ObjectID) error
	DeletePost(postID primitive.ObjectID) error
	GetCommentsByPostID(postID primitive.ObjectID) ([]Comment, error)
	CreatePost(post Post) (*Post, error)
	AddPostToUser(userID primitive.ObjectID, postID primitive.ObjectID) error
	GetAllPosts() ([]Post, error)
	GetPostsByUserID(authorID primitive.ObjectID) ([]Post, error)
	AddLikeToPost(postID primitive.ObjectID, userID primitive.ObjectID) error
}
