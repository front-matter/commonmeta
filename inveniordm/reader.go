// Package inveniordm provides functions to convert InvenioRDM metadata to/from the commonmeta metadata format.
package inveniordm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
)

// Content represents the InvenioRDM JSON API response.
type Content struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Inveniordm represents the InvenioRDM metadata.
type Inveniordm struct {
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
	Languages          []Language          `json:"language,omitempty"`
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
	ID          string       `json:"id,omitempty"`
	Title       AwardTitle   `json:"title,omitempty"`
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
	"Article":        "publication-preprint",
	"Book":           "book",
	"Dataset":        "dataset",
	"Image":          "image-other",
	"JournalArticle": "publication-article",
	"Presentation":   "presentation",
	"Software":       "software",
	"Other":          "other",
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

// CMToInvenioMappings maps Commonmeta identifier types to InvenioRDM identifier types
var CMToInvenioIdentifierMappings = map[string]string{
	"Ark":              "ark",
	"ArXiv":            "arxiv",
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

// Get retrieves InvenioRDM metadata.
func Get(id string) (Content, error) {
	var content Content
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	url := "https://zenodo.org/api/records/" + id
	resp, err := client.Get(url)
	if err != nil {
		return content, err
	}
	if resp.StatusCode != 200 {
		return content, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
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

// Read reads InvenioRDM JSON API response and converts it into Commonmeta metadata.
func Read(content Content) (commonmeta.Data, error) {
	var data commonmeta.Data
	data.ID = content.ID
	return data, nil
}
