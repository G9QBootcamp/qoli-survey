#!/bin/bash
set -e

# Create an empty log file if it doesn't exist
touch /var/log/cron.log

# Write the cron job to a temporary file
echo "$CRON_SCHEDULE /scripts/compress_and_remove.sh >> /var/log/cron.log 2>&1" > /tmp/cronjob

# Load the cron job using crontab
crontab /tmp/cronjob

# Start the cron daemon in the background
cron

# Tail the log file to keep the container running
tail -f /var/log/cron.log
