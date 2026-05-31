# Retry

## Concept
Re-attempt a transient operation with a backoff delay between tries, giving a
flaky dependency time to recover. Exponential backoff with jitter avoids
synchronised retry storms.

## When to Use
- Operations that fail transiently (network blips, throttling, contention).
- The operation is idempotent or safe to repeat.

## When NOT to Use
- Non-idempotent operations without deduplication.
- Permanent errors (validation, auth) — retrying wastes time.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Survives transient failures | Adds latency on failure paths |
| Jitter avoids thundering herds | Risk of retrying non-idempotent work |

## Go-Specific Notes
The advanced retry takes a `Retryable` predicate so callers decide which errors
are worth retrying, and respects `context.Context` cancellation between attempts.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Use exponential backoff with jitter.
- Let the caller classify retryable vs permanent errors.
- Always honour context cancellation between attempts.
