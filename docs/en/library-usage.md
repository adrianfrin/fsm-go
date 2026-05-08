# Library Usage

FSM Go is designed as a Go library. Business services load DSL files, register machines, provide a Repository, and trigger transitions through one runtime entrypoint.

## 1. Define DSL

```yaml
machine: order
version: v1
initial: PENDING

states:
  - name: PENDING
  - name: PAID
  - name: COMPLETED
    terminal: true

events:
  - name: PAY_SUCCESS
  - name: FINISH

transitions:
  - name: pay_success
    from: PENDING
    event: PAY_SUCCESS
    to: PAID
    priority: 10
    guard: "payload.paymentStatus == 'SUCCESS' && payload.amount > 0"
    idempotent: true
    actions:
      in_tx:
        - outbox.order_paid
```

## 2. Load and Register the Machine

```go
spec, err := fsm.LoadYAML("configs/order.v1.yaml")
if err != nil {
    return err
}

machine, err := fsm.Compile(spec)
if err != nil {
    return err
}

runtime := fsm.NewRuntime(repository, actionRegistry)
runtime.RegisterMachine(machine)
```

## 3. Initialize Entity State

```go
err := runtime.CreateEntity(ctx, fsm.StateEntity{
    Machine:        "order",
    MachineVersion: "v1",
    EntityID:       "order-10001",
    State:          "PENDING",
    Data:           map[string]any{},
})
```

## 4. Fire a Transition

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

When it succeeds, the state changes from `PENDING` to `PAID`. The runtime also writes a state log, an idempotency result, and an Outbox message when configured.

## 5. Use MySQL

```go
db, err := sql.Open("mysql", dsn)
if err != nil {
    return err
}

repository := mysqlrepo.NewRepository(db)
if err := repository.InitSchema(ctx); err != nil {
    return err
}
```

In production, prefer running the SQL schema through your migration tool instead of initializing schema at application startup.

## 6. Register Actions

```go
registry := fsm.NewActionRegistry()

actions.RegisterOutbox(registry, map[string]string{
    "outbox.order_paid": "order.paid",
})
```

Business services can register their own actions for audit records, notifications, or domain-specific writes.

## 7. More Examples

```bash
go run ./examples/order
go run ./examples/kafka_message
go run ./examples/agent_run
```

The corresponding DSL files are:

- `configs/order.v1.yaml`
- `configs/kafka-message.v1.yaml`
- `configs/agent-run.v1.yaml`
