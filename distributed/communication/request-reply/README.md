# Request-Reply

## Concept
Asynchronous request/response over a message channel using a correlation ID. The
requester registers a pending reply slot keyed by ID; the responder routes its
reply back by that same ID.

## When to Use
- RPC-style interaction over an async transport (message queue, event bus).
- You need to match replies to requests across goroutines or nodes.

## When NOT to Use
- Synchronous in-process calls — just call the function.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Works over async transports | Requires correlation bookkeeping |
| Decouples request from reply timing | Must handle timeouts and orphan replies |

## Go-Specific Notes
A `sync.Map` keyed by correlation ID holds a buffered reply channel per request.
`Request` blocks with a timeout; `HandleReply` delivers and returns whether a
waiter existed.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Correlation IDs match replies to requests.
- Always bound the wait with a timeout/context.
- Clean up the pending slot on completion or timeout.
