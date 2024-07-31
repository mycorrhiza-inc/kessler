#! /bin/bash
docker build -t fractalhuman1/kessler-frontend:latest --platform linux/amd64 ./frontend/
docker build -t fractalhuman1/kessler-backend:latest --platform linux/amd64 ./backend/

docker push fractalhuman1/kessler-frontend:latest
docker push fractalhuman1/kessler-backend:latest
