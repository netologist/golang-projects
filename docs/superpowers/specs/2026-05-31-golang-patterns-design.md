# Go Patterns & Distributed Systems Reference Repository

**Date:** 2026-05-31  
**Status:** Approved  
**Modules:** `github.com/veribaz/go-patterns` · `github.com/veribaz/distributed`  
**Go version:** 1.22+

---

## 1. Purpose & Success Criteria

This repository is a **runnable reference library** — not a framework, not a toolkit.
Every pattern exists to answer one question clearly: *"How do I do this idiomatically in Go?"*

A pattern implementation succeeds when:

1. `go run ./simple` runs with no setup and prints output that makes the pattern's mechanics obvious.
2. `go run ./advanced` runs with realistic scenario — real error paths, metrics, structured logs.
3. `go test ./advanced/...` passes with table-driven tests covering the happy path, the failure path, and at least one edge case.
4. A junior Go developer can read `README.md` and understand *why* this pattern exists before touching the code.
5. `go vet ./...` and `gofumpt -l .` produce zero output.

---

## 2. Repository Layout

```
golang-projects/
├── go.work                          # workspace: go 1.22
├── README.md                        # master index + quick-start
│
├── go-patterns/                     # 20 Go language patterns
│   ├── go.mod
│   ├── README.md
│   ├── internal/
│   │   └── testutil/
│   │       ├── assert.go            # thin assertion helpers (no testify)
│   │       └── clock.go             # fake clock for time-dependent tests
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
├── distributed/                     # 15 distributed system patterns
│   ├── go.mod
│   ├── README.md
│   ├── internal/
│   │   ├── testutil/
│   │   │   ├── server.go            # in-process test server helpers
│   │   │   └── clock.go
│   │   └── telemetry/
│   │       ├── metrics.go           # Prometheus registry bootstrap
│   │       └── tracing.go           # OTel provider bootstrap
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
    ├── rust-learning-prompt.md
    └── superpowers/specs/
        └── 2026-05-31-golang-patterns-design.md
```

---

## 3. Per-Pattern File Convention

```
<pattern>/
├── README.md
├── simple/
│   ├── main.go          # package main — runnable demo
│   └── <pattern>.go     # the pattern, no external deps, ~100-200 lines
└── advanced/
    ├── main.go          # realistic scenario, structured logs, graceful shutdown
    ├── <pattern>.go     # production-grade: configurable, observable
    ├── <pattern>_test.go # table-driven, uses internal/testutil
    └── metrics.go       # Prometheus metrics (where meaningful)
```

### README.md Structure (every pattern)

```markdown
# <Pattern Name>

## Concept
One paragraph. What problem it solves and the insight behind the solution.

## When to Use
Bulleted conditions. Be specific — not "when you need flexibility" but
"when you have ≥3 optional parameters and zero-value isn't always safe."

## When NOT to Use
The counter-cases. Most patterns are overused.

## Trade-offs
| Benefit                        | Cost                             |
|-------------------------------|----------------------------------|
| ...                           | ...                              |

## Go-Specific Notes
How Go's type system, goroutines, interfaces, or standard library
shape this pattern differently than in other languages.

## Running
\`\`\`bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
\`\`\`

## Key Takeaways
- 3–5 bullets. What to remember when you close the tab.

## Further Reading
- Link to relevant Go standard library source
- Link to canonical blog post or paper
```

---

## 4. Module Dependencies

### go-patterns/go.mod

```
module github.com/veribaz/go-patterns

go 1.22

require (
    golang.org/x/sync v0.7.0    // errgroup, semaphore
)
```

Rationale: keep deps minimal. `log/slog` (stdlib since 1.21) replaces zap in simple variants. The advanced decorator and middleware patterns use `net/http` from stdlib only.

### distributed/go.mod

```
module github.com/veribaz/distributed

go 1.22

require (
    google.golang.org/grpc              v1.63.0
    google.golang.org/protobuf          v1.34.0
    go.opentelemetry.io/otel            v1.26.0
    go.opentelemetry.io/otel/trace      v1.26.0
    go.opentelemetry.io/otel/sdk        v1.26.0
    go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.26.0
    github.com/prometheus/client_golang v1.19.0
    golang.org/x/sync                   v0.7.0
    go.uber.org/zap                     v1.27.0
)
```

