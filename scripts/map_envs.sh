#!/bin/bash

if [ "$#" -lt 2 ]; then
    echo "Description: creates new environment variables from existing ones"
    echo "Usage: $0 .env_file NEW_ENV1=EXISTING_ENV1 [NEW_ENV2=EXISTING_ENV2 ...]"
    exit 1
fi

env_file="$1"
shift

if [ ! -f "$env_file" ]; then
    echo "Error: .env file '$env_file' not found."
    exit 1
fi

export $(grep -v '^#' "$env_file" | xargs)

for arg in "$@"; do
    IFS='=' read -r new_env existing_env <<< "$arg"

    if [ -z "${!existing_env}" ]; then
        echo "Warning: Existing environment variable '$existing_env' is not set."
        continue
    fi

    # Get the value of the existing environment variable
    value="${!existing_env}"

    # Export the new environment variable with the value of the existing one
    export "$new_env=$value"
done

# Write the new environment variables to the .env file
{
    echo
    for arg in "$@"; do
        IFS='=' read -r new_env _ <<< "$arg"
        echo "$new_env=${!new_env}"
    done
} >> "$env_file"
