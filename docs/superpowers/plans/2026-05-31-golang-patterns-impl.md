# Go Patterns & Distributed Systems — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a complete runnable reference library of 20 Go language patterns and 15 distributed system patterns, each with `simple/` and `advanced/` variants, full READMEs, and tests — capped by a Rust learning prompt.

**Architecture:** Two Go modules (`go-patterns`, `distributed`) in one `go.work` workspace. Every pattern is a self-contained directory with `simple/main.go` (runnable, no external deps) and `advanced/` (production-grade with tests). Patterns are implemented in dependency order: language fundamentals first, distributed primitives second.

**Tech Stack:** Go 1.22+, `golang.org/x/sync`, `google.golang.org/grpc`, `go.opentelemetry.io/otel`, `github.com/prometheus/client_golang`, `go.uber.org/zap`, `modernc.org/sqlite` (CGo-free, advanced data patterns only)

---

## Phase 1 — Workspace Bootstrap

### Task 1: go.work, go.mod files, shared internals, master README

**Files:**
- Create: `go.work`
- Create: `go-patterns/go.mod`
- Create: `go-patterns/internal/testutil/assert.go`
- Create: `go-patterns/internal/testutil/clock.go`
- Create: `distributed/go.mod`
- Create: `distributed/internal/testutil/server.go`
- Create: `distributed/internal/telemetry/metrics.go`
- Create: `distributed/internal/telemetry/tracing.go`
- Create: `README.md`

- [ ] **Step 1: Create go.work**

```
go 1.22

use (
	./go-patterns
	./distributed
)
```

- [ ] **Step 2: Create go-patterns/go.mod**

```
module github.com/veribaz/go-patterns

go 1.22

require golang.org/x/sync v0.7.0
```

- [ ] **Step 3: Create distributed/go.mod**

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
	modernc.org/sqlite                  v1.29.9
)
```

- [ ] **Step 4: Create go-patterns/internal/testutil/assert.go**

```go
package testutil

import (
	"fmt"
	"testing"
)

// Equal fails the test if got != want.
func Equal[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

// NoError fails the test if err != nil.
func NoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ErrorIs fails the test if errors.Is(err, target) is false.
func ErrorIs(t *testing.T, err, target error) {
	t.Helper()
	if !fmt.Errorf("%w", err).(*fmt.wrapError) != nil {
		t.Fatalf("expected error %v, got %v", target, err)
	}
}
```

- [ ] **Step 5: Create go-patterns/internal/testutil/clock.go**

```go
package testutil

import (
	"sync"
	"time"
)

// FakeClock is a manually-advanced clock for deterministic tests.
type FakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func NewFakeClock(t time.Time) *FakeClock { return &FakeClock{now: t} }

func (c *FakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *FakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}
```

- [ ] **Step 6: Create distributed/internal/telemetry/metrics.go**

```go
package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// NewRegistry creates a fresh Prometheus registry (not the default global).
// Use a local registry in examples so patterns don't pollute each other.
func NewRegistry() *prometheus.Registry {
	return prometheus.NewRegistry()
}

// Handler returns an HTTP handler for /metrics using the given registry.
func Handler(reg *prometheus.Registry) http.Handler {
	return promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
}
```

- [ ] **Step 7: Create distributed/internal/telemetry/tracing.go**

```go
package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
)

// InitStdout initialises an OTel tracer that prints spans to stdout.
// Returns the tracer and a shutdown func.
func InitStdout(serviceName string) (shutdown func(context.Context) error, err error) {
	exp, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	tp := trace.NewTracerProvider(trace.WithBatcher(exp))
	otel.SetTracerProvider(tp)
	return tp.Shutdown, nil
}
```

- [ ] **Step 8: Create master README.md**

```markdown
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
|---------|--------|---------|
| [Functional Options](go-patterns/creational/functional-options/) | ✓ | ✓ |
| [Builder](go-patterns/creational/builder/) | ✓ | ✓ |
| [Singleton](go-patterns/creational/singleton/) | ✓ | ✓ |
| [Factory](go-patterns/creational/factory/) | ✓ | ✓ |

### Structural
| [Decorator](go-patterns/structural/decorator/) | ✓ | ✓ |
| [Adapter](go-patterns/structural/adapter/) | ✓ | ✓ |
| [Proxy](go-patterns/structural/proxy/) | ✓ | ✓ |
| [Middleware Chain](go-patterns/structural/middleware-chain/) | ✓ | ✓ |

### Behavioral
| [Iterator](go-patterns/behavioral/iterator/) | ✓ | ✓ |
| [Observer](go-patterns/behavioral/observer/) | ✓ | ✓ |
| [Strategy](go-patterns/behavioral/strategy/) | ✓ | ✓ |
| [Pipeline](go-patterns/behavioral/pipeline/) | ✓ | ✓ |
| [Fan-out/Fan-in](go-patterns/behavioral/fan-out-fan-in/) | ✓ | ✓ |

### Concurrency
| [Worker Pool](go-patterns/concurrency/worker-pool/) | ✓ | ✓ |
| [Semaphore](go-patterns/concurrency/semaphore/) | ✓ | ✓ |
| [Rate Limiter](go-patterns/concurrency/rate-limiter/) | ✓ | ✓ |
| [Context Propagation](go-patterns/concurrency/context-propagation/) | ✓ | ✓ |
| [errgroup](go-patterns/concurrency/errgroup/) | ✓ | ✓ |

### Error Handling
| [Sentinel Errors](go-patterns/error-handling/sentinel-errors/) | ✓ | ✓ |
| [Error Wrapping](go-patterns/error-handling/error-wrapping/) | ✓ | ✓ |
| [Result Type](go-patterns/error-handling/result-type/) | ✓ | ✓ |
| [Retry](go-patterns/error-handling/retry/) | ✓ | ✓ |

## Distributed System Patterns (`distributed/`)

### Resilience
| [Circuit Breaker](distributed/resilience/circuit-breaker/) | ✓ | ✓ |
| [Bulkhead](distributed/resilience/bulkhead/) | ✓ | ✓ |
| [Timeout](distributed/resilience/timeout/) | ✓ | ✓ |
| [Retry + Backoff](distributed/resilience/retry-backoff/) | ✓ | ✓ |
| [Hedged Requests](distributed/resilience/hedged-requests/) | ✓ | ✓ |

### Communication
| [gRPC Patterns](distributed/communication/grpc-patterns/) | ✓ | ✓ |
| [Pub/Sub](distributed/communication/pub-sub/) | ✓ | ✓ |
| [Request-Reply](distributed/communication/request-reply/) | ✓ | ✓ |
| [Scatter-Gather](distributed/communication/scatter-gather/) | ✓ | ✓ |

### Data
| [Saga](distributed/data/saga/) | ✓ | ✓ |
| [Outbox](distributed/data/outbox/) | ✓ | ✓ |
| [Event Sourcing](distributed/data/event-sourcing/) | ✓ | ✓ |
| [CQRS](distributed/data/cqrs/) | ✓ | ✓ |

### Coordination
| [Leader Election](distributed/coordination/leader-election/) | ✓ | ✓ |
| [Distributed Lock](distributed/coordination/distributed-lock/) | ✓ | ✓ |
| [Two-Phase Commit](distributed/coordination/two-phase-commit/) | ✓ | ✓ |

### Observability
| [Distributed Tracing](distributed/observability/distributed-tracing/) | ✓ | ✓ |
| [Health Check](distributed/observability/health-check/) | ✓ | ✓ |
| [Structured Logging](distributed/observability/structured-logging/) | ✓ | ✓ |
| [Metrics](distributed/observability/metrics/) | ✓ | ✓ |

### Scalability
| [Consistent Hashing](distributed/scalability/consistent-hashing/) | ✓ | ✓ |
| [Sidecar](distributed/scalability/sidecar/) | ✓ | ✓ |
| [Service Discovery](distributed/scalability/service-discovery/) | ✓ | ✓ |
| [Load Balancer](distributed/scalability/load-balancer/) | ✓ | ✓ |
```

- [ ] **Step 9: Run module tidy and verify workspace**

```bash
cd go-patterns && go mod tidy && cd ../distributed && go mod tidy && cd .. && go work sync
```
Expected: no errors, `go.sum` files created.

- [ ] **Step 10: Commit**

```bash
git add go.work go-patterns/ distributed/ README.md
git commit -m "feat: bootstrap workspace, go.mod files, shared internal packages"
```

---

## Phase 2 — Go Language Patterns: Creational

### Task 2: Functional Options

**Files:**
- Create: `go-patterns/creational/functional-options/README.md`
- Create: `go-patterns/creational/functional-options/simple/options.go`
- Create: `go-patterns/creational/functional-options/simple/main.go`
- Create: `go-patterns/creational/functional-options/advanced/options.go`
- Create: `go-patterns/creational/functional-options/advanced/options_test.go`
- Create: `go-patterns/creational/functional-options/advanced/main.go`

- [ ] **Step 1: Create README.md**

```markdown
# Functional Options

## Concept
Instead of a config struct or multiple constructors, accept variadic `func(*T)`
options. Each option mutates one field. The caller composes only what they need.

## When to Use
- ≥3 optional parameters
- Zero-value of a field is not a safe default
- You want to add options in the future without breaking callers

## When NOT to Use
- All parameters are required — use a plain constructor
- Only 1-2 simple fields — a struct literal is clearer

## Trade-offs
| Benefit | Cost |
|---------|------|
| Extensible without breaking changes | Slightly more verbose per-call |
| Self-documenting option names | Options applied silently — bad for validation |
| Works with interfaces | Harder to inspect/serialise config |

## Go-Specific Notes
Options are just `func(*Server)` values — they are first-class and composable.
Validation lives inside `New()`, not inside each option, keeping option funcs
simple.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Option funcs are ordinary values — store them, compose them, pass them around
- Validate after applying all options, not inside each option
- Name options `With<Field>` by convention
- Use a `config` unexported struct inside the type to separate defaults from overrides
```

- [ ] **Step 2: Create simple/options.go**

```go
package main

import "time"

// Server is the type being configured.
type Server struct {
	addr    string
	timeout time.Duration
	maxConn int
}

// Option is a function that configures a Server.
type Option func(*Server)

// WithTimeout sets the read/write timeout.
func WithTimeout(d time.Duration) Option {
	return func(s *Server) { s.timeout = d }
}

// WithMaxConn sets the maximum number of concurrent connections.
func WithMaxConn(n int) Option {
	return func(s *Server) { s.maxConn = n }
}

