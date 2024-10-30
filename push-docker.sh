#! /bin/bash
# IF YOU ARE ON AN ARM SYSTEM AND ITS TAKING A LONG TIME TO BUILD:
# 1. Ssh into amon-sul through tailscale.
# 2. sudo fish
# 3. cd into the kessler repo at /mycorrhiza/kessler 
# 4. git pull 
# 5. Run this script
set -e

docker build -t fractalhuman1/kessler-frontend:latest --platform linux/amd64 ./frontend/ &&
docker build -t fractalhuman1/kessler-backend-go:latest --platform linux/amd64 ./backend-go/ &&

docker push fractalhuman1/kessler-frontend:latest &&
docker push fractalhuman1/kessler-backend-go:latest 
