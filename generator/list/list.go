package list

import (
	"fmt"
	"math/rand"
	"time"

	"likha/generator/types"
)

// ListGenerator randomly selects a value from a list.
type ListGenerator struct {
	values []interface{}
	r      *rand.Rand
}

// New creates a new ListGenerator.
func New(settings map[string]interface{}) (types.Generator, error) {
	v, ok := settings["values"]
	if !ok {
		return nil, fmt.Errorf("list generator requires a 'values' setting")
	}

	values, ok := v.([]interface{})
	if !ok {
		return nil, fmt.Errorf("'values' setting must be a list")
	}

	return &ListGenerator{
		values: values,
		r:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

// Generate returns a random value from the list.
func (g *ListGenerator) Generate(row map[string]interface{}) (interface{}, error) {
	if len(g.values) == 0 {
		return nil, nil
	}
	return g.values[g.r.Intn(len(g.values))], nil
}
