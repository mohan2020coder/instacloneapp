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
	
}