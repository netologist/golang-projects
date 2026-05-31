# Rust Learning Prompt — A Go Developer's Guide to Rebuilding This Library in Rust

> Copy this file and use it as a prompt with any capable LLM (or work through it
> yourself) to build the Rust equivalent of the `golang-projects` reference
> library. It maps every Go pattern to idiomatic Rust, names the crates to use,
> and highlights the ownership/borrowing lessons each pattern teaches.

---

## How to Use This Prompt

For each pattern below, build:

1. **`simple/`** — minimal idiomatic Rust. No async, no extra crates unless noted.
2. **`advanced/`** — production-grade: async with `tokio`, real error handling, tests.
3. **`README.md`** — same sections as the Go version, plus a **Rust Notes** section
   explaining how ownership, borrowing, and lifetimes shaped the design differently
   than in Go.

Work through the patterns in the **Learning Sequence** at the end — it is ordered
so each pattern builds on concepts from the previous ones.

Use a **Cargo workspace** at the root:

```toml
# Cargo.toml
[workspace]
resolver = "2"
members = [
    "go-patterns/*/*/simple",
    "go-patterns/*/*/advanced",
    "distributed/*/*/simple",
    "distributed/*/*/advanced",
]
```

---

## Go → Rust Mental Model

| Go concept | Rust equivalent | Note |
|------------|----------------|------|
| `interface{}` / `any` | `Box<dyn Any>` / generics | Prefer generics; `dyn Any` only when you truly need type erasure |
| `interface` (implicit) | `trait` (explicit `impl Trait for T`) | Rust requires explicit `impl` |
| goroutine | `tokio::spawn` / `std::thread::spawn` | Tokio tasks for async, threads for blocking |
| `chan T` (unbuffered) | `tokio::sync::mpsc::channel(1)` | No truly-unbuffered async channel; cap 1 is closest |
| `chan T` (buffered) | `tokio::sync::mpsc::channel(n)` | |
| broadcast to many | `tokio::sync::broadcast` | Values must be `Clone` |
| single value handoff | `tokio::sync::oneshot` | Like a buffered `chan T` of cap 1, single use |
| `sync.Mutex` | `std::sync::Mutex<T>` / `tokio::sync::Mutex<T>` | Mutex *wraps the data* in Rust |
| `sync.RWMutex` | `std::sync::RwLock<T>` / `tokio::sync::RwLock<T>` | |
| `sync.Once` | `std::sync::OnceLock<T>` | |
| `sync.WaitGroup` | `tokio::task::JoinSet` / `JoinHandle` collection | |
| `errgroup.Group` | `tokio::task::JoinSet` + first-error logic | |
| `context.Context` (cancel) | `tokio_util::sync::CancellationToken` | |
| `context.Context` (deadline) | `tokio::time::timeout` / `timeout_at` | |
| `context.WithValue` | explicit struct fields or `tracing::Span` | Avoid `dyn Any` value bags |
| `defer` | `Drop` impl / `scopeguard::defer!` | RAII makes most `defer` unnecessary |
| `nil` | `Option<T>` | No null in safe Rust |
| `error` interface | `std::error::Error` / `thiserror` / `anyhow` | |
| `fmt.Errorf("...: %w", err)` | `anyhow::Context::context()` / `#[from]` | |
| `errors.Is` | `matches!` on an error enum / `downcast_ref` | |
| `errors.As` | `err.downcast_ref::<T>()` | |
| multiple return `(T, error)` | `Result<T, E>` | |
| `sync.Map` | `dashmap::DashMap<K, V>` | Lock-free concurrent map |
| `atomic.Int64` | `std::sync::atomic::AtomicI64` | |
| `select { }` | `tokio::select!` | |
| `time.After(d)` | `tokio::time::sleep(d)` | |
| `time.Ticker` | `tokio::time::interval` | |
| struct embedding | `Deref` impl or explicit delegation | No inheritance |
| `go:generate` | `build.rs` / proc-macros | |
| package `init()` | `inventory` / `ctor` crate, or explicit registration | No implicit init |

---

## Recommended Crates

