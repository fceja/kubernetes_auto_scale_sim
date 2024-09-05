#!/bin/bash

GREEN='\033[32m'
RESET='\033[0m'
filePath="${GREEN}${INIT_FILE_PATH}${RESET} -"

# Allow Zookeeper to start
echo -e "$filePath Allowing for Zookeeper to load. Sleep 3 secs."
sleep 3

# Connect to Zookeeper
maxAttempts=10
attempt=0
waitTime=5

echo -e "$filePath Checking for Zookeeper availability at '$ZOOKEEPER_HOSTNAME:$ZOOKEEPER_SERVER_PORT'."
while ! nc -zv $ZOOKEEPER_HOSTNAME $ZOOKEEPER_SERVER_PORT 2>/dev/null; do
    attempt=$((attempt + 1))
    if [ "$attempt" -ge "$maxAttempts" ]; then
        echo -e "$filePath Failed to connect to Zookeeper at '$ZOOKEEPER_HOSTNAME:$ZOOKEEPER_SERVER_PORT'. Exiting."
        exit 1
    fi
    echo -e "$filePath Attempting to connect..."
    sleep $waitTime
done
echo -e "$filePath Connected to Zookeeper."

# Start Kafka server
echo -e "$filePath Starting Kafka server."

${KAFKA_HOME}/bin/kafka-server-start.sh ${KAFKA_HOME}/config/server.kubernetes.properties &
KAFKA_PID=$!

# Wait for Kafka server to be ready
maxAttempts=10
attempt=0

echo -e "$filePath Allowing for Kafka server to load. Sleep 3 secs."
sleep 3

# Kubernetes automatically appends to 'KAFKA_SERVER_PORT'
# we only care about actual port number, parse only what we need
# e.g., "tcp://10.96.17.57:9092" -> "9092"
PORT=$(echo "${KAFKA_SERVER_PORT}" | sed 's|.*:||')

# Check if Kafka port ready
while ! nc -zv $KAFKA_HOSTNAME $PORT 2>/dev/null; do
    attempt=$((attempt + 1))
    if [ "$attempt" -gt "$maxAttempts" ]; then
        echo -e "$filePath Kafka server failed to load listening port. Exiting."
        exit 1
    fi
    echo -e "$filePath Pinging- '${KAFKA_HOSTNAME}:${KAFKA_SERVER_PORT}'"
    sleep $waitTime
done
echo -e "$filePath Kafka server is ready - '$KAFKA_HOSTNAME:$KAFKA_SERVER_PORT'."

# Remove netcat-openbsd since no longer needed
echo -e "$filePath Cleaning up."
apt-get remove -y netcat-openbsd >/dev/null 2>&1 &&
    apt-get autoremove -y >/dev/null 2>&1 &&
    apt-get clean >/dev/null 2>&1 &&
    rm -rf /var/lib/apt/lists/* >/dev/null 2>&1
echo -e "$filePath Done."

# Create topics on Kafka server start
echo -e "$filePath Executing './main'"
./main

# Keep Kafka server process running
wait $KAFKA_PID
