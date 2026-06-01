# Range Partitioning

## Concept
Data is partitioned by key range into shards. Each shard owns a contiguous
key space (e.g., [a, h), [h, p), [p, ∞)). Supports dynamic rebalancing
through shard splits — a shard that grows too large is split into two at a
chosen split key, redistributing only the split shard's keys.

## When to Use
- Queries frequently access ranges of keys (e.g., time-series, alphabetical).
- You need to scale by splitting hot shards on demand.

## When NOT to Use
- Access patterns are uniformly random — hash-based partitioning is simpler.
- You can't tolerate the overhead of rebalancing operations.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Efficient range scans | Can create hot spots if keys cluster |
| Dynamic splits enable elastic scaling | Rebalancing requires coordination |

## Go-Specific Notes
Shards are modelled as structs with Start/End boundaries (lexicographic).
`Route()` finds the owning shard via binary search. `RebalanceSplit()`
splits a shard at a midpoint and creates a new shard ID.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Range partitioning is the foundation of sorted distributed stores.
- Splits allow dynamic rebalancing without full redistribution.
- Hot spots are the main challenge — monitor shard sizes and split proactively.
