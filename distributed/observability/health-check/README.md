# Health Check

## Concept
Expose liveness and readiness endpoints so orchestrators (e.g., Kubernetes) know
whether to restart a pod (liveness) or route traffic to it (readiness). Each is a
set of named checks aggregated into a status.

## Liveness vs Readiness
- **Liveness** — is the process alive? Failing restarts the pod.
- **Readiness** — can it serve traffic now? Failing removes it from load balancing.

## When to Use
- Any service deployed under an orchestrator or behind a load balancer.

## When NOT to Use
- Short-lived batch jobs with no traffic.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Automated recovery and traffic gating | Bad checks cause flapping/restarts |
| Dependency visibility | Checks add a small runtime cost |

## Go-Specific Notes
Checks are `func(ctx) error`. The handler runs them with a timeout and returns
200 or 503 with a JSON body listing each check's status.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Separate liveness from readiness.
- Bound checks with a timeout.
- Return a machine-readable body for debugging.