```toml
[dependencies]
# Async runtime
tokio              = { version = "1", features = ["full"] }
tokio-util         = { version = "0.7", features = ["rt"] }

# Error handling
thiserror          = "1"
anyhow             = "1"

# Observability
tracing            = "0.1"
tracing-subscriber = { version = "0.3", features = ["env-filter", "json"] }
metrics            = "0.23"
metrics-exporter-prometheus = "0.15"
opentelemetry      = "0.24"
opentelemetry_sdk  = { version = "0.24", features = ["rt-tokio"] }

# gRPC
tonic              = "0.12"
prost              = "0.13"

# Serialisation
serde              = { version = "1", features = ["derive"] }
serde_json         = "1"

# Concurrent data structures
dashmap            = "6"
crossbeam          = "0.8"

# HTTP
axum               = "0.7"
reqwest            = { version = "0.12", features = ["json"] }
tower              = "0.5"
tower-http         = { version = "0.5", features = ["trace", "timeout"] }

# Resilience / utilities
governor           = "0.6"   # rate limiting (GCRA)
backon             = "1"      # retry with backoff policies
scopeguard         = "1"      # defer! macro

# Database (advanced data patterns)
sqlx               = { version = "0.8", features = ["sqlite", "runtime-tokio"] }

# Testing
mockall            = "0.13"
proptest           = "1"
tokio-test         = "0.4"

[build-dependencies]
tonic-build        = "0.12"
```

---

## Go Language Patterns

### Creational

#### 1. Functional Options — `go-patterns/creational/functional-options/`
**Rust crates:** none (simple); `derive_builder` (advanced).
```
Implement functional options in Rust two ways:
  (a) A builder struct with method chaining (the idiomatic Rust approach for config).
  (b) A Vec<Box<dyn FnOnce(&mut Config)>> approach (a direct translation of Go's options).
Then show the *typestate builder* that uses phantom types so that build() will not
compile unless all required fields are set. Include a test that demonstrates the
compile-time guarantee (use a doc-test with `compile_fail`).
Rust Notes: explain why Go's "apply functions to a pointer" is unusual in Rust and
how the builder owns its config until build() moves it out.
```

#### 2. Builder — `go-patterns/creational/builder/`
**Rust crates:** none.
```
Implement a SQL QueryBuilder.
  simple/: method chaining, validation at build() returning Result<Query, BuildError>.
  advanced/: a typestate builder where Table<Unset> vs Table<Set> makes calling
             build() without a table a compile error. Show a compile_fail doctest.
Rust Notes: contrast Go's runtime ErrMissingTable with Rust's compile-time enforcement.
```

#### 3. Singleton — `go-patterns/creational/singleton/`
**Rust crates:** `once_cell` or std `OnceLock`.
```
Implement a lazily-initialised singleton DB connection using std::sync::OnceLock.
  advanced/: a resettable-for-tests singleton using Mutex<Option<DB>>.
Rust Notes: explain why singletons are rarer in Rust (ownership discourages global
mutable state) and why OnceLock is preferred over a static mut. Discuss thread safety
guarantees vs Go's sync.Once.
```

#### 4. Factory — `go-patterns/creational/factory/`
**Rust crates:** `inventory` (advanced, for init()-style registration).
```
Implement a codec factory backed by a registry: HashMap<&'static str, fn() -> Box<dyn Codec>>.
  simple/: explicit registration in main().
  advanced/: use the `inventory` crate to register codecs at link time (Rust's
             closest analogue to Go's init() side-effect registration).
Rust Notes: Rust has no init(); explain inventory/ctor and why explicit wiring is
often preferred.
```

### Structural

#### 5. Decorator — `go-patterns/structural/decorator/`
**Rust crates:** none.
```
Define a `trait Store`. Implement two decorators:
  (a) LoggingStore<S: Store> — a generic (monomorphised, zero-cost) decorator.
  (b) A Box<dyn Store> dynamic-dispatch decorator.
Rust Notes: Go decorators always use dynamic dispatch via interfaces; Rust lets you
choose static (generics) or dynamic (dyn). Benchmark both and explain the trade-off.
```

#### 6. Adapter — `go-patterns/structural/adapter/`
**Rust crates:** `tokio`.
```
Adapt a blocking SyncReader into an async AsyncReader using tokio::task::spawn_blocking.
Rust Notes: this is the idiomatic way to bridge blocking code into async Rust —
contrast with Go where any function can block a goroutine cheaply.
```

