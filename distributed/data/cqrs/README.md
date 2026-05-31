# CQRS

## Concept
Command Query Responsibility Segregation separates the write model (commands that
mutate state) from the read model (queries that return data). Each side can be
optimised and scaled independently.

## When to Use
- Read and write workloads have very different shapes or scaling needs.
- Complex domains where write validation differs from read projections.
- Often paired with Event Sourcing.

## When NOT to Use
- Simple CRUD where one model serves both reads and writes.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Independent read/write optimisation | More moving parts |
| Clear separation of concerns | Eventual consistency between models |

## Go-Specific Notes
Separate command and query buses dispatch to registered handlers by name. The
write side updates a denormalised read projection the query side serves.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Commands mutate; queries read — never mix.
- A projection bridges the write model to the read model.
- The buses decouple senders from handlers.
