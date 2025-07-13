package json

import (
	"encoding/json"
	"io"

	"likha/output/types"
)

// JSONWriter writes data in JSON format.
type JSONWriter struct {
	writer  io.Writer
	encoder *json.Encoder
	isFirst bool
}

// New creates a new JSONWriter.
func New(w io.Writer, settings map[string]interface{}) (types.Writer, error) {
	jw := &JSONWriter{
		writer:  w,
		encoder: json.NewEncoder(w),
		isFirst: true,
	}
	if pretty, ok := settings["pretty"].(bool); ok && pretty {
		jw.encoder.SetIndent("", "  ")
	}
	return jw, nil
}

// WriteHeader starts the JSON array.
func (w *JSONWriter) WriteHeader(headers []string) error {
	_, err := w.writer.Write([]byte("[\n"))
	return err
}

// WriteRow writes a single row as a JSON object.
func (w *JSONWriter) WriteRow(row map[string]interface{}) error {
	if !w.isFirst {
		_, err := w.writer.Write([]byte(",\n"))
		if err != nil {
			return err
		}
	}
	w.isFirst = false
	return w.encoder.Encode(row)
}

// Close ends the JSON array.
func (w *JSONWriter) Close() error {
	_, err := w.writer.Write([]byte("\n]\n"))
	return err
}
