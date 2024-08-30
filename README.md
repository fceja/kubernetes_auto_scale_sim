# Init Container Stack

## Build Image - Kafka Producer

- Navigate to ~/project_root/kafka_producer
  - `docker build -t kafka_producer .`

## Deploy Stack

- Note: following images must first be built before deploying stack

  - Zookeeper

    - ```bash
       docker build -t zookeeper-3.9.2 ./zookeeper
      ```

  - Kafka Server

    - ```bash
       docker build -t kafka-3.4.1-server ./kafka_server
      ```

  - Kafka Producer

    - ```bash
       docker build -t kafka-3.4.1-producer ./kafka_producer
      ```

  - Kafka Worker

    - ```bash
       docker build -t kafka-3.4.1-worker ./kafka_worker
      ```

  - Kafka - All

    - ```bash
       docker build -t zookeeper-3.9.2 ./zookeeper \
       && docker build -t kafka-3.4.1-server ./kafka_server \
       && docker build -t kafka-3.4.1-producer ./kafka_producer \
       && docker build -t kafka-3.4.1-worker ./kafka_worker
      ```

  - Deploy stack

    - ```bash
       docker stack rm my_stack \
       && docker stack deploy -c docker-compose.yml my_stack
      ```

## Running Local Kafka Producer / Worker

Since Zookeepr and Kafka server are always ran in a docker container, they are are created via compose-docker stack and connect to the swarm managed network.

We need to add terminal ran `kafka-producer` and `kafka-worker` to the network

\*\* Note - localhost:9092 worked.

1. Identify Swarm Network

- Look for network associated with your stack (e.g. `stackname_default`)

- `docker network ls`

2. Find IP Address of the Containers

- Get IP address of the `kafka-server` container by inspecting
  - `docker ps`
  - `docker inspect <container_name_or_id>`
  - NetworkSettings > Networks > <my_stack_default> > IPAddress

3. Optional - Allow External Access to Containers

- If you can't directly connect via IP Address, you can expose the containers ports to your host machine. Ensure that the container service that ports exposed, e.g., `"9092:9092"` for Kafka.

4. Update Go application:

- In you Go code, set the Kafka Broker address using the container's IP address and port
  - e.g. `192.168.x.x:9092`
- Alternatively, if ports are exposed on the host, you can connect using `localhost:9092`
