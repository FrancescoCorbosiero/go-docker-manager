#!/bin/bash
TIMESTAMP=$(date +"%F")
BACKUP_DIR="/backups/mysql"
mkdir -p "$BACKUP_DIR"
docker exec website-alpacode-mariadb-1 \
  mysqldump -ualpacode -p'la-tua-password' alpacode_db > "$BACKUP_DIR/alpacode_db_$TIMESTAMP.sql"
