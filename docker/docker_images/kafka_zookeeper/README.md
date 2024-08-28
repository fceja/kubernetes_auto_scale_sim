# Create Kafka and Zookeeper Docker Image

## Build Image

`docker build -t kafka-zookeeper-arm64 -f Dockerfile .`

### Multi-platform Builds

For multi-platform builds, might need to use Docker Buildx.
To build an image for multiple platforms (e.g., linux/arm64), use:

`docker buildx build --platform linux/arm64 -t kafka-zookeeper-arm64 .`

## Run Docker Container

`docker run -d --name kafka-zookeeper-container -p 2181:2181 -p 9092:9092 kafka-zookeeper-arm64`
