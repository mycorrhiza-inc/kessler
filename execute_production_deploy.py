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
    template = ENV_FILE_TEMPLATE
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
    with open("deploy.env", "w") as f:
        f.write(template_content)

    # Stop and start services
    subprocess.run(
        ["podman-compose", "--env-file", "deploy.env", "-f", "prod.deploy.yaml", "up"],
        check=True,
    )


def main():
    """Main function that coordinates all operations."""

    args = parse_arguments()

    config = get_config_settings(args)

    docker_compose_content = generate_docker_compose_content(config)

    deploy_configuration(docker_compose_content)


ENV_FILE_TEMPLATE = """
DOMAIN={domain}
PUBLIC_KESSLER_API_URL={https_public_api_url}
VERSION_HASH={version_hash}
""".strip()


if __name__ == "__main__":
    main()
