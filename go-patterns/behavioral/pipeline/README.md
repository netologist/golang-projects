# Pipeline

## Concept
Chain processing stages where each stage reads from an input channel and writes
to an output channel. Stages run concurrently as goroutines connected by channels.

## When to Use
- Stream processing with independent transformation stages.
- You want back-pressure and bounded memory via channels.

## When NOT to Use
- A simple loop suffices — channels add overhead and complexity.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Natural concurrency between stages | Goroutine + channel overhead |
| Composable, testable stages | Must handle cancellation in every stage |

## Go-Specific Notes
Each stage owns its output channel and closes it on return. Every send is
guarded by a `select` on `ctx.Done()` to avoid goroutine leaks.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Each stage closes the channel it owns.
- Guard every send with `ctx.Done()` to prevent leaks.
- The advanced version parallelises a stage across N workers.
