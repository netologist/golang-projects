# Context Propagation

## Concept
Carry request-scoped values (request ID, user, tenant) and cancellation across
API boundaries and goroutines using `context.Context`. Typed keys prevent
collisions and accidental access.

## When to Use
- Pass request metadata through layers without changing every signature.
- Propagate deadlines and cancellation to downstream calls.

## When NOT to Use
- Storing optional function parameters — pass them explicitly instead.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Uniform cancellation + metadata flow | Values are untyped at the boundary |
| Decouples layers | Overuse hides real dependencies |

## Go-Specific Notes
Use an unexported `ctxKey` type so only this package can set/read its values.
Expose typed getters/setters; never store a context in a struct.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Define an unexported key type to avoid collisions.
- Provide typed `With*`/getter helpers.
- Context is for request scope and cancellation, not general storage.
