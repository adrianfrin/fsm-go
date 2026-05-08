package fsm

import (
	"context"
	"reflect"
	"strings"
	"time"
)

type Observer interface {
	TransitionStarted(context.Context, TransitionStarted)
	TransitionCompleted(context.Context, TransitionCompleted)
}

type TransitionStarted struct {
	Command FireCommand
}

type TransitionCompleted struct {
	Command   FireCommand
	Result    *TransitionResult
	Err       error
	Duration  time.Duration
	ErrorType string
	Status    string
}

type RuntimeOption func(*Runtime)

func WithObserver(observer Observer) RuntimeOption {
	return func(runtime *Runtime) {
		if observer != nil {
			runtime.observers = append(runtime.observers, observer)
		}
	}
}

func (r *Runtime) observeTransitionStarted(ctx context.Context, cmd FireCommand) {
	event := TransitionStarted{Command: cmd}
	for _, observer := range r.observers {
		observer.TransitionStarted(ctx, event)
	}
}

func (r *Runtime) observeTransitionCompleted(ctx context.Context, cmd FireCommand, result *TransitionResult, err error, duration time.Duration) {
	event := TransitionCompleted{
		Command:   cmd,
		Result:    result,
		Err:       err,
		Duration:  duration,
		ErrorType: errorType(err),
		Status:    transitionStatus(result, err),
	}
	for _, observer := range r.observers {
		observer.TransitionCompleted(ctx, event)
	}
}

func transitionStatus(result *TransitionResult, err error) string {
	if err != nil {
		return "error"
	}
	if result != nil && result.IdempotentHit {
		return "idempotent_hit"
	}
	return "success"
}

func errorType(err error) string {
	if err == nil {
		return ""
	}
	t := reflect.TypeOf(err)
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	name := t.String()
	if idx := strings.LastIndex(name, "."); idx >= 0 {
		return name[idx+1:]
	}
	return name
}
