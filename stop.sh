#!/bin/bash

# Stop and clean up
echo "🛑 Stopping services..."
docker-compose down

echo "🧹 Removing images..."
docker-compose down --rmi all

echo "✅ Cleanup complete!"
