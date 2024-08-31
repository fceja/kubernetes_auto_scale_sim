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

  - Kafka Consumer

    - ```bash
       docker build -t kafka-3.4.1-consumer ./kafka_consumer
      ```

  - Kafdrop Server

    - ```bash
       docker build -t kafdrop-4.0.3-snapshot ./kafdrop_server
      ```

  - Kafka - All

    - ```bash
       docker build -t zookeeper-3.9.2 ./zookeeper \
       && docker build -t kafka-3.4.1-server ./kafka_server \
       && docker build -t kafka-3.4.1-producer ./kafka_producer \
       && docker build -t kafka-3.4.1-consumer ./kafka_consumer \
       && docker build -t kafdrop-4.0.3-snapshot ./kafdrop_server
      ```

  - Deploy stack

    - ```bash
       docker stack rm my_stack \
       && docker stack deploy -c docker-compose.yml my_stack
      ```

## Running Local Kafka Producer / Consumer

Since Zookeepr and Kafka server are always ran in a docker container, they are are created via compose-docker stack and connect to the swarm managed network.

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
