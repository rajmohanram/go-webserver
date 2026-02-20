.PHONY: build run clean test fmt vet

# Build the server
build:
	@echo "Building server..."
	@go build -o bin/server ./cmd/server
	@GOOS=linux GOARCH=amd64 go build -o bin/server-linux-amd64 ./cmd/server
	@GOOS=linux GOARCH=arm64 go build -o bin/server-linux-arm64 ./cmd/server

# Run the server
run:
	@echo "Running server..."
	@go run ./cmd/server/main.go

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Build and run
dev: build
	@./bin/server

# Show help
help:
	@echo "Available commands:"
	@echo "  make build    - Build the server binary"
	@echo "  make run      - Run the server directly"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make test     - Run tests"
	@echo "  make fmt      - Format code"
	@echo "  make vet      - Run go vet"
	@echo "  make deps     - Install dependencies"
	@echo "  make dev      - Build and run the server"
	@echo "  make help     - Show this help message"
