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
	Type     string `xml:"type,attr,omitempty"`
	Provider string `xml:"provider,attr,omitempty"`
	Text     string `xml:",chardata"`
}

type DOIRecord struct {
	XMLName  xml.Name `xml:"doi_record"`
	Crossref Content  `xml:"crossref"`
}

// Content represents the Crossref XML metadata returned from Crossref. The type uses a struct
// instead of the slice for the Crossref metadata.
type Content struct {
	XMLName       xml.Name       `xml:"crossref"`
	Xmlns         string         `xml:"xmlns,attr,omitempty"`
	Book          *Book          `xml:"book,omitempty"`
	Conference    *Conference    `xml:"conference,omitempty"`
	Database      *Database      `xml:"database,omitempty"`
	Dissertation  *Dissertation  `xml:"dissertation,omitempty"`
	Journal       *Journal       `xml:"journal,omitempty"`
	PeerReview    *PeerReview    `xml:"peer_review,omitempty"`
	PostedContent *PostedContent `xml:"posted_content,omitempty"`
	SAComponent   *SAComponent   `xml:"sa_component,omitempty"`
	Standard      *Standard      `xml:"standard,omitempty"`
}

type Body struct {
	XMLName       xml.Name        `xml:"body"`
	Book          []Book          `xml:"book,omitempty"`
	Conference    []Conference    `xml:"conference,omitempty"`
	Database      []Database      `xml:"database,omitempty"`
	Dissertation  []Dissertation  `xml:"dissertation,omitempty"`
	Journal       []Journal       `xml:"journal,omitempty"`
	PeerReview    []PeerReview    `xml:"peer_review,omitempty"`
	PostedContent []PostedContent `xml:"posted_content,omitempty"`
	SAComponent   []SAComponent   `xml:"sa_component,omitempty"`
	Standard      []Standard      `xml:"standard,omitempty"`
}

