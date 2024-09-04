# Init Container Stack

## Build Image - Kafka Producer

- Navigate to ~/project_root/kafka_producer
  - `docker build -t kafka_producer .`

## Deploy Stack

- Note: following images must first be built before deploying stack

  - Grafana

    - ```bash
       docker build -t grafana-11.2.0 ./grafana
      ```

  - Kafdrop Server

    - ```bash
       docker build -t kafdrop-4.0.3-snapshot ./kafdrop_server
      ```

  - Kafka Consumer

    - ```bash
       docker build -t kafka-3.4.1-consumer ./kafka_consumer
      ```

  - Kafka Producer

    - ```bash
       docker build -t kafka-3.4.1-producer ./kafka_producer
      ```

  - Kafka Server

    - ```bash
       docker build -t kafka-3.4.1-server ./kafka_server
      ```

  - Prometheus

    - ```bash
       docker build -t prometheus-2.54.1 ./prometheus
      ```

  - Zookeeper

    - ```bash
       docker build -t zookeeper-3.9.2 ./zookeeper
      ```

  - All

    - ```bash
       && docker build -t grafana-11.2.0 ./grafana \
       && docker build -t kafdrop-4.0.3-snapshot:latest ./kafdrop_server \
       && docker build -t kafka-3.4.1-producer:latest ./kafka_producer \
       && docker build -t kafka-3.4.1-consumer:latest ./kafka_consumer \
       && docker build -t kafka-3.4.1-server.docker-compose:latest -f ./kafka_server/Dockerfile.docker-compose ./kafka_server \
       && docker build -t prometheus-2.54.1 ./prometheus \
       && docker build -t zookeeper-3.9.2:latest ./zookeeper \
      ```

  - Deploy with docker-compose

    - ```bash
      docker-compose -f docker-compose.yaml up
      ```

    - view config

    - ```bash
      docker-compose -f docker-compose.yaml config
      ```

## Running Local Kafka Producer / Consumer

Since Zookeeper and Kafka server are always ran in a docker container, they are are created via compose-docker stack and connect to the swarm managed network.

We need to add terminal ran `kafka-producer` and `kafka-consumer` to the network

\*\* Note - localhost:9092 worked.

- Identify Swarm Network
  - Look for network associated with your stack (e.g. `stackname_default`)
  - `docker network ls`
- Find IP Address of the Containers
  - Get IP address of the `kafka-server` container by inspecting
    - `docker ps`
    - `docker inspect <container_name_or_id>`
    - NetworkSettings > Networks > <my_stack_default> > IPAddress
- Optional - Allow External Access to Containers
- hat the container service that ports exposed, e.g., `"9092:9092"` for Kafka.
- Update Go application:
  - In you Go code, set the Kafka Broker address using the container's IP address and port
    - e.g. `192.168.x.x:9092`
  - Alternatively, if ports are exposed on the host, you can connect using `localhost:9092`
