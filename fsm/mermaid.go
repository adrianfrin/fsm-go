package fsm

import (
	"fmt"
	"strings"
)

func Mermaid(machine *Machine) string {
	var b strings.Builder
	b.WriteString("stateDiagram-v2\n")
	fmt.Fprintf(&b, "    [*] --> %s\n", machine.Initial)
	for _, transitions := range machine.Transitions {
		for _, transition := range transitions {
			fmt.Fprintf(&b, "    %s --> %s: %s\n", transition.From, transition.To, transition.Event)
		}
	}
	return b.String()
}
