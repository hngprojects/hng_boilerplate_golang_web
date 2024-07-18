#!/bin/bash

# Check if the script is run as root
if [ "$(id -u)" -ne 0 ]; then
    echo "This script must be run as root. Use sudo."
    exit 1
fi

chmod +x nginx.sh postgres_go.sh postgres_setup.sh

./nginx.sh
./postgres_go.sh
./postgres_setup.sh