Rationale: no Redis, no Postgres, no Kafka in simple variants. Advanced variants that need persistence use SQLite via `modernc.org/sqlite` (pure Go, no CGo). Advanced variants that need a message broker use an in-process implementation so `go run .` requires no Docker.

### go.work

```
go 1.22

use (
    ./go-patterns
    ./distributed
)
```

---

## 5. Pattern Catalogue with Technical Contracts

### 5.1 Go Language Patterns

#### Creational

**Functional Options**
```go
// Core contract
type Option func(*Server)
func WithTimeout(d time.Duration) Option
func WithMaxConns(n int) Option
func New(addr string, opts ...Option) *Server

// simple: 3 options, NewServer, print config
// advanced: 10 options, validation in Apply(), thread-safe after-build freeze
```

**Builder**
```go
// Core contract — builder returns error at Build(), not at each setter
type QueryBuilder struct { ... }
func (b *QueryBuilder) Select(cols ...string) *QueryBuilder
func (b *QueryBuilder) Where(cond string, args ...any) *QueryBuilder
func (b *QueryBuilder) Build() (Query, error)

// simple: SQL query builder
// advanced: + parameter binding validation, immutable Query result, benchmark
```

**Singleton**
```go
// Core contract
var instance *DB
var once sync.Once
func GetDB() *DB { once.Do(func() { instance = connect() }); return instance }

// simple: logger singleton
// advanced: resettable singleton for tests, health-aware getInstance
```

**Factory**
```go
// Core contract — registry-based, not switch-based
type Codec interface { Encode(any) ([]byte, error); Decode([]byte, any) error }
var registry = map[string]func() Codec{}
func Register(name string, fn func() Codec)
func New(name string) (Codec, error)

// simple: JSON/text codec factory
// advanced: plugin-style registry with validation, default fallback
```

---

#### Structural

**Decorator**
```go
// Core contract — wraps interface, adds behaviour
type Logger interface { Log(msg string) }
type TimedLogger struct { inner Logger; label string }
func (t *TimedLogger) Log(msg string) { start := time.Now(); t.inner.Log(msg); ... }

// simple: logging decorator around a Writer
// advanced: http.Handler decorator chain — auth, logging, metrics, recovery
```

**Adapter**
```go
// Core contract — our interface, their implementation
type OurStorage interface { Save(key, val string) error; Load(key string) (string, error) }
type RedisAdapter struct { client TheirRedisClient }

// simple: adapt os.File to an io.Writer interface subset
// advanced: adapt a legacy sync API to an async interface with context support
```

**Proxy**
```go
// Core contract — same interface, different behaviour
type ImageLoader interface { Load(url string) ([]byte, error) }
type CachingProxy struct { real ImageLoader; cache map[string][]byte; mu sync.RWMutex }

// simple: lazy-loading proxy
// advanced: caching proxy + access-control proxy composed, with TTL eviction
```

**Middleware Chain**
```go
// Core contract
type HandlerFunc func(ctx context.Context, req Request) (Response, error)
type Middleware func(HandlerFunc) HandlerFunc
func Chain(h HandlerFunc, mw ...Middleware) HandlerFunc

// simple: three middleware functions chained manually
// advanced: ordered chain with panic recovery, request-id injection, structured logging
```

---

#### Behavioral

**Iterator**
```go
// Core contract — Go 1.22 range-over-func
type Iter[T any] func(yield func(T) bool)
func PagedIter(fetch func(page int) ([]T, bool)) Iter[T]

// simple: slice iterator with filter/map
// advanced: paginated HTTP API iterator with backpressure, context cancellation
```

**Observer**
```go
// Core contract
type Event struct { Topic string; Payload any }
type Handler func(Event)
type Bus struct { ... }
func (b *Bus) Subscribe(topic string, h Handler) (cancel func())
func (b *Bus) Publish(ctx context.Context, e Event)

// simple: synchronous in-process event bus
// advanced: async bus with buffered delivery, slow-subscriber detection, metrics
```

**Strategy**
```go
// Core contract — function type IS the strategy
type SortStrategy[T any] func([]T) []T
type Sorter[T any] struct { strategy SortStrategy[T] }
func (s *Sorter[T]) SetStrategy(fn SortStrategy[T])
func (s *Sorter[T]) Sort(data []T) []T

// simple: sorting with swappable strategy
// advanced: compression strategy — gzip/snappy/zstd selected by payload size
```

