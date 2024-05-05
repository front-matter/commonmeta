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
	"strings"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/dateutils"
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
	Xmlns        string   `xml:"xmlns,attr"`
	Jats         string   `xml:"jats,attr"`
	AbstractType string   `xml:"abstract-type,attr"`
	Text         string   `xml:",chardata"`
	P            []string `xml:"p"`
	JATSP        []string `xml:"jats:p"`
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

type Archive struct {
	XMLName xml.Name `xml:"archive"`
	Text    string   `xml:",chardata"`
	Name    string   `xml:"name,attr"`
}

type ArchiveLocations struct {
	XMLName xml.Name  `xml:"archive_locations"`
	Text    string    `xml:",chardata"`
	Archive []Archive `xml:"archive"`
}

type Assertion struct {
	XMLName   xml.Name    `xml:"assertion"`
	Text      string      `xml:",chardata"`
	Name      string      `xml:"name,attr"`
	Provider  string      `xml:"provider,attr"`
	Assertion []Assertion `xml:"assertion"`
}

type Book struct {
	XMLName      xml.Name     `xml:"book"`
	BookType     string       `xml:"book_type,attr"`
	BookMetadata BookMetadata `xml:"book_metadata"`
	ContentItem  ContentItem  `xml:"content_item"`
}

type BookMetadata struct {
	XMLName         xml.Name          `xml:"book_metadata"`
	Language        string            `xml:"language,attr"`
	Contributors    Contributors      `xml:"contributors"`
	Titles          Titles            `xml:"titles"`
	PublicationDate []PublicationDate `xml:"publication_date"`
	ISBN            []ISBN            `xml:"isbn"`
	Publisher       Publisher         `xml:"publisher"`
	DOIData         DOIData           `xml:"doi_data"`
}

type ContentItem struct {
	XMLName             xml.Name          `xml:"content_item"`
	ComponentType       string            `xml:"component_type,attr"`
	LevelSequenceNumber string            `xml:"level_sequence_number,attr"`
	PublicationType     string            `xml:"publication_type,attr"`
	Contributors        Contributors      `xml:"contributors"`
	Titles              Titles            `xml:"titles"`
	PublicationDate     []PublicationDate `xml:"publication_date"`
	Pages               struct {
		FirstPage string `xml:"first_page"`
		LastPage  string `xml:"last_page"`
	} `xml:"pages"`
	DOIData      DOIData      `xml:"doi_data"`
	CitationList CitationList `xml:"citation_list"`
}

type Citation struct {
	XMLName      xml.Name `xml:"citation"`
	Key          string   `xml:"key,attr"`
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
	UnstructedCitation string `xml:"unstructured_citation"`
}

type CitationList struct {
	Citation []Citation `xml:"citation"`
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

type Contributors struct {
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
		Program []Program `xml:"program"`
	} `xml:"custom_metadata"`
}

