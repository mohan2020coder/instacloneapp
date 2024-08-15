package db

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB implements the Database interface for MongoDB
// type MongoDB struct {
// 	client     *mongo.Client
// 	collection *mongo.Collection
// }

// // NewMongoDB creates a new MongoDB connection
// func NewMongoDB(uri string, dbName string, collectionName string) (*MongoDB, error) {
// 	clientOptions := options.Client().ApplyURI(uri)
// 	client, err := mongo.Connect(context.Background(), clientOptions)
// 	if err != nil {
// 		return nil, err
// 	}

// 	collection := client.Database(dbName).Collection(collectionName)
// 	return &MongoDB{client: client, collection: collection}, nil
// }

type MongoDB struct {
	client      *mongo.Client
	collections map[string]*mongo.Collection
}

func NewMongoDB(uri string, dbName string, collectionNames string) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Initialize collections
	collections := make(map[string]*mongo.Collection)
	for _, name := range strings.Split(collectionNames, ",") {
		collections[name] = client.Database(dbName).Collection(name)
	}

	return &MongoDB{client: client, collections: collections}, nil
}

func (db *MongoDB) GetCollection(name string) (*mongo.Collection, bool) {
	collection, exists := db.collections[name]
	return collection, exists
}

// Close disconnects from MongoDB
func (db *MongoDB) Close() error {
	return db.client.Disconnect(context.Background())
}

func (db *MongoDB) GetUserByEmail(email string) (User, error) {
	var user User
	collection, exists := db.GetCollection("users") // Specify the collection name
	if !exists {
		return User{}, errors.New("collection 'users' does not exist")
	}

	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return User{}, nil
		}
		return User{}, err
	}
	return user, nil
}

