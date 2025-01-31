```
░▒▓█▓▒░░▒▓█▓▒░▒▓████████▓▒░░▒▓███████▓▒░▒▓███████▓▒░▒▓█▓▒░      ░▒▓████████▓▒░▒▓███████▓▒░
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░     ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░     ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░
░▒▓███████▓▒░░▒▓██████▓▒░  ░▒▓██████▓▒░░▒▓██████▓▒░░▒▓█▓▒░      ░▒▓██████▓▒░ ░▒▓███████▓▒░
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░             ░▒▓█▓▒░     ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░             ░▒▓█▓▒░     ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░
░▒▓█▓▒░░▒▓█▓▒░▒▓████████▓▒░▒▓███████▓▒░▒▓███████▓▒░░▒▓████████▓▒░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░
```

search, sorta

> warning!!: this code is pre-deployment. Use with anything you care about at your own peril

# Kubernetes Usage

WARNING k8s is still experiencing some bugs, bewarned 

use 

`./update-kubernetes.sh --nigthly-commit <hash>`



TODO: UPDATE DOCUMENTATION FOR EVERYTHING 

# Usage

## Prod

```bash
cd frontend
npm run build
cd ..
docker compose up
```

## Requirements

- docker
- docker compose

## Running the application (production)

To start the application run `docker compose --env-file ./config/global.env up`

# Development

## Setting Up your Dev Environment

### Environment Variables

Copy the environment files:

```bash
cp -r example-config config
```

And fill each config file out

You will need `docker` and `docker compose`

### Running it

To run the app in dev mode specify the dev compose

```bash
docker compose -f docker-compose.dev.yml up --force-recreate
```

### Adding packages

Adding packages to the frontend or the backend requires a rebuild for the
respective container.

```
docker compose up -d --no-deps --build frontend
```

and

```
docker compose up -d --no-deps --build backend
```

# Contribute

the software is currently under heavy development. No outside contributions will be accepted at this time.

# Volumes and Storage

All volumes are stored by default in the volumes folder off of main, this can be changed by changing the `VOLUMES_DIRECTORY` environment variable in ./config/global.env , the three volumes are currently mounted are

`${VOLUMES_DIRECTORY:-./volumes}/files`
`${VOLUMES_DIRECTORY:-./volumes}/tmp`
`${VOLUMES_DIRECTORY:-./volumes}/instance`

The last volume is where you should copy over database files to when running it on the backend.

# Debugging

If running into any weird issues with the software run these 2 commands first

```
docker rm $(docker ps -a -q) && docker rmi $(docker images -a -q) && docker system prune -a
```
