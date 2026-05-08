package fsm

import (
	"context"
	"fmt"
)

type ActionContext struct {
	Command    FireCommand
	Entity     *StateEntity
	Transition CompiledTransition
	Tx         TxRepository
	Result     *TransitionResult
}

type ActionFunc func(context.Context, ActionContext) error

type ActionRegistry struct {
	actions map[string]ActionFunc
}

func NewActionRegistry() *ActionRegistry {
	return &ActionRegistry{actions: map[string]ActionFunc{}}
}

func (r *ActionRegistry) Register(name string, fn ActionFunc) {
	r.actions[name] = fn
}

func (r *ActionRegistry) Run(ctx context.Context, name string, ac ActionContext) error {
	fn, ok := r.actions[name]
	if !ok {
		return fmt.Errorf("action not found: %s", name)
	}
	if err := fn(ctx, ac); err != nil {
		return fmt.Errorf("run action %s: %w", name, err)
	}
	return nil
}
