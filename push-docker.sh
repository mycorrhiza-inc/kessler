#! /usr/bin/fish
docker build -t fractalhuman1/kessler-frontend ./frontend/
docker build -t fractalhuman1/kessler-backend ./backend/

docker push fractalhuman1/kessler-frontend
docker push fractalhuman1/kessler-backend
