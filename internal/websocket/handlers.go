package websocket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// Handler manages WebSocket connections
type Handler struct {
	hub *Hub
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

// HandleWebSocket upgrades HTTP connection to WebSocket
func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Debug: Print all request headers
	log.Println("=== WebSocket Connection Attempt ===")
	log.Printf("Method: %s, URL: %s, Proto: %s", r.Method, r.URL.String(), r.Proto)
	log.Println("Request Headers:")
	for name, values := range r.Header {
		for _, value := range values {
			log.Printf("  %s: %s (len=%d)", name, value, len(value))
		}
	}
	log.Println("====================================")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	h.hub.Register(conn)

	go func() {
		defer func() {
			h.hub.Unregister(conn)
		}()

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket error: %v", err)
				}
				break
			}

			log.Printf("Received message: %s", message)

			// Broadcast message to all connected clients
			broadcastMsg := fmt.Sprintf("Broadcast: %s", string(message))
			h.hub.Broadcast([]byte(broadcastMsg))
		}
	}()
}
