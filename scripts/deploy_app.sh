#!/bin/bash

set -e

if [ $# -lt 2 ]; then
  echo "Usage: $0 <development|staging|production> <github.repository> [env_vars...]"
  exit 1
fi

ENV=$1
GITHUB_REPOSITORY=$2
shift 2

case $ENV in
  development)
    BRANCH="dev"
    APPROOT=~/deployments/development
    APP_NAME="development_app"
    ;;
  staging)
    BRANCH="staging"
    APPROOT=~/deployments/staging
    APP_NAME="staging_app"
    ;;
  production)
    BRANCH="main"
    APPROOT=~/deployments/production
    APP_NAME="production_app"
    ;;
  *)
    echo "Invalid environment: $ENV"
    exit 1
    ;;
esac

export APPROOT
export PATH=$PATH:~/.nvm/versions/node/v20.15.1/bin
export PATH=$PATH:/usr/local/go/bin

mkdir -p $APPROOT
cd $APPROOT

if [ -d "$APPROOT/.git" ]; then
  # Navigate to the repository directory and pull changes
  cd $APPROOT || { echo "Failed to navigate to web root directory"; exit 1; }
  git reset --hard HEAD || { echo "Failed to reset local changes"; exit 1; }
  git pull origin $BRANCH || { echo "Failed to pull latest changes"; exit 1; }
else
  git clone -b $BRANCH https://github.com/$GITHUB_REPOSITORY . || { echo "Failed to clone repository"; exit 1; }
fi

# Replace environment variables in app.env
cp app-sample.env app.env
for VAR in "$@"
do
  KEY=$(echo $VAR | cut -d '=' -f 1)
  VALUE=$(echo $VAR | cut -d '=' -f 2)
  sed -i "s|\${env.$KEY}|$VALUE|g" app.env
done

go build -o $APPROOT/$APP_NAME

# Start or restart the built binary with pm2
if pm2 list | grep -qw "$APP_NAME"; then
  echo "Process $APP_NAME is running. Restarting..."
  pm2 restart $APPROOT/$APP_NAME
else
  echo "Process $APP_NAME is not running. Starting..."
  pm2 start $APPROOT/$APP_NAME --name $APP_NAME
fi