# Docker Demo

Docker Demo 用于让用户快速看到状态机运行效果。

它会启动：

- MySQL 8.4
- `fsm-demo` HTTP 服务

## 启动

```bash
docker compose up -d --build
```

## API

| API | 说明 |
|---|---|
| `GET /healthz` | 健康检查 |
| `POST /demo/order/init` | 初始化订单 |
| `POST /demo/order/fire` | 触发订单状态流转 |
| `GET /demo/order/{id}` | 查询订单当前状态 |
| `GET /demo/order/{id}/logs` | 查询订单状态轨迹 |
| `GET /demo/outbox` | 查询 Outbox 消息 |

## 示例

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

## 清理

```bash
docker compose down -v
```
