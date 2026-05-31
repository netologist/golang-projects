# Service Discovery

## Concept
Let services find each other dynamically. Instances register with a registry
(with a TTL heartbeat); clients discover the current healthy instances and can
watch for changes instead of using hard-coded addresses.

## When to Use
- Dynamic environments where instances scale up/down or move.
- Client-side discovery without a central proxy.

## When NOT to Use
- Static, fixed infrastructure with known addresses.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Adapts to dynamic topology | Registry is a dependency to operate |
| Watches push changes to clients | TTL tuning affects staleness |

## Go-Specific Notes
The in-process registry stores instances with a TTL, runs a reaper to expire
stale ones, and supports `Watch` via a channel that receives the current instance
list on change.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Register with a TTL; heartbeat to stay alive.
- A reaper removes instances that stop heartbeating.
- Watch streams topology changes to clients.