**Pipeline**
```go
// Core contract — each stage is a func(in <-chan T) <-chan U
func Generate(ctx context.Context, vals ...int) <-chan int
func Square(ctx context.Context, in <-chan int) <-chan int
func Sink(ctx context.Context, in <-chan int) []int

// simple: integer transform pipeline
// advanced: file-processing pipeline with parallelism per stage, error channel, back-pressure
```

**Fan-out / Fan-in**
```go
// Core contract
func FanOut[T any](ctx context.Context, in <-chan T, n int, proc func(T) T) []<-chan T
func FanIn[T any](ctx context.Context, chans ...<-chan T) <-chan T

// simple: fan-out 1 producer to 3 workers, fan-in results
// advanced: dynamic worker scaling, ordered vs unordered merge modes
```

---

#### Concurrency

**Worker Pool**
```go
// Core contract
type Job struct { ID int; Payload any }
type Result struct { JobID int; Output any; Err error }
type Pool struct { ... }
func New(workers int) *Pool
func (p *Pool) Submit(ctx context.Context, job Job) error
func (p *Pool) Results() <-chan Result
func (p *Pool) Shutdown()

// simple: fixed goroutine pool, buffered job channel
// advanced: resizable pool, graceful drain on shutdown, Prometheus metrics (active workers, queue depth, job latency)
```

**Semaphore**
```go
// Core contract
type Semaphore struct { ch chan struct{} }
func New(n int) *Semaphore
func (s *Semaphore) Acquire(ctx context.Context) error
func (s *Semaphore) Release()
func (s *Semaphore) TryAcquire() bool

// simple: limit concurrent HTTP fetches
// advanced: weighted semaphore (golang.org/x/sync/semaphore), priority queue
```

**Rate Limiter**
```go
// Core contract — token bucket
type Limiter struct { ... }
func New(rate float64, burst int) *Limiter
func (l *Limiter) Wait(ctx context.Context) error  // blocks until token available
func (l *Limiter) Allow() bool                      // non-blocking check

// simple: manual token bucket with time.Ticker
// advanced: per-key limiter map with LRU eviction, HTTP middleware integration
```

**Context Propagation**
```go
// Core contract — custom context keys and value extraction
type ctxKey string
const RequestIDKey ctxKey = "request-id"
func WithRequestID(ctx context.Context, id string) context.Context
func RequestIDFrom(ctx context.Context) (string, bool)

// simple: trace a request ID through 3 layers
// advanced: context value hierarchy — auth user, request ID, tenant — with middleware injection
```

**errgroup**
```go
// Core contract
g, ctx := errgroup.WithContext(context.Background())
g.Go(func() error { return fetchA(ctx) })
g.Go(func() error { return fetchB(ctx) })
err := g.Wait()

// simple: 3 concurrent HTTP calls, first error cancels rest
// advanced: errgroup + semaphore for bounded parallelism, structured error aggregation
```

---

#### Error Handling

**Sentinel Errors**
```go
// Core contract
var ErrNotFound = errors.New("not found")
var ErrUnauthorised = errors.New("unauthorised")

func Get(id string) (*Record, error)  // wraps sentinel with context
errors.Is(err, ErrNotFound)           // caller inspection

// simple: store with typed sentinel errors
// advanced: error code registry, HTTP status mapping, Is/As chain
```

**Error Wrapping**
```go
// Core contract
fmt.Errorf("get user %d: %w", id, ErrNotFound)  // wrapping
errors.As(err, &target)                           // unwrapping to custom type

// simple: 3-layer call stack, unwrap at top
// advanced: custom error type with stack trace, JSON serialisation for API responses
```

**Result Type**
```go
// Core contract — Go generics
type Result[T any] struct { val T; err error }
func Ok[T any](v T) Result[T]
func Err[T any](e error) Result[T]
func (r Result[T]) Unwrap() (T, error)
func (r Result[T]) Map(fn func(T) T) Result[T]

// simple: Result[int] with Map/FlatMap/Or
// advanced: Result[T] pipeline — parse, validate, persist — short-circuits on first Err
```

