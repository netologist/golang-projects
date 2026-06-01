# Rendezvous Hashing

## Concept
Also known as Highest Random Weight (HRW) hashing. Each key is hashed against
every node; the node with the highest combined hash wins. When a node is added
or removed, only keys that mapped to that node are remapped — all others stay
put. Offers minimal disruption on membership changes.

## When to Use
- You want minimal key movement when nodes are added or removed.
- Your node set changes incrementally (not all at once).

## When NOT to Use
- You need O(1) lookup — rendezvous hashing is O(N) per key lookup.
- Your node count is very large — consistent hashing with virtual nodes scales better.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Minimal disruption on membership change | O(N) lookup (hash against every node) |
| No virtual nodes needed; uniform distribution | Slower than consistent hashing for large clusters |

## Go-Specific Notes
Each key is hashed against all nodes using a combined hash function.
The `Get()` method returns the node with the highest score. `GetN()`
returns the top-K nodes for replication.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Only keys that mapped to the removed/added node move — everything else stays.
- Uniform distribution without virtual nodes.
- Best for medium-sized clusters where node changes are infrequent.
