# Official Images for OpenJDK binaries built by Eclipse Temurin.
FROM arm64v8/eclipse-temurin:22.0.2_9-jdk

# Set environment variables
ENV ZOOKEEPER_VERSION=3.9.2 \
    ZOOKEEPER_HOME=/opt/zookeeper \
    ZOOKEEPER_DATA_DIR=/var/lib/zookeeper/data \
    ZOOKEEPER_DATA_LOG_DIR=/var/lib/zookeeper/log \
    ZOOKEEPER_PORT=2181 \
    ZOOKEEPER_TICK_TIME=2000

# Download and extract Zookeeper
RUN apt-get update && apt-get install -y wget gnupg \
    && wget https://archive.apache.org/dist/zookeeper/zookeeper-${ZOOKEEPER_VERSION}/apache-zookeeper-${ZOOKEEPER_VERSION}-bin.tar.gz \
    && tar -xzf apache-zookeeper-${ZOOKEEPER_VERSION}-bin.tar.gz -C /opt \
    && mv /opt/apache-zookeeper-${ZOOKEEPER_VERSION}-bin ${ZOOKEEPER_HOME} \
    && rm apache-zookeeper-${ZOOKEEPER_VERSION}-bin.tar.gz \
    && apt-get remove -y wget gnupg \
    && apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Create directories for Zookeeper data and logs
RUN mkdir -p ${ZOOKEEPER_DATA_DIR} ${ZOOKEEPER_DATA_LOG_DIR}

# Copy Zookeeper configuration files to container
COPY /zoo.cfg ${ZOOKEEPER_HOME}/conf/zoo.cfg

# Expose Zookeeper port
EXPOSE ${ZOOKEEPER_PORT}

# Set working directory
WORKDIR ${ZOOKEEPER_HOME}

# Define entrypoint to start Zookeeper
ENTRYPOINT ["bin/zkServer.sh", "start-foreground"]
