# Go Patterns

Idiomatic Go patterns covering concurrency, error handling, creational,
structural, and behavioural design patterns. Each pattern includes a simple
runnable example and an advanced version with tests.

## Categories

| Category | Description |
|----------|-------------|
| [Concurrency](./concurrency) | Worker pool, context propagation, errgroup, semaphore, rate limiter |
| [Error Handling](./error-handling) | Retry, result type, sentinel errors, error wrapping |
| [Creational](./creational) | Functional options, singleton, builder, factory |
| [Structural](./structural) | Proxy, decorator, adapter, middleware chain |
| [Behavioural](./behavioral) | Fan-out/fan-in, pipeline, observer, iterator, strategy |

## Running

```bash
make test          # Run all tests
make test-verbose  # Verbose with race detection
make lint          # Lint all packages
make fmt           # Format with gofumpt + goimports
```

Each pattern can be run individually:

```bash
go run ./concurrency/worker-pool/advanced
go test ./error-handling/retry/advanced/... -v
```

## Key Takeaways

- Patterns are idiomatic Go — using goroutines, channels, and interfaces.
- Simple examples focus on the core concept; advanced examples add tests and edge cases.
- Prioritises clarity and correctness over performance.