#### 7. Proxy — `go-patterns/structural/proxy/`
**Rust crates:** `dashmap` (advanced).
```
Implement a CachingProxy<L: Loader> with a TTL.
  simple/: Arc<Mutex<HashMap>> cache.
  advanced/: dashmap::DashMap for lock-free reads + an AuthProxy that reads a token
             from a request struct (Rust has no context.Value; pass it explicitly).
Rust Notes: explain why Rust prefers explicit request structs over context value bags.
```

#### 8. Middleware Chain — `go-patterns/structural/middleware-chain/`
**Rust crates:** `tower`.
```
Implement a middleware chain:
  simple/: manual composition with boxed closures.
  advanced/: tower::ServiceBuilder with custom Layers for logging, request-id, and
             panic recovery (tower::ServiceExt). 
Rust Notes: tower::Service + tower::Layer is the canonical Rust middleware abstraction;
compare to Go's func(http.Handler) http.Handler.
```

### Behavioral

#### 9. Iterator — `go-patterns/behavioral/iterator/`
**Rust crates:** none.
```
Reimplement filter/map/take.
  simple/: use Rust's built-in Iterator adaptors directly.
  advanced/: a custom struct implementing the Iterator trait + a lazy paginated iterator.
Rust Notes: Rust iterators are lazy and zero-cost by default — closer to Go 1.23
range-over-func than to channel-based iteration. No goroutine needed.
```

#### 10. Observer — `go-patterns/behavioral/observer/`
**Rust crates:** `tokio`.
```
Implement an event bus.
  simple/: Mutex<Vec<Box<dyn Fn(&Event) + Send + Sync>>> — synchronous.
  advanced/: tokio::sync::broadcast for async multi-consumer delivery.
Rust Notes: broadcast requires Event: Clone. For expensive payloads, send Arc<Event>.
Contrast with Go channels that copy by value.
```

#### 11. Strategy — `go-patterns/behavioral/strategy/`
**Rust crates:** `flate2` (advanced, compression).
```
Implement strategy three ways: trait object (Box<dyn Strategy>), generic (S: Strategy),
and a plain fn pointer. For advanced/, select a compression strategy by payload size.
Rust Notes: a function pointer / closure is the most idiomatic single-method strategy,
just like Go's func value.
```

#### 12. Pipeline — `go-patterns/behavioral/pipeline/`
**Rust crates:** `tokio`.
```
  simple/: Iterator chaining (single-threaded, lazy).
  advanced/: tokio mpsc channels with concurrent stages and an error channel.
Rust Notes: Go pipelines are always concurrent (goroutines+channels); Rust iterator
pipelines are single-threaded. Use tokio when you actually need stage concurrency.
```

#### 13. Fan-out / Fan-in — `go-patterns/behavioral/fan-out-fan-in/`
**Rust crates:** `tokio`.
```
Fan out across tokio tasks reading one mpsc receiver (wrapped in Arc<Mutex>>), fan in
by merging with tokio::select! or a JoinSet.
Rust Notes: show how the borrow checker forces explicit shared ownership (Arc<Mutex>)
where Go silently shares the channel.
```

### Concurrency

#### 14. Worker Pool — `go-patterns/concurrency/worker-pool/`
**Rust crates:** `tokio` (I/O); `rayon` (CPU).
```
  simple/: rayon::ThreadPool for CPU-bound work (automatic work stealing).
  advanced/: a tokio task pool with a bounded mpsc queue, graceful shutdown via
             CancellationToken, and metrics.
Rust Notes: rayon for CPU-bound, tokio for I/O-bound. Explain when to use each.
```

#### 15. Semaphore — `go-patterns/concurrency/semaphore/`
**Rust crates:** `tokio`.
```
Use tokio::sync::Semaphore. Show acquire_owned() returning an OwnedSemaphorePermit
that releases on Drop (RAII).
Rust Notes: Rust's permit auto-releases on drop — you cannot forget to Release()
like you can with Go's chan struct{}.
```

#### 16. Rate Limiter — `go-patterns/concurrency/rate-limiter/`
**Rust crates:** `governor`.
```
  simple/: a manual token bucket with std::time + Mutex<f64>.
  advanced/: the governor crate (GCRA algorithm — more accurate than a token bucket).
Rust Notes: explain GCRA vs token bucket and why governor is production-grade.
```

