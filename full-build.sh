#!/usr/bin/env bash
# Local equivalent of the CI pipeline (.github/workflows/build.yml):
# builds the frontend and backend the same way and stages the artefacts
# in deploy/ with the same layout as the CI Pi bundle.
set -e

echo "Building frontend..."
cd frontend
npm ci
npm run build
cd ..

echo "Vetting backend..."
cd backend
go vet ./...

echo "Building backend for Raspberry Pi 2B (linux/arm v7)..."
GOOS=linux GOARCH=arm GOARM=7 go build -o ../deploy/brightlight .
cd ..

echo "Copying frontend build to deploy..."
rsync -a --delete frontend/build/ deploy/frontend/build/

echo "Copying ui-config to deploy..."
rsync -a --delete backend/ui-config/ deploy/backend/ui-config/

echo "Copying systemd unit and installer to deploy..."
cp packaging/brightlight.service packaging/install.sh deploy/

echo "Done. Deploy directory matches the CI bundle layout, ready for scp to Pi."
