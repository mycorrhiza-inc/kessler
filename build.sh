#!/bin/bash

usage() {
  echo "Usage: [-d | --dev : development build ] [-p | --prod : production build]" 1>&2
  exit 0
}

dev() {
  COMPOSE_BAKE=true
  docker compose up --build --watch --remove-orphans
}
prod() {

  local current_hash=$(git rev-parse HEAD)
  echo "Current commit hash: $current_hash"

  sudo docker build -t "fractalhuman1/kessler-frontend:${current_hash}" --platform linux/amd64 --file ./frontend/prod.Dockerfile ./frontend/
  sudo docker build -t "fractalhuman1/kessler-backend-server:${current_hash}" --platform linux/amd64 --file ./backend/prod.server.Dockerfile ./backend
  sudo docker build -t "fractalhuman1/kessler-backend-ingest:${current_hash}" --platform linux/amd64 --file ./backend/prod.ingest.Dockerfile ./backend

  sudo docker push "fractalhuman1/kessler-frontend:${current_hash}"
  sudo docker push "fractalhuman1/kessler-backend-server:${current_hash}"
  sudo docker push "fractalhuman1/kessler-backend-ingest:${current_hash}"
}

# Initialize flags
build=false
prod=false

# Parse arguments
while [[ "$#" -gt 0 ]]; do
  case "$1" in
  --dev)
    build=true
    ;;
  --prod)
    prod=true
    ;;
  *)
    echo "Unknown option: $1"
    exit 1
    ;;
  esac
  shift
done

# Ensure only one of --build or --prod is used
if [ "$build" = true ] && [ "$prod" = true ]; then
  echo "Error: You cannot use --dev and --prod together."
  exit 1
elif [ "$build" = true ]; then
  echo "Running dev build..."
  dev
  # Add build commands here
elif [ "$prod" = true ]; then
  echo "Running production build..."
  prod
  # Add production commands here
else
  echo "Error: You must provide either --dev or --prod."
  exit 1
fi