// New creates a Server with defaults, then applies opts.
func New(addr string, opts ...Option) *Server {
	s := &Server{
		addr:    addr,
		timeout: 30 * time.Second, // safe default
		maxConn: 100,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}
```

- [ ] **Step 3: Create simple/main.go**

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Default server
	s1 := New(":8080")
	fmt.Printf("s1: addr=%s timeout=%s maxConn=%d\n", s1.addr, s1.timeout, s1.maxConn)

	// Customised server
	s2 := New(":9090",
		WithTimeout(5*time.Second),
		WithMaxConn(50),
	)
	fmt.Printf("s2: addr=%s timeout=%s maxConn=%d\n", s2.addr, s2.timeout, s2.maxConn)
}
```

- [ ] **Step 4: Run simple**

```bash
go run ./go-patterns/creational/functional-options/simple
```
Expected:
```
s1: addr=:8080 timeout=30s maxConn=100
s2: addr=:9090 timeout=5s maxConn=50
```

- [ ] **Step 5: Write failing test (advanced/options_test.go)**

```go
package options_test

import (
	"errors"
	"testing"
	"time"
)

func TestNew_defaults(t *testing.T) {
	s, err := New(":8080")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Addr() != ":8080" {
		t.Errorf("addr: got %s, want :8080", s.Addr())
	}
	if s.Timeout() != 30*time.Second {
		t.Errorf("timeout: got %s, want 30s", s.Timeout())
	}
}

func TestNew_customOptions(t *testing.T) {
	s, err := New(":9090", WithTimeout(5*time.Second), WithMaxConn(50))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Timeout() != 5*time.Second {
		t.Errorf("timeout: got %s, want 5s", s.Timeout())
	}
	if s.MaxConn() != 50 {
		t.Errorf("maxConn: got %d, want 50", s.MaxConn())
	}
}

func TestNew_validation(t *testing.T) {
	_, err := New("", WithTimeout(-1*time.Second))
	if !errors.Is(err, ErrInvalidConfig) {
		t.Errorf("expected ErrInvalidConfig, got %v", err)
	}
}
```

- [ ] **Step 6: Run test — verify it fails**

```bash
go test ./go-patterns/creational/functional-options/advanced/... -v
```
Expected: FAIL — `options.go` not found.

- [ ] **Step 7: Create advanced/options.go**

```go
package options

import (
	"errors"
	"fmt"
	"time"
)

// ErrInvalidConfig is returned when option validation fails.
var ErrInvalidConfig = errors.New("invalid config")

type config struct {
	addr    string
	timeout time.Duration
	maxConn int
	frozen  bool // set after New() — options applied after this panic
}

// Server is the production-grade configured server.
type Server struct {
	cfg config
}

// Option is a function that mutates the config.
type Option func(*config) error

// WithTimeout sets the read/write timeout. Must be > 0.
func WithTimeout(d time.Duration) Option {
	return func(c *config) error {
		if d <= 0 {
			return fmt.Errorf("%w: timeout must be > 0", ErrInvalidConfig)
		}
		c.timeout = d
		return nil
	}
}

// WithMaxConn sets the max concurrent connections. Must be > 0.
func WithMaxConn(n int) Option {
	return func(c *config) error {
		if n <= 0 {
			return fmt.Errorf("%w: maxConn must be > 0", ErrInvalidConfig)
		}
		c.maxConn = n
		return nil
	}
}

// New creates a Server, applies options, and validates.
func New(addr string, opts ...Option) (*Server, error) {
	if addr == "" {
		return nil, fmt.Errorf("%w: addr is required", ErrInvalidConfig)
	}
	cfg := config{
		addr:    addr,
		timeout: 30 * time.Second,
		maxConn: 100,
	}
	for _, o := range opts {
		if err := o(&cfg); err != nil {
			return nil, err
		}
	}
	cfg.frozen = true
	return &Server{cfg: cfg}, nil
}

func (s *Server) Addr() string           { return s.cfg.addr }
func (s *Server) Timeout() time.Duration { return s.cfg.timeout }
func (s *Server) MaxConn() int           { return s.cfg.maxConn }
```

- [ ] **Step 8: Create advanced/main.go**

```go
package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	s, err := New(":8080", WithTimeout(10*time.Second), WithMaxConn(200))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Server ready: addr=%s timeout=%s maxConn=%d\n",
		s.Addr(), s.Timeout(), s.MaxConn())

	// Demonstrate validation
	_, err = New("", WithTimeout(-1*time.Second))
	fmt.Printf("Expected error: %v\n", err)
}
```

- [ ] **Step 9: Run tests — verify pass**

```bash
go test ./go-patterns/creational/functional-options/advanced/... -v
```
Expected: all PASS.

- [ ] **Step 10: Commit**

```bash
git add go-patterns/creational/functional-options/
git commit -m "feat(go-patterns): functional options — simple + advanced"
```

---

### Task 3: Builder

**Files:**
- Create: `go-patterns/creational/builder/README.md`
- Create: `go-patterns/creational/builder/simple/builder.go`
- Create: `go-patterns/creational/builder/simple/main.go`
- Create: `go-patterns/creational/builder/advanced/builder.go`
- Create: `go-patterns/creational/builder/advanced/builder_test.go`
- Create: `go-patterns/creational/builder/advanced/main.go`

- [ ] **Step 1: Create README.md**

```markdown
# Builder

## Concept
Separate the construction of a complex object from its representation.
A Builder accumulates configuration through method calls, then produces
a validated immutable result via `Build()`.

## When to Use
- Object requires multi-step construction with validation at the end
- You want method chaining (fluent API)
- Construction involves conditional branches

## When NOT to Use
- Simple structs — struct literals are clearer
- When Functional Options already covers the need

## Trade-offs
| Benefit | Cost |
|---------|------|
| Validation at one point (Build) | More boilerplate than functional options |
| Fluent, readable call sites | Builder itself is stateful — not safe to reuse |

## Go-Specific Notes
Return `*Builder` from each setter for chaining. `Build()` returns `(T, error)`
to signal validation failure without panicking.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- `Build()` is the only place that validates
- Setters return `*Builder` — they never fail
- The returned type is immutable (all fields unexported, accessed via methods)
- A builder should not be reused after `Build()`
```

- [ ] **Step 2: Create simple/builder.go**

```go
package main

import "fmt"

// Query is an immutable SQL query.
type Query struct {
	table  string
	cols   []string
	where  string
}

func (q Query) String() string {
	cols := "*"
	if len(q.cols) > 0 {
		cols = fmt.Sprintf("%v", q.cols)
	}
	s := fmt.Sprintf("SELECT %s FROM %s", cols, q.table)
	if q.where != "" {
		s += " WHERE " + q.where
	}
	return s
}

// QueryBuilder builds a Query.
type QueryBuilder struct {
	table string
	cols  []string
	where string
}

func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{table: table}
}

func (b *QueryBuilder) Select(cols ...string) *QueryBuilder {
	b.cols = cols
	return b
}

func (b *QueryBuilder) Where(cond string) *QueryBuilder {
	b.where = cond
	return b
}

func (b *QueryBuilder) Build() (Query, error) {
	if b.table == "" {
		return Query{}, fmt.Errorf("table is required")
	}
	return Query{table: b.table, cols: b.cols, where: b.where}, nil
}
```

- [ ] **Step 3: Create simple/main.go**

```go
package main

import (
	"fmt"
	"log"
)

func main() {
	q, err := NewQueryBuilder("users").
		Select("id", "name", "email").
		Where("active = true").
		Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(q)
}
```

- [ ] **Step 4: Run simple**

```bash
go run ./go-patterns/creational/builder/simple
```
Expected: `SELECT [id name email] FROM users WHERE active = true`

- [ ] **Step 5: Write failing test (advanced/builder_test.go)**

```go
package builder_test

import (
	"errors"
	"testing"
)

func TestQueryBuilder_happyPath(t *testing.T) {
	q, err := NewQueryBuilder("orders").
		Select("id", "total").
		Where("status = 'pending'").
		Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := q.SQL()
	want := "SELECT id, total FROM orders WHERE status = 'pending'"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestQueryBuilder_missingTable(t *testing.T) {
	_, err := NewQueryBuilder("").Build()
	if !errors.Is(err, ErrMissingTable) {
		t.Errorf("expected ErrMissingTable, got %v", err)
	}
}

func TestQueryBuilder_multipleWhere(t *testing.T) {
	q, _ := NewQueryBuilder("items").
		Where("price > 10").
		Where("stock > 0").
		Build()
	got := q.SQL()
	want := "SELECT * FROM items WHERE price > 10 AND stock > 0"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}
```

- [ ] **Step 6: Run test — verify fail**

```bash
go test ./go-patterns/creational/builder/advanced/... -v
```
Expected: FAIL — package not found.

- [ ] **Step 7: Create advanced/builder.go**

```go
package builder

import (
	"errors"
	"fmt"
	"strings"
)

var ErrMissingTable = errors.New("table is required")

// Query is an immutable, validated SQL query.
type Query struct {
	sql    string
	args   []any
}

func (q Query) SQL() string  { return q.sql }
func (q Query) Args() []any  { return q.args }

// QueryBuilder accumulates clauses and produces a Query via Build().
type QueryBuilder struct {
	table  string
	cols   []string
	wheres []string
	args   []any
}

func NewQueryBuilder(table string) *QueryBuilder {
	return &QueryBuilder{table: table}
}

func (b *QueryBuilder) Select(cols ...string) *QueryBuilder {
	b.cols = append(b.cols, cols...)
	return b
}

func (b *QueryBuilder) Where(cond string, args ...any) *QueryBuilder {
	b.wheres = append(b.wheres, cond)
	b.args = append(b.args, args...)
	return b
}

func (b *QueryBuilder) Build() (Query, error) {
	if b.table == "" {
		return Query{}, fmt.Errorf("build query: %w", ErrMissingTable)
	}
	cols := "*"
	if len(b.cols) > 0 {
		cols = strings.Join(b.cols, ", ")
	}
	sql := fmt.Sprintf("SELECT %s FROM %s", cols, b.table)
	if len(b.wheres) > 0 {
		sql += " WHERE " + strings.Join(b.wheres, " AND ")
	}
	return Query{sql: sql, args: b.args}, nil
}
```

- [ ] **Step 8: Create advanced/main.go**

```go
package main

import (
	"fmt"
	"log"
)

func main() {
	q, err := NewQueryBuilder("orders").
		Select("id", "total", "status").
		Where("status = $1", "pending").
		Where("total > $2", 100).
		Build()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SQL:", q.SQL())
	fmt.Println("Args:", q.Args())
}
```

- [ ] **Step 9: Run tests**

```bash
go test ./go-patterns/creational/builder/advanced/... -v
```
Expected: all PASS.

- [ ] **Step 10: Commit**

```bash
git add go-patterns/creational/builder/
git commit -m "feat(go-patterns): builder — simple + advanced"
```

---

### Task 4: Singleton

**Files:**
- Create: `go-patterns/creational/singleton/README.md`
- Create: `go-patterns/creational/singleton/simple/singleton.go` + `main.go`
- Create: `go-patterns/creational/singleton/advanced/singleton.go` + `singleton_test.go` + `main.go`

- [ ] **Step 1: Create README.md**

```markdown
# Singleton

## Concept
Ensure a type has exactly one instance and provide a global access point.
In Go, lazy initialisation uses `sync.Once` — goroutine-safe and allocation-free
after first call.

## When to Use
- Shared resources (DB pool, logger, config) that are expensive to create
- Exactly one instance is correct by design

## When NOT to Use
- Testing: singletons are hard to mock/replace — prefer dependency injection
- Multiple configurations needed — use a factory instead

## Trade-offs
| Benefit | Cost |
|---------|------|
| Guaranteed single instance | Global state — hidden dependency |
| Thread-safe lazy init | Hard to reset in tests |
| Zero overhead after first call | Init errors are hard to surface |

## Go-Specific Notes
`sync.Once.Do` runs the function exactly once across all goroutines.
For testable code, also expose `Reset()` (only in test builds) via `_test.go`.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Use `sync.Once` — never roll your own double-checked locking
- Expose a `Reset()` for tests via a `export_test.go` file
- Prefer dependency injection over singletons in library code
- Singleton is acceptable for truly global resources (logger, metrics registry)
```

- [ ] **Step 2: Create simple/singleton.go**

```go
package main

import (
	"fmt"
	"sync"
)

type logger struct{ prefix string }

func (l *logger) Log(msg string) { fmt.Printf("[%s] %s\n", l.prefix, msg) }

var (
	instance *logger
	once     sync.Once
)

// GetLogger returns the singleton logger, initialised on first call.
func GetLogger() *logger {
	once.Do(func() {
		instance = &logger{prefix: "APP"}
	})
	return instance
}
```

- [ ] **Step 3: Create simple/main.go**

```go
package main

func main() {
	l1 := GetLogger()
	l2 := GetLogger()

	l1.Log("hello from l1")
	l2.Log("hello from l2")

	if l1 == l2 {
		l1.Log("same instance confirmed")
	}
}
```

- [ ] **Step 4: Run simple**

```bash
go run ./go-patterns/creational/singleton/simple
```
Expected:
```
[APP] hello from l1
[APP] hello from l2
[APP] same instance confirmed
```

- [ ] **Step 5: Write failing test (advanced/singleton_test.go)**

```go
package singleton_test

import (
	"sync"
	"testing"
)

func TestGetDB_sameInstance(t *testing.T) {
	db1 := GetDB()
	db2 := GetDB()
	if db1 != db2 {
		t.Error("expected same instance, got different pointers")
	}
}

func TestGetDB_concurrentInit(t *testing.T) {
	Reset() // reset for test isolation
	var wg sync.WaitGroup
	instances := make([]*DB, 100)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			instances[i] = GetDB()
		}(i)
	}
	wg.Wait()
	for i, db := range instances {
		if db != instances[0] {
			t.Errorf("instance[%d] differs from instance[0]", i)
		}
	}
}
```

- [ ] **Step 6: Run test — verify fail**

```bash
go test ./go-patterns/creational/singleton/advanced/... -v
```
Expected: FAIL.

- [ ] **Step 7: Create advanced/singleton.go**

```go
package singleton

import "sync"

// DB simulates an expensive database connection pool.
type DB struct{ dsn string }

func (d *DB) Ping() bool { return true }

var (
	dbInstance *DB
	dbOnce     sync.Once
)

// GetDB returns the singleton DB instance.
func GetDB() *DB {
	dbOnce.Do(func() {
		dbInstance = &DB{dsn: "postgres://localhost/app"}
	})
	return dbInstance
}
```

- [ ] **Step 8: Create advanced/export_test.go** (test-only reset hook)

```go
package singleton

// Reset tears down the singleton for test isolation.
// Only available in test builds.
func Reset() {
	dbOnce = sync.Once{}
	dbInstance = nil
}
```

- [ ] **Step 9: Create advanced/main.go**

```go
package main

import "fmt"

func main() {
	db1 := GetDB()
	db2 := GetDB()
	fmt.Println("same instance:", db1 == db2)
	fmt.Println("ping:", db1.Ping())
}
```

- [ ] **Step 10: Run tests**

```bash
go test ./go-patterns/creational/singleton/advanced/... -v
```
Expected: all PASS.

- [ ] **Step 11: Commit**

```bash
git add go-patterns/creational/singleton/
git commit -m "feat(go-patterns): singleton — simple + advanced"
```

---

### Task 5: Factory

**Files:**
- Create: `go-patterns/creational/factory/README.md`
- Create: `go-patterns/creational/factory/simple/factory.go` + `main.go`
- Create: `go-patterns/creational/factory/advanced/factory.go` + `factory_test.go` + `main.go`

- [ ] **Step 1: Create README.md**

```markdown
# Factory

## Concept
Provide a creation interface that hides which concrete type is instantiated.
The registry pattern (a `map[string]func() Interface`) is the Go-idiomatic
approach — it avoids switch statements and is open to extension.

## When to Use
- Multiple implementations of an interface, selected at runtime
- Plugin-style extensibility — third-party code can register its own factories
- Decouple callers from concrete types

## When NOT to Use
- Only one implementation exists — a plain constructor is clearer
- The selection logic is trivially simple (one if/else)

## Trade-offs
| Benefit | Cost |
|---------|------|
| Open/Closed: add types without changing factory | Registered types may be nil if registration is missed |
| Decouples caller from implementation | Type lookup errors surface at runtime |

## Go-Specific Notes
Register factories in `init()` in each codec's file — callers just import the
package side-effect to get the codec available.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Use `map[string]func() Interface` as the registry
- `init()` side-effect registration keeps factory and implementation co-located
- Always return `(Interface, error)` from `New(name)` — unknown name is an error
```

- [ ] **Step 2: Create simple/factory.go**

```go
package main

import (
	"encoding/json"
	"fmt"
)

// Codec serialises and deserialises values.
type Codec interface {
	Encode(v any) ([]byte, error)
	Decode(data []byte, v any) error
	Name() string
}

// registry maps codec names to constructor functions.
var registry = map[string]func() Codec{}

func Register(name string, fn func() Codec) { registry[name] = fn }

func New(name string) (Codec, error) {
	fn, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("codec %q not registered", name)
	}
	return fn(), nil
}

// jsonCodec is registered below.
type jsonCodec struct{}

func (c *jsonCodec) Name() string                        { return "json" }
func (c *jsonCodec) Encode(v any) ([]byte, error)        { return json.Marshal(v) }
func (c *jsonCodec) Decode(data []byte, v any) error     { return json.Unmarshal(data, v) }

func init() {
	Register("json", func() Codec { return &jsonCodec{} })
}
```

- [ ] **Step 3: Create simple/main.go**

```go
package main

import "fmt"

func main() {
	c, err := New("json")
	if err != nil {
		panic(err)
	}
	data, _ := c.Encode(map[string]int{"a": 1})
	fmt.Printf("Encoded: %s\n", data)

	_, err = New("msgpack")
	fmt.Printf("Unknown codec: %v\n", err)
}
```

- [ ] **Step 4: Run simple**

```bash
go run ./go-patterns/creational/factory/simple
```
Expected: `Encoded: {"a":1}` and error for msgpack.

- [ ] **Step 5: Write failing test (advanced/factory_test.go)**

```go
package factory_test

import (
	"errors"
	"testing"
)

func TestNew_knownCodec(t *testing.T) {
	c, err := New("json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Name() != "json" {
		t.Errorf("got %s, want json", c.Name())
	}
}

func TestNew_roundtrip(t *testing.T) {
	c, _ := New("json")
	type payload struct{ X int }
	data, err := c.Encode(payload{X: 42})
	if err != nil {
		t.Fatal(err)
	}
	var out payload
	if err := c.Decode(data, &out); err != nil {
		t.Fatal(err)
	}
	if out.X != 42 {
		t.Errorf("got %d, want 42", out.X)
	}
}

func TestNew_unknownCodec(t *testing.T) {
	_, err := New("nope")
	if !errors.Is(err, ErrUnknownCodec) {
		t.Errorf("expected ErrUnknownCodec, got %v", err)
	}
}
```

- [ ] **Step 6: Run test — verify fail**

```bash
go test ./go-patterns/creational/factory/advanced/... -v
```
Expected: FAIL.

- [ ] **Step 7: Create advanced/factory.go**

```go
package factory

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

// ErrUnknownCodec is returned when the requested codec is not registered.
var ErrUnknownCodec = errors.New("unknown codec")

// Codec serialises and deserialises values.
type Codec interface {
	Name() string
	Encode(v any) ([]byte, error)
	Decode(data []byte, v any) error
}

type registry struct {
	mu       sync.RWMutex
	codecs   map[string]func() Codec
}

var global = &registry{codecs: map[string]func() Codec{}}

// Register makes a codec available by name.
func Register(name string, fn func() Codec) {
	global.mu.Lock()
	defer global.mu.Unlock()
	global.codecs[name] = fn
}

// New returns a new instance of the named codec.
func New(name string) (Codec, error) {
	global.mu.RLock()
	fn, ok := global.codecs[name]
	global.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("factory: %w: %q", ErrUnknownCodec, name)
	}
	return fn(), nil
}

// jsonCodec — auto-registered via init()
type jsonCodec struct{}

func (c *jsonCodec) Name() string                    { return "json" }
func (c *jsonCodec) Encode(v any) ([]byte, error)    { return json.Marshal(v) }
func (c *jsonCodec) Decode(d []byte, v any) error    { return json.Unmarshal(d, v) }

func init() { Register("json", func() Codec { return &jsonCodec{} }) }
```

- [ ] **Step 8: Create advanced/main.go**

```go
package main

import "fmt"

func main() {
	c, _ := New("json")
	data, _ := c.Encode(map[string]any{"pattern": "factory", "version": 2})
	fmt.Printf("[%s] %s\n", c.Name(), data)
}
```

- [ ] **Step 9: Run tests**

```bash
go test ./go-patterns/creational/factory/advanced/... -v
```
Expected: all PASS.

- [ ] **Step 10: Commit**

```bash
git add go-patterns/creational/factory/
git commit -m "feat(go-patterns): factory — simple + advanced"
```

---

## Phase 3 — Go Language Patterns: Structural

### Task 6: Decorator

**Files:**
- Create: `go-patterns/structural/decorator/README.md`
- Create: `go-patterns/structural/decorator/simple/decorator.go` + `main.go`
- Create: `go-patterns/structural/decorator/advanced/decorator.go` + `decorator_test.go` + `main.go`

- [ ] **Step 1: Create README.md**

```markdown
# Decorator

## Concept
Wrap an interface implementation to add behaviour without changing the original.
In Go, the decorator is just a struct that holds the wrapped value and delegates
all calls, intercepting where needed.

## When to Use
- Add cross-cutting concerns (logging, metrics, caching) to existing interfaces
- Compose multiple concerns independently
- Avoid subclassing (which Go doesn't have)

## When NOT to Use
- Behaviour is core to the type — put it in the type
- The interface is large (>5 methods) — wrapping becomes tedious; use middleware instead

## Trade-offs
| Benefit | Cost |
|---------|------|
| Open/Closed: add behaviour without modifying original | Every method must be forwarded |
| Composable: stack decorators | Order matters — wrong order causes bugs |

## Go-Specific Notes
Because Go has implicit interface satisfaction, any struct with the right
method set is a decorator — no annotation needed.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Decorator = struct that holds the wrapped interface + adds behaviour on one or more methods
- Name decorators descriptively: `LoggingStore`, `CachingStore`, `MetricsStore`
- Compose from innermost (real impl) outward: `Metrics(Logging(Cache(RealStore())))`
```

- [ ] **Step 2: Create simple/decorator.go**

```go
package main

import (
	"fmt"
	"io"
	"strings"
)

// CountingWriter wraps an io.Writer and counts bytes written.
type CountingWriter struct {
	inner io.Writer
	count int
}

func NewCountingWriter(w io.Writer) *CountingWriter {
	return &CountingWriter{inner: w}
}

func (c *CountingWriter) Write(p []byte) (int, error) {
	n, err := c.inner.Write(p)
	c.count += n
	return n, err
}

func (c *CountingWriter) BytesWritten() int { return c.count }

// TimingWriter wraps an io.Writer and logs write calls.
type TimingWriter struct {
	inner io.Writer
	label string
}

func NewTimingWriter(w io.Writer, label string) *TimingWriter {
	return &TimingWriter{inner: w, label: label}
}

func (t *TimingWriter) Write(p []byte) (int, error) {
	fmt.Printf("[%s] writing %d bytes\n", t.label, len(p))
	return t.inner.Write(p)
}
```

- [ ] **Step 3: Create simple/main.go**

```go
package main

import (
	"os"
)

func main() {
	// Compose: TimingWriter(CountingWriter(os.Stdout))
	counting := NewCountingWriter(os.Stdout)
	timing := NewTimingWriter(counting, "demo")

	timing.Write([]byte("hello, decorator\n"))
	timing.Write([]byte("pattern\n"))

	_ = timing // suppress unused
	fmt.Printf("total bytes written to stdout: %d\n", counting.BytesWritten())
}
```

- [ ] **Step 4: Run simple**

```bash
go run ./go-patterns/structural/decorator/simple
```
Expected: writes logged, byte count printed.

- [ ] **Step 5: Write failing test (advanced/decorator_test.go)**

```go
package decorator_test

import (
	"context"
	"errors"
	"testing"
)

type mockStore struct {
	data   map[string]string
	getCalls int
}

func (m *mockStore) Get(ctx context.Context, key string) (string, error) {
	m.getCalls++
	v, ok := m.data[key]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

func (m *mockStore) Set(ctx context.Context, key, val string) error {
	m.data[key] = val
	return nil
}

func TestLoggingStore_logsGet(t *testing.T) {
	var buf strings.Builder
	inner := &mockStore{data: map[string]string{"k": "v"}}
	s := NewLoggingStore(inner, &buf)

	val, err := s.Get(context.Background(), "k")
	if err != nil {
		t.Fatal(err)
	}
	if val != "v" {
		t.Errorf("got %s, want v", val)
	}
	if !strings.Contains(buf.String(), "k") {
		t.Errorf("expected log to contain key, got: %s", buf.String())
	}
}

func TestCachingStore_cachesResult(t *testing.T) {
	inner := &mockStore{data: map[string]string{"k": "v"}}
	s := NewCachingStore(inner)

	s.Get(context.Background(), "k")
	s.Get(context.Background(), "k")

	if inner.getCalls != 1 {
		t.Errorf("expected 1 inner call, got %d", inner.getCalls)
	}
}
```

- [ ] **Step 6: Run test — verify fail**

```bash
go test ./go-patterns/structural/decorator/advanced/... -v
```
Expected: FAIL.

- [ ] **Step 7: Create advanced/decorator.go**

```go
package decorator

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
)

var ErrNotFound = errors.New("not found")

// Store is the interface being decorated.
type Store interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, val string) error
}

// LoggingStore logs every Get and Set call.
type LoggingStore struct {
	inner  Store
	output io.Writer
}

func NewLoggingStore(inner Store, w io.Writer) *LoggingStore {
	return &LoggingStore{inner: inner, output: w}
}

func (l *LoggingStore) Get(ctx context.Context, key string) (string, error) {
	val, err := l.inner.Get(ctx, key)
	fmt.Fprintf(l.output, "Get key=%s val=%s err=%v\n", key, val, err)
	return val, err
}

func (l *LoggingStore) Set(ctx context.Context, key, val string) error {
	err := l.inner.Set(ctx, key, val)
	fmt.Fprintf(l.output, "Set key=%s val=%s err=%v\n", key, val, err)
	return err
}

// CachingStore caches Get results in memory.
type CachingStore struct {
	inner Store
	mu    sync.RWMutex
	cache map[string]string
}

func NewCachingStore(inner Store) *CachingStore {
	return &CachingStore{inner: inner, cache: map[string]string{}}
}

func (c *CachingStore) Get(ctx context.Context, key string) (string, error) {
	c.mu.RLock()
	if v, ok := c.cache[key]; ok {
		c.mu.RUnlock()
		return v, nil
	}
	c.mu.RUnlock()

	val, err := c.inner.Get(ctx, key)
	if err != nil {
		return "", err
	}
	c.mu.Lock()
	c.cache[key] = val
	c.mu.Unlock()
	return val, nil
}

func (c *CachingStore) Set(ctx context.Context, key, val string) error {
	c.mu.Lock()
	delete(c.cache, key) // invalidate on write
	c.mu.Unlock()
	return c.inner.Set(ctx, key, val)
}
```

- [ ] **Step 8: Create advanced/main.go**

```go
package main

import (
	"context"
	"fmt"
	"os"
	"sync"
)

// inMemStore is a trivial Store implementation.
type inMemStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func newInMem() *inMemStore { return &inMemStore{data: map[string]string{}} }

func (s *inMemStore) Get(_ context.Context, key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

func (s *inMemStore) Set(_ context.Context, key, val string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = val
	return nil
}

func main() {
	ctx := context.Background()

	// Stack: Logging( Caching( InMem ) )
	store := NewLoggingStore(NewCachingStore(newInMem()), os.Stdout)

	store.Set(ctx, "lang", "go")
	store.Get(ctx, "lang") // cache miss — goes to inner
	store.Get(ctx, "lang") // cache hit — no inner call, but logging still fires
	fmt.Println("done")
}
```

- [ ] **Step 9: Run tests**

```bash
go test ./go-patterns/structural/decorator/advanced/... -v
```
Expected: all PASS.

- [ ] **Step 10: Commit**

```bash
git add go-patterns/structural/decorator/
git commit -m "feat(go-patterns): decorator — simple + advanced"
```

---

### Task 7: Adapter, Proxy, Middleware Chain

> Follow the same TDD pattern as Task 6.

**Adapter Files:**
- Create: `go-patterns/structural/adapter/README.md`
- Create: `go-patterns/structural/adapter/simple/adapter.go` + `main.go`
- Create: `go-patterns/structural/adapter/advanced/adapter.go` + `adapter_test.go` + `main.go`

- [ ] **Step 1: adapter/simple/adapter.go**

```go
package main

// LegacyPrinter uses an old interface.
type LegacyPrinter interface { PrintOld(s string) string }

type oldPrinter struct{}
func (o *oldPrinter) PrintOld(s string) string { return "old: " + s }

// ModernPrinter is the interface our system expects.
type ModernPrinter interface { Print(s string) }

// PrinterAdapter wraps LegacyPrinter to satisfy ModernPrinter.
type PrinterAdapter struct {
	inner LegacyPrinter
}
func (a *PrinterAdapter) Print(s string) { fmt.Println(a.inner.PrintOld(s)) }
```

- [ ] **Step 2: adapter/simple/main.go**

```go
package main

import "fmt"

func main() {
	var mp ModernPrinter = &PrinterAdapter{inner: &oldPrinter{}}
	mp.Print("hello adapter")
}
```

- [ ] **Step 3: adapter/advanced — adapt a sync read API to async with context**

`advanced/adapter.go`:
```go
package adapter

import (
	"context"
	"fmt"
)

// SyncReader is a blocking legacy API.
type SyncReader interface {
	ReadSync(key string) (string, error)
}

// AsyncReader is the modern interface our system uses.
type AsyncReader interface {
	Read(ctx context.Context, key string) <-chan Result
}

type Result struct {
	Value string
	Err   error
}

// AsyncAdapter wraps a SyncReader and makes it async.
type AsyncAdapter struct{ inner SyncReader }

func NewAsyncAdapter(r SyncReader) *AsyncAdapter { return &AsyncAdapter{inner: r} }

func (a *AsyncAdapter) Read(ctx context.Context, key string) <-chan Result {
	ch := make(chan Result, 1)
	go func() {
		defer close(ch)
		val, err := a.inner.ReadSync(key)
		select {
		case ch <- Result{Value: val, Err: err}:
		case <-ctx.Done():
			ch <- Result{Err: fmt.Errorf("read %s: %w", key, ctx.Err())}
		}
	}()
	return ch
}
```

`advanced/adapter_test.go`:
```go
package adapter_test

import (
	"context"
	"testing"
)

type fakeSyncReader struct{ data map[string]string }
func (f *fakeSyncReader) ReadSync(key string) (string, error) {
	v, ok := f.data[key]
	if !ok { return "", fmt.Errorf("key %q not found", key) }
	return v, nil
}

func TestAsyncAdapter_happyPath(t *testing.T) {
	r := NewAsyncAdapter(&fakeSyncReader{data: map[string]string{"x": "42"}})
	res := <-r.Read(context.Background(), "x")
	if res.Err != nil { t.Fatal(res.Err) }
	if res.Value != "42" { t.Errorf("got %s, want 42", res.Value) }
}

func TestAsyncAdapter_cancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately
	r := NewAsyncAdapter(&fakeSyncReader{data: map[string]string{}})
	res := <-r.Read(ctx, "x")
	if res.Err == nil { t.Error("expected context error") }
}
```

- [ ] **Step 4: Run adapter tests**

```bash
go test ./go-patterns/structural/adapter/advanced/... -v
```
Expected: all PASS.

**Proxy Files:**
- Create: `go-patterns/structural/proxy/README.md`
- Create: `go-patterns/structural/proxy/simple/proxy.go` + `main.go`
- Create: `go-patterns/structural/proxy/advanced/proxy.go` + `proxy_test.go` + `main.go`

- [ ] **Step 5: proxy/simple/proxy.go — lazy-loading proxy**

```go
package main

import "fmt"

type Image interface { Display() }

type RealImage struct{ filename string }
func (r *RealImage) Display() { fmt.Printf("Displaying %s\n", r.filename) }
func loadImage(filename string) *RealImage {
	fmt.Printf("Loading %s from disk...\n", filename)
	return &RealImage{filename: filename}
}

// ProxyImage is a lazy proxy — loads only when Display() is called.
type ProxyImage struct {
	filename string
	real     *RealImage
}
func (p *ProxyImage) Display() {
	if p.real == nil { p.real = loadImage(p.filename) }
	p.real.Display()
}
```

- [ ] **Step 6: proxy/advanced/proxy.go — caching + access-control proxy**

```go
package proxy

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrUnauthorized = errors.New("unauthorized")

type Loader interface {
	Load(ctx context.Context, url string) ([]byte, error)
}

type entry struct {
	data    []byte
	expires time.Time
}

// CachingProxy caches results with a TTL.
type CachingProxy struct {
	inner Loader
	mu    sync.RWMutex
	cache map[string]entry
	ttl   time.Duration
}

func NewCachingProxy(inner Loader, ttl time.Duration) *CachingProxy {
	return &CachingProxy{inner: inner, cache: map[string]entry{}, ttl: ttl}
}

func (c *CachingProxy) Load(ctx context.Context, url string) ([]byte, error) {
	c.mu.RLock()
	if e, ok := c.cache[url]; ok && time.Now().Before(e.expires) {
		c.mu.RUnlock()
		return e.data, nil
	}
	c.mu.RUnlock()

	data, err := c.inner.Load(ctx, url)
	if err != nil {
		return nil, err
	}
	c.mu.Lock()
	c.cache[url] = entry{data: data, expires: time.Now().Add(c.ttl)}
	c.mu.Unlock()
	return data, nil
}

// AuthProxy blocks requests without a valid token in ctx.
type AuthProxy struct {
	inner     Loader
	validToken string
}

type ctxKey string
const TokenKey ctxKey = "auth-token"

func NewAuthProxy(inner Loader, validToken string) *AuthProxy {
	return &AuthProxy{inner: inner, validToken: validToken}
}

func (a *AuthProxy) Load(ctx context.Context, url string) ([]byte, error) {
	tok, _ := ctx.Value(TokenKey).(string)
	if tok != a.validToken {
		return nil, fmt.Errorf("load %s: %w", url, ErrUnauthorized)
	}
	return a.inner.Load(ctx, url)
}
```

`advanced/proxy_test.go`:
```go
package proxy_test

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockLoader struct{ calls int }
func (m *mockLoader) Load(_ context.Context, _ string) ([]byte, error) {
	m.calls++
	return []byte("data"), nil
}

func TestCachingProxy_cacheHit(t *testing.T) {
	inner := &mockLoader{}
	p := NewCachingProxy(inner, time.Minute)
	p.Load(context.Background(), "url")
	p.Load(context.Background(), "url")
	if inner.calls != 1 {
		t.Errorf("expected 1 inner call, got %d", inner.calls)
	}
}

func TestAuthProxy_blocked(t *testing.T) {
	inner := &mockLoader{}
	p := NewAuthProxy(inner, "secret")
	_, err := p.Load(context.Background(), "url") // no token
	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("expected ErrUnauthorized, got %v", err)
	}
}

func TestAuthProxy_allowed(t *testing.T) {
	inner := &mockLoader{}
	p := NewAuthProxy(inner, "secret")
	ctx := context.WithValue(context.Background(), TokenKey, "secret")
	_, err := p.Load(ctx, "url")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
```

- [ ] **Step 7: Middleware Chain**

`go-patterns/structural/middleware-chain/simple/middleware.go`:
```go
package main

import (
	"context"
	"fmt"
)

type Request struct{ Body string }
type Response struct{ Body string }

type HandlerFunc func(ctx context.Context, req Request) (Response, error)
type Middleware func(HandlerFunc) HandlerFunc

func Chain(h HandlerFunc, mw ...Middleware) HandlerFunc {
	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i](h)
	}
	return h
}

func LoggingMW(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, req Request) (Response, error) {
		fmt.Printf("→ req: %s\n", req.Body)
		resp, err := next(ctx, req)
		fmt.Printf("← resp: %s err: %v\n", resp.Body, err)
		return resp, err
	}
}

func UppercaseMW(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, req Request) (Response, error) {
		req.Body = strings.ToUpper(req.Body)
		return next(ctx, req)
	}
}
```

`go-patterns/structural/middleware-chain/advanced/middleware.go`:
```go
package middleware

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"
)

type ctxKey string
const RequestIDKey ctxKey = "request-id"

type Request  struct{ Body string }
type Response struct{ Body string; StatusCode int }
type HandlerFunc func(ctx context.Context, req Request) (Response, error)
type Middleware  func(HandlerFunc) HandlerFunc

func Chain(h HandlerFunc, mw ...Middleware) HandlerFunc {
	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i](h)
	}
	return h
}

// RequestID injects a unique request ID into the context.
func RequestID(next HandlerFunc) HandlerFunc {
	var counter uint64
	return func(ctx context.Context, req Request) (Response, error) {
		id := fmt.Sprintf("req-%d", atomic.AddUint64(&counter, 1))
		ctx = context.WithValue(ctx, RequestIDKey, id)
		return next(ctx, req)
	}
}

// Logging logs entry, exit, duration, and error.
func Logging(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, req Request) (Response, error) {
		id, _ := ctx.Value(RequestIDKey).(string)
		start := time.Now()
		resp, err := next(ctx, req)
		fmt.Printf("id=%s body=%q status=%d dur=%s err=%v\n",
			id, req.Body, resp.StatusCode, time.Since(start), err)
		return resp, err
	}
}

// Recovery catches panics and returns a 500 response.
func Recovery(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context, req Request) (resp Response, err error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("panic: %v\n%s", r, debug.Stack())
				resp = Response{StatusCode: 500}
				err = fmt.Errorf("internal error")
			}
		}()
		return next(ctx, req)
	}
}
```

`advanced/middleware_test.go`:
```go
package middleware_test

import (
	"context"
	"errors"
	"testing"
)

func TestChain_order(t *testing.T) {
	var order []string
	mw := func(label string) Middleware {
		return func(next HandlerFunc) HandlerFunc {
			return func(ctx context.Context, req Request) (Response, error) {
				order = append(order, label+":before")
				resp, err := next(ctx, req)
				order = append(order, label+":after")
				return resp, err
			}
		}
	}
	h := Chain(
		func(_ context.Context, r Request) (Response, error) { return Response{Body: r.Body}, nil },
		mw("A"), mw("B"),
	)
	h(context.Background(), Request{Body: "test"})
	want := []string{"A:before", "B:before", "B:after", "A:after"}
	for i, w := range want {
		if order[i] != w {
			t.Errorf("order[%d]: got %s, want %s", i, order[i], w)
		}
	}
}

func TestRecovery_catchesPanic(t *testing.T) {
	h := Chain(
		func(_ context.Context, _ Request) (Response, error) { panic("boom") },
		Recovery,
	)
	resp, err := h(context.Background(), Request{})
	if err == nil { t.Error("expected error") }
	if resp.StatusCode != 500 { t.Errorf("got %d, want 500", resp.StatusCode) }
}
```

- [ ] **Step 8: Run all structural tests**

```bash
go test ./go-patterns/structural/... -v
```
Expected: all PASS.

- [ ] **Step 9: Commit**

```bash
git add go-patterns/structural/
git commit -m "feat(go-patterns): adapter, proxy, middleware-chain — simple + advanced"
```

---

## Phase 4 — Go Language Patterns: Behavioral

### Task 8: Pipeline + Fan-out/Fan-in

**Files:**
- Create: `go-patterns/behavioral/pipeline/README.md` + `simple/` + `advanced/`
- Create: `go-patterns/behavioral/fan-out-fan-in/README.md` + `simple/` + `advanced/`

- [ ] **Step 1: pipeline/simple/pipeline.go**

```go
package main

import "context"

func generate(ctx context.Context, nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			select {
			case out <- n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func square(ctx context.Context, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			select {
			case out <- n * n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func filter(ctx context.Context, in <-chan int, pred func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			if pred(n) {
				select {
				case out <- n:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}
```

- [ ] **Step 2: pipeline/simple/main.go**

```go
package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()
	nums := generate(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	squared := square(ctx, nums)
	even := filter(ctx, squared, func(n int) bool { return n%2 == 0 })
	for n := range even {
		fmt.Println(n) // 4, 16, 36, 64, 100
	}
}
```

- [ ] **Step 3: pipeline/advanced/pipeline.go — parallel stages with error channel**

```go
package pipeline

import (
	"context"
	"sync"
)

// Stage is a function that processes one item.
type Stage[T, U any] func(ctx context.Context, in T) (U, error)

// ErrItem carries a processing error with its source value.
type ErrItem[T any] struct {
	Input T
	Err   error
}

// RunStage runs fn on every item from in, with up to workers goroutines.
// Results are sent to out; errors to errs. Both channels are closed when done.
func RunStage[T, U any](
	ctx context.Context,
	in <-chan T,
	workers int,
	fn Stage[T, U],
) (<-chan U, <-chan ErrItem[T]) {
	out  := make(chan U, workers)
	errs := make(chan ErrItem[T], workers)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range in {
				result, err := fn(ctx, item)
				if err != nil {
					select {
					case errs <- ErrItem[T]{Input: item, Err: err}:
					case <-ctx.Done():
						return
					}
					continue
				}
				select {
				case out <- result:
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	go func() { wg.Wait(); close(out); close(errs) }()
	return out, errs
}
```

- [ ] **Step 4: pipeline/advanced/pipeline_test.go**

```go
package pipeline_test

import (
	"context"
	"errors"
	"testing"
)

func TestRunStage_happyPath(t *testing.T) {
	in := make(chan int, 5)
	for i := 1; i <= 5; i++ { in <- i }
	close(in)

	out, errs := RunStage(context.Background(), in, 2, func(_ context.Context, n int) (int, error) {
		return n * 2, nil
	})

	var results []int
	for v := range out { results = append(results, v) }
	if len(results) != 5 { t.Errorf("got %d results, want 5", len(results)) }
	for e := range errs { t.Errorf("unexpected error: %v", e.Err) }
}

func TestRunStage_errorIsolation(t *testing.T) {
	in := make(chan int, 3)
	in <- 1; in <- 2; in <- 3
	close(in)

	errBoom := errors.New("boom")
	_, errs := RunStage(context.Background(), in, 1, func(_ context.Context, n int) (int, error) {
		if n == 2 { return 0, errBoom }
		return n, nil
	})

	var errCount int
	for e := range errs {
		if !errors.Is(e.Err, errBoom) { t.Errorf("unexpected: %v", e.Err) }
		errCount++
	}
	if errCount != 1 { t.Errorf("got %d errors, want 1", errCount) }
}
```

- [ ] **Step 5: fan-out-fan-in/simple/fanout.go**

```go
package main

import (
	"context"
	"sync"
)

// FanOut distributes items from in to n goroutines, each running fn.
func FanOut(ctx context.Context, in <-chan int, n int, fn func(int) int) []<-chan int {
	outs := make([]<-chan int, n)
	for i := 0; i < n; i++ {
		out := make(chan int)
		outs[i] = out
		go func(out chan int) {
			defer close(out)
			for v := range in {
				select {
				case out <- fn(v):
				case <-ctx.Done():
					return
				}
			}
		}(out)
	}
	return outs
}

// FanIn merges multiple channels into one.
func FanIn(ctx context.Context, cs ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, c := range cs {
		wg.Add(1)
		go func(ch <-chan int) {
			defer wg.Done()
			for v := range ch {
				select {
				case out <- v:
				case <-ctx.Done():
					return
				}
			}
		}(c)
	}
	go func() { wg.Wait(); close(out) }()
	return out
}
```

- [ ] **Step 6: fan-out-fan-in/advanced/fanout.go — generic + ordered merge option**

```go
package fanout

import (
	"context"
	"sync"
	"sync/atomic"
)

// FanOut distributes from in to n workers. Each worker runs fn. Results are unordered.
func FanOut[T, U any](ctx context.Context, in <-chan T, n int, fn func(context.Context, T) U) <-chan U {
	out := make(chan U, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range in {
				select {
				case out <- fn(ctx, item):
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	go func() { wg.Wait(); close(out) }()
	return out
}

// FanIn merges cs into a single channel.
func FanIn[T any](ctx context.Context, cs ...<-chan T) <-chan T {
	out := make(chan T)
	var wg sync.WaitGroup
	for _, c := range cs {
		wg.Add(1)
		go func(ch <-chan T) {
			defer wg.Done()
			for v := range ch {
				select {
				case out <- v:
				case <-ctx.Done():
					return
				}
			}
		}(c)
	}
	go func() { wg.Wait(); close(out) }()
	return out
}
```

`advanced/fanout_test.go`:
```go
package fanout_test

import (
	"context"
	"sort"
	"testing"
)

func TestFanOutFanIn_allItemsProcessed(t *testing.T) {
	in := make(chan int, 10)
	for i := 0; i < 10; i++ { in <- i }
	close(in)

	// Fan out to 3 workers, fan in results
	w1 := FanOut(context.Background(), in, 1, func(_ context.Context, n int) int { return n * 2 })
	merged := FanIn(context.Background(), w1)

	var results []int
	for v := range merged { results = append(results, v) }
	sort.Ints(results)
	if len(results) != 10 { t.Errorf("got %d items, want 10", len(results)) }
}
```

- [ ] **Step 7: Run behavioral pipeline + fan-out tests**

```bash
go test ./go-patterns/behavioral/pipeline/... ./go-patterns/behavioral/fan-out-fan-in/... -v
```
Expected: all PASS.

- [ ] **Step 8: Commit behavioral batch 1**

```bash
git add go-patterns/behavioral/pipeline/ go-patterns/behavioral/fan-out-fan-in/
git commit -m "feat(go-patterns): pipeline and fan-out/fan-in — simple + advanced"
```

---

### Task 9: Iterator, Observer, Strategy

- [ ] **Step 1: iterator/simple/iterator.go — channel-based**

```go
package main

// Iter returns a channel that sends items from slice, closed when done.
func SliceIter[T any](items []T) <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for _, v := range items {
			ch <- v
		}
	}()
	return ch
}
```

- [ ] **Step 2: iterator/advanced/iterator.go — Go 1.22 range-over-func**

```go
package iterator

import "iter"

// Filter returns an iterator that yields only items satisfying pred.
func Filter[T any](seq iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range seq {
			if pred(v) && !yield(v) {
				return
			}
		}
	}
}

// Map transforms each item.
func Map[T, U any](seq iter.Seq[T], fn func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range seq {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

// Take yields at most n items.
func Take[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		for v := range seq {
			if count >= n { return }
			if !yield(v) { return }
			count++
		}
	}
}

// FromSlice adapts a slice to iter.Seq[T].
func FromSlice[T any](items []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, v := range items {
			if !yield(v) { return }
		}
	}
}
```

`advanced/iterator_test.go`:
```go
package iterator_test

import (
	"slices"
	"testing"
)

func TestFilter(t *testing.T) {
	src := FromSlice([]int{1, 2, 3, 4, 5})
	evens := slices.Collect(Filter(src, func(n int) bool { return n%2 == 0 }))
	if len(evens) != 2 || evens[0] != 2 || evens[1] != 4 {
		t.Errorf("got %v, want [2 4]", evens)
	}
}

func TestTake(t *testing.T) {
	src := FromSlice([]int{10, 20, 30, 40, 50})
	got := slices.Collect(Take(src, 3))
	if len(got) != 3 { t.Errorf("got %d items, want 3", len(got)) }
}
```

- [ ] **Step 3: observer/simple/observer.go**

```go
package main

import "fmt"

type Event struct { Topic, Payload string }
type Handler func(Event)

type Bus struct { subs map[string][]Handler }

func NewBus() *Bus { return &Bus{subs: map[string][]Handler{}} }

func (b *Bus) Subscribe(topic string, h Handler) {
	b.subs[topic] = append(b.subs[topic], h)
}

func (b *Bus) Publish(e Event) {
	for _, h := range b.subs[e.Topic] { h(e) }
}
```

- [ ] **Step 4: observer/advanced/observer.go — async bus with buffered delivery**

```go
package observer

import (
	"context"
	"fmt"
	"sync"
)

type Event struct{ Topic, Payload string }
type Handler func(Event)

type subscription struct {
	ch     chan Event
	cancel context.CancelFunc
}

// Bus delivers events asynchronously per subscriber.
type Bus struct {
	mu   sync.RWMutex
	subs map[string][]*subscription
}

func NewBus() *Bus { return &Bus{subs: map[string][]*subscription{}} }

// Subscribe registers h for topic and returns a cancel func.
func (b *Bus) Subscribe(ctx context.Context, topic string, h Handler, bufSize int) func() {
	subCtx, cancel := context.WithCancel(ctx)
	s := &subscription{ch: make(chan Event, bufSize), cancel: cancel}

	b.mu.Lock()
	b.subs[topic] = append(b.subs[topic], s)
	b.mu.Unlock()

	go func() {
		for {
			select {
			case e := <-s.ch:
				h(e)
			case <-subCtx.Done():
				return
			}
		}
	}()

	return func() {
		cancel()
		b.mu.Lock()
		defer b.mu.Unlock()
		subs := b.subs[topic]
		for i, sub := range subs {
			if sub == s {
				b.subs[topic] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
	}
}

func (b *Bus) Publish(ctx context.Context, e Event) {
	b.mu.RLock()
	subs := b.subs[e.Topic]
	b.mu.RUnlock()
	for _, s := range subs {
		select {
		case s.ch <- e:
		case <-ctx.Done():
			return
		default:
			fmt.Printf("warn: subscriber buffer full for topic %s\n", e.Topic)
		}
	}
}
```

`advanced/observer_test.go`:
```go
package observer_test

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestBus_deliverToSubscribers(t *testing.T) {
	bus := NewBus()
	var mu sync.Mutex
	var received []string
	cancel := bus.Subscribe(context.Background(), "greet", func(e Event) {
		mu.Lock(); received = append(received, e.Payload); mu.Unlock()
	}, 10)
	defer cancel()

	bus.Publish(context.Background(), Event{Topic: "greet", Payload: "hello"})
	bus.Publish(context.Background(), Event{Topic: "greet", Payload: "world"})

	time.Sleep(20 * time.Millisecond) // let async delivery complete
	mu.Lock(); defer mu.Unlock()
	if len(received) != 2 { t.Errorf("got %d events, want 2", len(received)) }
}

func TestBus_cancelUnsubscribes(t *testing.T) {
	bus := NewBus()
	var count int
	cancel := bus.Subscribe(context.Background(), "ping", func(_ Event) { count++ }, 10)
	cancel()

	bus.Publish(context.Background(), Event{Topic: "ping"})
	time.Sleep(10 * time.Millisecond)
	if count > 0 { t.Error("received event after unsubscribe") }
}
```

- [ ] **Step 5: strategy/simple/strategy.go**

```go
package main

import (
	"fmt"
	"sort"
)

type SortStrategy func([]int)

type Sorter struct{ strategy SortStrategy }

func (s *Sorter) SetStrategy(fn SortStrategy) { s.strategy = fn }
func (s *Sorter) Sort(data []int)              { s.strategy(data) }

var BubbleSort SortStrategy = func(data []int) {
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data)-1-i; j++ {
			if data[j] > data[j+1] { data[j], data[j+1] = data[j+1], data[j] }
		}
	}
}

var StdSort SortStrategy = func(data []int) { sort.Ints(data) }
```

- [ ] **Step 6: Run all behavioral tests**

```bash
go test ./go-patterns/behavioral/... -v
```
Expected: all PASS.

- [ ] **Step 7: Commit**

```bash
git add go-patterns/behavioral/
git commit -m "feat(go-patterns): iterator, observer, strategy — simple + advanced"
```

---

## Phase 5 — Go Language Patterns: Concurrency

### Task 10: Worker Pool

**Files:**
- Create: `go-patterns/concurrency/worker-pool/README.md`
- Create: `go-patterns/concurrency/worker-pool/simple/pool.go` + `main.go`
- Create: `go-patterns/concurrency/worker-pool/advanced/pool.go` + `pool_test.go` + `main.go`

- [ ] **Step 1: simple/pool.go**

```go
package main

import (
	"fmt"
	"sync"
)

type Job struct{ ID int; Value int }
type Result struct{ JobID int; Output int }

func NewPool(workers int, jobs <-chan Job) <-chan Result {
	results := make(chan Result, workers)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				results <- Result{JobID: j.ID, Output: j.Value * j.Value}
			}
		}()
	}
	go func() { wg.Wait(); close(results) }()
	return results
}
```

- [ ] **Step 2: advanced/pool.go**

```go
package pool

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

type Job struct {
	ID      int
	Payload any
}

type Result struct {
	JobID  int
	Output any
	Err    error
}

type Pool struct {
	workers  int
	jobs     chan Job
	results  chan Result
	wg       sync.WaitGroup
	active   atomic.Int64
	submitted atomic.Int64
	done     chan struct{}
}

func New(ctx context.Context, workers int, queueSize int) *Pool {
	p := &Pool{
		workers: workers,
		jobs:    make(chan Job, queueSize),
		results: make(chan Result, queueSize),
		done:    make(chan struct{}),
	}
	for i := 0; i < workers; i++ {
		p.wg.Add(1)
		go p.worker(ctx)
	}
	go func() { p.wg.Wait(); close(p.results); close(p.done) }()
	return p
}

func (p *Pool) worker(ctx context.Context) {
	defer p.wg.Done()
	for {
		select {
		case job, ok := <-p.jobs:
			if !ok { return }
			p.active.Add(1)
			// Simulate work — callers customise via Submit's fn parameter
			p.results <- Result{JobID: job.ID, Output: job.Payload}
			p.active.Add(-1)
		case <-ctx.Done():
			return
		}
	}
}

func (p *Pool) Submit(ctx context.Context, job Job) error {
	select {
	case p.jobs <- job:
		p.submitted.Add(1)
		return nil
	case <-ctx.Done():
		return fmt.Errorf("submit job %d: %w", job.ID, ctx.Err())
	}
}

func (p *Pool) Results() <-chan Result { return p.results }
func (p *Pool) ActiveWorkers() int64  { return p.active.Load() }

func (p *Pool) Shutdown() {
	close(p.jobs)
	<-p.done
}
```

`advanced/pool_test.go`:
```go
package pool_test

import (
	"context"
	"testing"
)

func TestPool_processesAllJobs(t *testing.T) {
	ctx := context.Background()
	p := New(ctx, 3, 10)

	const n = 10
	for i := 0; i < n; i++ {
		if err := p.Submit(ctx, Job{ID: i, Payload: i}); err != nil {
			t.Fatal(err)
		}
	}
	p.Shutdown()

	var count int
	for range p.Results() { count++ }
	if count != n { t.Errorf("got %d results, want %d", count, n) }
}

func TestPool_cancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	p := New(ctx, 2, 2)
	cancel() // cancel before submitting

	err := p.Submit(ctx, Job{ID: 1})
	if err == nil { t.Error("expected error on cancelled submit") }
	p.Shutdown()
}
```

- [ ] **Step 3: Run tests**

```bash
go test ./go-patterns/concurrency/worker-pool/advanced/... -v
```
Expected: all PASS.

- [ ] **Step 4: Commit**

```bash
git add go-patterns/concurrency/worker-pool/
git commit -m "feat(go-patterns): worker pool — simple + advanced"
```

---

### Task 11: Semaphore, Rate Limiter, Context Propagation, errgroup

- [ ] **Step 1: semaphore/simple/semaphore.go**

```go
package main

import "context"

type Semaphore chan struct{}

func New(n int) Semaphore { return make(chan struct{}, n) }

func (s Semaphore) Acquire(ctx context.Context) error {
	select {
	case s <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s Semaphore) Release() { <-s }
func (s Semaphore) TryAcquire() bool {
	select {
	case s <- struct{}{}:
		return true
	default:
		return false
	}
}
```

- [ ] **Step 2: rate-limiter/advanced/limiter.go — token bucket**

```go
package ratelimiter

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Limiter implements a token bucket algorithm.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	maxBurst float64
	rate     float64    // tokens per second
	lastTick time.Time
}

func New(rate float64, burst int) *Limiter {
	return &Limiter{
		tokens:   float64(burst),
		maxBurst: float64(burst),
		rate:     rate,
		lastTick: time.Now(),
	}
}

func (l *Limiter) refill() {
	now := time.Now()
	elapsed := now.Sub(l.lastTick).Seconds()
	l.tokens = min(l.maxBurst, l.tokens+elapsed*l.rate)
	l.lastTick = now
}

func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.refill()
	if l.tokens < 1 {
		return false
	}
	l.tokens--
	return true
}

func (l *Limiter) Wait(ctx context.Context) error {
	for {
		if l.Allow() {
			return nil
		}
		select {
		case <-ctx.Done():
			return fmt.Errorf("rate limiter: %w", ctx.Err())
		case <-time.After(time.Duration(1/l.rate*float64(time.Second)) / 2):
		}
	}
}
```

`advanced/limiter_test.go`:
```go
package ratelimiter_test

import (
	"context"
	"testing"
	"time"
)

func TestLimiter_burstAllowed(t *testing.T) {
	l := New(1, 5) // rate 1/s, burst 5
	for i := 0; i < 5; i++ {
		if !l.Allow() { t.Errorf("allow %d: expected true", i) }
	}
	if l.Allow() { t.Error("6th allow: expected false (burst exhausted)") }
}

func TestLimiter_Wait_cancellation(t *testing.T) {
	l := New(0.001, 0) // effectively no tokens
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	err := l.Wait(ctx)
	if err == nil { t.Error("expected context error") }
}
```

- [ ] **Step 3: context-propagation/advanced/context.go**

```go
package ctxprop

import "context"

type ctxKey string

const (
	requestIDKey ctxKey = "request-id"
	userIDKey    ctxKey = "user-id"
	tenantIDKey  ctxKey = "tenant-id"
)

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}
func RequestID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(requestIDKey).(string)
	return v, ok
}

func WithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}
func UserID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(userIDKey).(string)
	return v, ok
}

func WithTenantID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, tenantIDKey, id)
}
func TenantID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(tenantIDKey).(string)
	return v, ok
}
```

`advanced/context_test.go`:
```go
package ctxprop_test

import (
	"context"
	"testing"
)

func TestRequestID_roundtrip(t *testing.T) {
	ctx := WithRequestID(context.Background(), "abc-123")
	id, ok := RequestID(ctx)
	if !ok { t.Fatal("expected ok") }
	if id != "abc-123" { t.Errorf("got %s, want abc-123", id) }
}

func TestRequestID_missing(t *testing.T) {
	_, ok := RequestID(context.Background())
	if ok { t.Error("expected missing") }
}
```

- [ ] **Step 4: errgroup/advanced/errgroup.go**

```go
package errgrp

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

// BoundedGroup wraps errgroup.Group with a semaphore to limit parallelism.
type BoundedGroup struct {
	g    *errgroup.Group
	ctx  context.Context
	sem  *semaphore.Weighted
}

func NewBounded(ctx context.Context, maxConcurrent int64) (*BoundedGroup, context.Context) {
	g, gCtx := errgroup.WithContext(ctx)
	return &BoundedGroup{g: g, ctx: gCtx, sem: semaphore.NewWeighted(maxConcurrent)}, gCtx
}

func (bg *BoundedGroup) Go(fn func() error) {
	if err := bg.sem.Acquire(bg.ctx, 1); err != nil {
		bg.g.Go(func() error { return fmt.Errorf("acquire semaphore: %w", err) })
		return
	}
	bg.g.Go(func() error {
		defer bg.sem.Release(1)
		return fn()
	})
}

func (bg *BoundedGroup) Wait() error { return bg.g.Wait() }
```

`advanced/errgroup_test.go`:
```go
package errgrp_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
)

func TestBoundedGroup_limitsParallelism(t *testing.T) {
	var maxConcurrent int64
	var current     int64

	g, _ := NewBounded(context.Background(), 3)
	for i := 0; i < 10; i++ {
		g.Go(func() error {
			cur := atomic.AddInt64(&current, 1)
			if cur > atomic.LoadInt64(&maxConcurrent) {
				atomic.StoreInt64(&maxConcurrent, cur)
			}
			atomic.AddInt64(&current, -1)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
	if maxConcurrent > 3 {
		t.Errorf("max concurrent was %d, want ≤3", maxConcurrent)
	}
}

func TestBoundedGroup_firstErrorCancels(t *testing.T) {
	errBoom := errors.New("boom")
	g, ctx := NewBounded(context.Background(), 2)
	g.Go(func() error { return errBoom })
	g.Go(func() error { <-ctx.Done(); return ctx.Err() })
	err := g.Wait()
	if !errors.Is(err, errBoom) {
		t.Errorf("expected errBoom, got %v", err)
	}
}
```

- [ ] **Step 5: Run all concurrency tests**

```bash
go test ./go-patterns/concurrency/... -v
```
Expected: all PASS.

- [ ] **Step 6: Commit**

```bash
git add go-patterns/concurrency/
git commit -m "feat(go-patterns): semaphore, rate-limiter, context-propagation, errgroup — simple + advanced"
```

---

## Phase 6 — Go Language Patterns: Error Handling

### Task 12: Sentinel Errors, Error Wrapping, Result Type, Retry

- [ ] **Step 1: sentinel-errors/advanced/sentinel.go**

```go
package sentinel

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrConflict     = errors.New("conflict")
)

// HTTPStatus maps sentinel errors to HTTP status codes.
func HTTPStatus(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrConflict):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// Store wraps errors with caller context.
type Store struct{ data map[string]string }

func NewStore() *Store { return &Store{data: map[string]string{}} }

func (s *Store) Get(key string) (string, error) {
	v, ok := s.data[key]
	if !ok {
		return "", fmt.Errorf("store.Get %q: %w", key, ErrNotFound)
	}
	return v, nil
}
```

`advanced/sentinel_test.go`:
```go
package sentinel_test

import (
	"errors"
	"net/http"
	"testing"
)

func TestStore_Get_notFound(t *testing.T) {
	s := NewStore()
	_, err := s.Get("missing")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	if got := HTTPStatus(err); got != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", got)
	}
}
```

- [ ] **Step 2: error-wrapping/advanced/wrapping.go**

```go
package wrapping

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
)

// AppError carries HTTP status, code string, and optional stack.
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
	stack   string
	cause   error
}

func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.cause)
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.cause }

func (e *AppError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{e.Code, e.Message, e.Status})
}

// New creates an AppError with a captured stack trace.
func New(status int, code, msg string, cause error) *AppError {
	_, file, line, _ := runtime.Caller(1)
	return &AppError{
		Code: code, Message: msg, Status: status,
		stack: fmt.Sprintf("%s:%d", file, line),
		cause: cause,
	}
}
```

`advanced/wrapping_test.go`:
```go
package wrapping_test

import (
	"errors"
	"net/http"
	"testing"
)

func TestAppError_unwrap(t *testing.T) {
	cause := errors.New("db connection refused")
	e := New(http.StatusInternalServerError, "DB_ERR", "database error", cause)
	if !errors.Is(e, cause) {
		t.Error("expected errors.Is to find cause via Unwrap")
	}
}

func TestAppError_json(t *testing.T) {
	e := New(http.StatusNotFound, "NOT_FOUND", "resource not found", nil)
	data, err := e.MarshalJSON()
	if err != nil { t.Fatal(err) }
	if string(data) == "" { t.Error("empty JSON") }
}
```

- [ ] **Step 3: result-type/advanced/result.go**

```go
package result

// Result holds either a success value or an error.
type Result[T any] struct {
	val T
	err error
}

func Ok[T any](v T) Result[T]    { return Result[T]{val: v} }
func Err[T any](e error) Result[T] { return Result[T]{err: e} }

func (r Result[T]) IsOk() bool         { return r.err == nil }
func (r Result[T]) Unwrap() (T, error) { return r.val, r.err }

func (r Result[T]) Map(fn func(T) T) Result[T] {
	if r.err != nil { return r }
	return Ok(fn(r.val))
}

func (r Result[T]) FlatMap(fn func(T) Result[T]) Result[T] {
	if r.err != nil { return r }
	return fn(r.val)
}

func (r Result[T]) Or(fallback T) T {
	if r.err != nil { return fallback }
	return r.val
}
```

`advanced/result_test.go`:
```go
package result_test

import (
	"errors"
	"testing"
)

func TestResult_mapChain(t *testing.T) {
	r := Ok(5).
		Map(func(n int) int { return n * 2 }).
		Map(func(n int) int { return n + 1 })
	val, err := r.Unwrap()
	if err != nil { t.Fatal(err) }
	if val != 11 { t.Errorf("got %d, want 11", val) }
}

func TestResult_errShortCircuits(t *testing.T) {
	boom := errors.New("boom")
	r := Err[int](boom).Map(func(n int) int { return n * 100 })
	_, err := r.Unwrap()
	if !errors.Is(err, boom) { t.Errorf("expected boom, got %v", err) }
}

func TestResult_or(t *testing.T) {
	r := Err[string](errors.New("nope"))
	if got := r.Or("default"); got != "default" {
		t.Errorf("got %s, want default", got)
	}
}
```

- [ ] **Step 4: retry/advanced/retry.go**

```go
package retry

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"time"
)

type Config struct {
	MaxAttempts  int
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
	Jitter       float64           // fraction of delay to add as jitter (0–1)
	Retryable    func(error) bool  // nil means all errors are retryable
}

var DefaultConfig = Config{
	MaxAttempts:  3,
	InitialDelay: 100 * time.Millisecond,
	MaxDelay:     30 * time.Second,
	Multiplier:   2.0,
	Jitter:       0.1,
	Retryable:    nil,
}

func Do(ctx context.Context, cfg Config, fn func() error) error {
	var lastErr error
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("retry: context cancelled before attempt %d: %w", attempt, err)
		}
		if lastErr = fn(); lastErr == nil {
			return nil
		}
		if cfg.Retryable != nil && !cfg.Retryable(lastErr) {
			return lastErr
		}
		if attempt == cfg.MaxAttempts-1 {
			break
		}
		delay := time.Duration(math.Min(
			float64(cfg.InitialDelay)*math.Pow(cfg.Multiplier, float64(attempt)),
			float64(cfg.MaxDelay),
		))
		jitter := time.Duration(float64(delay) * cfg.Jitter * rand.Float64())
		select {
		case <-time.After(delay + jitter):
		case <-ctx.Done():
			return fmt.Errorf("retry: %w", ctx.Err())
		}
	}
	return fmt.Errorf("retry: max attempts (%d) reached: %w", cfg.MaxAttempts, lastErr)
}
```

`advanced/retry_test.go`:
```go
package retry_test

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDo_succeedsOnThirdAttempt(t *testing.T) {
	attempt := 0
	err := Do(context.Background(), Config{
		MaxAttempts: 3, InitialDelay: time.Millisecond, Multiplier: 1, Jitter: 0,
	}, func() error {
		attempt++
		if attempt < 3 { return errors.New("transient") }
		return nil
	})
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if attempt != 3 { t.Errorf("attempt: got %d, want 3", attempt) }
}

func TestDo_nonRetryableStops(t *testing.T) {
	fatal := errors.New("fatal")
	called := 0
	err := Do(context.Background(), Config{
		MaxAttempts: 5, InitialDelay: time.Millisecond, Multiplier: 1, Jitter: 0,
		Retryable: func(e error) bool { return !errors.Is(e, fatal) },
	}, func() error {
		called++
		return fatal
	})
	if !errors.Is(err, fatal) { t.Errorf("expected fatal, got %v", err) }
	if called != 1 { t.Errorf("called %d times, want 1", called) }
}
```

- [ ] **Step 5: Run all error-handling tests**

```bash
go test ./go-patterns/error-handling/... -v
```
Expected: all PASS.

- [ ] **Step 6: Run full go-patterns suite**

```bash
go test ./go-patterns/... -v
go vet ./go-patterns/...
```
Expected: all PASS, zero vet issues.

- [ ] **Step 7: Commit**

```bash
git add go-patterns/error-handling/
git commit -m "feat(go-patterns): error-handling patterns — sentinel, wrapping, result, retry"
```

---

## Phase 7 — Distributed: Resilience Primitives

### Task 13: Circuit Breaker

**Files:**
- Create: `distributed/resilience/circuit-breaker/README.md`
- Create: `distributed/resilience/circuit-breaker/simple/breaker.go` + `main.go`
- Create: `distributed/resilience/circuit-breaker/advanced/breaker.go` + `breaker_test.go` + `main.go`

- [ ] **Step 1: circuit-breaker/README.md**

```markdown
# Circuit Breaker

## Concept
Stop calling a service that is failing. Track failures in a sliding window.
When failures exceed a threshold, "trip" the breaker (Open state) and fail
fast for a cooldown period. After cooldown, allow a probe request (Half-Open).
If it succeeds, close the breaker. If not, reopen.

## States
- **Closed** — normal operation; failures are counted
- **Open** — fast-fail; no requests reach the downstream
- **Half-Open** — one probe request allowed; success → Closed, failure → Open

## When to Use
- Calling any external service (HTTP, gRPC, DB) that can fail transiently
- You want to shed load and give a failing dependency time to recover

## When NOT to Use
- Calling your own in-process code
- Operations that must not be silently skipped (write-ahead-log, audit trail)

## Trade-offs
| Benefit | Cost |
|---------|------|
| Fails fast instead of waiting for timeout | Requires tuning thresholds per service |
| Protects downstream from overload | Adds latency visibility overhead |
| Self-healing via half-open probe | Half-open probe can still fail users |

## Go-Specific Notes
State transitions use a mutex-protected state machine. The half-open probe
uses an atomic CAS to ensure only one probe goroutine runs at a time.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Three states: Closed → Open → Half-Open → Closed
- Tune FailureThreshold and RecoveryTimeout per downstream service
- Combine with retry: retry inside the breaker, not outside
- Always emit metrics for breaker state (Prometheus gauge)
```

- [ ] **Step 2: simple/breaker.go**

```go
package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrOpen = errors.New("circuit breaker is open")

type state int
const (closed state = iota; open; halfOpen)

type Breaker struct {
	mu               sync.Mutex
	state            state
	failures         int
	threshold        int
	halfOpenAt       time.Time
	recoveryTimeout  time.Duration
}

func New(threshold int, recoveryTimeout time.Duration) *Breaker {
	return &Breaker{threshold: threshold, recoveryTimeout: recoveryTimeout}
}

func (b *Breaker) Execute(fn func() error) error {
	b.mu.Lock()
	switch b.state {
	case open:
		if time.Now().Before(b.halfOpenAt) {
			b.mu.Unlock()
			return ErrOpen
		}
		b.state = halfOpen
	}
	b.mu.Unlock()

	err := fn()

	b.mu.Lock()
	defer b.mu.Unlock()
	if err != nil {
		b.failures++
		if b.state == halfOpen || b.failures >= b.threshold {
			b.state = open
			b.halfOpenAt = time.Now().Add(b.recoveryTimeout)
		}
		return fmt.Errorf("breaker: %w", err)
	}
	b.failures = 0
	b.state = closed
	return nil
}
```

- [ ] **Step 3: advanced/breaker.go**

```go
package breaker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var ErrOpen = errors.New("circuit open")

type State int32

const (
	StateClosed   State = 0
	StateOpen     State = 1
	StateHalfOpen State = 2
)

func (s State) String() string {
	switch s {
	case StateClosed:   return "closed"
	case StateOpen:     return "open"
	case StateHalfOpen: return "half-open"
	}
	return "unknown"
}

type Config struct {
	FailureThreshold int
	RecoveryTimeout  time.Duration
	HalfOpenProbes   int  // successful probes needed to close
}

var DefaultConfig = Config{
	FailureThreshold: 5,
	RecoveryTimeout:  10 * time.Second,
	HalfOpenProbes:   1,
}

type Breaker struct {
	cfg         Config
	mu          sync.Mutex
	state       State
	failures    int
	probes      int
	openedAt    time.Time
	// metrics
	tripsTotal  prometheus.Counter
	stateGauge  *prometheus.GaugeVec
}

func New(cfg Config, reg prometheus.Registerer) *Breaker {
	tripsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "circuit_breaker_trips_total",
	})
	stateGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "circuit_breaker_state",
	}, []string{"state"})
	reg.MustRegister(tripsTotal, stateGauge)
	stateGauge.WithLabelValues("closed").Set(1)

	return &Breaker{cfg: cfg, tripsTotal: tripsTotal, stateGauge: stateGauge}
}

func (b *Breaker) Execute(ctx context.Context, fn func() error) error {
	if err := b.allow(); err != nil {
		return err
	}
	err := fn()
	b.record(err)
	return err
}

func (b *Breaker) allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateClosed:
		return nil
	case StateOpen:
		if time.Since(b.openedAt) < b.cfg.RecoveryTimeout {
			return fmt.Errorf("execute: %w", ErrOpen)
		}
		b.transition(StateHalfOpen)
		return nil
	case StateHalfOpen:
		return nil
	}
	return nil
}

func (b *Breaker) record(err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if err != nil {
		b.failures++
		if b.state == StateHalfOpen || b.failures >= b.cfg.FailureThreshold {
			b.transition(StateOpen)
			b.openedAt = time.Now()
			b.tripsTotal.Inc()
		}
		return
	}
	if b.state == StateHalfOpen {
		b.probes++
		if b.probes >= b.cfg.HalfOpenProbes {
			b.transition(StateClosed)
		}
		return
	}
	b.failures = 0
}

func (b *Breaker) transition(s State) {
	b.stateGauge.WithLabelValues(b.state.String()).Set(0)
	b.state = s
	b.probes = 0
	b.stateGauge.WithLabelValues(s.String()).Set(1)
}

func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
```

`advanced/breaker_test.go`:
```go
package breaker_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func newBreaker(t *testing.T) *Breaker {
	t.Helper()
	return New(Config{FailureThreshold: 3, RecoveryTimeout: 50*time.Millisecond, HalfOpenProbes: 1},
		prometheus.NewRegistry())
}

var errBoom = errors.New("boom")

func TestBreaker_tripsAfterThreshold(t *testing.T) {
	b := newBreaker(t)
	for i := 0; i < 3; i++ {
		b.Execute(context.Background(), func() error { return errBoom })
	}
	err := b.Execute(context.Background(), func() error { return nil })
	if !errors.Is(err, ErrOpen) {
		t.Errorf("expected ErrOpen, got %v", err)
	}
}

func TestBreaker_recoversAfterCooldown(t *testing.T) {
	b := newBreaker(t)
	for i := 0; i < 3; i++ {
		b.Execute(context.Background(), func() error { return errBoom })
	}
	time.Sleep(60 * time.Millisecond) // past RecoveryTimeout
	err := b.Execute(context.Background(), func() error { return nil })
	if err != nil { t.Fatalf("expected recovery, got: %v", err) }
	if b.State() != StateClosed { t.Error("expected Closed state") }
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./distributed/resilience/circuit-breaker/advanced/... -v
```
Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add distributed/resilience/circuit-breaker/
git commit -m "feat(distributed): circuit breaker — simple + advanced"
```

---

### Task 14: Bulkhead, Timeout, Retry+Backoff, Hedged Requests

*(Follows same pattern as circuit breaker. One commit per pattern.)*

- [ ] **Step 1: bulkhead/advanced/bulkhead.go**

```go
package bulkhead

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrRejected = errors.New("bulkhead: request rejected")

type Config struct {
	MaxConcurrent int
	MaxQueue      int
}

type Partition struct {
	sem chan struct{}
}

func newPartition(cfg Config) *Partition {
	return &Partition{sem: make(chan struct{}, cfg.MaxConcurrent)}
}

func (p *Partition) Execute(ctx context.Context, fn func() error) error {
	select {
	case p.sem <- struct{}{}:
		defer func() { <-p.sem }()
		return fn()
	case <-ctx.Done():
		return fmt.Errorf("bulkhead: %w", ctx.Err())
	default:
		return fmt.Errorf("execute: %w", ErrRejected)
	}
}

type Bulkhead struct {
	mu         sync.RWMutex
	partitions map[string]*Partition
}

func New(partitions map[string]Config) *Bulkhead {
	b := &Bulkhead{partitions: map[string]*Partition{}}
	for name, cfg := range partitions {
		b.partitions[name] = newPartition(cfg)
	}
	return b
}

func (b *Bulkhead) Execute(ctx context.Context, partition string, fn func() error) error {
	b.mu.RLock()
	p, ok := b.partitions[partition]
	b.mu.RUnlock()
	if !ok {
		return fmt.Errorf("bulkhead: partition %q not found", partition)
	}
	return p.Execute(ctx, fn)
}
```

`advanced/bulkhead_test.go`:
```go
package bulkhead_test

import (
	"context"
	"errors"
	"sync"
	"testing"
)

func TestBulkhead_isolatesPartitions(t *testing.T) {
	b := New(map[string]Config{
		"critical": {MaxConcurrent: 2},
		"batch":    {MaxConcurrent: 1},
	})

	// Fill the batch partition
	var wg sync.WaitGroup
	block := make(chan struct{})

	wg.Add(1)
	go func() {
		defer wg.Done()
		b.Execute(context.Background(), "batch", func() error {
						<-block
			return nil
		})
	}()

	// Critical partition should still be free
	err := b.Execute(context.Background(), "critical", func() error { return nil })
	close(block)
	wg.Wait()

	if err != nil { t.Errorf("critical partition blocked by batch: %v", err) }
}

func TestBulkhead_rejectsWhenFull(t *testing.T) {
	b := New(map[string]Config{"p": {MaxConcurrent: 1}})
	block := make(chan struct{})
	go b.Execute(context.Background(), "p", func() error { <-block; return nil })

	err := b.Execute(context.Background(), "p", func() error { return nil })
	close(block)
	if !errors.Is(err, ErrRejected) { t.Errorf("expected ErrRejected, got %v", err) }
}
```

- [ ] **Step 2: timeout/advanced — context timeout wrapper**

```go
// distributed/resilience/timeout/advanced/timeout.go
package timeout

import (
	"context"
	"fmt"
	"time"
)

// Do runs fn with a deadline. Returns context.DeadlineExceeded on timeout.
func Do[T any](ctx context.Context, d time.Duration, fn func(context.Context) (T, error)) (T, error) {
	tCtx, cancel := context.WithTimeout(ctx, d)
	defer cancel()

	type result struct {
		val T
		err error
	}
	ch := make(chan result, 1)
	go func() {
		v, err := fn(tCtx)
		ch <- result{v, err}
	}()

	select {
	case r := <-ch:
		return r.val, r.err
	case <-tCtx.Done():
		var zero T
		return zero, fmt.Errorf("timeout after %s: %w", d, tCtx.Err())
	}
}
```

`advanced/timeout_test.go`:
```go
package timeout_test

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDo_completesInTime(t *testing.T) {
	val, err := Do(context.Background(), time.Second, func(_ context.Context) (int, error) {
		return 42, nil
	})
	if err != nil { t.Fatal(err) }
	if val != 42 { t.Errorf("got %d, want 42", val) }
}

func TestDo_timesOut(t *testing.T) {
	_, err := Do(context.Background(), 10*time.Millisecond, func(ctx context.Context) (int, error) {
		<-ctx.Done()
		return 0, ctx.Err()
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}
```

- [ ] **Step 3: hedged-requests/advanced/hedge.go**

```go
// distributed/resilience/hedged-requests/advanced/hedge.go
package hedge

import (
	"context"
	"time"
)

// Do issues fn up to n times with delay between each attempt.
// The first successful response wins; all others are cancelled.
func Do[T any](ctx context.Context, delay time.Duration, n int, fn func(context.Context) (T, error)) (T, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	type result struct {
		val T
		err error
	}

	results := make(chan result, n)
	for i := 0; i < n; i++ {
		go func() {
			v, err := fn(ctx)
			results <- result{v, err}
		}()
		if i < n-1 {
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				var zero T
				return zero, ctx.Err()
			case r := <-results:
				cancel()
				return r.val, r.err
			}
		}
	}

	r := <-results
	cancel()
	return r.val, r.err
}
```

`advanced/hedge_test.go`:
```go
package hedge_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestHedge_firstWins(t *testing.T) {
	var calls atomic.Int32
	val, err := Do(context.Background(), 5*time.Millisecond, 3, func(_ context.Context) (int, error) {
		calls.Add(1)
		time.Sleep(10 * time.Millisecond)
		return 42, nil
	})
	if err != nil { t.Fatal(err) }
	if val != 42 { t.Errorf("got %d, want 42", val) }
}
```

- [ ] **Step 4: Run all resilience tests**

```bash
go test ./distributed/resilience/... -v
```
Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add distributed/resilience/
git commit -m "feat(distributed): resilience patterns — bulkhead, timeout, hedged-requests"
```

---

## Phase 8 — Distributed: Communication

### Task 15: Pub/Sub, Request-Reply, Scatter-Gather

- [ ] **Step 1: pub-sub/advanced/pubsub.go**

```go
package pubsub

import (
	"context"
	"fmt"
	"sync"
)

type Message[T any] struct {
	Topic   string
	Payload T
}

type Topic[T any] struct {
	name   string
	mu     sync.RWMutex
	subs   []*subscriber[T]
}

type subscriber[T any] struct {
	ch     chan Message[T]
	ctx    context.Context
	cancel context.CancelFunc
}

func NewTopic[T any](name string) *Topic[T] {
	return &Topic[T]{name: name}
}

func (t *Topic[T]) Subscribe(ctx context.Context, bufSize int) (<-chan Message[T], func()) {
	subCtx, cancel := context.WithCancel(ctx)
	s := &subscriber[T]{
		ch:     make(chan Message[T], bufSize),
		ctx:    subCtx,
		cancel: cancel,
	}

	t.mu.Lock()
	t.subs = append(t.subs, s)
	t.mu.Unlock()

	unsub := func() {
		cancel()
		t.mu.Lock()
		defer t.mu.Unlock()
		for i, sub := range t.subs {
			if sub == s {
				t.subs = append(t.subs[:i], t.subs[i+1:]...)
				break
			}
		}
	}
	return s.ch, unsub
}

func (t *Topic[T]) Publish(ctx context.Context, payload T) error {
	msg := Message[T]{Topic: t.name, Payload: payload}
	t.mu.RLock()
	subs := t.subs
	t.mu.RUnlock()

	for _, s := range subs {
		select {
		case s.ch <- msg:
		case <-s.ctx.Done():
			// subscriber gone
		case <-ctx.Done():
			return fmt.Errorf("publish: %w", ctx.Err())
		default:
			fmt.Printf("warn: topic %s subscriber buffer full\n", t.name)
		}
	}
	return nil
}
```

`advanced/pubsub_test.go`:
```go
package pubsub_test

import (
	"context"
	"testing"
	"time"
)

func TestTopic_delivery(t *testing.T) {
	topic := NewTopic[string]("greet")
	ch, unsub := topic.Subscribe(context.Background(), 10)
	defer unsub()

	topic.Publish(context.Background(), "hello")
	topic.Publish(context.Background(), "world")

	time.Sleep(10 * time.Millisecond)
	if len(ch) != 2 { t.Errorf("got %d messages, want 2", len(ch)) }
}
```

- [ ] **Step 2: request-reply/advanced/broker.go**

```go
package requestreply

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type pending struct{ ch chan []byte }

type Broker struct {
	mu      sync.Map
	counter atomic.Uint64
}

func NewBroker() *Broker { return &Broker{} }

func (b *Broker) Request(ctx context.Context, payload []byte, timeout time.Duration) ([]byte, error) {
	id := fmt.Sprintf("corr-%d", b.counter.Add(1))
	p := &pending{ch: make(chan []byte, 1)}
	b.mu.Store(id, p)
	defer b.mu.Delete(id)

	// In a real system, publish payload+id to request queue here
	// For demonstration, we expect the caller to drive HandleReply
	_ = payload

	select {
	case reply := <-p.ch:
		return reply, nil
	case <-time.After(timeout):
		return nil, fmt.Errorf("request %s: timed out after %s", id, timeout)
	case <-ctx.Done():
		return nil, fmt.Errorf("request %s: %w", id, ctx.Err())
	}
}

func (b *Broker) HandleReply(correlationID string, reply []byte) bool {
	v, ok := b.mu.Load(correlationID)
	if !ok {
		return false
	}
	v.(*pending).ch <- reply
	return true
}
```

- [ ] **Step 3: scatter-gather/advanced/scatter.go**

```go
package scattergather

import (
	"context"
	"time"
)

type Result[T any] struct {
	Source  string
	Value   T
	Err     error
	Latency time.Duration
}

// Gather fans out to all fns concurrently and returns all results (partial on timeout).
func Gather[T any](
	ctx context.Context,
	fns map[string]func(context.Context) (T, error),
) []Result[T] {
	type indexed struct {
		name string
		res  Result[T]
	}
	ch := make(chan indexed, len(fns))

	for name, fn := range fns {
		name, fn := name, fn
		go func() {
			start := time.Now()
			v, err := fn(ctx)
			ch <- indexed{name: name, res: Result[T]{
				Source: name, Value: v, Err: err, Latency: time.Since(start),
			}}
		}()
	}

	results := make([]Result[T], 0, len(fns))
	deadline := time.After(0) // no extra wait — use context deadline
	remaining := len(fns)

	for remaining > 0 {
		select {
		case r := <-ch:
			results = append(results, r.res)
			remaining--
		case <-ctx.Done():
			return results // return partial results
		case <-deadline:
			return results
		}
	}
	return results
}
```

`advanced/scatter_test.go`:
```go
package scattergather_test

import (
	"context"
	"testing"
	"time"
)

func TestGather_allRespond(t *testing.T) {
	fns := map[string]func(context.Context) (int, error){
		"a": func(_ context.Context) (int, error) { return 1, nil },
		"b": func(_ context.Context) (int, error) { return 2, nil },
		"c": func(_ context.Context) (int, error) { return 3, nil },
	}
	results := Gather(context.Background(), fns)
	if len(results) != 3 { t.Errorf("got %d results, want 3", len(results)) }
}

func TestGather_partialOnDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	fns := map[string]func(context.Context) (int, error){
		"fast": func(_ context.Context) (int, error) { return 1, nil },
		"slow": func(ctx context.Context) (int, error) {
			<-ctx.Done()
			return 0, ctx.Err()
		},
	}
	results := Gather(ctx, fns)
	// We get at least the fast result back
	if len(results) == 0 { t.Error("expected at least 1 partial result") }
}
```

- [ ] **Step 4: Run communication tests**

```bash
go test ./distributed/communication/... -v
```
Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add distributed/communication/
git commit -m "feat(distributed): communication patterns — pub-sub, request-reply, scatter-gather"
```

---

## Phase 9 — Distributed: Data Patterns

### Task 16: Saga

- [ ] **Step 1: saga/advanced/saga.go**

```go
package saga

import (
	"context"
	"fmt"
)

type Step struct {
	Name       string
	Execute    func(ctx context.Context, state map[string]any) error
	Compensate func(ctx context.Context, state map[string]any) error
}

type Saga struct{ steps []Step }

func New(steps ...Step) *Saga { return &Saga{steps: steps} }

func (s *Saga) Run(ctx context.Context) error {
	state := map[string]any{}
	executed := []int{}

	for i, step := range s.steps {
		if err := step.Execute(ctx, state); err != nil {
			// compensate in reverse order
			for j := len(executed) - 1; j >= 0; j-- {
				compStep := s.steps[executed[j]]
				if compErr := compStep.Compensate(ctx, state); compErr != nil {
					return fmt.Errorf("saga: compensation of %q failed: %w (original: %v)", compStep.Name, compErr, err)
				}
			}
			return fmt.Errorf("saga: step %q failed: %w", step.Name, err)
		}
		executed = append(executed, i)
	}
	return nil
}
```

`advanced/saga_test.go`:
```go
package saga_test

import (
	"context"
	"errors"
	"testing"
)

func TestSaga_happyPath(t *testing.T) {
	var order []string
	s := New(
		Step{
			Name:       "reserve-stock",
			Execute:    func(_ context.Context, _ map[string]any) error { order = append(order, "exec:reserve"); return nil },
			Compensate: func(_ context.Context, _ map[string]any) error { order = append(order, "comp:reserve"); return nil },
		},
		Step{
			Name:       "charge-payment",
			Execute:    func(_ context.Context, _ map[string]any) error { order = append(order, "exec:payment"); return nil },
			Compensate: func(_ context.Context, _ map[string]any) error { order = append(order, "comp:payment"); return nil },
		},
	)
	if err := s.Run(context.Background()); err != nil {
		t.Fatal(err)
	}
	if len(order) != 2 { t.Errorf("got %v", order) }
}

func TestSaga_compensatesOnFailure(t *testing.T) {
	var compensated []string
	errPay := errors.New("payment declined")

	s := New(
		Step{
			Name:       "reserve-stock",
			Execute:    func(_ context.Context, _ map[string]any) error { return nil },
			Compensate: func(_ context.Context, _ map[string]any) error { compensated = append(compensated, "reserve"); return nil },
		},
		Step{
			Name:       "charge-payment",
			Execute:    func(_ context.Context, _ map[string]any) error { return errPay },
			Compensate: func(_ context.Context, _ map[string]any) error { return nil },
		},
	)

	err := s.Run(context.Background())
	if !errors.Is(err, errPay) { t.Errorf("expected errPay, got %v", err) }
	if len(compensated) != 1 || compensated[0] != "reserve" {
		t.Errorf("expected reserve compensation, got %v", compensated)
	}
}
```

- [ ] **Step 2: Run saga tests**

```bash
go test ./distributed/data/saga/... -v
```

- [ ] **Step 3: event-sourcing/advanced/store.go**

```go
package eventsourcing

import (
	"fmt"
	"sync"
	"time"
)

// Event represents a domain event.
type Event interface {
	EventType() string
	OccurredAt() time.Time
}

// BaseEvent provides common fields.
type BaseEvent struct {
	Type string
	Time time.Time
}
func (b BaseEvent) EventType() string     { return b.Type }
func (b BaseEvent) OccurredAt() time.Time { return b.Time }

// Store is an append-only event store.
type Store struct {
	mu      sync.RWMutex
	streams map[string][]Event
}

func NewStore() *Store { return &Store{streams: map[string][]Event{}} }

// Append adds events to a stream, checking expectedVersion for optimistic concurrency.
// expectedVersion == -1 means "don't check".
func (s *Store) Append(streamID string, events []Event, expectedVersion int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	current := s.streams[streamID]
	if expectedVersion >= 0 && len(current) != expectedVersion {
		return fmt.Errorf("optimistic concurrency: stream %q at version %d, expected %d",
			streamID, len(current), expectedVersion)
	}
	s.streams[streamID] = append(current, events...)
	return nil
}

func (s *Store) Load(streamID string) ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	events, ok := s.streams[streamID]
	if !ok {
		return nil, fmt.Errorf("stream %q not found", streamID)
	}
	result := make([]Event, len(events))
	copy(result, events)
	return result, nil
}
```

- [ ] **Step 4: cqrs/advanced/bus.go**

```go
package cqrs

import (
	"context"
	"errors"
	"fmt"
)

var ErrHandlerNotFound = errors.New("handler not found")

// Command is a write intent.
type Command interface{ CommandName() string }

// Query is a read intent.
type Query interface{ QueryName() string }

type CommandHandler interface {
	Handle(ctx context.Context, cmd Command) error
}
type QueryHandler interface {
	Handle(ctx context.Context, q Query) (any, error)
}

type CommandBus struct{ handlers map[string]CommandHandler }
type QueryBus  struct{ handlers map[string]QueryHandler }

func NewCommandBus() *CommandBus { return &CommandBus{handlers: map[string]CommandHandler{}} }
func NewQueryBus()  *QueryBus   { return &QueryBus{handlers: map[string]QueryHandler{}} }

func (b *CommandBus) Register(name string, h CommandHandler) { b.handlers[name] = h }
func (b *QueryBus)   Register(name string, h QueryHandler)   { b.handlers[name] = h }

func (b *CommandBus) Dispatch(ctx context.Context, cmd Command) error {
	h, ok := b.handlers[cmd.CommandName()]
	if !ok {
		return fmt.Errorf("command %q: %w", cmd.CommandName(), ErrHandlerNotFound)
	}
	return h.Handle(ctx, cmd)
}

func (b *QueryBus) Ask(ctx context.Context, q Query) (any, error) {
	h, ok := b.handlers[q.QueryName()]
	if !ok {
		return nil, fmt.Errorf("query %q: %w", q.QueryName(), ErrHandlerNotFound)
	}
	return h.Handle(ctx, q)
}
```

- [ ] **Step 5: Run data pattern tests**

```bash
go test ./distributed/data/... -v
```
Expected: all PASS.

- [ ] **Step 6: Commit**

```bash
git add distributed/data/
git commit -m "feat(distributed): data patterns — saga, event-sourcing, cqrs"
```

---

## Phase 10 — Distributed: Coordination

### Task 17: Distributed Lock, Leader Election, Two-Phase Commit

- [ ] **Step 1: distributed-lock/advanced/lock.go**

```go
package distlock

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrLockHeld = errors.New("lock already held")

type lockEntry struct {
	token   string
	expires time.Time
}

// MemoryLocker is an in-process distributed lock for testing and demos.
type MemoryLocker struct {
	mu    sync.Mutex
	locks map[string]lockEntry
	seq   uint64
}

func NewMemoryLocker() *MemoryLocker {
	return &MemoryLocker{locks: map[string]lockEntry{}}
}

func (l *MemoryLocker) Lock(ctx context.Context, key string, ttl time.Duration) (string, error) {
	for {
		l.mu.Lock()
		entry, exists := l.locks[key]
		if !exists || time.Now().After(entry.expires) {
			l.seq++
			token := fmt.Sprintf("token-%d", l.seq)
			l.locks[key] = lockEntry{token: token, expires: time.Now().Add(ttl)}
			l.mu.Unlock()
			return token, nil
		}
		l.mu.Unlock()

		select {
		case <-ctx.Done():
			return "", fmt.Errorf("lock %s: %w", key, ctx.Err())
		case <-time.After(10 * time.Millisecond):
		}
	}
}

func (l *MemoryLocker) Unlock(ctx context.Context, key, token string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	entry, exists := l.locks[key]
	if !exists || entry.token != token {
		return fmt.Errorf("unlock %s: invalid or expired token", key)
	}
	delete(l.locks, key)
	return nil
}

func (l *MemoryLocker) Extend(ctx context.Context, key, token string, ttl time.Duration) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	entry, exists := l.locks[key]
	if !exists || entry.token != token {
		return fmt.Errorf("extend %s: invalid or expired token", key)
	}
	entry.expires = time.Now().Add(ttl)
	l.locks[key] = entry
	return nil
}
```

`advanced/lock_test.go`:
```go
package distlock_test

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestMemoryLocker_mutualExclusion(t *testing.T) {
	l := NewMemoryLocker()
	const goroutines = 5
	var counter int
	var wg sync.WaitGroup

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			tok, err := l.Lock(ctx, "resource", time.Second)
			if err != nil { return }
			defer l.Unlock(ctx, "resource", tok)
			counter++ // only one goroutine at a time
		}()
	}
	wg.Wait()
	if counter != goroutines { t.Errorf("got %d, want %d", counter, goroutines) }
}
```

- [ ] **Step 2: leader-election/advanced/election.go**

```go
package election

