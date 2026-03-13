#!/bin/bash

set -e

echo "Building backend..."
cd backend
go mod tidy
go build -o bin/server ./cmd/server

echo "Building frontend web..."
cd ../frontend/web
npm install
npm run build

echo "Build completed!"
