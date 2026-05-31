# Middleware Chain

## Concept
Compose a request/response pipeline from small, single-purpose wrappers. Each
middleware takes the next handler and returns a new handler, so concerns like
logging, auth, and recovery stack cleanly.

## When to Use
- Cross-cutting concerns across many handlers (logging, metrics, auth, recovery).
- You want to add/remove/reorder concerns without touching handler logic.

## When NOT to Use
- A single handler with one concern — inline it.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Reusable, composable concerns | Execution order can be surprising |
| Handlers stay focused on business logic | Deep stacks complicate debugging |

## Go-Specific Notes
The idiom mirrors `net/http` middleware: `func(Handler) Handler`. `Chain`
applies middleware so the first listed runs outermost.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- A middleware is `func(next HandlerFunc) HandlerFunc`.
- `Chain` wraps in reverse so the first middleware is the outermost layer.
- Put recovery outermost so it catches panics from everything inside.