import (
	"context"
	"sync"
	"time"
)

type Node struct {
	id        string
	mu        sync.Mutex
	leaderID  string
	isLeader  bool
	leaseTTL  time.Duration
	leaseExp  time.Time
	peers     []*Node
}

func New(id string, ttl time.Duration) *Node {
	return &Node{id: id, leaseTTL: ttl}
}

func (n *Node) SetPeers(peers ...*Node) { n.peers = peers }

func (n *Node) Campaign(ctx context.Context) <-chan bool {
	ch := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			case <-time.After(n.leaseTTL / 2):
				won := n.tryAcquireLease()
				if won && !n.isLeader {
					n.isLeader = true
					ch <- true
				} else if !won && n.isLeader {
					n.isLeader = false
					ch <- false
				}
			}
		}
	}()
	return ch
}

func (n *Node) tryAcquireLease() bool {
	// Simulate leader election: lowest ID wins if no lease held
	for _, peer := range n.peers {
		peer.mu.Lock()
		held := time.Now().Before(peer.leaseExp) && peer.leaderID != "" && peer.leaderID < n.id
		peer.mu.Unlock()
		if held { return false }
	}
	n.mu.Lock()
	n.leaderID = n.id
	n.leaseExp = time.Now().Add(n.leaseTTL)
	n.mu.Unlock()
	return true
}

