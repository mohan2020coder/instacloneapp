package db

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Username       string               `bson:"username" json:"username" binding:"required"`
	Email          string               `bson:"email" json:"email" binding:"required,email"`
	Password       string               `bson:"password" json:"password" binding:"required"`
	ProfilePicture string               `bson:"profilePicture,omitempty" json:"profilePicture,omitempty"`
	Bio            string               `bson:"bio,omitempty" json:"bio,omitempty"`
	Gender         string               `bson:"gender,omitempty" json:"gender,omitempty"`
	Followers      []primitive.ObjectID `bson:"followers,omitempty" json:"followers,omitempty"`
	Following      []primitive.ObjectID `bson:"following,omitempty" json:"following,omitempty"`
	Posts          []primitive.ObjectID `bson:"posts,omitempty" json:"posts,omitempty"`
	Bookmarks      []primitive.ObjectID `bson:"bookmarks,omitempty" json:"bookmarks,omitempty"`
	CreatedAt      time.Time            `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt      time.Time            `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
}

// SeedUsers seeds the user table with initial data
func SeedUsers(database Database) {
	users := []User{
		{Username: "Alice"},
		{Username: "Bob"},
		{Username: "Charlie"},
	}

	for _, user := range users {
		_, err := database.CreateUser(user)
		if err != nil {
			log.Printf("Failed to seed user: %v", err)
		} else {
			log.Printf("User %s seeded successfully.", user.Username)
		}
	}
}

// SeedDatabase runs all the seeders
func SeedDatabase(database Database) {
	SeedUsers(database)
	// Add other seed functions here
}
