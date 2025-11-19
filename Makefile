.PHONY: help list logs down restart dock build clean network-create setup-networks

BINARY_NAME=go-docker-manager

# Default target
help:
	@echo "Docker Container Manager"
	@echo "------------------------"
	@echo "Available commands:"
	@echo "  make help     - Show this help message"
	@echo "  make build    - Build the Go application"
	@echo "  make setup-networks - Create required docker networks (traefik-network, wordpress-network)"
	@echo "  make network-create NETWORK=name - Create a docker network"
	@echo "  make list     - List running containers"
	@echo "  make logs CONTAINER=name   - Show logs for a specific container"
	@echo "  make down CONTAINER=name   - Stop and remove a container"
	@echo "  make restart CONTAINER=name - Restart a container"
	@echo "  make dock CONTAINER=name TEMPLATE=template - Create and start a new container"
	@echo "  make clean    - Remove the compiled binary"

# Build the Go application
build:
	@echo "Building Docker Manager..."
	go build -o $(BINARY_NAME) main.go

# Check if binary exists, build if not
$(BINARY_NAME):
	@echo "Binary not found, building..."
	@$(MAKE) build

# List running containers
list: $(BINARY_NAME)
	@./$(BINARY_NAME) -command=list

# Show logs for a container
logs: $(BINARY_NAME)
	@if [ -z "$(CONTAINER)" ]; then \
		echo "Error: CONTAINER parameter is required"; \
		echo "Usage: make logs CONTAINER=name"; \
		exit 1; \
	fi
	@./$(BINARY_NAME) -command=logs -container=$(CONTAINER)

# Stop and remove a container
down: $(BINARY_NAME)
	@if [ -z "$(CONTAINER)" ]; then \
		echo "Error: CONTAINER parameter is required"; \
		echo "Usage: make down CONTAINER=name"; \
		exit 1; \
	fi
	@./$(BINARY_NAME) -command=down -container=$(CONTAINER)

# Restart a container
restart: $(BINARY_NAME)
	@if [ -z "$(CONTAINER)" ]; then \
		echo "Error: CONTAINER parameter is required"; \
		echo "Usage: make restart CONTAINER=name"; \
		exit 1; \
	fi
	@./$(BINARY_NAME) -command=restart -container=$(CONTAINER)

# Create and start a new container
dock: $(BINARY_NAME)
	@if [ -z "$(CONTAINER)" ]; then \
		echo "Error: CONTAINER parameter is required"; \
		echo "Usage: make dock CONTAINER=name TEMPLATE=template"; \
		exit 1; \
	fi
	@if [ -z "$(TEMPLATE)" ]; then \
		echo "Error: TEMPLATE parameter is required"; \
		echo "Usage: make dock CONTAINER=name TEMPLATE=template"; \
		exit 1; \
	fi
	@./$(BINARY_NAME) -command=dock -container=$(CONTAINER) -template=$(TEMPLATE)

# Clean up compiled binary
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)

# Create a docker network
network-create: $(BINARY_NAME)
	@if [ -z "$(NETWORK)" ]; then \
		echo "Error: NETWORK parameter is required"; \
		echo "Usage: make network-create NETWORK=name"; \
		exit 1; \
	fi
	@./$(BINARY_NAME) -command=network-create -network=$(NETWORK)

# Setup required networks for the project
setup-networks: $(BINARY_NAME)
	@echo "Creating required docker networks..."
	@./$(BINARY_NAME) -command=network-create -network=traefik-network
	@./$(BINARY_NAME) -command=network-create -network=wordpress-network
	@echo "All required networks are ready!"