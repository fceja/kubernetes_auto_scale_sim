# Apache Kafka

## \*\*Important\*\* - Support for Zookeeper

As of Kafka version 3.5 and higher, Kafka removed support for zookeeper. Kafka is making push for own service, KRaft mode (Kafka Raft Metadata mode).

## What is Apache Kafka?

Kafka is a distributed event store and streams-processing platform, meaning simply it takes data from producers and streams them out to consumers.

These producers and consumers can also be though of "inputs" and "outputs", where data is taken from an "input" system and consumed by an "output" system.

### Kafka Brokers

[Ref - openlogic](https://www.openlogic.com/blog/using-kafka-zookeeper)

The main vehicle for this movement of data is the Kafka broker. The Kafka broker handles all requests from all clients (both producers and consumers as well as metadata). It also manages replication of data across a cluster as well as within topics and partitions.
