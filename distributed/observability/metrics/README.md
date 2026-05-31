# Metrics

## Concept
Instrument services with the RED method — Rate, Errors, Duration — using
Prometheus. A counter tracks request volume and errors; a histogram tracks
latency; a gauge tracks in-flight requests.

## When to Use
- Any production service needing dashboards and alerting.
- SLO tracking (latency percentiles, error rates).

## When NOT to Use
- Throwaway tools with no monitoring needs.

## Trade-offs
| Benefit | Cost |
|---------|------|
| Standard, queryable metrics | Cardinality must be controlled |
| Works with Grafana/Alertmanager | Histograms add memory per label set |

## Go-Specific Notes
The `Recorder` bundles RED instruments registered against a local registry. A
middleware auto-instruments any `http.Handler`.

## Running
```bash
go run ./simple
go run ./advanced
go test ./advanced/... -v
```

## Key Takeaways
- RED: Rate, Errors, Duration covers most service health needs.
- Keep label cardinality low (method, route, status class).
- Use a local registry in libraries/tests to avoid global pollution.
