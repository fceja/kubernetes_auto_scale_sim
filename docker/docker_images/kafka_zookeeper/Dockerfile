# Base image that supports arm64 architecture
FROM arm64v8/eclipse-temurin:22.0.2_9-jdk

# Set versions
ENV KAFKA_VERSION=3.4.1
ENV ZOOKEEPER_VERSION=3.9.2
ENV SCALA_VERSION=2.13

# Install packages
RUN apt-get update && \
    apt-get install -y wget tar && \
    rm -rf /var/lib/apt/lists/*

# Install ZooKeeper
RUN wget https://downloads.apache.org/zookeeper/zookeeper-${ZOOKEEPER_VERSION}/apache-zookeeper-${ZOOKEEPER_VERSION}.tar.gz && \
    tar -xzf apache-zookeeper-${ZOOKEEPER_VERSION}.tar.gz && \
    mv apache-zookeeper-${ZOOKEEPER_VERSION} /opt/zookeeper && \
    rm apache-zookeeper-${ZOOKEEPER_VERSION}.tar.gz

# Install Kafka
RUN wget https://archive.apache.org/dist/kafka/${KAFKA_VERSION}/kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz && \
    tar -xzf kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz && \
    mv kafka_${SCALA_VERSION}-${KAFKA_VERSION} /opt/kafka && \
    rm kafka_${SCALA_VERSION}-${KAFKA_VERSION}.tgz

# Set environment variables
ENV KAFKA_HOME=/opt/kafka
ENV ZOOKEEPER_HOME=/opt/zookeeper

# Copy configuration files to docker container
COPY zookeeper.properties /opt/zookeeper/conf/zookeeper.properties
COPY kafka.properties /opt/kafka/config/kafka.properties

# Expose necessary ports
EXPOSE 2181 9092

# Set default command to start ZooKeeper and Kafka
# CMD ["sh", "-c", "nohup $ZOOKEEPER_HOME/bin/zkServer.sh start $ZOOKEEPER_HOME/conf/zookeeper.properties & nohup $KAFKA_HOME/bin/kafka-server-start.sh $KAFKA_HOME/config/server.properties & tail -f /dev/null"]
CMD ["sh", "-c", "$ZOOKEEPER_HOME/bin/zkServer.sh start $ZOOKEEPER_HOME/conf/zookeeper.properties && until echo ruok | nc localhost 2181 | grep imok; do sleep 2; done && $KAFKA_HOME/bin/kafka-server-start.sh $KAFKA_HOME/config/server.properties"]
