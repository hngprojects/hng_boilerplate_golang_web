#!/bin/bash

# Check if the script is run as root
if [ "$(id -u)" -ne 0 ]; then
    echo "This script must be run as root. Use sudo."
    exit 1
fi

# Define the script directory
SCRIPT_DIR="$HOME/scripts"

# Ensure script directory exists
if [ ! -d "$SCRIPT_DIR" ]; then
    echo "Script directory $SCRIPT_DIR does not exist."
    exit 1
fi

# Array of scripts to execute
SCRIPTS=(
    "install_and_setup_nginx.sh"
    "install_postgres_and_go.sh"
    "install_pm2.sh"
    "setup_postgres.sh"
)

# Make all scripts executable and execute them
for script in "${SCRIPTS[@]}"; then
    SCRIPT_PATH="$SCRIPT_DIR/$script"
    
    if [ -f "$SCRIPT_PATH" ]; then
        chmod +x "$SCRIPT_PATH"
        echo "Executing $SCRIPT_PATH..."
        "$SCRIPT_PATH"
    else
        echo "Script $SCRIPT_PATH does not exist."
        exit 1
    fi
done