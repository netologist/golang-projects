# Quorum

## Concept
Configure a replicated store with N replicas. Writes require W acks
(write quorum); reads require R acks (read quorum). When W + R > N,
every read sees the latest write. Tuning W and R lets you trade off
write latency, read latency, and consistency.

## When to Use
- You need tunable consistency — pick W and R per operation.
- You want to survive (N - W) write failures or (N - R) read failures.

## When NOT to Use
- W = N and R = 1 (write-all-read-one) is too slow — use eventual consistency.
- You don't have enough replicas for meaningful quorum sizes.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Tunable consistency per operation | Higher latency than single-node reads/writes |
| Survives partial failures | Requires careful W/R/N configuration |

## Go-Specific Notes
The store has N nodes, each with versioned values. `Write()` collects W
acks; `Read()` collects R responses and returns the highest-version value.
Nodes can be marked offline to simulate failures.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- W + R > N guarantees read-your-writes consistency.
- W = R = (N+1)/2 gives majority quorum (Paxos-style).
- Version tracking ensures the reader picks the freshest value from the quorum.
