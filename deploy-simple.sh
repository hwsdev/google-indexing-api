#!/bin/bash

# Simple Google Indexing API Deployment Script
set -e

APP_NAME="google-indexing-api"

echo "🚀 Starting deployment of $APP_NAME..."

# Stop existing containers
echo "📦 Stopping existing containers..."
docker-compose down || true

# Build and start
echo "🔨 Building and starting containers..."
docker-compose up -d --build

# Wait for health check
echo "⏳ Waiting for application to be ready..."
sleep 10

# Check if running
if docker-compose ps | grep -q "Up"; then
    echo "✅ $APP_NAME deployed successfully!"
    echo "🌐 Application is running on http://localhost:8080"
    echo "🏥 Health check: http://localhost:8080/api/health"
else
    echo "❌ Deployment failed!"
    docker-compose logs
    exit 1
fi
