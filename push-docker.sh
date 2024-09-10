#! /bin/bash
docker build -t fractalhuman1/kessler-frontend:latest --platform linux/amd64 ./frontend/
docker build -t fractalhuman1/kessler-backend-python:latest --platform linux/amd64 ./backend-python/
docker build -t fractalhuman1/kessler-backend-go:latest --platform linux/amd64 ./backend-go/
# docker build -t fractalhuman1/kessler-thaumaturgy-python:latest --platform linux/amd64 ./thaumaturgy-python/

docker push fractalhuman1/kessler-frontend:latest
docker push fractalhuman1/kessler-backend-python:latest
docker push fractalhuman1/kessler-backend-go:latest 
# docker push fractalhuman1/kessler-thaumaturgy-python:latest
