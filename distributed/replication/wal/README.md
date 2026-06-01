# Write-Ahead Log (WAL)

## Concept
All state mutations are durably recorded in a sequential log before being
applied. On crash recovery, the log is replayed from the last checkpoint,
rebuilding in-memory state. Checkpoints periodically truncate the log to
control its size.

## When to Use
- You need durability — no committed data is lost on crash.
- You're building a database, message queue, or state machine.

## When NOT to Use
- Data is ephemeral or can be reconstructed from other sources.
- The overhead of sequential disk writes is unacceptable.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Durability guarantee | Sequential write overhead on every mutation |
| Simple recovery model | Unbounded log growth without checkpoints |

## Go-Specific Notes
Entries have an LSN (log sequence number) and operation type (SET/DEL).
`Append()` writes an entry and applies it. A checkpoint at a given LSN
truncates all prior entries. `Recover()` replays from checkpoint onwards.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- WAL is the foundation of durable storage — write to log, then apply to state.
- Checkpoints prevent unbounded log growth.
- Recovery replays only the post-checkpoint suffix, minimising restart time.