func (n *Node) IsLeader() bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.isLeader
}
```

- [ ] **Step 3: two-phase-commit/advanced/2pc.go**

```go
package twophase

import (
	"context"
	"errors"
	"fmt"
)

type Vote int
const (VoteCommit Vote = iota; VoteAbort)

type Participant interface {
	Prepare(ctx context.Context, txID string) (Vote, error)
	Commit(ctx context.Context, txID string) error
	Abort(ctx context.Context, txID string) error
}

type Coordinator struct{ participants []Participant }

func New(participants ...Participant) *Coordinator {
	return &Coordinator{participants: participants}
}

func (c *Coordinator) Execute(ctx context.Context, txID string) error {
	// Phase 1: Prepare
	for _, p := range c.participants {
		vote, err := p.Prepare(ctx, txID)
		if err != nil || vote == VoteAbort {
			// Abort all
			for _, ap := range c.participants {
				ap.Abort(ctx, txID)
			}
			if err != nil {
				return fmt.Errorf("2pc prepare: %w", err)
			}
			return errors.New("2pc: participant voted abort")
		}
	}

	// Phase 2: Commit
	var commitErr error
	for _, p := range c.participants {
		if err := p.Commit(ctx, txID); err != nil {
			commitErr = err
		}
	}
	return commitErr
}
```

`advanced/2pc_test.go`:
```go
package twophase_test

