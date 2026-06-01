# Fallback

## Concept
Graceful degradation when a primary operation fails. Attempts the primary
function first; if it returns an error, executes a fallback function.
Supports chaining multiple fallbacks in sequence (primary → fallback1 → fallback2)
until one succeeds or all are exhausted.

## When to Use
- You want to provide a degraded-but-working response when the primary is unavailable.
- You have cached/stale data that can serve as a backup.

## When NOT to Use
- The fallback itself is expensive — prefer a circuit breaker to fail fast.
- You need strict correctness — stale fallback data may be misleading.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Improved availability and user experience | Returns potentially stale data |
| Simple to implement with a generic type | Every fallback adds latency when primary fails |

## Go-Specific Notes
Uses Go generics (`Fallback[T]`) so the same pattern works for any
return type. The `Chain` function iterates through functions until one
succeeds.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Fallback is the simplest resilience pattern — try A, if it fails try B.
- Chaining allows multiple tiers of degradation (live → cache → default).
- Fallback latency is hidden when the primary succeeds.
