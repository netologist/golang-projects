# Chord DHT

## Concept
A distributed hash table with O(log N) lookups. Nodes are arranged on a
circular identifier ring. Each node maintains a finger table — pointers to
nodes at exponentially increasing distances. Key lookups hop through the
finger table, halving the remaining distance each step.

## When to Use
- You need a scalable, self-organising key-value store in a P2P network.
- You want O(log N) lookup complexity without a central directory.

## When NOT to Use
- Your network is small and stable — a central directory is simpler.
- Nodes churn rapidly — finger table maintenance adds overhead.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Decentralised, self-healing | Churn causes finger table maintenance cost |
| O(log N) lookups | Inefficient for small networks |

## Go-Specific Notes
Nodes are integers on a 2^M ring. The finger table is precomputed on join.
`Lookup(key)` walks the ring via finger table hops and returns the
successor node plus hop count.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- The finger table provides exponential "shortcuts" — O(log N) hops.
- Each node only needs to know about O(log N) other nodes.
- Chord is the canonical academic DHT — widely referenced in distributed systems.
