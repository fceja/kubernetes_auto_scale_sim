# Finish Docker Compose Setup

## Build Remaining Docker Image

- Navigate to project Root and run:

  - ```bash
     docker build -t kafka-3.4.1-server.docker-compose:latest -f ./kafka/server/Dockerfile.docker-compose ./kafka/server
    ```

## Deploy with Docker Compose

- Deploy with docker-compose, at `~/project_root` run:

  - ```bash
    docker-compose -f docker-compose.yaml up
    ```

- Debugging: view config

  - ```bash
    docker-compose -f docker-compose.yaml config
    ```
