# Event Sourcing

## Concept
Persist state as an ordered log of immutable events rather than current values.
Current state is derived by replaying events. This yields a full audit trail and
the ability to rebuild or time-travel state.

## When to Use
- You need a complete audit history of every change.
- Temporal queries ("what was the balance last Tuesday?").
- Pairs naturally with CQRS for read projections.

## When NOT to Use
- Simple CRUD with no audit requirement.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Full audit log, replayable | More storage; eventual read models |
| Optimistic concurrency via versions | Steeper learning curve |

## Go-Specific Notes
Events are an interface; the store appends with an expected-version check for
optimistic concurrency. Aggregates apply events to fold state.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- State = fold over the event log.
- Append-only store with expected-version concurrency control.
- Snapshots optimise replay for long streams.
