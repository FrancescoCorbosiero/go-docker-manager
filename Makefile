.PHONY: help list logs down restart dock build

# Default target
help:
	@echo "Docker Container Manager"
	@echo "------------------------"
	@echo "Available commands:"
	@echo "  make help     - Show this help message"
	@echo "  make list     - List running containers"
	@echo "  make logs CONTAINER=name   - Show logs for a specific container"
	@echo "  make down CONTAINER=name   - Stop and remove a container"
	@echo "  make restart CONTAINER=name - Restart a container"
	@echo "  make dock CONTAINER=name TEMPLATE=template - Create and start a new container"
	@echo "  make build    - Build the Go application"

# Build the Go application
build:
	@echo "Building Docker Manager..."
	go build -o go-docker-manager main.go

# List running containers
list:
	@./go-docker-manager -command=list

# Show logs for a container
logs:
	@if [ -z "$(CONTAINER)" ]; then \
		echo "Error: CONTAINER parameter is required"; \
		echo "Usage: make logs CONTAINER=name"; \
		exit 1; \
	fi
	@./go-docker-manager -command=logs -container=$(CONTAINER)

# Stop and remove a container
down:
	@if [ -z "$(CONTAINER)" ]; then \
		echo "Error: CONTAINER parameter is required"; \
		echo "Usage: make down CONTAINER=name"; \
		exit 1; \
	fi
	@./go-docker-manager -command=down -container=$(CONTAINER)

# Restart a container
restart:
	@if [ -z "$(CONTAINER)" ]; then \
		echo "Error: CONTAINER parameter is required"; \
		echo "Usage: make restart CONTAINER=name"; \
		exit 1; \
	fi
	@./go-docker-manager -command=restart -container=$(CONTAINER)

# Create and start a new container
dock:
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
	@./go-docker-manager -command=dock -container=$(CONTAINER) -template=$(TEMPLATE)