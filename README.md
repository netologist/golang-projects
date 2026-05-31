# Go Patterns & Distributed Systems

A runnable reference library of Go language patterns and distributed system
patterns. Each pattern has a `simple/` (minimal, zero external deps) and
`advanced/` (production-grade, tested) variant.

## Quick Start

```bash
go run ./go-patterns/creational/functional-options/simple
go run ./go-patterns/concurrency/worker-pool/advanced
go run ./distributed/resilience/circuit-breaker/advanced
```

## Go Language Patterns (`go-patterns/`)

### Creational
| Pattern | simple | advanced |
|---------|--------|----------|
| [Functional Options](go-patterns/creational/functional-options/) | ✓ | ✓ |
| [Builder](go-patterns/creational/builder/) | ✓ | ✓ |
| [Singleton](go-patterns/creational/singleton/) | ✓ | ✓ |
| [Factory](go-patterns/creational/factory/) | ✓ | ✓ |

### Structural
| Pattern | simple | advanced |
|---------|--------|----------|
| [Decorator](go-patterns/structural/decorator/) | ✓ | ✓ |
| [Adapter](go-patterns/structural/adapter/) | ✓ | ✓ |
| [Proxy](go-patterns/structural/proxy/) | ✓ | ✓ |
| [Middleware Chain](go-patterns/structural/middleware-chain/) | ✓ | ✓ |

### Behavioral
| Pattern | simple | advanced |
|---------|--------|----------|
| [Iterator](go-patterns/behavioral/iterator/) | ✓ | ✓ |
| [Observer](go-patterns/behavioral/observer/) | ✓ | ✓ |
| [Strategy](go-patterns/behavioral/strategy/) | ✓ | ✓ |
| [Pipeline](go-patterns/behavioral/pipeline/) | ✓ | ✓ |
| [Fan-out/Fan-in](go-patterns/behavioral/fan-out-fan-in/) | ✓ | ✓ |

### Concurrency
| Pattern | simple | advanced |
|---------|--------|----------|
| [Worker Pool](go-patterns/concurrency/worker-pool/) | ✓ | ✓ |
| [Semaphore](go-patterns/concurrency/semaphore/) | ✓ | ✓ |
| [Rate Limiter](go-patterns/concurrency/rate-limiter/) | ✓ | ✓ |
| [Context Propagation](go-patterns/concurrency/context-propagation/) | ✓ | ✓ |
| [errgroup](go-patterns/concurrency/errgroup/) | ✓ | ✓ |

### Error Handling
| Pattern | simple | advanced |
|---------|--------|----------|
| [Sentinel Errors](go-patterns/error-handling/sentinel-errors/) | ✓ | ✓ |
| [Error Wrapping](go-patterns/error-handling/error-wrapping/) | ✓ | ✓ |
| [Result Type](go-patterns/error-handling/result-type/) | ✓ | ✓ |
| [Retry](go-patterns/error-handling/retry/) | ✓ | ✓ |

## Distributed System Patterns (`distributed/`)

### Resilience
| Pattern | simple | advanced |
|---------|--------|----------|
| [Circuit Breaker](distributed/resilience/circuit-breaker/) | ✓ | ✓ |
| [Bulkhead](distributed/resilience/bulkhead/) | ✓ | ✓ |
| [Timeout](distributed/resilience/timeout/) | ✓ | ✓ |
| [Retry + Backoff](distributed/resilience/retry-backoff/) | ✓ | ✓ |
| [Hedged Requests](distributed/resilience/hedged-requests/) | ✓ | ✓ |

### Communication
| Pattern | simple | advanced |
|---------|--------|----------|
| [gRPC Patterns](distributed/communication/grpc-patterns/) | ✓ | ✓ |
| [Pub/Sub](distributed/communication/pub-sub/) | ✓ | ✓ |
| [Request-Reply](distributed/communication/request-reply/) | ✓ | ✓ |
| [Scatter-Gather](distributed/communication/scatter-gather/) | ✓ | ✓ |

### Data
| Pattern | simple | advanced |
|---------|--------|----------|
| [Saga](distributed/data/saga/) | ✓ | ✓ |
| [Outbox](distributed/data/outbox/) | ✓ | ✓ |
| [Event Sourcing](distributed/data/event-sourcing/) | ✓ | ✓ |
| [CQRS](distributed/data/cqrs/) | ✓ | ✓ |

### Coordination
| Pattern | simple | advanced |
|---------|--------|----------|
| [Leader Election](distributed/coordination/leader-election/) | ✓ | ✓ |
| [Distributed Lock](distributed/coordination/distributed-lock/) | ✓ | ✓ |
| [Two-Phase Commit](distributed/coordination/two-phase-commit/) | ✓ | ✓ |

### Observability
| Pattern | simple | advanced |
|---------|--------|----------|
| [Distributed Tracing](distributed/observability/distributed-tracing/) | ✓ | ✓ |
| [Health Check](distributed/observability/health-check/) | ✓ | ✓ |
| [Structured Logging](distributed/observability/structured-logging/) | ✓ | ✓ |
| [Metrics](distributed/observability/metrics/) | ✓ | ✓ |

### Scalability
| Pattern | simple | advanced |
|---------|--------|----------|
| [Consistent Hashing](distributed/scalability/consistent-hashing/) | ✓ | ✓ |
| [Sidecar](distributed/scalability/sidecar/) | ✓ | ✓ |
| [Service Discovery](distributed/scalability/service-discovery/) | ✓ | ✓ |
| [Load Balancer](distributed/scalability/load-balancer/) | ✓ | ✓ |

## Layout Convention

Every pattern directory contains:

```
<pattern>/
├── README.md          concept, when to use, trade-offs, Go-specific notes
├── simple/            minimal runnable demo, zero external deps
└── advanced/          production-grade implementation + table-driven tests
```

## Running Tests

```bash
go test ./go-patterns/... ./distributed/...
go vet ./go-patterns/... ./distributed/...
```

## Rust Equivalent

See [`docs/rust-learning-prompt.md`](docs/rust-learning-prompt.md) for a
comprehensive prompt to build the Rust equivalent of this library.
