#!/bin/bash

# Output file
output_file=".local.env"

# Clear the output file if it exists
> "$output_file"

# Iterate over each environment variable
while IFS='=' read -r name value; do
    echo "$name=$value" >> "$output_file"
done < <(env)

echo "Environment variables have been written to $output_file"
