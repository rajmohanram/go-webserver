package websocket

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Hub manages WebSocket connections and broadcasts messages
type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()
			log.Printf("Client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error broadcasting to client: %v", err)
					client.Close()
					h.mu.RUnlock()
					h.mu.Lock()
					delete(h.clients, client)
					h.mu.Unlock()
					h.mu.RLock()
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register adds a new client
func (h *Hub) Register(client *websocket.Conn) {
	h.register <- client
}

// Unregister removes a client
func (h *Hub) Unregister(client *websocket.Conn) {
	h.unregister <- client
}

// Broadcast sends a message to all clients
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

// getRandomMessage returns a random message from a predefined list
func getRandomMessage() string {
	messages := []string{
		"Hello from the server!",
		"Random message: The weather is nice today",
		"Did you know? Go is a great language!",
		"Server notification: Everything is running smoothly",
		"Fun fact: WebSockets enable real-time communication",
		"Tip of the day: Stay hydrated!",
		"Random thought: Coffee makes everything better",
		"Server says: Keep coding!",
		"Interesting: This message was sent automatically",
		"Reminder: Take a break every hour",
		"Alert: This is a random server message",
		"Quote: 'Code is poetry' - Unknown",
	}
	return messages[rand.Intn(len(messages))]
}

// StartRandomMessageSender starts sending random messages at random intervals
func (h *Hub) StartRandomMessageSender() {
	go func() {
		for {
			// Random interval between 1 and 5 seconds
			interval := time.Duration(rand.Intn(4)+1) * time.Second
			time.Sleep(interval)

			h.mu.RLock()
			hasClients := len(h.clients) > 0
			h.mu.RUnlock()

			// Only send if there are connected clients
			if hasClients {
				message := fmt.Sprintf("🤖 Server: %s", getRandomMessage())
				h.Broadcast([]byte(message))
				log.Printf("Sent random message: %s", message)
			}
		}
	}()
}