**Retry**
```go
// Core contract
type BackoffStrategy func(attempt int) time.Duration
func ExponentialBackoff(base time.Duration, maxJitter time.Duration) BackoffStrategy
func Retry(ctx context.Context, op func() error, strategy BackoffStrategy, maxAttempts int) error

// simple: retry with fixed delay
// advanced: exponential backoff + jitter + circuit-breaker integration + Prometheus attempt counter
```

---

### 5.2 Distributed System Patterns

#### Resilience

**Circuit Breaker**
```go
// State machine: Closed → Open → Half-Open → Closed
type State int
const (Closed State = iota; Open; HalfOpen)

type Breaker struct { ... }
func New(opts ...Option) *Breaker
func (b *Breaker) Execute(ctx context.Context, fn func() error) error
// Emits: prometheus gauge (state), counter (trips, successes, failures)

// simple: threshold-based trip (N failures in window → Open)
// advanced: configurable per-state timeout, half-open probe count, event hooks,
//           integration with retry pattern
```

**Bulkhead**
```go
// Isolate caller groups — each gets its own semaphore + worker pool
type Bulkhead struct { pools map[string]*Pool }
func New(partitions map[string]Config) *Bulkhead
func (b *Bulkhead) Execute(ctx context.Context, partition string, fn func() error) error

// simple: two partitions — "critical" (10 concurrent) vs "batch" (2 concurrent)
// advanced: dynamic partition creation, overflow to dead-letter handler, metrics per partition
```

**Timeout**
```go
// Context-based timeout wrapping
func WithTimeout[T any](ctx context.Context, d time.Duration, fn func(context.Context) (T, error)) (T, error)

// simple: context.WithTimeout around an HTTP call
// advanced: cascading timeouts (total budget decomposed across steps), timeout telemetry
```

**Retry + Backoff**
```go
// Full production retry: exponential backoff, jitter, context-aware, retryable error predicate
type Config struct {
    MaxAttempts  int
    InitialDelay time.Duration
    MaxDelay     time.Duration
    Multiplier   float64
    Jitter       float64
    Retryable    func(error) bool  // caller decides what's retryable
}
func Do[T any](ctx context.Context, cfg Config, fn func() (T, error)) (T, error)

// simple: fixed-delay retry with max attempts
// advanced: full config above + circuit breaker wrapping + per-attempt span in OTel
```

**Hedged Requests**
```go
// Issue duplicate request after hedge delay; first response wins
func Hedge[T any](ctx context.Context, delay time.Duration, fn func(context.Context) (T, error), n int) (T, error)

// simple: two goroutines, first wins via select
// advanced: configurable n, latency percentile trigger, cancel all losers, metrics (hedge rate, win rate)
```

---

#### Communication

**gRPC Patterns**
```
proto/echo.proto  →  generated stubs
unary/            →  simple unary RPC
streaming/        →  server-streaming, client-streaming, bidirectional
interceptors/     →  unary + stream interceptor chain (logging, auth, metrics)

// advanced: OTel tracing interceptor, retry interceptor, deadline propagation
```

**Pub/Sub**
```go
// Core contract
type Topic[T any] struct { ... }
func NewTopic[T any](name string) *Topic[T]
func (t *Topic[T]) Publish(ctx context.Context, msg T) error
func (t *Topic[T]) Subscribe(ctx context.Context) (<-chan T, func())

// simple: in-process broker, 1 publisher / 2 subscribers
// advanced: async delivery with per-subscriber buffer, slow-consumer drop policy,
//           dead-letter queue, Prometheus lag metric
```

**Request-Reply**
```go
// Async request with correlation ID, reply via channel registry
type Broker struct { pending sync.Map }
func (b *Broker) Request(ctx context.Context, payload []byte) ([]byte, error)
func (b *Broker) HandleReply(correlationID string, reply []byte)

// simple: single-node in-process request-reply with timeout
// advanced: multi-node simulation — requests routed by topic, reply routing by correlationID
```

**Scatter-Gather**
```go
// Fan-out to N backends, collect all or first-N results within deadline
type Result[T any] struct { Source string; Value T; Err error; Latency time.Duration }
func ScatterGather[T any](ctx context.Context, fns map[string]func(context.Context) (T, error)) []Result[T]

// simple: 3 mock price feeds, gather all
// advanced: partial results on deadline (return what arrived), fastest-N mode, per-source circuit breaker
```

