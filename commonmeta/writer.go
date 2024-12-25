package commonmeta

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"github.com/front-matter/commonmeta/schemautils"
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
func Write(data Data) ([]byte, error) {
	output, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	err = schemautils.SchemaErrors(output)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// WriteAll writes commonmeta metadata in slice format.
func WriteAll(list []Data) ([]byte, error) {
	output, err := json.Marshal(list)
	if err != nil {
		fmt.Println(err)
	}
	err = schemautils.SchemaErrors(output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
