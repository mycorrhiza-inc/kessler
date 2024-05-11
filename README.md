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

## Running the application

To start the application run `docker compose up -d`

Once up, the Kessler API can be found at backend.docker.localhost or
`localhost:5055`

# Development

## Setting Up your Dev Environment

### Environment Variables

Copy the environment files:

_TODO : Refactor to throw all the enviornment variables in 1 file_

```bash
cp .env.example .env
cp frontend/.env.local.example frontend/.env.local

```

And fill them out

You will need `docker` and `docker compose`

### Chroma

To run chroma on a non-local basis you will need to set up basic auth for chroma's docker container as specified in the [Authentication Documentation](https://docs.trychroma.com/usage-guide#authentication).

At the root of the project run

```bash
htpasswd -Bbn admin admin > server.htpasswd
```

and in the `.env` file add:

```bash
export CHROMA_SERVER_AUTH_CREDENTIALS_FILE="server.htpasswd"
export CHROMA_SERVER_AUTH_CREDENTIALS_PROVIDER="chromadb.auth.providers.HtpasswdFileServerAuthCredentialsProvider"
export CHROMA_SERVER_AUTH_PROVIDER="chromadb.auth.basic.BasicAuthServerProvider"
```

### Running it

Currently the Kessler docker compose is in dev mode.

To run it:

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

All volumes are stored in the volumes folder off of main, the three volumes are currently mounted are

`./volumes/files`
`./volumes/tmp`
`./volumes/instance`

The last volume is where you should copy over database files to when running it on the backend.
