#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <site-dir>"
  exit 1
fi

SITE_NAME=$(basename "$1")

echo "💣 Destroying site: $SITE_NAME (containers, networks, volumes)"
docker compose -p $SITE_NAME down -v
