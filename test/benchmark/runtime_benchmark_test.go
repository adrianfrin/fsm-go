package benchmark_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/flandersrin/fsm-go/fsm"
	"github.com/flandersrin/fsm-go/fsmtest"
	fsmprom "github.com/flandersrin/fsm-go/observability/prometheus"
)

const benchmarkEntityCount = 100_000

func BenchmarkRuntimeFire100K(b *testing.B) {
	b.Run("without_observability", func(b *testing.B) {
		benchmarkRuntimeFire(b, false)
	})
	b.Run("with_prometheus_observability", func(b *testing.B) {
		benchmarkRuntimeFire(b, true)
	})
}

func benchmarkRuntimeFire(b *testing.B, withObservability bool) {
	b.ReportAllocs()
	ctx := context.Background()

	for n := 0; n < b.N; n++ {
		runtime, commands := newBenchmarkRuntime(b, ctx, withObservability)

		b.StartTimer()
		startedAt := time.Now()
		for _, cmd := range commands {
			if _, err := runtime.Fire(ctx, cmd); err != nil {
				b.Fatal(err)
			}
		}
		elapsed := time.Since(startedAt)
		b.StopTimer()

		b.ReportMetric(float64(benchmarkEntityCount)/elapsed.Seconds(), "transitions/s")
		b.ReportMetric(float64(elapsed.Nanoseconds())/float64(benchmarkEntityCount), "ns/transition")
	}
}

func newBenchmarkRuntime(b *testing.B, ctx context.Context, withObservability bool) (*fsm.Runtime, []fsm.FireCommand) {
	b.Helper()
	b.StopTimer()

	machine, err := fsm.Compile(benchmarkSpec())
	if err != nil {
		b.Fatal(err)
	}

	repo := fsmtest.NewMemoryRepository()
	opts := []fsm.RuntimeOption{}
	if withObservability {
		opts = append(opts, fsm.WithObserver(fsmprom.NewObserver()))
	}
	runtime := fsm.NewRuntime(repo, fsm.NewActionRegistry(), opts...)
	runtime.RegisterMachine(machine)

	commands := make([]fsm.FireCommand, benchmarkEntityCount)
	for i := range benchmarkEntityCount {
		entityID := "entity-" + strconv.Itoa(i)
		if err := runtime.CreateEntity(ctx, fsm.StateEntity{
			Machine:        "benchmark",
			MachineVersion: "v1",
			EntityID:       entityID,
			State:          "PENDING",
			Data:           map[string]any{"index": i},
		}); err != nil {
			b.Fatal(err)
		}
		commands[i] = fsm.FireCommand{
			Machine:        "benchmark",
			MachineVersion: "v1",
			EntityID:       entityID,
			Event:          "APPROVE",
			RequestID:      "req-" + strconv.Itoa(i),
			IdempotencyKey: "idem-" + strconv.Itoa(i),
			Payload:        map[string]any{"approved": true, "amount": i + 1},
		}
	}

	return runtime, commands
}

func benchmarkSpec() *fsm.MachineSpec {
	return &fsm.MachineSpec{
		Machine: "benchmark",
		Version: "v1",
		Initial: "PENDING",
		States: []fsm.StateSpec{
			{Name: "PENDING"},
			{Name: "APPROVED", Terminal: true},
		},
		Events: []fsm.EventSpec{
			{Name: "APPROVE"},
		},
		Transitions: []fsm.TransitionSpec{
			{
				Name:       "approve",
				From:       "PENDING",
				Event:      "APPROVE",
				To:         "APPROVED",
				Priority:   10,
				Guard:      "payload.approved == true && payload.amount > 0",
				Idempotent: true,
			},
		},
	}
}
