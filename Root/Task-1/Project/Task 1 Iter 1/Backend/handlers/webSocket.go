package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	role := c.GetString("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not allowed"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

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
