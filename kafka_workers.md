# Workers

## Worker Failure

When working with kafka consumers, tracking offsets and handling worker failures are key concerns.

- Offset Storage

  - `Consumer Offsets` - kafka stores offsets in special internal topic called `__consumer_offsets`. Each consumer group maintains its offsets in this topic. This allows kafka to track which messages have been process by which consumer group.

  - `Commit Offsets` - consumers periodically commit their offsets to kafka. This can be done automatically or manually. When using automatic commits, kafka manages offsets for you, but you can configure the frequency of commits. `Manual commits give you more control`, allowing you to commit offsets only after you've successfully processed messages.

- Worker Failure and Offset Recovery:

  - `Consumer Groups and Failover`: in kafka consumer group, each partition is assigned to only one consumer (worker) at a time. If a worker fails, kafka's consumer group coordinator will reassign the partitions to other available consumers in the group.

  - `Rebalancing`: when a worker fails, kafka triggers a rebalance process. During rebalance, kafka reassigns the partitions from the failed worker to other workers in the consumer group. The new worker will resume consuming message starting from the last committed offset.

  - `Offset Tracking`: since offsets are committed to kafka, when a new worker takes over a partition, it reads the committed offset from kafka's `__consumer_offsets` topic and resumes from there. This ensures that the new worker picks up from where the failed worker left off, minimizing the risk of message loss.

- Handling Offset Commit Failures

  - `Retry Logic`: Implement retry logic in you consumer code for committing offsets. If committing offsets fails due to temporary issues, retries ensure that the commit eventually succeeds.

  - `Idempotent Processing`: ensure you message processing is idempotent, meaning that processing the same message multiple times doesn't have adverse effects. This is crucial for scenarios where a message might be reprocessed after a worker failure.

Example Scenario:

1. `Normal Operation`:

   - Worker A is processing messages from partition 0 and has committed an offset of 100.
   - Worker B is processing messages from partition 1.

2. `Worker A Fails`:

   - Kafka detects the failure and triggers a rebalance.
   - Worker B or another available worker takes over processing partition 0.

3. `Resuming Processing`:

   - The new worker reads the last committed offset (100) from `__consumer_offsets`.
   - The new worker resumes processing from offset 100.
