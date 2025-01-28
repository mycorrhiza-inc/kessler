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
    local tag=""
    
    # Set API URL based on branch
    if [ "$branch" = "release" ]; then
        api_url="https://api.kessler.xyz"
        tag="latest"
    else
        api_url="https://nightly-api.kessler.xyz"
        tag="nightly"
    fi
    local api_version_hash_url="${api_url}/v2/version_hash"

    echo "Processing branch: $branch"
    sudo mkdir -p /mycorrhiza
    sudo chmod 777 -R /mycorrhiza
    cd /mycorrhiza
    
    if [ ! -d "/mycorrhiza/kessler" ]; then
        git clone https://github.com/mycorrhiza-inc/kessler
    fi
    
    cd kessler

    # Checkout specific commit or update branch
    if [ -n "$commit_hash" ]; then
        git checkout "$commit_hash"
        echo "Checked out specific commit: $commit_hash"
    else
        git switch "$branch"
        git pull
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
        
        # Determine deployment environment
        if [ "$tag" = "latest" ]; then
            ssh root@nightly.kessler.xyz "export KUBECONFIG=/etc/rancher/k3s/k3s.yaml && cd /mycorrhiza/infra && helm upgrade kessler-prod ./helm -f helm/values-prod.yaml --set versionHash=${current_hash}"
        else
            ssh root@nightly.kessler.xyz "export KUBECONFIG=/etc/rancher/k3s/k3s.yaml && cd /mycorrhiza/infra && helm upgrade kessler-nightly ./helm -f helm/values-nightly.yaml --set versionHash=${current_hash}"
        fi
    else
        echo "No changes detected for $branch, skipping deployment."
    fi
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case "$1" in
        --remote)
            remote_target="$2"
            shift 2
            ;;
        --remote=*)
            remote_target="${1#*=}"
            shift
            ;;
        --prod-commit)
            prod_commit="$2"
            shift 2
            ;;
        --prod-commit=*)
            prod_commit="${1#*=}"
            shift
            ;;
        --nightly-commit)
            nightly_commit="$2"
            shift 2
            ;;
        --nightly-commit=*)
            nightly_commit="${1#*=}"
            shift
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Handle remote execution
if [ -n "$remote_target" ]; then
    user=$(whoami)
    host="$remote_target"
    # Split user@host if provided
    if [[ "$remote_target" == *"@"* ]]; then
        user=$(echo "$remote_target" | cut -d@ -f1)
        host=$(echo "$remote_target" | cut -d@ -f2)
    fi

    # Pass all remaining arguments to remote execution
    ssh "$user@$host" /bin/bash -s -- --no-remote "$prod_commit" "$nightly_commit" <<'EOF'
    # Remote script execution starts here
    set -e
    prod_commit=""
    nightly_commit=""
    
    # Re-parse arguments on remote side
    while [[ $# -gt 0 ]]; do
        case "$1" in
            --prod-commit)
                prod_commit="$2"
                shift 2
                ;;
            --nightly-commit)
                nightly_commit="$2"
                shift 2
                ;;
            *)
                shift
                ;;
        esac
    done

    if [[ -n "$prod_commit" || -n "$nightly_commit" ]]; then
        if [[ -n "$prod_commit" ]]; then
            process_branch "release" "$prod_commit"
        fi
        if [[ -n "$nightly_commit" ]]; then
            process_branch "main" "$nightly_commit"
        fi
    else
        process_branch "main"
        process_branch "release"
    fi
EOF
    exit 0
fi

# Local execution
if [[ -n "$prod_commit" || -n "$nightly_commit" ]]; then
    if [[ -n "$prod_commit" ]]; then
        process_branch "release" "$prod_commit"
    fi
    if [[ -n "$nightly_commit" ]]; then
        process_branch "main" "$nightly_commit"
    fi
else
    process_branch "main"
    process_branch "release"
fi
#!/bin/bash
set -e


