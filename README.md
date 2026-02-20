# Go Web Server

This project demonstrates a complete web server implementation in Go following the standard Go project layout.

## Project Overview

This repository implements a Web server in Go with the following features:

- **Index/Homepage**: A simple Single Page Application using Bootstrap 5 templates
- **REST API**: Server at `/api/v1/` with endpoints for CRUD operations on users (in-memory data storage with JSON request/response handling)
- **WebSocket**: Real-time communication endpoint at `/ws/` for broadcasting messages to connected clients

## Project Structure

```
.
├── cmd/
│   └── server/          # Main application entry point
│       └── main.go
├── internal/            # Private application code
│   ├── api/            # REST API handlers
│   │   └── handlers.go
│   ├── store/          # In-memory data store
│   │   └── store.go
│   └── websocket/      # WebSocket functionality
│       ├── hub.go
│       └── handlers.go
├── web/
│   └── templates/      # HTML templates
│       └── index.html
├── bin/                # Compiled binaries (generated)
├── go.mod              # Go module file
├── Makefile            # Build automation
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.23 or higher
- Make (optional, for using Makefile commands)

### Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd go-webserver
```

2. Install dependencies:

```bash
go mod download
# or
make deps
```

### Running the Server

#### Using Make:

```bash
make run
```

#### Using Go directly:

```bash
go run ./cmd/server/main.go
```

#### Building and running:

```bash
make build
./bin/server
```

The server will start on `http://localhost:8080`

### Available Endpoints

#### Homepage

- `GET /` - Main page with Bootstrap 5 UI

#### REST API (Users)

- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/{id}` - Get a specific user
- `POST /api/v1/users` - Create a new user
- `PUT /api/v1/users/{id}` - Update a user
- `DELETE /api/v1/users/{id}` - Delete a user

#### WebSocket

- `WS /ws/` - WebSocket endpoint for real-time messaging

### API Examples

#### Create a new user:

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com"}'
```

#### Get all users:

```bash
curl http://localhost:8080/api/v1/users
```

#### Update a user:

```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"John Updated","email":"john.updated@example.com"}'
```

#### Delete a user:

```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## Make Commands

- `make build` - Build the server binary
- `make run` - Run the server directly
- `make clean` - Remove build artifacts
- `make test` - Run tests
- `make fmt` - Format code
- `make vet` - Run go vet
- `make deps` - Install dependencies
- `make dev` - Build and run the server
- `make help` - Show available commands

## Features

### REST API

- Full CRUD operations on user resources
- JSON request/response handling
- In-memory data storage with concurrent access protection
- Gorilla Mux for routing

### WebSocket

- Real-time bidirectional communication
- Message broadcasting to all connected clients
- Concurrent connection management
- Automatic connection cleanup
- **Automatic random messages**: Server sends random messages to connected clients at random intervals (1-5 seconds)

### Frontend

- Bootstrap 5 responsive UI
- Interactive user management
- Real-time chat interface
- AJAX API calls

## Technology Stack

- **Language**: Go 1.23
- **Router**: Gorilla Mux
- **WebSocket**: Gorilla WebSocket
- **Frontend**: Bootstrap 5, Vanilla JavaScript
- **Storage**: In-memory (with concurrent access protection)

## Environment Variables

- `PORT` - Server port (default: 8080)

Example:

```bash
PORT=3000 go run ./cmd/server/main.go
```

## License

This project is for demonstration purposes.
