#!/bin/bash

# Variables
DOMAIN_OR_IP="91.229.239.238"
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

# Obtain an SSL certificate from Let's Encrypt using Certbot
sudo certbot --nginx -d $DOMAIN_OR_IP --non-interactive --agree-tos --email $EMAIL

# # Set up automatic renewal of the certificate
# echo "0 0,12 * * * root certbot renew --quiet" | sudo tee -a /etc/crontab > /dev/null

# Display the final server URL
echo "SSL certificate configured for: https://$DOMAIN_OR_IP"
