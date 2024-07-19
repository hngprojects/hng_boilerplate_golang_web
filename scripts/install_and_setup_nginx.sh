#!/bin/bash

# Define environment configurations
declare -A environments
environments=( 
    ["development"]="7000 deployment.api-golang.boilerplate.hng.tech"
    ["staging"]="8000 staging.api-golang.boilerplate.hng.tech"
    ["production"]="9000 api-golang.boilerplate.hng.tech"
)

# General setup
sudo apt update
sudo apt install -y nginx

# Create directories and configure Nginx for each environment
for env in "${!environments[@]}"; do
    # Split the value into port and domain using parameter expansion
    port="${environments[$env]%% *}"
    domain="${environments[$env]#* }"
    web_root="~/deployments/$env"
    config_file="/etc/nginx/conf.d/$env.conf"

    # Create web root directory
    sudo mkdir -p $web_root

    # Nginx configuration content
    content="server {
        listen 80;
        server_name $domain;
        location / {
            proxy_pass http://127.0.0.1:$port;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
        }
    }"

    # Write Nginx configuration file
    echo "$content" | sudo tee $config_file > /dev/null
    sudo chmod 664 $config_file
done

# Delete Default Nginx Webpage
sudo rm /etc/nginx/sites-available/default && sudo rm /etc/nginx/sites-enabled/default

# Validate the Configuration
sudo nginx -t
sudo systemctl restart nginx

# Output
echo "Nginx setup for the applications to reverse proxy requests to them"
