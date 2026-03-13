#!/bin/bash

set -e

echo "Deploying with Docker Compose..."
cd deploy/docker
docker-compose up -d --build

echo "Deployment completed!"
echo "Backend API: http://localhost:8080"
echo "Frontend Web: http://localhost:80"
