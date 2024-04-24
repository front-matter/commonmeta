// Package cff provides a function to read CFF and convert it to commonmeta.
package cff

import "github.com/front-matter/commonmeta/commonmeta"

type content struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Read reads CFF and converts it to commonmeta.
func Read(content content) (commonmeta.Data, error) {
	var data commonmeta.Data
	data.ID = content.ID
	return data, nil
}
