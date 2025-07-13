package yaml

import (
	"io"

	"likha/output/types"

	"gopkg.in/yaml.v3"
)

// YAMLWriter writes data in YAML format.
type YAMLWriter struct {
	writer  io.Writer
	encoder *yaml.Encoder
}

// New creates a new YAMLWriter.
func New(w io.Writer, settings map[string]interface{}) (types.Writer, error) {
	return &YAMLWriter{
		writer:  w,
		encoder: yaml.NewEncoder(w),
	}, nil
}

// WriteHeader is a no-op for YAML stream.
func (w *YAMLWriter) WriteHeader(headers []string) error {
	return nil
}

// WriteRow writes a single row as a YAML document.
func (w *YAMLWriter) WriteRow(row map[string]interface{}) error {
	return w.encoder.Encode(row)
}

// Close closes the YAML encoder.
func (w *YAMLWriter) Close() error {
	return w.encoder.Close()
}
