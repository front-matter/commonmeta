// Package schemaorg converts Schema.org metadata to/from the commonmeta metadata format.
package schemaorg

import "github.com/front-matter/commonmeta/commonmeta"

// Content represents the Schema.org metadata.
type Content struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Get retrieves Schema.org metadata.
// func Get(id string) (Content, error) {
// 	var content Content
// 	return content, nil
// }

// Read reads Schema.org metadata and converts it to commonmeta.
func Read(content Content) (commonmeta.Data, error) {
	var data commonmeta.Data
	data.ID = content.ID
	return data, nil
}
