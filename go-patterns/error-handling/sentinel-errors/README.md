# Sentinel Errors

## Concept
Define named, package-level error values that callers compare against with
`errors.Is`. Wrapping preserves the sentinel while adding context.

## When to Use
- Callers must branch on specific, well-known error conditions.
- You want a stable error API (e.g., `ErrNotFound`).

## When NOT to Use
- Errors carry structured data — use a custom error type with `errors.As`.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Simple, comparable error identity | Only identity, no structured fields |
| Stable public API | Proliferation of sentinels can sprawl |

## Go-Specific Notes
Wrap sentinels with `fmt.Errorf("...: %w", ErrNotFound)` so `errors.Is` still
matches. Map sentinels to transport codes (e.g., HTTP) in one place.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Compare with `errors.Is`, never `==` on wrapped errors.
- Keep sentinels few and meaningful.
- Centralise sentinel-to-status mapping.
