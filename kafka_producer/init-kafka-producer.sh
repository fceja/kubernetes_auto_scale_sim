#!/bin/bash

GREEN='\033[32m'
RESET='\033[0m'
filePath="${GREEN}${INIT_FILE_PATH}${RESET} -"

# Wait for Kafka server to be ready
echo -e "$filePath Waiting for Kafka server availability at '${KAFKA_HOSTNAME}:${KAFKA_SERVER_PORT}'"
maxAttempts=10
attempt=0
waitTime=2

while ! nc -zv "${KAFKA_HOSTNAME}" "${KAFKA_SERVER_PORT}" 2>/dev/null; do
    attempt=$((attempt + 1))
    if [ "$attempt" -gt "$maxAttempts" ]; then
        echo "Failed to connect to Kafka server at '${KAFKA_HOSTNAME}:${KAFKA_SERVER_PORT}'. Exiting."
        exit 1
    fi
    echo -e "$filePath Attempting to connect..."
    sleep $waitTime
done
echo -e "$filePath Kafka connection successful."

# Remove netcat-openbsd since no longer needed
echo -e "$filePath Cleaning up."
apt-get remove -y netcat-openbsd >/dev/null 2>&1 &&
    apt-get autoremove -y >/dev/null 2>&1 &&
    apt-get clean >/dev/null 2>&1 &&
    rm -rf /var/lib/apt/lists/* >/dev/null 2>&1
echo -e "$filePath Done."

# Run producer
echo -e "$filePath Starting Kafka producer."
./kafka-producer &
KAFKA_PRODUCER=$!

wait $KAFKA_PRODUCER
