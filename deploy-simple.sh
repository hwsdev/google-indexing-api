#!/bin/bash

# Simple deployment script for Google Indexing API

set -e

echo "🚀 Starting deployment..."

# Build and start with docker-compose
echo "📦 Building and starting containers..."
docker-compose up --build -d

echo "⏳ Waiting for service to be ready..."
sleep 10

# Check if service is running
if curl -f http://localhost:8080/api/health > /dev/null 2>&1; then
    echo "✅ Service is running successfully!"
    echo "🌐 API available at: http://localhost:8080"
    echo ""
    echo "📚 Endpoints:"
    echo "  Health Check: GET  http://localhost:8080/api/health"
    echo "  Single URL:   POST http://localhost:8080/api/v1/index"
    echo "  Batch URLs:   POST http://localhost:8080/api/v1/index/batch"
    echo "  URL Status:   GET  http://localhost:8080/api/v1/status/{url}"
    echo "  Cache Stats:  GET  http://localhost:8080/api/v1/cache/stats"
    echo "  Clear Cache:  POST http://localhost:8080/api/v1/cache/clear"
else
    echo "❌ Service failed to start properly"
    echo "📋 Checking logs..."
    docker-compose logs
    exit 1
fi
