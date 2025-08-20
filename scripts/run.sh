#!/bin/bash

# Development run script for POS System

set -e

echo "🚀 Starting POS System Development Environment..."

# Check if .env exists
if [ ! -f .env ]; then
    echo "📝 Creating .env from .env.example..."
    cp .env.example .env
fi

# Start infrastructure services
echo "🐳 Starting infrastructure services..."
docker-compose up -d postgres redis rabbitmq

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if services are healthy
echo "🔍 Checking service health..."
docker-compose ps

# Build and run the API
echo "🔨 Building API..."
go build -o bin/api cmd/api/main.go

echo "🌟 Starting API server..."
./bin/api

echo "✅ POS System is running at http://localhost:8080"