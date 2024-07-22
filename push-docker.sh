#! /usr/bin/fish
docker build -t fractalhuman1/kessler-frontend:latest ./frontend/
docker build -t fractalhuman1/kessler-backend:latest ./backend/

docker push fractalhuman1/kessler-frontend:latest
docker push fractalhuman1/kessler-backend:latest
