# Error Handling

Idiomatic Go error handling patterns that go beyond basic `if err != nil`.

| Pattern | Description |
|---------|-------------|
| [Retry](./retry) | Retry with exponential backoff, jitter, and retryable classification |
| [Result Type](./result-type) | Functional Result[T] with Map, FlatMap, and Or combinators |
| [Sentinel Errors](./sentinel-errors) | Predefined error values for structured error comparison |
| [Error Wrapping](./error-wrapping) | Rich error types with cause chaining, JSON serialization, and `errors.As` support |
