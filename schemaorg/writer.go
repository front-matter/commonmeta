// Package schemaorg converts Schema.org metadata to/from the commonmeta metadata format.
package schemaorg

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/front-matter/commonmeta/commonmeta"
)

// Convert converts commonmeta metadata to Schema.org metadata.
func Convert(data commonmeta.Data) (SchemaOrg, error) {
	var schemaorg SchemaOrg
	schemaorg.Context = "http://schema.org"
	schemaorg.ID = data.ID
	schemaorg.Type = commonmeta.CMToSOMappings[data.Type]

	schemaorg.AdditionalType = data.AdditionalType
	if len(data.Contributors) > 0 {
		for _, c := range data.Contributors {
			if slices.Contains(c.ContributorRoles, "Author") {
				if c.Type == "Person" {
					var affiliation Organization
					if len(c.Affiliations) > 0 {
						affiliation = Organization{
							ID:   c.Affiliations[0].ID,
							Name: c.Affiliations[0].Name,
						}
					}
					schemaorg.Author = append(schemaorg.Author, Contributor{
						ID:          c.ID,
						Type:        "Person",
						GivenName:   c.GivenName,
						FamilyName:  c.FamilyName,
						Affiliation: affiliation,
					})
				} else if c.Type == "Organization" {
					schemaorg.Author = append(schemaorg.Author, Contributor{
						ID:   c.ID,
						Type: "Organization",
						Name: c.Name,
					})
				}
			} else if slices.Contains(c.ContributorRoles, "Editor") {
				if c.Type == "Person" {
					var affiliation Organization
					if len(c.Affiliations) > 0 {
						affiliation = Organization{
							ID:   c.Affiliations[0].ID,
							Name: c.Affiliations[0].Name,
						}
					}
					schemaorg.Editor = append(schemaorg.Editor, Editor{
						ID:          c.ID,
						GivenName:   c.GivenName,
						FamilyName:  c.FamilyName,
						Affiliation: affiliation,
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
			if reference.ID != "" {
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
		for _, description := range data.Descriptions {
			schemaorg.Description = append(schemaorg.Description, description.Description)
		}
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
func Write(data commonmeta.Data) ([]byte, error) {
	schemaorg, err := Convert(data)
	if err != nil {
		fmt.Println(err)
	}
	output, err := json.Marshal(schemaorg)
	return output, err
}

// WriteAll writes a list of schemaorg metadata.
func WriteAll(list []commonmeta.Data) ([]byte, error) {
	var schemaorgList []SchemaOrg
	for _, data := range list {
		csl, err := Convert(data)
		if err != nil {
			fmt.Println(err)
		}
		schemaorgList = append(schemaorgList, csl)
	}
	output, err := json.Marshal(schemaorgList)
	return output, err
}
