package simple

import (
	"fmt"

	"likha/generator/types"
)

// SimpleGenerator generates a static value.
type SimpleGenerator struct {
	value interface{}
}

// New creates a new SimpleGenerator.
func New(settings map[string]interface{}) (types.Generator, error) {
	val, ok := settings["value"]
	if !ok {
		return nil, fmt.Errorf("simple generator requires a 'value' setting")
	}
	return &SimpleGenerator{value: val}, nil
}

// Generate returns the static value.
func (g *SimpleGenerator) Generate(row map[string]interface{}) (interface{}, error) {
	return g.value, nil
}
