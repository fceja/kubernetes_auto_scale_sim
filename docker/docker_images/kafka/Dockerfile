# Official Images for OpenJDK binaries built by Eclipse Temurin.
FROM arm64v8/eclipse-temurin:22.0.2_9-jdk

# NOTE - Kafka 3.4.1 supports zookeeper

# Set environment variables for Kafka version and home directory
ENV KAFKA_VERSION=3.4.1 \
    SCALA_VERSION=2.13 \
    KAFKA_HOME=/opt/kafka \
    KAFKA_DATA_DIR=/var/lib/kafka/data \
    KAFKA_LOG_DIR=/var/lib/kafka/logs

# Install required packages and clean up
RUN apt-get update && apt-get install -y wget curl gnupg \
    && wget https://archive.apache.org/dist/kafka/${KAFKA_VERSION}/kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz \
    && tar -xzf kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz -C /opt \
    && mv /opt/kafka_${SCALA_VERSION}-${KAFKA_VERSION} ${KAFKA_HOME} \
    && rm kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz \
    && apt-get remove -y wget curl gnupg \
    && apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create directories for Kafka data and logs
RUN mkdir -p ${KAFKA_DATA_DIR} ${KAFKA_LOG_DIR}

# Expose Kafka port
EXPOSE 9092

# Set working dir
WORKDIR ${KAFKA_HOME}

# Copy Kafka config file
COPY /server.properties ${KAFKA_HOME}/config/server.properties

# Define entrypoint to start Kafka
ENTRYPOINT ["bin/kafka-server-start.sh", "config/server.properties"]