#### 17. Context Propagation — `go-patterns/concurrency/context-propagation/`
**Rust crates:** `tokio-util`, `tracing`.
```
Show three replacements for context.Context:
  (a) CancellationToken for cancellation.
  (b) tokio::time::timeout for deadlines.
  (c) tracing::Span (or an explicit RequestContext struct) for request-scoped values.
Rust Notes: Rust has no single Context type. Explain why this is more explicit/safer
and how tracing::Span propagates request metadata through async call trees.
```

#### 18. errgroup — `go-patterns/concurrency/errgroup/`
**Rust crates:** `tokio`.
```
Implement a bounded task group with tokio::task::JoinSet + a tokio::sync::Semaphore.
Return the first error and abort the rest (drop the JoinSet).
Rust Notes: dropping a JoinSet aborts its tasks — contrast with errgroup's
context-cancellation model.
```

### Error Handling

#### 19. Sentinel Errors — `go-patterns/error-handling/sentinel-errors/`
**Rust crates:** `thiserror`.
```
Model sentinels as an enum: 
  #[derive(thiserror::Error, Debug)] enum StoreError { #[error("not found")] NotFound, ... }
Match with `matches!` or pattern matching; map variants to HTTP status codes.
Rust Notes: Rust enums make the set of errors exhaustive and checkable at compile time —
a strict improvement over Go's open set of sentinel values.
```

#### 20. Error Wrapping — `go-patterns/error-handling/error-wrapping/`
**Rust crates:** `thiserror`, `anyhow`.
```
Show both layers:
  - Library errors with thiserror + #[from] for automatic wrapping.
  - Application code with anyhow::Context for fmt.Errorf("%w")-style context.
Implement a custom AppError that serialises to JSON (serde) for API responses.
Rust Notes: map fmt.Errorf("%w") to #[from] and .context(); errors.As to downcast_ref.
```

#### 21. Result Type — `go-patterns/error-handling/result-type/`
**Rust crates:** none.
```
This pattern is *built into Rust* as Result<T, E>. Reimplement the Go Map/FlatMap/Or
chain using Result::map / and_then / unwrap_or, then show the `?` operator as the
idiomatic short-circuit.
Rust Notes: explain why Go needed a custom Result type and Rust does not.
```

#### 22. Retry — `go-patterns/error-handling/retry/`
**Rust crates:** `backon`, `tokio`.
```
  simple/: a manual async retry loop with exponential backoff + tokio::time::sleep.
  advanced/: the backon crate with ExponentialBuilder + jitter + a retryable predicate.
Rust Notes: contrast with the distributed/retry-backoff pattern; same idea, async.
```

---

## Distributed System Patterns

### Resilience

#### 23. Circuit Breaker — `distributed/resilience/circuit-breaker/`
**Rust crates:** `tokio`, `metrics`.
```
  simple/: a state machine guarded by std::sync::Mutex.
  advanced/: Arc<BreakerState> with an AtomicU8 state and compare_exchange for
             lock-free transitions; emit metrics via the `metrics` crate.
Rust Notes: multiple callers share the breaker — show Arc<...> and why the borrow
checker forces this explicit sharing.
```

#### 24. Bulkhead — `distributed/resilience/bulkhead/`
**Rust crates:** `tokio`, `tower`.
```
  simple/: one tokio::sync::Semaphore per partition.
  advanced/: tower::limit::ConcurrencyLimitLayer per partition.
Rust Notes: show how tower layers compose bulkheads with other middleware.
```

#### 25. Timeout — `distributed/resilience/timeout/`
**Rust crates:** `tokio`.
```
Wrap any future with tokio::time::timeout. advanced/: cascading deadlines with
tokio::time::Instant + timeout_at and a remaining-budget calculation.
Rust Notes: async cancellation in Rust is cooperative (futures stop being polled),
similar to honouring ctx.Done().
```

#### 26. Retry + Backoff — `distributed/resilience/retry-backoff/`
**Rust crates:** `backon`.
```
Use backon's ExponentialBuilder with jitter and a retryable predicate, returning a
typed Result<T, E>.
Rust Notes: same engine as the language-level retry, but used for distributed calls.
```

