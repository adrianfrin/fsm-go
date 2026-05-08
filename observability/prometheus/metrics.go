package prometheus

import (
	"context"

	"github.com/flandersrin/fsm-go/fsm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type Observer struct {
	registry            *prometheus.Registry
	transitionTotal     *prometheus.CounterVec
	transitionDuration  *prometheus.HistogramVec
	transitionErrors    *prometheus.CounterVec
	idempotencyHits     *prometheus.CounterVec
	inFlightTransitions *prometheus.GaugeVec
}

func NewObserver() *Observer {
	registry := prometheus.NewRegistry()
	observer := &Observer{
		registry: registry,
		transitionTotal: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "fsm_transition_total",
			Help: "Total number of FSM transitions by machine, event, transition, and status.",
		}, []string{"machine", "machine_version", "event", "from_state", "to_state", "transition", "status"}),
		transitionDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "fsm_transition_duration_seconds",
			Help:    "FSM transition execution duration in seconds.",
			Buckets: prometheus.DefBuckets,
		}, []string{"machine", "machine_version", "event", "from_state", "to_state", "transition", "status"}),
		transitionErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "fsm_transition_errors_total",
			Help: "Total number of FSM transition errors by machine, event, and error type.",
		}, []string{"machine", "machine_version", "event", "error_type"}),
		idempotencyHits: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "fsm_idempotency_hits_total",
			Help: "Total number of FSM idempotency hits.",
		}, []string{"machine", "machine_version", "event"}),
		inFlightTransitions: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "fsm_in_flight_transitions",
			Help: "Current number of in-flight FSM transitions.",
		}, []string{"machine", "machine_version", "event"}),
	}
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		observer.transitionTotal,
		observer.transitionDuration,
		observer.transitionErrors,
		observer.idempotencyHits,
		observer.inFlightTransitions,
	)
	return observer
}

func (o *Observer) Registry() *prometheus.Registry {
	return o.registry
}

func (o *Observer) TransitionStarted(_ context.Context, event fsm.TransitionStarted) {
	cmd := event.Command
	o.inFlightTransitions.WithLabelValues(cmd.Machine, cmd.MachineVersion, cmd.Event).Inc()
}

func (o *Observer) TransitionCompleted(_ context.Context, event fsm.TransitionCompleted) {
	cmd := event.Command
	result := event.Result
	fromState := ""
	toState := ""
	transition := ""
	if result != nil {
		fromState = result.FromState
		toState = result.ToState
		transition = result.TransitionName
	}

	o.inFlightTransitions.WithLabelValues(cmd.Machine, cmd.MachineVersion, cmd.Event).Dec()
	o.transitionTotal.WithLabelValues(cmd.Machine, cmd.MachineVersion, cmd.Event, fromState, toState, transition, event.Status).Inc()
	o.transitionDuration.WithLabelValues(cmd.Machine, cmd.MachineVersion, cmd.Event, fromState, toState, transition, event.Status).Observe(event.Duration.Seconds())
	if event.Err != nil {
		o.transitionErrors.WithLabelValues(cmd.Machine, cmd.MachineVersion, cmd.Event, event.ErrorType).Inc()
	}
	if result != nil && result.IdempotentHit {
		o.idempotencyHits.WithLabelValues(cmd.Machine, cmd.MachineVersion, cmd.Event).Inc()
	}
}
