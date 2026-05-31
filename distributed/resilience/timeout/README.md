# Timeout

## Concept
Bound how long an operation may run using `context.WithTimeout`. When the
deadline passes, the operation is cancelled and the caller gets a timeout error
instead of blocking indefinitely.

## When to Use
- Any call that could hang: network, DB, external services.
- To enforce a latency budget across a request.

## When NOT to Use
- Fast, purely local computation with no blocking.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Prevents unbounded blocking | The callee must honour context cancellation |
| Enables latency budgets | A too-tight timeout causes spurious failures |

## Go-Specific Notes
The generic `Do` runs fn in a goroutine and selects on the result or the
timeout. The callee should observe `ctx.Done()` to actually stop work.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Use `context.WithTimeout` and always `defer cancel()`.
- The callee must respect the context to stop early.
- Distinguish deadline-exceeded from other errors.
