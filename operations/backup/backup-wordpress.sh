docker run --rm -v root_alpacode-wordpress-data:/data -v $(pwd)/backups:/backup alpine \
  tar czf /backup/wp-content-alpacode-$(date +%F).tar.gz -C /data .