type Dissertation struct {
	XMLName      xml.Name     `xml:"dissertation"`
	PersonName   []PersonName `xml:"person_name"`
	Titles       Titles       `xml:"titles"`
	Institution  Institution  `xml:"institution"`
	Degree       string       `xml:"degree"`
	DOIData      DOIData      `xml:"doi_data"`
	CitationList CitationList `xml:"citation_list"`
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

// ISBN represents a ISSN in Crossref XML metadata.
type ISBN struct {
	XMLName   xml.Name `xml:"isbn"`
	MediaType string   `xml:"media_type,attr"`
	Text      string   `xml:",chardata"`
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
	XMLName         xml.Name `xml:"journal_article"`
	Text            string   `xml:",chardata"`
	PublicationType string   `xml:"publication_type,attr"`
	Pages           struct {
		FirstPage string `xml:"first_page"`
		LastPage  string `xml:"last_page"`
	} `xml:"pages"`
	ReferenceDistributionOpts string            `xml:"reference_distribution_opts,attr"`
	Titles                    Titles            `xml:"titles"`
	Contributors              Contributors      `xml:"contributors"`
	PublicationDate           []PublicationDate `xml:"publication_date"`
	Abstract                  []Abstract        `xml:"jats:abstract"`
	ISSN                      string            `xml:"issn"`
	ItemNumber                ItemNumber        `xml:"item_number"`
	Program                   []Program         `xml:"program"`
	Crossmark                 Crossmark         `xml:"crossmark"`
	ArchiveLocations          ArchiveLocations  `xml:"archive_locations"`
	DOIData                   DOIData           `xml:"doi_data"`
	CitationList              CitationList      `xml:"citation_list"`
}

type JournalIssue struct {
	XMLName         xml.Name          `xml:"journal_issue"`
	PublicationDate []PublicationDate `xml:"publication_date"`
	JournalVolume   JournalVolume     `xml:"journal_volume"`
	Issue           string            `xml:"issue"`
}

// JournalMetadata represents journal metadata in Crossref XML metadata.
type JournalMetadata struct {
	XMLName   xml.Name `xml:"journal_metadata"`
	Language  string   `xml:"language,attr"`
	FullTitle string   `xml:"full_title"`
	ISSN      ISSN     `xml:"issn"`
}

type JournalVolume struct {
	XMLName xml.Name `xml:"journal_volume"`
	Volume  string   `xml:"volume"`
}

type LicenseRef struct {
	Text      string `xml:",chardata"`
	AppliesTo string `xml:"applies_to,attr"`
}

type PeerReview struct {
	XMLName      xml.Name     `xml:"peer_review"`
	Stage        string       `xml:"stage,attr"`
	Type         string       `xml:"type,attr"`
	Contributors Contributors `xml:"contributors"`
	Titles       Titles       `xml:"titles"`
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
	Contributors   Contributors   `xml:"contributors"`
	Titles         Titles         `xml:"titles"`
	PostedDate     PostedDate     `xml:"posted_date"`
	AcceptanceDate AcceptanceDate `xml:"acceptance_date"`
	Institution    Institution    `xml:"institution"`
	Abstract       []Abstract     `xml:"abstract"`
	ItemNumber     ItemNumber     `xml:"item_number"`
	Program        []Program      `xml:"program"`
	DOIData        DOIData        `xml:"doi_data"`
	CitationList   struct {
		Citation []Citation `xml:"citation"`
	} `xml:"citation_list"`
}

type PostedDate struct {
	XMLName   xml.Name `xml:"posted_date"`
	MediaType string   `xml:"media_type,attr"`
	Month     string   `xml:"month"`
	Day       string   `xml:"day"`
	Year      string   `xml:"year"`
}

type Program struct {
	XMLName     xml.Name      `xml:"program"`
	Text        string        `xml:",chardata"`
	Fr          string        `xml:"fr,attr"`
	Name        string        `xml:"name,attr"`
	Ai          string        `xml:"ai,attr"`
	Rel         string        `xml:"rel,attr"`
	Assertion   []Assertion   `xml:"assertion"`
	LicenseRef  []LicenseRef  `xml:"license_ref"`
	RelatedItem []RelatedItem `xml:"related_item"`
}

type PublicationDate struct {
	XMLName   xml.Name `xml:"publication_date"`
	MediaType string   `xml:"media_type,attr"`
	Month     string   `xml:"month"`
	Day       string   `xml:"day"`
	Year      string   `xml:"year"`
}

type Publisher struct {
	XMLName        xml.Name `xml:"publisher"`
	PublisherName  string   `xml:"publisher_name"`
	PublisherPlace string   `xml:"publisher_place"`
}

type PublisherItem struct {
	XMLName    xml.Name `xml:"publisher_item"`
	Text       string   `xml:",chardata"`
	ItemNumber struct {
		Text           string `xml:",chardata"`
		ItemNumberType string `xml:"item_number_type,attr"`
	} `xml:"item_number"`
	Identifier struct {
		Text   string `xml:",chardata"`
		IDType string `xml:"id_type,attr"`
	} `xml:"identifier"`
}

type RelatedItem struct {
	Text              string `xml:",chardata"`
	Description       string `xml:"description"`
	InterWorkRelation struct {
		Text             string `xml:",chardata"`
		IdentifierType   string `xml:"identifier-type,attr"`
		RelationshipType string `xml:"relationship-type,attr"`
	} `xml:"inter_work_relation"`
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

// Titles represents the titles in Crossref XML metadata.
type Titles struct {
	XMLName               xml.Name `xml:"titles"`
	Title                 string   `xml:"title"`
	Subtitle              string   `xml:"subtitle"`
	OriginalLanguageTitle struct {
		Text     string `xml:",chardata"`
		Language string `xml:"language,attr"`
	} `xml:"original_language_title"`
}

// CRToCMMappings maps Crossref XML types to Commonmeta types
// source: http://api.crossref.org/types,
// Crossref XML naming conventions are different from the REST API
var CRToCMMappings = map[string]string{
	"book_chapter":        "BookChapter",
	"book_content":        "BookChapter",
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

	var containerTitle, firstPage, issn, issue, language, lastPage, volume string
	var accessIndicators Program
	var abstract []Abstract
	var archiveLocations ArchiveLocations
	var citationList CitationList
	var contributors Contributors
	var doiData DOIData
	var fundref Program
	var isbn []ISBN
	var publicationDate []PublicationDate
	var program []Program
	var subjects []string
	var titles Titles

	meta := content.Query.DOIRecord.Crossref

	data.ID = doiutils.NormalizeDOI(content.Query.DOI.Text)
	data.Type = CRToCMMappings[content.Query.DOI.Type]
	if data.Type == "" {
		data.Type = "Other"
	}

	// fetch metadata depending on Crossref type (using the commonmeta vocabulary)
	switch data.Type {
	case "JournalArticle":
		abstract = meta.Journal.JournalArticle.Abstract
		citationList = meta.Journal.JournalArticle.CitationList
		contributors = meta.Journal.JournalArticle.Contributors
		archiveLocations = meta.Journal.JournalArticle.ArchiveLocations
		containerTitle = meta.Journal.JournalMetadata.FullTitle
		doiData = meta.Journal.JournalArticle.DOIData
		firstPage = meta.Journal.JournalArticle.Pages.FirstPage
		issn = meta.Journal.JournalMetadata.ISSN.Text
		issue = meta.Journal.JournalIssue.Issue
		language = meta.Journal.JournalMetadata.Language
		lastPage = meta.Journal.JournalArticle.Pages.LastPage
		program = append(program, meta.Journal.JournalArticle.Program...)
		// program metadata is also found in crossmark custom metadata
		program = append(program, meta.Journal.JournalArticle.Crossmark.CustomMetadata.Program...)
		publicationDate = meta.Journal.JournalArticle.PublicationDate
		titles = meta.Journal.JournalArticle.Titles
		volume = meta.Journal.JournalIssue.JournalVolume.Volume
	case "JournalIssue":
		language = meta.Journal.JournalMetadata.Language
	case "Journal":
		language = meta.Journal.JournalMetadata.Language
	case "Article":
		abstract = meta.PostedContent.Abstract
		// archiveLocations not supported
		citationList = meta.PostedContent.CitationList
		contributors = meta.PostedContent.Contributors
		doiData = meta.PostedContent.DOIData
		language = meta.PostedContent.Language
		program = meta.PostedContent.Program
		// use posted date as publication date
		publicationDate = append(publicationDate, PublicationDate{
			Year:  meta.PostedContent.PostedDate.Year,
			Month: meta.PostedContent.PostedDate.Month,
			Day:   meta.PostedContent.PostedDate.Day,
		})
		subjects = append(subjects, meta.PostedContent.GroupTitle)
		titles = meta.PostedContent.Titles
	case "BookChapter":
		citationList = meta.Book.ContentItem.CitationList
		contributors = meta.Book.ContentItem.Contributors
		doiData = meta.Book.BookMetadata.DOIData
		isbn = meta.Book.BookMetadata.ISBN
		language = meta.Book.BookMetadata.Language
		publicationDate = meta.Book.ContentItem.PublicationDate
		titles = meta.Book.ContentItem.Titles
	case "BookSeries":
	case "BookSet":
	case "Book":
		contributors = meta.Book.BookMetadata.Contributors
		citationList = meta.Book.ContentItem.CitationList
		doiData = meta.Book.BookMetadata.DOIData
		language = meta.Book.BookMetadata.Language
		titles = meta.Book.BookMetadata.Titles
	case "ProceedingsArticle":
	case "Component":
	case "Dataset":
	case "Report":
	case "PeerReview":
		doiData = meta.PeerReview.DOIData
	case "Dissertation":
		doiData = meta.Dissertation.DOIData
	}

	if len(archiveLocations.Archive) > 0 {
		var al []string
		for _, v := range archiveLocations.Archive {
			al = append(al, v.Name)
		}
		data.ArchiveLocations = al
	}

	containerType := commonmeta.ContainerTypes[data.Type]
	var identifier, identifierType string
	if issn != "" {
		identifier = issn
		identifierType = "ISSN"
	} else if len(isbn) > 0 {
		// find the first electronic ISBN, use the first ISBN if no electronic ISBN is found
		i := slices.IndexFunc(isbn, func(c ISBN) bool { return c.MediaType == "electronic" })
		if i == -1 {
			i = 0
		}
		identifier = isbn[i].Text
		identifierType = "ISBN"
	}

	data.Container = commonmeta.Container{
		Identifier:     identifier,
		IdentifierType: identifierType,
		Type:           containerType,
		Title:          containerTitle,
		Volume:         volume,
		Issue:          issue,
		FirstPage:      firstPage,
		LastPage:       lastPage,
	}

	if len(contributors.PersonName) > 0 {
		contrib, err := GetContributors(contributors)
		if err != nil {
			return data, err
		}
		data.Contributors = append(data.Contributors, contrib...)
	}

	if len(publicationDate) > 0 {
		i := slices.IndexFunc(publicationDate, func(c PublicationDate) bool { return c.MediaType == "online" })
		if i == -1 {
			i = 0
		}
		data.Date.Published = dateutils.GetDateFromCrossrefParts(publicationDate[i].Year, publicationDate[i].Month, publicationDate[i].Day)
	}

	if len(abstract) > 0 {
		for _, v := range abstract {
			d := strings.Join(v.P, " ")
			t := v.AbstractType
			if t == "" {
				t = "Abstract"
			}
			data.Descriptions = append(data.Descriptions, commonmeta.Description{
				Description: utils.Sanitize(d),
				Type:        t,
			})
		}
	}

	if len(doiData.Collection.Item) > 0 {
		for _, v := range doiData.Collection.Item {
			if v.Resource.Text != "" && v.Resource.MimeType != "" {
				file := commonmeta.File{
					URL:      v.Resource.Text,
					MimeType: v.Resource.MimeType,
				}
				data.Files = append(data.Files, file)
			}
		}
		data.Files = utils.DedupeSlice(data.Files)
	}

	if len(program) > 0 {
		i := slices.IndexFunc(program, func(c Program) bool { return c.Name == "AccessIndicators" })
		if i != -1 {
			accessIndicators = program[i]
		}
		j := slices.IndexFunc(program, func(c Program) bool { return c.Name == "fundref" })
		if j != -1 {
			fundref = program[j]
		}
	}

	if len(fundref.Assertion) > 0 {
		fundingReferences, err := GetFundingReferences(fundref)
		if err != nil {
			return data, err
		}
		data.FundingReferences = append(data.FundingReferences, fundingReferences...)
	}

	data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
		Identifier:     data.ID,
		IdentifierType: "DOI",
	})

	data.Language = language

	if len(accessIndicators.LicenseRef) > 0 {
		// find the first license that applies to the version of record, use the first license if no license is found
		i := slices.IndexFunc(accessIndicators.LicenseRef, func(c LicenseRef) bool { return c.AppliesTo == "vor" })
		if i == -1 {
			i = 0
		}
		url, _ := utils.NormalizeCCUrl(accessIndicators.LicenseRef[i].Text)
		id := utils.URLToSPDX(url)
		data.License = commonmeta.License{
			ID:  id,
			URL: url,
		}
	}

	if len(subjects) > 0 {
		for _, v := range subjects {
			data.Subjects = append(data.Subjects, commonmeta.Subject{
				Subject: v,
			})
		}
	}

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

	if len(citationList.Citation) > 0 {
		for _, v := range citationList.Citation {
			reference := commonmeta.Reference{
				Key:             v.Key,
				ID:              doiutils.NormalizeDOI(v.Doi.Text),
				Title:           v.ArticleTitle,
				PublicationYear: v.CYear,
				Unstructured:    v.UnstructedCitation,
			}
			containsKey := slices.ContainsFunc(data.References, func(e commonmeta.Reference) bool {
				return e.Key != "" && e.Key == reference.Key
			})
			if !containsKey {
				data.References = append(data.References, reference)
			}
		}
	}

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

	if titles.Title != "" {
		data.Titles = append(data.Titles, commonmeta.Title{
			Title: titles.Title,
		})
	}
	if titles.Subtitle != "" {
		data.Titles = append(data.Titles, commonmeta.Title{
			Title: titles.Subtitle,
			Type:  "Subtitle",
		})
	}
	if titles.OriginalLanguageTitle.Text != "" {
		data.Titles = append(data.Titles, commonmeta.Title{
			Title:    titles.OriginalLanguageTitle.Text,
			Type:     "TranslatedTitle",
			Language: titles.OriginalLanguageTitle.Language,
		})
	}

	data.URL = doiData.Resource

	return data, nil
}

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