type Abstract struct {
	XMLName      xml.Name `xml:"abstract"`
	Xmlns        string   `xml:"xmlns,attr"`
	Title        string   `xml:"title,attr,omitempty"`
	AbstractType string   `xml:"abstract-type,attr,omitempty"`
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

type Affiliations struct {
	XMLName     xml.Name      `xml:"affiliations"`
	Institution []Institution `xml:"institution,omitempty"`
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
	Archive []Archive `xml:"archive,omitempty"`
}

type Assertion struct {
	XMLName    xml.Name    `xml:"assertion"`
	Text       string      `xml:",chardata"`
	Name       string      `xml:"name,attr"`
	Provider   string      `xml:"provider,attr,omitempty"`
	Label      string      `xml:"label,attr,omitempty"`
	GroupName  string      `xml:"group_name,attr,omitempty"`
	GroupLabel string      `xml:"group_label,attr,omitempty"`
	Order      string      `xml:"order,attr,omitempty"`
	Assertion  []Assertion `xml:"assertion"`
}

type Book struct {
	XMLName         xml.Name        `xml:"book"`
	BookType        string          `xml:"book_type,attr"`
	BookMetadata    BookMetadata    `xml:"book_metadata"`
	BookSetMetadata BookSetMetadata `xml:"book_set_metadata"`
	ContentItem     ContentItem     `xml:"content_item"`
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

type BookSetMetadata struct {
	XMLName         xml.Name          `xml:"book_set_metadata"`
	Language        string            `xml:"language,attr"`
	SetMetadata     SetMetadata       `xml:"set_metadata"`
	Volume          string            `xml:"volume"`
	EditionNumber   string            `xml:"edition_number"`
	PublicationDate []PublicationDate `xml:"publication_date"`
	ISBN            []ISBN            `xml:"isbn"`
	Publisher       Publisher         `xml:"publisher"`
}

type Citation struct {
	XMLName            xml.Name `xml:"citation"`
	Key                string   `xml:"key,attr"`
	JournalTitle       string   `xml:"journal_title,omitempty"`
	Author             string   `xml:"author,omitempty"`
	Volume             string   `xml:"volume,omitempty"`
	FirstPage          string   `xml:"first_page,omitempty"`
	CYear              string   `xml:"cYear,omitempty"`
	ArticleTitle       string   `xml:"article_title,omitempty"`
	DOI                *DOI     `xml:"doi,omitempty"`
	UnstructedCitation string   `xml:"unstructured_citation,omitempty"`
}

type CitationList struct {
	XMLName  xml.Name   `xml:"citation_list"`
	Citation []Citation `xml:"citation,omitempty"`
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
	CitationList    CitationList      `xml:"citation_list,omitempty"`
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
	CitationList CitationList `xml:"citation_list,omitempty"`
}

type CreationDate struct {
	XMLName   xml.Name `xml:"creation_date"`
	MediaType string   `xml:"media_type,attr"`
	Month     string   `xml:"month"`
	Day       string   `xml:"day"`
	Year      string   `xml:"year"`
}

type Contributors struct {
	XMLName      xml.Name       `xml:"contributors"`
	Organization []Organization `xml:"organization,omitempty"`
	PersonName   []PersonName   `xml:"person_name,omitempty"`
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
	CitationList    CitationList `xml:"citation_list,omitempty"`
}

type DOIData struct {
	XMLName    xml.Name   `xml:"doi_data"`
	DOI        string     `xml:"doi"`
	Timestamp  string     `xml:"timestamp,omitempty"`
	Resource   string     `xml:"resource"`
	Collection Collection `xml:"collection,omitempty"`
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
	XMLName          xml.Name       `xml:"institution"`
	InstitutionName  string         `xml:"institution_name,omitempty"`
	InstitutionPlace string         `xml:"institution_place,omitempty"`
	InstitutionID    *InstitutionID `xml:"institution_id,omitempty"`
}

type InstitutionID struct {
	XMLName xml.Name `xml:"institution_id"`
	Type    string   `xml:"type,attr,omitempty"`
	Text    string   `xml:",chardata"`
}

type InterWorkRelation struct {
	XMLName          xml.Name `xml:"inter_work_relation"`
	RelationshipType string   `xml:"relationship-type,attr"`
	IdentifierType   string   `xml:"identifier-type,attr"`
	Text             string   `xml:",chardata"`
}

type IntraWorkRelation struct {
	Text             string `xml:",chardata"`
	RelationshipType string `xml:"relationship-type,attr"`
	IdentifierType   string `xml:"identifier-type,attr"`
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
	Crawler  string   `xml:"crawler,attr,omitempty"`
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
	JournalMetadata JournalMetadata `xml:"journal_metadata,omitempty"`
	JournalIssue    JournalIssue    `xml:"journal_issue,omitempty"`
	JournalArticle  JournalArticle  `xml:"journal_article,omitempty"`
}

// JournalArticle represents a journal article in Crossref XML metadata.
type JournalArticle struct {
	XMLName                   xml.Name          `xml:"journal_article"`
	Text                      string            `xml:",chardata"`
	PublicationType           string            `xml:"publication_type,attr,omitempty"`
	ReferenceDistributionOpts string            `xml:"reference_distribution_opts,attr,omitempty"`
	Titles                    Titles            `xml:"titles,omitempty"`
	Contributors              Contributors      `xml:"contributors,omitempty"`
	PublicationDate           []PublicationDate `xml:"publication_date"`
	PublisherItem             *PublisherItem    `xml:"publisher_item,omitempty"`
	Abstract                  []Abstract        `xml:"jats:abstract"`
	Pages                     *Pages            `xml:"pages,omitempty"`
	ISSN                      []ISSN            `xml:"issn"`
	Program                   []Program         `xml:"program"`
	Crossmark                 *Crossmark        `xml:"crossmark,omitempty"`
	ArchiveLocations          ArchiveLocations  `xml:"archive_locations"`
	DOIData                   DOIData           `xml:"doi_data"`
	CitationList              CitationList      `xml:"citation_list,omitempty"`
}

type JournalIssue struct {
	XMLName         xml.Name          `xml:"journal_issue"`
	PublicationDate []PublicationDate `xml:"publication_date"`
	JournalVolume   JournalVolume     `xml:"journal_volume,omitempty"`
	Issue           string            `xml:"issue,omitempty"`
	DOIData         *DOIData          `xml:"doi_data,omitempty"`
}

// JournalMetadata represents journal metadata in Crossref XML metadata.
type JournalMetadata struct {
	XMLName   xml.Name `xml:"journal_metadata"`
	Language  string   `xml:"language,attr,omitempty"`
	FullTitle string   `xml:"full_title,omitempty"`
	ISSN      []ISSN   `xml:"issn,omitempty"`
	DOIData   *DOIData `xml:"doi_data,omitempty"`
}

type JournalVolume struct {
	XMLName xml.Name `xml:"journal_volume"`
	Volume  string   `xml:"volume"`
}

type LicenseRef struct {
	Text      string `xml:",chardata"`
	AppliesTo string `xml:"applies_to,attr"`
}

// OrganizationName represents an organization in Crossref XML metadata.
type Organization struct {
	XMLName         xml.Name `xml:"organization"`
	ContributorRole string   `xml:"contributor_role,attr"`
	Sequence        string   `xml:"sequence,attr"`
	Text            string   `xml:",chardata"`
}

type OriginalLanguageTitle struct {
	Text     string `xml:",chardata"`
	Language string `xml:"language,attr"`
}

type P struct {
	XMLName xml.Name `xml:"p"`
	Xmlns   string   `xml:"xmlns,attr,omitempty"`
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
	Affiliations    *Affiliations `xml:"affiliations,omitempty"`
	Affiliation     string        `xml:"affiliation,omitempty"`
	ORCID           string        `xml:"ORCID,omitempty"`
}

// PostedContent represents posted content in Crossref XML metadata.
type PostedContent struct {
	XMLName        xml.Name        `xml:"posted_content"`
	Type           string          `xml:"type,attr"`
	Language       string          `xml:"language,attr,omitempty"`
	GroupTitle     string          `xml:"group_title,omitempty"`
	Contributors   Contributors    `xml:"contributors,omitempty"`
	Titles         Titles          `xml:"titles,omitempty"`
	PostedDate     PostedDate      `xml:"posted_date,omitempty"`
	AcceptanceDate *AcceptanceDate `xml:"acceptance_date,omitempty"`
	Institution    *Institution    `xml:"institution,omitempty"`
	ItemNumber     ItemNumber      `xml:"item_number,omitempty"`
	Abstract       []Abstract      `xml:"abstract"`
	Program        []Program       `xml:"program"`
	DOIData        DOIData         `xml:"doi_data"`
	CitationList   CitationList    `xml:"citation_list,omitempty"`
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
	Xmlns       string        `xml:"xmlns,attr"`
	Fr          string        `xml:"fr,attr,omitempty"`
	Name        string        `xml:"name,attr,omitempty"`
	Ai          string        `xml:"ai,attr,omitempty"`
	Rel         string        `xml:"rel,attr,omitempty"`
	Text        string        `xml:",chardata"`
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
	ItemNumber ItemNumber `xml:"item_number,omitempty"`
	Identifier struct {
		Text   string `xml:",chardata"`
		IDType string `xml:"id_type,attr"`
	} `xml:"identifier"`
}

type RelatedItem struct {
	XMLName           xml.Name           `xml:"related_item"`
	Text              string             `xml:",chardata"`
	Description       string             `xml:"description,omitempty"`
	InterWorkRelation *InterWorkRelation `xml:"inter_work_relation,omitempty"`
	IntraWorkRelation *IntraWorkRelation `xml:"intra_work_relation,omitempty"`
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

type SetMetadata struct {
	XMLName      xml.Name     `xml:"set_metadata"`
	Titles       Titles       `xml:"titles"`
	ISBN         []ISBN       `xml:"isbn"`
	Contributors Contributors `xml:"contributors"`
	DOIData      DOIData      `xml:"doi_data"`
}

type Standard struct {
	XMLName xml.Name `xml:"standard"`
	DOIData DOIData  `xml:"doi_data"`
}

// Titles represents the titles in Crossref XML metadata.
type Titles struct {
	XMLName               xml.Name               `xml:"titles"`
	Title                 string                 `xml:"title,omitempty"`
	Subtitle              string                 `xml:"subtitle,omitempty"`
	OriginalLanguageTitle *OriginalLanguageTitle `xml:"original_language_title,omitempty"`
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

var InterWorkRelationTypes = []string{
	"IsPartOf",
	"HasPart",
	"IsReviewOf",
	"HasReview",
	"IsRelatedMaterial",
	"HasRelatedMaterial",
}

var IntraWorkRelationTypes = []string{
	"IsIdenticalTo",
	"IsPreprintOf",
	"HasPreprint",
	"IsTranslationOf",
	"HasTranslation",
	"IsVersionOf",
	"HasVersion",
}

var OFRToRORMappings = map[string]string{
	"https://doi.org/10.13039/100000001":    "https://ror.org/021nxhr62",
	"https://doi.org/10.13039/501100000780": "https://ror.org/00k4n6c32",
	"https://doi.org/10.13039/501100007601": "https://ror.org/00k4n6c32",
	"https://doi.org/10.13039/501100001659": "https://ror.org/018mejw64",
	"https://doi.org/10.13039/501100006390": "https://ror.org/019whta54",
	"https://doi.org/10.13039/501100001711": "https://ror.org/00yjd3n13",
	"https://doi.org/10.13039/501100003043": "https://ror.org/04wfr2810",
}

var RORToOFRMappings = map[string]string{
	"https://ror.org/021nxhr62": "https://doi.org/10.13039/100000001",
	"https://ror.org/00k4n6c32": "https://doi.org/10.13039/501100000780",
	"https://ror.org/018mejw64": "https://doi.org/10.13039/501100001659",
	"https://ror.org/019whta54": "https://doi.org/10.13039/501100006390",
	"https://ror.org/00yjd3n13": "https://doi.org/10.13039/501100001711",
	"https://ror.org/04wfr2810": "https://doi.org/10.13039/501100003043",
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
func Get(pid string) (Query, error) {
	var query Query

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
			Body struct {
				Query Query `xml:"query"`
			} `xml:"body"`
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
		return query, errors.New("invalid DOI")
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	url := "https://api.crossref.org/works/" + doi + "/transform/application/vnd.crossref.unixsd+xml"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	v := "0.1"
	u := "info@front-matter.io"
	userAgent := fmt.Sprintf("commonmeta/%s (https://commonmeta.org/; mailto: %s)", v, u)
	req.Header.Set("User-Agent", userAgent)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return query, err
	}
	if resp.StatusCode >= 400 {
		return query, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error:", err)
		return query, err
	}
	err = xml.Unmarshal(body, &crossrefResult)
	if err != nil {
		fmt.Println("error:", err)
	}
	query = crossrefResult.QueryResult.Body.Query
	return query, err
}

// Read Crossref XML response and return work struct in Commonmeta format
func Read(query Query) (commonmeta.Data, error) {
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

	meta := query.DOIRecord.Crossref

	data.ID = doiutils.NormalizeDOI(query.DOI.Text)
	data.Type = CRToCMMappings[query.DOI.Type]
	if data.Type == "" {
		data.Type = "Other"
	}

	// fetch metadata depending on Crossref type (using the commonmeta vocabulary)
	switch data.Type {
	case "Article": // posted-content
		postedContent := meta.PostedContent
		abstract = append(abstract, postedContent.Abstract...)
		// archiveLocations not supported
		citationList = postedContent.CitationList
		contributors = postedContent.Contributors
		doiData = postedContent.DOIData
		itemNumber = postedContent.ItemNumber
		language = postedContent.Language
		// pages not supported
		program = append(program, postedContent.Program...)
		// use posted date as publication date
		publicationDate = append(publicationDate, PublicationDate{
			Year:  postedContent.PostedDate.Year,
			Month: postedContent.PostedDate.Month,
			Day:   postedContent.PostedDate.Day,
		})
		// use group title for subjects
		subjects = append(subjects, postedContent.GroupTitle)
		titles = postedContent.Titles
	case "Book":
		book := meta.Book
		abstract = book.BookMetadata.Abstract
		contributors = book.BookMetadata.Contributors
		citationList = book.ContentItem.CitationList
		doiData = book.BookMetadata.DOIData
		isbn = book.BookMetadata.ISBN
		language = book.BookMetadata.Language
		publicationDate = book.BookMetadata.PublicationDate
		pages = book.ContentItem.Pages
		titles = book.BookMetadata.Titles
	case "BookChapter":
		book := meta.Book
		abstract = book.BookMetadata.Abstract
		citationList = book.ContentItem.CitationList
		contributors = book.ContentItem.Contributors
		doiData = book.BookMetadata.DOIData
		isbn = book.BookMetadata.ISBN
		language = book.BookMetadata.Language
		publicationDate = book.ContentItem.PublicationDate
		pages = book.ContentItem.Pages
		titles = book.ContentItem.Titles
	case "BookPart":
	case "BookSection":
	case "BookSeries":
	case "BookSet":
	case "BookTrack":
	case "Component":
		component := meta.SAComponent
		doiData = component.ComponentList.Component[0].DOIData
	case "Database":
	case "Dataset":
		database := meta.Database
		containerTitle = database.DatabaseMetadata.Titles.Title
		contributors = database.Dataset.Contributors
		titles = database.Dataset.Titles
		// use creation date as publication date
		publicationDate = append(publicationDate, PublicationDate{
			Year:  database.Dataset.DatabaseDate.CreationDate.Year,
			Month: database.Dataset.DatabaseDate.CreationDate.Month,
			Day:   database.Dataset.DatabaseDate.CreationDate.Day,
		})
		doiData = database.Dataset.DOIData
	case "Dissertation":
		dissertation := meta.Dissertation
		contributors = Contributors{
			PersonName: dissertation.PersonName,
		}
		doiData = dissertation.DOIData
		// use approval date as publication date
		publicationDate = append(publicationDate, PublicationDate{
			Year:  dissertation.ApprovalDate.Year,
			Month: dissertation.ApprovalDate.Month,
			Day:   dissertation.ApprovalDate.Day,
		})
		titles = dissertation.Titles
	case "Entry":
	case "Grant":
	case "Journal":
		journal := meta.Journal
		containerTitle = journal.JournalMetadata.FullTitle
		language = journal.JournalMetadata.Language
		// doiData = journal.JournalMetadata.DOIData
	case "JournalArticle":
		journal := meta.Journal
		abstract = journal.JournalArticle.Abstract
		archiveLocations = journal.JournalArticle.ArchiveLocations
		citationList = journal.JournalArticle.CitationList
		containerTitle = journal.JournalMetadata.FullTitle
		contributors = journal.JournalArticle.Contributors
		// customMetadata = journal.JournalArticle.Crossmark.CustomMetadata
		doiData = journal.JournalArticle.DOIData
		issn = journal.JournalMetadata.ISSN
		issue = journal.JournalIssue.Issue
		itemNumber = journal.JournalArticle.PublisherItem.ItemNumber
		language = journal.JournalMetadata.Language
		// pages = *journal.JournalArticle.Pages
		program = append(program, journal.JournalArticle.Program...)
		publicationDate = journal.JournalArticle.PublicationDate
		titles = journal.JournalArticle.Titles
		volume = journal.JournalIssue.JournalVolume.Volume
	case "JournalIssue":
		journal := meta.Journal
		containerTitle = journal.JournalMetadata.FullTitle
		// doiData = journal.JournalIssue.DOIData
		issn = journal.JournalMetadata.ISSN
		language = journal.JournalMetadata.Language
	case "JournalVolume":
	case "Other":
	case "PeerReview":
		peerReview := meta.PeerReview
		contributors = peerReview.Contributors
		titles = peerReview.Titles
		program = append(program, peerReview.Program...)
		// use review date as publication date
		publicationDate = append(publicationDate, PublicationDate{
			Year:  peerReview.ReviewDate.Year,
			Month: peerReview.ReviewDate.Month,
			Day:   peerReview.ReviewDate.Day,
		})
		doiData = peerReview.DOIData
	case "Proceedings":
	case "ProceedingsArticle":
		conference := meta.Conference
		citationList = conference.ConferencePaper.CitationList
		containerTitle = conference.EventMetadata.ConferenceName
		contributors = conference.ConferencePaper.Contributors
		doiData = conference.ConferencePaper.DOIData
		isbn = conference.ProceedingsMetadata.ISBN
		program = append(program, conference.ConferencePaper.Crossmark.CustomMetadata.Program...)
		pages = conference.ConferencePaper.Pages
		publicationDate = conference.ConferencePaper.PublicationDate
		titles = conference.ConferencePaper.Titles
	case "ProceedingsSeries":
	case "ReferenceBook":
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
	for _, v := range query.CRMItem {
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
	// workaround until Crossref supports BlogPost as posted-content type
	if data.Type == "Article" && data.Publisher.Name == "Front Matter" {
		data.Type = "BlogPost"
	}

	if len(citationList.Citation) > 0 {
		for _, v := range citationList.Citation {
			var id string
			if v.DOI != nil {
				id = doiutils.NormalizeDOI(v.DOI.Text)
			}
			if id != "" {
				reference := commonmeta.Reference{
					Key:             v.Key,
					ID:              id,
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
	// if titles.OriginalLanguageTitle.Text != "" {
	// 	data.Titles = append(data.Titles, commonmeta.Title{
	// 		Title:    titles.OriginalLanguageTitle.Text,
	// 		Type:     "TranslatedTitle",
	// 		Language: titles.OriginalLanguageTitle.Language,
	// 	})
	// }

	data.URL = doiData.Resource

	return data, nil
}

// ReadAll reads a list of Crossref XML responses and returns a list of works in Commonmeta format
func ReadAll(query []Query) ([]commonmeta.Data, error) {
	var data []commonmeta.Data
	for _, v := range query {
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
	var query Query

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
	err = decoder.Decode(&query)
	if err != nil {
		return data, err
	}
	data, err = Read(query)
	if err != nil {
		return data, err
	}
	return data, nil
}

// LoadAll loads the metadata for a list of works from an XML file and converts it to the Commonmeta format
func LoadAll(filename string) ([]commonmeta.Data, error) {
	type Body struct {
		XMLName          xml.Name `xml:"body"`
		CrossrefMetadata Query    `xml:"crossref_metadata"`
	}

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
							Body Body `xml:"body"`
						} `xml:"query_result"`
					}
				} `xml:"metadata"`
			} `xml:"record"`
		} `xml:"ListRecords"`
	}

	var data []commonmeta.Data
	var query []Query
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
		q := Query{
			Status:    "resolved",
			DOI:       v.Metadata.CrossrefResult.QueryResult.Body.CrossrefMetadata.DOI,
			DOIRecord: v.Metadata.CrossrefResult.QueryResult.Body.CrossrefMetadata.DOIRecord,
			CRMItem:   v.Metadata.CrossrefResult.QueryResult.Body.CrossrefMetadata.CRMItem,
		}
		query = append(query, q)
	}

	data, err = ReadAll(query)
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
			if v.Affiliations != nil || v.Affiliation != "" {
				var affiliations []*commonmeta.Affiliation
				if v.Affiliations != nil {
					for _, i := range v.Affiliations.Institution {
						if i.InstitutionName != "" {
							if i.InstitutionID != nil && i.InstitutionID.Text != "" {
								InstitutionID := utils.NormalizeROR(i.InstitutionID.Text)
								affiliations = append(affiliations, &commonmeta.Affiliation{
									ID:   InstitutionID,
									Name: i.InstitutionName,
								})
							} else {
								affiliations = append(affiliations, &commonmeta.Affiliation{
									Name: i.InstitutionName,
								})
							}
						}
					}
				} else if v.Affiliation != "" {
					affiliations = append(affiliations, &commonmeta.Affiliation{
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
				contributors = append(contributors, contributor)
			} else {
				contributor := commonmeta.Contributor{
					ID:               ID,
					Type:             Type,
					GivenName:        v.GivenName,
					FamilyName:       v.Surname,
					Name:             "",
					ContributorRoles: []string{"Author"},
				}
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

// Type represents the Crossref type of a work
func (q Query) Type() string {
	switch q.DOI.Type {
	case "book_title":
		if q.DOIRecord.Crossref.Book.BookSetMetadata.SetMetadata.DOIData.DOI != "" {
			return "BookSet"
		}
		return "Book"
	case "book_content":
		switch q.DOIRecord.Crossref.Book.ContentItem.ComponentType {
		case "other":
			return "Other"
		case "part":
			return "BookPart"
		case "reference_entry":
			return "Entry"
		case "section":
			return "BookSection"
		case "track":
			return "BookTrack"
		default:
			return "BookChapter"
		}
	case "book_series":
		return "BookSeries"
	case "component":
		return "Component"
	case "conference":
		return "Proceedings"
	case "conference_paper":
		return "ProceedingsArticle"
	case "conference_series":
		return "ProceedingsSeries"
	case "conference_title":
		return "Proceedings"
	case "database_title":
		return "Database"
	case "dataset":
		return "Dataset"
	case "dissertation":
		return "Dissertation"
	case "grant":
		return "Grant"
	case "journal_article":
		return "JournalArticle"
	case "journal_issue":
		return "JournalIssue"
	case "journal_volume":
		return "JournalVolume"
	case "journal_title":
		return "Journal"
	case "peer_review":
		return "PeerReview"
	case "posted_content":
		return "Article"
	case "report-paper_content":
		return "ReportComponent"
	case "report-paper_series":
		return "ReportSeries"
	case "report-paper_title":
		return "Report"
	case "standard_title":
		return "Standard"
	default:
		return "Other"
	}
}
