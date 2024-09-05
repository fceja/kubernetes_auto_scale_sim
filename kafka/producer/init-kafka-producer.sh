#!/bin/bash

GREEN='\033[32m'
RESET='\033[0m'
filePath="${GREEN}${INIT_FILE_PATH}${RESET} -"

# Kubernetes automatically appends to 'KAFKA_SERVER_PORT'
# we only care about actual port number, parse only what we need
# e.g., "tcp://10.96.17.57:9092" -> "9092"
PORT=$(echo "${KAFKA_SERVER_PORT}" | sed 's|.*:||')

# Wait for Kafka server to be ready
echo -e "$filePath Waiting for Kafka server availability at '$KAFKA_HOSTNAME:$PORT'. Sleep 10 secs."
sleep 10

maxAttempts=10
attempt=0
waitTime=2

while ! nc -zv $KAFKA_HOSTNAME $PORT 2>/dev/null; do
    attempt=$((attempt + 1))
    if [ "$attempt" -gt "$maxAttempts" ]; then
        echo -e "Failed to connect to Kafka server at '$KAFKA_HOSTNAME:$PORT'. Exiting."
        exit 1
    fi
    echo -e "$filePath Attempting to connect..."
    sleep $waitTime
done
echo -e "$filePath Kafka connection successful at '$KAFKA_HOSTNAME:$PORT'."

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
