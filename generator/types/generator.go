package types

import (
	"fmt"

	"likha/config"
)

// Generator is the interface that all value generators must implement.
type Generator interface {
	Generate(row map[string]interface{}) (interface{}, error)
}

// GeneratorFactory creates a Generator based on the provided configuration.
func GeneratorFactory(cfg config.GeneratorConfig, generators map[string]Generator) (Generator, error) {
	// This is a placeholder for the actual factory implementation
	// which will be in the main generator package to avoid circular dependencies.
	return nil, fmt.Errorf("not implemented in types package")
}
