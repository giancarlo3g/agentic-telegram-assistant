.PHONY: build run test clean docker-build docker-run help

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the Go application"
	@echo "  run          - Run the application locally"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run with Docker Compose"
	@echo "  deps         - Download Go dependencies"

# Build the application
build:
	@echo "Building calendar assistant bot..."
	go build -o calendar-bot ./cmd/bot

# Run the application locally
run: build
	@echo "Running calendar assistant bot..."
	./calendar-bot

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -f calendar-bot
	go clean
	docker rmi calendar-assistant-bot

# Download dependencies
deps:
	@echo "Downloading Go dependencies..."
	go mod tidy
	go mod download

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t calendar-assistant-bot .

# Run with Docker Compose
docker-run:
	@echo "Starting services with Docker Compose..."
	docker-compose -f docker-compose.go.yml up -d

# Stop Docker Compose services
docker-stop:
	@echo "Stopping Docker Compose services..."
	docker-compose -f docker-compose.go.yml down

# View logs
view-logs:
	@echo "Viewing logs..."
	docker-compose -f docker-compose.go.yml logs -f

# Setup development environment
setup: deps
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then \
		echo "Creating .env file from template..."; \
		cp env.example .env; \
		echo "Please edit .env file with your credentials"; \
	else \
		echo ".env file already exists"; \
	fi
	@echo "Development environment setup complete!"
	@echo "Next steps:"
	@echo "1. Edit .env file with your credentials"
	@echo "2. Run 'make build' to build the application"
	@echo "3. Run 'make run' to start the bot locally"
	@echo "4. Or run 'make docker-run' to start with Docker Compose"
