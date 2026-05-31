# Iterator

## Concept
Provide sequential access to elements without exposing the underlying structure.
Go 1.23+ formalises this with `iter.Seq[T]` (range-over-func), enabling lazy,
composable iterators.

## When to Use
- Lazily produce or transform sequences (pagination, streams, filters).
- Hide the backing collection from consumers.

## When NOT to Use
- A plain slice and `range` are enough.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Lazy evaluation, composable adaptors | Slightly more abstract than a slice |
| Uniform interface across sources | Requires Go 1.23+ for iter.Seq |

## Go-Specific Notes
`iter.Seq[T]` is `func(yield func(T) bool)`. Returning `false` from `yield`
stops iteration early, which adaptors must honour.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- An iterator is `func(yield func(T) bool)`.
- Always check `yield`'s bool return and stop when it is false.
- `slices.Collect` materialises an iterator into a slice.
