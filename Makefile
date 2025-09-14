# Google Indexing API Makefile

.PHONY: build run test clean docker-build docker-run help

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=google-indexing-api
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/main.go

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/main.go
	./$(BINARY_NAME)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f coverage.out

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Build for Linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/main.go

# Docker build
docker-build:
	docker build -t $(BINARY_NAME) .

# Docker run
docker-run:
	docker run -d --name $(BINARY_NAME) -p 8080:8080 $(BINARY_NAME)

# Docker compose up
compose-up:
	docker-compose up -d

# Docker compose down
compose-down:
	docker-compose down

# Format code
fmt:
	$(GOCMD) fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Check for security issues (requires gosec)
security:
	gosec ./...

# Install development tools
install-tools:
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint
	$(GOGET) -u github.com/securecodewarrior/gosec/v2/cmd/gosec

# Development setup
dev-setup: deps install-tools
	cp .env.example .env
	@echo "Please edit .env file with your configuration"

# Help
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean         - Clean build files"
	@echo "  deps          - Download dependencies"
	@echo "  build-linux   - Build for Linux"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run Docker container"
	@echo "  compose-up    - Start with Docker Compose"
	@echo "  compose-down  - Stop Docker Compose"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  security      - Check for security issues"
	@echo "  install-tools - Install development tools"
	@echo "  dev-setup     - Setup development environment"
	@echo "  help          - Show this help message"
