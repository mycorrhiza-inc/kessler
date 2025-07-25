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
        git clone https://github.com/mycorrhiza-inc/kessler --recurse-submodules
        git config --global --add safe.directory /mycorrhiza/kessler
    fi
    
    cd kessler

       # Checkout specific commit or update branch
        if [ -n "$commit_hash" ]; then
            git clean -fd
            git fetch
            git reset --hard HEAD
            git clean -fd
            git checkout "$commit_hash" --recurse-submodules
            echo "Checked out specific commit: $commit_hash"
        else
            git clean -fd
            git fetch
            git reset --hard HEAD
            git clean -fd
            git checkout "$branch" --recurse-submodules
            git reset --hard origin/"$branch"
            echo "Updated branch $branch to latest"
        fi

    local current_hash=$(git rev-parse HEAD)
    echo "Current commit hash: $current_hash"
    
    local deployed_hash=$(curl -s "$api_version_hash_url" || echo "")
    echo "Deployed version hash: $deployed_hash"
    
    if [ -n "$commit_hash" ] || [ "$current_hash" != "$deployed_hash" ]; then
        echo "Rebuilding and deploying images..."
        
        echo "Building Frontend Image"
        # Build and push Docker images
        sudo docker build -t "fractalhuman1/kessler-frontend:${current_hash}" --platform linux/amd64 --file ./frontend/prod.Dockerfile ./frontend/
        echo "Building Backend Server Image"
        sudo docker build -t "fractalhuman1/kessler-backend-server:${current_hash}" --platform linux/amd64 --file ./backend/prod.server.Dockerfile ./backend
        echo "Building Backend Ingest Image"
        sudo docker build -t "fractalhuman1/kessler-backend-ingest:${current_hash}" --platform linux/amd64 --file ./backend/prod.ingest.Dockerfile ./backend

        echo "Building Fugu Database Image"
        sudo docker build -t "fractalhuman1/kessler-fugudb:${current_hash}" --platform linux/amd64 --file ./fugu/Dockerfile ./fugu

        sudo docker push "fractalhuman1/kessler-frontend:${current_hash}"
        sudo docker push "fractalhuman1/kessler-backend-server:${current_hash}"
        sudo docker push "fractalhuman1/kessler-backend-ingest:${current_hash}"
        sudo docker push "fractalhuman1/kessler-fugudb:${current_hash}"

        # Update docker-compose.yml on the server
        # Set deployment variables based on environment
        local deploy_flag=""
        local deploy_host=""
        if [ "$is_prod" = true ]; then
          deploy_flag="production"
          deploy_host="kessler.xyz"
        else
          deploy_flag="nightly"
          deploy_host="nightly.kessler.xyz"
        fi


        ssh "root@${deploy_host}" "cd /mycorrhiza/kessler && git reset --hard HEAD && git clean -fd && git switch main && git pull && python3 execute_production_deploy.py --${deploy_flag} --version ${current_hash}"
    else
        echo "No changes detected, deployemnt already on provided hash, skipping deployment: ${current_hash}"
    fi
}

# Parse command line arguments
is_prod=false
while [[ $# -gt 0 ]]; do
    case "$1" in
        --commit)
            commit="$2"
            shift 2
            ;;
        --commit=*)
            commit="${1#*=}"
            shift
            ;;
        --prod)
            is_prod=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

if [[  -n "$commit" ]]; then
    process_branch "main" "$commit"
else
    process_branch "main"
    # process_branch "release"
fi
#!/bin/bash
set -e


