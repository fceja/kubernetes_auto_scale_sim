# Description

A `Kubernetes cluster` project that simulates `Pod` auto-scaling. Written in Go and a couple of bash scripts.

## Additional Info
\*\*NOTE - There are two ways to run:

Option 1: using `docker-compose` to run as standalone Docker containers.

Option 2: using `kind` to run `local Kubernetes clusters` using Docker containers.

This page contains common installation steps for both options. For remaining installation/deployment steps:

- For Option 1, view `~/project_root/README-Docker-Compose-Installation.md`.
- For Option 2, view `~/project_root/README-Kubernetes-Installation.md`.

## Common Installation Steps

### Install Docker Desktop

- https://www.docker.com/products/docker-desktop/

#### Optional

- `kafdrop`, `prometheus` and `grafana` require additional installation steps.

  - You can simply skip these by commenting them out from the services section in `docker-compose.yaml`

- To install:

  - Kafdrop

    - To run Kafdrop, you will need to prepare a `kafdrop .bin.tar.gz` file.
      - Follow the `Building` instructions at [GitHub-obsidiandynamics](https://github.com/obsidiandynamics/kafdrop).
      - The output will produce a `/target` directory.
      - Copy and paste `kafdrop-4.0.3-SNAPSHOT-bin.tar.gz` into `~/project_root/kafdrop`

  - Prometheus
    - To run Prometheus, download `jmx_prometheus_javaagent-1.0.1.jar` from [Github-prometheus](https://github.com/prometheus/jmx_exporter/releases) and place into `~/project_root/kafka/server`.

## Build Docker Images

- Navigate to project Root and run:

  - Note - if you skipped the `kafdrop` or `prometheus` or `grafana` optional setups, you can also skip their build steps below.

  - ```bash
     docker build -t grafana-11.2.0 ./grafana \
     && docker build -t kafdrop-4.0.3-snapshot:latest ./kafdrop \
     && docker build -t kafka-3.4.1-producer:latest ./kafka/producer \
     && docker build -t kafka-3.4.1-consumer:latest ./kafka/consumer \
     && docker build -t prometheus-2.54.1 ./prometheus \
     && docker build -t zookeeper-3.9.2:latest ./zookeeper
    ```

## Remaining Installation / Deployment

For remaining installation/deployment steps:

- For Option 1, view `~/project_root/README-Docker-Compose-Installation.md`.
- For Option 2, view `~/project_root/README-Kubernetes-Installation.md`.

## Screenshots 

### Docker Compose

<img src="/screenshots/docker-compose.png" alt="project screenshot for docker-compose" />

### Kubernetes Cluster

<img src="/screenshots/kubernetes-cluster.png" alt="project screenshot for kubernetes cluster" />
