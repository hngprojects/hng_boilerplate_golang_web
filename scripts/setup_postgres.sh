#! /bin/bash

# Variables for credentials
DB_PASSWORD="password"

# Define the databases and users
DATABASES=("development_db" "staging_db" "production_db")
USERS=("development_user" "staging_user" "production_user")

# Check if the script is running as root (necessary for changing PostgreSQL settings)
if [ "$(id -u)" -ne "0" ]; then
  echo "This script must be run as root."
  exit 1
fi

# Check if the DB_USER is postgres and alter the password of this user if necessary
if echo "${USERS[@]}" | grep -qw "postgres"; then
  sudo -i -u postgres psql <<EOF
ALTER USER postgres WITH PASSWORD '$DB_PASSWORD';
EOF
  echo "Password for 'postgres' user has been updated."
fi

# Create databases and users, and grant permissions
for i in ${!DATABASES[@]}; do
  DB_NAME=${DATABASES[$i]}
  DB_USER=${USERS[$i]}

  sudo -i -u postgres psql <<EOF
    -- Create a database named '$DB_NAME'
    CREATE DATABASE $DB_NAME;

    -- Create a user named '$DB_USER' with password '$DB_PASSWORD'
    CREATE USER $DB_USER WITH PASSWORD '$DB_PASSWORD';

    -- Grant all privileges on the database '$DB_NAME' to the user '$DB_USER'
    GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;

    -- Grant all privileges on the public schema of '$DB_NAME' to the user '$DB_USER'
    GRANT ALL ON SCHEMA public to $DB_USER;
EOF
  echo "Database '$DB_NAME' and user '$DB_USER' created with full access."
done

# Modify pg_hba.conf to allow password authentication
PG_HBA_FILE=$(sudo -i -u postgres psql -t -P format=unaligned -c 'SHOW hba_file')
sudo sed -i "s/^local\s\+all\s\+all\s\+peer/local   all             all                                     md5/" $PG_HBA_FILE
sudo sed -i "s/^host\s\+all\s\+all\s\+127.0.0.1\/32\s\+ident/host    all             all             127.0.0.1\/32            md5/" $PG_HBA_FILE
sudo sed -i "s/^host\s\+all\s\+all\s\+::1\/128\s\+ident/host    all             all             ::1\/128                 md5/" $PG_HBA_FILE
sudo bash -c "echo 'host    all             all             0.0.0.0/0               md5' >> $PG_HBA_FILE"
sudo bash -c "echo 'host    all             all             ::/0                    md5' >> $PG_HBA_FILE"

# Modify postgresql.conf to listen on all addresses
POSTGRESQL_CONF=$(sudo -i -u postgres psql -t -P format=unaligned -c 'SHOW config_file')
sudo sed -i "s/^#listen_addresses = 'localhost'/listen_addresses = '*'/" $POSTGRESQL_CONF

# Restart PostgreSQL to apply changes
sudo systemctl restart postgresql

echo "PostgreSQL setup is complete. Databases and users have been created and configured."
