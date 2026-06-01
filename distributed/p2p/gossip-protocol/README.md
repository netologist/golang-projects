# Gossip Protocol

## Concept
Epidemic-style information dissemination. Each node forwards information to a
random subset of peers (fanout) each round. Information spreads exponentially,
reaching all nodes in O(log N) rounds. Used for membership, failure detection,
and data propagation in large clusters.

## When to Use
- You need to disseminate information across a large, dynamic cluster.
- You want decentralised, fault-tolerant communication (no single coordinator).

## When NOT to Use
- You need strict ordering or exactly-once delivery.
- The cluster is small enough that direct messaging is simpler.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Scalable to thousands of nodes | Eventual — no guarantee of when all nodes receive |
| Decentralised, no single point of failure | Redundant messages (same info from multiple peers) |

## Go-Specific Notes
The cluster tracks infection state per node. Each round, infected nodes
forward to a random subset of peers. Convergence is detected when all
live nodes are infected.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Exponential spread: with fanout f, the cluster converges in O(log_f N) rounds.
- Decentralised — no coordinator needed; every node acts independently.
- Trade off latency for robustness and scale.
