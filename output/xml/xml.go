package xml

import (
	"encoding/xml"
	"fmt"
	"io"

	"likha/output/types"
)

// XMLWriter writes data in XML format.
type XMLWriter struct {
	writer   io.Writer
	encoder  *xml.Encoder
	rootNode string
}

// New creates a new XMLWriter.
func New(w io.Writer, settings map[string]interface{}) (types.Writer, error) {
	xw := &XMLWriter{
		writer:   w,
		encoder:  xml.NewEncoder(w),
		rootNode: "data",
	}
	if root, ok := settings["root_node"].(string); ok {
		xw.rootNode = root
	}
	xw.encoder.Indent("", "  ")
	return xw, nil
}

// WriteHeader writes the XML header and root element.
func (w *XMLWriter) WriteHeader(headers []string) error {
	_, err := w.writer.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	return w.encoder.EncodeToken(xml.StartElement{Name: xml.Name{Local: w.rootNode}})
}

// WriteRow writes a single row as an XML element.
func (w *XMLWriter) WriteRow(row map[string]interface{}) error {
	start := xml.StartElement{Name: xml.Name{Local: "row"}}
	if err := w.encoder.EncodeToken(start); err != nil {
		return err
	}

	for key, val := range row {
		elemStart := xml.StartElement{Name: xml.Name{Local: key}}
		if err := w.encoder.EncodeToken(elemStart); err != nil {
			return err
		}
		if err := w.encoder.EncodeToken(xml.CharData(fmt.Sprintf("%v", val))); err != nil {
			return err
		}
		if err := w.encoder.EncodeToken(elemStart.End()); err != nil {
			return err
		}
	}

	return w.encoder.EncodeToken(start.End())
}

// Close closes the root XML element.
func (w *XMLWriter) Close() error {
	if err := w.encoder.EncodeToken(xml.EndElement{Name: xml.Name{Local: w.rootNode}}); err != nil {
		return err
	}
	return w.encoder.Flush()
}
