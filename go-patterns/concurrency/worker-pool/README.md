# Worker Pool

## Concept
Bound parallelism by running a fixed number of worker goroutines that pull jobs
from a shared queue (a buffered channel) and push results to a result channel.

## When to Use
- Many independent tasks where unbounded goroutines would exhaust resources.
- You need a steady, controllable level of concurrency.

## When NOT to Use
- Few tasks, or tasks that must run sequentially.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Bounded resource usage | Fixed pool can underutilise bursty load |
| Backpressure via channel capacity | Requires careful shutdown handling |

## Go-Specific Notes
The advanced pool exposes `Submit`, `Results`, and a `Shutdown` that drains
remaining work and closes the result channel exactly once.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Size the pool to the bottleneck resource (CPU cores, DB connections).
- Close the jobs channel to signal shutdown; a WaitGroup closes results.
- Always provide a context to abort in-flight submissions.
