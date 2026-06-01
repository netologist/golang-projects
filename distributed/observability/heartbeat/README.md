# Heartbeat

## Concept
Fault detection through periodic heartbeat signals. Each node sends regular
heartbeats to a monitor. If a node misses a configurable number of consecutive
heartbeats (threshold), the monitor flags it as dead. Provides a status
callback on transitions (alive → suspect → dead).

## When to Use
- You need to detect node failures in a distributed system.
- You want configurable sensitivity (interval, threshold) for different environments.

## When NOT to Use
- The system can rely on external health check infrastructure (e.g., Kubernetes probes).
- The network is so unreliable that heartbeats would generate excessive false positives.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Simple, widely understood mechanism | False positives under network jitter |
| Configurable sensitivity | Adds constant background traffic |

## Go-Specific Notes
The monitor runs a goroutine checking heartbeats at a fixed tick interval.
Nodes send beats via a thread-safe `Beat()` method. A status callback
reports state transitions.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Heartbeats are the simplest form of fault detection.
- Threshold (missed beats) and interval control sensitivity vs. noise.
- Combine with health checks for comprehensive node monitoring.
