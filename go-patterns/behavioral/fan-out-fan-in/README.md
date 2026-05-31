# Fan-out / Fan-in

## Concept
Fan-out distributes work from one channel across multiple worker goroutines.
Fan-in merges multiple result channels back into one. Together they parallelise
a workload and then collect the results.

## When to Use
- A CPU- or IO-bound workload that is embarrassingly parallel.
- You want to bound parallelism with a fixed number of workers.

## When NOT to Use
- The work is trivial — the coordination overhead dominates.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Parallel throughput | Result order is not preserved |
| Simple to reason about with channels | Requires careful channel closing |

## Go-Specific Notes
Each worker shares one input channel (fan-out). A `sync.WaitGroup` closes the
merged output once all workers finish (fan-in). The advanced version is generic.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Fan-out: N goroutines read from the same input channel.
- Fan-in: a WaitGroup closes the merged channel after all sources drain.
- Results are unordered — add an index if order matters.
