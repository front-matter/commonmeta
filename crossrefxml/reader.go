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
	Query   Query    `xml:"query"`
}

// Content represents the Crossref XML metadata returned from Crossref. The type is more
// flexible than the Crossrefxml type, allowing for different formats of some metadata.
type Content struct {
	Crossrefxml
	Query            Query `xml:"query"`
	CrossrefMetadata Query `xml:"crossref_metadata"`
}

type Query struct {
	Status    string    `xml:"status,attr"`
	DOI       DOI       `xml:"doi"`
	CRMItem   []CRMItem `xml:"crm-item"`
	DOIRecord DOIRecord `xml:"doi_record"`
}

type CRMItem struct {
	XMLName xml.Name `xml:"crm-item"`
	Text    string   `xml:",chardata"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Claim   string   `xml:"claim,attr"`
}

type DOI struct {
	Type string `xml:"type,attr"`
	Text string `xml:",chardata"`
}

type DOIRecord struct {
	XMLName  xml.Name `xml:"doi_record"`
	Crossref struct {
		Xmlns          string        `xml:"xmlns,attr"`
		SchemaLocation string        `xml:"schemaLocation,attr"`
		Book           Book          `xml:"book,omitempty"`
		Conference     Conference    `xml:"conference,omitempty"`
		Database       Database      `xml:"database,omitempty"`
		Dissertation   Dissertation  `xml:"dissertation,omitempty"`
		Journal        Journal       `xml:"journal,omitempty"`
		PeerReview     PeerReview    `xml:"peer_review,omitempty"`
		PostedContent  PostedContent `xml:"posted_content,omitempty"`
		SAComponent    SAComponent   `xml:"sa_component,omitempty"`
		Standard       Standard      `xml:"standard,omitempty"`
	} `xml:"crossref"`
}

type Abstract struct {
	XMLName      xml.Name `xml:"abstract"`
	Xmlns        string   `xml:"xmlns,attr"`
	Title        string   `xml:"title"`
	AbstractType string   `xml:"abstract-type,attr"`
	Text         string   `xml:",chardata"`
	P            []P      `xml:"p"`
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

type ApprovalDate struct {
	XMLName xml.Name `xml:"approval_date"`
	Month   string   `xml:"month"`
	Day     string   `xml:"day"`
	Year    string   `xml:"year"`
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
	XMLName    xml.Name    `xml:"assertion"`
	Text       string      `xml:",chardata"`
	Name       string      `xml:"name,attr"`
	Provider   string      `xml:"provider,attr"`
	Label      string      `xml:"label,attr"`
	GroupName  string      `xml:"group_name,attr"`
	GroupLabel string      `xml:"group_label,attr"`
	Order      string      `xml:"order,attr"`
	Assertion  []Assertion `xml:"assertion"`
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
	Abstract        []Abstract        `xml:"abstract"`
	EditionNumber   int               `xml:"edition_number"`
	PublicationDate []PublicationDate `xml:"publication_date"`
	ISBN            []ISBN            `xml:"isbn"`
	Publisher       Publisher         `xml:"publisher"`
	DOIData         DOIData           `xml:"doi_data"`
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

type Component struct {
	XMLName        xml.Name `xml:"component"`
	RegAgency      string   `xml:"reg-agency,attr"`
	ParentRelation string   `xml:"parent_relation,attr"`
	Text           string   `xml:",chardata"`
	Titles         Titles   `xml:"titles"`
	Format         Format   `xml:"format"`
	DOIData        DOIData  `xml:"doi_data"`
}

type ComponentList struct {
	XMLName   xml.Name    `xml:"component_list"`
	Text      string      `xml:",chardata"`
	Component []Component `xml:"component"`
}

type Conference struct {
	XMLName             xml.Name            `xml:"conference"`
	Contributors        Contributors        `xml:"contributors"`
	EventMetadata       EventMetadata       `xml:"event_metadata"`
	ProceedingsMetadata ProceedingsMetadata `xml:"proceedings_metadata"`
	ConferencePaper     ConferencePaper     `xml:"conference_paper"`
}

type ConferencePaper struct {
	XMLName         xml.Name          `xml:"conference_paper"`
	PublicationType string            `xml:"publication_type,attr"`
	Contributors    Contributors      `xml:"contributors"`
	Titles          Titles            `xml:"titles"`
	PublicationDate []PublicationDate `xml:"publication_date"`
	Pages           Pages             `xml:"pages"`
	PublisherItem   PublisherItem     `xml:"publisher_item"`
	Crossmark       Crossmark         `xml:"crossmark"`
	DOIData         DOIData           `xml:"doi_data"`
	CitationList    CitationList      `xml:"citation_list"`
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

type CreationDate struct {
	XMLName   xml.Name `xml:"creation_date"`
	MediaType string   `xml:"media_type,attr"`
	Month     string   `xml:"month"`
	Day       string   `xml:"day"`
	Year      string   `xml:"year"`
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
	CrossmarkDomainExclusive string         `xml:"crossmark_domain_exclusive"`
	CustomMetadata           CustomMetadata `xml:"custom_metadata"`
}

type CustomMetadata struct {
	XMLName   xml.Name    `xml:"custom_metadata"`
	Text      string      `xml:",chardata"`
	Assertion []Assertion `xml:"assertion"`
	Program   []Program   `xml:"program"`
}

type Database struct {
	DatabaseMetadata DatabaseMetadata `xml:"database_metadata"`
	Dataset          Dataset          `xml:"dataset"`
}

type DatabaseMetadata struct {
	Titles struct {
		Title string `xml:"title"`
	} `xml:"titles"`
}

type Dataset struct {
	Contributors Contributors `xml:"contributors"`
	Titles       Titles       `xml:"titles"`
	DatabaseDate struct {
		CreationDate CreationDate `xml:"creation_date"`
	} `xml:"database_date"`
	DOIData DOIData `xml:"doi_data"`
}

type Dissertation struct {
	XMLName         xml.Name     `xml:"dissertation"`
	Language        string       `xml:"language"`
	PublicationType string       `xml:"publication_type"`
	PersonName      []PersonName `xml:"person_name"`
	Titles          Titles       `xml:"titles"`
	ApprovalDate    ApprovalDate `xml:"approval_date"`
	Institution     Institution  `xml:"institution"`
	Degree          string       `xml:"degree"`
	DOIData         DOIData      `xml:"doi_data"`
	CitationList    CitationList `xml:"citation_list"`
}

type DOIData struct {
	XMLName    xml.Name   `xml:"doi_data"`
	DOI        string     `xml:"doi"`
	Timestamp  string     `xml:"timestamp"`
	Resource   string     `xml:"resource"`
	Collection Collection `xml:"collection"`
}

type EventMetadata struct {
	XMLName            xml.Name `xml:"event_metadata"`
	ConferenceName     string   `xml:"conference_name"`
	ConferenceAcronym  string   `xml:"conference_acronym"`
	ConferenceSponsor  string   `xml:"conference_sponsor"`
	ConferenceLocation string   `xml:"conference_location"`
	ConferenceDate     string   `xml:"conference_date"`
}

type Format struct {
	XMLName  xml.Name `xml:"format"`
	MimeType string   `xml:"mime_type,attr"`
	Text     string   `xml:",chardata"`
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
	XMLName        xml.Name `xml:"item_number"`
	ItemNumberType string   `xml:"item_number_type,attr"`
	Text           string   `xml:",chardata"`
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
	XMLName                   xml.Name          `xml:"journal_article"`
	Text                      string            `xml:",chardata"`
	PublicationType           string            `xml:"publication_type,attr"`
	ReferenceDistributionOpts string            `xml:"reference_distribution_opts,attr"`
	Titles                    Titles            `xml:"titles"`
	Contributors              Contributors      `xml:"contributors"`
	PublicationDate           []PublicationDate `xml:"publication_date"`
	PublisherItem             struct {
		ItemNumber ItemNumber `xml:"item_number"`
	} `xml:"publisher_item"`
	Abstract         []Abstract       `xml:"jats:abstract"`
	Pages            Pages            `xml:"pages"`
	ISSN             []ISSN           `xml:"issn"`
	Program          []Program        `xml:"program"`
	Crossmark        Crossmark        `xml:"crossmark"`
	ArchiveLocations ArchiveLocations `xml:"archive_locations"`
	DOIData          DOIData          `xml:"doi_data"`
	CitationList     CitationList     `xml:"citation_list"`
}

type JournalIssue struct {
	XMLName         xml.Name          `xml:"journal_issue"`
	PublicationDate []PublicationDate `xml:"publication_date"`
	JournalVolume   JournalVolume     `xml:"journal_volume"`
	Issue           string            `xml:"issue"`
	DOIData         DOIData           `xml:"doi_data"`
}

// JournalMetadata represents journal metadata in Crossref XML metadata.
type JournalMetadata struct {
	XMLName   xml.Name `xml:"journal_metadata"`
	Language  string   `xml:"language,attr"`
	FullTitle string   `xml:"full_title"`
	ISSN      []ISSN   `xml:"issn"`
	DOIData   DOIData  `xml:"doi_data"`
}

type JournalVolume struct {
	XMLName xml.Name `xml:"journal_volume"`
	Volume  string   `xml:"volume"`
}

type LicenseRef struct {
	Text      string `xml:",chardata"`
	AppliesTo string `xml:"applies_to,attr"`
}

type P struct {
	XMLName xml.Name `xml:"p"`
	Xmlns   string   `xml:"xmlns,attr"`
	Text    string   `xml:",chardata"`
}

type Pages struct {
	FirstPage string `xml:"first_page"`
	LastPage  string `xml:"last_page"`
}

type PeerReview struct {
	XMLName                    xml.Name     `xml:"peer_review"`
	Stage                      string       `xml:"stage,attr"`
	RevisionRound              string       `xml:"revision-round,attr"`
	Recommendation             string       `xml:"recommendation,attr"`
	Type                       string       `xml:"type,attr"`
	Contributors               Contributors `xml:"contributors"`
	Titles                     Titles       `xml:"titles"`
	ReviewDate                 ReviewDate   `xml:"review_date"`
	CompetingInterestStatement string       `xml:"competing_interest_statement"`
	Program                    []Program    `xml:"program"`
	DOIData                    DOIData      `xml:"doi_data"`
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
	Affiliation     string        `xml:"affiliation"`
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

type ProceedingsMetadata struct {
	XMLName          xml.Name          `xml:"proceedings_metadata"`
	Language         string            `xml:"language,attr"`
	ProceedingsTitle string            `xml:"proceedings_title"`
	Publisher        Publisher         `xml:"publisher"`
	PublicationDate  []PublicationDate `xml:"publication_date"`
	ISBN             []ISBN            `xml:"isbn"`
	PublisherItem    PublisherItem     `xml:"publisher_item"`
	DOIData          DOIData           `xml:"doi_data"`
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
	XMLName    xml.Name   `xml:"publisher_item"`
	Text       string     `xml:",chardata"`
	ItemNumber ItemNumber `xml:"item_number"`
	Identifier struct {
		Text   string `xml:",chardata"`
		IDType string `xml:"id_type,attr"`
	} `xml:"identifier"`
}

type RelatedItem struct {
	XMLName           xml.Name `xml:"related_item"`
	Text              string   `xml:",chardata"`
	Description       string   `xml:"description"`
	InterWorkRelation struct {
		Text             string `xml:",chardata"`
		IdentifierType   string `xml:"identifier-type,attr"`
		RelationshipType string `xml:"relationship-type,attr"`
	} `xml:"inter_work_relation"`
	IntraWorkRelation struct {
		Text             string `xml:",chardata"`
		IdentifierType   string `xml:"identifier-type,attr"`
		RelationshipType string `xml:"relationship-type,attr"`
	} `xml:"intra_work_relation"`
}

// Resource represents a resource in Crossref XML metadata.
type Resource struct {
	XMLName  xml.Name `xml:"resource"`
	Text     string   `xml:",chardata"`
	MimeType string   `xml:"mime_type,attr"`
}

type ReviewDate struct {
	XMLDate xml.Name `xml:"review_date"`
	Month   string   `xml:"month"`
	Day     string   `xml:"day"`
	Year    string   `xml:"year"`
}

type SAComponent struct {
	XMLName       xml.Name      `xml:"sa_component"`
	ComponentList ComponentList `xml:"component_list"`
}

type Standard struct {
	XMLName xml.Name `xml:"standard"`
	DOIData DOIData  `xml:"doi_data"`
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

// CRToCMMappings maps Crossref Query types to Commonmeta types
// source: https://www.crossref.org/schemas/crossref_query_output3.0.xsd
// Crossref Query naming conventions are different from the REST API
var CRToCMMappings = map[string]string{
	"journal_title":       "Journal",
	"journal_issue":       "JournalIssue",
	"journal_volume":      "JournalVolume",
	"journal_article":     "JournalArticle",
	"conference_title":    "Proceedings",
	"conference_series":   "ProceedingsSeries",
	"conference_paper":    "ProceedingsArticle",
	"book_title":          "Book",
	"book_series":         "BookSeries",
	"book_content":        "BookChapter",
	"component":           "Component",
	"dissertation":        "Dissertation",
	"peer_review":         "PeerReview",
	"posted_content":      "Article",
	"report-paper_title":  "Report",
	"report-paper_series": "ReportSeries",
	//"report-paper_content": "ReportComponent",
	"standard_title":  "Standard",
	"standard_series": "StandardSeries",
	// "standard_content": "StandardComponent",
	// "prepublication":
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
	return crossrefResult.QueryResult.Body, err
}

// Read Crossref XML response and return work struct in Commonmeta format
func Read(content Content) (commonmeta.Data, error) {
	var data = commonmeta.Data{}

	var containerTitle, issue, language, volume string
	var accessIndicators Program
	var abstract []Abstract
	var archiveLocations ArchiveLocations
	var citationList CitationList
	var contributors Contributors
	var customMetadata CustomMetadata
	var doiData DOIData
	var fundref Program
	var isbn []ISBN
	var issn []ISSN
	var itemNumber ItemNumber
	var pages Pages
	var publicationDate []PublicationDate
	var program []Program
	var relations Program
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
	case "Article": // posted-content
		abstract = meta.PostedContent.Abstract
		// archiveLocations not supported
		citationList = meta.PostedContent.CitationList
		contributors = meta.PostedContent.Contributors
		doiData = meta.PostedContent.DOIData
		itemNumber = meta.PostedContent.ItemNumber
		language = meta.PostedContent.Language
		// pages not supported
		program = meta.PostedContent.Program
		// use posted date as publication date
		publicationDate = append(publicationDate, PublicationDate{
			Year:  meta.PostedContent.PostedDate.Year,
			Month: meta.PostedContent.PostedDate.Month,
			Day:   meta.PostedContent.PostedDate.Day,
		})
		// use group title for subjects
		subjects = append(subjects, meta.PostedContent.GroupTitle)
		titles = meta.PostedContent.Titles
	case "Book":
		abstract = meta.Book.BookMetadata.Abstract
		contributors = meta.Book.BookMetadata.Contributors
		citationList = meta.Book.ContentItem.CitationList
		doiData = meta.Book.BookMetadata.DOIData
		isbn = meta.Book.BookMetadata.ISBN
		language = meta.Book.BookMetadata.Language
		publicationDate = meta.Book.BookMetadata.PublicationDate
		pages = meta.Book.ContentItem.Pages
		titles = meta.Book.BookMetadata.Titles
	case "BookChapter":
		abstract = meta.Book.BookMetadata.Abstract
		citationList = meta.Book.ContentItem.CitationList
		contributors = meta.Book.ContentItem.Contributors
		doiData = meta.Book.BookMetadata.DOIData
		isbn = meta.Book.BookMetadata.ISBN
		language = meta.Book.BookMetadata.Language
		publicationDate = meta.Book.ContentItem.PublicationDate
		pages = meta.Book.ContentItem.Pages
		titles = meta.Book.ContentItem.Titles
	case "BookPart":
	case "BookSection":
	case "BookSeries":
	case "BookSet":
	case "BookTrack":
	case "Component":
		doiData = meta.SAComponent.ComponentList.Component[0].DOIData
	case "Database":
	case "Dataset":
		containerTitle = meta.Database.DatabaseMetadata.Titles.Title
		contributors = meta.Database.Dataset.Contributors
		titles = meta.Database.Dataset.Titles
		// use creation date as publication date
		publicationDate = append(publicationDate, PublicationDate{
			Year:  meta.Database.Dataset.DatabaseDate.CreationDate.Year,
			Month: meta.Database.Dataset.DatabaseDate.CreationDate.Month,
			Day:   meta.Database.Dataset.DatabaseDate.CreationDate.Day,
		})
		doiData = meta.Database.Dataset.DOIData
	case "Dissertation":
		contributors = Contributors{
			PersonName: meta.Dissertation.PersonName,
		}
		doiData = meta.Dissertation.DOIData
		// use approval date as publication date
		publicationDate = append(publicationDate, PublicationDate{
			Year:  meta.Dissertation.ApprovalDate.Year,
			Month: meta.Dissertation.ApprovalDate.Month,
			Day:   meta.Dissertation.ApprovalDate.Day,
		})
		titles = meta.Dissertation.Titles
	case "EditedBook":
	case "Grant":
	case "Journal":
		containerTitle = meta.Journal.JournalMetadata.FullTitle
		language = meta.Journal.JournalMetadata.Language
		doiData = meta.Journal.JournalMetadata.DOIData
	case "JournalArticle":
		abstract = meta.Journal.JournalArticle.Abstract
		archiveLocations = meta.Journal.JournalArticle.ArchiveLocations
		citationList = meta.Journal.JournalArticle.CitationList
		containerTitle = meta.Journal.JournalMetadata.FullTitle
		contributors = meta.Journal.JournalArticle.Contributors
		customMetadata = meta.Journal.JournalArticle.Crossmark.CustomMetadata
		doiData = meta.Journal.JournalArticle.DOIData
		issn = meta.Journal.JournalMetadata.ISSN
		issue = meta.Journal.JournalIssue.Issue
		itemNumber = meta.Journal.JournalArticle.PublisherItem.ItemNumber
		language = meta.Journal.JournalMetadata.Language
		pages = meta.Journal.JournalArticle.Pages
		program = append(program, meta.Journal.JournalArticle.Program...)
		publicationDate = meta.Journal.JournalArticle.PublicationDate
		titles = meta.Journal.JournalArticle.Titles
		volume = meta.Journal.JournalIssue.JournalVolume.Volume
	case "JournalIssue":
		containerTitle = meta.Journal.JournalMetadata.FullTitle
		doiData = meta.Journal.JournalIssue.DOIData
		issn = meta.Journal.JournalMetadata.ISSN
		language = meta.Journal.JournalMetadata.Language
	case "JournalVolume":
	case "Monograph":
	case "Other":
	case "PeerReview":
		contributors = meta.PeerReview.Contributors
		titles = meta.PeerReview.Titles
		program = meta.PeerReview.Program
		// use review date as publication date
		publicationDate = append(publicationDate, PublicationDate{
			Year:  meta.PeerReview.ReviewDate.Year,
			Month: meta.PeerReview.ReviewDate.Month,
			Day:   meta.PeerReview.ReviewDate.Day,
		})
		doiData = meta.PeerReview.DOIData
	case "Proceedings":
	case "ProceedingsArticle":
		citationList = meta.Conference.ConferencePaper.CitationList
		containerTitle = meta.Conference.EventMetadata.ConferenceName
		contributors = meta.Conference.ConferencePaper.Contributors
		doiData = meta.Conference.ConferencePaper.DOIData
		isbn = meta.Conference.ProceedingsMetadata.ISBN
		program = meta.Conference.ConferencePaper.Crossmark.CustomMetadata.Program
		pages = meta.Conference.ConferencePaper.Pages
		publicationDate = meta.Conference.ConferencePaper.PublicationDate
		titles = meta.Conference.ConferencePaper.Titles
	case "ProceedingsSeries":
	case "ReferenceBook":
	case "ReferenceEntry":
	case "Report":
	case "ReportComponent":
	case "ReportSeries":
	case "Standard":
	}

	// program metadata is also found in crossmark custom metadata
	program = append(program, customMetadata.Program...)

	if len(program) > 0 {
		i := slices.IndexFunc(program, func(c Program) bool { return c.Name == "AccessIndicators" })
		if i != -1 {
			accessIndicators = program[i]
		}
		j := slices.IndexFunc(program, func(c Program) bool { return c.Name == "fundref" })
		if j != -1 {
			fundref = program[j]
		}
		k := slices.IndexFunc(program, func(c Program) bool { return c.Name == "" })
		if k != -1 {
			relations = program[k]
		}
	}

	// submission and acceptance dates may be found in crossmark custom metadata
	if len(customMetadata.Assertion) > 0 {
		s := slices.IndexFunc(customMetadata.Assertion, func(c Assertion) bool { return c.Name == "received" })
		if s != -1 {
			dateSubmitted := customMetadata.Assertion[s]
			data.Date.Submitted = dateutils.ParseDate(dateSubmitted.Text)
		}
		a := slices.IndexFunc(customMetadata.Assertion, func(c Assertion) bool { return c.Name == "accepted" })
		if a != -1 {
			dateAccepted := customMetadata.Assertion[a]
			data.Date.Accepted = dateutils.ParseDate(dateAccepted.Text)
		}
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
	if len(issn) > 0 {
		// find the first electronic ISSN, use the first ISSN if no electronic ISSN is found
		i := slices.IndexFunc(issn, func(c ISSN) bool { return c.MediaType == "electronic" })
		if i == -1 {
			i = 0
		}
		identifier = issn[i].Text
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
		FirstPage:      pages.FirstPage,
		LastPage:       pages.LastPage,
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
			var str []string
			for _, p := range v.P {
				str = append(str, p.Text)
			}
			d := strings.TrimSpace(strings.Join(str, " "))
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
	if len(itemNumber.Text) > 0 {
		identifier = itemNumber.Text
		identifierType := "Other"
		// use only known identifier types, otherwise use "Other". Case insensitive.
		if slices.ContainsFunc(commonmeta.IdentifierTypes, func(s string) bool {
			return strings.EqualFold(s, itemNumber.ItemNumberType)
		}) {
			identifierType = strings.ToUpper(itemNumber.ItemNumberType)
		}
		if identifierType == "UUID" {
			// workaround for UUID as Crossref item number can only be 32 characters long
			if len(identifier) == 32 {
				identifier = strings.Join([]string{identifier[:8], identifier[8:12], identifier[12:16], identifier[16:20], identifier[20:]}, "-")
			}
		}

		data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
			Identifier:     identifier,
			IdentifierType: identifierType,
		})
	}

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
			// } else if v.Name == "created" {
			// 	data.Date.Created = v.Text
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

	if len(relations.RelatedItem) > 0 {
		for _, v := range relations.RelatedItem {
			var id, t string
			if v.InterWorkRelation.Text != "" {
				if v.InterWorkRelation.IdentifierType == "doi" {
					id = doiutils.NormalizeDOI(v.InterWorkRelation.Text)
				} else if v.InterWorkRelation.IdentifierType == "issn" {
					id = utils.ISSNAsURL(v.InterWorkRelation.Text)
				} else if utils.ValidateURL(v.InterWorkRelation.Text) == "URL" {
					id = v.InterWorkRelation.Text
				} else {
					id = v.InterWorkRelation.Text
				}
				t = utils.TitleCase(v.InterWorkRelation.RelationshipType)
			} else if v.IntraWorkRelation.Text != "" {
				if v.IntraWorkRelation.IdentifierType == "doi" {
					id = doiutils.NormalizeDOI(v.IntraWorkRelation.Text)
				} else if v.IntraWorkRelation.IdentifierType == "issn" {
					id = utils.ISSNAsURL(v.IntraWorkRelation.Text)
				} else if utils.ValidateURL(v.IntraWorkRelation.Text) == "URL" {
					id = v.IntraWorkRelation.Text
				} else {
					id = v.IntraWorkRelation.Text
				}
				t = utils.TitleCase(v.IntraWorkRelation.RelationshipType)
			}
			relation := commonmeta.Relation{
				ID:   id,
				Type: t,
			}
			data.Relations = append(data.Relations, relation)
		}
	}
	if data.Container.IdentifierType == "ISSN" {
		data.Relations = append(data.Relations, commonmeta.Relation{
			ID:   utils.ISSNAsURL(data.Container.Identifier),
			Type: "IsPartOf",
		})
	}

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

// ReadList reads a list of Crossref XML responses and returns a list of works in Commonmeta format
func ReadList(content []Content) ([]commonmeta.Data, error) {
	var data []commonmeta.Data
	for _, v := range content {
		d, err := Read(v)
		if err != nil {
			log.Println(err)
		}
		data = append(data, d)
	}
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

// LoadList loads the metadata for a list of works from an XML file and converts it to the Commonmeta format
func LoadList(filename string) ([]commonmeta.Data, error) {
	type Response struct {
		ListRecords struct {
			Record []struct {
				Metadata struct {
					CrossrefResult struct {
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
				} `xml:"metadata"`
			} `xml:"record"`
		} `xml:"ListRecords"`
	}

	var data []commonmeta.Data
	var content []Content
	var err error

	extension := path.Ext(filename)
	if extension != ".xml" {
		return data, errors.New("unsupported file format")
	}

	var response Response
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("error reading file")
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&response)
	if err != nil {
		return data, err
	}
	for _, v := range response.ListRecords.Record {
		// rewrite XML structure to match the Crossref API response
		c := Content{
			Query: Query{
				Status:    "resolved",
				DOI:       v.Metadata.CrossrefResult.QueryResult.Body.CrossrefMetadata.DOI,
				DOIRecord: v.Metadata.CrossrefResult.QueryResult.Body.CrossrefMetadata.DOIRecord,
				CRMItem:   v.Metadata.CrossrefResult.QueryResult.Body.CrossrefMetadata.CRMItem,
			},
		}
		content = append(content, c)
	}

	data, err = ReadList(content)
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
			} else if v.Affiliation != "" {
				affiliations = append(affiliations, commonmeta.Affiliation{
					Name: v.Affiliation,
				})
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
	fundingReferences = utils.DedupeSlice(fundingReferences)
	return fundingReferences, nil
}
