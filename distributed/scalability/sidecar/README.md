# Sidecar

## Concept
Attach a helper component alongside a main service to handle cross-cutting
concerns (logging, retries, TLS, header injection) without changing the service
itself. The sidecar sits on the request path as a proxy.

## When to Use
- Add capabilities (observability, resilience) uniformly across services.
- Keep language-agnostic concerns out of application code (service mesh model).

## When NOT to Use
- A single in-process middleware suffices and no isolation is needed.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Decouples infra concerns from app | Extra network hop / process |
| Reusable across services/languages | More moving parts to operate |

## Go-Specific Notes
This example models the sidecar as an in-process reverse proxy that enriches
requests (adds headers, logs) and retries idempotent upstream failures, avoiding
the complexity of a real subprocess while showing the pattern.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- The sidecar intercepts traffic to add cross-cutting behaviour.
- It is transparent to the upstream service.
- Real sidecars run as separate processes/containers (e.g., Envoy).
