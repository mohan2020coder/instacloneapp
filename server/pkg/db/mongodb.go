// package db

// import (
// 	"context"
// 	"errors"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // MongoDB implements the Database interface for MongoDB
// type MongoDB struct {
// 	conn *mongo.Collection
// }

// // NewMongoDB creates a new MongoDB connection
// func NewMongoDB(uri string, dbName string, collectionName string) (*MongoDB, error) {
// 	clientOptions := options.Client().ApplyURI(uri)
// 	client, err := mongo.Connect(context.Background(), clientOptions)
// 	if err != nil {
// 		return nil, err
// 	}

// 	collection := client.Database(dbName).Collection(collectionName)
// 	return &MongoDB{conn: collection}, nil
// }

// func (db *MongoDB) GetUsers() ([]User, error) {
// 	cursor, err := db.conn.Find(context.Background(), bson.M{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer cursor.Close(context.Background())

// 	var users []User
// 	for cursor.Next(context.Background()) {
// 		var user User
// 		if err := cursor.Decode(&user); err != nil {
// 			return nil, err
// 		}
// 		users = append(users, user)
// 	}

// 	return users, nil
// }

// func (db *MongoDB) CreateUser(u User) (User, error) {
// 	// Insert the user into the collection
// 	result, err := db.conn.InsertOne(context.Background(), u)
// 	if err != nil {
// 		return User{}, err
// 	}

// 	// Extract the inserted ID and assert it to primitive.ObjectID
// 	objectID, ok := result.InsertedID.(primitive.ObjectID)
// 	if !ok {
// 		return User{}, errors.New("failed to convert inserted ID to ObjectID")
// 	}

// 	// Set the ID field of the user
// 	u.ID = objectID

//		// Return the user with the MongoID set
//		return u, nil
//	}
package db

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// User represents the user model in MongoDB
// type User struct {
// 	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
// 	Username       string             `bson:"username" json:"username"`
// 	Email          string             `bson:"email" json:"email"`
// 	Password       string             `bson:"password" json:"password"`
// 	Bio            string             `bson:"bio,omitempty" json:"bio"`
// 	Gender         string             `bson:"gender,omitempty" json:"gender"`
// 	ProfilePicture string             `bson:"profilePicture,omitempty" json:"profilePicture"`
// 	Following      []primitive.ObjectID `bson:"following,omitempty" json:"following"`
// }

// MongoDB implements the Database interface for MongoDB
type MongoDB struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(uri string, dbName string, collectionName string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	collection := client.Database(dbName).Collection(collectionName)
	return &MongoDB{client: client, collection: collection}, nil
}

// Close disconnects from MongoDB
func (db *MongoDB) Close() error {
	return db.client.Disconnect(context.Background())
}

// GetUserByEmail retrieves a user by their email
func (db *MongoDB) GetUserByEmail(email string) (User, error) {
	var user User
	err := db.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, nil
		}
		return User{}, err
	}
	return user, nil
}

// CreateUser inserts a new user into the database
func (db *MongoDB) CreateUser(u User) (User, error) {
	result, err := db.collection.InsertOne(context.Background(), u)
	if err != nil {
		return User{}, err
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return User{}, errors.New("failed to convert inserted ID to ObjectID")
	}

	u.ID = objectID
	return u, nil
}

// UpdateUser updates a user's profile information
func (db *MongoDB) UpdateUser(id primitive.ObjectID, update interface{}) error {
	filter := bson.M{"_id": id}
	_, err := db.collection.UpdateOne(context.Background(), filter, update)
	return err
}

// GetAllUsers retrieves all users
func (db *MongoDB) GetAllUsers() ([]User, error) {
	cursor, err := db.collection.Find(context.Background(), bson.M{})
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

// FollowUser updates the following list of a user
func (db *MongoDB) FollowUser(followingUserID, targetUserID primitive.ObjectID) error {
	filter := bson.M{"_id": followingUserID}
	update := bson.M{"$addToSet": bson.M{"following": targetUserID}}
	_, err := db.collection.UpdateOne(context.Background(), filter, update)
	return err
}

// UnfollowUser updates the following list of a user
func (db *MongoDB) UnfollowUser(followingUserID, targetUserID primitive.ObjectID) error {
	filter := bson.M{"_id": followingUserID}
	update := bson.M{"$pull": bson.M{"following": targetUserID}}
	_, err := db.collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (db *MongoDB) GetUsers(filter interface{}) (*mongo.Cursor, error) {
	return db.collection.Find(context.Background(), filter)
}

func (db *MongoDB) GetUserByID(id primitive.ObjectID) (User, error) {
	var user User
	err := db.collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	return user, err
}

func (db *MongoDB) DeleteUser(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": id}
	return db.collection.DeleteOne(context.Background(), filter)
}

func (db *MongoDB) FollowOrUnfollowUser(followingUserID, targetUserID primitive.ObjectID, action string) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": followingUserID}
	var update bson.M

	if action == "follow" {
		update = bson.M{"$addToSet": bson.M{"following": targetUserID}}
	} else if action == "unfollow" {
		update = bson.M{"$pull": bson.M{"following": targetUserID}}
	} else {
		return nil, errors.New("invalid action")
	}

	return db.collection.UpdateOne(context.Background(), filter, update)
}
