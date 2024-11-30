FROM bitnami/minideb:latest

# Install cron (includes crontab) and other necessary utilities
RUN apt-get update && apt-get install -y \
    cron gzip findutils && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Create necessary directories for scripts and logs
RUN mkdir -p /scripts /var/log


# Set the default command to bash
ENTRYPOINT ["bash"]
