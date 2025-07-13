package factory

import (
	"fmt"

	"likha/config"
	"likha/generator/builtin"
	"likha/generator/custom"
	"likha/generator/expression"
	"likha/generator/foreignkey"
	"likha/generator/list"
	"likha/generator/simple"
	"likha/generator/types"
)

// NewGenerator creates a new generator based on the provided configuration.
// It acts as a factory routing to the specific generator implementations.
func NewGenerator(cfg config.GeneratorConfig, allGenerators map[string]types.Generator) (types.Generator, error) {
	switch cfg.Type {
	case "simple":
		return simple.New(cfg.Settings)
	case "list":
		return list.New(cfg.Settings)
	case "builtin":
		return builtin.New(cfg.Settings)
	case "expression":
		return expression.New(cfg.Settings)
	case "custom":
		return custom.New(cfg.Settings)
	case "foreignkey":
		return foreignkey.New(cfg, allGenerators, NewGenerator)
	default:
		return nil, fmt.Errorf("unknown generator type: %s", cfg.Type)
	}
}
