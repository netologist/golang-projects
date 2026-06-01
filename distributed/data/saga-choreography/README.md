# Saga Choreography

## Concept
A decentralised saga pattern where services react to events rather than
being orchestrated by a central coordinator. Each service subscribes to
relevant events, performs its local transaction, and publishes new events
for downstream services. On failure, compensating events (refunds, cancels)
propagate through the same event bus.

## When to Use
- Services are autonomous and prefer event-driven communication.
- The saga flow is simple enough that event chains are easy to reason about.
- You want loose coupling — no orchestrator knows the full workflow.

## When NOT to Use
- The workflow is complex with many branches — orchestration is easier to trace.
- You need a single place to monitor saga state.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Decentralised, loosely coupled | Harder to visualise the full workflow |
| No single point of failure | Event chains can become convoluted |

## Go-Specific Notes
An event bus manages subscriptions and event routing. Each service subscribes
and returns zero or more follow-up events. The saga tracks each step's
success/failure in a record. Compensation handlers (refund, cancel) are
also event subscribers.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Choreography replaces a central coordinator with event chains.
- Each service knows only its own step and what events to emit next.
- Compensation events undo work on failure — same pattern, opposite direction.
