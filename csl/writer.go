package csl

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/dateutils"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/schemautils"
)

var CMToCSLMappings = map[string]string{
	"Article":               "article",
	"JournalArticle":        "article-journal",
	"BlogPost":              "post-weblog",
	"Book":                  "book",
	"BookChapter":           "chapter",
	"Collection":            "collection",
	"Dataset":               "dataset",
	"Document":              "document",
	"Entry":                 "entry",
	"Event":                 "event",
	"Figure":                "figure",
	"Image":                 "graphic",
	"LegalDocument":         "legal_case",
	"Manuscript":            "manuscript",
	"Map":                   "map",
	"Audiovisual":           "motion_picture",
	"Patent":                "patent",
	"Performance":           "performance",
	"Journal":               "periodical",
	"PersonalCommunication": "personal_communication",
	"Report":                "report",
	"Review":                "review",
	"Software":              "software",
	"Presentation":          "speech",
	"Standard":              "standard",
	"Dissertation":          "thesis",
	"WebPage":               "webpage",
}

// Convert converts commonmeta metadata to CSL JSON.
func Convert(data commonmeta.Data) (CSL, error) {
	var csl CSL

	csl.ID = data.ID
	csl.Type = CMToCSLMappings[data.Type]
	if data.Type == "Software" && data.Version != "" {
		csl.Type = "book"
	} else if csl.Type == "" {
		csl.Type = "document"
	}
	csl.ContainerTitle = data.Container.Title
	doi, _ := doiutils.ValidateDOI(data.ID)
	csl.DOI = doi
	csl.Issue = data.Container.Issue
	if len(data.Subjects) > 0 {
		var keywords []string
		for _, subject := range data.Subjects {
			if subject.Subject != "" {
				keywords = append(keywords, subject.Subject)
			}
		}
		csl.Keyword = strings.Join(keywords[:], ", ")
	}
	csl.Language = data.Language
	csl.Page = data.Container.Pages()
	if len(data.Titles) > 0 {
		csl.Title = data.Titles[0].Title
	}
	csl.URL = data.URL
	csl.Volume = data.Container.Volume
	if len(data.Contributors) > 0 {
		var author Author
		for _, contributor := range data.Contributors {
			if slices.Contains(contributor.ContributorRoles, "Author") {
				if contributor.FamilyName != "" {
					author = Author{
						Given:  contributor.GivenName,
						Family: contributor.FamilyName,
					}
				} else {
					author = Author{
						Literal: contributor.Name,
					}

				}
				csl.Author = append(csl.Author, author)
			}
		}
	}

	if data.Date.Published != "" {
		csl.Issued.DateAsParts = dateutils.GetDateParts(data.Date.Published)
	}
	if data.Date.Submitted != "" {
		csl.Submitted.DateAsParts = dateutils.GetDateParts(data.Date.Submitted)
	}
	if data.Date.Accessed != "" {
		csl.Accessed.DateAsParts = dateutils.GetDateParts(data.Date.Accessed)
	}

	if len(data.Descriptions) > 0 {
		csl.Abstract = data.Descriptions[0].Description
	}
	csl.Publisher = data.Publisher.Name
	csl.Version = data.Version

	return csl, nil
}

// Write writes CSL metadata.
func Write(data commonmeta.Data) ([]byte, error) {
	csl, err := Convert(data)
	if err != nil {
		fmt.Println(err)
	}
	output, err := json.Marshal(csl)
	if err != nil {
		fmt.Println(err)
	}
	err = schemautils.JSONSchemaErrors(output, "csl-data")
	if err != nil {
		return nil, err
	}

	return output, nil
}

// WriteAll writes a list of CSL metadata.
func WriteAll(list []commonmeta.Data) ([]byte, error) {
	var cslList []CSL
	for _, data := range list {
		csl, err := Convert(data)
		if err != nil {
			fmt.Println(err)
		}
		cslList = append(cslList, csl)
	}
	output, err := json.Marshal(cslList)
	if err != nil {
		fmt.Println(err)
	}
	err = schemautils.JSONSchemaErrors(output, "csl-data")
	return output, err
}
