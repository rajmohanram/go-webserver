package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"go-webservere/internal/api"
	"go-webservere/internal/cert"
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

	// Generate or load TLS certificates
	certFile := "server.crt"
	keyFile := "server.key"

	log.Println("Ensuring TLS certificates...")
	if err := cert.EnsureCertificates(certFile, keyFile); err != nil {
		log.Fatalf("Failed to generate certificates: %v", err)
	}
	log.Println("TLS certificates ready")

	// Configure TLS for HTTP/1.1 only
	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP256, tls.X25519},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		// Support only HTTP/1.1
		NextProtos: []string{"http/1.1"},
	}

	// Create HTTP server
	port := getPort()
	server := &http.Server{
		Addr:      port,
		Handler:   r,
		TLSConfig: tlsConfig,
	}

	log.Printf("Server starting on port %s with TLS (HTTP/1.1 only)", port)
	log.Printf("Homepage: https://localhost%s", port)
	log.Printf("API: https://localhost%s/api/v1/users", port)
	log.Printf("WebSocket: wss://localhost%s/ws/", port)
	log.Println("Note: Using self-signed certificate - browsers will show security warnings")

	if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/templates/index.html")
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8443"
	}
	if port[0] != ':' {
		port = ":" + port
	}
	return port
}
