# Coordination

Algorithms and patterns for reaching agreement and coordinating actions
across distributed nodes.

| Pattern | Description |
|---------|-------------|
| [Paxos](./paxos) | Classic consensus algorithm — proposers, acceptors, learners |
| [Raft](./raft) | Leader-based consensus with log replication and safety guarantees |
| [Two-Phase Commit](./two-phase-commit) | Atomic commit across participants (prepare → commit) |
| [Three-Phase Commit](./three-phase-commit) | Non-blocking atomic commit with pre-commit phase |
| [Leader Election](./leader-election) | Bully algorithm style leader election among peers |
| [Distributed Lock](./distributed-lock) | Mutex across processes with lease-based expiry |
