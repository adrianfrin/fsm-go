# Testing

FSM Go covers both core logic and real dependencies.

## Unit Tests

```bash
go test ./...
```

Covered areas:

- DSL validation.
- Guard evaluation.
- State transitions.
- Idempotency hits.
- Action execution.
- State log and Outbox writes.

## Race Tests

```bash
go test -race ./...
```

This checks the runtime and test helpers for data races.

## Integration Tests

```bash
go test -count=1 -tags=integration ./test/integration/...
```

Integration tests use Testcontainers to start a real MySQL instance.

They cover:

- MySQL schema initialization.
- CAS state updates.
- Only one concurrent transition winner.
- State log writes.
- Idempotency result reuse.
- Transactional Outbox writes.

## Benchmark

```bash
go test -run '^$' -bench BenchmarkRuntimeFire100K -benchtime=1x -benchmem ./test/benchmark
```

The benchmark runs 100,000 transitions and compares latency and memory allocation with and without Prometheus observability.

See [Benchmark](benchmark.md) for details.

## Local Check

If Taskfile is installed:

```bash
task check
```

This runs formatting checks, module tidy checks, `go vet`, `golangci-lint`, unit tests, race tests, Testcontainers integration tests, and the 100K transition benchmark.

If you do not want to install Taskfile globally:

```bash
go run github.com/go-task/task/v3/cmd/task@v3.50.0 check
```
