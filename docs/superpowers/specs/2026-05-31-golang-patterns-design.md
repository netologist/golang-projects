# Go Patterns & Distributed Systems Reference Repository — Design Spec

**Date:** 2026-05-31  
**Status:** Approved  
**Scope:** Go language patterns + Distributed system patterns, implemented as a mono-workspace reference library with simple and advanced variants per pattern.

---

## 1. Goals

1. Build a comprehensive, runnable reference library of Go language patterns and distributed system patterns.
2. Every pattern must be self-contained: `go run ./simple` and `go run ./advanced` both work with no external infrastructure (or with clearly documented local deps).
3. Each pattern includes a README that explains the concept, when to use it, trade-offs, and Go-specific notes — suitable as a learning resource.
4. After all Go patterns are complete, produce a detailed Rust learning prompt (`docs/rust-learning-prompt.md`) that maps every pattern to its Rust equivalent.

---

## 2. Repository Structure

```
golang-projects/
├── go.work                          # Go workspace root (go 1.22+)
├── README.md                        # Master index linking all patterns
│
├── go-patterns/                     # Module: github.com/veribaz/go-patterns
│   ├── go.mod
│   ├── README.md
│   ├── internal/
│   │   └── testutil/                # Shared test helpers
│   ├── creational/
│   │   ├── functional-options/
│   │   ├── builder/
│   │   ├── singleton/
│   │   └── factory/
│   ├── structural/
│   │   ├── decorator/
│   │   ├── adapter/
│   │   ├── proxy/
│   │   └── middleware-chain/
│   ├── behavioral/
│   │   ├── iterator/
│   │   ├── observer/
│   │   ├── strategy/
│   │   ├── pipeline/
│   │   └── fan-out-fan-in/
│   ├── concurrency/
│   │   ├── worker-pool/
│   │   ├── semaphore/
│   │   ├── rate-limiter/
│   │   ├── context-propagation/
│   │   └── errgroup/
│   └── error-handling/
│       ├── sentinel-errors/
│       ├── error-wrapping/
│       ├── result-type/
│       └── retry/
│
├── distributed/                     # Module: github.com/veribaz/distributed
│   ├── go.mod
│   ├── README.md
│   ├── internal/
│   │   ├── testutil/
│   │   └── metrics/                 # Shared Prometheus bootstrap
│   ├── resilience/
│   │   ├── circuit-breaker/
│   │   ├── bulkhead/
│   │   ├── timeout/
│   │   ├── retry-backoff/
│   │   └── hedged-requests/
│   ├── communication/
│   │   ├── grpc-patterns/
│   │   ├── pub-sub/
│   │   ├── request-reply/
│   │   └── scatter-gather/
│   ├── data/
│   │   ├── saga/
│   │   ├── outbox/
│   │   ├── event-sourcing/
│   │   └── cqrs/
│   ├── coordination/
│   │   ├── leader-election/
│   │   ├── distributed-lock/
│   │   └── two-phase-commit/
│   ├── observability/
│   │   ├── distributed-tracing/
│   │   ├── health-check/
│   │   ├── structured-logging/
│   │   └── metrics/
│   └── scalability/
│       ├── consistent-hashing/
│       ├── sidecar/
│       ├── service-discovery/
│       └── load-balancer/
│
└── docs/
    ├── rust-learning-prompt.md      # Rust equivalent prompt (post Go implementation)
    └── superpowers/
        └── specs/
            └── 2026-05-31-golang-patterns-design.md
```

---

## 3. Per-Pattern Layout Convention

Every pattern directory follows this exact structure:

```
<pattern-name>/
├── README.md
├── simple/
│   ├── main.go          # go run . — demonstrates pattern core
│   └── <pattern>.go     # clean, minimal implementation (~100-200 lines)
└── advanced/
    ├── main.go          # go run . — realistic production scenario
    ├── <pattern>.go     # production-grade implementation
    ├── <pattern>_test.go # table-driven tests
    └── metrics.go       # Prometheus metrics (where applicable)
```

### README.md Template (per pattern)

