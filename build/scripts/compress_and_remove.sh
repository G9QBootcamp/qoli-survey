#!/bin/bash
set -e

TARGET_DIR="${TARGET_DIR:-/data}"
LOG_FILE="/var/log/file-maintenance.log"

# Compress and remove old files
find "$TARGET_DIR" -type f -mtime +7 -not -name "*.gz" -exec gzip {} \; -exec echo "Compressed {}" >> "$LOG_FILE" \;
find "$TARGET_DIR" -type f -name "*.gz" -mtime +30 -exec rm -f {} \; -exec echo "Removed {}" >> "$LOG_FILE" \;
