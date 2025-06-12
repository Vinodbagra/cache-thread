#!/bin/bash

# Cache API Test Runner Script
# This script runs the comprehensive API test suite

echo "🧪 Cache API Test Runner"
echo "========================"

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed or not in PATH"
    exit 1
fi

# Check if we're in the right directory
if [ ! -f "test_apis.go" ]; then
    echo "❌ Error: test_apis.go not found. Please run this script from the cmd/test directory"
    exit 1
fi

# Check if server is running (optional check)
echo "🔍 Checking if server is running on http://localhost:8080..."
if curl -s http://localhost:8080/api/cache/health > /dev/null 2>&1; then
    echo "✅ Server is running"
else
    echo "⚠️  Warning: Server might not be running. Make sure to start it with:"
    echo "   go run cmd/api/main.go"
    echo ""
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ Test cancelled"
        exit 1
    fi
fi

echo ""
echo "🚀 Starting API tests..."
echo ""

# Run the tests
if go run test_apis.go; then
    echo ""
    echo "✅ Tests completed successfully!"
    exit 0
else
    echo ""
    echo "❌ Tests failed!"
    exit 1
fi 