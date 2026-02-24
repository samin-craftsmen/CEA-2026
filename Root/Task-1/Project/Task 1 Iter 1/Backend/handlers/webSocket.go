package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
	"github.com/samin-craftsmen/gin-project/utils"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins (adjust for production)
		},
	}
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan interface{})
	mu        sync.Mutex
)

func init() {
	go broadcastHandler()
}

func broadcastHandler() {
	for {
		message := <-broadcast
		mu.Lock()
		for client := range clients {
			err := client.WriteJSON(message)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

func WebSocketHeadcount(c *gin.Context) {
	// Extract token from query params
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	// Validate token
	parsedToken, err := utils.ValidateToken(token)
	if err != nil {
		log.Println("Token validation error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Extract claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	// Check if user is admin
	role, ok := claims["role"].(string)
	if !ok || role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin only"})
		return
	}

	// Get date from URL params
	date := c.Param("date")
	if date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing date"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// Register client
	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	log.Printf("WebSocket client connected for date: %s\n", date)

	// Keep connection alive
	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			break
		}
	}
}

func BroadcastHeadcountUpdate(date string, headcount interface{}) {
	broadcast <- gin.H{
		"type": "headcount_update",
		"date": date,
		"data": headcount,
	}
}
