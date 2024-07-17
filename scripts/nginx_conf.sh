#!/bin/bash

# Check if nginx is installed
if ! command -v nginx &> /dev/null
then
    echo "Nginx is not installed. Please install Nginx and try again."
    exit 1
fi

# Function to create nginx configuration files
create_nginx_conf() {
    local env=$1
    local port=$2
    local server_name=$3
    local config_file="/etc/nginx/conf.d/${env}.conf"

    cat <<EOL > "$config_file"
server {
    listen $port;
    server_name $server_name;

    location / {
        root /var/www/$env;
        index index.html index.htm;
    }

    error_page 404 /404.html;
    location = /404.html {
        internal;
    }

    error_page 500 502 503 504 /50x.html;
    location = /50x.html {
        internal;
    }
}
EOL
    echo "Nginx configuration for $env created at $config_file"
}

# Create directories for web roots if they don't exist
mkdir -p /var/www/dev
mkdir -p /var/www/prod

# Create nginx configuration files
create_nginx_conf "dev" 8000 "dev.domain.com"
create_nginx_conf "prod" 9000 "domain.com"

# Reload nginx to apply changes
nginx -t && systemctl reload nginx

echo "Nginx configuration setup completed successfully."
