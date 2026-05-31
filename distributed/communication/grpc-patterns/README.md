# gRPC Patterns

## Concept
gRPC is a high-performance RPC framework using Protocol Buffers for typed
contracts and HTTP/2 for transport. It supports four call types: unary,
server-streaming, client-streaming, and bidirectional streaming, plus
interceptors for cross-cutting concerns.

## When to Use
- Typed, efficient service-to-service communication.
- Streaming workloads (telemetry, feeds, chat).
- Polyglot environments needing a shared contract (.proto).

## When NOT to Use
- Public browser-facing APIs (REST/JSON is friendlier without gRPC-Web).
- Very simple internal calls where HTTP/JSON suffices.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Typed contracts, codegen | Requires protoc toolchain |
| Streaming + HTTP/2 efficiency | Harder to debug than JSON |
| Interceptors for cross-cutting | Steeper learning curve |

## Go-Specific Notes
`protoc` generates message and service stubs. The server implements the generated
interface; interceptors wrap unary/stream handlers (logging, auth, metrics).

## Layout
```
proto/            generated stubs (echo.pb.go, echo_grpc.pb.go) + echo.proto
simple/server/    unary server
simple/client/    unary client
advanced/         interceptors (logging) + an in-process round-trip test
```

## Regenerating Stubs
```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/echo.proto
```

## Running
```bash
# Terminal 1
go run ./simple/server
# Terminal 2
go run ./simple/client

# Advanced (in-process, no ports):
go test ./advanced/... -v
```

## Key Takeaways
- Define the contract in .proto; generate typed stubs.
- Implement the generated server interface.
- Use interceptors for logging, auth, retries, and tracing.