func (db *MongoDB) CreateUser(u User) (User, error) {
	collection, exists := db.GetCollection("users") // Specify the collection name
	if !exists {
		return User{}, errors.New("collection 'users' does not exist")
	}

	result, err := collection.InsertOne(context.Background(), u)
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

func (db *MongoDB) UpdateUser(id primitive.ObjectID, update interface{}) error {
	collection, exists := db.GetCollection("users") // Specify the collection name
	if !exists {
		return errors.New("collection 'users' does not exist")
	}

	filter := bson.M{"_id": id}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (db *MongoDB) GetAllUsers() ([]User, error) {
	collection, exists := db.GetCollection("users") // Specify the collection name
	if !exists {
		return nil, errors.New("collection 'users' does not exist")
	}

	cursor, err := collection.Find(context.Background(), bson.M{})
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

func (db *MongoDB) FollowUser(followingUserID, targetUserID primitive.ObjectID) error {
	collection, exists := db.GetCollection("users") // Specify the collection name
	if !exists {
		return errors.New("collection 'users' does not exist")
	}

	filter := bson.M{"_id": followingUserID}
	update := bson.M{"$addToSet": bson.M{"following": targetUserID}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (db *MongoDB) UnfollowUser(followingUserID, targetUserID primitive.ObjectID) error {
	collection, exists := db.GetCollection("users") // Specify the collection name
	if !exists {
		return errors.New("collection 'users' does not exist")
	}

	filter := bson.M{"_id": followingUserID}
	update := bson.M{"$pull": bson.M{"following": targetUserID}}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

func (db *MongoDB) GetUsers(filter interface{}) (*mongo.Cursor, error) {
	collection, exists := db.GetCollection("users")
	if !exists {
		return nil, errors.New("collection 'users' does not exist")
	}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	// Do not close the cursor here; let the caller handle it.
	return cursor, nil
}

func (db *MongoDB) GetUserByID(id primitive.ObjectID) (User, error) {
	var user User
	collection, exists := db.GetCollection("users") // Specify the collection name
	if !exists {
		return User{}, errors.New("collection 'users' does not exist")
	}

	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	return user, err
}

func (db *MongoDB) DeleteUser(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	collection, exists := db.GetCollection("users") // Specify the collection name
	if !exists {
		return nil, errors.New("collection 'users' does not exist")
	}

	filter := bson.M{"_id": id}
	return collection.DeleteOne(context.Background(), filter)
}

func (db *MongoDB) FollowOrUnfollowUser(followingUserID, targetUserID primitive.ObjectID, action string) (*mongo.UpdateResult, error) {
	collection, exists := db.GetCollection("users") // Specify the collection name
	if !exists {
		return nil, errors.New("collection 'users' does not exist")
	}

	filter := bson.M{"_id": followingUserID}
	var update bson.M

	if action == "follow" {
		update = bson.M{"$addToSet": bson.M{"following": targetUserID}}
	} else if action == "unfollow" {
		update = bson.M{"$pull": bson.M{"following": targetUserID}}
	} else {
		return nil, errors.New("invalid action")
	}

	return collection.UpdateOne(context.Background(), filter, update)
}

// GetConversation retrieves a conversation by participants' IDs
func (db *MongoDB) GetConversation(senderID, receiverID primitive.ObjectID) (*Conversation, error) {
	collection, exists := db.GetCollection("conversations") // Get the collection and existence flag
	if !exists {
		return nil, errors.New("collection 'conversations' does not exist")
	}

	filter := bson.M{
		"participants": bson.M{"$all": []primitive.ObjectID{senderID, receiverID}},
	}
	var conversation Conversation
	err := collection.FindOne(context.Background(), filter).Decode(&conversation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No conversation found
		}
		return nil, err
	}
	return &conversation, nil
}

// UpdateConversation updates a conversation with the provided data
func (db *MongoDB) UpdateConversation(id primitive.ObjectID, update interface{}) error {
	collection, exists := db.GetCollection("conversations") // Get the collection and existence flag
	if !exists {
		return errors.New("collection 'conversations' does not exist")
	}

	filter := bson.M{"_id": id}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	return err
}

// CreateMessage creates a new message
func (db *MongoDB) CreateMessage(senderID, receiverID primitive.ObjectID, messageText string) (*Message, error) {
	collection, exists := db.GetCollection("messages") // Get the collection and existence flag
	if !exists {
		return nil, errors.New("collection 'messages' does not exist")
	}

	message := Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Message:    messageText,
	}
	result, err := collection.InsertOne(context.Background(), message)
	if err != nil {
		return nil, err
	}
	message.ID = result.InsertedID.(primitive.ObjectID)
	return &message, nil
}

// CreateConversation creates a new conversation
func (db *MongoDB) CreateConversation(participant1, participant2 primitive.ObjectID) (*Conversation, error) {
	collection, exists := db.GetCollection("conversations") // Get the collection and existence flag
	if !exists {
		return nil, errors.New("collection 'conversations' does not exist")
	}

	conversation := Conversation{
		Participants: []primitive.ObjectID{participant1, participant2},
	}
	result, err := collection.InsertOne(context.Background(), conversation)
	if err != nil {
		return nil, err
	}
	conversation.ID = result.InsertedID.(primitive.ObjectID)
	return &conversation, nil
}

// RemoveBookmarkFromUser removes a bookmark from a user
func (db *MongoDB) RemoveBookmarkFromUser(userID, postID primitive.ObjectID) error {
	collection, exists := db.GetCollection("users") // Get the collection and existence flag
	if !exists {
		return errors.New("collection 'users' does not exist")
	}

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": userID},
		bson.M{"$pull": bson.M{"bookmarks": postID}},
	)
	return err
}

// AddBookmarkToUser adds a bookmark to a user
func (db *MongoDB) AddBookmarkToUser(userID, postID primitive.ObjectID) error {
	collection, exists := db.GetCollection("users") // Get the collection and existence flag
	if !exists {
		return errors.New("collection 'users' does not exist")
	}

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": userID},
		bson.M{"$addToSet": bson.M{"bookmarks": postID}},
	)
	return err
}

// RemovePostFromUser removes a post ID from the user's list of posts
func (db *MongoDB) RemovePostFromUser(userID, postID primitive.ObjectID) error {
	collection, exists := db.GetCollection("users") // Get the collection and existence flag
	if !exists {
		return errors.New("collection 'users' does not exist")
	}

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": userID},
		bson.M{"$pull": bson.M{"posts": postID}},
	)
	return err
}

// CreateComment creates a new comment
func (db *MongoDB) CreateComment(authorID, postID primitive.ObjectID, text string) (*Comment, error) {
	collection, exists := db.GetCollection("comments") // Get the collection and existence flag
	if !exists {
		return nil, errors.New("collection 'comments' does not exist")
	}

	comment := Comment{
		Author: authorID,
		Post:   postID,
		Text:   text,
	}
	result, err := collection.InsertOne(context.Background(), comment)
	if err != nil {
		return nil, err
	}
	comment.ID = result.InsertedID.(primitive.ObjectID)
	return &comment, nil
}

// DeleteCommentsByPostID deletes comments by post ID
func (db *MongoDB) DeleteCommentsByPostID(postID primitive.ObjectID) error {
	collection, exists := db.GetCollection("comments") // Get the collection and existence flag
	if !exists {
		return errors.New("collection 'comments' does not exist")
	}

	_, err := collection.DeleteMany(context.Background(), bson.M{"post": postID})
	return err
}

// GetPostByID retrieves a post by its ID
func (db *MongoDB) GetPostByID(postID primitive.ObjectID) (*Post, error) {
	collection, exists := db.GetCollection("posts") // Get the collection and existence flag
	if !exists {
		return nil, errors.New("collection 'posts' does not exist")
	}

	var post Post
	err := collection.FindOne(context.Background(), bson.M{"_id": postID}).Decode(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// RemoveLikeFromPost removes a like from a post
func (db *MongoDB) RemoveLikeFromPost(postID, userID primitive.ObjectID) error {
	collection, exists := db.GetCollection("posts") // Get the collection and existence flag
	if !exists {
		return errors.New("collection 'posts' does not exist")
	}

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": postID},
		bson.M{"$pull": bson.M{"likes": userID}},
	)
	return err
}

// AddCommentToPost adds a comment to a post
func (db *MongoDB) AddCommentToPost(postID, commentID primitive.ObjectID) error {
	collection, exists := db.GetCollection("posts") // Get the collection and existence flag
	if !exists {
		return errors.New("collection 'posts' does not exist")
	}

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": postID},
		bson.M{"$push": bson.M{"comments": commentID}},
	)
	return err
}

// DeletePost deletes a post by its ID
func (db *MongoDB) DeletePost(postID primitive.ObjectID) error {
	collection, exists := db.GetCollection("posts") // Get the collection and existence flag
	if !exists {
		return errors.New("collection 'posts' does not exist")
	}

	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": postID})
	return err
}

// GetCommentsByPostID retrieves comments for a post by its ID
func (db *MongoDB) GetCommentsByPostID(postID primitive.ObjectID) ([]Comment, error) {
	collection, exists := db.GetCollection("comments") // Get the collection and existence flag
	if !exists {
		return nil, errors.New("collection 'comments' does not exist")
	}

	cursor, err := collection.Find(context.Background(), bson.M{"post": postID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var comments []Comment
	for cursor.Next(context.Background()) {
		var comment Comment
		if err := cursor.Decode(&comment); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

// GetMessagesByIDs retrieves messages by their IDs
func (db *MongoDB) GetMessagesByIDs(ids []primitive.ObjectID) ([]Message, error) {
	collection, exists := db.GetCollection("messages") // Get the collection and existence flag
	if !exists {
		return nil, errors.New("collection 'messages' does not exist")
	}

	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var messages []Message
	for cursor.Next(context.Background()) {
		var message Message
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// / CreatePost creates a new post in the database
func (db *MongoDB) CreatePost(post Post) (*Post, error) {
	// Set the created time for the post
	post.CreatedAt = time.Now()

	// Get the collection
	collection, exists := db.GetCollection("posts")
	if !exists {
		return nil, errors.New("collection 'posts' does not exist")
	}

	// Insert the post into the collection
	result, err := collection.InsertOne(context.Background(), post)
	if err != nil {
		return nil, err
	}

	// Assign the generated ID to the post
	post.ID = result.InsertedID.(primitive.ObjectID)
	return &post, nil
}

// AddPostToUser adds a post ID to the user's list of posts
func (db *MongoDB) AddPostToUser(userID primitive.ObjectID, postID primitive.ObjectID) error {
	// Get the collection
	collection, exists := db.GetCollection("users")
	if !exists {
		return errors.New("collection 'users' does not exist")
	}

	filter := bson.M{"_id": userID}
	update := bson.M{
		"$push": bson.M{"posts": postID},
	}

	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

// GetAllPosts retrieves all posts from the database
func (db *MongoDB) GetAllPosts() ([]Post, error) {
	// Get the collection
	collection, exists := db.GetCollection("posts")
	if !exists {
		return nil, errors.New("collection 'posts' does not exist")
	}

	var posts []Post
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Iterate through the cursor and decode each document into a Post struct
	for cursor.Next(context.Background()) {
		var post Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// GetPostsByUserID retrieves all posts by a specific user
func (db *MongoDB) GetPostsByUserID(authorID primitive.ObjectID) ([]Post, error) {
	// Get the collection
	collection, exists := db.GetCollection("posts")
	if !exists {
		return nil, errors.New("collection 'posts' does not exist")
	}

	var posts []Post

	// Define the filter to match the authorID field
	filter := bson.M{"author": authorID}

	// Query the database for posts that match the authorID
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	// Iterate through the cursor and decode each document into a Post struct
	for cursor.Next(context.Background()) {
		var post Post
		if err := cursor.Decode(&post); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	// Check if there were any errors during cursor iteration
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

// AddLikeToPost adds a like to a post by updating the post's likes list
func (db *MongoDB) AddLikeToPost(postID primitive.ObjectID, userID primitive.ObjectID) error {
	// Get the collection
	collection, exists := db.GetCollection("posts")
	if !exists {
		return errors.New("collection 'posts' does not exist")
	}

	// Define the filter to match the postID
	filter := bson.M{"_id": postID}

	// Define the update to add the userID to the likes array
	update := bson.M{
		"$addToSet": bson.M{"likes": userID},
	}

	// Perform the update operation
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
