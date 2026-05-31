# Adapter

## Concept
Bridge two incompatible interfaces by wrapping one to satisfy the other.
Our system speaks one interface; a legacy or third-party component speaks
another. The adapter translates between them.

## When to Use
- Integrate a legacy or third-party API without changing your call sites.
- Convert a synchronous API into an asynchronous one (or vice versa).

## When NOT to Use
- You control both sides — change one interface directly.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Isolates third-party changes behind one type | Extra indirection layer |
| Lets you unit-test against your own interface | Adapter can hide latency/behaviour quirks |

## Go-Specific Notes
Because interface satisfaction is implicit, the adapter only needs the right
method set. The async adapter uses a goroutine plus a channel and honours
`context.Context` cancellation.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- The adapter implements *your* interface and delegates to *their* implementation.
- For sync→async, wrap the blocking call in a goroutine and respect context.
