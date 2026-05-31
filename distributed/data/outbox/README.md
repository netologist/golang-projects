# Outbox Pattern

## Concept
Reliably publish events by writing them to an "outbox" within the same
transaction that mutates business state. A relay process later reads pending
outbox rows and publishes them, guaranteeing the event is never lost even if the
broker is down at write time.

## When to Use
- You must atomically change state AND emit an event (no dual-write race).
- At-least-once delivery is acceptable (consumers deduplicate).

## When NOT to Use
- Events are best-effort and loss is acceptable.

## Trade-offs
| Benefit | Cost |
|---------|------|
| No lost events on broker downtime | At-least-once (needs idempotent consumers) |
| Atomic state + event write | Requires a relay and polling |

## Implementation Note
To keep this example infrastructure-free and CGo-free, the store is in-memory
but models a single transaction that writes the record and the outbox event
together. In production this is a database table written in the business
transaction, with a relay polling for `published_at IS NULL`.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Write the event in the same transaction as the state change.
- A relay publishes pending events and marks them published.
- Consumers must be idempotent (at-least-once delivery).
