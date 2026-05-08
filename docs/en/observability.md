# Observability

FSM Go provides an optional Prometheus observability implementation. The core library exposes transition events through the `fsm.Observer` interface, and the Prometheus implementation lives in `observability/prometheus`.

## Metrics

| Metric | Type | Description |
|---|---|---|
| `fsm_transition_total` | Counter | Total number of transitions by machine, version, event, transition, and status |
| `fsm_transition_duration_seconds` | Histogram | Transition execution duration |
| `fsm_transition_errors_total` | Counter | Transition errors by error type |
| `fsm_idempotency_hits_total` | Counter | Idempotency hit count |
| `fsm_in_flight_transitions` | Gauge | Current in-flight transition count |

The Prometheus implementation also exposes Go runtime and process metrics.

## Enable in Code

```go
metrics := fsmprom.NewObserver()
runtime := fsm.NewRuntime(repo, registry, fsm.WithObserver(metrics))
```

Expose metrics over HTTP:

```go
mux.Handle("GET /metrics", promhttp.HandlerFor(metrics.Registry(), promhttp.HandlerOpts{}))
```

The demo service already exposes `/metrics`.

## Docker Compose

```bash
docker compose up -d --build
```

Available endpoints:

- demo: `http://127.0.0.1:8080`
- metrics: `http://127.0.0.1:8080/metrics`
- Prometheus: `http://127.0.0.1:9090`
- Grafana: `http://127.0.0.1:3000`

Default Grafana login:

```text
admin / admin
```

The dashboard is provisioned automatically under the `FSM Go` folder.

## Dashboard

Built-in Grafana dashboard file:

```text
observability/grafana/dashboards/fsm-go.json
```

Panels include:

- Transition rate
- Transition p95 duration
- Error rate
- In-flight transitions
- Idempotency hits in the last hour
