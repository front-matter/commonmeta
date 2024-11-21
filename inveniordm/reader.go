// Package inveniordm provides functions to convert InvenioRDM metadata to/from the commonmeta metadata format.
package inveniordm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/utils"
)

// Content represents the InvenioRDM JSON API response.
type Content struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Query represents the InvenioRDM JSON API query.
type Query struct {
	Hits struct {
		Hits  []Inveniordm `json:"hits"`
		Total int          `json:"total"`
	} `json:"hits"`
}

// Inveniordm represents the InvenioRDM metadata.
type Inveniordm struct {
	ID           string       `json:"id,omitempty"`
	Pids         Pids         `json:"pids"`
	Access       Access       `json:"access"`
	Files        Files        `json:"files"`
	Metadata     Metadata     `json:"metadata"`
	CustomFields CustomFields `json:"custom_fields"`
}

type Affiliation struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

type Pids struct {
	DOI DOI `json:"doi"`
}

type Access struct {
	Record string `json:"record"`
	Files  string `json:"files"`
}

type Files struct {
	Enabled bool `json:"enabled"`
}

type Metadata struct {
	ResourceType       ResourceType        `json:"resource_type"`
	Creators           []Creator           `json:"creators"`
	Title              string              `json:"title"`
	Publisher          string              `json:"publisher,omitempty"`
	PublicationDate    string              `json:"publication_date"`
	Subjects           []Subject           `json:"subjects,omitempty"`
	Dates              []Date              `json:"dates,omitempty"`
	Description        string              `json:"description,omitempty"`
	Rights             []Right             `json:"rights,omitempty"`
	Languages          []Language          `json:"languages,omitempty"`
	Identifiers        []Identifier        `json:"identifiers,omitempty"`
	RelatedIdentifiers []RelatedIdentifier `json:"related_identifiers,omitempty"`
	Funding            []Funding           `json:"funding,omitempty"`
	Version            string              `json:"version,omitempty"`
}

type CustomFields struct {
	Journal      Journal `json:"journal:journal,omitempty"`
	ContentText  string  `json:"rs:content_text,omitempty"`
	FeatureImage string  `json:"rs:image,omitempty"`
}

type Award struct {
	ID string `json:"id,omitempty"`
	// Title       AwardTitle   `json:"title,omitempty"`
	Number      string       `json:"number,omitempty"`
	Identifiers []Identifier `json:"identifiers,omitempty"`
}

type AwardTitle struct {
	En string `json:"en,omitempty"`
}

type Date struct {
	Date string `json:"date"`
	Type Type   `json:"type"`
}
type DOI struct {
	Identifier string `json:"identifier"`
	Provider   string `json:"provider"`
}

type Funder struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

type Funding struct {
	Funder Funder `json:"funder"`
	Award  Award  `json:"award,omitempty"`
}

type ResourceType struct {
	ID string `json:"id"`
}

type Creator struct {
	PersonOrOrg  PersonOrOrg   `json:"person_or_org"`
	Affiliations []Affiliation `json:"affiliations,omitempty"`
}

type Identifier struct {
	Identifier string `json:"identifier"`
	Scheme     string `json:"scheme,omitempty"`
}

type PersonOrOrg struct {
	Type        string       `json:"type"`
	Name        string       `json:"name,omitempty"`
	GivenName   string       `json:"given_name,omitempty"`
	FamilyName  string       `json:"family_name,omitempty"`
	Identifiers []Identifier `json:"identifiers,omitempty"`
}

type RelatedIdentifier struct {
	Identifier   string `json:"identifier"`
	Scheme       string `json:"scheme"`
	RelationType Type   `json:"relation_type"`
}

type Subject struct {
	ID      string `json:"id,omitempty"`
	Subject string `json:"subject,omitempty"`
	Scheme  string `json:"scheme,omitempty"`
}

type Right struct {
	ID string `json:"id"`
}

type Language struct {
	ID string `json:"id"`
}

type Journal struct {
	Title  string `json:"title,omitempty"`
	Volume string `json:"volume,omitempty"`
	Issue  string `json:"issue,omitempty"`
	Pages  string `json:"pages,omitempty"`
	ISSN   string `json:"issn,omitempty"`
}

type Type struct {
	ID string `json:"id"`
}

