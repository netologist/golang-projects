# Error Wrapping

## Concept
Add context to errors as they propagate up the stack without discarding the
original. `fmt.Errorf("%w", err)` wraps; `errors.Is`/`errors.As` inspect.

## When to Use
- You want a readable error chain ("load config: open file: permission denied").
- Callers need to inspect either identity (`Is`) or a typed cause (`As`).

## When NOT to Use
- The extra context adds no value — return the original error.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Preserves cause for inspection | Over-wrapping yields noisy messages |
| Rich, layered diagnostics | Custom types need Unwrap implemented |

## Go-Specific Notes
A custom error type implements `Error()` and `Unwrap()` so `errors.As` and
`errors.Is` traverse the chain. The advanced type also serialises to JSON for
API responses.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Wrap with `%w` to keep the cause inspectable.
- Implement `Unwrap()` on custom error types.
- Add context at each layer, but only meaningful context.
