# Two-Phase Commit (2PC)

## Concept
Atomically commit a transaction across multiple participants. In phase 1
(Prepare) the coordinator asks every participant to vote. If all vote commit, the
coordinator tells everyone to Commit in phase 2; if any votes abort (or errors),
everyone Aborts.

## When to Use
- Strong atomicity is required across heterogeneous resource managers.
- All participants are available and reasonably reliable.

## When NOT to Use
- High availability matters more than strict atomicity — prefer Saga.
- Participants may be slow/unreliable (2PC blocks on the coordinator).

## Trade-offs
| Benefit | Cost |
|---------|------|
| Strong atomic guarantee | Blocking; coordinator is a failure point |
| Simple mental model | Poor availability under partitions |

## Go-Specific Notes
The coordinator drives the protocol over a `Participant` interface. Any prepare
error or abort vote aborts the whole transaction.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Phase 1 collects votes; phase 2 commits or aborts.
- A single abort vote aborts the transaction.
- 2PC trades availability for atomicity — consider Saga when that is wrong.
