# Result Type

## Concept
Encapsulate either a success value or an error in a single generic type, enabling
functional-style chaining (`Map`, `FlatMap`) that short-circuits on the first error.

## When to Use
- Pipelines of fallible transformations where you want to defer error handling.
- You prefer composition over repeated `if err != nil` at each step.

## When NOT to Use
- Idiomatic Go `(T, error)` is clearer for most code — use this sparingly.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Clean chaining, single error path | Non-idiomatic for many Go teams |
| Generic and reusable | Hides error sites if overused |

## Go-Specific Notes
Built on Go generics. `Map` skips when the Result already holds an error, so an
error propagates unchanged to the end of the chain.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- `Ok`/`Err` constructors; `Unwrap` returns `(T, error)`.
- `Map`/`FlatMap` short-circuit on error.
- Use as a focused tool, not a replacement for idiomatic Go errors.
