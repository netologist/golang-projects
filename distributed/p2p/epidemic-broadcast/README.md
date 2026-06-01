# Epidemic Broadcast

## Concept
Models information spread as an epidemic using the SI (Susceptible-Infected)
model. A seed node starts infected. Each round, infected nodes spread to a
random subset of susceptible peers (fanout). Dead nodes never become infected.
The epidemic converges when all live nodes are infected.

## When to Use
- You're modelling information propagation in unreliable networks.
- You need to understand convergence behaviour with node failures.

## When NOT to Use
- All nodes are always alive — gossip protocol is more practical.
- You need exact delivery guarantees — use reliable messaging.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Models real epidemic dynamics | Some messages wasted on already-infected nodes |
| Handles dead nodes gracefully | Slower convergence with high dead-node ratio |

## Go-Specific Notes
Uses the SI model (no recovery). `Spread()` iterates infected nodes and
forwards to random peers. Dead nodes are tracked and excluded from
infection and forwarding.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Epidemic protocols are self-stabilising — eventually all live nodes converge.
- Dead nodes naturally drop out; no explicit removal protocol needed.
- Higher fanout = faster convergence but more redundant messages.
