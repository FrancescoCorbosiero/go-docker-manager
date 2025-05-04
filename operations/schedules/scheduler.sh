crontab -e
# esempio: backup ogni notte alle 3
0 3 * * * /root/scripts/backup-alpacode-db.sh >> /var/log/backup.log 2>&1
