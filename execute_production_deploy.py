import argparse
import os
import sys
import subprocess


def parse_arguments():
    """Check if the script is being run as root."""
    if os.geteuid() != 0:
        print("This script must be run as root")
        sys.exit(1)
    """Parse command line arguments."""
    parser = argparse.ArgumentParser(
        description="Generate Docker Compose configuration"
    )
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument(
        "--production", action="store_true", help="Generate production configuration"
    )
    group.add_argument(
        "--nightly", action="store_true", help="Generate nightly configuration"
    )
    parser.add_argument(
        "--version", default="prod", help="Version hash of the containers"
    )
    return parser.parse_args()


def get_config_settings(args):
    """Determine configuration settings based on arguments."""
    template = COMPOSE_TEMPLATE
    if args.production:
        domain = "kessler.xyz"
        public_api_url = "api.kessler.xyz"
    else:
        domain = "nightly.kessler.xyz"
        public_api_url = f"nightly-api.kessler.xyz"

    https_public_api_url = f"https://{public_api_url}"
    return {
        "version_hash": args.version,
        "domain": domain,
        "public_api_url": public_api_url,
        "https_public_api_url": https_public_api_url,
        "template": template,
    }


def generate_docker_compose_content(config):
    """Generate the Docker Compose file content by replacing variables in the template."""
    return config["template"].format(**config)


def deploy_configuration(template_content):
    """Deploy the Docker Compose configuration."""
    subprocess.run(["git", "reset", "--hard", "HEAD"], check=True)
    subprocess.run(["git", "clean", "-fd"], check=True)
    subprocess.run(["git", "pull"], check=True)
    # Write template content to file
    with open("docker-compose.deploy.yaml", "w") as f:
        f.write(template_content)

    # Stop and start services
    subprocess.run(
        ["podman-compose", "-f", "docker-compose.deploy.yaml", "down"], check=True
    )
    subprocess.run(
        ["podman-compose", "-f", "docker-compose.deploy.yaml", "up"], check=True
    )


def main():
    """Main function that coordinates all operations."""

    args = parse_arguments()

    config = get_config_settings(args)

    docker_compose_content = generate_docker_compose_content(config)

    deploy_configuration(docker_compose_content)


COMPOSE_TEMPLATE = """
x-common-env: &common-env
  OTEL_EXPORTER_OTLP_ENDPOINT: otel-collector:4317
  VERSION_HASH: {version_hash}
  DOMAIN: "{domain}"
  PUBLIC_KESSLER_API_URL: "{https_public_api_url}"
  INTERNAL_KESSLER_API_URL: "http://backend-go:4041"
  INTERNAL_REDIS_ADDRESS: "valkey:6379"

services:
  reverse-proxy:
    image: traefik:v3.0
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--providers.docker.constraints=Label(`traefik.namespace`,`kessler`)"
      - "--providers.docker.exposedbydefault=false"
      - "--entryPoints.websecure.address=:443"
      - "--certificatesresolvers.myresolver.acme.tlschallenge=true"
      - "--certificatesresolvers.myresolver.acme.email=mbright@kessler.xyz"
      - "--certificatesresolvers.myresolver.acme.storage=/letsencrypt/acme.json"
      - "--providers.file.filename=/etc/traefik/traefik_dynamic.yaml"
    ports:
      - "80:80"
      - "443:443"
      - "8083:8080"
    volumes:
      - "./volumes/letsencrypt:/letsencrypt"
      - "/run/podman/podman.sock:/var/run/docker.sock:ro"
      - "./traefik_dynamic.yaml:/etc/traefik/traefik_dynamic.yaml:ro"

  frontend:
    image: fractalhuman1/kessler-frontend:{version_hash}
    env_file:
      - config/global.env
    environment:
      <<: *common-env
    volumes:
      - ./config/frontend.env.local:/app.env.local
      - ./config/global.env:/app/.env
    labels:
      - "traefik.enable=true"
      - "traefik.namespace=kessler"
      - "traefik.http.routers.frontend.rule=Host(`{domain}`) && PathPrefix(`/`)"
      - "traefik.http.routers.blog.tls.domains[0].main={domain}"
      - "traefik.http.routers.frontend.entrypoints=websecure"
      - "traefik.http.routers.frontend.tls.certresolver=myresolver"
      - "traefik.http.routers.whoami.rule=Host(`{domain}`)"
      - "traefik.http.routers.whoami.entrypoints=websecure"
      - "traefik.http.routers.whoami.tls.certresolver=myresolver"
    expose:
      - 3000
    command:
      - "npm"
      - "run"
      - "start"

  backend-go:
    image: fractalhuman1/kessler-backend-go:{version_hash}
    command:
      - "./kessler-backend-go"
    env_file:
      - config/global.env
    environment:
      <<: *common-env
    expose:
      - 4041
    labels:
      - "traefik.enable=true"
      - "traefik.namespace=kessler"
      - "traefik.http.routers.backend-go.rule=Host(`{public_api_url}`) && PathPrefix(`/v2`)"
      - "traefik.http.routers.backend-go.entrypoints=websecure"
      - "traefik.http.routers.backend-go.tls.certresolver=myresolver"

  ingest:
    image: fractalhuman1/kessler-ingest:{version_hash}
    command:
      - "./kessler-ingest"
    env_file:
      - config/global.env
    environment:
      <<: *common-env
    expose:
      - 4042
    labels:
      - "traefik.enable=true"
      - "traefik.namespace=kessler"
      - "traefik.http.routers.ingest_v1.rule=Host(`{public_api_url}`) && PathPrefix(`/ingest_v1`)"
      - "traefik.http.routers.ingest_v1.entrypoints=websecure"
      - "traefik.http.routers.ingest_v1.tls.certresolver=myresolver"
  valkey:
    hostname: valkey
    image: valkey/valkey:7.2.5
    volumes:
      - ./volumes/valkey.conf:/etc/valkey/valkey.conf
      - ./volumes/valkey-data:/data
    command: valkey-server /etc/valkey/valkey.conf
    healthcheck:
      test: ["CMD-SHELL", "redis-cli ping | grep PONG"]
      interval: 1s
      timeout: 3s
      retries: 5
    expose:
      - 6379
""".strip()


if __name__ == "__main__":
    main()
