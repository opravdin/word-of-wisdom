# Word of Wisdom Project Makefile

# Variables
SERVER_DIR = server
CLIENT_DIR = client
SERVER_BIN = word-of-wisdom-server
CLIENT_BIN = word-of-wisdom-client
BUILD_DIR = build

# Default target
.PHONY: all
all: build

# Help target
.PHONY: help
help:
	@echo "Word of Wisdom Project Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Common targets:"
	@echo "  build         Build both server and client"
	@echo "  run           Run both server and client using Docker Compose"
	@echo "  stop          Stop all services running in Docker Compose"
	@echo "  logs          Show logs from Docker Compose services"
	@echo "  test          Run all tests"
	@echo "  clean         Remove build artifacts"
	@echo "  help          Show this help message"

# Build targets
.PHONY: build
build:
	@echo "Building server..."
	@mkdir -p $(BUILD_DIR)
	@cd $(SERVER_DIR) && go build -o ../$(BUILD_DIR)/$(SERVER_BIN) ./cmd/service
	@echo "Building client..."
	@mkdir -p $(BUILD_DIR)
	@cd $(CLIENT_DIR) && go build -o ../$(BUILD_DIR)/$(CLIENT_BIN) ./cmd/client

# Docker targets
.PHONY: run stop logs
run:
	@echo "Running server and client with Docker Compose..."
	@docker-compose up -d

stop:
	@echo "Stopping all services..."
	@docker-compose down

logs:
	@echo "Showing logs from all services..."
	@docker-compose logs -f

# Test targets
.PHONY: test
test:
	@echo "Running all tests..."
	@cd $(SERVER_DIR) && go test -v ./...
	@cd $(CLIENT_DIR) && go test -v ./...


# Clean target
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@cd $(SERVER_DIR) && go clean
	@cd $(CLIENT_DIR) && go clean