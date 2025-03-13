package crossrefxml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/dateutils"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/roguescholar"
	"github.com/front-matter/commonmeta/utils"
	"github.com/google/uuid"
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
	Body           Body     `xml:"body"`
}

type Account struct {
	LoginID     string `xml:"login_id"`
	LoginPasswd string `xml:"login_passwd"`
	Depositor   string `xml:"depositor"`
	Email       string `xml:"email"`
	Registrant  string `xml:"registrant"`
}

// CMToCRMappings maps Commonmeta types to Crossref types
// source: http://api.crossref.org/types
var CMToCRMappings = map[string]string{
	"Article":            "PostedContent",
	"BlogPost":           "PostedContent",
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
func Convert(data commonmeta.Data) (Body, error) {
	c := Body{}

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
	organization := []Organization{}
	if len(data.Contributors) > 0 {
		for i, contributor := range data.Contributors {
			contributorRole := "author"
			sequence := "first"
			if i > 0 {
				sequence = "additional"
			}
			if contributor.Type == "Organization" {
				organization = append(organization, Organization{
					ContributorRole: contributorRole,
					Sequence:        sequence,
					Text:            contributor.Name,
				})
			} else {
				if len(contributor.Affiliations) > 0 {
					institution := []Institution{}
					for _, a := range contributor.Affiliations {
						if a.Name != "" {
							if a.ID != "" {
								institutionID := InstitutionID{
									Type: "ror",
									Text: a.ID,
								}
								institution = append(institution, Institution{
									InstitutionID:   &institutionID,
									InstitutionName: a.Name,
								})
							} else {
								institution = append(institution, Institution{
									InstitutionName: a.Name,
								})
							}
						}
					}
					affiliations := Affiliations{
						Institution: institution,
					}
					personName = append(personName, PersonName{
						ContributorRole: contributorRole,
						Sequence:        sequence,
						ORCID:           contributor.ID,
						GivenName:       contributor.GivenName,
						Surname:         contributor.FamilyName,
						Affiliations:    &affiliations,
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
			item := Item{
				Resource: Resource{
					Text:     file.URL,
					MimeType: file.MimeType,
				},
			}
			// as text/html item may already have been added
			if !slices.Contains(items, item) {
				items = append(items, item)
			}
		}
	}

	var issn []ISSN
	if data.Container.IdentifierType == "issn" {
		issn = append(issn, ISSN{
			MediaType: "electronic",
			Text:      data.Container.Identifier,
		})
	}

	doiData := DOIData{
		DOI:      doi,
		Resource: data.URL,
		Collection: Collection{
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

	program := []Program{}
	if len(data.FundingReferences) > 0 {
		assertion := []Assertion{}
		for _, fundingReference := range data.FundingReferences {
			a := []Assertion{}
			f := Assertion{}
			if fundingReference.FunderIdentifier != "" {
				_, type_ := utils.ValidateID(fundingReference.FunderIdentifier)
				if type_ == "ROR" {
					f = Assertion{
						Name: "ror",
						Text: fundingReference.FunderIdentifier,
					}
				} else if type_ == "Crossref Funder ID" {
					fi := Assertion{
						Name: "funder_identifier",
						Text: fundingReference.FunderIdentifier,
					}
					f = Assertion{
						Name:      "funder_name",
						Text:      fundingReference.FunderName,
						Assertion: []Assertion{fi},
					}
				}
			} else {
				f = Assertion{
					Name:      "funder_name",
					Text:      fundingReference.FunderName,
					Assertion: []Assertion{},
				}
			}
			a = append(a, f)
			if fundingReference.AwardNumber != "" {
				f = Assertion{
					Name: "award_number",
					Text: fundingReference.AwardNumber,
				}
				a = append(a, f)
			}
			if len(data.FundingReferences) > 1 {
				fg := Assertion{
					Name:      "fundgroup",
					Assertion: a,
				}
				assertion = append(assertion, fg)
			} else {
				assertion = append(assertion, a...)
			}
		}
		program = append(program, Program{
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
		program = append(program, Program{
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
			identifierTypes := []string{
				"doi",
				"issn",
				"isbn",
				"uri",
				"pmid",
				"pmcid",
				"purl",
				"arxiv",
				"ark",
				"handle",
				"uuid",
				"ecli",
				"accession",
				"other",
			}
			if slices.Contains(InterWorkRelationTypes, relation.Type) && slices.Contains(identifierTypes, strings.ToLower(identifierType)) && id != "" {
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
			if slices.Contains(IntraWorkRelationTypes, relation.Type) && slices.Contains(identifierTypes, strings.ToLower(identifierType)) && id != "" {
				intraWorkRelation := &IntraWorkRelation{
					RelationshipType: utils.CamelCaseString(relation.Type),
					IdentifierType:   strings.ToLower(identifierType),
					Text:             id,
				}
				r := RelatedItem{
					IntraWorkRelation: intraWorkRelation,
				}
				relatedItem = append(relatedItem, r)
			}
		}
		program = append(program, Program{
			Name:        "relations",
			Xmlns:       "http://www.crossref.org/relations.xsd",
			RelatedItem: relatedItem,
		})
	}

	citationList := CitationList{}
	if len(data.References) > 0 {
		for i, v := range data.References {
			key := v.Key
			if v.Key == "" {
				key = fmt.Sprintf("ref%d", i+1)
			}
			d, _ := doiutils.ValidateDOI(v.ID)
			if d != "" {
				citationList.Citation = append(citationList.Citation, Citation{
					Key: key,
					DOI: &DOI{
						Text: d,
					},
					ArticleTitle:       v.Title,
					CYear:              v.PublicationYear,
					UnstructedCitation: v.Unstructured,
				})
			} else if v.Unstructured != "" {
				citationList.Citation = append(citationList.Citation, Citation{
					Key:                key,
					ArticleTitle:       v.Title,
					CYear:              v.PublicationYear,
					UnstructedCitation: v.Unstructured,
				})
			}
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
	case "Article", "BlogPost":
		var groupTitle string
		if len(data.Subjects) > 0 {
			for _, v := range data.Subjects {
				if commonmeta.FOSMappings[v.Subject] != "" {
					groupTitle = v.Subject
					break
				}
			}
		}
		var postedDate PostedDate
		if len(data.Date.Published) > 0 {
			datePublished := dateutils.GetDateStruct(data.Date.Published)
			postedDate = PostedDate{
				MediaType: "online",
				Year:      fmt.Sprintf("%04d", datePublished.Year),
				Month:     fmt.Sprintf("%02d", datePublished.Month),
				Day:       fmt.Sprintf("%02d", datePublished.Day),
			}
		}
		c.PostedContent = append(c.PostedContent, PostedContent{
			Type:       "other",
			Language:   data.Language,
			GroupTitle: groupTitle,
			Contributors: Contributors{
				Organization: organization,
				PersonName:   personName,
			},
			Titles:       titles,
			PostedDate:   postedDate,
			Institution:  institution,
			ItemNumber:   itemNumber,
			Abstract:     abstract,
			Program:      program,
			DOIData:      doiData,
			CitationList: citationList,
		})
	case "Book":
		c.Book = append(c.Book, Book{})
	case "BookChapter":
		c.Book = append(c.Book, Book{})
	case "Component":
		c.SAComponent = append(c.SAComponent, SAComponent{})
	case "Dataset":
		c.Database = append(c.Database, Database{})
	case "Dissertation":
		c.Dissertation = append(c.Dissertation, Dissertation{
			Language:        data.Language,
			PublicationType: "thesis",
			//PersonName: 		personName,
			Titles: titles,
			// ApprovalDate: ApprovalDate{
			// 	Year:  data.Date.Published.Year,
			// 	Month: data.Date.Published.Month,
			// 	Day:   data.Date.Published.Day,
			// },
			DOIData:      doiData,
			CitationList: citationList,
		})
	case "JournalArticle":
		c.Journal = append(c.Journal, Journal{
			JournalArticle: JournalArticle{
				PublicationType: "full_text",
				Abstract:        abstract,
				// ArchiveLocations: ArchiveLocations{}
				CitationList: citationList,
				Contributors: Contributors{
					Organization: organization,
					PersonName:   personName},
				// Crossmark: Crossmark{
				// 	CustomMetadata: customMetadata,
				// },
				DOIData: doiData,
				// Pages:
				Program: program,
				// PublicationDate: data.Date.Published,
				// PublisherItem: PublisherItem{
				// 	ItemNumber: itemNumber,
				// },
				Titles: titles,
			},
			JournalMetadata: JournalMetadata{
				Language:  data.Language,
				FullTitle: data.Container.Title,
				ISSN:      issn,
			},
			JournalIssue: JournalIssue{
				JournalVolume: JournalVolume{
					Volume: data.Container.Volume,
				},
				Issue: data.Container.Issue,
			},
		})
	case "PeerReview":
		c.PeerReview = append(c.PeerReview, PeerReview{})
	case "ProceedingsArticle":
		c.Conference = append(c.Conference, Conference{})
	case "Standard":
		c.Standard = append(c.Standard, Standard{})
	}

	return c, nil
}

// Write writes Crossrefxml metadata.
func Write(data commonmeta.Data, account Account) ([]byte, error) {
	body, err := Convert(data)
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
		Body:    body,
	}

	output, _ := xml.MarshalIndent(doiBatch, "", "  ")
	// TODO: handle error
	// if err == nil {
	// 	fmt.Println(err)
	// }
	output = []byte(xml.Header + string(output))
	return output, nil
}

// WriteAll writes a list of commonmeta metadata.
func WriteAll(list []commonmeta.Data, account Account) ([]byte, error) {
	var body Body
	for _, data := range list {
		ifCrossref, ok := doiutils.GetDOIRA(data.ID)
		if !ok {
			fmt.Println("DOI is not a valid DOI:", data.ID)
			continue
		} else if ifCrossref != "Crossref" {
			continue
		}
		crossref, err := Convert(data)
		if err != nil {
			fmt.Println(err)
		}
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
	// TODO: handle error
	// if err == nil {
	// 	fmt.Println(err)
	// }
	output = []byte(xml.Header + string(output))
	return output, nil
}

// Upsert updates or creates Crossrefxml metadata.
func Upsert(record commonmeta.APIResponse, account Account, legacyKey string, data commonmeta.Data) (commonmeta.APIResponse, error) {
	isCrossref, ok := doiutils.GetDOIRA(data.ID)
	if !ok {
		return record, errors.New("DOI is not a valid DOI")
	} else if isCrossref != "Crossref" {
		return record, nil
	}

	record.DOI = data.ID

	// provide and check UUID for Rogue Scholar DOIs
	if doiutils.IsRogueScholarDOI(data.ID, "crossref") {
		for _, identifier := range data.Identifiers {
			if identifier.IdentifierType == "UUID" {
				record.UUID = identifier.Identifier
			}
		}
	}

	type HTML struct {
		Head struct {
			Title string `xml:"title"`
		} `xml:"head"`
		Body struct {
			H2 string `xml:"h2"`
			P  string `xml:"p"`
		} `xml:"body"`
	}
	type Response HTML
	var response Response

	crossrefxml, err := Write(data, account)
	if err != nil {
		return record, errors.New("JSON schema validation failed")
	}
	// the filename displayed in the Crossref admin interface, using the current UNIX timestamp
	filename := strconv.FormatInt(time.Now().Unix(), 10)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	part, _ := w.CreateFormFile("fname", filename)
	_, err = part.Write(crossrefxml)
	if err != nil {
		return record, err
	}
	w.WriteField("operation", "doMDUpload")
	w.WriteField("login_id", account.LoginID)
	w.WriteField("login_passwd", account.LoginPasswd)
	w.Close()

	postUrl := "https://doi.crossref.org/servlet/deposit"
	req, err := http.NewRequest(http.MethodPost, postUrl, strings.NewReader(b.String()))
	req.Header.Add("Content-Type", w.FormDataContentType())
	if err != nil {
		return record, err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error uploading batch", err)
		return record, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return record, err
	}
	err = xml.Unmarshal(body, &response)
	if err != nil {
		return record, err
	}
	if response.Body.H2 == "FAILURE" {
		return record, errors.New(response.Body.P)
	}
	record.Status = "submitted"

	// update rogue-scholar legacy record if legacy key is provided
	if doiutils.IsRogueScholarDOI(data.ID, "crossref") && legacyKey != "" {
		record, err = roguescholar.UpdateLegacyRecord(record, legacyKey, "doi")
		if err != nil {
			return record, err
		}
		record.Status = "submitted_and_updated_legacy"
	}

	return record, nil
}

// UpsertAll updates or creates a list of Crossrefxml metadata.
func UpsertAll(list []commonmeta.Data, account Account, legacyKey string) ([]commonmeta.APIResponse, error) {
	var records []commonmeta.APIResponse
	for _, data := range list {
		isCrossref, ok := doiutils.GetDOIRA(data.ID)
		if !ok {
			fmt.Println("DOI is not a valid DOI:", data.ID)
			continue
		} else if isCrossref != "Crossref" {
			continue
		}

		record := commonmeta.APIResponse{
			DOI: data.ID,
		}
		if doiutils.IsRogueScholarDOI(data.ID, "crossref") {
			for _, identifier := range data.Identifiers {
				if identifier.IdentifierType == "UUID" {
					record.UUID = identifier.Identifier
				}
			}
		}
		records = append(records, record)
	}

	// if no metadata to write, return empty list
	if len(records) == 0 {
		return records, nil
	}

	type HTML struct {
		Head struct {
			Title string `xml:"title"`
		} `xml:"head"`
		Body struct {
			H2 string `xml:"h2"`
			P  string `xml:"p"`
		} `xml:"body"`
	}
	type Response HTML
	var response Response

	crossrefxml, err := WriteAll(list, account)
	if err != nil {
		return records, errors.New("JSON schema validation failed")
	}
	// the filename displayed in the Crossref admin interface, using the current UNIX timestamp
	filename := strconv.FormatInt(time.Now().Unix(), 10)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	part, _ := w.CreateFormFile("fname", filename)
	_, err = part.Write(crossrefxml)
	if err != nil {
		return records, err
	}
	w.WriteField("operation", "doMDUpload")
	w.WriteField("login_id", account.LoginID)
	w.WriteField("login_passwd", account.LoginPasswd)
	w.Close()

	postUrl := "https://doi.crossref.org/servlet/deposit"
	req, err := http.NewRequest(http.MethodPost, postUrl, strings.NewReader(b.String()))
	req.Header.Add("Content-Type", w.FormDataContentType())
	if err != nil {
		return records, err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error uploading batch", err)
		return records, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return records, err
	}
	err = xml.Unmarshal(body, &response)
	if err != nil {
		return records, err
	}
	if response.Body.H2 == "FAILURE" {
		return records, errors.New(response.Body.P)
	}

	// update rogue-scholar legacy record with doi if legacy key is provided
	for i := range records {
		records[i].Status = "submitted"
		if doiutils.IsRogueScholarDOI(records[i].DOI, "crossref") && legacyKey != "" {
			records[i], err = roguescholar.UpdateLegacyRecord(records[i], legacyKey, "doi")
			if err != nil {
				return records, err
			}
			records[i].Status = "submitted_and_updated_legacy"
		}
	}
	return records, nil
}
