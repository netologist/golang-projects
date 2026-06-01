# Three-Phase Commit (3PC)

## Concept
A non-blocking extension of 2PC. Adds a PreCommit phase between Prepare and
Commit so that participants know the transaction outcome even if the
coordinator crashes. Uses timeouts to make progress without the coordinator.

## When to Use
- You need non-blocking atomic commit (participants don't block on coordinator failure).
- The network is reliable enough that timeouts work correctly.

## When NOT to Use
- Network partitions are common — 3PC can still block in partition scenarios.
- 2PC suffices and you value simplicity.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Non-blocking under coordinator failure | Extra round-trip (3 phases vs 2 in 2PC) |
| Participants can make progress via timeout | More complex implementation |

## Go-Specific Notes
The coordinator drives Prepare → PreCommit → Commit phases. An optional
crash recovery path simulates coordinator failure after PreCommit.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- The PreCommit phase lets participants proceed after coordinator failure.
- 3PC adds a phase to avoid the blocking problem of 2PC.
- Timeouts are critical — misconfigured timeouts can still cause blocking.
