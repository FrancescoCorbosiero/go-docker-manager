-- reset-credentials.sql

SET @username := '${WP_USER}';
SET @new_password := '${WP_PASS}';

UPDATE wp_users 
SET user_pass = MD5(@new_password) 
WHERE user_login = @username;
