package main

import (
	"log"
	"net/http"
	"os"

	"go-webservere/internal/api"
	"go-webservere/internal/store"
	"go-webservere/internal/websocket"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize store
	userStore := store.NewStore()

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()
	hub.StartRandomMessageSender()

	// Initialize handlers
	apiHandler := api.NewHandler(userStore)
	wsHandler := websocket.NewHandler(hub)

	// Create router
	r := mux.NewRouter()

	// Homepage
	r.HandleFunc("/", homeHandler).Methods("GET")

	// REST API endpoints
	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/users", apiHandler.GetUsers).Methods("GET")
	apiRouter.HandleFunc("/users/{id}", apiHandler.GetUser).Methods("GET")
	apiRouter.HandleFunc("/users", apiHandler.CreateUser).Methods("POST")
	apiRouter.HandleFunc("/users/{id}", apiHandler.UpdateUser).Methods("PUT")
	apiRouter.HandleFunc("/users/{id}", apiHandler.DeleteUser).Methods("DELETE")

	// WebSocket endpoint
	r.HandleFunc("/ws/", wsHandler.HandleWebSocket)

	// Start server
	port := getPort()
	log.Printf("Server starting on port %s", port)
	log.Printf("Homepage: http://localhost%s", port)
	log.Printf("API: http://localhost%s/api/v1/users", port)
	log.Printf("WebSocket: ws://localhost%s/ws/", port)

	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/index.html")
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if port[0] != ':' {
		port = ":" + port
	}
	return port
}
