package main

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/rajmohanram/go-webserver/internal/api"
	"github.com/rajmohanram/go-webserver/internal/cert"
	"github.com/rajmohanram/go-webserver/internal/store"
	"github.com/rajmohanram/go-webserver/internal/websocket"

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

	// Favicon
	r.HandleFunc("/favicon.ico", faviconHandler).Methods("GET")

	// REST API endpoints
	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/users", apiHandler.GetUsers).Methods("GET")
	apiRouter.HandleFunc("/users/{id}", apiHandler.GetUser).Methods("GET")
	apiRouter.HandleFunc("/users", apiHandler.CreateUser).Methods("POST")
	apiRouter.HandleFunc("/users/{id}", apiHandler.UpdateUser).Methods("PUT")
	apiRouter.HandleFunc("/users/{id}", apiHandler.DeleteUser).Methods("DELETE")

	// WebSocket endpoint
	r.HandleFunc("/ws/", wsHandler.HandleWebSocket)

	// Blog page
	r.HandleFunc("/blog", blogHandler).Methods("GET")

	// Static files for blog (CSS, JS, images)
	blogStaticFS := http.FileServer(http.Dir("web/blog/static"))
	r.PathPrefix("/blog/static/").Handler(http.StripPrefix("/blog/static/", blogStaticFS))

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
		Handler:   loggingMiddleware(r),
		TLSConfig: tlsConfig,
		ErrorLog:  log.New(&tlsErrorFilter{}, "", log.LstdFlags),
	}

	log.Printf("Server starting on port %s with TLS (HTTP/1.1 only)", port)
	log.Printf("Homepage: https://localhost%s", port)
	log.Printf("Blog: https://localhost%s/blog", port)
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

func blogHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "web/blog/index.html")
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	http.ServeFile(w, r, "web/favicon.svg")
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

// loggingMiddleware logs all HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the ResponseWriter to capture status code
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(lrw, r)

		// Log the request
		duration := time.Since(start)
		log.Printf("[ACCESS] %s %s %d %v %s",
			r.Method,
			r.RequestURI,
			lrw.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}

// loggingResponseWriter wraps http.ResponseWriter to capture status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Hijack implements http.Hijacker interface for WebSocket support
func (lrw *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := lrw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return hijacker.Hijack()
}

// tlsErrorFilter filters out TLS handshake errors from server logs
type tlsErrorFilter struct{}

func (f *tlsErrorFilter) Write(p []byte) (n int, err error) {
	msg := string(p)
	// Filter out TLS handshake errors but keep other errors
	if strings.Contains(msg, "TLS handshake error") {
		return len(p), nil
	}
	return os.Stderr.Write(p)
}
