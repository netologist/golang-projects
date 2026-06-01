# Scalability

Patterns for distributing load, data, and services across nodes to
handle growth.

| Pattern | Description |
|---------|-------------|
| [Consistent Hashing](./consistent-hashing) | Hash ring with virtual nodes for minimal key remapping |
| [Rendezvous Hashing](./rendezvous-hashing) | Highest random weight (HRW) hashing for deterministic assignment |
| [Load Balancer](./load-balancer) | Round-robin, least connections, and weighted strategies |
| [Service Discovery](./service-discovery) | Registry-based dynamic endpoint resolution |
| [Sidecar](./sidecar) | Out-of-process proxy for cross-cutting concerns |
| [Range Partitioning](./range-partitioning) | Key-range sharding with rebalancing via splits |
