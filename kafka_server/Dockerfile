# Official Images for OpenJDK binaries, built by Eclipse Temurin.
FROM arm64v8/eclipse-temurin:22.0.2_9-jdk

# NOTE - Kafka 3.4.1 supports zookeeper 3.9.2
# NOTE - Kafka versions after 3.4.1 no longer support zookeeper

# Set environment variables
ENV KAFKA_VERSION=3.4.1 \
    SCALA_VERSION=2.13 \
    KAFKA_HOME=/opt/kafka \
    KAFKA_DATA_DIR=/var/lib/kafka/data \
    KAFKA_LOG_DIR=/var/lib/kafka/logs \
    KAFKA_HOSTNAME=localhost \
    KAFKA_SERVER_PORT=9092 \
    ZOOKEEPER_HOSTNAME=zookeeper \
    ZOOKEEPER_SERVER_PORT=2181 \
    INIT_FILE_PATH=init-kafka-server.sh

# Install required packages and clean up
RUN apt-get update && apt-get install -y wget curl gnupg netcat-openbsd \
    && wget https://archive.apache.org/dist/kafka/${KAFKA_VERSION}/kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz \
    && tar -xzf kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz -C /opt \
    && mv /opt/kafka_${SCALA_VERSION}-${KAFKA_VERSION} ${KAFKA_HOME} \
    && rm kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz \
    && apt-get remove -y wget gnupg curl \
    && apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create directories for Kafka data and logs
RUN mkdir -p ${KAFKA_DATA_DIR} ${KAFKA_LOG_DIR}

# Expose Kafka port
EXPOSE 9092

# Set working dir
WORKDIR ${KAFKA_HOME}

# Copy config file
COPY /server.properties ${KAFKA_HOME}/config/server.properties

# Copy init script
COPY init-kafka-server.sh /usr/local/bin/init-kafka-server.sh
RUN chmod +x /usr/local/bin/init-kafka-server.sh

# Define entrypoint to Kafka server init script
ENTRYPOINT ["/usr/local/bin/init-kafka-server.sh"]