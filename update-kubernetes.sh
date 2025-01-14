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
    local tag=""
    
    # Set tag based on branch
    if [ "$branch" = "release" ]; then
        tag="latest"
    else
        tag="nightly"
    fi
    
    echo "Processing branch: $branch with tag: $tag"
    
    # Switch to branch
    git switch "$branch"
    git pull

    echo "Switched to branch: $branch and checked for new commits."
    
    # Get latest commit timestamp
    commit_timestamp=$(git log -1 --format=%ct)
    
    # Get Docker image timestamp
    image_timestamp=$(docker inspect -f '{{ .Created }}' "fractalhuman1/kessler-frontend:${tag}" || echo "1970-01-01T00:00:00Z")
    image_unix_timestamp=$(date -d "$image_timestamp" +%s)
    
    if [ $commit_timestamp -gt $image_unix_timestamp ]; then
        echo "New commits found, rebuilding images..."
        
        # Build and push Docker images
        sudo docker build -t "fractalhuman1/kessler-frontend:${tag}" --platform linux/amd64 ./frontend/
        sudo docker build -t "fractalhuman1/kessler-backend-go:${tag}" --platform linux/amd64 ./backend-go/
        
        sudo docker push "fractalhuman1/kessler-frontend:${tag}"
        sudo docker push "fractalhuman1/kessler-backend-go:${tag}"
        
        # Deploy to appropriate environment
        if [ "$tag" = "latest" ]; then
            ssh root@nightly.kessler.xyz "cd /mycorrhiza/infra && helm upgrade kessler-prod ./helm -f helm/values-prod.yaml"
        else
            ssh root@nightly.kessler.xyz "cd /mycorrhiza/infra && helm upgrade kessler-nightly ./helm -f helm/values-nightly.yaml"
        fi
    else
        echo "No new commits found for $branch, skipping..."
    fi
}

# Process both branches
process_branch "main"
process_branch "release"
