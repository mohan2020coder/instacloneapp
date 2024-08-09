package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB implements the Database interface for MongoDB
type MongoDB struct {
	conn *mongo.Collection
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(uri string, dbName string, collectionName string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)
	return &MongoDB{conn: collection}, nil
}

func (db *MongoDB) GetUsers() ([]User, error) {
	cursor, err := db.conn.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []User
	for cursor.Next(context.Background()) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (db *MongoDB) CreateUser(u User) (User, error) {
	result, err := db.conn.InsertOne(context.Background(), User{Name: u.Name})
	if err != nil {
		return User{}, err
	}

	id := result.InsertedID.(int64) // Assuming ID is an int64
	return User{ID: id, Name: u.Name}, nil
}
