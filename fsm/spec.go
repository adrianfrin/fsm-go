package fsm

type MachineSpec struct {
	Machine     string           `yaml:"machine"`
	Version     string           `yaml:"version"`
	Initial     string           `yaml:"initial"`
	States      []StateSpec      `yaml:"states"`
	Events      []EventSpec      `yaml:"events"`
	Transitions []TransitionSpec `yaml:"transitions"`
}

type StateSpec struct {
	Name     string `yaml:"name"`
	Terminal bool   `yaml:"terminal"`
}

type EventSpec struct {
	Name string `yaml:"name"`
}

type TransitionSpec struct {
	Name       string     `yaml:"name"`
	From       string     `yaml:"from"`
	Event      string     `yaml:"event"`
	To         string     `yaml:"to"`
	Priority   int        `yaml:"priority"`
	Guard      string     `yaml:"guard"`
	Idempotent bool       `yaml:"idempotent"`
	Actions    ActionSpec `yaml:"actions"`
}

type ActionSpec struct {
	InTx        []string `yaml:"in_tx"`
	AfterCommit []string `yaml:"after_commit"`
}
