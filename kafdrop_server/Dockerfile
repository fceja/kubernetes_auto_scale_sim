FROM eclipse-temurin:17.0.12_7-jdk

# Set environment variables
ENV KAFKA_HOSTNAME=kafka-server \
    KAFKA_SERVER_PORT=9092 \
    INIT_FILE_PATH=init-kafdrop-server.sh

WORKDIR /app

# Add Kafdrop files
ADD kafdrop*tar.gz /
ADD kafdrop.sh /
RUN chmod +x /kafdrop.sh

# Copy init script
COPY init-kafdrop-server.sh .
RUN chmod +x init-kafdrop-server.sh


EXPOSE 9000

ENTRYPOINT ["sh", "/app/init-kafdrop-server.sh"]
