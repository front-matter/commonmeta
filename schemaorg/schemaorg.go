// Package schemaorg converts Schema.org metadata to/from the commonmeta metadata format.
package schemaorg

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/xeipuuv/gojsonschema"
)

// SchemaOrg represents the Schema.org metadata.
type SchemaOrg struct {
	Context               string        `json:"@context"`
	ID                    string        `json:"@id"`
	Type                  string        `json:"@type"`
	AdditionalType        string        `json:"additionalType,omitempty"`
	Author                []Author      `json:"author,omitempty"`
	Citation              []Citation    `json:"citation,omitempty"`
	CodeRepository        string        `json:"codeRepository,omitempty"`
	DateCreated           string        `json:"dateCreated,omitempty"`
	DatePublished         string        `json:"datePublished,omitempty"`
	DateModified          string        `json:"dateModified,omitempty"`
	Description           string        `json:"description,omitempty"`
	Distribution          []MediaObject `json:"distribution,omitempty"`
	Editor                []Editor      `json:"editor,omitempty"`
	Encoding              []MediaObject `json:"encoding,omitempty"`
	Identifier            []string      `json:"identifier,omitempty"`
	IncludedInDataCatalog DataCatalog   `json:"includedInDataCatalog,omitempty"`
	InLanguage            string        `json:"inLanguage,omitempty"`
	Keywords              string        `json:"keywords,omitempty"`
	License               string        `json:"license,omitempty"`
	Name                  string        `json:"name,omitempty"`
	PageStart             string        `json:"pageStart,omitempty"`
	PageEnd               string        `json:"pageEnd,omitempty"`
	Periodical            Periodical    `json:"periodical,omitempyt"`
	Provider              Provider      `json:"provider,omitempty"`
	Publisher             Publisher     `json:"publisher,omitempty"`
	URL                   string        `json:"url,omitempty"`
	Version               string        `json:"version,omitempty"`
}

// Content represents the SchemaOrg metadata returned from SchemaOrg sources. The type is more
// flexible than the SchemaOrg type, allowing for different formats of some metadata.
// Identifier can be string or []string.
type Content struct {
	*SchemaOrg
	Identifier json.RawMessage `json:"identifier"`
}

// Author represents the author of this CreativeWork.
type Author struct {
	ID           string         `json:"@id,omitempty"`
	Type         string         `json:"@type,omitempty"`
	GivenName    string         `json:"givenName,omitempty"`
	FamilyName   string         `json:"familyName"`
	Name         string         `json:"name,omitempty"`
	Affiliations []Organization `json:"affiliations,omitempty"`
}

// Citation represents a citation or reference to another creative work, such as another publication, web page, scholarly article, etc.
type Citation struct {
	ID   string `json:"@id,omitempty"`
	Type string `json:"@type,omitempty"`
	Name string `json:"name,omitempty"`
}

// Coderepository
type CodeRepository struct {
}

// Datacatalog represents a collection of datasets.
type DataCatalog struct {
	ID   string `json:"@id,omitempty"`
	Type string `json:"@type,omitempty"`
	Name string `json:"name,omitempty"`
}

// Editor represents
type Editor struct {
	ID           string         `json:"@id,omitempty"`
	Type         string         `json:"@type,omitempty"`
	GivenName    string         `json:"givenName,omitempty"`
	FamilyName   string         `json:"familyName"`
	Name         string         `json:"name,omitempty"`
	Affiliations []Organization `json:"affiliations,omitempty"`
}

// MediaObject represents a media object, such as an image, video, audio, or text object
// embedded in a web page or a downloadable dataset i.e. DataDownload.
type MediaObject struct {
	Type           string `json:"@type"`
	ContentURL     string `json:"contentUrl"`
	EncodingFormat string `json:"encodingFormat,omitempty"`
	Name           string `json:"name,omitempty"`
	SHA256         string `json:"sha256,omitempty"`
	Size           string `json:"size,omitempty"`
}

// Organization represents an organization such as a school, NGO, corporation, club, etc.
type Organization struct {
	ID   string `json:"@id,omitempty"`
	Name string `json:"name"`
}

