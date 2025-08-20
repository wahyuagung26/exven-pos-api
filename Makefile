.PHONY: build run test clean docker-build docker-run dev

# Variables
APP_NAME = pos-api
MAIN_PATH = cmd/api/main.go
BUILD_PATH = bin/$(APP_NAME)

# Default target
dev: build run

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_PATH) $(MAIN_PATH)
	@echo "Build completed: $(BUILD_PATH)"

# Run the application
run:
	@echo "Starting $(APP_NAME)..."
	@./$(BUILD_PATH)

# Run with hot reload (requires air)
watch:
	@echo "Starting development server with hot reload..."
	@air

# Test the application
test:
	@echo "Running tests..."
	@go test -v ./...

# Test with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -cover ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf tmp/

# Generate mocks (requires mockery)
mocks:
	@echo "Generating mocks..."
	@go generate ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Vet code
vet:
	@echo "Vetting code..."
	@go vet ./...

# Lint code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@golangci-lint run

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

# Docker build
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME) .

# Docker run
docker-run:
	@echo "Running with Docker Compose..."
	@docker-compose up --build

# Start infrastructure services only
infra:
	@echo "Starting infrastructure services..."
	@docker-compose up -d postgres redis rabbitmq

# Stop all services
stop:
	@echo "Stopping all services..."
	@docker-compose down

# Database migration
migrate-up:
	@echo "Running database migrations..."
	@go run cmd/migration/main.go up

# Database migration rollback
migrate-down:
	@echo "Rolling back database migrations..."
	@go run cmd/migration/main.go down

# Database reset
migrate-reset:
	@echo "Resetting database..."
	@go run cmd/migration/main.go reset

# Setup development environment
setup: infra migrate-up
	@echo "Development environment ready!"

# Full development workflow
dev-full: clean tidy build setup test
	@echo "Full development setup complete!"

# Production build
prod-build:
	@echo "Building for production..."
	@CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o $(BUILD_PATH) $(MAIN_PATH)

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application"
	@echo "  dev           - Build and run (default)"
	@echo "  watch         - Run with hot reload"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  mocks         - Generate mocks"
	@echo "  fmt           - Format code"
	@echo "  vet           - Vet code"
	@echo "  lint          - Lint code"
	@echo "  tidy          - Tidy dependencies"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker Compose"
	@echo "  infra         - Start infrastructure services"
	@echo "  stop          - Stop all services"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback migrations"
	@echo "  migrate-reset - Reset database"
	@echo "  setup         - Setup development environment"
	@echo "  dev-full      - Full development setup"
	@echo "  prod-build    - Build for production"
	@echo "  help          - Show this help"