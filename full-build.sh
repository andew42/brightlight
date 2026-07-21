#!/usr/bin/env bash
set -e

echo "Building frontend..."
cd frontend
npm run build
cd ..

echo "Building backend for Linux ARM (Raspberry Pi)..."
cd backend
GOOS=linux GOARCH=arm GOARM=5 go build -o ../deploy/brightlight
cd ..

echo "Copying frontend build to deploy..."
rsync -a --delete frontend/build/ deploy/frontend/build/

echo "Copying ui-config to deploy..."
rsync -a --delete backend/ui-config/ deploy/backend/ui-config/

echo "Done. Deploy directory is ready for scp to Pi."
