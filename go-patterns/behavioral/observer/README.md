# Observer

## Concept
Decouple event producers from consumers. Subscribers register interest in a
topic; publishers emit events without knowing who listens. In Go this is built
naturally with channels.

## When to Use
- One-to-many event notification (UI events, domain events, cache invalidation).
- You want producers and consumers to evolve independently.

## When NOT to Use
- A direct function call is simpler and the coupling is acceptable.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Loose coupling between producer/consumer | Harder to trace flow at runtime |
| Many subscribers without producer changes | Slow subscribers can cause back-pressure |

## Go-Specific Notes
The advanced bus delivers asynchronously via a buffered channel per subscriber
and drops events for full buffers (configurable), preventing one slow consumer
from blocking the publisher.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- A subscription returns a cancel func to unsubscribe.
- Async delivery needs a per-subscriber buffer and a drop/block policy.
- Always tie subscriber lifetime to a context.
