# Rate Limiter

## Concept
Control the rate of operations using a token bucket: tokens refill at a fixed
rate up to a burst capacity, and each operation consumes one token.

## When to Use
- Protect downstream services or APIs from overload.
- Enforce per-client quotas.

## When NOT to Use
- You only need to bound concurrency (use a semaphore) rather than rate.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Smooths bursts to a steady rate | Requires tuning rate and burst |
| Allows short bursts up to capacity | Time-based; needs a monotonic clock |

## Go-Specific Notes
The advanced limiter implements a token bucket with `sync.Mutex` and refills
lazily on each call. `Wait` blocks until a token is available or the context ends.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Token bucket: refill rate + burst capacity.
- `Allow` is non-blocking; `Wait` blocks with context cancellation.
- For production, consider `golang.org/x/time/rate`.
