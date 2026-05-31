# Strategy

## Concept
Define a family of interchangeable algorithms and select one at runtime. In Go,
a function type is the most idiomatic strategy — first-class and zero-ceremony.

## When to Use
- Multiple algorithms for the same task (sorting, compression, routing).
- The choice depends on runtime conditions (payload size, config).

## When NOT to Use
- Only one algorithm exists and will not change.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Swap behaviour without touching callers | Indirection can obscure which runs |
| Easy to unit-test each strategy | Too many strategies fragment logic |

## Go-Specific Notes
A strategy is just a `func` value. The advanced example selects a compression
strategy based on payload size.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Prefer a function type over an interface for single-method strategies.
- Selection logic lives in one place; strategies stay independent.
