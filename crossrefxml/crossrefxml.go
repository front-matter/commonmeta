// Package crossrefxml provides function to convert Crossref XML metadata to/from the commonmeta metadata format.
package crossrefxml

import "github.com/front-matter/commonmeta/commonmeta"

// Content represents the Crossref XML metadata.
type Content struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Read reads Crossref XML metadata and converts it to commonmeta.
func Read(content Content) (commonmeta.Data, error) {
	var data commonmeta.Data
	data.ID = content.ID
	return data, nil
}
