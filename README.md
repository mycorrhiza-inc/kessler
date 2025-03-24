```
░▒▓█▓▒░░▒▓█▓▒░▒▓████████▓▒░░▒▓███████▓▒░▒▓███████▓▒░▒▓█▓▒░      ░▒▓████████▓▒░▒▓███████▓▒░
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░     ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░     ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░
░▒▓███████▓▒░░▒▓██████▓▒░  ░▒▓██████▓▒░░▒▓██████▓▒░░▒▓█▓▒░      ░▒▓██████▓▒░ ░▒▓███████▓▒░
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░             ░▒▓█▓▒░     ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░             ░▒▓█▓▒░     ░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░
░▒▓█▓▒░░▒▓█▓▒░▒▓████████▓▒░▒▓███████▓▒░▒▓███████▓▒░░▒▓████████▓▒░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░
```

search, for goverment stuff, but actually good. Check it out at [kessler.xyz](https://kessler.xyz) or our nightly branch [nightly.kessler.xyz](https://nightly.kessler.xyz).

# Contributions

We are always looking for contributors! If you want to contribute to the project feel free to read CONTRIBUTING.md for more info or just make an issue or pull request!


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

# Running the application 

To update the nightly enviornment execute the following script:

`./update-docker-compose.sh`

This will go ahead and download the repo to your computer at /mycorrhiza/kessler. Navigate to the latest version of the main branch. Build the containers, and push them under the tag kessler-imagetype:commit-hash-of-main. Then execute a remote command to change the deployment at kessler.xyz to whatever commit you just pushed. (this does require root access on said remote)

If you want to specify a certain commit you can run 

`./update-docker-compose.sh --commit <commit-hash>`

or to push whatever commit your dev environment is on run

```bash
./update-docker-compose.sh --commit $(git rev-parse HEAD)
```

To update prod instead of nightly run 
`./update-docker-compose.sh --prod`
or
`./update-docker-compose.sh --commit <commit-hash> --prod`

To set prod to a specific version. This does work with rollbacks

IMPORTANT: Its recommended to run this command on an x86 machine. Cross compilation of docker issues to x86 has been known to take a long time.

By default the remote is set to `kessler.xyz` to set a different remote use 

`./update-docker-compose.sh --commit <commit-hash>`



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

```bash
docker compose up -d --no-deps --build frontend
```

and

```bash
docker compose up -d --no-deps --build backend
```


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
