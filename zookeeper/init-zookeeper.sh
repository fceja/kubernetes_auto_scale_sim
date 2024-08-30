#!/bin/bash

GREEN='\033[32m'
RESET='\033[0m'

filePath="${GREEN}${INIT_FILE_PATH}${RESET} -"

# Start Zookeper
echo -e "${filePath} Starting Zookeeper."
bin/zkServer.sh start-foreground &
ZOOKEEPER_PID=$!

# Wait for Zookeeper to load
maxAttempts=10
attempt=0
waitTime=2

echo -e "$filePath Waiting for Zookeeper to load."
while ! nc -zv "${ZOOKEEPER_HOSTNAME}" "${ZOOKEEPER_SERVER_PORT}" 2>/dev/null; do
    attempt=$((attempt + 1))
    if [ "$attempt" -ge "$maxAttempts" ]; then
        echo "$filePath Zookeeper failed to load. Exiting."
        exit 1
    fi
    echo "$filePath Loading..."
    sleep $waitTime
done
echo -e "$filePath Zookeeper is ready."

# Remove netcat-openbsd since no longer needed
echo -e "$filePath Cleaning up."
apt-get remove -y netcat-openbsd >/dev/null 2>&1 &&
    apt-get autoremove -y >/dev/null 2>&1 &&
    apt-get clean >/dev/null 2>&1 &&
    rm -rf /var/lib/apt/lists/* >/dev/null 2>&1

echo -e "$filePath Done."

# Keep Zookeeper process running
wait $ZOOKEEPER_PID
