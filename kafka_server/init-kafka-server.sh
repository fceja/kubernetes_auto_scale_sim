#!/bin/bash

GREEN='\033[32m'
RESET='\033[0m'

filePath="${GREEN}${INIT_FILE_PATH}${RESET} -"
waitTime=2

# Connect to Zookeeper
echo -e "$filePath Waiting for Zookeeper to be available."
while ! nc -zv "${ZOOKEEPER_HOSTNAME}" "${ZOOKEEPER_SERVER_PORT}" 2>/dev/null; do
  echo "$filePath Checking again in $waitTime seconds."
  sleep $waitTime
done
echo -e "$filePath Connected to Zookeeper."

# Start Kafka server
echo -e "$filePath Starting Kafka server."
${KAFKA_HOME}/bin/kafka-server-start.sh ${KAFKA_HOME}/config/server.properties &

# Save Kafka server PID
KAFKA_PID=$!

# Wait for Kafka server to be ready
echo -e "$filePath Waiting for Kafka server to be available."
while ! nc -zv "${KAFKA_HOSTNAME}" "${KAFKA_SERVER_PORT}" 2>/dev/null; do
  echo -e "$filePath Checking again in $waitTime seconds."
  sleep $waitTime
done
echo -e "$filePath Kafka server is ready."

# Create Kafka topic(s)
echo -e "$filePath Creating topic(s) in Kafka server."
${KAFKA_HOME}/bin/kafka-topics.sh --create --bootstrap-server ${KAFKA_HOSTNAME}:${KAFKA_SERVER_PORT} --replication-factor 1 --partitions 1 --topic example_topic_1
${KAFKA_HOME}/bin/kafka-topics.sh --create --bootstrap-server ${KAFKA_HOSTNAME}:${KAFKA_SERVER_PORT} --replication-factor 1 --partitions 1 --topic example_topic_2
echo -e "$filePath Topic(s) have been created."

# Remove netcat-openbsd since no longer needed
echo -e "$filePath Cleaning up."
apt-get remove -y netcat-openbsd >/dev/null 2>&1 &&
  apt-get autoremove -y >/dev/null 2>&1 &&
  apt-get clean >/dev/null 2>&1 &&
  rm -rf /var/lib/apt/lists/* >/dev/null 2>&1

echo -e "$filePath Done."

# keep Kafka server running
wait $KAFKA_PID