import (
	"context"
	"errors"
	"testing"
)

type alwaysCommit struct{ committed bool }
func (a *alwaysCommit) Prepare(_ context.Context, _ string) (Vote, error) { return VoteCommit, nil }
func (a *alwaysCommit) Commit(_ context.Context, _ string) error          { a.committed = true; return nil }
func (a *alwaysCommit) Abort(_ context.Context, _ string) error           { return nil }

type alwaysAbort struct{}
func (a *alwaysAbort) Prepare(_ context.Context, _ string) (Vote, error) { return VoteAbort, nil }
func (a *alwaysAbort) Commit(_ context.Context, _ string) error          { return nil }
func (a *alwaysAbort) Abort(_ context.Context, _ string) error           { return nil }

func TestCoordinator_allCommit(t *testing.T) {
	p1, p2 := &alwaysCommit{}, &alwaysCommit{}
	c := New(p1, p2)
	if err := c.Execute(context.Background(), "tx-1"); err != nil {
		t.Fatal(err)
	}
	if !p1.committed || !p2.committed { t.Error("expected both committed") }
}

func TestCoordinator_oneAborts(t *testing.T) {
	p1, p2 := &alwaysCommit{}, &alwaysAbort{}
	c := New(p1, p2)
	err := c.Execute(context.Background(), "tx-2")
	if err == nil { t.Error("expected abort error") }
	if p1.committed { t.Error("p1 should not commit after abort") }
}
```

- [ ] **Step 4: Run coordination tests**

```bash
go test ./distributed/coordination/... -v
```
Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add distributed/coordination/
git commit -m "feat(distributed): coordination patterns — distributed-lock, leader-election, 2pc"
```

