# Consistent Hashing

## Concept
Map keys to nodes on a hash ring so that adding or removing a node only remaps a
small fraction of keys (roughly 1/N), instead of nearly all of them as with
modulo hashing. Virtual nodes smooth the key distribution.

## When to Use
- Distributed caches and sharded stores where nodes come and go.
- Any partitioning that must minimise reshuffling on membership change.

## When NOT to Use
- A fixed, never-changing set of nodes.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Minimal remap on node changes | More complex than modulo |
| Even load via virtual nodes | Ring lookups cost O(log n) |

## Go-Specific Notes
The ring is a sorted slice of hashed virtual-node points; `Get` binary-searches
for the first point >= the key's hash, wrapping around at the end.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Use many virtual nodes per physical node for even distribution.
- Lookup is a binary search for the successor on the ring.
- Removing a node only affects its own key range.
