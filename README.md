# Description

A `Kubernetes cluster` project that simulates `Pod` auto-scaling.

There are two ways to run:

Option 1: using `docker-compose` to run as standalone Docker containers.

Option 2: using `kind` to run local Kubernetes clusters using Docker containers.

## Install Docker Desktop

- https://www.docker.com/products/docker-desktop/

### Optional

Some services require additional steps for installation.

- To skip additional installation for `Kafdrop`, you can simply comment out the `kafdrop-server` section from services in `docker-compose.yaml`
- To skip additional installation for `Prometheus`, you can simply comment out the `prometheus` and `grafana` sections from services in `docker-compose.yaml`

- Kafdrop setup

  - To run Kafdrop, you will need to prepare a `kafdrop .bin.tar.gz` file.

    - Follow the `Building` instructions at [GitHub-obsidiandynamics](https://github.com/obsidiandynamics/kafdrop).
    - The output will produce a `/target` directory.
      - Copy and paste `kafdrop-4.0.3-SNAPSHOT-bin.tar.gz` into `~/project_root/kafdrop`

- Prometheus setup

  - To run Prometheus, download `jmx_prometheus_javaagent-1.0.1.jar` from [Github-prometheus](https://github.com/prometheus/jmx_exporter/releases) and place into `~/project_root/kafka/server`.

## Build Docker Images

- Navigate to project Root and run:

  - Note - if you skipped the `kafdrop` or `prometheus` or `grafana` optional setups, you can also skip their build steps below.

  - ```bash
     docker build -t grafana-11.2.0 ./grafana \
     && docker build -t kafdrop-4.0.3-snapshot:latest ./kafdrop \
     && docker build -t kafka-3.4.1-producer:latest ./kafka/producer \
     && docker build -t kafka-3.4.1-consumer:latest ./kafka/consumer \
     && docker build -t kafka-3.4.1-server.docker-compose:latest -f ./kafka/server/Dockerfile.docker-compose ./kafka/server \
     && docker build -t prometheus-2.54.1 ./prometheus \
     && docker build -t zookeeper-3.9.2:latest ./zookeeper
    ```

## Deploy with Docker Compose

- Deploy with docker-compose

  - ```bash
    docker-compose -f docker-compose.yaml up
    ```

  - Debugging: view config

    - ```bash
      docker-compose -f docker-compose.yaml config
      ```

### Screenshot
<img src="/screenshots/docker-compose.png" alt="project screenshot for docker-compose" />
