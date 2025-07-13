package csv

import (
	"encoding/csv"
	"fmt"
	"io"

	"likha/output/types"
)

// CSVWriter writes data in CSV format.
type CSVWriter struct {
	writer  *csv.Writer
	headers []string
}

// New creates a new CSVWriter.
func New(w io.Writer, settings map[string]interface{}) (types.Writer, error) {
	writer := csv.NewWriter(w)

	if delim, ok := settings["delimiter"].(string); ok && len(delim) == 1 {
		writer.Comma = rune(delim[0])
	}

	return &CSVWriter{writer: writer}, nil
}

// WriteHeader writes the CSV header row.
func (w *CSVWriter) WriteHeader(headers []string) error {
	w.headers = headers
	includeHeaders := true

	if w.writer.Comma == 't' {
		includeHeaders = false
	}

	if includeHeaders {
		return w.writer.Write(headers)
	}
	return nil
}

// WriteRow writes a single row to the CSV file.
func (w *CSVWriter) WriteRow(row map[string]interface{}) error {
	record := make([]string, len(w.headers))
	for i, h := range w.headers {
		record[i] = fmt.Sprintf("%v", row[h])
	}
	return w.writer.Write(record)
}

// Close flushes the writer.
func (w *CSVWriter) Close() error {
	w.writer.Flush()
	return w.writer.Error()
}
