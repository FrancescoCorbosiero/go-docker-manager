#!/bin/bash

set -e

# Environment variables
TEMPLATES_DIR="templates"

# Ensure required networks exist
for net in traefik-network wordpress-network; do
  if ! docker network inspect "$net" >/dev/null 2>&1; then
    echo "🔧 Creating missing network: $net"
    docker network create "$net"
  fi
done

# Check if Traefik container is up and healthy

TRAEFIK_CONTAINER_NAME=traefik-traefik-1

if ! docker ps --format '{{.Names}}' | grep -q "^$TRAEFIK_CONTAINER_NAME$"; then
  echo "❌ Traefik container \"$TRAEFIK_CONTAINER_NAME\" is not running."
  echo "Showing last 20 log lines (if available):"
  docker logs --tail 20 "$TRAEFIK_CONTAINER_NAME" 2>/dev/null || echo "⚠️ No logs found. Container may not exist."
  exit 1
fi

echo "⏳ Checking health of Traefik container..."

STATUS=$(docker inspect --format='{{.State.Health.Status}}' "$TRAEFIK_CONTAINER_NAME" 2>/dev/null || echo "none")

if [ "$STATUS" == "healthy" ]; then
  echo "✅ Traefik is healthy."
elif [ "$STATUS" == "none" ]; then
  echo "⚠️ Traefik does not have a health check defined. Proceeding anyway."
else
  echo "❌ Traefik container is not healthy (status: $STATUS)."
  echo "🪵 Showing last 20 log lines:"
  docker logs --tail 20 "$TRAEFIK_CONTAINER_NAME"
  exit 1
fi


# Accept name from Makefile or ask interactively
if [ -z "$1" ]; then
  read -p "Site name (no spaces, e.g. website): " NAME
else
  NAME="$1"
fi

read -p "Domain (e.g. ${NAME}.com): " DOMAIN
read -p "Admin email: " ADMIN_EMAIL
read -p "Admin username: " ADMIN_USER
read -p "Admin password: " ADMIN_PASS
read -p "SMTP host: " STMP_HOST

SITE_DIR="~/containers/compose/$NAME"
mkdir -p "$SITE_DIR"

cat > "$SITE_DIR/.env" <<EOF
PROJECT_NAME=$NAME
WORDPRESS_HOSTNAME=$DOMAIN

WORDPRESS_DB_NAME=$ADMIN_USER
WORDPRESS_DB_USER=$ADMIN_USER
WORDPRESS_DB_PASSWORD=$ADMIN_PASS
WORDPRESS_DB_ADMIN_PASSWORD=$ADMIN_PASS

WORDPRESS_ADMIN_USERNAME=$ADMIN_USER
WORDPRESS_ADMIN_PASSWORD=$ADMIN_PASS
WORDPRESS_ADMIN_EMAIL=$ADMIN_EMAIL

WORDPRESS_BLOG_NAME=$NAME
WORDPRESS_ADMIN_NAME=$ADMIN_USER
WORDPRESS_ADMIN_LASTNAME=admin
WORDPRESS_TABLE_PREFIX=wp_

WORDPRESS_SMTP_ADDRESS=$STMP_HOST
WORDPRESS_SMTP_PORT=587
WORDPRESS_SMTP_USER_NAME=$ADMIN_EMAIL
WORDPRESS_SMTP_PASSWORD=$ADMIN_PASS

WORDPRESS_IMAGE_TAG=bitnami/wordpress:latest
WORDPRESS_MARIADB_IMAGE_TAG=bitnami/mariadb:latest
EOF

# Process docker-compose from template
envsubst < $TEMPLATES_DIR/wordpress.yml > "$SITE_DIR/docker-compose.yml"
#cp "$TEMPLATES_DIR/wordpress.yml" "$SITE_DIR/docker-compose.yml"

echo "✅ Project created in $SITE_DIR"
echo "➡️  To start it: cd $SITE_DIR && docker compose --env-file .env up -d"
