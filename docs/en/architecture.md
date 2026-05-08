# Architecture

## Goal

FSM Go is not just a status enum helper. It is a state governance library.

The goal is to move state changes out of scattered business code and centralize:

- Transition rules.
- Invalid transition rejection.
- Guard evaluation.
- Concurrency control.
- Idempotency.
- State logs.
- Outbox consistency.

## Core Flow

```text
Business Service
  -> Runtime.Fire
  -> Machine Registry
  -> Compiled DSL
  -> Guard Engine
  -> Repository Transaction
      -> CAS state update
      -> insert state log
      -> insert Outbox
      -> save idempotency result
```

## Modules

| Module | Responsibility |
|---|---|
| `fsm` | Core FSM library, DSL, Runtime, Repository interface |
| `actions` | Reusable actions |
| `persistence/mysql` | MySQL Repository implementation |
| `fsmtest` | In-memory Repository for tests and examples |
| `observability/prometheus` | Prometheus metrics implementation |
| `cmd/fsm-demo` | Runnable demo service |
| `test/integration` | Testcontainers integration tests |

## Built-in Examples

| Example | DSL | Purpose |
|---|---|---|
| `examples/order` | `configs/order.v1.yaml` | Payment, shipping, completion, cancellation |
| `examples/kafka_message` | `configs/kafka-message.v1.yaml` | Kafka processing, retry, dead letter |
| `examples/agent_run` | `configs/agent-run.v1.yaml` | Agent planning, running, tool waiting, completion |

## Storage Boundary

The core library only depends on the Repository interface.

The default MySQL implementation provides four tables:

- `fsm_entity`
- `fsm_state_log`
- `fsm_idempotency`
- `fsm_outbox`

The default schema does not include `tenant_id` or `sub_tenant_id`. Tenant isolation should be implemented through a custom Repository or plugin.
