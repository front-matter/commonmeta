package crossrefxml

import (
	"encoding/xml"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/dateutils"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/utils"
	"github.com/google/uuid"
	"github.com/xeipuuv/gojsonschema"
)

type StringMap map[string]string

type Account struct {
	Depositor  string `xml:"depositor"`
	Email      string `xml:"email"`
	Registrant string `xml:"registrant"`
}

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
	c := &Crossref{}
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
	if len(data.FundingReferences) > 0 {
		assertion := []Assertion{}
		for _, fundingReference := range data.FundingReferences {
			a := []Assertion{}
			f := Assertion{}
			fi := Assertion{}
			if fundingReference.FunderIdentifier != "" {
				fi = Assertion{
					Name: "funder_identifier",
					Text: fundingReference.FunderIdentifier,
				}
			}
			f = Assertion{
				Name:      "funder_name",
				Text:      fundingReference.FunderName,
				Assertion: []Assertion{fi},
			}
			a = append(a, f)
			if fundingReference.AwardNumber != "" {
				f = Assertion{
					Name: "award_number",
					Text: fundingReference.AwardNumber,
				}
				a = append(a, f)
			}
			fg := Assertion{
				Name:      "fundgroup",
				Assertion: a,
			}
			assertion = append(assertion, fg)
		}
		program = append(program, &Program{
			Name:      "fundref",
			Assertion: assertion,
		})
	}

	if data.License.URL != "" {
		licenseRef := []LicenseRef{}
		licenseRef = append(licenseRef, LicenseRef{
			AppliesTo: "vor",
			Text:      data.License.URL,
		})
		licenseRef = append(licenseRef, LicenseRef{
			AppliesTo: "tdm",
			Text:      data.License.URL,
		})
		program = append(program, &Program{
			Name:       "AccessIndicators",
			LicenseRef: licenseRef,
		})
	}
	if len(data.Relations) > 0 {
		relatedItem := []RelatedItem{}
		for _, relation := range data.Relations {
			id, identifierType := utils.ValidateID(relation.ID)
			if identifierType == "URL" {
				identifierType = "uri"
			}
			if slices.Contains(InterWorkRelationTypes, relation.Type) && id != "" {
				interWorkRelation := &InterWorkRelation{
					RelationshipType: relation.Type,
					IdentifierType:   strings.ToLower(identifierType),
					Text:             id,
				}
				r := RelatedItem{
					InterWorkRelation: interWorkRelation,
				}
				relatedItem = append(relatedItem, r)
			}
			if slices.Contains(IntraWorkRelationTypes, relation.Type) && id != "" {
				intraWorkRelation := &IntraWorkRelation{
					RelationshipType: relation.Type,
					IdentifierType:   strings.ToLower(identifierType),
					Text:             id,
				}
				r := RelatedItem{
					IntraWorkRelation: intraWorkRelation,
				}
				relatedItem = append(relatedItem, r)
			}
		}
		program = append(program, &Program{
			Name:        "relations",
			RelatedItem: relatedItem,
		})
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

	return c, nil
}

// Write writes Crossrefxml metadata.
func Write(data commonmeta.Data, account Account) ([]byte, []gojsonschema.ResultError) {
	type Depositor struct {
		DepositorName string `xml:"depositor_name"`
		Email         string `xml:"email_address"`
	}

	type Head struct {
		DOIBatchID string    `xml:"doi_batch_id"`
		Timestamp  string    `xml:"timestamp"`
		Depositor  Depositor `xml:"depositor"`
		Registrant string    `xml:"registrant"`
	}

	type DOIBatch struct {
		XMLName        xml.Name `xml:"doi_batch"`
		Xmlns          string   `xml:"xmlns,attr"`
		Version        string   `xml:"version,attr"`
		Xsi            string   `xml:"xsi,attr"`
		SchemaLocation string   `xml:"schemaLocation,attr"`
		Head           Head     `xml:"head"`
		Body           Crossref `xml:"body"`
	}

	crossref, err := Convert(data)
	if err != nil {
		fmt.Println(err)
	}

	depositor := Depositor{
		DepositorName: account.Depositor,
		Email:         account.Email,
	}
	uuid, _ := uuid.NewRandom()
	head := Head{
		DOIBatchID: uuid.String(),
		Timestamp:  time.Now().Format(dateutils.CrossrefDateTimeFormat),
		Depositor:  depositor,
		Registrant: account.Registrant,
	}
	doiBatch := DOIBatch{
		Xmlns:          "http://www.crossref.org/schema/5.3.1",
		Version:        "5.3.1",
		Xsi:            "http://www.w3.org/2001/XMLSchema-instance",
		SchemaLocation: "http://www.crossref.org/schema/5.3.1 ",
		Head:           head,
		Body:           *crossref,
	}

	output, err := xml.MarshalIndent(doiBatch, "", "  ")
	if err == nil {
		fmt.Println(err)
	}
	output = []byte(xml.Header + string(output))
	return output, nil
}
