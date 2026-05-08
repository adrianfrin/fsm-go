# 可观测性

FSM Go 提供可选的 Prometheus 可观测性实现。核心库通过 `fsm.Observer` 接口暴露状态流转事件，Prometheus 实现位于 `observability/prometheus`。

## 指标

| 指标 | 类型 | 说明 |
|---|---|---|
| `fsm_transition_total` | Counter | 状态流转总次数，按状态机、版本、事件、迁移、结果区分 |
| `fsm_transition_duration_seconds` | Histogram | 状态流转耗时 |
| `fsm_transition_errors_total` | Counter | 状态流转错误次数，按错误类型区分 |
| `fsm_idempotency_hits_total` | Counter | 幂等命中次数 |
| `fsm_in_flight_transitions` | Gauge | 当前正在执行的状态流转数量 |

Prometheus 实现还会暴露 Go runtime 和进程指标。

## 在代码里启用

```go
metrics := fsmprom.NewObserver()
runtime := fsm.NewRuntime(repo, registry, fsm.WithObserver(metrics))
```

HTTP 暴露：

```go
mux.Handle("GET /metrics", promhttp.HandlerFor(metrics.Registry(), promhttp.HandlerOpts{}))
```

demo 服务已经内置 `/metrics`。

## Docker Compose

```bash
docker compose up -d --build
```

启动后可以访问：

- demo: `http://127.0.0.1:8080`
- metrics: `http://127.0.0.1:8080/metrics`
- Prometheus: `http://127.0.0.1:9090`
- Grafana: `http://127.0.0.1:3000`

Grafana 默认账号：

```text
admin / admin
```

仪表盘会自动加载到 `FSM Go` 文件夹，名称为 `FSM Go`。

## 仪表盘

内置 Grafana 仪表盘文件：

```text
observability/grafana/dashboards/fsm-go.json
```

面板包含：

- 状态流转速率
- 状态流转 p95 耗时
- 错误速率
- 正在执行的状态流转数量
- 最近一小时幂等命中次数