func GetContributors(contrib Contributors) ([]commonmeta.Contributor, error) {
	var contributors []commonmeta.Contributor

	if len(contrib.PersonName) > 0 {
		for _, v := range contrib.PersonName {
			var ID string
			if v.GivenName != "" || v.Surname != "" {
				if v.ORCID != "" {
					ID, _ = utils.NormalizeURL(v.ORCID, true, false) // enforce HTTPS
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
			containsName := slices.ContainsFunc(contributors, func(e commonmeta.Contributor) bool {
				return e.GivenName == contributor.GivenName && e.FamilyName != "" && e.FamilyName == contributor.FamilyName
			})
			if !containsName {
				contributors = append(contributors, contributor)
			}
		}
	}
	return contributors, nil
}

func GetFundingReferences(fundref Program) ([]commonmeta.FundingReference, error) {
	var fundingReferences []commonmeta.FundingReference

	var fundGroups []Assertion
	for _, v := range fundref.Assertion {
		if v.Name == "fundgroup" {
			fundGroups = append(fundGroups, v)
		}
	}

	var funderName, funderIdentifier, funderIdentifierType string
	for _, fundgroup := range fundGroups {
		var awardNumbers []Assertion
		for _, awardNumber := range fundgroup.Assertion {
			if awardNumber.Name == "award_number" {
				awardNumbers = append(awardNumbers, awardNumber)
			}
		}
		for _, v := range fundgroup.Assertion {
			if v.Name == "funder_name" {
				if v.Assertion != nil {
					for _, a := range v.Assertion {
						if a.Name == "funder_identifier" {
							if a.Provider == "crossref" {
								funderIdentifierType = "Crossref Funder ID"
								funderIdentifier = doiutils.NormalizeDOI("10.13039/" + a.Text)
							} else {
								funderIdentifier = doiutils.NormalizeDOI(a.Text)
								if funderIdentifier == "" {
									funderIdentifier = a.Text
								}
							}
						}
					}
				}
				funderName = strings.TrimSpace(v.Text)
			}
			if len(awardNumbers) > 0 {
				for _, awardNumber := range awardNumbers {
					fundingReference := commonmeta.FundingReference{
						FunderIdentifier:     funderIdentifier,
						FunderIdentifierType: funderIdentifierType,
						FunderName:           funderName,
						AwardNumber:          awardNumber.Text,
					}
					fundingReferences = append(fundingReferences, fundingReference)
				}
			} else {
				// no award numbers
				fundingReference := commonmeta.FundingReference{
					FunderIdentifier:     funderIdentifier,
					FunderIdentifierType: funderIdentifierType,
					FunderName:           funderName,
				}
				fundingReferences = append(fundingReferences, fundingReference)
			}
		}
	}
	return fundingReferences, nil
}
