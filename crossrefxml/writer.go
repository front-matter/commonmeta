package crossrefxml

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/dateutils"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/utils"
	"github.com/xeipuuv/gojsonschema"
)

type StringMap map[string]string

// CMToCRMappings maps Commonmeta types to Crossref types
// source: http://api.crossref.org/types
var CMToCRMappings = map[string]string{
	"Article":            "PostedContent",
	"BookChapter":        "BookChapter",
	"BookSeries":         "BookSeries",
	"Book":               "Book",
	"Component":          "Component",
	"Dataset":            "Dataset",
	"Dissertation":       "Dissertation",
	"Grant":              "Grant",
	"JournalArticle":     "JournalArticle",
	"JournalIssue":       "JournalIssue",
	"JournalVolume":      "JournalVolume",
	"Journal":            "Journal",
	"ProceedingsArticle": "ProceedingsArticle",
	"ProceedingsSeries":  "ProceedingsSeries",
	"Proceedings":        "Proceedings",
	"ReportComponent":    "ReportComponent",
	"ReportSeries":       "ReportSeries",
	"Report":             "Report",
	"Review":             "PeerReview",
	"Other":              "Other",
}

// Convert converts Commonmeta metadata to Crossrefxml metadata
func Convert(data commonmeta.Data) (*Crossref, error) {
	c := &Crossref{
		Xmlns:          "http://www.crossref.org/schema/5.3.1",
		SchemaLocation: "http://www.crossref.org/schema/5.3.1 ",
		Version:        "5.3.1",
	}
	abstract := []Abstract{}
	if len(data.Descriptions) > 0 {
		for _, description := range data.Descriptions {
			if description.Type == "Abstract" {
				abstract = append(abstract, Abstract{
					Xmlns: "http://www.ncbi.nlm.nih.gov/JATS1",
					Text:  description.Description,
				})
			}
		}
	}
	personName := []PersonName{}
	if len(data.Contributors) > 0 {
		for i, contributor := range data.Contributors {
			contributorRole := "author"
			sequence := "first"
			if i > 0 {
				sequence = "additional"
			}
			institution := []Institution{}
			for _, a := range contributor.Affiliations {
				if a.Name != "" {
					institutionID := InstitutionID{}
					if a.ID != "" {
						institutionID = InstitutionID{
							IDType: "ror",
							Text:   a.ID,
						}
					}
					institution = append(institution, Institution{
						InstitutionID:   &institutionID,
						InstitutionName: a.Name,
					})
				}
			}
			affiliations := &Affiliations{
				Institution: institution,
			}
			personName = append(personName, PersonName{
				ContributorRole: contributorRole,
				Sequence:        sequence,
				ORCID:           contributor.ID,
				GivenName:       contributor.GivenName,
				Surname:         contributor.FamilyName,
				Affiliations:    affiliations,
			})
		}
	}

	doi, _ := doiutils.ValidateDOI(data.ID)
	var items []Item
	items = append(items, Item{
		Resource: Resource{
			Text:     data.URL,
			MimeType: "text/html",
		},
	})
	if len(data.Files) > 0 {
		for _, file := range data.Files {
			items = append(items, Item{
				Resource: Resource{
					Text:     file.URL,
					MimeType: file.MimeType,
				},
			})
		}
	}

	doiData := DOIData{
		DOI:      doi,
		Resource: data.URL,
		Collection: &Collection{
			Property: "text-mining",
			Item:     items,
		},
	}

	var itemNumber ItemNumber
	if len(data.Identifiers) > 0 {
		for _, identifier := range data.Identifiers {
			if identifier.IdentifierType == "UUID" {
				text := strings.Replace(identifier.Identifier, "-", "", 4)
				itemNumber = ItemNumber{
					Text:           text,
					ItemNumberType: "UUID",
				}
			}
		}
	}

	institution := &Institution{
		InstitutionName: data.Publisher.Name,
	}

	program := []*Program{}
	if len(data.Relations) > 0 {
		for _, relation := range data.Relations {
			if relation.Type == "IsPartOf" {
				program = append(program, &Program{})
			}
		}
	}

	citationList := CitationList{}
	if len(data.References) > 0 {
		for _, v := range data.References {
			var doi DOI
			d, _ := doiutils.ValidateDOI(v.ID)
			if d != "" {
				doi = DOI{
					Text: d,
				}
			}
			citationList.Citation = append(citationList.Citation, Citation{
				Key:                v.Key,
				DOI:                &doi,
				ArticleTitle:       v.Title,
				CYear:              v.PublicationYear,
				UnstructedCitation: v.Unstructured,
			})
		}
	}

	titles := Titles{}
	if len(data.Titles) > 0 {
		for _, title := range data.Titles {
			if title.Type == "Subtitle" {
				titles.Subtitle = title.Title
			} else if title.Type == "TranslatedTitle" {
				titles.OriginalLanguageTitle.Text = title.Title
				titles.OriginalLanguageTitle.Language = title.Language
			} else {
				titles.Title = title.Title
			}
		}
	}

	switch data.Type {
	case "Article":
		var groupTitle string
		if len(data.Subjects) > 0 {
			groupTitle = utils.CamelCaseToWords(data.Subjects[0].Subject)
		}
		var postedDate PostedDate
		if len(data.Date.Published) > 0 {
			datePublished := dateutils.GetDateStruct(data.Date.Published)
			postedDate = PostedDate{
				MediaType: "online",
				Year:      datePublished.Year,
				Month:     datePublished.Month,
				Day:       datePublished.Day,
			}
		}
		c.PostedContent = &PostedContent{
			Type:       "other",
			Language:   data.Language,
			GroupTitle: groupTitle,
			Contributors: &Contributors{
				PersonName: personName},
			Titles:       &titles,
			PostedDate:   postedDate,
			Institution:  institution,
			ItemNumber:   itemNumber,
			Abstract:     &abstract,
			Program:      program,
			DOIData:      doiData,
			CitationList: &citationList,
		}
	case "JournalArticle":
		c.Journal = &Journal{}
	}
	// required properties

	// optional properties

	return c, nil
}

// Write writes Crossrefxml metadata.
func Write(data commonmeta.Data) ([]byte, []gojsonschema.ResultError) {
	crossref, err := Convert(data)
	if err != nil {
		fmt.Println(err)
	}
	output, err := xml.MarshalIndent(crossref, "", "  ")
	if err == nil {
		fmt.Println(err)
	}
	output = []byte(xml.Header + string(output))
	return output, nil
}