// Periodical represents a publication in any medium issued in successive parts bearing numerical or chronological designations and intended to continue indefinitely, such as a magazine, scholarly journal, or newspaper.
type Periodical struct {
	ID   string `json:"@id,omitempty"`
	Type string `json:"@type"`
	Name string `json:"name,omitempty"`
	ISSN string `json:"issn,omitempty"`
}

// Person represents a person (alive, dead, undead, or fictional).
type Person struct {
	ID         string `json:"@id,omitempty"`
	GivenName  string `json:"givenName,omitempty"`
	FamilyName string `json:"familyName"`
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
func Read(content Content) (commonmeta.Data, error) {
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
	if len(data.Contributors) > 0 {
		for _, c := range data.Contributors {
			if slices.Contains(c.ContributorRoles, "Author") {
				if c.Type == "Person" {
					var affiliations []Organization
					for _, affiliation := range c.Affiliations {
						affiliations = append(affiliations, Organization{
							ID:   affiliation.ID,
							Name: affiliation.Name,
						})
					}
					schemaorg.Author = append(schemaorg.Author, Author{
						ID:           c.ID,
						Type:         "Person",
						GivenName:    c.GivenName,
						FamilyName:   c.FamilyName,
						Affiliations: affiliations,
					})
				} else if c.Type == "Organization" {
					schemaorg.Author = append(schemaorg.Author, Author{
						ID:   c.ID,
						Type: "Organization",
						Name: c.Name,
					})
				}
			} else if slices.Contains(c.ContributorRoles, "Editor") {
				if c.Type == "Person" {
					var affiliations []Organization
					for _, affiliation := range c.Affiliations {
						affiliations = append(affiliations, Organization{
							ID:   affiliation.ID,
							Name: affiliation.Name,
						})
					}
					schemaorg.Editor = append(schemaorg.Editor, Editor{
						ID:           c.ID,
						GivenName:    c.GivenName,
						FamilyName:   c.FamilyName,
						Affiliations: affiliations,
					})
				} else if c.Type == "Organization" {
					schemaorg.Editor = append(schemaorg.Editor, Editor{
						ID:   c.ID,
						Name: c.Name,
					})
				}
			}
		}
	}

	if len(data.References) > 0 {
		for _, reference := range data.References {
			t := "CreativeWork"
			if reference.Type == "JournalArticle" {
				t = "ScholarlyArticle"
			}
			schemaorg.Citation = append(schemaorg.Citation, Citation{
				ID:   reference.ID,
				Type: t,
				Name: reference.Title,
			})
		}
	}

	if data.Type == "Dataset" {
		schemaorg.IncludedInDataCatalog = DataCatalog{
			ID:   data.Container.Identifier,
			Type: "DataCatalog",
			Name: data.Container.Title,
		}
	} else {
		var ISSN string
		var ID string
		if data.Container.IdentifierType == "ISSN" {
			ISSN = data.Container.Identifier
			ID = ""
		}
		schemaorg.Periodical = Periodical{
			ID:   ID,
			Type: "Periodical",
			Name: data.Container.Title,
			ISSN: ISSN,
		}
	}

	schemaorg.DateCreated = data.Date.Created
	schemaorg.DatePublished = data.Date.Published
	schemaorg.DateModified = data.Date.Updated
	if len(data.Descriptions) > 0 {
		schemaorg.Description = data.Descriptions[0].Description
	}
	var mediaObjects []MediaObject
	if len(data.Files) > 0 {
		for _, file := range data.Files {
			var size string
			if file.Size > 0 {
				size = strconv.Itoa(file.Size)
			}
			mediaObjects = append(mediaObjects, MediaObject{
				Type:           "MediaObject",
				ContentURL:     file.URL,
				EncodingFormat: file.MimeType,
				Name:           file.Key,
				SHA256:         file.Checksum,
				Size:           size,
			})
		}
	}
	if data.Type == "Dataset" {
		schemaorg.Distribution = mediaObjects
	} else {
		schemaorg.Encoding = mediaObjects
	}

	if len(data.Identifiers) > 0 {
		schemaorg.Identifier = []string{}
		for _, identifier := range data.Identifiers {
			schemaorg.Identifier = append(schemaorg.Identifier, identifier.Identifier)
		}
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
	if len(data.Titles) > 0 {
		schemaorg.Name = data.Titles[0].Title
	}
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