```markdown
# <Pattern Name>

## Concept
What problem this pattern solves and the core idea behind it.

## When to Use
Conditions and scenarios where this pattern is the right choice.

## When NOT to Use
Counter-indications and simpler alternatives.

## Trade-offs
| Benefit | Cost |
|---------|------|
| ...     | ...  |

## Go-Specific Notes
How Go's features (goroutines, interfaces, channels, context) shape this pattern.

## Running
\`\`\`bash
# Simple version
go run ./simple

# Advanced version
go run ./advanced
\`\`\`

## Key Takeaways
3–5 bullet points summarising what to remember.

## Further Reading
Links to blog posts, papers, or Go standard library source.
```

---

## 4. Module Configuration

### go.work
```
go 1.22

use (
    ./go-patterns
    ./distributed
)
```

### go-patterns/go.mod
```
module github.com/veribaz/go-patterns

go 1.22

require (
    golang.org/x/sync v0.7.0
    go.uber.org/zap v1.27.0         // advanced logging examples only
)
```

### distributed/go.mod
```
module github.com/veribaz/distributed

go 1.22

require (
    google.golang.org/grpc v1.63.0
    go.opentelemetry.io/otel v1.26.0
    go.opentelemetry.io/otel/trace v1.26.0
    github.com/prometheus/client_golang v1.19.0
    golang.org/x/sync v0.7.0
    go.uber.org/zap v1.27.0
)
```

---

## 5. Pattern Catalogue

### 5.1 Go Language Patterns (`go-patterns/`)

#### Creational
| Pattern | Core Idea | Key Go Feature |
|---------|-----------|---------------|
| Functional Options | Configure structs without constructor explosion | Variadic func options |
| Builder | Step-by-step object construction with validation | Method chaining on pointer receiver |
| Singleton | Single instance with lazy initialisation | `sync.Once` |
| Factory | Create objects behind an interface | Interface + type switch / registry map |

#### Structural
| Pattern | Core Idea | Key Go Feature |
|---------|-----------|---------------|
| Decorator | Wrap behaviour without modifying the original | Interface embedding |
| Adapter | Bridge incompatible interfaces | Struct wrapping |
| Proxy | Control access — lazy load, auth, caching | Interface with forwarding struct |
| Middleware Chain | Composable request/response pipeline | `http.Handler` chaining / function wrapping |

#### Behavioral
| Pattern | Core Idea | Key Go Feature |
|---------|-----------|---------------|
| Iterator | Sequential access without exposing internals | `iter.Seq` (Go 1.22) or channel-based |
| Observer | Decouple event producers from consumers | Channel fan-out |
| Strategy | Swap algorithms at runtime | Function types as first-class values |
| Pipeline | Chain processing stages | Channels + goroutines |
| Fan-out/Fan-in | Distribute then collect | Multiple goroutines + `sync.WaitGroup` |

#### Concurrency
| Pattern | Core Idea | Key Go Feature |
|---------|-----------|---------------|
| Worker Pool | Bound parallelism | Buffered channel as job queue |
| Semaphore | Limit concurrent access | `chan struct{}` or `golang.org/x/sync/semaphore` |
| Rate Limiter | Control throughput | `time.Ticker` / token bucket |
| Context Propagation | Cancellation + deadlines across goroutines | `context.Context` |
| errgroup | Wait for goroutines + collect first error | `golang.org/x/sync/errgroup` |

#### Error Handling
| Pattern | Core Idea | Key Go Feature |
|---------|-----------|---------------|
| Sentinel Errors | Named error values for comparison | `errors.Is` |
| Error Wrapping | Add context without losing the original | `fmt.Errorf("%w")` + `errors.As` |
| Result Type | Explicit success/failure return | Multiple return values |
| Retry | Re-attempt on transient failures | Backoff loop + context cancellation |

---

### 5.2 Distributed System Patterns (`distributed/`)

#### Resilience
| Pattern | Core Idea | Advanced Deps |
|---------|-----------|--------------|
| Circuit Breaker | Stop calling a failing service | State machine (Closed/Open/Half-Open), Prometheus |
| Bulkhead | Isolate failures to one partition | Worker pools per service, semaphores |
| Timeout | Bound wait time | `context.WithTimeout` |
| Retry + Backoff | Retry with exponential backoff + jitter | Configurable strategy, circuit breaker integration |
| Hedged Requests | Duplicate slow requests | Context cancellation race |