---

#### Data

**Saga**
```go
// Orchestration-based saga — coordinator drives steps with compensation
type Step struct {
    Name      string
    Execute   func(ctx context.Context, state *State) error
    Compensate func(ctx context.Context, state *State) error
}
type Saga struct { steps []Step }
func (s *Saga) Run(ctx context.Context) error  // runs forward; on failure, compensates in reverse

// simple: 3-step order saga (reserve stock → charge payment → dispatch)
// advanced: persistent saga state (SQLite), idempotency keys, step-level retry, OTel spans
```

**Outbox Pattern**
```go
// Transactionally write event to outbox table alongside business record
// Relay goroutine polls and publishes
type OutboxStore interface {
    SaveWithOutbox(ctx context.Context, record any, event OutboxEvent) error
    PendingEvents(ctx context.Context, limit int) ([]OutboxEvent, error)
    MarkPublished(ctx context.Context, id string) error
}

// simple: SQLite-backed outbox, polling relay, in-process publisher
// advanced: exactly-once delivery guarantee, relay with advisory lock, dead-letter handling
```

**Event Sourcing**
```go
// All state changes are events; current state is projection of event log
type Event interface { EventType() string; OccurredAt() time.Time }
type Store interface {
    Append(ctx context.Context, streamID string, events []Event, expectedVersion int) error
    Load(ctx context.Context, streamID string) ([]Event, error)
}
type Aggregate interface { Apply(Event); Version() int }

// simple: bank account — Deposited/Withdrawn events, balance projection
// advanced: snapshot every N events, optimistic concurrency via expectedVersion,
//           read-model projection (CQRS bridge)
```

**CQRS**
```go
// Command side: validates + mutates; Query side: reads from projection
type CommandBus struct { handlers map[string]CommandHandler }
type QueryBus  struct { handlers map[string]QueryHandler  }
type CommandHandler interface { Handle(ctx context.Context, cmd Command) error }
type QueryHandler  interface { Handle(ctx context.Context, q Query) (any, error) }

// simple: separate command/query structs, in-memory projection
// advanced: event-sourced command side, denormalised read model, async projection update,
//           versioned projections with rebuild support
```

---

#### Coordination

**Leader Election**
```go
// Lease-based election among N in-process nodes
type Node struct { id string; lease *Lease; peers []*Node }
type Lease struct { holder string; expiry time.Time; mu sync.Mutex }
func (n *Node) Campaign(ctx context.Context) (<-chan bool, error) // sends true on election win

// simple: 3 nodes, lease heartbeat, leader re-election on crash
// advanced: split-brain prevention, follower → candidate → leader state machine,
//           OTel event on every transition
```

**Distributed Lock**
```go
// Core contract — same interface, two implementations
type Locker interface {
    Lock(ctx context.Context, key string, ttl time.Duration) (token string, err error)
    Unlock(ctx context.Context, key string, token string) error
    Extend(ctx context.Context, key string, token string, ttl time.Duration) error
}
// simple: in-memory implementation with goroutine-safe map
// advanced: "Redis-like" implementation using a single goroutine as lock server,
//           fencing token for safe resource access, lock watchdog goroutine
```

**Two-Phase Commit**
```go
// Coordinator drives Prepare → Commit/Abort across participants
type Participant interface {
    Prepare(ctx context.Context, txID string) (Vote, error)  // Vote: Commit | Abort
    Commit(ctx context.Context, txID string) error
    Abort(ctx context.Context, txID string) error
}
type Coordinator struct { participants []Participant }
func (c *Coordinator) Execute(ctx context.Context, txID string) error

// simple: 2 participants, happy path + one abort
// advanced: coordinator crash recovery (WAL), participant timeout, heuristic completion
```

---

#### Observability

**Distributed Tracing**
```go
// OTel SDK — trace across simulated service boundary
func InitTracer(serviceName string) (trace.Tracer, func())
func InjectHTTP(ctx context.Context, req *http.Request)
func ExtractHTTP(req *http.Request) context.Context

// simple: single-service span with attributes and events
// advanced: two in-process services communicating via HTTP, trace propagated,
//           child spans per operation, error recording, stdout OTLP exporter
```

