# Benchmark

FSM Go includes a 100K transition benchmark for measuring runtime latency, throughput, and memory allocation under batch workloads.

The benchmark covers two scenarios:

- `without_observability`: core runtime transitions only.
- `with_prometheus_observability`: the same workload with Prometheus observability enabled.

## Command

```bash
go test -run '^$' -bench BenchmarkRuntimeFire100K -benchtime=1x -benchmem ./test/benchmark
```

With Taskfile:

```bash
task test:benchmark
```

Without installing Taskfile globally:

```bash
go run github.com/go-task/task/v3/cmd/task@v3.50.0 test:benchmark
```

## What It Measures

Each benchmark run preloads 100,000 state entities, then fires one transition for each entity.

Every transition goes through:

- Compiled transition rule matching.
- Guard expression evaluation.
- CAS state update.
- State log write.
- Idempotency result save.

With observability enabled, it also records:

- Transition count.
- Transition duration.
- Error count.
- Idempotency hit count.
- In-flight transition count.

## Local Sample Result

The following result was measured on an Apple M1 Pro. It is a local sample only. Actual results depend on hardware, Go version, and runtime load.

| Scenario | Total Time | Per Transition | Throughput | Allocation |
|---|---:|---:|---:|---:|
| Without observability | 561.80 ms | 5,618 ns | 178,005 transitions/s | 699.27 MB |
| With Prometheus observability | 601.15 ms | 6,011 ns | 166,354 transitions/s | 699.28 MB |

In this run, Prometheus observability increased per-transition latency by about 7.0%. The overhead mainly comes from metric counting, label writes, and duration recording.

## Guidance

Use the benchmark to track trends, not as a fixed performance guarantee. For production use, rerun it with your own DSL complexity, actions, storage implementation, database latency, and concurrency model.
