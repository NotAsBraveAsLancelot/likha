package foreignkey

import (
	"fmt"

	"likha/config"

	"likha/generator/types"
)

// ForeignKeyGenerator generates a value based on the value of another field.
type ForeignKeyGenerator struct {
	sourceField string
	valueMap    map[interface{}]types.Generator
}

// New creates a new ForeignKeyGenerator.
func New(
	cfg config.GeneratorConfig,
	allGenerators map[string]types.Generator,
	factoryFn func(config.GeneratorConfig, map[string]types.Generator) (types.Generator, error),
) (types.Generator, error) {
	if cfg.SourceField == "" {
		return nil, fmt.Errorf("foreignkey generator requires a 'source_field'")
	}
	if len(cfg.Map) == 0 {
		return nil, fmt.Errorf("foreignkey generator requires a 'map' of values to generators")
	}

	g := &ForeignKeyGenerator{
		sourceField: cfg.SourceField,
		valueMap:    make(map[interface{}]types.Generator),
	}

	for key, genCfg := range cfg.Map {
		gen, err := factoryFn(genCfg, allGenerators)
		if err != nil {
			return nil, fmt.Errorf("failed to create generator for foreign key map value '%v': %w", key, err)
		}
		g.valueMap[key] = gen
	}

	return g, nil
}

// Generate looks up the source field's value and uses the corresponding generator.
func (g *ForeignKeyGenerator) Generate(row map[string]interface{}) (interface{}, error) {
	sourceValue, ok := row[g.sourceField]
	if !ok {
		// This can happen if the source field hasn't been generated yet for this row.
		// The runner must ensure correct ordering.
		return nil, fmt.Errorf("source field '%s' not found in current row", g.sourceField)
	}

	mappedGenerator, ok := g.valueMap[sourceValue]
	if !ok {
		// Return nil or a default if no mapping exists for the source value
		return nil, nil
	}

	return mappedGenerator.Generate(row)
}
