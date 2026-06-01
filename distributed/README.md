# Distributed Systems Patterns in Go

A collection of distributed systems patterns implemented in idiomatic Go. Each
pattern includes a simple runnable example and an advanced version with tests.

## Categories

| Category | Description |
|----------|-------------|
| [Coordination](./coordination) | Paxos, Raft, 2PC, 3PC, leader election, distributed locks |
| [Observability](./observability) | Metrics, health checks, heartbeats, tracing, structured logging |
| [Resilience](./resilience) | Circuit breaker, retry/backoff, timeout, bulkhead, fallback, hedged requests |
| [Scalability](./scalability) | Consistent hashing, rendezvous hashing, load balancing, service discovery, sidecar, range partitioning |
| [P2P](./p2p) | Gossip protocol, Chord DHT, epidemic broadcast |
| [Replication](./replication) | CRDTs, WAL, primary-replica, read repair, multi-leader, quorum |
| [Communication](./communication) | Scatter-gather, pub-sub, request-reply, gRPC patterns |
| [Data](./data) | Outbox, CQRS, saga, saga choreography, event sourcing |

## Running

```bash
make test          # Run all tests
make test-verbose  # Verbose with race detection
make lint          # Lint all packages
make fmt           # Format with gofumpt + goimports
```

Each pattern can be run individually:

```bash
go run ./coordination/paxos/advanced
go test ./replication/quorum/advanced/... -v
```

## Key Takeaways

- Each pattern is self-contained with no cross-pattern dependencies.
- Simple examples focus on the core algorithm; advanced examples add tests and edge cases.
- Patterns are educational — not production-ready — and prioritise clarity over performance.
