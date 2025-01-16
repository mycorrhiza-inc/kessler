# Do the following for both branches "main" and "release"
# 1. Set a variable called "TAG" to "nightly" or "latest" depending on the branch.
# 2. Git switch to the branch.
# 3. Check the date of the docker image "fractalhuman1/kessler-frontend:${TAG}" and check the date of the latest git commit. 
# IF there are new commits do the following:
# 4. Build the docker image and push it to "fractalhuman1/kessler-frontend:${TAG}" using the script below and stuff
# 5. ssh into nightly.kessler.xyz and depending on what TAG run the command:
# cd /mycorrhiza/infra
# # If release run :
# helm upgrade kessler-prod ./helm -f helm/values-prod.yaml
# # If nightly run 
# helm upgrade kessler-nightly ./helm -f helm/values-nightly.yaml



#!/bin/bash
set -e


function process_branch() {
    local branch=$1
    # local tag=""
    local api_url=""
    
    # Set tag based on branch
    if [ "$branch" = "release" ]; then
        api_url="https://api.kessler.xyz"
        # tag="latest"
    else
        api_url="https://nightly-api.kessler.xyz"
        # tag="nightly"
    fi
    local api_version_hash_url="${api_url}/v2/version_hash"

    echo "Processing branch: $branch with tag: $tag"
    # Create /mycorrhiza directory if it doesn't exist
    sudo mkdir -p /mycorrhiza
    sudo chmod 777 -R /mycorrhiza
    cd /mycorrhiza
    # Clone kessler repo if it doesn't exist
    if [ ! -d "/mycorrhiza/kessler" ]; then
        git clone https://github.com/mycorrhiza-inc/kessler
    fi
    cd /mycorrhiza/kessler
    
    # Switch to branch
    git switch "$branch"
    git pull

    echo "Switched to branch: $branch and checked for new commits."
    
    # Get current commit hash
    current_hash=$(git rev-parse HEAD)
    echo "Current commit hash: $current_hash"
    
    # Get deployed version hash
    deployed_hash=$(curl -s "$api_version_hash_url" || echo "")
    echo "Deployed version hash: $deployed_hash"
    
    if [ "$current_hash" != "$deployed_hash" ]; then
        echo "New commits found, rebuilding images..."
        
        # # Build and push Docker images
        # sudo docker build -t "fractalhuman1/kessler-frontend:${current_hash}" --platform linux/amd64 ./frontend/
        # sudo docker build -t "fractalhuman1/kessler-backend-go:${current_hash}" --platform linux/amd64 ./backend-go/
        #
        # sudo docker push "fractalhuman1/kessler-frontend:${current_hash}"
        # sudo docker push "fractalhuman1/kessler-backend-go:${current_hash}"
        
        # Deploy to appropriate environment
        #
           if [ "$tag" = "latest" ]; then
                ssh root@nightly.kessler.xyz "export KUBECONFIG=/etc/rancher/k3s/k3s.yaml && cd /mycorrhiza/infra && helm upgrade kessler-prod ./helm -f helm/values-prod.yaml --set versionHash=${current_hash}"
            else
                ssh root@nightly.kessler.xyz "export KUBECONFIG=/etc/rancher/k3s/k3s.yaml && cd /mycorrhiza/infra && helm upgrade kessler-nightly ./helm -f helm/values-nightly.yaml --set versionHash=${current_hash}"
            fi
    else
        echo "No new commits found for $branch, skipping..."
    fi
}

# Process both branches
process_branch "main"
# process_branch "release"
