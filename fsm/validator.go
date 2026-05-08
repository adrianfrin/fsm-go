package fsm

import "fmt"

func Validate(spec *MachineSpec) error {
	if spec == nil {
		return fmt.Errorf("machine spec is required")
	}
	if spec.Machine == "" {
		return fmt.Errorf("machine is required")
	}
	if spec.Version == "" {
		return fmt.Errorf("version is required")
	}
	if spec.Initial == "" {
		return fmt.Errorf("initial state is required")
	}

	states := make(map[string]StateSpec, len(spec.States))
	for _, state := range spec.States {
		if state.Name == "" {
			return fmt.Errorf("state name is required")
		}
		if _, ok := states[state.Name]; ok {
			return fmt.Errorf("duplicate state: %s", state.Name)
		}
		states[state.Name] = state
	}
	if _, ok := states[spec.Initial]; !ok {
		return fmt.Errorf("initial state not found: %s", spec.Initial)
	}

	events := make(map[string]bool, len(spec.Events))
	for _, event := range spec.Events {
		if event.Name == "" {
			return fmt.Errorf("event name is required")
		}
		if events[event.Name] {
			return fmt.Errorf("duplicate event: %s", event.Name)
		}
		events[event.Name] = true
	}

	transitionNames := make(map[string]bool, len(spec.Transitions))
	for _, transition := range spec.Transitions {
		if transition.Name == "" {
			return fmt.Errorf("transition name is required")
		}
		if transitionNames[transition.Name] {
			return fmt.Errorf("duplicate transition: %s", transition.Name)
		}
		transitionNames[transition.Name] = true
		if _, ok := states[transition.From]; !ok {
			return fmt.Errorf("transition %s from state not found: %s", transition.Name, transition.From)
		}
		if _, ok := states[transition.To]; !ok {
			return fmt.Errorf("transition %s to state not found: %s", transition.Name, transition.To)
		}
		if !events[transition.Event] {
			return fmt.Errorf("transition %s event not found: %s", transition.Name, transition.Event)
		}
		if states[transition.From].Terminal {
			return fmt.Errorf("terminal state has outgoing transition: %s", transition.From)
		}
	}

	return nil
}
