# Singleton

## Concept
Ensure a type has exactly one instance and provide a global access point.
In Go, lazy initialisation uses `sync.Once` — goroutine-safe and allocation-free
after the first call.

## When to Use
- Shared resources (DB pool, logger, config) that are expensive to create.
- Exactly one instance is correct by design.

## When NOT to Use
- Testing: singletons are hard to mock — prefer dependency injection.
- Multiple configurations needed — use a factory instead.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Guaranteed single instance | Global state — a hidden dependency |
| Thread-safe lazy init | Hard to reset in tests |
| Zero overhead after first call | Init errors are hard to surface |

## Go-Specific Notes
`sync.Once.Do` runs the function exactly once across all goroutines. For
testable code, expose a `Reset()` helper guarded for test use.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Use `sync.Once` — never roll your own double-checked locking.
- Expose a `Reset()` for tests.
- Prefer dependency injection over singletons in library code.
- Singletons are acceptable for truly global resources (logger, metrics registry).
