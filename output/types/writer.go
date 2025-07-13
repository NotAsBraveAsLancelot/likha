package types

import "io"

// Writer is the interface that all data writers must implement.
type Writer interface {
	// WriteHeader writes the header of the file, if applicable.
	WriteHeader(headers []string) error
	// WriteRow writes a single row of data.
	WriteRow(row map[string]interface{}) error
	// Close finalizes the writing process and closes the underlying writer.
	Close() error
}

// WriterFactory creates a Writer based on the provided configuration.
type WriterFactory func(w io.Writer, settings map[string]interface{}) (Writer, error)
