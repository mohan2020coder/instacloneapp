package controller

import (
	"instacloneapp/server/socket"
	"net/http"

	// "instacloneapp/server/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SendMessage handles sending a message
func SendMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		senderID := c.Param("id")
		receiverID := c.Param("receiver_id")
		var req struct {
			TextMessage string `json:"textMessage"`
		}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
			return
		}

		senderObjectID, err := primitive.ObjectIDFromHex(senderID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid sender ID"})
			return
		}
		receiverObjectID, err := primitive.ObjectIDFromHex(receiverID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid receiver ID"})
			return
		}

		// Check if conversation exists
		conversation, err := dbInstance.GetConversation(senderObjectID, receiverObjectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving conversation"})
			return
		}

		if conversation == nil {
			conversation, err = dbInstance.CreateConversation(senderObjectID, receiverObjectID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating conversation"})
				return
			}
		}

		// Create a new message
		newMessage, err := dbInstance.CreateMessage(senderObjectID, receiverObjectID, req.TextMessage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating message"})
			return
		}

		conversation.Messages = append(conversation.Messages, newMessage.ID)
		err = dbInstance.UpdateConversation(conversation.ID, bson.M{"$set": bson.M{"messages": conversation.Messages}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating conversation"})
			return
		}

		// Send real-time notification
		receiverSocketID := socket.GetReceiverSocketID(receiverObjectID.Hex())
		if receiverSocketID != "" {
			socket.BroadcastMessageToUser(receiverSocketID, "newMessage", newMessage)
		}

		c.JSON(http.StatusCreated, gin.H{
			"success":    true,
			"newMessage": newMessage,
		})
	}
}

// GetMessages handles retrieving messages from a conversation
func GetMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		senderID := c.Param("id")
		receiverID := c.Param("receiver_id")

		senderObjectID, err := primitive.ObjectIDFromHex(senderID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid sender ID"})
			return
		}
		receiverObjectID, err := primitive.ObjectIDFromHex(receiverID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid receiver ID"})
			return
		}

		conversation, err := dbInstance.GetConversation(senderObjectID, receiverObjectID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving conversation"})
			return
		}

		if conversation == nil {
			c.JSON(http.StatusOK, gin.H{"success": true, "messages": []interface{}{}})
			return
		}

		messages, err := dbInstance.GetMessagesByIDs(conversation.Messages)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error retrieving messages"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"messages": messages,
		})
	}
}
