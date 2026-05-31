# Circuit Breaker

## Concept
Stop calling a failing service. Track failures; when they exceed a threshold,
"trip" to Open and fail fast for a cooldown. After cooldown, allow a probe
(Half-Open); success closes the breaker, failure reopens it.

## States
- **Closed** — normal; failures are counted.
- **Open** — fail fast; no calls reach the downstream.
- **Half-Open** — a probe is allowed; success closes, failure reopens.

## When to Use
- Calling external services (HTTP, gRPC, DB) that fail transiently.
- You want to shed load and let a dependency recover.

## When NOT to Use
- In-process calls, or operations that must never be skipped.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Fails fast instead of waiting on timeouts | Thresholds need per-service tuning |
| Protects downstream from overload | Half-open probe may still hit users |

## Go-Specific Notes
State transitions are mutex-guarded. The advanced breaker emits Prometheus
metrics (state gauge, trip counter) using a local registry.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Three states: Closed -> Open -> Half-Open -> Closed.
- Tune failure threshold and recovery timeout per dependency.
- Emit metrics for breaker state.
