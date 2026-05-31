# errgroup

## Concept
Run a group of goroutines, wait for them all, and capture the first error.
`golang.org/x/sync/errgroup` cancels a shared context on the first failure,
stopping the remaining work early.

## When to Use
- Concurrent subtasks where any failure should abort the rest.
- Fan-out calls (e.g., parallel HTTP requests) that share a deadline.

## When NOT to Use
- You must collect every error, not just the first (aggregate manually).

## Trade-offs
| Benefit | Cost |
|---------|------|
| First error cancels the group | Only the first error is returned |
| Clean Wait() semantics | Need a semaphore to bound parallelism |

## Go-Specific Notes
The advanced wrapper adds a weighted semaphore so the group runs at most N tasks
at once while preserving first-error cancellation.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- `errgroup.WithContext` cancels siblings on first error.
- Combine with `semaphore.Weighted` to bound concurrency.
- Return wrapped errors so callers can inspect the cause.
