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
			break
		}
		if messageType == websocket.TextMessage {
			fmt.Printf("Received message : %s\n", msg)
		}
	}

	// Clean up on disconnect
	mu.Lock()
	delete(userSocketMap, userID)
	mu.Unlock()
	broadcastOnlineUsers()
}

func broadcastOnlineUsers() {
	mu.Lock()
	defer mu.Unlock()

	onlineUsers := make([]string, 0, len(userSocketMap))
	for userID := range userSocketMap {
		onlineUsers = append(onlineUsers, userID)
	}

	for _, conn := range userSocketMap {
		err := conn.WriteJSON(onlineUsers)
		if err != nil {
			conn.Close()
		}
	}
}
