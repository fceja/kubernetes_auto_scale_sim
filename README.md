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

  - Kafka - All

    - ```bash
       docker build -t zookeeper-3.9.2 ./zookeeper \
       && docker build -t kafka-3.4.1-server ./kafka_server \
       && docker build -t kafka-3.4.1-producer ./kafka_producer
      ```

  - Deploy stack

    - ```bash
       docker stack rm my_stack \
       && docker stack deploy -c docker-compose.yml my_stack
      ```
