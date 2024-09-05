#!/bin/bash

GREEN='\033[32m'
RESET='\033[0m'
filePath="${GREEN}${INIT_FILE_PATH}${RESET} -"

# Connect to Zookeeper
maxAttempts=10
attempt=0
waitTime=2

echo -e "$filePath Allow Zookeeper to load..."
sleep 7

echo -e "$filePath Checking for Zookeeper availability at '${ZOOKEEPER_HOSTNAME}:${ZOOKEEPER_SERVER_PORT}'."
while ! nc -zv "${ZOOKEEPER_HOSTNAME}" "${ZOOKEEPER_SERVER_PORT}" 2>/dev/null; do
    attempt=$((attempt + 1))
    if [ "$attempt" -ge "$maxAttempts" ]; then
        echo -e "$filePath Failed to connect to Zookeeper at '${ZOOKEEPER_HOSTNAME}:${ZOOKEEPER_SERVER_PORT}'. Exiting."
        exit 1
    fi
    echo -e "$filePath Attempting to connect..."
    sleep $waitTime
done
echo -e "$filePath Connected to Zookeeper."

# Start Kafka server
echo -e "$filePath Starting Kafka server."
${KAFKA_HOME}/bin/kafka-server-start.sh ${KAFKA_HOME}/config/server.properties.docker-compose &
KAFKA_PID=$!

# Wait for Kafka server to be ready
maxAttempts=10
attempt=0

echo -e "$filePath Allow Kafka server to load..."
sleep 1

# Check if inside port ready
while ! nc -zv "${KAFKA_INSIDE_HOSTNAME}" "${KAFKA_INSIDE_SERVER_PORT}" 2>/dev/null; do
    attempt=$((attempt + 1))
    if [ "$attempt" -gt "$maxAttempts" ]; then
        echo -e "$filePath Kafka server failed to start INSIDE listening port. Exiting."
        exit 1
    fi
    echo -e "$filePath Pinging INSIDE-'${KAFKA_INSIDE_HOSTNAME}:${KAFKA_INSIDE_SERVER_PORT}'"
    sleep $waitTime
done

# Check if outside port ready
while ! nc -zv "${KAFKA_OUTSIDE_HOSTNAME}" "${KAFKA_OUTSIDE_SERVER_PORT}" 2>/dev/null; do
    attempt=$((attempt + 1))
    if [ "$attempt" -gt "$maxAttempts" ]; then
        echo -e "$filePath Kafka server failed to load OUTSIDE listening port. Exiting."
        exit 1
    fi
    echo -e "$filePath Pinging OUTSIDE-'${KAFKA_OUTSIDE_HOSTNAME}:${KAFKA_OUTSIDE_SERVER_PORT}'"
    sleep $waitTime
done
echo -e "$filePath Kafka server is ready - INSIDE: '${KAFKA_INSIDE_HOSTNAME}:${KAFKA_INSIDE_SERVER_PORT}' OUTSIDE: '${KAFKA_OUTSIDE_HOSTNAME}:${KAFKA_OUTSIDE_SERVER_PORT}'."

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
