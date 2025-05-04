#!/bin/bash

# Usage:
# ./reset-admin.sh admin newpassword wordpress ireandlunaslife_mariadb-1

USER="${1:-admin}"
PASS="${2:-changeme123}"
DB="${3:-wordpress}"
CONTAINER="${4:-ireandlunaslife_mariadb-1}"

# SQL template with variable placeholders
SQL_TEMPLATE=$(cat <<EOF
SET @username := '${USER}';
SET @new_password := '${PASS}';
UPDATE wp_users SET user_pass = MD5(@new_password) WHERE user_login = @username;
EOF
)

# Run the SQL command inside the container
echo "🔐 Resetting WordPress admin password..."
docker exec -i "$CONTAINER" mysql -u root -p"$MYSQL_ROOT_PASSWORD" "$DB" -e "$SQL_TEMPLATE"

echo "✅ Done! Admin password for '$USER' has been updated."
