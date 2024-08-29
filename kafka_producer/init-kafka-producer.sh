#!/bin/bash

GREEN='\033[32m'
RESET='\033[0m'

filePath="${GREEN}${INIT_FILE_PATH}${RESET} -"
waitTime=2

# Wait for Kafka server to be ready
echo -e "$filePath Waiting for Kafka server to be available."
# while ! nc -zv kafka 9092 2>/dev/null; do
while ! nc -zv "${KAFKA_HOSTNAME}" "${KAFKA_SERVER_PORT}" 2>/dev/null; do
    echo -e "$filePath Checking again in $waitTime seconds."
    sleep $waitTime
done
echo -e "$filePath Kafka server is ready."

# Remove netcat-openbsd since no longer needed
echo -e "$filePath Cleaning up."
apt-get remove -y netcat-openbsd >/dev/null 2>&1 &&
    apt-get autoremove -y >/dev/null 2>&1 &&
    apt-get clean >/dev/null 2>&1 &&
    rm -rf /var/lib/apt/lists/* >/dev/null 2>&1

echo -e "$filePath Done."

# Run producer
echo "Starting Kafka producer."

./kafka-producer &

wait # Do we want want producer to run, or terminate after completing task?
