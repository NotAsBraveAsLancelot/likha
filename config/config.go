package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure.
type Config struct {
	Fields []Field      `yaml:"fields"`
	Output OutputConfig `yaml:"output"`
}

// Field represents a single data field to be generated.
type Field struct {
	Name      string          `yaml:"name"`
	Generator GeneratorConfig `yaml:"generator"`
}

// GeneratorConfig holds the configuration for a value generator.
type GeneratorConfig struct {
	Type        string                          `yaml:"type"`
	Settings    map[string]interface{}          `yaml:"settings"`
	SourceField string                          `yaml:"source_field"` // For foreignkey
	Map         map[interface{}]GeneratorConfig `yaml:"map"`          // For foreignkey
}

// OutputConfig defines the output format and its settings.
type OutputConfig struct {
	Type     string                 `yaml:"type"`
	File     string                 `yaml:"file"`
	Settings map[string]interface{} `yaml:"settings"`
}

// LoadConfig reads a YAML configuration file from the given path and unmarshals it.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config file %s: %w", path, err)
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal config file %s: %w", path, err)
	}

	return &cfg, nil
}
