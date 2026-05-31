# Scatter-Gather

## Concept
Send a request to multiple backends in parallel (scatter) and combine their
responses (gather). Return all results, or whatever arrived before the deadline.

## When to Use
- Query multiple replicas/shards and aggregate (search, price comparison).
- Aggregate data from several services for one response.

## When NOT to Use
- A single backend can answer — no need to fan out.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Lower latency via parallelism | Higher total load (N calls) |
| Partial results on deadline | Aggregation logic adds complexity |

## Go-Specific Notes
`Gather` launches one goroutine per backend and collects results on a channel,
returning partial results if the context deadline fires first.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Fan out with goroutines; collect on a buffered channel.
- Return partial results on deadline rather than failing entirely.
- Record per-source latency and error for observability.
