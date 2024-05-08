package commonmeta

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/front-matter/commonmeta/schemautils"
	"github.com/xeipuuv/gojsonschema"
)

type Writer struct {
	w *bufio.Writer
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: bufio.NewWriter(w),
	}
}

// Write writes commonmeta metadata.
func Write(data Data) ([]byte, []gojsonschema.ResultError) {
	output, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	validation := schemautils.JSONSchemaErrors(output)
	if !validation.Valid() {
		return nil, validation.Errors()
	}
	return output, nil
}

// WriteAll writes commonmeta metadata in slice format.
func WriteAll(list []Data) ([]byte, []gojsonschema.ResultError) {
	output, err := json.Marshal(list)
	if err != nil {
		fmt.Println(err)
	}
	validation := schemautils.JSONSchemaErrors(output)
	if !validation.Valid() {
		return nil, validation.Errors()
	}
	return output, nil
}
