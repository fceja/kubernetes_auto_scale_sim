# Create Zookeeper Docker Image

## Build Image

`docker build -t zookeeper-arm64v8 -f Dockerfile .`

### Multi-platform Builds

For multi-platform builds, might need to use Docker Buildx.
To build an image for multiple platforms (e.g., linux/arm64), use:

`docker buildx build --platform linux/arm64 -t zookeeper-arm64 .`

## Run Docker Container

`docker run -d --name zookeeper-container -p 2181:2181 -p 9092:9092 zookeeper-arm64v8`