// Awards represents the InvenioRDM awards.yaml file.
type AwardVocabulary struct {
	ID    string `yaml:"id"`
	Title struct {
		En string `yaml:"en"`
	} `yaml:"title"`
	Number  string `yaml:"number"`
	Acronym string `yaml:"acronym"`
	Funder  struct {
		ID   string `yaml:"id"`
		Name string `yaml:"name"`
	} `yaml:"funder"`
	Identifiers []struct {
		Identifier string `yaml:"identifier"`
		Scheme     string `yaml:"scheme"`
	}
}

// InvenioToCMMappings maps InvenioRDM resource types to Commonmeta types
// source: https://github.com/zenodo/zenodo/blob/master/zenodo/modules/records/data/objecttypes.json
var InvenioToCMMappings = map[string]string{
	"book":                  "Book",
	"section":               "BookChapter",
	"conferencepaper":       "ProceedingsArticle",
	"patent":                "Patent",
	"publication":           "JournalArticle",
	"publication-preprint":  "Article",
	"report":                "Report",
	"softwaredocumentation": "Software",
	"thesis":                "Dissertation",
	"technicalnote":         "Report",
	"workingpaper":          "Report",
	"datamanagementplan":    "OutputManagementPlan",
	"annotationcollection":  "Collection",
	"taxonomictreatment":    "Collection",
	"peerreview":            "PeerReview",
	"poster":                "Presentation",
	"presentation":          "Presentation",
	"dataset":               "Dataset",
	"figure":                "Image",
	"plot":                  "Image",
	"drawing":               "Image",
	"photo":                 "Image",
	"image":                 "Image",
	"video":                 "Audiovisual",
	"software":              "Software",
	"lesson":                "InteractiveResource",
	"physicalobject":        "PhysicalObject",
	"workflow":              "Workflow",
	"other":                 "Other",
}

// CMTOInvenioMappings maps Commonmeta types to InvenioRDM resource types
var CMToInvenioMappings = map[string]string{
	"Article":               "publication-preprint",
	"Audiovisual":           "video",
	"Book":                  "publication-book",
	"BookChapter":           "publication-section",
	"Collection":            "publication-annotationcollection",
	"ComputationalNotebook": "software-computationalnotebook",
	"Dataset":               "dataset",
	"Dissertation":          "publication-thesis",
	"Document":              "publication",
	"Entry":                 "publication",
	"Event":                 "event",
	"Figure":                "image-figure",
	"Image":                 "image",
	"Instrument":            "other",
	"Journal":               "publication-journal",
	"JournalArticle":        "publication-article",
	"LegalDocument":         "publication",
	"Manuscript":            "publication",
	"Map":                   "other",
	"Patent":                "patent",
	"PersonalCommunication": "publication",
	"PhysicalObject":        "physicalobject",
	"Post":                  "publication",
	"Presentation":          "presentation",
	"ProceedingsArticle":    "publication-conferencepaper",
	"Proceedings":           "publication-conferenceproceeding",
	"Report":                "publication-report",
	"Review":                "publication-peerreview",
	"Software":              "software",
	"Sound":                 "audio",
	"Standard":              "publication-standard",
	"WebPage":               "publication",
	"Workflow":              "workflow",
	"Other":                 "other",
}

