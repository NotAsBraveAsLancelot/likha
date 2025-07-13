package custom

import (
	"fmt"
	"os/exec"
	"strings"

	"likha/generator/types"
)

// CustomGenerator executes an external command to generate a value.
type CustomGenerator struct {
	command string
	args    []string
}

// New creates a new CustomGenerator.
func New(settings map[string]interface{}) (types.Generator, error) {
	cmdPath, ok := settings["command"].(string)
	if !ok {
		return nil, fmt.Errorf("custom generator requires a 'command' string setting")
	}

	parts := strings.Fields(cmdPath)
	if len(parts) == 0 {
		return nil, fmt.Errorf("custom generator 'command' cannot be empty")
	}

	return &CustomGenerator{
		command: parts[0],
		args:    parts[1:],
	}, nil
}

// Generate executes the external command and returns its standard output.
func (g *CustomGenerator) Generate(row map[string]interface{}) (interface{}, error) {
	cmd := exec.Command(g.command, g.args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("custom generator command failed: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}
