package socket

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var userSocketMap = make(map[string]*websocket.Conn) // Stores WebSocket connections corresponding to user IDs
var mu sync.Mutex                                    // Mutex for synchronizing access to userSocketMap

// GetReceiverSocketID retrieves the WebSocket ID for a given user ID
func GetReceiverSocketID(userID string) string {
	mu.Lock()
	defer mu.Unlock()
	conn, ok := userSocketMap[userID]
	if ok {
		return conn.RemoteAddr().String() // This might not be an ideal unique identifier
	}
	return ""
}

// BroadcastMessageToUser sends a message to a specific user
func BroadcastMessageToUser(userID, event string, message interface{}) {
	mu.Lock()
	defer mu.Unlock()
	conn, ok := userSocketMap[userID]
	if ok {
		err := conn.WriteJSON(map[string]interface{}{
			"event":   event,
			"message": message,
		})
		if err != nil {
			fmt.Printf("Error writing message to user %s: %v\n", userID, err)
			conn.Close()
			delete(userSocketMap, userID)
		}
	}
}

// HandleConnection handles WebSocket connections
func HandleConnection(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	defer conn.Close()

	userID := c.Query("userId")
	if userID != "" {
		mu.Lock()
		userSocketMap[userID] = conn
		mu.Unlock()
	}

	// Broadcast online users
	broadcastOnlineUsers()

	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("Error reading message: %v\n", err)
			break
		}
		if messageType == websocket.TextMessage {
			fmt.Printf("Received message from %s: %s\n", userID, msg)
		}
	}

	// Clean up on disconnect
	mu.Lock()
	delete(userSocketMap, userID)
	mu.Unlock()
	broadcastOnlineUsers()
}

// broadcastOnlineUsers broadcasts the list of online users to all connected clients
func broadcastOnlineUsers() {
	mu.Lock()
	defer mu.Unlock()

	onlineUsers := make([]string, 0, len(userSocketMap))
	for userID := range userSocketMap {
		onlineUsers = append(onlineUsers, userID)
	}

	for userID, conn := range userSocketMap {
		err := conn.WriteJSON(onlineUsers)
		if err != nil {
			fmt.Printf("Error broadcasting online users to %s: %v\n", userID, err)
			conn.Close()
		}
	}
}