// FOSMappings maps OECD FOS strings to OECD FOS identifiers
var FOSMappings = map[string]string{
	"Natural sciences":                         "http://www.oecd.org/science/inno/38235147.pdf?1",
	"Mathematics":                              "http://www.oecd.org/science/inno/38235147.pdf?1.1",
	"Computer and information sciences":        "http://www.oecd.org/science/inno/38235147.pdf?1.2",
	"Physical sciences":                        "http://www.oecd.org/science/inno/38235147.pdf?1.3",
	"Chemical sciences":                        "http://www.oecd.org/science/inno/38235147.pdf?1.4",
	"Earth and related environmental sciences": "http://www.oecd.org/science/inno/38235147.pdf?1.5",
	"Biological sciences":                      "http://www.oecd.org/science/inno/38235147.pdf?1.6",
	"Other natural sciences":                   "http://www.oecd.org/science/inno/38235147.pdf?1.7",
	"Engineering and technology":               "http://www.oecd.org/science/inno/38235147.pdf?2",
	"Civil engineering":                        "http://www.oecd.org/science/inno/38235147.pdf?2.1",
	"Electrical engineering, electronic engineering, information engineering": "http://www.oecd.org/science/inno/38235147.pdf?2.2",
	"Mechanical engineering":               "http://www.oecd.org/science/inno/38235147.pdf?2.3",
	"Chemical engineering":                 "http://www.oecd.org/science/inno/38235147.pdf?2.4",
	"Materials engineering":                "http://www.oecd.org/science/inno/38235147.pdf?2.5",
	"Medical engineering":                  "http://www.oecd.org/science/inno/38235147.pdf?2.6",
	"Environmental engineering":            "http://www.oecd.org/science/inno/38235147.pdf?2.7",
	"Environmental biotechnology":          "http://www.oecd.org/science/inno/38235147.pdf?2.8",
	"Industrial biotechnology":             "http://www.oecd.org/science/inno/38235147.pdf?2.9",
	"Nano technology":                      "http://www.oecd.org/science/inno/38235147.pdf?2.10",
	"Other engineering and technologies":   "http://www.oecd.org/science/inno/38235147.pdf?2.11",
	"Medical and health sciences":          "http://www.oecd.org/science/inno/38235147.pdf?3",
	"Basic medicine":                       "http://www.oecd.org/science/inno/38235147.pdf?3.1",
	"Clinical medicine":                    "http://www.oecd.org/science/inno/38235147.pdf?3.2",
	"Health sciences":                      "http://www.oecd.org/science/inno/38235147.pdf?3.3",
	"Health biotechnology":                 "http://www.oecd.org/science/inno/38235147.pdf?3.4",
	"Other medical sciences":               "http://www.oecd.org/science/inno/38235147.pdf?3.5",
	"Agricultural sciences":                "http://www.oecd.org/science/inno/38235147.pdf?4",
	"Agriculture, forestry, and fisheries": "http://www.oecd.org/science/inno/38235147.pdf?4.1",
	"Animal and dairy science":             "http://www.oecd.org/science/inno/38235147",
	"Veterinary science":                   "http://www.oecd.org/science/inno/38235147",
	"Agricultural biotechnology":           "http://www.oecd.org/science/inno/38235147",
	"Other agricultural sciences":          "http://www.oecd.org/science/inno/38235147",
	"Social science":                       "http://www.oecd.org/science/inno/38235147.pdf?5",
	"Psychology":                           "http://www.oecd.org/science/inno/38235147.pdf?5.1",
	"Economics and business":               "http://www.oecd.org/science/inno/38235147.pdf?5.2",
	"Educational sciences":                 "http://www.oecd.org/science/inno/38235147.pdf?5.3",
	"Sociology":                            "http://www.oecd.org/science/inno/38235147.pdf?5.4",
	"Law":                                  "http://www.oecd.org/science/inno/38235147.pdf?5.5",
	"Political science":                    "http://www.oecd.org/science/inno/38235147.pdf?5.6",
	"Social and economic geography":        "http://www.oecd.org/science/inno/38235147.pdf?5.7",
	"Media and communications":             "http://www.oecd.org/science/inno/38235147.pdf?5.8",
	"Other social sciences":                "http://www.oecd.org/science/inno/38235147.pdf?5.9",
	"Humanities":                           "http://www.oecd.org/science/inno/38235147.pdf?6",
	"History and archaeology":              "http://www.oecd.org/science/inno/38235147.pdf?6.1",
	"Languages and literature":             "http://www.oecd.org/science/inno/38235147.pdf?6.2",
	"Philosophy, ethics and religion":      "http://www.oecd.org/science/inno/38235147.pdf?6.3",
	"Arts (arts, history of arts, performing arts, music)": "http://www.oecd.org/science/inno/38235147.pdf?6.4",
	"Other humanities": "http://www.oecd.org/science/inno/38235147.pdf?6.5",
}

