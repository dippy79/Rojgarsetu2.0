package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rojgarsetu/backend/internal/services"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type WSHandler struct {
	notificationService *services.NotificationService
	clients              map[string]*websocket.Conn
	clientsMutex         sync.RWMutex
}

func NewWSHandler(notificationService *services.NotificationService) *WSHandler {
	return &WSHandler{
		notificationService: notificationService,
		clients:              make(map[string]*websocket.Conn),
	}
}

// HandleWebSocket upgrades HTTP connection to WebSocket
// @Summary WebSocket connection for real-time notifications
// @Description Establish WebSocket connection for real-time notifications
// @Tags websocket
// @Success 101 {string} string "Switching Protocols"
// @Router /api/v1/ws [get]
func (h *WSHandler) HandleWebSocket(c *gin.Context) {
	// Get user ID from context (should be set by auth middleware)
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Store client connection
	h.clientsMutex.Lock()
	h.clients[userID] = conn
	h.clientsMutex.Unlock()

	log.Printf("WebSocket client connected: %s", userID)

	// Send welcome message
	conn.WriteJSON(gin.H{
		"type": "connected",
		"message": "WebSocket connection established",
		"user_id": userID,
	})

	// Keep connection alive and handle incoming messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		// Handle incoming messages (pings, etc.)
		if messageType == websocket.TextMessage {
			log.Printf("Received message from %s: %s", userID, string(message))
		}
	}

	// Clean up on disconnect
	h.clientsMutex.Lock()
	delete(h.clients, userID)
	h.clientsMutex.Unlock()
	log.Printf("WebSocket client disconnected: %s", userID)
}

// BroadcastNotification sends a notification to a specific user
func (h *WSHandler) BroadcastNotification(userID string, notification map[string]interface{}) error {
	h.clientsMutex.RLock()
	conn, exists := h.clients[userID]
	h.clientsMutex.RUnlock()

	if !exists {
		return nil // User not connected, that's okay
	}

	return conn.WriteJSON(notification)
}

// BroadcastToAll sends a notification to all connected clients
func (h *WSHandler) BroadcastToAll(notification map[string]interface{}) {
	h.clientsMutex.RLock()
	defer h.clientsMutex.RUnlock()

	for userID, conn := range h.clients {
		if err := conn.WriteJSON(notification); err != nil {
			log.Printf("Error sending to client %s: %v", userID, err)
			// Remove disconnected client
			h.clientsMutex.Lock()
			delete(h.clients, userID)
			h.clientsMutex.Unlock()
		}
	}
}

// GetConnectedClients returns the number of connected clients
func (h *WSHandler) GetConnectedClients() int {
	h.clientsMutex.RLock()
	defer h.clientsMutex.RUnlock()
	return len(h.clients)
}