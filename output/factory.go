package output

import (
	"fmt"
	"io"

	"likha/config"

	"likha/output/csv"
	"likha/output/json"
	"likha/output/types"
	"likha/output/xml"
	"likha/output/yaml"
)

// NewWriter creates a new data writer based on the output configuration.
func NewWriter(cfg *config.OutputConfig, w io.Writer) (types.Writer, error) {
	switch cfg.Type {
	case "csv":
		return csv.New(w, cfg.Settings)
	case "json":
		return json.New(w, cfg.Settings)
	case "xml":
		return xml.New(w, cfg.Settings)
	case "yaml":
		return yaml.New(w, cfg.Settings)
	default:
		return nil, fmt.Errorf("unknown output type: %s", cfg.Type)
	}
}