#### 27. Hedged Requests — `distributed/resilience/hedged-requests/`
**Rust crates:** `tokio`.
```
Spawn staggered attempts and race them with tokio::select!. The first success wins;
losers are aborted by dropping their JoinHandles / cancelling a shared token.
Rust Notes: dropping a future cancels it — the natural way to cancel losing hedges.
```

### Communication

#### 28. gRPC Patterns — `distributed/communication/grpc-patterns/`
**Rust crates:** `tonic`, `prost`, `tonic-build`.
```
Define echo.proto (unary + server-streaming + client-streaming). Generate with
tonic-build in build.rs. Implement the service and a logging interceptor
(tower::Layer / tonic Interceptor).
Rust Notes: tonic is built on tower + hyper; interceptors are tower layers.
```

#### 29. Pub/Sub — `distributed/communication/pub-sub/`
**Rust crates:** `tokio`.
```
  simple/: a single-topic broadcaster.
  advanced/: typed topics with tokio::sync::broadcast and serde for payloads.
Rust Notes: broadcast clones each message — wrap large payloads in Arc<T>.
```

#### 30. Request-Reply — `distributed/communication/request-reply/`
**Rust crates:** `tokio`, `dashmap`.
```
Correlate requests to replies with DashMap<String, oneshot::Sender<Bytes>>.
Rust Notes: tokio::sync::oneshot is the direct equivalent of Go's buffered chan(1).
```

#### 31. Scatter-Gather — `distributed/communication/scatter-gather/`
**Rust crates:** `tokio`.
```
Fan out with a JoinSet; gather results until a deadline (tokio::time::timeout on
JoinSet::join_next), returning partial results.
Rust Notes: JoinSet makes per-task lifecycle explicit; contrast with Go's errgroup.
```

### Data

#### 32. Saga — `distributed/data/saga/`
**Rust crates:** `tokio`, `async-trait` (or 2024-edition native async traits).
```
Orchestration saga with async steps. Each step is a struct implementing an async
trait with execute()/compensate(). Compensate completed steps in reverse on failure.
Rust Notes: discuss async closures vs async-trait vs the 2024 edition's native async
fn in traits; show the Box<dyn Future> workaround if needed.
```

#### 33. Outbox — `distributed/data/outbox/`
**Rust crates:** `sqlx` (advanced).
```
  simple/: an in-memory transactional outbox.
  advanced/: SQLx + SQLite, writing the record and the outbox row in one transaction,
             with a relay task that polls and publishes.
Rust Notes: SQLx checks SQL at compile time; show a transaction with sqlx::query!.
```

#### 34. Event Sourcing — `distributed/data/event-sourcing/`
**Rust crates:** `serde`.
```
Represent events as an enum (idiomatic) rather than Box<dyn Event>. Implement an
append-only store with an expected-version optimistic-concurrency check, plus
snapshotting.
Rust Notes: prefer an event enum + match over trait-object downcasting; explain why.
```

#### 35. CQRS — `distributed/data/cqrs/`
**Rust crates:** `async-trait`.
```
Command and query buses dispatching to handlers by name. Use an event-sourced write
side feeding a denormalised read projection.
Rust Notes: trait objects for handlers (Box<dyn CommandHandler>); discuss async traits.
```

### Coordination

#### 36. Leader Election — `distributed/coordination/leader-election/`
**Rust crates:** `tokio`.
```
Lease-based election among in-process nodes sharing an Arc<Mutex<Lease>>. Renew on a
tokio::time::interval; fail over when the lease expires.
Rust Notes: shared lease requires Arc<Mutex<...>>; the borrow checker makes the
shared-state contract explicit.
```

#### 37. Distributed Lock — `distributed/coordination/distributed-lock/`
**Rust crates:** `tokio`.
```
In-memory locker with TTL and a fencing token. advanced/: a RAII LockGuard that
releases on Drop.
Rust Notes: RAII means you cannot forget to unlock — a key safety win over Go's
manual Lock/Unlock.
```

#### 38. Two-Phase Commit — `distributed/coordination/two-phase-commit/`
**Rust crates:** `tokio`, `async-trait`.
```
Coordinator drives an async Participant trait (prepare/commit/abort). Any abort vote
or prepare error aborts all.
Rust Notes: model Vote as an enum; exhaustive matching guarantees both branches handled.
```

