#!/bin/bash

# Get site directory (you run this from inside parent/)
if [ -z "$1" ]; then
  echo "Usage: $0 <site-dir>"
  exit 1
fi

SITE_DIR=$1
SITE_NAME=$(basename "$SITE_DIR")

COMPOSE_FILE="$SITE_DIR/docker-compose.yml"
ENV_FILE="$SITE_DIR/.env"

if [ ! -f "$COMPOSE_FILE" ]; then
  echo "❌ Compose file not found: $COMPOSE_FILE"
  exit 1
fi

if [ ! -f "$ENV_FILE" ]; then
  echo "❌ .env file not found: $ENV_FILE"
  exit 1
fi

echo "🚀 Launching site: $SITE_NAME"
docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d
