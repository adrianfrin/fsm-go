# FSM Go

[![CI](https://github.com/flandersrin/fsm-go/actions/workflows/ci.yml/badge.svg)](https://github.com/flandersrin/fsm-go/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

FSM Go is a production-oriented finite state machine library for Go.

It describes states, events, transitions, guards, and actions through YAML DSL. State changes are executed through one controlled runtime entrypoint, with optional MySQL persistence, state logs, idempotency, and transactional Outbox support.

中文文档见 [README.md](README.md)。

## Use Cases

- Order lifecycle management.
- Approval workflows.
- Kafka consumer state tracking.
- Async task recovery.
- Saga workflows.
- AI Agent workflow control.

## Features

- YAML-based FSM DSL.
- Guard expressions.
- Single `Fire` entrypoint for transitions.
- Repository interface for storage isolation.
- Default MySQL Repository.
- CAS-based concurrency control.
- State transition logs.
- Idempotency result reuse.
- Transactional Outbox writes.
- Prometheus metrics.
- Grafana dashboard.
- Docker Compose demo.
- Testcontainers integration tests.

## Installation

```bash
go get github.com/flandersrin/fsm-go
```

## Quick Start

Run tests:

```bash
go test ./...
go test -race ./...
go test -count=1 -tags=integration ./test/integration/...
go test -run '^$' -bench BenchmarkRuntimeFire100K -benchtime=1x -benchmem ./test/benchmark
```

If Taskfile is installed:

```bash
task check
```

If you do not want to install Taskfile globally:

```bash
go run github.com/go-task/task/v3/cmd/task@v3.50.0 check
```

Run examples:

```bash
go run ./examples/order
go run ./examples/kafka_message
go run ./examples/agent_run
```

The order example prints:

```text
PENDING -> PAID
logs=1 outbox=1
```

See [English Quick Start](docs/en/quickstart.md) for details.

## Library Usage

Load and compile a DSL file:

```go
spec, err := fsm.LoadYAML("configs/order.v1.yaml")
if err != nil {
    return err
}

machine, err := fsm.Compile(spec)
if err != nil {
    return err
}
```

Register the runtime:

```go
runtime := fsm.NewRuntime(repository, actionRegistry)
runtime.RegisterMachine(machine)
```

Fire a transition:

```go
result, err := runtime.Fire(ctx, fsm.FireCommand{
    Machine:        "order",
    MachineVersion: "v1",
    EntityID:       "order-10001",
    Event:          "PAY_SUCCESS",
    Actor:          fsm.Actor{ID: "user-1", Role: "customer"},
    RequestID:      "req-1",
    IdempotencyKey: "pay-10001",
    Payload: map[string]any{
        "paymentStatus": "SUCCESS",
        "amount":        100,
    },
})
```

See [Library Usage](docs/en/library-usage.md) for a fuller example.

## Docker Demo

Start the demo service and MySQL:

```bash
docker compose up -d --build
```

Initialize an order:

```bash
curl -X POST http://127.0.0.1:8080/demo/order/init \
  -H 'Content-Type: application/json' \
  -d '{"entity_id":"order-10001","data":{}}'
```

Fire a successful payment event:

```bash
curl -X POST http://127.0.0.1:8080/demo/order/fire \
  -H 'Content-Type: application/json' \
  -d '{
    "entity_id":"order-10001",
    "event":"PAY_SUCCESS",
    "actor_id":"user-1",
    "actor_role":"customer",
    "request_id":"req-1",
    "idempotency_key":"pay-10001",
    "payload":{"paymentStatus":"SUCCESS","amount":100}
  }'
```

Inspect state, logs, and Outbox:

```bash
curl http://127.0.0.1:8080/demo/order/order-10001
curl http://127.0.0.1:8080/demo/order/order-10001/logs
curl http://127.0.0.1:8080/demo/outbox
```

Clean up:

```bash
docker compose down -v
```

## Observability

The demo service exposes Prometheus metrics:

```bash
curl http://127.0.0.1:8080/metrics
```

Docker Compose also starts Prometheus and Grafana:

```text
Prometheus: http://127.0.0.1:9090
Grafana:    http://127.0.0.1:3000
```

Default Grafana login is `admin / admin`, and the built-in dashboard is provisioned automatically.

See [Observability](docs/en/observability.md) for details.

## Docker Package

Release tags publish the demo image to GitHub Container Registry:

```text
ghcr.io/flandersrin/fsm-go
```

Example:

```bash
docker pull ghcr.io/flandersrin/fsm-go:v0.1.0
```

## Documentation

- [Architecture](docs/en/architecture.md)
- [Quick Start](docs/en/quickstart.md)
- [Library Usage](docs/en/library-usage.md)
- [Docker Demo](docs/en/docker-demo.md)
- [Observability](docs/en/observability.md)
- [Testing](docs/en/testing.md)
- [Benchmark](docs/en/benchmark.md)

## Project Layout

```text
fsm/                    Core FSM library
actions/                Reusable actions
persistence/mysql/      MySQL Repository
fsmtest/                In-memory Repository for tests and examples
observability/          Prometheus and Grafana setup
configs/                Example DSL files
examples/order/         Order FSM example
examples/kafka_message/ Kafka consumer FSM example
examples/agent_run/     Agent run FSM example
cmd/fsm-demo/           Runnable demo service
test/integration/       Testcontainers integration tests
test/benchmark/         100K transition benchmarks
docs/                   Chinese and English documentation
```

## Contributing

Issues and pull requests are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) and [CONTRIBUTORS.md](CONTRIBUTORS.md).

Security reports should follow [SECURITY.md](SECURITY.md).

## License

FSM Go is released under the [MIT License](LICENSE).
