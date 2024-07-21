#!/bin/bash

# Variables
DOMAINS=(
  "api-golang.boilerplate.hng.tech"  # prod
  "deployment.api-golang.boilerplate.hng.tech"  # dev
  "staging.api-golang.boilerplate.hng.tech"  # staging
)
EMAIL="osinachi.chukwujama@gmail.com"  # Used for certbot notifications and recovery

# Update package lists
sudo apt update

# Install Snapd if not already installed
sudo apt install -y snapd

# Install the core snap & Ensure Snapd is up to date
sudo snap install core
sudo snap refresh core

# Remove any existing Certbot installations
sudo apt-get remove certbot

# Install Certbot using Snap
sudo snap install --classic certbot

# Create a symbolic link to make Certbot command globally available
sudo ln -s /snap/bin/certbot /usr/bin/certbot

# Obtain SSL certificates for each domain
for DOMAIN in "${DOMAINS[@]}"; do
  echo "Configuring SSL for domain: $DOMAIN"
  sudo certbot --nginx -d "$DOMAIN" --non-interactive --agree-tos --email "$EMAIL"
done

# Set up automatic renewal of the certificates
sudo crontab -l | { cat; echo "0 0,12 * * * root certbot renew --quiet"; } | sudo crontab -

# Display the configured domains
for DOMAIN in "${DOMAINS[@]}"; do
  echo "SSL certificate configured for: https://$DOMAIN"
done
