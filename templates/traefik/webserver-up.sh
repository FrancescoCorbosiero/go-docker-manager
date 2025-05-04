#!/bin/bash
set -e

# Load config
CONFIG_FILE="/config.sh"
if [ ! -f "$CONFIG_FILE" ]; then
  echo "❌ Missing $CONFIG_FILE. Please ensure it exists."
  exit 1
fi
source "$CONFIG_FILE"

# Ensure .env exists
if [ ! -f "$TRAEFIK_ENV" ]; then
  echo "❌ Missing $TRAEFIK_ENV. Please create it manually based on documentation."
  exit 1
fi

# Ensure docker-compose.yml exists
if [ ! -f "$TRAEFIK_COMPOSE" ]; then
  echo "❌ Missing $TRAEFIK_COMPOSE. Please create it manually based on documentation."
  exit 1
fi

# Ensure required networks exist
for net in "${REQUIRED_NETWORKS[@]}"; do
  if ! docker network inspect "$net" >/dev/null 2>&1; then
    echo "🔧 Creating network: $net"
    docker network create "$net"
  fi
done

# Launch Traefik
echo "🚀 Running Traefik from $TRAEFIK_DIR..."
docker compose --env-file "$TRAEFIK_ENV" -f "$TRAEFIK_COMPOSE" -p "$TRAEFIK_CONTAINER_NAME" up -d

# Wait for container health
echo "⏳ Waiting for Traefik to be healthy..."
MAX_RETRIES=30
RETRY_COUNT=0

while true; do
  STATUS=$(docker inspect --format='{{.State.Health.Status}}' "$TRAEFIK_CONTAINER_NAME" 2>/dev/null || echo "none")
  if [ "$STATUS" == "healthy" ]; then
    echo "✅ Traefik is healthy and running."
    break
  elif [ "$STATUS" == "none" ]; then
    echo "⚠️  No health check defined. Proceeding anyway."
    break
  else
    if [ "$RETRY_COUNT" -ge "$MAX_RETRIES" ]; then
      echo "❌ Traefik failed to become healthy."
      echo "🪵 Showing last 20 log lines:"
      docker logs --tail 20 "$TRAEFIK_CONTAINER_NAME"
      exit 1
    fi
    echo "🔄 Waiting for health (attempt $((RETRY_COUNT + 1))/$MAX_RETRIES)..."
    sleep 2
    ((RETRY_COUNT++))
  fi
done
