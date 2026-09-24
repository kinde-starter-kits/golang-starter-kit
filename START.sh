#!/bin/bash

echo "🔧 Installing Go dependencies..."

# Download all dependencies
go mod download

# Tidy up go.mod and create go.sum
go mod tidy

echo ""
echo "✅ Dependencies installed successfully!"
echo ""
echo "🚀 Starting the Kinde Golang Starter Kit..."
echo ""

# Run the application
go run main.go


