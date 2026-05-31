# Decorator

## Concept
Wrap an interface implementation to add behaviour without changing the original.
In Go, the decorator is a struct that holds the wrapped value and delegates all
calls, intercepting where needed.

## When to Use
- Add cross-cutting concerns (logging, metrics, caching) to existing interfaces.
- Compose multiple concerns independently.
- Avoid subclassing (which Go does not have).

## When NOT to Use
- Behaviour is core to the type — put it in the type.
- The interface is large (>5 methods) — wrapping becomes tedious; use middleware.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Open/Closed: add behaviour without modifying original | Every method must be forwarded |
| Composable: stack decorators | Order matters — wrong order causes bugs |

## Go-Specific Notes
Because Go has implicit interface satisfaction, any struct with the right method
set is a decorator — no annotation needed.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- A decorator holds the wrapped interface and adds behaviour on one or more methods.
- Name decorators descriptively: `LoggingStore`, `CachingStore`.
- Compose from innermost (real impl) outward.