### Observability

#### 39. Structured Logging — `distributed/observability/structured-logging/`
**Rust crates:** `tracing`, `tracing-subscriber`.
```
Use tracing with a JSON subscriber. Propagate request-scoped fields with tracing
spans (#[instrument] and span.record).
Rust Notes: tracing spans are the idiomatic replacement for context-carried log fields.
```

#### 40. Health Check — `distributed/observability/health-check/`
**Rust crates:** `axum`, `tower-http`.
```
Expose /healthz/live and /healthz/ready with axum. Checks are async closures /
Box<dyn Check + Send + Sync>. Bound each with tower-http::timeout.
Rust Notes: axum handlers are async fns; show extractors and JSON responses.
```

#### 41. Metrics — `distributed/observability/metrics/`
**Rust crates:** `metrics`, `metrics-exporter-prometheus`.
```
Record RED metrics with the `metrics` crate macros (counter!, histogram!, gauge!) and
expose them via metrics-exporter-prometheus. Add an axum/tower layer that
auto-instruments requests.
Rust Notes: the `metrics` facade decouples instrumentation from the exporter.
```

#### 42. Distributed Tracing — `distributed/observability/distributed-tracing/`
**Rust crates:** `opentelemetry`, `opentelemetry_sdk`, `tracing-opentelemetry`.
```
Wire tracing + opentelemetry. Create a parent span and child spans for sub-operations;
record errors. Export to stdout (or an OTLP collector).
Rust Notes: tracing-opentelemetry bridges the tracing ecosystem to OTel.
```

### Scalability

#### 43. Consistent Hashing — `distributed/scalability/consistent-hashing/`
**Rust crates:** none (std `BTreeMap`).
```
Implement the ring with BTreeMap<u64, String>; use .range(hash..).next() (wrapping to
.next() on the whole map) for successor lookup — cleaner than Go's sorted-slice binary
search. Use 150 virtual nodes; include a distribution test.
Rust Notes: BTreeMap::range gives O(log n) successor lookup directly.
```

#### 44. Sidecar — `distributed/scalability/sidecar/`
**Rust crates:** `tower`, `tower-http`.
```
Model the sidecar as a tower layer wrapping an inner Service: inject headers, log,
and retry on 503.
Rust Notes: tower layers are the Rust way to express sidecar-style interception in-process.
```

#### 45. Service Discovery — `distributed/scalability/service-discovery/`
**Rust crates:** `tokio`, `dashmap`.
```
In-process registry: DashMap of instances with TTL, a reaper task, and Watch via
tokio::sync::watch::Receiver<Vec<Instance>>.
Rust Notes: tokio::sync::watch is the idiomatic equivalent of Go's Watch channel —
it always holds the latest value.
```

#### 46. Load Balancer — `distributed/scalability/load-balancer/`
**Rust crates:** `tower`.
```
Implement RoundRobin (AtomicUsize) and LeastConns (DashMap<String, AtomicI64>) as a
Balancer trait, then wrap them as a tower load-balance layer.
Rust Notes: tower::balance exists in the ecosystem; building your own teaches the
Service/Layer model.
```

---

## Learning Sequence

Build in this order — each phase reinforces the previous one:

