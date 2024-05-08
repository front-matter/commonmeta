// Package codemeta provides a function to read Codemeta metadata and convert it to commonmeta.
package codemeta

import (
	"github.com/front-matter/commonmeta/commonmeta"
)

type content struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Read reads codemeta and converts it to commonmeta.
func Read(content content) (commonmeta.Data, error) {
	var data commonmeta.Data

	data.ID = content.ID
	return data, nil
}
