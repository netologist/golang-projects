# Distributed Tracing

## Concept
Trace a request as it flows across functions and service boundaries by recording
spans linked into a single trace. OpenTelemetry (OTel) is the standard API; spans
carry attributes, events, and parent/child relationships.

## When to Use
- Diagnosing latency and errors across multiple services.
- Understanding causal request flow in a distributed system.

## When NOT to Use
- A single, simple process where logs suffice.

## Trade-offs
| Benefit | Cost |
|---------|------|
| End-to-end request visibility | Exporter/collector infrastructure |
| Pinpoints slow spans | Sampling needed at high volume |

## Go-Specific Notes
The example uses the OTel SDK with a stdout exporter (no collector needed). A
parent span wraps child spans for sub-operations; context carries the active span.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Spans nest via context propagation.
- Add attributes and record errors on spans.
- Use a real exporter (OTLP/Jaeger) in production; stdout for local dev.
