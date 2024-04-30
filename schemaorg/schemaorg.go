// Package schemaorg converts Schema.org metadata to/from the commonmeta metadata format.
package schemaorg

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/xeipuuv/gojsonschema"
)

// SchemaOrg represents the Schema.org metadata.
type SchemaOrg struct {
	Context        string    `json:"@context"`
	ID             string    `json:"@id"`
	Type           string    `json:"@type"`
	AdditionalType string    `json:"additionalType,omitempty"`
	DateCreated    string    `json:"dateCreated,omitempty"`
	DatePublished  string    `json:"datePublished,omitempty"`
	DateModified   string    `json:"dateModified,omitempty"`
	Description    string    `json:"description,omitempty"`
	Identifier     []string  `json:"identifier,omitempty"`
	InLanguage     string    `json:"inLanguage,omitempty"`
	Keywords       string    `json:"keywords,omitempty"`
	License        string    `json:"license,omitempty"`
	PageStart      string    `json:"pageStart,omitempty"`
	PageEnd        string    `json:"pageEnd,omitempty"`
	Provider       Provider  `json:"provider,omitempty"`
	Publisher      Publisher `json:"publisher,omitempty"`
	URL            string    `json:"url,omitempty"`
	Version        string    `json:"version,omitempty"`
}

// Provider represents the provider of the metadata.
type Provider struct {
	Type string `json:"@type"`
	Name string `json:"name"`
}

// Publisher represents the publisher of the metadata.
type Publisher struct {
	Type string `json:"@type"`
	Name string `json:"name"`
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

// Read reads Schema.org metadata and converts it to commonmeta.
func Read(content SchemaOrg) (commonmeta.Data, error) {
	var data commonmeta.Data
	data.ID = content.ID
	return data, nil
}

// Convert converts commonmeta metadata to Schema.org metadata.
func Convert(data commonmeta.Data) (SchemaOrg, error) {
	var schemaorg SchemaOrg
	schemaorg.Context = "http://schema.org"
	schemaorg.ID = data.ID
	schemaorg.Type = CMToSOMappings[data.Type]

	schemaorg.AdditionalType = data.AdditionalType
	schemaorg.DateCreated = data.Date.Created
	schemaorg.DatePublished = data.Date.Published
	schemaorg.DateModified = data.Date.Updated
	if len(data.Descriptions) > 0 {
		schemaorg.Description = data.Descriptions[0].Description
	}
	for _, identifier := range data.Identifiers {
		schemaorg.Identifier = append(schemaorg.Identifier, identifier.Identifier)
	}
	schemaorg.InLanguage = data.Language
	if len(data.Subjects) > 0 {
		var keywords []string
		for _, subject := range data.Subjects {
			if subject.Subject != "" {
				keywords = append(keywords, subject.Subject)
			}
		}
		schemaorg.Keywords = strings.Join(keywords[:], ", ")
	}
	schemaorg.License = data.License.URL
	schemaorg.PageStart = data.Container.FirstPage
	schemaorg.PageEnd = data.Container.LastPage
	schemaorg.Provider = Provider{
		Type: "Organization",
		Name: data.Provider,
	}
	schemaorg.Publisher = Publisher{
		Type: "Organization",
		Name: data.Publisher.Name,
	}
	schemaorg.URL = data.URL
	schemaorg.Version = data.Version

	return schemaorg, nil
}

// 				"identifier": metadata.id,
// 				"@type": schema_org,
//

// 				"additionalType": additional_type,
// 				"name": parse_attributes(metadata.titles, content="title", first=True),
// 				"author": to_schema_org_creators(authors),
// 				"editor": to_schema_org_creators(editors),

// 				"periodical": periodical if periodical else None,
// 				"includedInDataCatalog": data_catalog if data_catalog else None,
// 				"distribution": media_objects if metadata.type == "Dataset" else None,
// 				"encoding": media_objects if metadata.type != "Dataset" else None,
// 				"codeRepository": code_repository,

// Write writes schemaorg metadata.
func Write(data commonmeta.Data) ([]byte, []gojsonschema.ResultError) {
	schemaorg, err := Convert(data)
	if err != nil {
		fmt.Println(err)
	}
	output, err := json.Marshal(schemaorg)
	if err != nil {
		fmt.Println(err)
	}

	return output, nil
}

// WriteList writes a list of schemaorg metadata.
func WriteList(list []commonmeta.Data) ([]byte, []gojsonschema.ResultError) {
	var schemaorgList []SchemaOrg
	for _, data := range list {
		csl, err := Convert(data)
		if err != nil {
			fmt.Println(err)
		}
		schemaorgList = append(schemaorgList, csl)
	}
	output, err := json.Marshal(schemaorgList)
	if err != nil {
		fmt.Println(err)
	}

	return output, nil
}
