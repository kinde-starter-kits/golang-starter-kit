#!/bin/bash
set -e

echo "================================================"
echo "🔧 Kinde Golang Starter Kit - Installation"
echo "================================================"
echo ""


echo "📦 Step 1: Cleaning up old files..."
rm -f go.sum

echo "✅ Step 2: Downloading dependencies..."
go get github.com/gin-contrib/sessions
go get github.com/gin-contrib/sessions/cookie  
go get github.com/gin-gonic/gin
go get github.com/joho/godotenv
go get golang.org/x/oauth2

echo "✅ Step 3: Tidying go.mod..."
go mod tidy

echo "✅ Step 4: Verifying installation..."
go mod download

echo ""
echo "================================================"
echo "✅ Installation complete!"
echo "================================================"
echo ""
echo "🚀 Starting the server..."
echo "Visit: http://localhost:3000"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

# Run the application
go run main.go


