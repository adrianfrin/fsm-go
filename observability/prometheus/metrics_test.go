package prometheus

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/flandersrin/fsm-go/fsm"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestObserverRecordsTransitionMetrics(t *testing.T) {
	observer := NewObserver()
	ctx := context.Background()
	cmd := fsm.FireCommand{
		Machine:        "order",
		MachineVersion: "v1",
		EntityID:       "order-1",
		Event:          "PAY_SUCCESS",
	}
	result := &fsm.TransitionResult{
		Machine:        "order",
		MachineVersion: "v1",
		EntityID:       "order-1",
		Event:          "PAY_SUCCESS",
		FromState:      "PENDING",
		ToState:        "PAID",
		TransitionName: "pay_success",
	}

	observer.TransitionStarted(ctx, fsm.TransitionStarted{Command: cmd})
	if got := testutil.ToFloat64(observer.inFlightTransitions.WithLabelValues("order", "v1", "PAY_SUCCESS")); got != 1 {
		t.Fatalf("in-flight transitions before completion = %v, want 1", got)
	}

	observer.TransitionCompleted(ctx, fsm.TransitionCompleted{
		Command:  cmd,
		Result:   result,
		Duration: 10 * time.Millisecond,
		Status:   "success",
	})

	if got := testutil.ToFloat64(observer.inFlightTransitions.WithLabelValues("order", "v1", "PAY_SUCCESS")); got != 0 {
		t.Fatalf("in-flight transitions after completion = %v, want 0", got)
	}
	if got := testutil.ToFloat64(observer.transitionTotal.WithLabelValues("order", "v1", "PAY_SUCCESS", "PENDING", "PAID", "pay_success", "success")); got != 1 {
		t.Fatalf("transition total = %v, want 1", got)
	}
}

func TestObserverRecordsErrorAndIdempotencyMetrics(t *testing.T) {
	observer := NewObserver()
	ctx := context.Background()
	cmd := fsm.FireCommand{
		Machine:        "order",
		MachineVersion: "v1",
		EntityID:       "order-1",
		Event:          "PAY_SUCCESS",
	}

	observer.TransitionStarted(ctx, fsm.TransitionStarted{Command: cmd})
	observer.TransitionCompleted(ctx, fsm.TransitionCompleted{
		Command:   cmd,
		Err:       errors.New("boom"),
		Duration:  5 * time.Millisecond,
		ErrorType: "errorString",
		Status:    "error",
	})

	if got := testutil.ToFloat64(observer.transitionErrors.WithLabelValues("order", "v1", "PAY_SUCCESS", "errorString")); got != 1 {
		t.Fatalf("transition errors = %v, want 1", got)
	}

	observer.TransitionStarted(ctx, fsm.TransitionStarted{Command: cmd})
	observer.TransitionCompleted(ctx, fsm.TransitionCompleted{
		Command: cmd,
		Result: &fsm.TransitionResult{
			Machine:        "order",
			MachineVersion: "v1",
			EntityID:       "order-1",
			Event:          "PAY_SUCCESS",
			FromState:      "PAID",
			ToState:        "PAID",
			IdempotentHit:  true,
		},
		Duration: 3 * time.Millisecond,
		Status:   "idempotent_hit",
	})

	if got := testutil.ToFloat64(observer.idempotencyHits.WithLabelValues("order", "v1", "PAY_SUCCESS")); got != 1 {
		t.Fatalf("idempotency hits = %v, want 1", got)
	}
}
