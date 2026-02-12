#!/bin/bash

# Scalable Learning Platform - Quick Start Script

echo "🚀 Scalable Learning Platform - Quick Start"
echo "==========================================="

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

echo "✅ Docker and Docker Compose are installed"
echo ""

# Check if .env file exists
if [ ! -f ".env" ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env
    echo "✅ .env file created. Please review and update if needed."
fi

echo ""
echo "🐳 Starting Docker Compose services..."
echo ""

# Start services
docker-compose up

echo ""
echo "✅ Services are running!"
echo ""
echo "📍 Access the application at:"
echo "   Frontend: http://localhost:5173"
echo "   Backend:  http://localhost:8080"
echo "   Database: localhost:5432"
echo ""
echo "To stop services, press Ctrl+C"
