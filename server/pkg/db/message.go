package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SenderID   primitive.ObjectID `bson:"senderId,omitempty" json:"senderId,omitempty"`
	ReceiverID primitive.ObjectID `bson:"receiverId,omitempty" json:"receiverId,omitempty"`
	Message    string             `bson:"message" json:"message"`
}
