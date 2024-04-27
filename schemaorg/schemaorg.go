// Package schemaorg converts Schema.org metadata to/from the commonmeta metadata format.
package schemaorg

import "github.com/front-matter/commonmeta/commonmeta"

// Content represents the Schema.org metadata.
type Content struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// CMToSOMappings maps Commonmeta types to Schema.org types.
var CMToSOMappings = map[string]string{
	"Article":        "Article",
	"Audiovisual":    "CreativeWork",
	"Book":           "Book",
	"BookChapter":    "BookChapter",
	"Collection":     "CreativeWork",
	"Dataset":        "Dataset",
	"Dissertation":   "Dissertation",
	"Document":       "CreativeWork",
	"Entry":          "CreativeWork",
	"Event":          "CreativeWork",
	"Figure":         "CreativeWork",
	"Image":          "CreativeWork",
	"Instrument":     "Instrument",
	"JournalArticle": "ScholarlyArticle",
	"LegalDocument":  "Legislation",
	"Software":       "SoftwareSourceCode",
	"Presentation":   "PresentationDigitalDocument",
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
