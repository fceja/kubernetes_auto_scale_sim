# Kafka Messages

## Unlimited Capacity

In Kafka, a topic can technically hold an unlimited number of messages.

- The number of messages that a topic can hold is constrained only by the dis space available on the Kafka brokers.
- Messages in a topic are stored in partitions, and as long as there's enough disk space, you can keep adding messages.

## Retention Policy

Kafka also has a configurable retention policy that determines how long messages are kept before they are deleted.

- You can set this policy based on
  - time (e.g. retain messages for 7 days)
  - based on size (e.g. retain up to 10 GB of data).
- Once the retention limit is reached older messages are deleted to make room for new ones.
