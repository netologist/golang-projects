# Factory

## Concept
Provide a creation interface that hides which concrete type is instantiated.
The registry pattern (a `map[string]func() Interface`) is the Go-idiomatic
approach — it avoids switch statements and is open to extension.

## When to Use
- Multiple implementations of an interface, selected at runtime.
- Plugin-style extensibility — third-party code can register its own factories.
- Decouple callers from concrete types.

## When NOT to Use
- Only one implementation exists — a plain constructor is clearer.
- The selection logic is trivially simple (one if/else).

## Trade-offs
| Benefit | Cost |
|---------|------|
| Open/Closed: add types without changing factory | Unregistered names fail at runtime |
| Decouples caller from implementation | Type lookup errors surface at runtime |

## Go-Specific Notes
Register factories in `init()` in each codec's file — callers import the
package for its side effect to make the codec available.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Use `map[string]func() Interface` as the registry.
- `init()` side-effect registration keeps factory and implementation together.
- Always return `(Interface, error)` from `New(name)` — unknown name is an error.