1. **error-handling/** — `Result`, `Option`, `thiserror`, `anyhow`, `?` are the
   foundation of every later pattern.
2. **creational/** — builder + functional-options show how ownership shapes API design.
3. **structural/** — decorator + adapter teach the trait-object vs generics trade-off.
4. **behavioral/iterator** — internalise Rust's lazy, zero-cost iterator model.
5. **concurrency/** — make the goroutine → `tokio::spawn` mental shift; meet `Arc`,
   `Mutex`, channels, and `tokio::select!`.
6. **distributed/resilience/** — apply the concurrency primitives in realistic scenarios.
7. **distributed/observability/** — wire `tracing`, `metrics`, and OpenTelemetry.
8. **distributed/communication/** — tonic gRPC, broadcast pub/sub, oneshot request-reply.
9. **distributed/data/** — wrestle with async traits, trait objects, and SQLx.
10. **distributed/coordination/ & scalability/** — RAII locking, watch channels, tower layers.

---

## Key Rust Concepts Per Pattern

| Pattern | Key Rust concept to master |
|---------|----------------------------|
| Functional Options | Builder pattern, `impl Trait`, phantom-type typestate |
| Builder | Typestate, compile-time required-field enforcement |
| Singleton | `OnceLock`, why global mutable state is rare |
| Factory | `inventory`/`ctor`, link-time registration |
| Decorator | Trait objects vs generics, monomorphisation |
| Adapter | `spawn_blocking`, bridging sync↔async |
| Proxy | `Arc<Mutex>` vs `DashMap`, explicit request context |
| Middleware Chain | `tower::Service` / `tower::Layer` |
| Iterator | `Iterator` trait, laziness, zero-cost adaptors |
| Observer | `broadcast`, `Arc<T>` for cheap clones |
| Strategy | fn pointers, closures, generic vs `dyn` |
| Pipeline | iterator chains vs tokio channels |
| Fan-out/Fan-in | `JoinSet`, `Arc<Mutex<Receiver>>`, `select!` |
| Worker Pool | `rayon` vs `tokio`, `CancellationToken` |
| Semaphore | `OwnedSemaphorePermit`, RAII release |
| Rate Limiter | `governor`, GCRA |
| Context Propagation | `CancellationToken`, `timeout`, `tracing::Span` |
| errgroup | `JoinSet` abort-on-drop |
| Sentinel Errors | error enums, exhaustive matching |
| Error Wrapping | `thiserror` `#[from]`, `anyhow` context |
| Result Type | built-in `Result`, the `?` operator |
| Retry | `backon`, async backoff |
| Circuit Breaker | `AtomicU8` + `compare_exchange` |
| Bulkhead | `tower` concurrency limit layers |
| Timeout | `tokio::time::timeout`, cooperative cancellation |
| Hedged Requests | future cancellation on drop |
| gRPC | `tonic`, `tonic-build`, tower interceptors |
| Pub/Sub | `broadcast`, `Arc<T>` payloads |
| Request-Reply | `oneshot`, `DashMap` correlation |
| Scatter-Gather | `JoinSet` + deadline |
| Saga | async traits, `Box<dyn Future>`, `Pin` |
| Outbox | SQLx transactions, compile-checked SQL |
| Event Sourcing | event enums vs trait objects |
| CQRS | async-trait handlers |
| Leader Election | `Arc<Mutex<Lease>>`, intervals |
| Distributed Lock | RAII `LockGuard`, fencing tokens |
| Two-Phase Commit | vote enums, exhaustive matching |
| Structured Logging | `tracing`, `#[instrument]` |
| Health Check | `axum`, async checks |
| Metrics | `metrics` facade + Prometheus exporter |
| Distributed Tracing | `tracing-opentelemetry` |
| Consistent Hashing | `BTreeMap::range` successor lookup |
| Sidecar | `tower::Layer` interception |
| Service Discovery | `tokio::sync::watch`, `DashMap`, reaper task |
| Load Balancer | `Balancer` trait, `tower::balance` |

---

## Ownership Gotchas (Traps for Go Developers)

1. **One `&mut` at a time.** You cannot hold two mutable references to the same data.
   Plan your data-access patterns before writing the code.
2. **`Rc<T>` is single-threaded.** Use `Arc<T>` to share ownership across threads/tasks.
3. **Closures capture by reference by default.** Add `move` for closures passed to
   `tokio::spawn` or `std::thread::spawn`.
4. **Async functions are lazy.** A `Future` does nothing until `.await`ed (or spawned).
5. **Trait objects crossing threads need `Send + Sync`.** Expect `Box<dyn Trait + Send + Sync>`.
6. **`clone()` is explicit.** Go copies structs implicitly; in Rust you design for cloning
   (often via `Arc` to make clones cheap).
7. **The `?` operator uses `From`.** Design your error enums with `#[from]` so `?`
   converts automatically up the call stack.
8. **Drop order is deterministic (reverse declaration order).** Use `Drop` for cleanup
   instead of `defer`; this powers RAII guards (locks, permits, spans).
9. **`Mutex` wraps the data.** There is no separate "lock then touch the field" — you
   lock to *get* the data, and the guard releasing on drop unlocks it.
10. **No nil.** Model absence with `Option<T>` and the compiler forces you to handle `None`.
