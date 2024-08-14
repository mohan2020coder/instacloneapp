package controller

// import (
// 	"context"
// 	"net/http"

// 	// "instacloneapp/server/services"
// 	"instacloneapp/server/socket"

// 	"github.com/gin-gonic/gin"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/bson/primitive"
// 	"go.mongodb.org/mongo-driver/mongo"
// )

// var conversationCollection *mongo.Collection
// var messageCollection *mongo.Collection

// func init() {
// 	// Initialize collections (ensure `config` is properly set up)
// 	conversationCollection = services.GetCollection("conversations")
// 	messageCollection = services.GetCollection("messages")
// }

// func SendMessage(c *gin.Context) {
// 	var requestBody struct {
// 		TextMessage string `json:"textMessage" binding:"required"`
// 	}

// 	if err := c.BindJSON(&requestBody); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	senderID, _ := primitive.ObjectIDFromHex(c.GetString("userID"))
// 	receiverID, _ := primitive.ObjectIDFromHex(c.Param("id"))

// 	var conversation models.Conversation
// 	filter := bson.M{"participants": bson.M{"$all": []primitive.ObjectID{senderID, receiverID}}}
// 	err := conversationCollection.FindOne(context.Background(), filter).Decode(&conversation)

// 	if err == mongo.ErrNoDocuments {
// 		conversation = models.Conversation{
// 			Participants: []primitive.ObjectID{senderID, receiverID},
// 		}
// 		_, err = conversationCollection.InsertOne(context.Background(), conversation)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
// 			return
// 		}
// 	}

// 	message := models.Message{
// 		SenderID:   senderID,
// 		ReceiverID: receiverID,
// 		Message:    requestBody.TextMessage,
// 	}
// 	_, err = messageCollection.InsertOne(context.Background(), message)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
// 		return
// 	}

// 	// Update conversation with the new message
// 	conversation.Messages = append(conversation.Messages, message.ID)
// 	_, err = conversationCollection.UpdateOne(
// 		context.Background(),
// 		bson.M{"_id": conversation.ID},
// 		bson.M{"$set": bson.M{"messages": conversation.Messages}},
// 	)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update conversation"})
// 		return
// 	}

// 	// Emit message through WebSocket
// 	receiverSocketID := socket.GetReceiverSocketID(receiverID.Hex())
// 	if receiverSocketID != "" {
// 		socket.BroadcastToSocket(receiverSocketID, "newMessage", message)
// 	}

// 	c.JSON(http.StatusCreated, gin.H{"success": true, "newMessage": message})
// }

// func GetMessages(c *gin.Context) {
// 	senderID, _ := primitive.ObjectIDFromHex(c.GetString("userID"))
// 	receiverID, _ := primitive.ObjectIDFromHex(c.Param("id"))

// 	var conversation models.Conversation
// 	filter := bson.M{"participants": bson.M{"$all": []primitive.ObjectID{senderID, receiverID}}}
// 	err := conversationCollection.FindOne(context.Background(), filter).Decode(&conversation)
// 	if err != nil {
// 		if err == mongo.ErrNoDocuments {
// 			c.JSON(http.StatusOK, gin.H{"success": true, "messages": []interface{}{}})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve conversation"})
// 		return
// 	}

// 	var messages []models.Message
// 	cursor, err := messageCollection.Find(context.Background(), bson.M{"_id": bson.M{"$in": conversation.Messages}})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
// 		return
// 	}
// 	defer cursor.Close(context.Background())

// 	for cursor.Next(context.Background()) {
// 		var message models.Message
// 		err := cursor.Decode(&message)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode message"})
// 			return
// 		}
// 		messages = append(messages, message)
// 	}

// 	c.JSON(http.StatusOK, gin.H{"success": true, "messages": messages})
// }
