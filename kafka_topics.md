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
