package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post represents the MongoDB schema for a Post.

type Post struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Caption   string               `bson:"caption,omitempty" json:"caption,omitempty"`
	Image     string               `bson:"image" json:"image"`
	Author    primitive.ObjectID   `bson:"author,omitempty" json:"author,omitempty"`
	Likes     []primitive.ObjectID `bson:"likes,omitempty" json:"likes,omitempty"`
	Comments  []primitive.ObjectID `bson:"comments,omitempty" json:"comments,omitempty"`
	CreatedAt time.Time            `bson:"createdAt,omitempty"`
	UpdatedAt time.Time            `bson:"updatedAt,omitempty"`
}
