package store

import (
	"sync"
	"time"
)

// User represents a user resource for CRUD operations
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Store holds in-memory data and manages concurrent access
type Store struct {
	mu     sync.RWMutex
	users  map[int]User
	nextID int
}

// NewStore creates a new Store instance with sample data
func NewStore() *Store {
	s := &Store{
		users:  make(map[int]User),
		nextID: 1,
	}
	s.initializeData()
	return s
}

func (s *Store) initializeData() {
	s.users[1] = User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
	}
	s.users[2] = User{
		ID:        2,
		Name:      "Jane Smith",
		Email:     "jane@example.com",
		CreatedAt: time.Now(),
	}
	s.nextID = 3
}

// GetAll returns all users
func (s *Store) GetAll() []User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

// Get returns a user by ID
func (s *Store) Get(id int) (User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	return user, exists
}

// Create adds a new user
func (s *Store) Create(name, email string) User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := User{
		ID:        s.nextID,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}
	s.users[user.ID] = user
	s.nextID++
	return user
}

// Update modifies an existing user
func (s *Store) Update(id int, name, email string) (User, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.users[id]
	if !exists {
		return User{}, false
	}

	user.Name = name
	user.Email = email
	s.users[id] = user
	return user, true
}

// Delete removes a user
func (s *Store) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[id]; !exists {
		return false
	}

	delete(s.users, id)
	return true
}