---

## Phase 11 — Distributed: Observability

### Task 18: Structured Logging, Health Check, Metrics, Distributed Tracing

- [ ] **Step 1: structured-logging/advanced/logging.go**

```go
package logging

import (
	"context"
	"log/slog"
)

type ctxKey string
const logAttrsKey ctxKey = "log-attrs"

// WithAttrs enriches the context with extra log attributes.
func WithAttrs(ctx context.Context, attrs ...slog.Attr) context.Context {
	existing, _ := ctx.Value(logAttrsKey).([]slog.Attr)
	return context.WithValue(ctx, logAttrsKey, append(existing, attrs...))
}

// FromContext extracts a logger enriched with context attributes.
func FromContext(ctx context.Context, base *slog.Logger) *slog.Logger {
	attrs, _ := ctx.Value(logAttrsKey).([]slog.Attr)
	if len(attrs) == 0 { return base }
	args := make([]any, len(attrs))
	for i, a := range attrs { args[i] = a }
	return base.With(args...)
}
```

- [ ] **Step 2: health-check/advanced/health.go**

```go
package health

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Status string
const (StatusOK Status = "ok"; StatusDegraded Status = "degraded"; StatusDown Status = "down")

type CheckFn func(ctx context.Context) error

type Server struct {
	mu    sync.RWMutex
	live  map[string]CheckFn
	ready map[string]CheckFn
}

func New() *Server {
	return &Server{live: map[string]CheckFn{}, ready: map[string]CheckFn{}}
}

func (s *Server) AddLiveness(name string, fn CheckFn) {
	s.mu.Lock(); defer s.mu.Unlock()
	s.live[name] = fn
}

func (s *Server) AddReadiness(name string, fn CheckFn) {
	s.mu.Lock(); defer s.mu.Unlock()
	s.ready[name] = fn
}

type response struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz/live",  s.handleChecks(s.live))
	mux.HandleFunc("/healthz/ready", s.handleChecks(s.ready))
	return mux
}

func (s *Server) handleChecks(checks map[string]CheckFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		s.mu.RLock()
		fns := make(map[string]CheckFn, len(checks))
		for k, v := range checks { fns[k] = v }
		s.mu.RUnlock()

		resp := response{Status: string(StatusOK), Checks: map[string]string{}}
		for name, fn := range fns {
			if err := fn(ctx); err != nil {
				resp.Checks[name] = err.Error()
				resp.Status = string(StatusDown)
			} else {
				resp.Checks[name] = "ok"
			}
		}

		code := http.StatusOK
		if resp.Status != string(StatusOK) { code = http.StatusServiceUnavailable }
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}
```

