# Stage 1: Build Go binary for ARM64v8 using Alpine
FROM arm64v8/golang:alpine3.19 AS builder

# Set working directory inside temp container
WORKDIR /app

# Copy source code into container
COPY . .

# Build Go binary
RUN go build -o kafka-worker main.go

# Stage 2: Create a minimal image to run Go binary
FROM arm64v8/alpine:3.19

# Set working directory inside final container
WORKDIR /app

# Install certificates
RUN apk add --no-cache ca-certificates

# Copy Go binary from builder stage
COPY --from=builder /app/kafka-worker .

# Command to run binary
ENTRYPOINT ["sh", "-c", "./kafka-worker"]
