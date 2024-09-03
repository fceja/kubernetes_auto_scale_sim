#!/bin/bash

GREEN='\033[32m'
RESET='\033[0m'

filePath="${GREEN}${INIT_FILE_PATH}${RESET} -"

# Start Zookeper
echo -e "${filePath} Starting Zookeeper at '${ZOOKEEPER_HOSTNAME}:${ZOOKEEPER_SERVER_PORT}'."
bin/zkServer.sh start-foreground &
ZOOKEEPER_PID=$! # save PID

# Wait for Zookeeper to load
maxAttempts=10
attempt=0
waitTime=2

echo -e "$filePath Loading..."
sleep 3

while ! nc -zv "${ZOOKEEPER_HOSTNAME}" "${ZOOKEEPER_SERVER_PORT}" 2>/dev/null; do
    attempt=$((attempt + 1))
    if [ "$attempt" -ge "$maxAttempts" ]; then
        echo -e "$filePath Zookeeper failed to load. Exiting."
        exit 1
    fi
    echo -e "$filePath Pinging - '${ZOOKEEPER_HOSTNAME}:${ZOOKEEPER_SERVER_PORT}'"
    sleep $waitTime
done
echo -e "$filePath Zookeeper is ready at '${ZOOKEEPER_HOSTNAME}:${ZOOKEEPER_SERVER_PORT}'."

# Remove netcat-openbsd since no longer needed
echo -e "$filePath Cleaning up."
apt-get remove -y netcat-openbsd >/dev/null 2>&1 &&
    apt-get autoremove -y >/dev/null 2>&1 &&
    apt-get clean >/dev/null 2>&1 &&
    rm -rf /var/lib/apt/lists/* >/dev/null 2>&1

echo -e "$filePath Done."

# Keep Zookeeper process running
wait $ZOOKEEPER_PID
