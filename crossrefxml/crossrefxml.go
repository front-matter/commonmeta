// Package crossrefxml provides function to convert Crossref XML metadata to/from the commonmeta metadata format.
package crossrefxml

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"slices"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/utils"
)

// Crossrefxml represents the Crossref XML metadata.
type Crossrefxml struct {
	XMLName xml.Name `xml:"body"`
	Query   struct {
		Status string `xml:"status,attr"`
		DOI    struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
		} `xml:"doi"`
		CRMItem   []CRMItem `xml:"crm-item"`
		DOIRecord DOIRecord `xml:"doi_record"`
	} `xml:"query"`
}

// Content represents the Crossref XML metadata returned from Crossref. The type is more
// flexible than the Crossrefxml type, allowing for different formats of some metadata.
type Content struct {
	*Crossrefxml
}

type CRMItem struct {
	XMLName xml.Name `xml:"crm-item"`
	Text    string   `xml:",chardata"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Claim   string   `xml:"claim,attr"`
}

type DOIRecord struct {
	XMLName  xml.Name `xml:"doi_record"`
	Crossref struct {
		Xmlns          string        `xml:"xmlns,attr"`
		SchemaLocation string        `xml:"schemaLocation,attr"`
		Book           Book          `xml:"book,omitempty"`
		Dissertation   Dissertation  `xml:"dissertation,omitempty"`
		Journal        Journal       `xml:"journal,omitempty"`
		PeerReview     PeerReview    `xml:"peer_review,omitempty"`
		PostedContent  PostedContent `xml:"posted_content,omitempty"`
	} `xml:"crossref"`
}

type Abstract struct {
	XMLName      xml.Name `xml:"abstract"`
	Jats         string   `xml:"jats,attr"`
	AbstractType string   `xml:"abstract-type,attr"`
	Text         string   `xml:",chardata"`
	P            []string `xml:"p"`
}

type AcceptanceDate struct {
	XMLName   xml.Name `xml:"acceptance_date"`
	MediaType string   `xml:"media_type,attr"`
	Month     string   `xml:"month"`
	Day       string   `xml:"day"`
	Year      string   `xml:"year"`
}

// Affiliation represents an affiliation in Crossref XML metadata.
type Affiliation struct {
	XMLName xml.Name `xml:"affiliation"`
	Text    string   `xml:",chardata"`
}

type ArchiveLocations struct {
	XMLName xml.Name `xml:"archive_locations"`
	Text    string   `xml:",chardata"`
	Archive struct {
		Text string `xml:",chardata"`
		Name string `xml:"name,attr"`
	} `xml:"archive"`
}

type Assertion struct {
	XMLName   xml.Name `xml:"assertion"`
	Text      string   `xml:",chardata"`
	Name      string   `xml:"name,attr"`
	Assertion struct {
		Text      string `xml:",chardata"` // SystemsX, EMBO longterm p...
		Name      string `xml:"name,attr"`
		Assertion struct {
			Text     string `xml:",chardata"` // 501100006390
			Name     string `xml:"name,attr"`
			Provider string `xml:"provider,attr"`
		} `xml:"assertion"`
	} `xml:"assertion"`
}

type Book struct {
	XMLName      xml.Name `xml:"book"`
	BookType     string   `xml:"book_type,attr"`
	BookMetadata struct {
		Language     string        `xml:"language,attr"`
		Contributors []Contributor `xml:"contributors"`
		Titles       struct {
			Title    string `xml:"title"`
			Subtitle string `xml:"subtitle"`
		} `xml:"titles"`
		PublicationDate PublicationDate `xml:"publication_date"`
		DOIData         DOIData         `xml:"doi_data"`
	} `xml:"book_metadata"`
}

type Citation struct {
	XMLName      xml.Name `xml:"citation"`
	Key          string   `xml:"key,attr"`
	Text         string   `xml:",chardata"`
	JournalTitle string   `xml:"journal_title"`
	Author       string   `xml:"author"`
	Volume       string   `xml:"volume"`
	FirstPage    string   `xml:"first_page"`
	CYear        string   `xml:"cYear"`
	ArticleTitle string   `xml:"article_title"`
	Doi          struct {
		Text     string `xml:",chardata"`
		Provider string `xml:"provider,attr"`
	} `xml:"doi"`
}

type Collection struct {
	XMLName  xml.Name `xml:"collection"`
	Property string   `xml:"property,attr"`
	Item     []Item   `xml:"item"`
}

type ComponentList struct {
	XMLName   xml.Name `xml:"component_list"`
	Text      string   `xml:",chardata"`
	Component []struct {
		Text           string `xml:",chardata"`
		ParentRelation string `xml:"parent_relation,attr"`
		Titles         struct {
			Text     string `xml:",chardata"`
			Title    string `xml:"title"`
			Subtitle string `xml:"subtitle"`
		} `xml:"titles"`
		Format struct {
			Text     string `xml:",chardata"`
			MimeType string `xml:"mime_type,attr"`
		} `xml:"format"`
		DOIData DOIData `xml:"doi_data"`
	}
}

type Contributor struct {
	XMLName    xml.Name     `xml:"contributors"`
	PersonName []PersonName `xml:"person_name"`
}

type Crossmark struct {
	XMLName          xml.Name `xml:"crossmark"`
	Text             string   `xml:",chardata"`
	CrossmarkVersion string   `xml:"crossmark_version"`
	CrossmarkPolicy  string   `xml:"crossmark_policy"`
	CrossmarkDomains struct {
		Text            string `xml:",chardata"`
		CrossmarkDomain struct {
			Text   string `xml:",chardata"`
			Domain string `xml:"domain"`
		} `xml:"crossmark_domain"`
	} `xml:"crossmark_domains"`
	CrossmarkDomainExclusive string `xml:"crossmark_domain_exclusive"`
	CustomMetadata           struct {
		Text      string `xml:",chardata"`
		Assertion []struct {
			Text       string `xml:",chardata"`
			Name       string `xml:"name,attr"`
			Label      string `xml:"label,attr"`
			GroupName  string `xml:"group_name,attr"`
			GroupLabel string `xml:"group_label,attr"`
			Order      string `xml:"order,attr"`
		} `xml:"assertion"`
		Program Program `xml:"program"`
	} `xml:"custom_metadata"`
}

type Dissertation struct {
	XMLName    xml.Name     `xml:"dissertation"`
	PersonName []PersonName `xml:"person_name"`
	Titles     struct {
		Text     string `xml:",chardata"`
		Title    string `xml:"title"`
		Subtitle string `xml:"subtitle"`
	} `xml:"titles"`
	Institution Institution `xml:"institution"`
	Degree      string      `xml:"degree"`
	DOIData     DOIData     `xml:"doi_data"`
}

type DOIData struct {
	XMLName    xml.Name   `xml:"doi_data"`
	DOI        string     `xml:"doi"`
	Timestamp  string     `xml:"timestamp"`
	Resource   string     `xml:"resource"`
	Collection Collection `xml:"collection"`
}

type Institution struct {
	XMLName          xml.Name `xml:"institution"`
	InstitutionName  string   `xml:"institution_name"`
	InstitutionPlace string   `xml:"institution_place"`
	InstitutionID    struct {
		Text   string `xml:",chardata"`
		IDType string `xml:"id_type,attr"`
	} `xml:"institution_id"`
}

// ISSN represents a ISSN in Crossref XML metadata.
type ISSN struct {
	XMLName   xml.Name `xml:"issn"`
	MediaType string   `xml:"media_type,attr"`
	Text      string   `xml:",chardata"`
}

type Item struct {
	XMLName  xml.Name `xml:"item"`
	Crawler  string   `xml:"crawler,attr"`
	Resource Resource `xml:"resource"`
}

// ItemNumber represents an item number in Crossref XML metadata.
type ItemNumber struct {
	XMLName          xml.Name `xml:"item_number"`
	ItemNumberType   string   `xml:"item_number_type,attr"`
	ItemNumberString string   `xml:",innerxml"`
}

// Journal represents a journal in Crossref XML metadata.
type Journal struct {
	XMLName         xml.Name        `xml:"journal"`
	JournalIssue    JournalIssue    `xml:"journal_issue"`
	JournalMetadata JournalMetadata `xml:"journal_metadata"`
	JournalArticle  JournalArticle  `xml:"journal_article"`
}

// JournalArticle represents a journal article in Crossref XML metadata.
type JournalArticle struct {
	XMLName                   xml.Name         `xml:"journal_article"`
	Text                      string           `xml:",chardata"`
	PublicationType           string           `xml:"publication_type,attr"`
	ReferenceDistributionOpts string           `xml:"reference_distribution_opts,attr"`
	Titles                    []Title          `xml:"titles>title"`
	Contributors              []PersonName     `xml:"contributors>person_name"`
	PublicationDate           PublicationDate  `xml:"publication_date"`
	Abstract                  Abstract         `xml:"abstract"`
	ISSN                      string           `xml:"issn"`
	ItemNumber                ItemNumber       `xml:"item_number"`
	Crossmark                 Crossmark        `xml:"crossmark"`
	ArchiveLocations          ArchiveLocations `xml:"archive_locations"`
	DOIData                   DOIData          `xml:"doi_data"`
	CitationList              struct {
		Text     string     `xml:",chardata"`
		Citation []Citation `xml:"citation"`
	} `xml:"citation_list"`
}

type JournalIssue struct {
	Text            string          `xml:",chardata"`
	PublicationDate PublicationDate `xml:"publication_date"`
	JournalVolume   JournalVolume   `xml:"journal_volume"`
}

// JournalMetadata represents journal metadata in Crossref XML metadata.
type JournalMetadata struct {
	XMLName   xml.Name `xml:"journal_metadata"`
	Language  string   `xml:"language,attr"`
	FullTitle string   `xml:"full_title"`
	ISSN      ISSN     `xml:"issn"`
}

type JournalVolume struct {
	Text   string `xml:",chardata"`
	Volume string `xml:"volume"`
}

type PeerReview struct {
	XMLName      xml.Name     `xml:"peer_review"`
	Stage        string       `xml:"stage,attr"`
	Type         string       `xml:"type,attr"`
	Contributors []PersonName `xml:"contributors>person_name"`
	Titles       []Title      `xml:"titles>title"`
	DOIData      DOIData      `xml:"doi_data"`
}

// PersonName represents a person in Crossref XML metadata.
type PersonName struct {
	XMLName         xml.Name      `xml:"person_name"`
	ContributorRole string        `xml:"contributor_role,attr"`
	Sequence        string        `xml:"sequence,attr"`
	Text            string        `xml:",chardata"`
	GivenName       string        `xml:"given_name"`
	Surname         string        `xml:"surname"`
	ORCID           string        `xml:"ORCID"`
	Affiliations    []Institution `xml:"affiliations>institution"`
}

// PostedContent represents posted content in Crossref XML metadata.
type PostedContent struct {
	XMLName        xml.Name       `xml:"posted_content"`
	Type           string         `xml:"type,attr"`
	Language       string         `xml:"language,attr"`
	GroupTitle     string         `xml:"group_title"`
	Contributors   []PersonName   `xml:"contributors>person_name"`
	Titles         []Title        `xml:"titles>title"`
	PostedDate     PostedDate     `xml:"posted_date"`
	AcceptanceDate AcceptanceDate `xml:"acceptance_date"`
	Institution    Institution    `xml:"institution"`
	ItemNumber     ItemNumber     `xml:"item_number"`
	DOIData        DOIData        `xml:"doi_data"`
}

type PostedDate struct {
	XMLName   xml.Name `xml:"posted_date"`
	MediaType string   `xml:"media_type,attr"`
	Month     string   `xml:"month"`
	Day       string   `xml:"day"`
	Year      string   `xml:"year"`
}

type Program struct {
	XMLName    xml.Name    `xml:"program"`
	Text       string      `xml:",chardata"`
	Fr         string      `xml:"fr,attr"`
	Name       string      `xml:"name,attr"`
	Ai         string      `xml:"ai,attr"`
	Rel        string      `xml:"rel,attr"`
	Assertion  []Assertion `xml:"assertion"`
	LicenseRef []struct {
		Text      string `xml:",chardata"`
		AppliesTo string `xml:"applies_to,attr"`
	} `xml:"license_ref"`
	RelatedItem struct {
		Text              string `xml:",chardata"`
		Description       string `xml:"description"`
		InterWorkRelation struct {
			Text             string `xml:",chardata"`
			IdentifierType   string `xml:"identifier-type,attr"`
			RelationshipType string `xml:"relationship-type,attr"`
		} `xml:"inter_work_relation"`
	} `xml:"related_item"`
}

type PublicationDate struct {
	XMLName   xml.Name `xml:"publication_date"`
	Text      string   `xml:",chardata"`
	MediaType string   `xml:"media_type,attr"`
	Month     string   `xml:"month"`
	Day       string   `xml:"day"`
	Year      string   `xml:"year"`
}

type PublisherItem struct {
	XMLName    xml.Name `xml:"publisher_item"`
	Text       string   `xml:",chardata"`
	ItemNumber struct {
		Text           string `xml:",chardata"` // e01567
		ItemNumberType string `xml:"item_number_type,attr"`
	} `xml:"item_number"`
	Identifier struct {
		Text   string `xml:",chardata"` // 10.7554/eLife.01567
		IDType string `xml:"id_type,attr"`
	} `xml:"identifier"`
}

// Resource represents a resource in Crossref XML metadata.
type Resource struct {
	XMLName  xml.Name `xml:"resource"`
	Text     string   `xml:",chardata"`
	MimeType string   `xml:"mime_type,attr"`
}

type SAComponent struct {
	XMLName       xml.Name      `xml:"sa_component"`
	ComponentList ComponentList `xml:"component_list"`
}

// Title represents the title in Crossref XML metadata.
type Title struct {
	XMLName xml.Name `xml:"title"`
	Text    string   `xml:",chardata"`
}

// CRToCMMappings maps Crossref XML types to Commonmeta types
// source: http://api.crossref.org/types,
// Crossref XML naming conventions are different from the REST API
var CRToCMMappings = map[string]string{
	"book_chapter":        "BookChapter",
	"book_part":           "BookPart",
	"book_section":        "BookSection",
	"book_series":         "BookSeries",
	"book_set":            "BookSet",
	"book_track":          "BookTrack",
	"book":                "Book",
	"component":           "Component",
	"database":            "Database",
	"dataset":             "Dataset",
	"dissertation":        "Dissertation",
	"edited_book":         "Book",
	"grant":               "Grant",
	"journal_article":     "JournalArticle",
	"journal_issue":       "JournalIssue",
	"journal_volume":      "JournalVolume",
	"journal":             "Journal",
	"monograph":           "Book",
	"other":               "Other",
	"peer_review":         "PeerReview",
	"posted_content":      "Article",
	"proceedings_article": "ProceedingsArticle",
	"proceedings_series":  "ProceedingsSeries",
	"proceedings":         "Proceedings",
	"reference_book":      "Book",
	"reference_entry":     "Entry",
	"report_component":    "ReportComponent",
	"report_series":       "ReportSeries",
	"report":              "Report",
	"standard":            "Standard",
}

// Fetch gets the metadata for a single work from the Crossref API and converts it to the Commonmeta format
func Fetch(str string) (commonmeta.Data, error) {
	var data commonmeta.Data
	id, ok := doiutils.ValidateDOI(str)
	if !ok {
		return data, errors.New("invalid DOI")
	}
	content, err := Get(id)
	if err != nil {
		return data, err
	}
	data, err = Read(content)
	if err != nil {
		return data, err
	}
	return data, nil
}

// Get gets the metadata for a single work from the Crossref API in Crossref XML format.
func Get(pid string) (Content, error) {

	// the envelope for the XML response from the Crossref API
	type CrossrefResult struct {
		XMLName        xml.Name `xml:"crossref_result"`
		Xmlns          string   `xml:"xmlns,attr"`
		Version        string   `xml:"version,attr"`
		Xsi            string   `xml:"xsi,attr"`
		SchemaLocation string   `xml:"schemaLocation,attr"`
		QueryResult    struct {
			Head struct {
				DoiBatchID string `xml:"doi_batch_id"`
			} `xml:"head"`
			Body Content `xml:"body"`
		} `xml:"query_result"`
	}

	crossrefResult := CrossrefResult{
		Xmlns:          "http://www.crossref.org/qrschema/3.0",
		Version:        "3.0",
		Xsi:            "http://www.w3.org/2001/XMLSchema-instance",
		SchemaLocation: "http://www.crossref.org/qrschema/3.0 http://www.crossref.org/qrschema/crossref_query_output3.0.xsd",
	}
	doi, ok := doiutils.ValidateDOI(pid)
	if !ok {
		return crossrefResult.QueryResult.Body, errors.New("invalid DOI")
	}
	url := "https://api.crossref.org/works/" + doi + "/transform/application/vnd.crossref.unixsd+xml"
	req, err := http.NewRequest("GET", url, nil)
	v := "0.1"
	u := "info@front-matter.io"
	userAgent := fmt.Sprintf("commonmeta/%s (https://commonmeta.org/; mailto: %s)", v, u)
	req.Header.Set("User-Agent", userAgent)
	if err != nil {
		log.Fatalln(err)
	}
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return crossrefResult.QueryResult.Body, err
	}
	if resp.StatusCode >= 400 {
		return crossrefResult.QueryResult.Body, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error:", err)
		return crossrefResult.QueryResult.Body, err
	}
	err = xml.Unmarshal(body, &crossrefResult)
	if err != nil {
		fmt.Println("error:", err)
	}
	log.Println(crossrefResult.QueryResult.Body.Query.DOIRecord.Crossref)
	return crossrefResult.QueryResult.Body, err
}

// Read Crossref JSON response and return work struct in Commonmeta format
func Read(content Content) (commonmeta.Data, error) {
	var data = commonmeta.Data{}

	var contributors []PersonName
	var doiData DOIData

	meta := content.Query.DOIRecord.Crossref

	data.ID = doiutils.NormalizeDOI(content.Query.DOI.Text)
	data.Type = CRToCMMappings[content.Query.DOI.Type]
	if data.Type == "" {
		data.Type = "Other"
	}

	// fetch metadata depending of Crossref type
	switch data.Type {
	case "JournalArticle":
		doiData = meta.Journal.JournalArticle.DOIData
		contributors = meta.Journal.JournalArticle.Contributors
		data.Language = meta.Journal.JournalMetadata.Language
	case "JournalIssue":
		// doiData = meta.Journal.JournalIssue.JournalVolume.DOIData
		data.Language = meta.Journal.JournalMetadata.Language
	case "Journal":
		// doiData = meta.Journal.JournalMetadata
		data.Language = meta.Journal.JournalMetadata.Language
	case "Article":
		doiData = meta.PostedContent.DOIData
		contributors = meta.PostedContent.Contributors
		data.Language = meta.PostedContent.Language
	case "BookChapter":
		data.Language = ""
		data.URL = ""
	case "BookSeries":
		data.Language = ""
		data.URL = ""
	case "BookSet":
		data.Language = ""
		data.URL = ""
	case "Book":
		doiData = meta.Book.BookMetadata.DOIData
		data.Language = meta.Book.BookMetadata.Language
	case "ProceedingsArticle":
		data.Language = ""
		data.URL = ""
	case "Component":
		data.Language = ""
		data.URL = ""
	case "Dataset":
		data.Language = ""
		data.URL = ""
	case "Report":
		data.Language = ""
		data.URL = ""
	case "PeerReview":
		doiData = meta.PeerReview.DOIData
		data.Language = ""
	case "Dissertation":
		doiData = meta.Dissertation.DOIData
		data.Language = ""
	}

	// containerType := crossref.CrossrefContainerTypes[content.Type]
	// containerType = crossref.CRToCMContainerTranslations[containerType]

	// for _, v := range content.Archive {
	// 	if !slices.Contains(data.ArchiveLocations, v) {
	// 		data.ArchiveLocations = append(data.ArchiveLocations, v)
	// 	}
	// }

	// var identifier, identifierType string
	// if len(content.ISSNType) > 0 {
	// 	i := make(map[string]string)
	// 	for _, issn := range content.ISSNType {
	// 		i[issn.Type] = issn.Value
	// 	}
	// 	if i["electronic"] != "" {
	// 		identifier = i["electronic"]
	// 		identifierType = "ISSN"
	// 	} else if i["print"] != "" {
	// 		identifier = i["print"]
	// 		identifierType = "ISSN"
	// 	}
	// } else if len(content.ISBNType) > 0 {
	// 	i := make(map[string]string)
	// 	for _, isbn := range content.ISBNType {
	// 		i[isbn.Type] = isbn.Value
	// 	}
	// 	if i["electronic"] != "" {
	// 		identifier = i["electronic"]
	// 		identifierType = "ISBN"
	// 	} else if i["print"] != "" {
	// 		identifier = i["print"]
	// 		identifierType = "ISBN"
	// 	}
	// }
	// var containerTitle string
	// if len(content.ContainerTitle) > 0 {
	// 	containerTitle = content.ContainerTitle[0]
	// }
	// var lastPage string
	// pages := strings.Split(content.Page, "-")
	// firstPage := pages[0]
	// if len(pages) > 1 {
	// 	lastPage = pages[1]
	// }

	// data.Container = commonmeta.Container{
	// 	Identifier:     identifier,
	// 	IdentifierType: identifierType,
	// 	Type:           containerType,
	// 	Title:          containerTitle,
	// 	Volume:         content.Volume,
	// 	Issue:          content.Issue,
	// 	FirstPage:      firstPage,
	// 	LastPage:       lastPage,
	// }

	if len(contributors) > 0 {
		for _, v := range contributors {
			var ID string
			if v.GivenName != "" || v.Surname != "" {
				if v.ORCID != "" {
					// enforce HTTPS
					ID, _ = utils.NormalizeURL(v.ORCID, true, false)
				}
			}
			Type := "Person"
			var affiliations []commonmeta.Affiliation
			if len(v.Affiliations) > 0 {
				for _, a := range v.Affiliations {
					var ID string
					if a.InstitutionID.Text != "" {
						ID = utils.NormalizeROR(a.InstitutionID.Text)
					}
					if a.InstitutionName != "" {
						affiliations = append(affiliations, commonmeta.Affiliation{
							ID:   ID,
							Name: a.InstitutionName,
						})
					}
				}
			}

			contributor := commonmeta.Contributor{
				ID:               ID,
				Type:             Type,
				GivenName:        v.GivenName,
				FamilyName:       v.Surname,
				Name:             "",
				ContributorRoles: []string{"Author"},
				Affiliations:     affiliations,
			}
			containsName := slices.ContainsFunc(data.Contributors, func(e commonmeta.Contributor) bool {
				return e.GivenName == contributor.GivenName && e.FamilyName != "" && e.FamilyName == contributor.FamilyName
			})
			if !containsName {
				data.Contributors = append(data.Contributors, contributor)
			}
		}
	}

	// if content.Abstract != "" {
	// 	abstract := utils.Sanitize(content.Abstract)
	// 	data.Descriptions = append(data.Descriptions, commonmeta.Description{
	// 		Description: abstract,
	// 		Type:        "Abstract",
	// 	})
	// }

	// for _, v := range content.Link {
	// 	if v.ContentType != "unspecified" {
	// 		data.Files = append(data.Files, commonmeta.File{
	// 			URL:      v.URL,
	// 			MimeType: v.ContentType,
	// 		})
	// 	}
	// }
	// if len(content.Link) > 1 {
	// 	data.Files = utils.DedupeSlice(data.Files)
	// }

	// if len(content.Funder) > 1 {
	// 	for _, v := range content.Funder {
	// 		funderIdentifier := doiutils.NormalizeDOI(v.DOI)
	// 		var funderIdentifierType string
	// 		if strings.HasPrefix(v.DOI, "10.13039") {
	// 			funderIdentifierType = "Crossref Funder ID"
	// 		}
	// 		if len(v.Award) > 0 {
	// 			for _, award := range v.Award {
	// 				data.FundingReferences = append(data.FundingReferences, commonmeta.FundingReference{
	// 					FunderIdentifier:     funderIdentifier,
	// 					FunderIdentifierType: funderIdentifierType,
	// 					FunderName:           v.Name,
	// 					AwardNumber:          award,
	// 				})
	// 			}
	// 		} else {
	// 			data.FundingReferences = append(data.FundingReferences, commonmeta.FundingReference{
	// 				FunderIdentifier:     funderIdentifier,
	// 				FunderIdentifierType: funderIdentifierType,
	// 				FunderName:           v.Name,
	// 			})
	// 		}
	// 	}
	// 	// if len(content.Funder) > 1 {
	// 	data.FundingReferences = utils.DedupeSlice(data.FundingReferences)
	// 	// }
	// }

	// data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
	// 	Identifier:     data.ID,
	// 	IdentifierType: "DOI",
	// })

	// 	license_ = (
	// 		py_.get(bibmeta, "program.0.license_ref")
	// 		or py_.get(bibmeta, "crossmark.custom_metadata.program.0.license_ref")
	// 		or py_.get(bibmeta, "crossmark.custom_metadata.program.1.license_ref")
	// )
	// if content.License != nil && len(content.License) > 0 {
	// 	url, _ := utils.NormalizeCCUrl(content.License[0].URL)
	// 	id := utils.URLToSPDX(url)
	// 	data.License = commonmeta.License{
	// 		ID:  id,
	// 		URL: url,
	// 	}
	// }

	data.Provider = "Crossref"

	var publisherID, publisherName string
	for _, v := range content.Query.CRMItem {
		if v.Name == "member-id" {
			memberID := v.Text
			publisherID = "https://api.crossref.org/members/" + memberID
		} else if v.Name == "publisher-name" {
			publisherName = v.Text
		} else if v.Name == "created" {
			data.Date.Created = v.Text
		} else if v.Name == "last-update" {
			data.Date.Updated = v.Text
		}
	}
	data.Publisher = commonmeta.Publisher{
		ID:   publisherID,
		Name: publisherName,
	}

	// for _, v := range content.Reference {
	// 	reference := commonmeta.Reference{
	// 		Key:             v.Key,
	// 		ID:              doiutils.NormalizeDOI(v.DOI),
	// 		Title:           v.ArticleTitle,
	// 		PublicationYear: v.Year,
	// 		Unstructured:    v.Unstructured,
	// 	}
	// 	containsKey := slices.ContainsFunc(data.References, func(e commonmeta.Reference) bool {
	// 		return e.Key != "" && e.Key == reference.Key
	// 	})
	// 	if !containsKey {
	// 		data.References = append(data.References, reference)
	// 	}
	// }

	// fields := reflect.VisibleFields(reflect.TypeOf(content.Relation))
	// for _, field := range fields {
	// 	if slices.Contains(relationTypes, field.Name) {
	// 		relationByType := reflect.ValueOf(content.Relation).FieldByName(field.Name)
	// 		for _, v := range relationByType.Interface().([]struct {
	// 			ID     string `json:"id"`
	// 			IDType string `json:"id-type"`
	// 		}) {
	// 			var id string
	// 			if v.IDType == "doi" {
	// 				id = doiutils.NormalizeDOI(v.ID)
	// 			} else if v.IDType == "issn" {
	// 				id = utils.ISSNAsURL(v.ID)
	// 			} else if utils.ValidateURL(v.ID) == "URL" {
	// 				id = v.ID
	// 			}
	// 			relation := commonmeta.Relation{
	// 				ID:   id,
	// 				Type: field.Name,
	// 			}
	// 			if id != "" && !slices.Contains(data.Relations, relation) {
	// 				data.Relations = append(data.Relations, relation)
	// 			}
	// 		}
	// 		sort.Slice(data.Relations, func(i, j int) bool {
	// 			return data.Relations[i].Type < data.Relations[j].Type
	// 		})
	// 	}
	// }
	// if data.Container.IdentifierType == "ISSN" {
	// 	data.Relations = append(data.Relations, commonmeta.Relation{
	// 		ID:   utils.ISSNAsURL(data.Container.Identifier),
	// 		Type: "IsPartOf",
	// 	})
	// }

	// for _, v := range content.Subject {
	// 	subject := commonmeta.Subject{
	// 		Subject: v,
	// 	}
	// 	if !slices.Contains(data.Subjects, subject) {
	// 		data.Subjects = append(data.Subjects, subject)
	// 	}
	// }

	// if content.GroupTitle != "" {
	// 	data.Subjects = append(data.Subjects, commonmeta.Subject{
	// 		Subject: content.GroupTitle,
	// 	})
	// }

	// if len(content.Title) > 0 && content.Title[0] != "" {
	// 	data.Titles = append(data.Titles, commonmeta.Title{
	// 		Title: content.Title[0],
	// 	})
	// }
	// if len(content.Subtitle) > 0 && content.Subtitle[0] != "" {
	// 	data.Titles = append(data.Titles, commonmeta.Title{
	// 		Title: content.Subtitle[0],
	// 		Type:  "Subtitle",
	// 	})
	// }
	// if len(content.OriginalTitle) > 0 && content.OriginalTitle[0] != "" {
	// 	data.Titles = append(data.Titles, commonmeta.Title{
	// 		Title: content.OriginalTitle[0],
	// 		Type:  "TranslatedTitle",
	// 	})
	// }

	data.URL = doiData.Resource

	return data, nil
}

// def generate_crossref_xml(metadata: Commonmeta) -> Optional[str]:
//     """Generate Crossref XML. First checks for write errors (JSON schema validation)"""
//     xml = crossref_root()
//     head = etree.SubElement(xml, "head")
//     # we use a uuid as batch_id
//     etree.SubElement(head, "doi_batch_id").text = str(uuid.uuid4())
//     etree.SubElement(head, "timestamp").text = datetime.now().strftime("%Y%m%d%H%M%S")
//     depositor = etree.SubElement(head, "depositor")
//     etree.SubElement(depositor, "depositor_name").text = metadata.depositor
//     etree.SubElement(depositor, "email_address").text = metadata.email
//     etree.SubElement(head, "registrant").text = metadata.registrant

//     body = etree.SubElement(xml, "body")
//     body = insert_crossref_work(metadata, body)
//     return etree.tostring(
//         xml,
//         doctype='<?xml version="1.0" encoding="UTF-8"?>',
//         pretty_print=True,
//     )

// def insert_crossref_work(metadata, xml):
//     """Insert crossref work"""
//     if metadata.type not in ["JournalArticle", "Article"]:
//         return xml
//     if doi_from_url(metadata.id) is None or metadata.url is None:
//         return xml
//     if metadata.type == "JournalArticle":
//         xml = insert_journal(metadata, xml)
//     elif metadata.type == "Article":
//         xml = insert_posted_content(metadata, xml)

// def insert_journal(metadata, xml):
//     """Insert journal"""
//     journal = etree.SubElement(xml, "journal")
//     if metadata.language is not None:
//         journal_metadata = etree.SubElement(
//             journal, "journal_metadata", {"language": metadata.language[:2]}
//         )
//     else:
//         journal_metadata = etree.SubElement(journal, "journal_metadata")
//     if (
//         metadata.container is not None
//         and metadata.container.get("title", None) is not None
//     ):
//         etree.SubElement(journal_metadata, "full_title").text = metadata.container.get(
//             "title"
//         )
//     journal_metadata = insert_group_title(metadata, journal_metadata)
//     journal_article = etree.SubElement(
//         journal, "journal_article", {"publication_type": "full_text"}
//     )
//     journal_article = insert_crossref_titles(metadata, journal_article)
//     journal_article = insert_crossref_contributors(metadata, journal_article)
//     journal_article = insert_crossref_publication_date(metadata, journal_article)
//     journal_article = insert_crossref_abstract(metadata, journal_article)
//     journal_article = insert_crossref_issn(metadata, journal_article)
//     journal_article = insert_item_number(metadata, journal_article)
//     journal_article = insert_funding_references(metadata, journal_article)
//     journal_article = insert_crossref_access_indicators(metadata, journal_article)
//     journal_article = insert_crossref_relations(metadata, journal_article)
//     journal_article = insert_archive_locations(metadata, journal_article)
//     journal_article = insert_doi_data(metadata, journal_article)
//     journal_article = insert_citation_list(metadata, journal_article)

//     return journal

// def insert_posted_content(metadata, xml):
//     """Insert posted content"""
//     if metadata.language is not None:
//         posted_content = etree.SubElement(
//             xml, "posted_content", {"type": "other", "language": metadata.language[:2]}
//         )
//     else:
//         posted_content = etree.SubElement(xml, "posted_content", {"type": "other"})

//     posted_content = insert_group_title(metadata, posted_content)
//     posted_content = insert_crossref_contributors(metadata, posted_content)
//     posted_content = insert_crossref_titles(metadata, posted_content)
//     posted_content = insert_posted_date(metadata, posted_content)
//     posted_content = insert_institution(metadata, posted_content)
//     posted_content = insert_item_number(metadata, posted_content)
//     posted_content = insert_crossref_abstract(metadata, posted_content)
//     posted_content = insert_funding_references(metadata, posted_content)
//     posted_content = insert_crossref_access_indicators(metadata, posted_content)
//     posted_content = insert_crossref_relations(metadata, posted_content)
//     posted_content = insert_archive_locations(metadata, posted_content)
//     posted_content = insert_doi_data(metadata, posted_content)
//     posted_content = insert_citation_list(metadata, posted_content)

//     return xml

// def insert_group_title(metadata, xml):
//     """Insert group title"""
//     if metadata.subjects is None or len(metadata.subjects) == 0:
//         return xml
//     etree.SubElement(xml, "group_title").text = metadata.subjects[0].get(
//         "subject", None
//     )
//     return xml

// def insert_crossref_contributors(metadata, xml):
//     """Insert crossref contributors"""
//     if metadata.contributors is None or len(metadata.contributors) == 0:
//         return xml
//     contributors = etree.SubElement(xml, "contributors")
//     con = [
//         c
//         for c in metadata.contributors
//         if c.get("contributorRoles", None) == ["Author"]
//         or c.get("contributorRoles", None) == ["Editor"]
//     ]
//     for num, contributor in enumerate(con):
//         contributor_role = (
//             "author" if "Author" in contributor.get("contributorRoles") else None
//         )
//         if contributor_role is None:
//             contributor_role = (
//                 "editor" if "Editor" in contributor.get("contributorRoles") else None
//             )
//         sequence = "first" if num == 0 else "additional"
//         if (
//             contributor.get("type", None) == "Organization"
//             and contributor.get("name", None) is not None
//         ):
//             etree.SubElement(
//                 contributors,
//                 "organization",
//                 {"contributor_role": contributor_role, "sequence": sequence},
//             ).text = contributor.get("name")
//         elif (
//             contributor.get("givenName", None) is not None
//             or contributor.get("familyName", None) is not None
//         ):
//             person_name = etree.SubElement(
//                 contributors,
//                 "person_name",
//                 {"contributor_role": contributor_role, "sequence": sequence},
//             )
//             person_name = insert_crossref_person(contributor, person_name)
//         elif contributor.get("affiliations", None) is not None:
//             anonymous = etree.SubElement(
//                 contributors,
//                 "anonymous",
//                 {"contributor_role": contributor_role, "sequence": sequence},
//             )
//             anonymous = insert_crossref_anonymous(contributor, anonymous)
//         else:
//             etree.SubElement(
//                 contributors,
//                 "anonymous",
//                 {"contributor_role": contributor_role, "sequence": sequence},
//             )
//     return xml

// def insert_crossref_person(contributor, xml):
//     """Insert crossref person"""
//     if contributor.get("givenName", None) is not None:
//         etree.SubElement(xml, "given_name").text = contributor.get("givenName")
//     if contributor.get("familyName", None) is not None:
//         etree.SubElement(xml, "surname").text = contributor.get("familyName")

//     if contributor.get("affiliations", None) is not None:
//         affiliations = etree.SubElement(xml, "affiliations")
//         institution = etree.SubElement(affiliations, "institution")
//         if py_.get(contributor, "affiliations.0.name") is not None:
//             etree.SubElement(institution, "institution_name").text = py_.get(
//                 contributor, "affiliations.0.name"
//             )
//         if py_.get(contributor, "affiliations.0.id") is not None:
//             etree.SubElement(
//                 institution, "institution_id", {"type": "ror"}
//             ).text = py_.get(contributor, "affiliations.0.id")
//     orcid = normalize_orcid(contributor.get("id", None))
//     if orcid is not None:
//         etree.SubElement(xml, "ORCID").text = orcid
//     return xml

// def insert_crossref_anonymous(contributor, xml):
//     """Insert crossref anonymous"""
//     if contributor.get("affiliations", None) is None:
//         return xml
//     affiliations = etree.SubElement(xml, "affiliations")
//     institution = etree.SubElement(affiliations, "institution")
//     if py_.get(contributor, "affiliations.0.name") is not None:
//         etree.SubElement(institution, "institution_name").text = py_.get(
//             contributor, "affiliations.0.name"
//         )
//     return xml

// def insert_crossref_titles(metadata, xml):
//     """Insert crossref titles"""
//     titles = etree.SubElement(xml, "titles")
//     for title in wrap(metadata.titles):
//         if isinstance(title, dict):
//             etree.SubElement(titles, "title").text = title.get("title", None)
//         else:
//             etree.SubElement(titles, "title").text = title
//     return xml

// def insert_citation_list(metadata, xml):
//     """Insert citation list"""
//     if metadata.references is None or len(metadata.references) == 0:
//         return xml

//     citation_list = etree.SubElement(xml, "citation_list")
//     for ref in metadata.references:
//         citation = etree.SubElement(
//             citation_list, "citation", {"key": ref.get("key", None)}
//         )
//         if ref.get("journal_title", None) is not None:
//             etree.SubElement(citation, "journal_article").text = ref.get(
//                 "journal_title"
//             )
//         if ref.get("author", None) is not None:
//             etree.SubElement(citation, "author").text = ref.get("author")
//         if ref.get("volume", None) is not None:
//             etree.SubElement(citation, "volume").text = ref.get("volume")
//         if ref.get("first_page", None) is not None:
//             etree.SubElement(citation, "first_page").text = ref.get("first_page")
//         if ref.get("publicationYear", None) is not None:
//             etree.SubElement(citation, "cYear").text = ref.get("publicationYear")
//         if ref.get("title", None) is not None:
//             etree.SubElement(citation, "article_title").text = ref.get("title")
//         if ref.get("doi", None) is not None:
//             etree.SubElement(citation, "doi").text = doi_from_url(ref.get("doi"))
//         if ref.get("url", None) is not None:
//             etree.SubElement(citation, "unstructured_citation").text = ref.get("url")
//     return xml

// def insert_crossref_access_indicators(metadata, xml):
//     """Insert crossref access indicators"""
//     rights_uri = (
//         metadata.license.get("url", None) if metadata.license is not None else None
//     )
//     if rights_uri is None:
//         return xml
//     program = etree.SubElement(
//         xml,
//         "program",
//         {
//             "xmlns": "http://www.crossref.org/AccessIndicators.xsd",
//             "name": "AccessIndicators",
//         },
//     )
//     etree.SubElement(program, "license_ref", {"applies_to": "vor"}).text = rights_uri
//     etree.SubElement(program, "license_ref", {"applies_to": "tdm"}).text = rights_uri
//     return xml

// def insert_crossref_relations(metadata, xml):
//     """Insert crossref relations"""
//     if metadata.relations is None or len(metadata.relations) == 0:
//         return xml
//     program = etree.SubElement(
//         xml,
//         "program",
//         {
//             "xmlns": "http://www.crossref.org/relations.xsd",
//             "name": "relations",
//         },
//     )
//     for relation in metadata.relations:
//         if relation.get("type", None) in [
//             "IsPartOf",
//             "HasPart",
//             "IsReviewOf",
//             "HasReview",
//             "IsRelatedMaterial",
//             "HasRelatedMaterial",
//         ]:
//             group = "inter_work_relation"
//         elif relation.get("type", None) in [
//             "IsIdenticalTo",
//             "IsPreprintOf",
//             "HasPreprint",
//             "IsTranslationOf",
//             "HasTranslation",
//             "IsVersionOf",
//             "HasVersion",
//         ]:
//             group = "intra_work_relation"
//         else:
//             continue

//         related_item = etree.SubElement(program, "related_item")
//         f = furl(relation.get("id", None))
//         if validate_doi(relation.get("id", None)):
//             identifier_type = "doi"
//             _id = doi_from_url(relation.get("id", None))
//         elif f.host == "portal.issn.org":
//             identifier_type = "issn"
//             _id = f.path.segments[-1]
//         elif validate_url(relation.get("id", None)) == "URL":
//             identifier_type = "uri"
//             _id = relation.get("id", None)
//         else:
//             identifier_type = "other"
//             _id = relation.get("id", None)

//         etree.SubElement(
//             related_item,
//             group,
//             {
//                 "relationship-type": py_.lower_first(relation.get("type"))
//                 if relation.get("type", None) is not None
//                 else None,
//                 "identifier-type": identifier_type,
//             },
//         ).text = _id

//     return xml

// def insert_funding_references(metadata, xml):
//     """Insert funding references"""
//     if metadata.funding_references is None or len(metadata.funding_references) == 0:
//         return xml
//     program = etree.SubElement(
//         xml,
//         "program",
//         {
//             "xmlns": "http://www.crossref.org/fundref.xsd",
//             "name": "fundref",
//         },
//     )
//     for funding_reference in metadata.funding_references:
//         assertion = etree.SubElement(program, "assertion", {"name": "fundgroup"})
//         funder_name = etree.SubElement(
//             assertion,
//             "assertion",
//             {"name": "funder_name"},
//         )
//         if funding_reference.get("funderIdentifier", None) is not None:
//             etree.SubElement(
//                 funder_name,
//                 "assertion",
//                 {"name": "funder_identifier"},
//             ).text = funding_reference.get("funderIdentifier", None)
//         if funding_reference.get("awardNumber", None) is not None:
//             etree.SubElement(
//                 assertion,
//                 "assertion",
//                 {"name": "award_number"},
//             ).text = funding_reference.get("awardNumber", None)
//         funder_name.text = funding_reference.get("funderName", None)
//     return xml

// def insert_crossref_subjects(metadata, xml):
//     """Insert crossref subjects"""
//     if metadata.subjects is None:
//         return xml
//     subjects = etree.SubElement(xml, "subjects")
//     for subject in metadata.subjects:
//         if isinstance(subject, dict):
//             etree.SubElement(subjects, "subject").text = subject.get("subject", None)
//         else:
//             etree.SubElement(subjects, "subject").text = subject
//     return xml

// def insert_crossref_language(metadata, xml):
//     """Insert crossref language"""
//     if metadata.language is None:
//         return xml
//     etree.SubElement(xml, "language").text = metadata.language
//     return xml

// def insert_crossref_publication_date(metadata, xml):
//     """Insert crossref publication date"""
//     pub_date = parse(metadata.date.get("published", None))
//     if pub_date is None:
//         return xml

//     publication_date = etree.SubElement(
//         xml, "publication_date", {"media_type": "online"}
//     )
//     etree.SubElement(publication_date, "month").text = f"{pub_date.month:d}"
//     etree.SubElement(publication_date, "day").text = f"{pub_date.day:d}"
//     etree.SubElement(publication_date, "year").text = str(pub_date.year)
//     return xml

// def insert_posted_date(metadata, xml):
//     """Insert posted date"""
//     pub_date = parse(metadata.date.get("published", None))
//     if pub_date is None:
//         return xml

//     posted_date = etree.SubElement(xml, "posted_date", {"media_type": "online"})
//     etree.SubElement(posted_date, "month").text = f"{pub_date.month:d}"
//     etree.SubElement(posted_date, "day").text = f"{pub_date.day:d}"
//     etree.SubElement(posted_date, "year").text = str(pub_date.year)
//     return xml

// def insert_institution(metadata, xml):
//     """Insert institution"""
//     if metadata.publisher.get("name", None) is None:
//         return xml
//     institution = etree.SubElement(xml, "institution")
//     etree.SubElement(institution, "institution_name").text = metadata.publisher.get(
//         "name"
//     )
//     return xml

// def insert_item_number(metadata, xml):
//     """Insert item number"""
//     if metadata.identifiers is None:
//         return xml
//     for identifier in metadata.identifiers:
//         if identifier.get("identifier", None) is None:
//             continue
//         if identifier.get("identifierType", None) is not None:
//             # strip hyphen from UUIDs, as item_number can only be 32 characters long (UUIDv4 is 36 characters long)
//             if identifier.get("identifierType", None) == "UUID":
//                 identifier["identifier"] = identifier.get(
//                     "identifier", ""
//                 ).replace("-", "")
//             etree.SubElement(
//                 xml,
//                 "item_number",
//                 {
//                     "item_number_type": identifier.get(
//                         "identifierType", ""
//                     ).lower()
//                 },
//             ).text = identifier.get("identifier", None)
//         else:
//             etree.SubElement(xml, "item_number").text = identifier.get(
//                 "identifier", None
//             )
//     return xml

// def insert_archive_locations(metadata, xml):
//     """Insert archive locations"""
//     if metadata.archive_locations is None:
//         return xml
//     archive_locations = etree.SubElement(xml, "archive_locations")
//     for archive_location in metadata.archive_locations:
//         etree.SubElement(archive_locations, "archive", {"name": archive_location})
//     return xml

// def insert_doi_data(metadata, xml):
//     """Insert doi data"""
//     if doi_from_url(metadata.id) is None or metadata.url is None:
//         return xml
//     doi_data = etree.SubElement(xml, "doi_data")
//     etree.SubElement(doi_data, "doi").text = doi_from_url(metadata.id)
//     etree.SubElement(doi_data, "resource").text = metadata.url
//     collection = etree.SubElement(doi_data, "collection", {"property": "text-mining"})
//     item = etree.SubElement(collection, "item")
//     etree.SubElement(item, "resource", {"mime_type": "text/html"}).text = metadata.url
//     if metadata.files is None:
//         return xml
//     for file in metadata.files:
//         # Crossref schema currently doesn't support text/markdown
//         if file.get("mimeType", None) == "text/markdown":
//             file["mimeType"] = "text/plain"
//         item = etree.SubElement(collection, "item")
//         etree.SubElement(
//             item, "resource", {"mime_type": file.get("mimeType", "")}
//         ).text = file.get("url", None)
//     return xml

// def insert_crossref_license(metadata, xml):
//     """Insert crossref license"""
//     if metadata.license is None:
//         return xml
//     license_ = etree.SubElement(xml, "license")
//     if isinstance(metadata.license, dict):
//         r = metadata.license
//     else:
//         r = {}
//         r["rights"] = metadata.license
//         r["rightsUri"] = normalize_id(metadata.license)
//     attributes = compact(
//         {
//             "rightsURI": r.get("rightsUri", None),
//             "rightsIdentifier": r.get("rightsIdentifier", None),
//             "rightsIdentifierScheme": r.get("rightsIdentifierScheme"),
//             "schemeURI": r.get("schemeUri", None),
//             "xml:lang": r.get("lang", None),
//         }
//     )
//     etree.SubElement(license_, "rights", attributes).text = r.get("rights", None)
//     return xml

// def insert_crossref_issn(metadata, xml):
//     """Insert crossref issn"""
//     if (
//         metadata.container is None
//         or metadata.container.get("identifierType", None) != "ISSN"
//     ):
//         return xml
//     etree.SubElement(xml, "issn").text = metadata.container.get("identifier", None)
//     return xml

// def insert_crossref_abstract(metadata, xml):
//     """Insert crossref abstrac"""
//     if metadata.descriptions is None:
//         return xml
//     if isinstance(metadata.descriptions[0], dict):
//         d = metadata.descriptions[0]
//     else:
//         d = {}
//         d["description"] = metadata.descriptions[0]
//     abstract = etree.SubElement(
//         xml, "abstract", {"xmlns": "http://www.ncbi.nlm.nih.gov/JATS1"}
//     )
//     etree.SubElement(abstract, "p").text = d.get("description", None)
//     return xml

// def crossref_root():
//     """Crossref root with namespaces"""
//     doi_batch = """<doi_batch xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns="http://www.crossref.org/schema/5.3.1" xmlns:jats="http://www.ncbi.nlm.nih.gov/JATS1" xmlns:fr="http://www.crossref.org/fundref.xsd" xmlns:mml="http://www.w3.org/1998/Math/MathML" xsi:schemaLocation="http://www.crossref.org/schema/5.3.1 https://www.crossref.org/schemas/crossref5.3.1.xsd" version="5.3.1"></doi_batch>"""
//     return etree.fromstring(doi_batch)

// def generate_crossref_xml_list(metalist) -> Optional[str]:
//     """Generate Crossref XML list."""
//     if not metalist.is_valid:
//         return None
//     xml = crossref_root()
//     head = etree.SubElement(xml, "head")
//     # we use a uuid as batch_id
//     etree.SubElement(head, "doi_batch_id").text = str(uuid.uuid4())
//     etree.SubElement(head, "timestamp").text = datetime.now().strftime("%Y%m%d%H%M%S")
//     depositor = etree.SubElement(head, "depositor")
//     etree.SubElement(depositor, "depositor_name").text = metalist.depositor or "test"
//     etree.SubElement(depositor, "email_address").text = (
//         metalist.email or "info@example.org"
//     )
//     etree.SubElement(head, "registrant").text = metalist.registrant or "test"

//     body = etree.SubElement(xml, "body")
//     body = [insert_crossref_work(item, body) for item in metalist.items]
//     return etree.tostring(
//         xml,
//         doctype='<?xml version="1.0" encoding="UTF-8"?>',
//         pretty_print=True,
//     )

// Load loads the metadata for a single work from a XML file
func Load(filename string) (commonmeta.Data, error) {
	var data commonmeta.Data
	var content Content

	extension := path.Ext(filename)
	if extension != ".xml" {
		return data, errors.New("invalid file extension")
	}
	file, err := os.Open(filename)
	if err != nil {
		return data, errors.New("error reading file")
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&content)
	if err != nil {
		return data, err
	}
	data, err = Read(content)
	if err != nil {
		return data, err
	}
	return data, nil
}
