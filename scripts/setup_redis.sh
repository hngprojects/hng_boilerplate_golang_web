#!/bin/bash

# Variables
CONFIG_DIR="/etc/redis"
LOG_DIR="/etc/server.d"
LOG_FILE="$LOG_DIR/golang.conf"
PORTS=("6130" "6131" "6132")
INSTANCES=("golang_redis_development" "golang_redis_staging" "golang_redis_production")
USERS=("golang_redis_development_user" "golang_redis_staging_user" "golang_redis_production_user")

# Create the config directory if it does not exist
mkdir -p $CONFIG_DIR

# Create the log directory if it does not exist
mkdir -p $LOG_DIR

# Clear the log file
> $LOG_FILE

# Function to generate a random password
generate_password() {
    tr -dc A-Za-z0-9 < /dev/urandom | head -c 12
}

# Loop over instances and create configuration files
for i in ${!INSTANCES[@]}; do
    INSTANCE=${INSTANCES[$i]}
    USER=${USERS[$i]}
    PORT=${PORTS[$i]}
    PASSWORD=$(generate_password)
    
    CONFIG_FILE="$CONFIG_DIR/$INSTANCE.conf"
    
    # Create and write the configuration file
    cat <<EOL > $CONFIG_FILE
# Redis Configuration File for $INSTANCE

# Port to listen on
port $PORT

# Require clients to issue AUTH <PASSWORD> before processing any other commands.
# requirepass $PASSWORD

# ACL Configuration
user $USER on +@all

# Log level
loglevel notice

# Log file
logfile ""

# Databases
databases 16

# Bind address
bind 127.0.0.1

# Timeout
timeout 0

# Save on disk
save 900 1
save 300 10
save 60 10000

# Append-only file
appendonly no

# RDB/AOF snapshot
appendfsync everysec
EOL

    echo "Redis configuration file '$CONFIG_FILE' has been created with user '$USER' and password '$PASSWORD' on port $PORT."
    
    # Log the username and password
    echo "$USER => $PASSWORD" >> $LOG_FILE
    
    # Start Redis server in the background
    redis-server $CONFIG_FILE &
done

echo "Redis instances started and passwords logged to '$LOG_FILE'."
