# Stage 1: Build Go binary for ARM64v8 using Alpine
FROM arm64v8/golang:alpine3.19 AS builder

# Set working directory inside container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Build the Go binary
RUN go build -o kafka-producer main.go

# Stage 2: Create a minimal image to run the Go binary
FROM arm64v8/alpine:3.19

# Install necessary certificates
RUN apk add --no-cache ca-certificates

# Set working directory inside the final container
WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/kafka-producer .

# Command to run the binary
ENTRYPOINT ["/entrypoint.sh"]
