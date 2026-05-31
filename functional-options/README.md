# Functional Options

## Concept
Instead of a config struct or multiple constructors, accept variadic `func(*T)`
options. Each option mutates one field. The caller composes only what they need.

## When to Use
- Three or more optional parameters.
- The zero value of a field is not a safe default.
- You want to add options later without breaking existing callers.

## When NOT to Use
- All parameters are required — use a plain constructor.
- Only one or two simple fields — a struct literal is clearer.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Extensible without breaking changes | Slightly more verbose per call |
| Self-documenting option names | Options applied silently unless validated |
| Composable and reusable | Harder to inspect/serialise config |

## Go-Specific Notes
Options are ordinary `func(*config) error` values — first-class and composable.
Validation lives inside `New()` after all options are applied, keeping each
option function trivial.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- Option funcs are ordinary values — store, compose, and pass them around.
- Validate after applying all options, not inside each option.
- Name options `With<Field>` by convention.
- Use an unexported `config` struct to separate defaults from overrides.