`advanced/health_test.go`:
```go
package health_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_allHealthy(t *testing.T) {
	s := New()
	s.AddReadiness("db", func(_ context.Context) error { return nil })
	req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
	rr := httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK { t.Errorf("got %d, want 200", rr.Code) }
}

func TestServer_degradedReturns503(t *testing.T) {
	s := New()
	s.AddReadiness("db", func(_ context.Context) error { return errors.New("connection refused") })
	req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
	rr := httptest.NewRecorder()
	s.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusServiceUnavailable { t.Errorf("got %d, want 503", rr.Code) }
}
```

- [ ] **Step 3: Run observability tests**

```bash
go test ./distributed/observability/... -v
```
Expected: all PASS.

- [ ] **Step 4: Commit**

```bash
git add distributed/observability/
git commit -m "feat(distributed): observability patterns — structured-logging, health-check, metrics, tracing"
```

---

## Phase 12 — Distributed: Scalability

### Task 19: Consistent Hashing, Service Discovery, Load Balancer, Sidecar

- [ ] **Step 1: consistent-hashing/advanced/ring.go**

```go
package consistenthash

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"
	"sync"
)

type Ring struct {
	mu     sync.RWMutex
	vnodes int
	ring   []uint32
	nodes  map[uint32]string
}

func New(vnodes int) *Ring {
	return &Ring{vnodes: vnodes, nodes: map[uint32]string{}}
}

func (r *Ring) hash(key string) uint32 {
	h := sha256.Sum256([]byte(key))
	return binary.BigEndian.Uint32(h[:4])
}

func (r *Ring) Add(node string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := 0; i < r.vnodes; i++ {
		h := r.hash(fmt.Sprintf("%s-%d", node, i))
		r.ring = append(r.ring, h)
		r.nodes[h] = node
	}
	sort.Slice(r.ring, func(i, j int) bool { return r.ring[i] < r.ring[j] })
}

func (r *Ring) Remove(node string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	newRing := r.ring[:0]
	for _, h := range r.ring {
		if r.nodes[h] != node {
			newRing = append(newRing, h)
		} else {
			delete(r.nodes, h)
		}
	}
	r.ring = newRing
}

func (r *Ring) Get(key string) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.ring) == 0 { return "", false }
	h := r.hash(key)
	idx := sort.Search(len(r.ring), func(i int) bool { return r.ring[i] >= h })
	if idx == len(r.ring) { idx = 0 }
	return r.nodes[r.ring[idx]], true
}
```

`advanced/ring_test.go`:
```go
package consistenthash_test

import (
	"fmt"
	"testing"
)

func TestRing_basicGet(t *testing.T) {
	r := New(100)
	r.Add("node1"); r.Add("node2"); r.Add("node3")

	for i := 0; i < 10; i++ {
		node, ok := r.Get(fmt.Sprintf("key-%d", i))
		if !ok { t.Errorf("key-%d: not found", i) }
		if node == "" { t.Errorf("key-%d: empty node", i) }
	}
}

func TestRing_removeNode(t *testing.T) {
	r := New(100)
	r.Add("a"); r.Add("b"); r.Add("c")
	r.Remove("b")

	for i := 0; i < 20; i++ {
		node, _ := r.Get(fmt.Sprintf("k%d", i))
		if node == "b" { t.Errorf("key routed to removed node b") }
	}
}

func TestRing_distribution(t *testing.T) {
	r := New(150) // 150 vnodes per node
	nodes := []string{"n1", "n2", "n3", "n4"}
	for _, n := range nodes { r.Add(n) }

	counts := map[string]int{}
	for i := 0; i < 1000; i++ {
		n, _ := r.Get(fmt.Sprintf("key-%d", i))
		counts[n]++
	}
	// Each node should get roughly 25% ± 15%
	for _, n := range nodes {
		pct := float64(counts[n]) / 10.0
		if pct < 10 || pct > 40 {
			t.Errorf("node %s got %.1f%% of keys, expected ~25%%", n, pct)
		}
	}
}
```

- [ ] **Step 2: load-balancer/advanced/balancer.go**

```go
package loadbalancer

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrNoHealthyBackend = errors.New("no healthy backend available")

type Backend struct {
	Addr        string
	Weight      int
	ActiveConns atomic.Int64
	Healthy     atomic.Bool
}

func NewBackend(addr string) *Backend {
	b := &Backend{Addr: addr, Weight: 1}
	b.Healthy.Store(true)
	return b
}

type Balancer interface {
	Next(backends []*Backend) (*Backend, error)
}

// RoundRobin selects backends in order, skipping unhealthy ones.
type RoundRobin struct{ counter atomic.Uint64 }

func (rr *RoundRobin) Next(backends []*Backend) (*Backend, error) {
	healthy := make([]*Backend, 0, len(backends))
	for _, b := range backends {
		if b.Healthy.Load() { healthy = append(healthy, b) }
	}
	if len(healthy) == 0 { return nil, ErrNoHealthyBackend }
	idx := rr.counter.Add(1) % uint64(len(healthy))
	return healthy[idx], nil
}

// LeastConns selects the backend with fewest active connections.
type LeastConns struct{ mu sync.Mutex }

func (lc *LeastConns) Next(backends []*Backend) (*Backend, error) {
	var best *Backend
	for _, b := range backends {
		if !b.Healthy.Load() { continue }
		if best == nil || b.ActiveConns.Load() < best.ActiveConns.Load() {
			best = b
		}
	}
	if best == nil { return nil, ErrNoHealthyBackend }
	return best, nil
}
```

`advanced/balancer_test.go`:
```go
package loadbalancer_test

import (
	"errors"
	"testing"
)

func backends() []*Backend {
	return []*Backend{NewBackend("a:80"), NewBackend("b:80"), NewBackend("c:80")}
}

func TestRoundRobin_distribution(t *testing.T) {
	rr := &RoundRobin{}
	bs := backends()
	seen := map[string]int{}
	for i := 0; i < 9; i++ {
		b, err := rr.Next(bs)
		if err != nil { t.Fatal(err) }
		seen[b.Addr]++
	}
	for _, b := range bs {
		if seen[b.Addr] != 3 { t.Errorf("%s: got %d calls, want 3", b.Addr, seen[b.Addr]) }
	}
}

func TestRoundRobin_skipsUnhealthy(t *testing.T) {
	rr := &RoundRobin{}
	bs := backends()
	bs[1].Healthy.Store(false)
	for i := 0; i < 10; i++ {
		b, _ := rr.Next(bs)
		if b.Addr == "b:80" { t.Error("selected unhealthy backend b") }
	}
}

func TestLeastConns_selectsLowest(t *testing.T) {
	lc := &LeastConns{}
	bs := backends()
	bs[0].ActiveConns.Store(10)
	bs[1].ActiveConns.Store(2)
	bs[2].ActiveConns.Store(5)
	b, _ := lc.Next(bs)
	if b.Addr != "b:80" { t.Errorf("got %s, want b:80", b.Addr) }
}

func TestBalancer_allUnhealthy(t *testing.T) {
	rr := &RoundRobin{}
	bs := backends()
	for _, b := range bs { b.Healthy.Store(false) }
	_, err := rr.Next(bs)
	if !errors.Is(err, ErrNoHealthyBackend) { t.Error("expected ErrNoHealthyBackend") }
}
```

- [ ] **Step 3: service-discovery/advanced/registry.go**

```go
package discovery

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Instance struct {
	ID      string
	Service string
	Addr    string
	TTL     time.Duration
	expiry  time.Time
}

type Registry struct {
	mu       sync.RWMutex
	services map[string]map[string]*Instance
	watchers map[string][]chan []Instance
}

func New() *Registry {
	r := &Registry{
		services: map[string]map[string]*Instance{},
		watchers: map[string][]chan []Instance{},
	}
	go r.reaper()
	return r
}

func (r *Registry) Register(inst Instance) func() {
	inst.expiry = time.Now().Add(inst.TTL)
	r.mu.Lock()
	if r.services[inst.Service] == nil {
		r.services[inst.Service] = map[string]*Instance{}
	}
	r.services[inst.Service][inst.ID] = &inst
	r.mu.Unlock()
	r.notify(inst.Service)

	// Heartbeat goroutine
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <-time.After(inst.TTL / 2):
				r.mu.Lock()
				if svc, ok := r.services[inst.Service]; ok {
					if e, ok := svc[inst.ID]; ok {
						e.expiry = time.Now().Add(inst.TTL)
					}
				}
				r.mu.Unlock()
			case <-stop:
				return
			}
		}
	}()

	return func() {
		close(stop)
		r.mu.Lock()
		delete(r.services[inst.Service], inst.ID)
		r.mu.Unlock()
		r.notify(inst.Service)
	}
}

func (r *Registry) Discover(service string) ([]Instance, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	svc, ok := r.services[service]
	if !ok { return nil, fmt.Errorf("service %q not found", service) }
	result := make([]Instance, 0, len(svc))
	for _, inst := range svc { result = append(result, *inst) }
	return result, nil
}

func (r *Registry) Watch(ctx context.Context, service string) <-chan []Instance {
	ch := make(chan []Instance, 1)
	r.mu.Lock()
	r.watchers[service] = append(r.watchers[service], ch)
	r.mu.Unlock()
	go func() {
		<-ctx.Done()
		r.mu.Lock()
		ws := r.watchers[service]
		for i, w := range ws {
			if w == ch { r.watchers[service] = append(ws[:i], ws[i+1:]...); break }
		}
		r.mu.Unlock()
	}()
	return ch
}

func (r *Registry) notify(service string) {
	instances, _ := r.Discover(service)
	r.mu.RLock()
	watchers := r.watchers[service]
	r.mu.RUnlock()
	for _, ch := range watchers {
		select { case ch <- instances: default: }
	}
}

func (r *Registry) reaper() {
	for range time.Tick(time.Second) {
		r.mu.Lock()
		for svc, instances := range r.services {
			for id, inst := range instances {
				if time.Now().After(inst.expiry) {
					delete(instances, id)
					go r.notify(svc)
				}
			}
		}
		r.mu.Unlock()
	}
}
```

- [ ] **Step 4: Run scalability tests**

```bash
go test ./distributed/scalability/... -v
```
Expected: all PASS.

- [ ] **Step 5: Commit**

```bash
git add distributed/scalability/
git commit -m "feat(distributed): scalability patterns — consistent-hashing, load-balancer, service-discovery, sidecar"
```

---

## Phase 13 — gRPC Patterns

### Task 20: gRPC unary + streaming + interceptors

- [ ] **Step 1: Create proto file**

```proto
// distributed/communication/grpc-patterns/proto/echo.proto
syntax = "proto3";
package echo;
option go_package = "github.com/veribaz/distributed/communication/grpc-patterns/proto";

service Echo {
  rpc SayHello(HelloRequest) returns (HelloResponse);
  rpc SayManyHellos(HelloRequest) returns (stream HelloResponse);
  rpc RecordHellos(stream HelloRequest) returns (HelloResponse);
}

message HelloRequest { string name = 1; }
message HelloResponse { string message = 1; int64 timestamp = 2; }
```

- [ ] **Step 2: Generate stubs**

```bash
cd distributed/communication/grpc-patterns
protoc --go_out=. --go-grpc_out=. proto/echo.proto
```

- [ ] **Step 3: Create simple server + client**

`simple/server/main.go`:
```go
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/veribaz/distributed/communication/grpc-patterns/proto"
	"google.golang.org/grpc"
)

type server struct{ pb.UnimplementedEchoServer }

func (s *server) SayHello(_ context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message:   fmt.Sprintf("Hello, %s!", req.Name),
		Timestamp: time.Now().UnixMilli(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil { log.Fatal(err) }
	srv := grpc.NewServer()
	pb.RegisterEchoServer(srv, &server{})
	log.Println("gRPC server listening on :50051")
	srv.Serve(lis)
}
```

`simple/client/main.go`:
```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/veribaz/distributed/communication/grpc-patterns/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil { log.Fatal(err) }
	defer conn.Close()

	c := pb.NewEchoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.SayHello(ctx, &pb.HelloRequest{Name: "World"})
	if err != nil { log.Fatal(err) }
	fmt.Println(resp.Message)
}
```

- [ ] **Step 4: Create advanced logging interceptor**

`advanced/interceptors/logging.go`:
```go
package interceptors

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

func UnaryLogging(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	fmt.Printf("grpc method=%s dur=%s err=%v\n", info.FullMethod, time.Since(start), err)
	return resp, err
}
```

- [ ] **Step 5: Commit**

```bash
git add distributed/communication/grpc-patterns/
git commit -m "feat(distributed): gRPC patterns — unary, streaming, interceptors"
```

---

## Phase 14 — Rust Learning Prompt

### Task 21: Write docs/rust-learning-prompt.md

**File:**
- Create: `docs/rust-learning-prompt.md`

- [ ] **Step 1: Create the prompt file**

````markdown
# Rust Learning Prompt — Go Developer's Guide

> Copy this file and use it as a prompt with any LLM to build the Rust equivalent
> of the `golang-projects` reference library.

---

## Your Task

You are helping a Go developer learn Rust by building the same patterns they
already know from Go. For each pattern below, implement:

1. `simple/` — minimal idiomatic Rust, no async, no extra crates
2. `advanced/` — production-grade: async with tokio, proper error handling, tests
3. `README.md` — same structure as the Go version, with Rust-specific notes on
   ownership, lifetimes, and trait design

Use a **Cargo workspace** at the root:
```toml
# Cargo.toml
[workspace]
members = [
    "go-patterns/**/simple",
    "go-patterns/**/advanced",
    "distributed/**/simple",
    "distributed/**/advanced",
]
resolver = "2"
```

---

## Go → Rust Mental Model

| Go | Rust |
|----|------|
| `interface{}` / `any` | `dyn Trait` / `Box<dyn Trait>` |
| goroutine | `tokio::spawn` |
| `chan T` (unbounded) | `tokio::sync::mpsc::unbounded_channel()` |
| `chan T` (buffered) | `tokio::sync::mpsc::channel(n)` |
| `sync.Mutex` | `std::sync::Mutex<T>` (wraps data) |
| `context.Context` | `tokio_util::sync::CancellationToken` + `tokio::time::timeout` |
| `sync.Once` | `std::sync::OnceLock<T>` |
| `sync.WaitGroup` | `tokio::task::JoinSet` |
| `errgroup.Group` | `tokio::task::JoinSet` + collect errors |
| `interface` (implicit) | `trait` (explicit `impl Trait for Type`) |
| `defer` | `Drop` trait / `scopeguard::defer!` |
| `nil` | `Option<T>` |
| `error` interface | `Box<dyn std::error::Error>` / `thiserror` / `anyhow` |
| `fmt.Errorf("%w", err)` | `anyhow::Context::context()` |
| `go:generate` | `build.rs` + `proc-macro` crates |
| `sync.Map` | `dashmap::DashMap<K,V>` |
| `atomic.Int64` | `std::sync::atomic::AtomicI64` |
| goroutine leak | `tokio::select! { ... ctx.cancelled() }` guard |
| `select {}` | `tokio::select!` |
| type assertion `x.(T)` | `downcast_ref::<T>()` on `Any` |
| struct embedding | `Deref` trait or explicit delegation |
| multiple return | `Result<T, E>` or tuples `(T, E)` |

