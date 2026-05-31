# Builder

## Concept
Separate the construction of a complex object from its representation.
A Builder accumulates configuration through method calls, then produces a
validated immutable result via `Build()`.

## When to Use
- Object requires multi-step construction with validation at the end.
- You want method chaining (a fluent API).
- Construction involves conditional branches.

## When NOT to Use
- Simple structs — struct literals are clearer.
- When Functional Options already covers the need.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Validation at one point (Build) | More boilerplate than functional options |
| Fluent, readable call sites | Builder is stateful — not safe to reuse |

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
- `Build()` is the only place that validates.
- Setters return `*Builder` — they never fail.
- The returned type is immutable (fields unexported, accessed via methods).
- A builder should not be reused after `Build()`.
