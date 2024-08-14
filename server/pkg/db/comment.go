package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Text   string             `bson:"text" json:"text"`
	Author primitive.ObjectID `bson:"author,omitempty" json:"author,omitempty"`
	Post   primitive.ObjectID `bson:"post,omitempty" json:"post,omitempty"`
}
