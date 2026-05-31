# Leader Election

## Concept
Elect a single leader among a set of nodes so exactly one performs coordinated
work. Leadership is held via a time-bounded lease that the leader renews; if it
stops renewing (crash), another node takes over.

## When to Use
- A cluster needs one coordinator (scheduler, compactor, sequencer).
- Avoid duplicate work across replicas.

## When NOT to Use
- All nodes can safely act independently.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Single active coordinator | Brief gaps during failover |
| Automatic failover via lease | Needs careful lease/heartbeat tuning |

## Go-Specific Notes
This in-process simulation uses a shared `Cluster` lease. Each node campaigns on
a ticker; the lowest-ID live candidate wins. Crossing a lease boundary lets a new
leader emerge after the current one stops renewing.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Leadership is a renewable, time-bounded lease.
- Failover happens when the lease expires without renewal.
- Keep lease TTL > renewal interval to avoid flapping.
