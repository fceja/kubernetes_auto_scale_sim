# Kafka Topics

## Scaling Horizontally

Horizontal Scaling refers to adding more workers (consumers) to process the data in a topic.

- When you scale horizontally, you distribute the workload across multiple workers.
  - Adding more workers to a consumer group allows you to process more messages simultaneously.
  - Each worker will be assigned one or more partitions, and they will read and process the messages independently.

## Impact of Horizontal Scaling

- `Partition Handling`
  - Kafka automatically rebalances the partitions among the available workers when you add or remove workers from a consumer group.
- `Fault Tolerance`
  - If one worker fails, Kafka can redistribute the partitions it was responsible for to the remaining workers, ensuring the processing continues without interruption.

## Partitions

In Kafka, the number of partitions for a topic must be defined at the time of topic creation. It cannot be dynamically adjusted based on the number of consumer group workers. However, you can manually increase the number of partitions after the topic has been created. But this is a manual process and does not happen automatically as more consumers are added to the group.

### Why Partitions Are Static At Creation

Partitions are a fundamental unit of parallelism in Kafka. When a topic is created, you decide the number or partitions, which dictates how many consumers in a consumer group can consume in parallel. Once a partition is created, Kafka doesn't automatically split or merge partitions.

### Making Partition Dynamic

While Kafka itself doesn't support dynamic partitioning, based on the number of consumers, you can implement a mechanism to monitor the load and manually adjust partitions as needed. Here are some strategies:

- `Set a Higher Number of Partitions Initially`
  - When creating a topic, you can set a higher number of partitions that initially needed. This allows Kafka to distribute the load among more consumers if you consumer group scales up in the future.
- `Manually Increase Partitions`
  - If you notice that your application need more parallelism, you can manually increase the number of partitions using Kafka's CLI tool or programmatically using the admin client.
  - `kafka-topics.sh --bootstrap-server localhost:9092 --alter --topic example_topic_2 --partitions 4`
- `Monitor and Adjust`
  - You can write a monitoring tool that checks the load on your consumers and triggers a parititon increase if the load becomes too high. This would be an external process that periodically checks Kafka metrics and makes adjustments as necessary.

## Replication Factor

- `ReplicationFactor`
  - In Kafka defines the number of replicas each partition will be maintained across different brokers.
  - Each replica is stored on a different broker.
  - Replication Factor 1
    - means on the leader replica exists for a partition.
    - only one copy of the partition's data stored in the Kafka cluster.
  - Replication Factor 2
    - means on the one leader replica and one follower replica for each partition.
    - there are two copies of the partition's data stored in the Kafka cluster.
    - if broker holding the leader fails, the follower can be promoted to leader.
  - Replication Factor n
    - means one leader and n - 1 followers for each partition
    - n cannot be zero
      - n > 0
      - zero is not allowed
- Purpose
  - This is an aspect of Kafka's design to ensure fault tolerance and high availability.

### Roles

#### Leader Replica

- Responsibility
  - Handles all read and write requests for the partition. When a producer writes messages to a partition, it writes to the leader replica. Consumers also read messages from the leader.
- Writes
  - All new messages are written to the leader replica first.

#### Follower Replicas

- Responsibility
  - Replicate the messages from the leader to ensure data redundancy and fault tolerance.
- Replication
  - Followers fetch data from the leader and keep their data in sync. They do not handle client requests directly.
- Reads

  - Followers do not handle reads; they ony store copies of the data for redundancy. Reads are always served by the leader.

### Process Flow

- Message Write
  - A producer sends a message to the leader of a partition.
  - The leader writes the message to its log and then replicates it to its followers.
- Replication
  - Followers periodically fetch messages from the leaders and append them to their logs.
  - This replication process ensures that the follower replicas have the same data as the leader.
- Failover

  - If the leader fails, one of the followers can be promoted to leader. This ensures that the partition remains available and continues to process requests.

### Key Points

- Synchronization
  - Followers are synchronized with the leader. They contain the same messages but do not process requests themselves.
- Fault Tolerance
  - If the leader fails, one of the followers will be promoted to leader, ensuring continue availability of the partition.
- No Independent Processing
  - Followers do not process messages independently or perform any operations on the messages. They are purely for replication and redundancy.
