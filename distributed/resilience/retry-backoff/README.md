# Retry + Backoff

## Concept
A production retry primitive: exponential backoff with jitter, a caller-supplied
retryable predicate, and context-aware cancellation. Distinct from the language
retry pattern by adding generic typed results.

## When to Use
- Idempotent calls to flaky downstreams (network, throttled APIs).
- You want exponential backoff to avoid hammering a recovering service.

## When NOT to Use
- Non-idempotent operations; permanent errors.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Survives transient faults gracefully | Adds tail latency on failures |
| Jitter prevents retry storms | Misclassified errors waste attempts |

## Go-Specific Notes
`Do[T]` returns a typed result, retries based on `cfg.Retryable`, and sleeps with
`select` on the context so cancellation is immediate.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Exponential backoff + jitter is the default for distributed retries.
- Let the caller decide retryable vs permanent.
- Respect context cancellation between attempts.
