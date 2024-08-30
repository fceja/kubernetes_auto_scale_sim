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
echo -e "$filePath Starting Kafdrop server."

./kafdrop.sh &
KAFDROP_PID=$!

wait $KAFDROP_PID
