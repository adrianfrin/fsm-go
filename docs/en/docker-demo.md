# Docker Demo

The Docker demo lets users quickly see FSM Go running with real MySQL.

It starts:

- MySQL 8.4
- `fsm-demo` HTTP service
- Prometheus
- Grafana

## Start

```bash
docker compose up -d --build
```

## API

| API | Description |
|---|---|
| `GET /healthz` | Health check |
| `POST /demo/order/init` | Initialize an order |
| `POST /demo/order/fire` | Fire an order transition |
| `GET /demo/order/{id}` | Read current order state |
| `GET /demo/order/{id}/logs` | Read order state history |
| `GET /demo/outbox` | Read Outbox messages |
| `GET /metrics` | Prometheus metrics |

## Example

```bash
curl -X POST http://127.0.0.1:8080/demo/order/init \
  -H 'Content-Type: application/json' \
  -d '{"entity_id":"order-10001","data":{}}'
```

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

## Clean Up

```bash
docker compose down -v
```

## Observability Endpoints

- Prometheus: `http://127.0.0.1:9090`
- Grafana: `http://127.0.0.1:3000`

Default Grafana login:

```text
admin / admin
```
