#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <site-dir>"
  exit 1
fi

SITE_NAME=$(basename "$1")

echo "🛑 Stopping site: $SITE_NAME"
docker compose -p $SITE_NAME down
