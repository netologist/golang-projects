# Distributed Lock

## Concept
Provide mutual exclusion across processes or nodes. A lock is acquired with a TTL
and a fencing token; only the holder of the current token may unlock or extend.
The TTL ensures a crashed holder's lock eventually frees.

## When to Use
- Only one node should perform an action at a time (leader task, migration).
- Coordinating access to a shared external resource.

## When NOT to Use
- In-process synchronisation — use `sync.Mutex`.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Cross-node mutual exclusion | TTL tuning: too short risks double-execution |
| Self-healing via TTL expiry | Clock skew can cause edge cases |

## Go-Specific Notes
The in-memory `MemoryLocker` models a lock server: a mutex-guarded map of
key->{token, expiry}. A fencing token protects against stale holders.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Always acquire with a TTL so crashed holders free the lock.
- Use a fencing token to reject stale unlock/extend calls.
- Extend the lease for long-running work (a watchdog).
