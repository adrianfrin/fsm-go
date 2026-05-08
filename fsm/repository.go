package fsm

import (
	"context"
	"time"
)

type StateEntity struct {
	Machine        string
	MachineVersion string
	EntityID       string
	State          string
	Revision       int64
	Data           map[string]any
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type StateLog struct {
	Machine        string
	MachineVersion string
	EntityID       string
	Event          string
	FromState      string
	ToState        string
	TransitionName string
	ActorID        string
	RequestID      string
	IdempotencyKey string
	Payload        map[string]any
}

type OutboxMessage struct {
	ID      int64
	Topic   string
	Key     string
	Payload map[string]any
}

type IdempotencyResult struct {
	Hit    bool
	Result *TransitionResult
}

type Repository interface {
	WithTx(ctx context.Context, fn func(context.Context, TxRepository) error) error
}

type TxRepository interface {
	CreateEntity(ctx context.Context, entity StateEntity) error
	GetEntity(ctx context.Context, machine string, entityID string) (*StateEntity, error)
	UpdateStateCAS(ctx context.Context, machine string, entityID string, fromState string, toState string, revision int64) (bool, error)
	InsertStateLog(ctx context.Context, log StateLog) error
	TryGetIdempotency(ctx context.Context, machine string, idempotencyKey string) (*IdempotencyResult, error)
	SaveIdempotencyResult(ctx context.Context, machine string, idempotencyKey string, entityID string, event string, result TransitionResult) error
	InsertOutbox(ctx context.Context, msg OutboxMessage) error
}
