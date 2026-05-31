# Load Balancer

## Concept
Distribute requests across a set of backends. Strategies include round-robin
(even rotation) and least-connections (route to the least busy). Health awareness
skips unhealthy backends.

## When to Use
- Spreading traffic across replicas of a stateless service.
- Client-side load balancing in a service mesh or RPC client.

## When NOT to Use
- A single backend, or when an external LB already handles distribution.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Even utilisation, failover | Needs health signals to be effective |
| Pluggable strategies | Least-conns requires connection tracking |

## Go-Specific Notes
A `Balancer` interface picks the next backend. Backends track health and active
connections with atomics, so selection is lock-light.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Round-robin is simple and fair for uniform backends.
- Least-connections adapts to uneven request costs.
- Always skip unhealthy backends.
