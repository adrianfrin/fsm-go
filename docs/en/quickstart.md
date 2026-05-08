# Quick Start

This guide shows how to run FSM Go locally and execute a state transition.

## 1. Run Tests

```bash
go test ./...
go test -race ./...
```

If Docker is available, run the integration tests:

```bash
go test -count=1 -tags=integration ./test/integration/...
```

Integration tests use Testcontainers to start a real MySQL instance and clean it up after the tests finish.

## 2. Run Examples

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

The Kafka example demonstrates a message moving from `RUNNING` to `RETRY`.

The Agent example demonstrates a run moving through planning, execution, tool waiting, and completion.

## 3. Run Docker Demo

```bash
docker compose up -d --build
```

Health check:

```bash
curl http://127.0.0.1:8080/healthz
```

Initialize an order:

```bash
curl -X POST http://127.0.0.1:8080/demo/order/init \
  -H 'Content-Type: application/json' \
  -d '{"entity_id":"order-10001","data":{}}'
```

Fire a payment success event:

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
