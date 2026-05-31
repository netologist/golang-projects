# Bulkhead

## Concept
Isolate resources into independent partitions so a failure or overload in one
partition cannot exhaust capacity for others — like watertight compartments in
a ship's hull.

## When to Use
- Multiple workloads share a process and must not starve each other.
- Protect critical traffic from being crowded out by batch work.

## When NOT to Use
- A single uniform workload — partitioning adds needless complexity.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Failure isolation between partitions | Lower peak utilisation per partition |
| Predictable capacity per workload | Requires sizing each partition |

## Go-Specific Notes
Each partition is a semaphore (`chan struct{}`). `Execute` rejects immediately
when a partition is saturated, preventing pile-ups.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- One semaphore per partition isolates concurrency.
- Reject fast when a partition is full (fail fast, not queue forever).
- Size partitions by workload criticality.
