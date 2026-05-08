package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/flandersrin/fsm-go/fsm"
	"github.com/flandersrin/fsm-go/fsmtest"
)

func main() {
	ctx := context.Background()

	spec, err := fsm.LoadYAML("configs/agent-run.v1.yaml")
	if err != nil {
		slog.Error("load dsl", "error", err)
		os.Exit(1)
	}
	machine, err := fsm.Compile(spec)
	if err != nil {
		slog.Error("compile dsl", "error", err)
		os.Exit(1)
	}

	repo := fsmtest.NewMemoryRepository()
	runtime := fsm.NewRuntime(repo, fsm.NewActionRegistry())
	runtime.RegisterMachine(machine)

	if err := runtime.CreateEntity(ctx, fsm.StateEntity{
		Machine:        "agent_run",
		MachineVersion: "v1",
		EntityID:       "agent-run-example-1",
		State:          "DRAFT",
		Data:           map[string]any{"retryCount": 0},
	}); err != nil {
		slog.Error("create entity", "error", err)
		os.Exit(1)
	}

	events := []string{"SUBMIT", "PLAN_DONE", "NEED_TOOL", "TOOL_DONE", "FINISH"}
	var result *fsm.TransitionResult
	for _, event := range events {
		result = mustFire(ctx, runtime, fsm.FireCommand{
			Machine:        "agent_run",
			MachineVersion: "v1",
			EntityID:       "agent-run-example-1",
			Event:          event,
		})
	}

	fmt.Printf("%s -> %s\n", result.FromState, result.ToState)
	fmt.Printf("logs=%d outbox=%d\n", len(repo.Logs()), len(repo.Outbox()))
}

func mustFire(ctx context.Context, runtime *fsm.Runtime, cmd fsm.FireCommand) *fsm.TransitionResult {
	result, err := runtime.Fire(ctx, cmd)
	if err != nil {
		slog.Error("fire transition", "error", err)
		os.Exit(1)
	}
	return result
}
