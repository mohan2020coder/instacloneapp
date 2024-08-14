package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Participants []primitive.ObjectID `bson:"participants,omitempty" json:"participants,omitempty"`
	Messages     []primitive.ObjectID `bson:"messages,omitempty" json:"messages,omitempty"`
}
