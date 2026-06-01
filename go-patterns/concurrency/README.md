# Concurrency

Go-idiomatic concurrency patterns leveraging goroutines and channels.

| Pattern | Description |
|---------|-------------|
| [Worker Pool](./worker-pool) | Fixed-size pool of goroutines processing from a job channel |
| [Context Propagation](./context-propagation) | Passing context, request IDs, and values through the call chain |
| [Errgroup](./errgroup) | Group of goroutines with synchronised error collection |
| [Semaphore](./semaphore) | Weighted semaphore for bounded concurrency |
| [Rate Limiter](./rate-limiter) | Token bucket and sliding window rate limiters |
