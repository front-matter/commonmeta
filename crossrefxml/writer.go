package crossrefxml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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
	Xmlns          string   `xml:"xmlns,attr,omitempty"`
	Version        string   `xml:"version,attr,omitempty"`
	Xsi            string   `xml:"xsi,attr,omitempty"`
	SchemaLocation string   `xml:"schemaLocation,attr,omitempty"`
	Head           Head     `xml:"head"`
	Body           Crossref `xml:"body"`
}

type Account struct {
	LoginID       string `xml:"login_id"`
	LoginPassword string `xml:"login_passwd"`
	Depositor     string `xml:"depositor"`
	Email         string `xml:"email"`
	Registrant    string `xml:"registrant"`
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
func Convert(data commonmeta.Data) (Crossref, error) {
	c := Crossref{}
	abstract := []Abstract{}
	if len(data.Descriptions) > 0 {
		for _, description := range data.Descriptions {
			if description.Type == "Abstract" {
				p := []P{}
				p = append(p, P{
					Text: description.Description,
				})
				abstract = append(abstract, Abstract{
					Xmlns: "http://www.ncbi.nlm.nih.gov/JATS1",
					P:     p,
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
			if len(contributor.Affiliations) > 0 {
				institution := []Institution{}
				for _, a := range contributor.Affiliations {
					if a.Name != "" {
						if a.ID != "" {
							institutionID := &InstitutionID{
								IDType: "ror",
								Text:   a.ID,
							}
							institution = append(institution, Institution{
								InstitutionID:   institutionID,
								InstitutionName: a.Name,
							})
						} else {
							institution = append(institution, Institution{
								InstitutionName: a.Name,
							})
						}
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
			} else {
				personName = append(personName, PersonName{
					ContributorRole: contributorRole,
					Sequence:        sequence,
					ORCID:           contributor.ID,
					GivenName:       contributor.GivenName,
					Surname:         contributor.FamilyName,
				})
			}
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
			if file.MimeType == "text/markdown" {
				// Crossref schema currently doesn't support text/markdown
				file.MimeType = "text/plain"
			}
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
			Xmlns:     "http://www.crossref.org/fundref.xsd",
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
			Xmlns:      "http://www.crossref.org/AccessIndicators.xsd",
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
				// Crossref relation types are camel case rather than pascal case
				interWorkRelation := &InterWorkRelation{
					RelationshipType: utils.CamelCaseString(relation.Type),
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
			Xmlns:       "http://www.crossref.org/relations.xsd",
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
		c.PostedContent = append(c.PostedContent, &PostedContent{
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
		})
	case "JournalArticle":
		c.Journal = append(c.Journal, &Journal{})
	}

	return c, nil
}

// Write writes Crossrefxml metadata.
func Write(data commonmeta.Data, account Account) ([]byte, []gojsonschema.ResultError) {
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
		Xmlns:   "http://www.crossref.org/schema/5.3.1",
		Version: "5.3.1",
		Head:    head,
		Body:    crossref,
	}

	output, err := xml.MarshalIndent(doiBatch, "", "  ")
	if err == nil {
		// TODO: handle error
		// fmt.Println(err)
	}
	output = []byte(xml.Header + string(output))
	return output, nil
}

// WriteAll writes a list of commonmeta metadata.
func WriteAll(list []commonmeta.Data, account Account) ([]byte, []gojsonschema.ResultError) {
	var body Crossref
	for _, data := range list {
		crossref, _ := Convert(data)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// workaround to handle the different content types
		body.Book = append(body.Book, crossref.Book...)
		body.Conference = append(body.Conference, crossref.Conference...)
		body.Database = append(body.Database, crossref.Database...)
		body.Dissertation = append(body.Dissertation, crossref.Dissertation...)
		body.Journal = append(body.Journal, crossref.Journal...)
		body.PeerReview = append(body.PeerReview, crossref.PeerReview...)
		body.PostedContent = append(body.PostedContent, crossref.PostedContent...)
		body.SAComponent = append(body.SAComponent, crossref.SAComponent...)
		body.Standard = append(body.Standard, crossref.Standard...)
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
		Xmlns:   "http://www.crossref.org/schema/5.3.1",
		Version: "5.3.1",
		Head:    head,
		Body:    body,
	}

	output, _ := xml.MarshalIndent(doiBatch, "", "  ")
	// if err == nil {
	// 	// fmt.Println(err)
	// }
	output = []byte(xml.Header + string(output))
	return output, nil
}

func Upload(content []byte, account Account) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	postUrl := "https://doi.crossref.org/servlet/deposit"

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	part, _ := w.CreateFormFile("fname", "output.xml")
	_, err := part.Write(content)
	if err != nil {
		return "", err
	}
	w.WriteField("operation", "doMDUpload")
	w.WriteField("login_id", account.LoginID)
	w.WriteField("login_passwd", account.LoginPassword)
	defer w.Close()

	req, err := http.NewRequest(http.MethodPost, postUrl, &b)
	req.Header.Add("Content-Type", w.FormDataContentType())
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error uploading batch", err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(body))
	message := "Your batch submission was successfully received. " + resp.Status
	return message, nil
}