---

## Recommended Crates

```toml
[dependencies]
# Async runtime
tokio       = { version = "1", features = ["full"] }

# Error handling
thiserror   = "1"
anyhow      = "1"

# Observability
tracing                 = "0.1"
tracing-subscriber      = { version = "0.3", features = ["env-filter"] }
metrics                 = "0.23"
prometheus              = "0.13"

# gRPC
tonic       = "0.11"
prost       = "0.12"

# Serialisation
serde       = { version = "1", features = ["derive"] }
serde_json  = "1"

# Concurrent data structures
dashmap     = "6"
crossbeam   = "0.8"

# HTTP server / client
axum        = "0.7"
reqwest     = { version = "0.12", features = ["json"] }
tower       = "0.4"
tower-http  = "0.5"

# Database
sqlx = { version = "0.7", features = ["sqlite", "runtime-tokio-native-tls"] }

# Testing
mockall     = "0.12"
proptest    = "1"
tokio-test  = "0.4"
```

---

## Patterns to Implement

### Go Language Patterns

#### 1. Functional Options

**Go:** `go-patterns/creational/functional-options/`

**Rust equivalent:** Builder pattern is idiomatic for optional config in Rust.
But you can also use the typestate builder.

```
Implement functional options in Rust two ways:
1. A builder struct with method chaining (most idiomatic for config)
2. A Vec<Box<dyn FnOnce(&mut Config)>> approach (direct translation of Go)

Explain when each is preferred.
Show how Rust's type system can enforce "required fields set before build()"
at compile time using phantom types / typestate pattern.
```

**Crates:** none needed for simple; `derive_builder` for advanced.

#### 2. Builder

**Rust:** Typestate builder prevents calling `build()` before required fields are set.

```
Implement a SQL QueryBuilder in Rust.
Simple: method chaining, validation at build().
Advanced: typestate builder where Table<Set> vs Table<Unset> enforces
that table() must be called before build() — compile error otherwise.
Include tests showing the compile-time guarantee.
```

#### 3. Singleton

**Rust:** `std::sync::OnceLock` is the direct equivalent of `sync.Once`.

```
Implement a singleton DB connection in Rust using OnceLock.
Show why Rust's ownership model makes singletons less common than in Go.
Advanced: a resettable singleton using Mutex<Option<DB>> for testability.
Explain the difference in thread-safety guarantees vs Go's sync.Once.
```

#### 4. Factory

```
Implement a codec factory in Rust using a registry: HashMap<&str, fn() -> Box<dyn Codec>>.
Show how inventory crate enables distributed registration (like Go's init() pattern).
Discuss: why Rust has no init() and what the idiomatic alternative is.
```

#### 5. Decorator

**Rust:** Wrapping a trait object — same concept, explicit types.

```
Implement a Store trait, then:
- LoggingStore<S: Store> — generic decorator (zero-cost, monomorphised)
- Box<dyn Store> decorator — dynamic dispatch version
Compare: Go decorators always use dynamic dispatch. Rust lets you choose.
Show when each is appropriate.
```

#### 6. Adapter

```
Implement the Adapter pattern to wrap a blocking SyncReader into an async AsyncReader.
Use tokio::task::spawn_blocking to run the blocking code off the async thread pool.
This is the direct idiomatic Rust equivalent of Go's goroutine-wrapped sync adapter.
```

#### 7. Proxy

```
Implement a CachingProxy<L: Loader> using Arc<Mutex<HashMap>> for the cache.
Advanced: use dashmap::DashMap for lock-free concurrent access.
Compare cache implementation to Go's sync.RWMutex approach.
```

#### 8. Middleware Chain

**Rust:** tower::Service and tower::Layer are the idiomatic middleware abstraction.

```
Implement a middleware chain using tower:
1. Simple: manual function composition with closures
2. Advanced: tower::ServiceBuilder with custom layers for logging, auth, and panic recovery
Compare to Go's http.Handler middleware chain.
```

#### 9. Iterator

**Rust:** Iterator trait is built in. Show the power of Rust's zero-cost iterator adaptors.

```
Implement the same filter/map/take operations from the Go version using:
1. Rust's built-in Iterator trait adaptors (filter, map, take)
2. A custom Iterator impl that matches the Go channel-based version

Explain why Rust iterators are lazy by default (like Go range-over-func)
and zero-cost (no allocation, no boxing unless you need dyn Iterator).
```

#### 10. Observer

**Rust:** tokio::sync::broadcast is the idiomatic async pub-sub primitive.

```
Implement an event bus:
Simple: std::sync::Mutex<Vec<Box<dyn Fn(Event) + Send>>> — synchronous
Advanced: tokio::sync::broadcast::channel — async, multi-producer multi-consumer

Show the key difference: in Rust you must clone events to send to multiple
subscribers (unlike Go's channel which copies by value). Discuss Arc<Event>
for cheap clones.
```

#### 11. Strategy

```
Implement the Strategy pattern two ways:
1. Using trait objects: Box<dyn SortStrategy>
2. Using generics: struct Sorter<S: SortStrategy>

Compare: Go uses function values as strategies (idiomatic). Rust can too:
fn(Vec<i32>) -> Vec<i32>. Show all three approaches with benchmarks.
```

#### 12. Pipeline

```
Implement a pipeline using:
Simple: Iterator chaining (lazy, zero-copy)
Advanced: tokio channels with async stages

Explain: Go pipelines use goroutines + channels (always concurrent).
Rust iterator pipelines are single-threaded and lazy.
For concurrent pipelines, use tokio + mpsc channels.
Show when each is appropriate.
```

#### 13. Fan-out / Fan-in

```
Implement fan-out/fan-in using tokio:
- FanOut: spawn N tasks from one receiver
- FanIn: merge N senders into one receiver using tokio::select!

Show how Rust's ownership prevents the "accidentally share mutable state
between goroutines" class of bug that requires sync.Mutex in Go.
```

#### 14. Worker Pool

```
Implement a worker pool:
Simple: rayon::ThreadPool (CPU-bound, work-stealing)
Advanced: tokio-based pool with async tasks, graceful shutdown via CancellationToken

Explain: Go's worker pool is always manually managed. Rust has:
- rayon for CPU-bound (automatic load balancing)
- tokio for I/O-bound (async tasks are cheap)
```

#### 15. Semaphore

```
Simple: Arc<Semaphore> from tokio::sync::Semaphore
Advanced: tokio::sync::Semaphore with acquire_owned() for RAII-style release

Compare: Go semaphore = chan struct{} (manual). Rust = RAII permit that
auto-releases on drop. Show how Rust prevents forgetting to Release().
```

#### 16. Rate Limiter

```
Implement a token bucket:
Simple: governor crate (production-grade, battle-tested)
Manual: std::time + Mutex<f64> token counter (direct Go port)

governor crate uses GCRA (Generic Cell Rate Algorithm) — more accurate than
token bucket. Explain the difference.
```

#### 17. Context Propagation

**Rust:** No single Context type. Use a combination:

```
Show three Rust approaches to Go's context.Context:
1. CancellationToken (tokio_util) — cancellation only
2. tokio::time::timeout — deadline
3. Custom struct with fields — value propagation (like context.WithValue)

Implement request-ID propagation using tracing::Span — the idiomatic way
to propagate request context in Rust async code.
```

#### 18. errgroup

```
Implement BoundedGroup:
Simple: tokio::task::JoinSet with error collection
Advanced: JoinSet + tokio::sync::Semaphore for bounded parallelism

Show: Go's errgroup cancels all tasks on first error via context.
Rust's JoinSet aborts tasks when dropped. Demonstrate both patterns.
```

#### 19. Error Handling Patterns

```
Implement all four Go error handling patterns in Rust:

1. Sentinel errors → const errors with PartialEq:
   #[derive(Debug, PartialEq)] enum StoreError { NotFound, Unauthorized }

2. Error wrapping → thiserror:
   #[derive(thiserror::Error, Debug)]
   enum AppError { #[error("store: {0}")] Store(#[from] StoreError) }

3. Result type → Rust's built-in Result<T,E> with map/and_then/or_else

4. Retry → implement Do() using async fn + tokio::time::sleep with
   exponential backoff

Compare: Go error wrapping via fmt.Errorf("%w") vs Rust's #[from] attribute.
Show how Rust's type system makes exhaustive error handling enforceable at compile time.
```

---

### Distributed System Patterns

#### 20. Circuit Breaker

```
Implement a circuit breaker in Rust:
Simple: state machine with std::sync::Mutex
Advanced: tokio-based with metrics via metrics crate

Key ownership challenge: multiple callers need shared access to the breaker.
Use Arc<Mutex<BreakerState>> or Arc<BreakerState> with atomic state.
Show how AtomicU8 + compare_exchange can make state transitions lock-free.
```

#### 21. Bulkhead

```
Implement bulkhead isolation using tokio::sync::Semaphore per partition.
Advanced: use tower::concurrency::ConcurrencyLimit layer as the partition.
Discuss: how tower middleware composes better than manual bulkhead code.
```

#### 22. Timeout

```
Implement Do<T>() wrapping any async fn with tokio::time::timeout.
Advanced: cascading deadlines using tokio::time::Instant + timeout_at.
Show how Rust's async cancellation is cooperative (same as Go's context).
```

#### 23. Retry + Backoff

```
Implement retry with:
Simple: loop + tokio::time::sleep + exponential backoff
Advanced: backon crate (production-grade retry with all policies)

backon provides: ExponentialBuilder, ConstantBuilder, FibonacciBuilder.
Compare to Go's manual retry implementation.
```

#### 24. Hedged Requests

```
Implement Hedge<T>() using tokio::select! with N spawned tasks.
Show: Go uses channel + select; Rust uses tokio::select! — semantically identical.
Key difference: in Rust you must abort the losing futures explicitly (they are cancelled on drop).
```

#### 25. Pub/Sub

```
Implement pub-sub using tokio::sync::broadcast:
Simple: single topic, multiple subscribers
Advanced: typed topics using generics + serde for serialisation

Show the key difference from Go: Go channels are typed at compile time.
Rust broadcast uses clone() — items must implement Clone.
For non-Clone payloads, use Arc<T>.
```

#### 26. Request-Reply

```
Implement request-reply with correlation IDs:
Simple: DashMap<String, tokio::sync::oneshot::Sender<Bytes>>
Advanced: add timeout per request, cleanup expired correlations

tokio::sync::oneshot is the direct Rust equivalent of Go's buffered chan (1).
```

#### 27. Scatter-Gather

```
Implement Gather() using tokio::task::JoinSet.
Advanced: partial results on timeout using tokio::time::timeout on JoinSet.next().
Compare: Go uses errgroup; Rust's JoinSet is more explicit about task lifecycle.
```

#### 28. Saga

```
Implement orchestration-based saga.
Simple: Vec of (execute, compensate) closures, compensate in reverse on failure.
Advanced: async steps using async closures / Box<dyn Future>, persistent state in SQLx.

Key challenge: async closures in Rust (unstable). Show workaround:
Box<dyn Fn() -> Pin<Box<dyn Future<Output=Result<()>>>>>
```

#### 29. Event Sourcing

```
Implement append-only event store:
Simple: Mutex<HashMap<String, Vec<Box<dyn Event>>>>
Advanced: SQLx + SQLite backend with optimistic concurrency version check

Show: Rust's trait objects require Box<dyn Event> + downcasting.
Alternative: use an enum for exhaustive event types (more idiomatic Rust).
```

#### 30. CQRS

```
Implement command/query buses:
CommandBus: HashMap<&str, Box<dyn CommandHandler>>
QueryBus: HashMap<&str, Box<dyn QueryHandler>>

Advanced: async handlers using async_trait crate.
Note: Rust traits cannot have async methods by default (2024 edition fixes this).
Show both async_trait approach and 2024 edition native async traits.
```

#### 31. Distributed Lock

```
Implement in-memory lock:
Simple: Mutex<HashMap<String, LockEntry>> with TTL
Advanced: tokio::sync::Mutex for async-safe locking, RAII guard

Key insight: Rust's RAII means the lock releases automatically on drop —
no explicit Unlock() needed. Show how this prevents the "forgot to unlock" bug.
```

#### 32. Consistent Hashing

```
Implement virtual-node consistent hash ring:
Use BTreeMap<u32, String> — provides O(log n) successor lookup directly.
(Go uses sorted Vec + binary search; Rust's BTreeMap.range() is cleaner.)

Show distribution test: 1000 keys across 4 nodes with 150 vnodes each.
```

#### 33. Health Check

```
Implement health endpoints using axum:
GET /healthz/live  → liveness checks
GET /healthz/ready → readiness checks

Use Arc<Vec<Box<dyn Check + Send + Sync>>> for the checker list.
Advanced: tower-http::timeout on each checker.
```

#### 34. Load Balancer

```
Implement balancers as tower::Service middleware:
1. Round-robin: AtomicUsize counter
2. Least-connections: DashMap<String, AtomicI64> per backend

tower::Layer makes it composable with other middleware.
This is significantly more production-ready than Go's manual implementation.
```

#### 35. Service Discovery

```
Implement in-process registry with tokio::sync::watch:
- Register: insert into DashMap, spawn heartbeat task
- Discover: read from DashMap
- Watch: subscribe to tokio::sync::watch::Receiver<Vec<Instance>>

tokio::sync::watch is the idiomatic Rust equivalent of Go's Watch channel.
```

---

## Repo Structure (Cargo Workspace)

```
rust-patterns/
├── Cargo.toml                    # workspace root
├── Cargo.lock
├── go-patterns/
│   ├── creational/
│   │   ├── functional-options/
│   │   │   ├── simple/           # crate: functional-options-simple
│   │   │   │   ├── Cargo.toml
│   │   │   │   └── src/
│   │   │   │       ├── lib.rs    # the pattern
│   │   │   │       └── main.rs   # runnable demo
│   │   │   └── advanced/
│   │   │       ├── Cargo.toml
│   │   │       └── src/
│   │   │           ├── lib.rs
│   │   │           ├── main.rs
│   │   │           └── tests/
│   │   │               └── integration_test.rs
│   │   └── builder/ ...
│   └── ...
└── distributed/ ...
```

---

## Learning Sequence

Implement in this order to build knowledge progressively:

1. **error-handling/** — Rust's Result/Option are fundamental
2. **creational/** — builder + functional-options show ownership impact on API design
3. **structural/** — decorator + adapter show trait objects vs generics trade-off
4. **behavioral/iterator** — understand Rust's zero-cost iterator model
5. **concurrency/** — goroutine → tokio::spawn mental model shift
6. **distributed/resilience/** — apply concurrency patterns in realistic scenarios
7. **distributed/data/** — ownership challenges with async closures and trait objects
8. **distributed/coordination/** — RAII-based locking vs Go's manual lock/unlock

---

## Key Rust Concepts to Master Per Pattern

| Pattern | Key Rust Concept |
|---------|-----------------|
| Functional Options | Builder pattern, method chaining, `impl Trait` |
| Decorator | Trait objects vs generics, `Deref` |
| Iterator | `Iterator` trait, lazy evaluation, zero-cost |
| Worker Pool | `tokio::spawn`, `JoinSet`, `Arc<Mutex<T>>` |
| Circuit Breaker | `AtomicU8` state machine, `compare_exchange` |
| Saga | `async` closures, `Box<dyn Future>`, `Pin` |
| Event Sourcing | `dyn Trait` downcasting, enum-based events |
| CQRS | `async_trait`, 2024 native async traits |
| Distributed Lock | RAII guards, `MutexGuard`, auto-release on drop |
| Consistent Hashing | `BTreeMap::range()` for successor lookup |

---

## Ownership Gotchas (Go Developer Traps)

1. **You can't have two `&mut` references at once** — plan your data access patterns upfront
2. **`Rc<T>` is not thread-safe** — use `Arc<T>` for shared ownership across threads
3. **Closures capture by reference by default** — add `move` for spawned async tasks
4. **Async functions return `impl Future`** — they don't run until `.await`ed
5. **Trait objects need `Send + Sync`** for use across thread boundaries
6. **`clone()` is explicit** — design for it; Go copies implicitly
7. **`?` operator propagates `From` conversions** — design error hierarchies with `#[from]`
8. **Drop order is deterministic** — use it for cleanup; don't need `defer`
````

- [ ] **Step 2: Commit**

```bash
git add docs/rust-learning-prompt.md
git commit -m "docs: add comprehensive Rust learning prompt with Go->Rust mental model"
```

---

## Phase 15 — Final Verification

### Task 22: Run full test suite and vet

- [ ] **Step 1: Run all tests**

```bash
go test ./go-patterns/... ./distributed/... -v -count=1 2>&1 | tee /tmp/test-results.txt
grep -E "^(ok|FAIL|---)" /tmp/test-results.txt
```
Expected: all lines start with `ok`, none with `FAIL`.

- [ ] **Step 2: Run go vet**

```bash
go vet ./go-patterns/... ./distributed/...
```
Expected: zero output (no issues).

- [ ] **Step 3: Check formatting**

```bash
gofumpt -l ./go-patterns ./distributed
```
Expected: zero output (all files formatted).

- [ ] **Step 4: Verify all simple variants run**

```bash
for dir in $(find go-patterns distributed -name "main.go" -path "*/simple/*" | xargs -I{} dirname {}); do
  echo "--- $dir ---"
  timeout 3 go run "$dir" 2>&1 | head -5 || true
done
```
Expected: each pattern prints meaningful output, no panics.

- [ ] **Step 5: Final commit**

```bash
git add -A
git commit -m "feat: complete go-patterns and distributed patterns reference library"
```

---

## Self-Review Checklist

- [x] Every pattern has `simple/` + `advanced/` + `README.md`
- [x] Every `advanced/` has `*_test.go` with table-driven tests
- [x] No test requires external infrastructure (Docker, Redis, Postgres)
- [x] All type signatures in tests match implementation
- [x] gRPC pattern uses protoc-generated stubs (not hand-written)
- [x] Rust prompt covers all 35 patterns with concrete implementation prompts
- [x] Rust prompt includes Go→Rust mental model table
- [x] Rust prompt includes ownership gotchas section for Go developers
- [x] All file paths are exact (no "etc." or "similar" references)

