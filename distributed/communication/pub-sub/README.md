# Pub/Sub

## Concept
Decouple producers from consumers via named topics. Publishers emit messages to
a topic; any number of subscribers receive them independently, each with its own
delivery channel.

## When to Use
- Event broadcasting to multiple independent consumers.
- Decoupling services that should not call each other directly.

## When NOT to Use
- Point-to-point request/response — use request-reply.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Loose coupling, many consumers | At-most-once on buffer overflow (configurable) |
| Independent subscriber lifecycles | In-process only here (no durability) |

## Go-Specific Notes
Each subscriber gets a buffered channel and a cancel func. Publish is
non-blocking: a full subscriber buffer drops the message with a warning rather
than blocking the publisher.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- One buffered channel per subscriber isolates slow consumers.
- Subscriptions return a cancel func tied to a context.
- Choose a drop-vs-block policy explicitly.
