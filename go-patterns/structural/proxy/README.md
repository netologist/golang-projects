# Proxy

## Concept
Provide a surrogate with the same interface as the real object, controlling
access to it. Common variants: lazy-loading, caching, and access-control proxies.

## When to Use
- Defer expensive creation until first use (lazy/virtual proxy).
- Add caching, authorisation, or rate limiting transparently.

## When NOT to Use
- No access control or laziness is needed — call the object directly.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Transparent to callers (same interface) | Adds a layer of indirection |
| Centralises cross-cutting access control | Can mask the cost of the real call |

## Go-Specific Notes
A proxy implements the same interface and forwards to the real implementation.
The caching proxy uses `sync.RWMutex` and a TTL; the auth proxy reads a token
from `context.Context`.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Proxy shares the subject's interface and decides whether/when to forward.
- Compose proxies: auth in front of caching in front of the real loader.
