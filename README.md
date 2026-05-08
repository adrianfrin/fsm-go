# FSM Go

[![CI](https://github.com/flandersrin/fsm-go/actions/workflows/ci.yml/badge.svg)](https://github.com/flandersrin/fsm-go/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

FSM Go 是一个面向生产环境的 Go 状态机库。

它用 DSL 描述状态、事件和流转规则，通过统一入口完成状态变化，并提供 MySQL、状态日志、幂等和 Outbox 示例实现。

English documentation: [README.en.md](README.en.md)

## 适用场景

- 订单状态流转。
- 审批流。
- Kafka 消费状态治理。
- 异步任务恢复。
- Saga 流程。
- AI Agent Workflow。

## 核心能力

- YAML DSL 定义状态机。
- Guard 条件判断。
- 统一 `Fire` 入口。
- Repository 接口隔离存储实现。
- 默认 MySQL Repository。
- CAS 并发控制。
- 状态迁移日志。
- 幂等请求复用历史结果。
- Outbox 事务写入。
- Prometheus 指标。
- Grafana 仪表盘。
- Docker Compose demo。
- Testcontainers 集成测试。

## 安装

```bash
go get github.com/flandersrin/fsm-go
```

## 快速开始

运行单元测试：

```bash
go test ./...
go test -race ./...
```

运行 Go 示例：

```bash
go run ./examples/order
go run ./examples/kafka_message
go run ./examples/agent_run
```

订单示例预期输出类似：

```text
PENDING -> PAID
logs=1 outbox=1
```

更多说明见 [快速开始](docs/quickstart.md)。

## 作为库使用

加载 DSL：

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

注册运行时：

```go
runtime := fsm.NewRuntime(repository, actionRegistry)
runtime.RegisterMachine(machine)
```

触发状态流转：

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

完整接入说明见 [库接入示例](docs/library-usage.md)。

## Docker Demo

启动 demo 服务和 MySQL：

```bash
docker compose up -d --build
```

初始化订单：

```bash
curl -X POST http://127.0.0.1:8080/demo/order/init \
  -H 'Content-Type: application/json' \
  -d '{"entity_id":"order-10001","data":{}}'
```

触发支付成功：

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

查看状态、日志和 Outbox：

```bash
curl http://127.0.0.1:8080/demo/order/order-10001
curl http://127.0.0.1:8080/demo/order/order-10001/logs
curl http://127.0.0.1:8080/demo/outbox
```

清理环境：

```bash
docker compose down -v
```

更多说明见 [Docker Demo](docs/docker-demo.md)。

## 可观测性

demo 服务暴露 Prometheus 指标：

```bash
curl http://127.0.0.1:8080/metrics
```

Docker Compose 会同时启动 Prometheus 和 Grafana：

```text
Prometheus: http://127.0.0.1:9090
Grafana:    http://127.0.0.1:3000
```

Grafana 默认账号是 `admin / admin`，内置仪表盘会自动加载。

更多说明见 [可观测性](docs/observability.md)。

## Docker Package

发布版本会把 demo 镜像推送到 GitHub Container Registry：

```text
ghcr.io/flandersrin/fsm-go
```

拉取示例：

```bash
docker pull ghcr.io/flandersrin/fsm-go:v0.1.0
```

## 测试

```bash
go test ./...
go test -race ./...
go test -count=1 -tags=integration ./test/integration/...
go test -run '^$' -bench BenchmarkRuntimeFire100K -benchtime=1x -benchmem ./test/benchmark
```

集成测试使用 Testcontainers 启动真实 MySQL。

如果已安装 Taskfile：

```bash
task check
task test:integration
```

如果不想全局安装 Taskfile，也可以直接运行：

```bash
go run github.com/go-task/task/v3/cmd/task@v3.50.0 check
```

更多说明见 [测试说明](docs/testing.md) 和 [Benchmark](docs/benchmark.md)。

## 项目结构

```text
fsm/                  核心状态机库
actions/              可复用 Action
persistence/mysql/    MySQL Repository
fsmtest/              测试和示例用内存 Repository
observability/        Prometheus 和 Grafana 配置
configs/              示例 DSL
examples/order/       订单状态机示例
examples/kafka_message/ Kafka 消费状态机示例
examples/agent_run/   Agent Run 状态机示例
cmd/fsm-demo/         可运行 demo 服务
test/integration/     Testcontainers 集成测试
test/benchmark/       10 万级状态流转 Benchmark
docs/                 使用和架构文档
```

更多说明见 [架构说明](docs/architecture.md)。

## 贡献

欢迎提交 issue 和 pull request。提交前请阅读 [贡献指南](CONTRIBUTING.md)。

贡献者列表见 [CONTRIBUTORS.md](CONTRIBUTORS.md)。

安全问题请参考 [安全说明](SECURITY.md)。

## 许可证

本项目使用 [MIT License](LICENSE)。
