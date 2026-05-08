package fsm

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadYAML(path string) (*MachineSpec, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read dsl file: %w", err)
	}

	var spec MachineSpec
	if err := yaml.Unmarshal(raw, &spec); err != nil {
		return nil, fmt.Errorf("parse dsl yaml: %w", err)
	}

	return &spec, nil
}
