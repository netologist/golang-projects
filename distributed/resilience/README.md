# Resilience

Patterns for building fault-tolerant systems that degrade gracefully
under failure.

| Pattern | Description |
|---------|-------------|
| [Circuit Breaker](./circuit-breaker) | Trip after threshold failures; prevent cascading failures |
| [Retry / Backoff](./retry-backoff) | Exponential backoff with jitter and retry budgets |
| [Timeout](./timeout) | Deadline propagation to prevent resource leaks |
| [Bulkhead](./bulkhead) | Resource isolation to contain failures to a partition |
| [Fallback](./fallback) | Graceful degradation with fallback responses |
| [Hedged Requests](./hedged-requests) | Race multiple requests; take the first success |
