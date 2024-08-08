#!/bin/bash

echo "DON'T RUN THIS SCRIPT MULTIPLE TIMES WITHOUT KILLING OLDER PROCESSES. VIEW PROCESSES USING ps aux | grep redis"

CONFIG_DIR="/etc/redis"
LOG_DIR="/var/log"
LOG_FILE="$LOG_DIR/golang.log"
REDIS_DIR="/var/lib/redis"
PORTS=("6130" "6131" "6132")
INSTANCES=("golang_redis_development" "golang_redis_staging" "golang_redis_production")
USERS=("golang_redis_development_user" "golang_redis_staging_user" "golang_redis_production_user")

# Create the config directory if it does not exist
sudo mkdir -p $CONFIG_DIR

# Create the log directory if it does not exist
sudo mkdir -p $LOG_DIR
touch $LOG_FILE

generate_password() {
    tr -dc A-Za-z0-9 < /dev/urandom | head -c 12
}


# Create Redis instances
for i in "${!INSTANCES[@]}"
do
    PORT="${PORTS[$i]}"
    INSTANCE="${INSTANCES[$i]}"
    USER="${USERS[$i]}"
    PASSWORD=$(generate_password)

    # Create Redis configuration file
    CONFIG_FILE="${CONFIG_DIR}/${INSTANCE}.conf"
    
    # Create directory for this instance
    INSTANCE_DIR="${REDIS_DIR}/${INSTANCE}"
    sudo mkdir -p $INSTANCE_DIR
    
    # Create log file for this instance
    REDIS_LOG_FILE="$LOG_DIR/$INSTANCE.log"
    sudo touch $REDIS_LOG_FILE
    sudo chown redis:redis $REDIS_LOG_FILE
    sudo chmod 664 $REDIS_LOG_FILE
    
    echo "# Redis Configuration File for $INSTANCE
bind 127.0.0.1
port ${PORT}

# Require clients to issue AUTH <PASSWORD> before processing any other commands. (not enforced)
# requirepass $PASSWORD

# ACL Configuration
user $USER on +@all

logfile ${LOG_DIR}/${INSTANCE}.log
dir ${INSTANCE_DIR}

" | sudo tee $CONFIG_FILE > /dev/null
    
    # Set appropriate permissions
    sudo chown -R redis:redis $INSTANCE_DIR
    
    # Create systemd service file
    SERVICE_FILE="/etc/systemd/system/${INSTANCE}.service"
    
    echo "[Unit]
Description=Redis instance for ${INSTANCE}
After=network.target

[Service]
Type=simple
User=redis
Group=redis
ExecStart=/usr/bin/redis-server ${CONFIG_FILE}
ExecStop=/usr/bin/redis-cli -p ${PORT} shutdown
Restart=always

[Install]
WantedBy=multi-user.target
" | sudo tee $SERVICE_FILE > /dev/null
    
    # Reload systemd to recognize the new service
    sudo systemctl daemon-reload
    
    # Enable and start the service
    sudo systemctl enable ${INSTANCE}
    sudo systemctl start ${INSTANCE}

    echo "Redis configuration file '$CONFIG_FILE' has been created with user '$USER' and password '$PASSWORD' on port $PORT."
    
    # Log the username and password
    echo "$USER => $PASSWORD" >> $LOG_FILE

    echo "Redis instance ${INSTANCE} started on port ${PORT}"
done

echo "All Redis instances have been created and started as systemd services."
echo "Redis instances started and passwords logged to '$LOG_FILE'."