#!/bin/bash
set -e

# Global variables
remote_target=""
prod_commit=""
nightly_commit=""

function process_branch() {
    local branch=$1
    local commit_hash=$2  # Optional commit hash
    local api_url=""
    
    # Set API URL based on branch
    # if [ "$branch" = "main" ]; then
    #     api_url="https://api.kessler.xyz"
    #     tag="latest"
    # else
    #     api_url="https://nightly-api.kessler.xyz"
    #     tag="nightly"
    # fi
    api_url="https://api.kessler.xyz"
    local api_version_hash_url="${api_url}/v2/version_hash"

    echo "Processing branch: $branch"
    sudo mkdir -p /mycorrhiza
    sudo chmod 777 -R /mycorrhiza
    cd /mycorrhiza
    
    if [ ! -d "/mycorrhiza/kessler" ]; then
        git clone https://github.com/mycorrhiza-inc/kessler
        git config --global --add safe.directory /mycorrhiza/kessler
    fi
    
    cd kessler

       # Checkout specific commit or update branch
        if [ -n "$commit_hash" ]; then
            git clean -fd
            git fetch
            git reset --hard HEAD
            git clean -fd
            git checkout "$commit_hash"
            echo "Checked out specific commit: $commit_hash"
        else
            git clean -fd
            git fetch
            git reset --hard HEAD
            git clean -fd
            git switch "$branch"
            git reset --hard origin/"$branch"
            echo "Updated branch $branch to latest"
        fi

    local current_hash=$(git rev-parse HEAD)
    echo "Current commit hash: $current_hash"
    
    local deployed_hash=$(curl -s "$api_version_hash_url" || echo "")
    echo "Deployed version hash: $deployed_hash"
    
    if [ -n "$commit_hash" ] || [ "$current_hash" != "$deployed_hash" ]; then
        echo "Rebuilding and deploying images..."
        
        # Build and push Docker images
        sudo docker build -t "fractalhuman1/kessler-frontend:${current_hash}" --platform linux/amd64 ./frontend/
        sudo docker build -t "fractalhuman1/kessler-backend-go:${current_hash}" --platform linux/amd64 ./backend-go/

        sudo docker push "fractalhuman1/kessler-frontend:${current_hash}"
        sudo docker push "fractalhuman1/kessler-backend-go:${current_hash}"

        # Update docker-compose.yml on the server
        ssh root@kessler.xyz "sed -i \"s|image: fractalhuman1/kessler-frontend:.*|image: fractalhuman1/kessler-frontend:${current_hash}|\" /mycorrhiza/kessler/docker-compose.deploy.yaml && \
        sed -i \"s|image: fractalhuman1/kessler-backend-go:.*|image: fractalhuman1/kessler-backend-go:${current_hash}|\" /mycorrhiza/kessler/docker-compose.deploy.yaml && \
        sed -i \"s|CONFIG_HASH:.*|CONFIG_HASH: ${current_hash}|\" /mycorrhiza/kessler/docker-compose.deploy.yaml && \
        cd /mycorrhiza/kessler && git clean -fd && git pull && docker compose -f docker-compose.deploy.yaml down && docker compose -f docker-compose.deploy.yaml up -d"
    else
        echo "No changes detected, deployemnt already on provided hash, skipping deployment: ${current_hash}"
    fi
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        --prod-commit)
            commit="$2"
            shift 2
            ;;
        --prod-commit=*)
            commit="${1#*=}"
            shift
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

if [[ -n "$prod_commit" || -n "$nightly_commit" ]]; then
    if [[ -n "$prod_commit" ]]; then
        process_branch "main" "$prod_commit"
    fi
else
    process_branch "main"
    # process_branch "release"
fi
#!/bin/bash
set -e


