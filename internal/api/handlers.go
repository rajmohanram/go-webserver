package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rajmohanram/go-webserver/internal/store"

	"github.com/gorilla/mux"
)

// Handler manages API requests
type Handler struct {
	store *store.Store
}

// NewHandler creates a new API handler
func NewHandler(s *store.Store) *Handler {
	return &Handler{store: s}
}

// GetUsers returns all users
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := h.store.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUser returns a specific user by ID
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := 0
	fmt.Sscanf(vars["id"], "%d", &id)

	user, exists := h.store.Get(id)
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateUser creates a new user
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := h.store.Create(input.Name, input.Email)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser updates an existing user
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := 0
	fmt.Sscanf(vars["id"], "%d", &id)

	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, exists := h.store.Update(id, input.Name, input.Email)
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DeleteUser deletes a user
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := 0
	fmt.Sscanf(vars["id"], "%d", &id)

	if !h.store.Delete(id) {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
