package fsm

import "fmt"

type ErrInvalidTransition struct {
	State string
	Event string
}

func (e ErrInvalidTransition) Error() string {
	return fmt.Sprintf("invalid transition: state=%s event=%s", e.State, e.Event)
}

type ErrGuardRejected struct {
	State string
	Event string
}

func (e ErrGuardRejected) Error() string {
	return fmt.Sprintf("guard rejected: state=%s event=%s", e.State, e.Event)
}

type ErrConcurrentTransition struct {
	EntityID string
}

func (e ErrConcurrentTransition) Error() string {
	return fmt.Sprintf("concurrent transition conflict: entity_id=%s", e.EntityID)
}

type ErrTerminalState struct {
	State string
}

func (e ErrTerminalState) Error() string {
	return fmt.Sprintf("terminal state cannot transition: state=%s", e.State)
}
