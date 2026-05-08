package main

import (
	"context"
	"fmt"
	"log"

	"github.com/flandersrin/fsm-go/actions"
	"github.com/flandersrin/fsm-go/fsm"
	"github.com/flandersrin/fsm-go/fsmtest"
)

func main() {
	ctx := context.Background()

	spec, err := fsm.LoadYAML("configs/kafka-message.v1.yaml")
	if err != nil {
		log.Fatal(err)
	}
	machine, err := fsm.Compile(spec)
	if err != nil {
		log.Fatal(err)
	}

	repo := fsmtest.NewMemoryRepository()
	registry := fsm.NewActionRegistry()
	actions.RegisterOutbox(registry, map[string]string{
		"outbox.kafka_retry":       "kafka.retry",
		"outbox.kafka_dead_letter": "kafka.dead_letter",
	})

	runtime := fsm.NewRuntime(repo, registry)
	runtime.RegisterMachine(machine)

	if err := runtime.CreateEntity(ctx, fsm.StateEntity{
		Machine:        "kafka_message",
		MachineVersion: "v1",
		EntityID:       "message-example-1",
		State:          "PENDING",
		Data:           map[string]any{"retryCount": 0},
	}); err != nil {
		log.Fatal(err)
	}

	mustFire(ctx, runtime, fsm.FireCommand{
		Machine:        "kafka_message",
		MachineVersion: "v1",
		EntityID:       "message-example-1",
		Event:          "START",
	})
	result := mustFire(ctx, runtime, fsm.FireCommand{
		Machine:        "kafka_message",
		MachineVersion: "v1",
		EntityID:       "message-example-1",
		Event:          "FAIL",
	})

	fmt.Printf("%s -> %s\n", result.FromState, result.ToState)
	fmt.Printf("logs=%d outbox=%d\n", len(repo.Logs()), len(repo.Outbox()))
}

func mustFire(ctx context.Context, runtime *fsm.Runtime, cmd fsm.FireCommand) *fsm.TransitionResult {
	result, err := runtime.Fire(ctx, cmd)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
