package expression

import (
	"fmt"

	"likha/expression"

	"likha/generator/types"
)

// ExpressionGenerator evaluates an expression string.
type ExpressionGenerator struct {
	evaluator *expression.Evaluator
	template  string
}

// New creates a new ExpressionGenerator.
func New(settings map[string]interface{}) (types.Generator, error) {
	template, ok := settings["expression"].(string)
	if !ok {
		return nil, fmt.Errorf("expression generator requires an 'expression' string setting")
	}
	return &ExpressionGenerator{
		evaluator: expression.NewEvaluator(),
		template:  template,
	}, nil
}

// Generate evaluates the expression.
func (g *ExpressionGenerator) Generate(row map[string]interface{}) (interface{}, error) {
	return g.evaluator.Evaluate(g.template, row)
}
