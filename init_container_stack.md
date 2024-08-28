# Init Container Stack

## Build Image - Kafka Producer

- Navigate to ~/project_root/kafka_producer
  - `docker build -t kafka_producer .`

## Deploy Stack

- Note: following images must first be built before deploying stack

  - Kafka Producer

    - ```bash
       docker build -t kafka_producer ./kafka_producer
      ```

  - Deploy stack

    - ```bash
       docker stack rm my_stack \
       && docker stack deploy -c docker-compose.yml my_stack
      ```

- All-in-one

  - Builds docker images, deletes previous stack, re-deploys stack.
  - Navigate to ~/project root

    - ```bash
       docker build -t kafka_producer ./kafka_producer \
       && docker stack rm my_stack \
       && docker stack deploy -c docker-compose.yml my_stack
      ```