#### Communication
| Pattern | Core Idea | Advanced Deps |
|---------|-----------|--------------|
| gRPC Patterns | Unary + streaming + interceptors | `google.golang.org/grpc`, OTel interceptor |
| Pub/Sub | Decouple producer/consumer via topic | In-memory broker (simple), Redis-backed (advanced) |
| Request-Reply | Async request with reply correlation | Correlation ID + channel registry |
| Scatter-Gather | Fan-out to N services, aggregate results | `errgroup` + deadline |

#### Data
| Pattern | Core Idea | Advanced Deps |
|---------|-----------|--------------|
| Saga | Distributed transaction with compensation | Orchestration-based state machine |
| Outbox | Reliable event publishing via DB | SQLite (simple), Postgres (advanced) |
| Event Sourcing | State as ordered event log | Append-only store + snapshot |
| CQRS | Separate read and write models | Command/query bus + projection |

#### Coordination
| Pattern | Core Idea | Advanced Deps |
|---------|-----------|--------------|
| Leader Election | One node acts as primary | Heartbeat + lease (in-memory cluster sim) |
| Distributed Lock | Mutual exclusion across nodes | In-memory (simple), Redis SETNX (advanced) |
| Two-Phase Commit | Atomic across multiple participants | Coordinator + participant state machine |

#### Observability
| Pattern | Core Idea | Advanced Deps |
|---------|-----------|--------------|
| Distributed Tracing | Trace request across services | OTel SDK + Jaeger exporter |
| Health Check | Liveness + readiness endpoints | `net/http` handler + dependency checks |
| Structured Logging | Machine-readable, context-enriched logs | `log/slog` (simple), `zap` (advanced) |
| Metrics | RED metrics (Rate, Errors, Duration) | Prometheus client + Grafana dashboard config |

#### Scalability
| Pattern | Core Idea | Advanced Deps |
|---------|-----------|--------------|
| Consistent Hashing | Minimal key remapping on node change | Virtual nodes ring |
| Sidecar | Co-located helper process | In-process sidecar (simple), subprocess (advanced) |
| Service Discovery | Dynamic node registration/lookup | In-memory registry with TTL |
| Load Balancer | Distribute requests across backends | Round-robin (simple), least-connections (advanced) |

---

## 6. Quality Constraints

- Every `simple/` and `advanced/` must compile and run with `go run .` — no broken examples.
- All `advanced/` implementations must include `go test ./...` passing tests.
- No test may require external infrastructure unless clearly documented with a `docker-compose.yml` or `make setup` target.
- All code formatted with `gofumpt`. `go vet ./...` must pass.
- Imports ordered: stdlib → external → internal.
- Error wrapping follows `fmt.Errorf("operation: %w", err)` throughout.

---

## 7. Rust Learning Prompt

File: `docs/rust-learning-prompt.md`

Content:
- Full list of all 35 patterns with their Rust equivalents and recommended crates
- Go → Rust mental model mapping table (goroutines → tokio tasks, channels → `tokio::sync::mpsc`, interfaces → traits, `context.Context` → `CancellationToken`, etc.)
- Per-pattern prompt template: "Implement <pattern> in Rust using <crate>. Show simple version and production-grade version. Include a README explaining ownership/lifetime considerations."
- Curated crate list by category
- Repo structure recommendation for the Rust equivalent workspace (Cargo workspace)

---

## 8. Implementation Order

Implement in this sequence to build dependency knowledge progressively:

1. Workspace bootstrap (`go.work`, `go.mod` files, `README.md`, master index)
2. `go-patterns/` — all 20 patterns (creational → structural → behavioral → concurrency → error-handling)
3. `distributed/resilience/` — foundational (circuit-breaker, retry, timeout first)
4. `distributed/communication/`
5. `distributed/data/`
6. `distributed/coordination/`
7. `distributed/observability/`
8. `distributed/scalability/`
9. `docs/rust-learning-prompt.md`

---

## 9. Out of Scope

- No Kubernetes operator code in this repo (separate concern).
- No production deployment manifests.
- No external infrastructure required for `simple/` variants.
- No UI or HTTP service for the Go patterns module (distributed module has HTTP examples only where the pattern requires it).