// CMToInvenioIdentifierMappings maps Commonmeta identifier types to InvenioRDM identifier types
var CMToInvenioIdentifierMappings = map[string]string{
	"Ark":              "ark",
	"arXiv":            "arxiv",
	"Bibcode":          "ads",
	"CrossrefFunderID": "crossreffunderid",
	"DOI":              "doi",
	"EAN13":            "ean13",
	"EISSN":            "eissn",
	"GRID":             "grid",
	"Handle":           "handle",
	"IGSN":             "igsn",
	"ISBN":             "isbn",
	"ISNI":             "isni",
	"ISSN":             "issn",
	"ISTC":             "istc",
	"LISSN":            "lissn",
	"LSID":             "lsid",
	"PMID":             "pmid",
	"PURL":             "purl",
	"UPC":              "upc",
	"URL":              "url",
	"URN":              "urn",
	"W3ID":             "w3id",
	"GUID":             "guid",
	"UUID":             "uuid",
	"Other":            "other",
}

// CMToInvenioRelationTypeMappings maps Commonmeta identifier types to InvenioRDM identifier types
var CMToInvenioRelationTypeMappings = map[string]string{
	"IsCitedBy":         "iscitedby",
	"Cites":             "cites",
	"IsSupplementTo":    "issupplementto",
	"IsSupplementedBy":  "issupplementedby",
	"IsContinuedBy":     "iscontinuedby",
	"Continues":         "continues",
	"IsNewVersionOf":    "isnewversionof",
	"IsPreviousVersion": "ispreviousversion",
	"IsPartOf":          "ispartof",
	"HasPart":           "haspart",
	"IsReferencedBy":    "isreferencedby",
	"References":        "references",
	"IsDocumentedBy":    "isdocumentedby",
	"Documents":         "documents",
	"IsCompiledBy":      "iscompiledby",
	"Compiles":          "compiles",
	"IsVariantFormOf":   "isvariantformof",
	"IsOriginalFormOf":  "isoriginalformof",
	"IsIdenticalTo":     "isidenticalto",
	"IsReviewedBy":      "isreviewedby",
	"Reviews":           "reviews",
	"IsDerivedFrom":     "isderivedfrom",
	"IsSourceOf":        "issourceof",
	"Describes":         "describes",
	"IsDescribedBy":     "isdescribedby",
	"IsMetadataFor":     "ismetadatafor",
	"HasMetadata":       "hasmetadata",
	"IsAnnotatedBy":     "isannotatedby",
	"Annotates":         "annotates",
	"IsCorrectedBy":     "iscorrectedby",
	"Corrects":          "corrects",
}

// Fetch fetches InvenioRDM metadata and returns Commonmeta metadata.
func Fetch(str string) (commonmeta.Data, error) {
	var data commonmeta.Data
	id, _ := utils.ValidateID(str)
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

// Get retrieves InvenioRDM metadata.
func Get(id string) (Content, error) {
	var content Content
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(id)
	if err != nil {
		return content, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return content, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return content, err
	}
	err = json.Unmarshal(body, &content)
	if err != nil {
		fmt.Println("error:", err)
	}
	return content, err
}

// SearchByDOI searches InvenioRDM records by external DOI.
func SearchByDOI(doi string, client *InvenioRDMClient) (string, error) {
	var query Query
	doistr := doiutils.EscapeDOI(doi)
	requestURL := fmt.Sprintf("https://%s/api/records?q=doi:%s", client.Host, doistr)
	req, _ := http.NewRequest(http.MethodGet, requestURL, nil)
	req.Header = http.Header{
		"Content-Type": {"application/json"},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &query)
	if err != nil {
		return "", err
	}
	if query.Hits.Total == 0 {
		return "", nil
	} else {
		return query.Hits.Hits[0].ID, nil
	}
}

// SearchBySlug searches InvenioRDM communities by slug.
func SearchBySlug(slug string, client *InvenioRDMClient) (string, error) {
	var query Query
	requestURL := fmt.Sprintf("https://%s/api/communities?q=slug:%s", client.Host, slug)
	req, _ := http.NewRequest(http.MethodGet, requestURL, nil)
	req.Header = http.Header{
		"Content-Type": {"application/json"},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	err = json.Unmarshal(body, &query)
	if err != nil {
		return "", err
	}
	if query.Hits.Total == 0 {
		return "", nil
	} else {
		return query.Hits.Hits[0].ID, nil
	}
}

// Read reads InvenioRDM JSON API response and converts it into Commonmeta metadata.
func Read(content Content) (commonmeta.Data, error) {
	var data commonmeta.Data
	data.ID = content.ID
	return data, nil
}