**Health Check**
```go
// Liveness + readiness separation
type Checker interface { Check(ctx context.Context) error }
type Server struct {
    live  []Checker  // liveness: is the process alive?
    ready []Checker  // readiness: can it serve traffic?
}
func (s *Server) Handler() http.Handler  // GET /healthz/live, /healthz/ready

// simple: always-healthy endpoints
// advanced: dependency checkers (DB ping, downstream HTTP), slow checker with timeout,
//           K8s-compatible JSON response body
```

**Structured Logging**
```go
// log/slog with context propagation (request-id, trace-id extracted from ctx)
type contextHandler struct { slog.Handler }
func (h contextHandler) Handle(ctx context.Context, r slog.Record) error

// simple: slog with JSON handler, field extraction
// advanced: context-enriched handler, sampling handler (drop DEBUG in production),
//           zap integration adapter
```

**Metrics**
```go
// RED method: Rate, Errors, Duration per operation
type Recorder struct {
    requestsTotal   *prometheus.CounterVec    // labels: method, status
    requestDuration *prometheus.HistogramVec  // labels: method
    inFlight        *prometheus.GaugeVec      // labels: method
}

// simple: manual counter + histogram for a mock handler
// advanced: middleware that auto-instruments http.Handler, custom buckets,
//           Grafana dashboard JSON (pre-built)
```

---

#### Scalability

**Consistent Hashing**
```go
// Virtual nodes ring — minimal remapping on node add/remove
type Ring struct { vnodes int; ring []uint32; nodes map[uint32]string; mu sync.RWMutex }
func New(vnodes int) *Ring
func (r *Ring) Add(node string)
func (r *Ring) Remove(node string)
func (r *Ring) Get(key string) string      // returns responsible node
func (r *Ring) GetN(key string, n int) []string  // n replicas

// simple: 3 nodes, show key distribution
// advanced: 150 vnodes, rebalancing simulation, key distribution histogram
```

**Sidecar**
```go
// Co-located process that enriches/intercepts traffic
// Modelled as in-process interceptor to avoid subprocess complexity
type Sidecar struct { port int; upstream string }
func (s *Sidecar) ServeHTTP(w http.ResponseWriter, r *http.Request)  // proxy + enrich

// simple: logging sidecar that prints request/response
// advanced: mTLS-terminating sidecar, header injection, retry on 503, metrics
```

**Service Discovery**
```go
// In-process registry with TTL-based deregistration
type Registry struct { services sync.Map }
type Instance struct { ID string; Addr string; Tags []string; TTL time.Duration }
func (r *Registry) Register(ctx context.Context, inst Instance) (deregister func(), err error)
func (r *Registry) Discover(serviceName string) ([]Instance, error)
func (r *Registry) Watch(ctx context.Context, serviceName string) <-chan []Instance

// simple: register + discover, no TTL
// advanced: TTL heartbeat, watch stream, round-robin client-side load balancing
```

**Load Balancer**
```go
// Strategy interface, multiple implementations
type Balancer interface { Next(backends []*Backend) *Backend }
type RoundRobin  struct { counter atomic.Uint64 }
type LeastConns  struct {}
type WeightedRR  struct {}
type Backend     struct { Addr string; Weight int; ActiveConns atomic.Int64; Healthy atomic.Bool }

// simple: round-robin across 3 backends
// advanced: least-connections + health-aware (skip unhealthy), circuit-breaker per backend,
//           Prometheus per-backend request counter
```

---

## 6. Code Quality Rules

| Rule | Enforcement |
|------|-------------|
| Format | `gofumpt -l .` → zero output |
| Vet | `go vet ./...` → zero output |
| Imports | stdlib → external → internal (goimports order) |
| Error wrapping | `fmt.Errorf("verb noun: %w", err)` — always lowercase, no period |
| Context first | Every function that does I/O takes `ctx context.Context` as first arg |
| No global state in patterns | Each `New()` returns self-contained instance |
| Test naming | `TestFunctionName_scenario_expectation` |
| No t.Fatal in goroutines | Use `t.Errorf` + channel or `sync.WaitGroup` |
| Goroutine leak prevention | Every goroutine started in `advanced/` has a documented shutdown path |

