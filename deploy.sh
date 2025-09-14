#!/bin/bash

# Simple Google Indexing API Deployment Script
# Usage: ./deploy.sh [start|stop|restart|logs|build]

set -e

APP_NAME="google-indexing-api"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Build application
build() {
    info "Building $APP_NAME..."
    docker build -t $APP_NAME .
    info "Build completed!"
}

# Start application
start() {
    info "Starting $APP_NAME..."
    
    # Create .env if not exists
    if [ ! -f ".env" ]; then
        warn ".env not found, creating from .env.example..."
        cp .env.example .env
    fi
    
    # Start with docker-compose
    docker-compose up -d
    
    # Wait and check
    sleep 5
    if docker-compose ps | grep -q "Up"; then
        info "$APP_NAME started successfully!"
        info "Access: http://localhost:8080"
        info "Health: http://localhost:8080/api/health"
    else
        error "Failed to start $APP_NAME"
        docker-compose logs
    fi
}

# Stop application
stop() {
    info "Stopping $APP_NAME..."
    docker-compose down
    info "$APP_NAME stopped!"
}

# Restart application
restart() {
    stop
    start
}

# Show logs
logs() {
    info "Showing logs..."
    docker-compose logs -f
}

# Install Docker (Ubuntu)
install_docker() {
    info "Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker $USER
    
    # Install Docker Compose
    sudo curl -L "https://github.com/docker/compose/releases/download/v2.21.0/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    sudo chmod +x /usr/local/bin/docker-compose
    
    info "Docker installed! Please logout and login again."
}

# Main script
case "${1:-start}" in
    build)
        build
        ;;
    start)
        start
        ;;
    stop)
        stop
        ;;
    restart)
        restart
        ;;
    logs)
        logs
        ;;
    install-docker)
        install_docker
        ;;
    *)
        echo "Usage: $0 {build|start|stop|restart|logs|install-docker}"
        echo ""
        echo "Commands:"
        echo "  build         - Build Docker image"
        echo "  start         - Start application"
        echo "  stop          - Stop application"
        echo "  restart       - Restart application"
        echo "  logs          - Show logs"
        echo "  install-docker- Install Docker (Ubuntu)"
        ;;
esac
