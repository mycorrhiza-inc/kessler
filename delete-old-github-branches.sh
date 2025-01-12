#!/bin/bash

# Start the SSH agent and load the key
eval "$(ssh-agent -s)"
# change this to your key
ssh-add ~/.ssh/id_ed25519

# Get list of all remote branches and filter out main and release branches
git branch -r | grep -v "main\|release" > ~/branches_to_check.txt

# Create empty file for branches to delete
> ~/branches_to_delete.txt

# Check each branch for merge status and differences
while IFS= read -r ref_branch; do
    ref_branch=$(echo "$ref_branch" | tr -d '[:space:]')
    branch=$(echo "$ref_branch" | grep -Eo '[^/]+$')
    
    # Check if branch is merged into main
    if git branch -r --merged origin/main | grep -q "$ref_branch"; then
        # Check if branch has any unique commits compared to main
        if [ -z "$(git log origin/main.."$ref_branch" --oneline)" ]; then
            echo "$ref_branch" >> ~/branches_to_delete.txt
        fi
    fi
done < ~/branches_to_check.txt

# Delete the branches
while IFS= read -r ref_branch; do
    ref="$(echo "$ref_branch" | grep -Eo '^[^/]+')"
    branch="$(echo "$ref_branch" | grep -Eo '[^/]+$')"
    echo "Deleting branch: $branch from $ref"
    git push "$ref" --delete "$branch"
done < ~/branches_to_delete.txt

# Kill the SSH agent
ssh-agent -k

# Cleanup temporary files
rm ~/branches_to_check.txt ~/branches_to_delete.txt
