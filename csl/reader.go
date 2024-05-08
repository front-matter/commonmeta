// Package csl provides a function to read CSL JSON and convert it to commonmeta.
package csl

import (
	"github.com/front-matter/commonmeta/commonmeta"
)

type content struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Read reads CSL JSON and converts it to commonmeta.
func Read(content content) (commonmeta.Data, error) {
	var data commonmeta.Data

	data.ID = content.ID
	return data, nil
}
