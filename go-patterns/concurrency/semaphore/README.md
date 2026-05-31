# Semaphore

## Concept
Limit the number of goroutines accessing a resource concurrently. A buffered
channel of empty structs is the canonical Go semaphore: acquiring sends, releasing
receives.

## When to Use
- Cap concurrent access to a limited resource (DB connections, file handles, API).
- Throttle fan-out without a full worker pool.

## When NOT to Use
- You need weighted permits across many resource units — use x/sync/semaphore.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Tiny, dependency-free | Counting only (no weights) in the simple form |
| Context-aware acquire | Must always pair Acquire with Release |

## Go-Specific Notes
`chan struct{}` of capacity N is the simple semaphore. The advanced version uses
`golang.org/x/sync/semaphore` for weighted acquisition.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Acquire = send, Release = receive on a buffered channel.
- Always `defer Release()` after a successful Acquire.
- Use weighted semaphores when units cost differently.
