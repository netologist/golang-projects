# Structured Logging

## Concept
Emit logs as machine-readable key/value records (JSON) instead of free text, and
enrich them with request-scoped fields (request ID, trace ID) pulled from the
context. Go's standard `log/slog` is the foundation.

## When to Use
- Production services where logs are ingested by a log platform.
- You need to correlate logs by request/trace ID.

## When NOT to Use
- Throwaway scripts where plain text is fine.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Queryable, correlatable logs | Slightly more verbose call sites |
| Consistent fields across services | Requires discipline on field names |

## Go-Specific Notes
A custom `slog.Handler` reads attributes stashed in the context and adds them to
every record, so handlers log request IDs without threading them manually.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Prefer `log/slog` with a JSON handler in production.
- Enrich records from context (request ID, trace ID).
- Standardise attribute keys across services.
