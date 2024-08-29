# Stage 1: Build Go binary for ARM64v8 using Alpine
FROM arm64v8/golang:alpine3.19 AS builder

# Set working directory inside temp container
WORKDIR /app

# Copy source code into container
COPY . .

# Build Go binary
RUN go build -o kafka-producer main.go

# Stage 2: Create minimal image to run Go binary
FROM arm64v8/alpine:3.19

# Set environment variables
ENV KAFKA_HOSTNAME=kafka-server \
    KAFKA_SERVER_PORT=9092 \
    INIT_FILE_PATH=init-kafka-producer.sh

# Set working directory inside final container
WORKDIR /app

# Install certificates
RUN apk add --no-cache ca-certificates netcat-openbsd

# Copy Go binary from builder stage and init script
COPY --from=builder /app/kafka-producer .
COPY init-kafka-producer.sh .

# Run Kafka producer startup script
ENTRYPOINT ["sh", "/app/init-kafka-producer.sh"]