# Apache Kafka

## \*\*Important\*\* - Support for Zookeeper

As of Kafka version 3.5 and higher, Kafka removed support for zookeeper. Kafka is making push for own service, KRaft mode (Kafka Raft Metadata mode).

## What is Apache Kafka?

Kafka is a distributed event store and streams-processing platform, meaning simply it takes data from producers and streams them out to consumers.

These producers and consumers can also be though of "inputs" and "outputs", where data is taken from an "input" system and consumed by an "output" system.

### Kafka Brokers

[Ref - openlogic](https://www.openlogic.com/blog/using-kafka-zookeeper)

The main vehicle for this movement of data is the Kafka broker. The Kafka broker handles all requests from all clients (both producers and consumers as well as metadata). It also manages replication of data across a cluster as well as within topics and partitions.

### Kafka Client

Go (AKA golang)

Kafka Version: 0.8.x
Maintainer: IBM
License: MIT
[GitHub - sarama](https://github.com/IBM/sarama)

\*\* used to be maintained by Shopify, however transferred to IBM.
Shopify favored bindings of `librdkafka` instead (C/C++)

### Topics

A `topic` is a name stream of records (or messages).
When you refer to `topics` in your code, you're dealing with the names of these streams, which organize and categorize data flowing through Kafka.

Topics:
These are the categories or channels to which data records are sent. i.g. `user-signups` or `order-updates`

Messages:
These are the actual data records or event that are published to a topic. Each message Typically contains a key, a value, and potentially some metadata.

### Producer

A `producer` is a component in the Apache Kafka ecosystem that is responsible fo sending records (messages) to Kafka to topics. `Producers` are essential the data entry point into Kafka.

### Consumer

A Kafka `consumer` subscribes to one or more Kafka topics and reads messages from them.

Consumer Group: a `consumer group` is a group of consumers that work together to consume messages from a Kafka topic. Each consumer in the group is assigned a partition of the topic to read from. Kafka ensures each partition is consumed by ony one consumer in the group at a time.

### Consuming Messages

You can consume messages from a topic in several ways:

- From the `beginning`: consume the earliest (oldest) message.
- From the `latest`: consume the latest (youngest) message.
- From a specific `offset`: consume specific offset (index)

### Managing Message Retention

To control how long messages are kept:
By default, kakfa retains all messages forever.

- `Retention Time`: Messages are keeps for a specific period (e.g., 7 days.)
- `Retention Size`: MEssages are kept until the log reaches a specific size.
- `Log Compaction`: Keeps only the latest message for each key, effectively "deleting" older versions.

### Balance strategies