---

## 7. Implementation Order

Build in dependency order — earlier patterns are used in later ones:

| Phase | Patterns | Why this order |
|-------|----------|----------------|
| 1 | workspace bootstrap, go.mod files, master README | foundation |
| 2 | go-patterns: functional-options, builder, singleton, factory | no deps |
| 3 | go-patterns: decorator, adapter, proxy, middleware-chain | build on interfaces |
| 4 | go-patterns: pipeline, fan-out-fan-in, worker-pool, semaphore | channels |
| 5 | go-patterns: rate-limiter, context-propagation, errgroup | context |
| 6 | go-patterns: all error-handling, iterator, observer, strategy | complete language patterns |
| 7 | distributed: timeout, retry-backoff, circuit-breaker | resilience primitives |
| 8 | distributed: bulkhead, hedged-requests | resilience composites |
| 9 | distributed: structured-logging, metrics, health-check, distributed-tracing | observability |
| 10 | distributed: pub-sub, request-reply, scatter-gather, grpc-patterns | communication |
| 11 | distributed: outbox, event-sourcing, cqrs, saga | data patterns |
| 12 | distributed: distributed-lock, leader-election, two-phase-commit | coordination |
| 13 | distributed: consistent-hashing, service-discovery, load-balancer, sidecar | scalability |
| 14 | docs/rust-learning-prompt.md | capstone |

---

## 8. Rust Learning Prompt (`docs/rust-learning-prompt.md`)

The file is structured as a self-contained AI prompt. It contains:

### 8.1 Go → Rust Mental Model Table

| Go concept | Rust equivalent |
|------------|----------------|
| `interface{}` / `any` | `dyn Trait` / `Box<dyn Trait>` |
| goroutine | `tokio::spawn` / `async fn` |
| channel `chan T` | `tokio::sync::mpsc`, `broadcast`, `watch` |
| `sync.Mutex` | `std::sync::Mutex` / `tokio::sync::Mutex` |
| `context.Context` | `tokio_util::sync::CancellationToken` + `tokio::time::timeout` |
| `sync.Once` | `std::sync::OnceLock` |
| `sync.WaitGroup` | `tokio::task::JoinSet` |
| `errgroup.Group` | `tokio::task::JoinSet` with error collection |
| `interface` (implicit) | `trait` (explicit `impl`) |
| `defer` | `Drop` trait / `scopeguard` crate |
| nil | `Option<T>` |
| `error` interface | `std::error::Error` trait / `thiserror` / `anyhow` |
| `fmt.Errorf("%w")` | `anyhow::Context::context()` |
| `go:generate` | `build.rs` / `proc-macro` |
| `sync.Map` | `DashMap` crate |

### 8.2 Per-Pattern Prompt Template

Each pattern section contains:
```
## <Pattern Name>

**Go pattern:** go-patterns/<category>/<pattern>/
**Rust equivalent crates:** <list>

### Prompt
"Implement the <Pattern Name> pattern in Rust.

Requirements:
- Simple version: <specific constraints — no async, no extra crates>
- Advanced version: <specific constraints — async with tokio, production features>
- Include a README.md with the same sections as the Go version.
- Show how ownership and borrowing affect the design differently from Go.
- Highlight where Rust's type system provides stronger guarantees than Go's.

Specific implementation notes:
<2-3 Go-specific things to watch for in the Rust port>
"
```

### 8.3 Recommended Crate List by Category

```
Async runtime:        tokio (1.x)
Error handling:       thiserror, anyhow
Observability:        tracing, tracing-subscriber, metrics, prometheus
gRPC:                 tonic
Serialisation:        serde, serde_json
Concurrent maps:      dashmap
Lock-free:            crossbeam
HTTP server:          axum
HTTP client:          reqwest
Testing:              tokio::test, mockall, proptest
Config:               config, dotenvy
Database:             sqlx (async, compile-time checked)
```

---

## 9. Out of Scope

- Kubernetes operator code (separate repository).
- External infrastructure in `simple/` (no Docker required).
- UI, dashboards, or frontend code.
- Performance benchmarks beyond `_test.go` benchmark functions in `advanced/`.
- Production deployment manifests.
- Pattern variants beyond `simple/` + `advanced/`.
