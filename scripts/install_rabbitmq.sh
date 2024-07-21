#!/bin/bash

# Update package lists
sudo apt-get update

# Install RabbitMQ server
sudo apt-get install -y rabbitmq-server

# Start & Enable RabbitMQ
sudo systemctl start rabbitmq-server
sudo systemctl enable rabbitmq-server