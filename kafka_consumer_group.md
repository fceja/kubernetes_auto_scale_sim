# Consumer Group

In Kafka, consumer groups is a set of consumers that cooperate to consume data from one or more topics. Each consumer in the group has a unique ID, known as the group ID, which identifies the group.

## Creating a Consumer Group

To "create" a `consumer group`, you just need to start a Kafka consumer with a specified `group.id`. The consumer group is automatically created by Kafka when the consumer starts consuming message from a topic.

## Check If a Consumer Group Exists

You can check if a `consumer group` exists using the Kafka command-line tools or through a Kafka client library.

- Kafka Command-Line tools
  - List all `consumer groups`
    - `kafka-consumer-groups.sh --bootstrap-server localhost:9092 --group my-consumer-group --describe`
  - Describe a specific `consumer group`
    - `kafka-consumer-groups.sh --bootstrap-server localhost:9092 --group my-consumer-group --describe`

### Consumer Group Worker Assignment

When you add new workers with the same `consumer group ID`, Kafka will automatically handle the assignment of partitions to those workers.

- `Automatic Partition Assignment`
  - When workers (consumers) join the same `consumer group`, Kafka automatically rebalances the partitions among the available workers.
  - Means that each worker will be assigned one or more partitions, ensuring that all partitions are consumed but by only one worker in the group at a time.
- `Partition Rebalancing`
  - As you add more workers to the `consumer group`, Kafka will trigger a rebalance process.
  - The partitions will be redistributed among the new set of workers. This is done automatically, so you don't have to manually configure with worker gets which partition.
- `Manual Partitioning`
  - While Kafka handles this automatically, you can also manually assign specific partitions to specific consumers if you need fine-grained control. However, this is not necessary for most use cases where automatic partition assignment works well.

### Message Retrieval

It is standard practice to retrieve and process all available messages from the assigned partition in Kafka. Kafka's design is optimized for high-throughput, streaming data, where consumers typically process messages as they arrive in real-time. By consumning all messages in a partition, you ensure that no data is missed and that the consumer group can keep up with the flow of data.

- In `func ConsumeClaim(...)` -> `claim.Messages()`
