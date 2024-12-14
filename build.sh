#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Build the React frontend
echo "Building frontend..."
cd frontend
npm install
npm run build

# Copy the frontend build to the backend static directory
echo "Copying build files to backend..."
cd ..
rm -rf backend/public
mkdir -p backend/public
cp -r frontend/dist/* backend/public/

echo "Build process completed successfully."
