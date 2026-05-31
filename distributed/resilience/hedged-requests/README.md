# Hedged Requests

## Concept
Reduce tail latency by issuing a backup request if the first has not returned
within a hedge delay. The first response wins; the rest are cancelled.

## When to Use
- Read-only, idempotent requests where p99 latency matters.
- Replicated backends where any replica can serve the request.

## When NOT to Use
- Non-idempotent writes; expensive operations (hedging doubles load).

## Trade-offs
| Benefit | Cost |
|---------|------|
| Cuts tail latency significantly | Extra load from duplicate requests |
| Simple with context cancellation | Only safe for idempotent calls |

## Go-Specific Notes
`Do[T]` spawns attempts staggered by the hedge delay, races them with
`select`, and cancels the losers via a shared cancellable context.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Hedge only idempotent requests.
- The first success cancels all other attempts.
- Tune the hedge delay to your p95 latency.
