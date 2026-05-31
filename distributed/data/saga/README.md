# Saga

## Concept
Manage a distributed transaction as a sequence of local steps, each with a
compensating action. If a step fails, previously completed steps are compensated
in reverse order, achieving eventual consistency without a global lock.

## When to Use
- A business transaction spans multiple services/resources.
- Two-phase commit is too costly or unavailable.

## When NOT to Use
- A single ACID transaction can do the job.

## Trade-offs
| Benefit | Cost |
|---------|------|
| No distributed locks | Only eventual consistency |
| Clear failure compensation | Compensations must be designed carefully |

## Go-Specific Notes
This is an orchestration-based saga: a coordinator runs steps forward and, on
failure, calls `Compensate` for each completed step in reverse.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Each step has Execute and Compensate.
- Compensate completed steps in reverse on failure.
- Steps and compensations should be idempotent.
