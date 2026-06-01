# Replication

Patterns for copying and synchronising data across multiple nodes.

| Pattern | Description |
|---------|-------------|
| [CRDT](./crdt) | Conflict-free replicated data types (G-Counter, PN-Counter, LWW-Register, 2P-Set) |
| [WAL](./wal) | Write-ahead log for crash recovery with checkpoint support |
| [Primary-Replica](./primary-replica) | Single-primary replication with async propagation |
| [Read Repair](./read-repair) | Detect and fix stale replicas on read |
| [Multi-Leader](./multi-leader) | Concurrent writes with conflict resolution (LWW) |
| [Quorum](./quorum) | N/R/W quorum reads and writes for consistency |
