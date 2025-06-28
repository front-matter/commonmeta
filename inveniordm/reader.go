// Package inveniordm provides functions to convert InvenioRDM metadata to/from the commonmeta metadata format.
package inveniordm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"slices"
	"strconv"
	"time"

	"github.com/muesli/cache2go"

	"github.com/front-matter/commonmeta/authorutils"
	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/ror"
	"github.com/front-matter/commonmeta/spdx"
	"github.com/front-matter/commonmeta/utils"
)

// Query represents the InvenioRDM JSON API query.
type Query struct {
	Hits struct {
		Hits  []Content `json:"hits"`
		Total int       `json:"total"`
	} `json:"hits"`
}

// Inveniordm represents the InvenioRDM metadata.
type Inveniordm struct {
	ID           string       `json:"id,omitempty"`
	Parent       Parent       `json:"parent"`
	Pids         Pids         `json:"pids"`
	Access       Access       `json:"access"`
	Files        Files        `json:"files"`
	Metadata     Metadata     `json:"metadata"`
	CustomFields CustomFields `json:"custom_fields"`
}

// Content represents the Inveniordm metadata returned from an Inveniordm API. The type is more
// flexible than the Inveniordm type, allowing for different formats of some metadata, e.g.
// customized instances such as Zenodo.
type Content struct {
	*Inveniordm
	ID       interface{}     `json:"id,omitempty"`
	DOI      string          `json:"doi,omitempty"`
	Files    json.RawMessage `json:"files,omitempty"`
	Metadata MetadataJSON    `json:"metadata"`
}

type Affiliation struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

type Parent struct {
	ID          string      `json:"id"`
	Communities Communities `json:"communities"`
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
	Funding            []Funding           `json:"funding,omitempty"`
	Dates              []Date              `json:"dates,omitempty"`
	Description        string              `json:"description,omitempty"`
	Grants             []Grant             `json:"grants,omitempty"`
	Identifiers        []Identifier        `json:"identifiers,omitempty"`
	Keywords           []string            `json:"keywords,omitempty"`
	Language           string              `json:"language,omitempty"`
	Languages          []Language          `json:"languages,omitempty"`
	License            *License            `json:"license,omitempty"`
	Publisher          string              `json:"publisher,omitempty"`
	PublicationDate    string              `json:"publication_date"`
	References         []Reference         `json:"references,omitempty"`
	RelatedIdentifiers []RelatedIdentifier `json:"related_identifiers,omitempty"`
	Rights             []Right             `json:"rights,omitempty"`
	Subjects           []Subject           `json:"subjects,omitempty"`
	Title              string              `json:"title"`
	Version            string              `json:"version,omitempty"`
}

type MetadataJSON struct {
	*Metadata
	Dates []DateJSON `json:"dates,omitempty"`
}

type Award struct {
	ID          string       `json:"id,omitempty"`
	Number      string       `json:"number,omitempty"`
	Title       AwardTitle   `json:"title,omitempty"`
	Identifiers []Identifier `json:"identifiers,omitempty"`
}

type AwardTitle struct {
	En string `json:"en,omitempty"`
}

type Communities struct {
	IDS     []string    `json:"ids"`
	Default string      `json:"default"`
	Entries []Community `json:"entries"`
}

type Community struct {
	ID           string                `json:"id"`
	Slug         string                `json:"slug"`
	Created      string                `json:"created"`
	Updated      string                `json:"updated"`
	Metadata     CommunityMetadata     `json:"metadata"`
	CustomFields CommunityCustomFields `json:"custom_fields"`
}

type CommunityCustomFields struct {
	ISSN     string `json:"rs:issn,omitempty"`
	FeedURL  string `json:"rs:feed_url,omitempty"`
	Language string `json:"rs:language,omitempty"`
	License  string `json:"rs:license,omitempty"`
	Category string `json:"rs:category,omitempty"`
}

type CommunityMetadata struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Type        struct {
		ID string `json:"id"`
	}
	Website string `json:"website,omitempty"`
}

type Creator struct {
	PersonOrOrg  PersonOrOrg   `json:"person_or_org"`
	Affiliations []Affiliation `json:"affiliations,omitempty"`
	Name         string        `json:"name"`
	ORCID        string        `json:"orcid,omitempty"`
	Affiliation  string        `json:"affiliation,omitempty"`
}

type CustomFields struct {
	Journal      Journal `json:"journal:journal,omitempty"`
	ContentHTML  string  `json:"rs:content_html,omitempty"`
	FeatureImage string  `json:"rs:image,omitempty"`
	Generator    string  `json:"rs:generator,omitempty"`
}

type Date struct {
	Date string `json:"date"`
	Type Type   `json:"type"`
}

type DateJSON struct {
	Date string          `json:"date"`
	Type json.RawMessage `json:"type"`
}

type DOI struct {
	Identifier string `json:"identifier"`
	Provider   string `json:"provider"`
}

type Funder struct {
	ID   string `json:"id,omitempty"`
	DOI  string `json:"doi,omitempty"`
	Name string `json:"name"`
}

type Funding struct {
	Funder Funder `json:"funder"`
	Award  Award  `json:"award"`
}

type Grant struct {
	Code   string `json:"code,omitempty"`
	Funder Funder `json:"funder"`
	Title  string `json:"title,omitempty"`
	URL    string `json:"url,omitempty"`
}

type Identifier struct {
	Identifier string `json:"identifier"`
	Scheme     string `json:"scheme,omitempty"`
}

type License struct {
	ID string `json:"id,omitempty"`
}

type PersonOrOrg struct {
	Type        string       `json:"type"`
	Name        string       `json:"name,omitempty"`
	GivenName   string       `json:"given_name,omitempty"`
	FamilyName  string       `json:"family_name,omitempty"`
	Identifiers []Identifier `json:"identifiers,omitempty"`
}

type Reference struct {
	Reference  string `json:"reference"`
	Scheme     string `json:"scheme"`
	Identifier string `json:"identifier"`
}

type RelatedIdentifier struct {
	Identifier   string `json:"identifier"`
	Scheme       string `json:"scheme"`
	RelationType Type   `json:"relation_type"`
}

type ResourceType struct {
	ID      string `json:"id,omitempty"`
	Subtype string `json:"subtype,omitempty"`
	Type    string `json:"type,omitempty"`
}

type Subject struct {
	ID      string `json:"id,omitempty"`
	Subject string `json:"subject,omitempty"`
	Scheme  string `json:"scheme,omitempty"`
}

type Right struct {
	ID    string `json:"id"`
	Props struct {
		URL string `json:"url,omitempty"`
	} `json:"props,omitempty"`
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
	"annotationcollection":  "Collection",
	"book":                  "Book",
	"conferencepaper":       "ProceedingsArticle",
	"datamanagementplan":    "OutputManagementPlan",
	"dataset":               "Dataset",
	"drawing":               "Image",
	"figure":                "Image",
	"image":                 "Image",
	"lesson":                "InteractiveResource",
	"patent":                "Patent",
	"peerreview":            "PeerReview",
	"photo":                 "Image",
	"physicalobject":        "PhysicalObject",
	"plot":                  "Image",
	"poster":                "Poster",
	"presentation":          "Presentation",
	"preprint":              "Article",
	"publication":           "JournalArticle",
	"publication-blogpost":  "BlogPost",
	"publication-preprint":  "BlogPost", //"Article"
	"report":                "Report",
	"section":               "BookChapter",
	"software":              "Software",
	"softwaredocumentation": "Software",
	"taxonomictreatment":    "Collection",
	"technicalnote":         "Report",
	"thesis":                "Dissertation",
	"video":                 "Audiovisual",
	"workflow":              "Workflow",
	"workingpaper":          "Report",
	"other":                 "Other",
}

// CMTOInvenioMappings maps Commonmeta types to InvenioRDM resource types
var CMToInvenioMappings = map[string]string{
	"Article":               "publication-preprint",
	"Audiovisual":           "video",
	"BlogPost":              "publication-blogpost",
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
	"Poster":                "poster",
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

// InvenioToCMIdentifierMappings maps Commonmeta identifier types to InvenioRDM identifier types
var InvenioToCMIdentifierMappings = map[string]string{
	"ark":              "Ark",
	"arxiv":            "arXiv",
	"ads":              "Bibcode",
	"crossreffunderid": "CrossrefFunderID",
	"doi":              "DOI",
	"ean13":            "EAN13",
	"eissn":            "EISSN",
	"grid":             "GRID",
	"handle":           "Handle",
	"igsn":             "IGSN",
	"isbn":             "ISBN",
	"isni":             "ISNI",
	"issn":             "ISSN",
	"istc":             "ISTC",
	"lissn":            "LISSN",
	"lsid":             "LSID",
	"pmid":             "PMID",
	"purl":             "PURL",
	"upc":              "UPC",
	"url":              "URL",
	"urn":              "URN",
	"w3id":             "W3ID",
	"guid":             "GUID",
	"uuid":             "UUID",
	"other":            "Other",
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

// InvenioToCMRelationTypeMappings maps Commonmeta identifier types to InvenioRDM identifier types
var InvenioToCMRelationTypeMappings = map[string]string{
	"iscitedby":         "IsCitedBy",
	"issupplementto":    "IsSupplementTo",
	"issupplementedby":  "IsSupplementedBy",
	"iscontinuedby":     "IsContinuedBy",
	"continues":         "Continues",
	"isnewversionof":    "IsNewVersionOf",
	"ispreviousversion": "IsPreviousVersion",
	"ispartof":          "IsPartOf",
	"haspart":           "HasPart",
	"isreferencedby":    "IsReferencedBy",
	"isdocumentedby":    "IsDocumentedBy",
	"documents":         "Documents",
	"iscompiledby":      "IsCompiledBy",
	"compiles":          "Compiles",
	"isvariantformof":   "IsVariantFormOf",
	"isoriginalformof":  "IsOriginalFormOf",
	"isidenticalto":     "IsIdenticalTo",
	"isreviewedby":      "IsReviewedBy",
	"reviews":           "Reviews",
	"isderivedfrom":     "IsDerivedFrom",
	"issourceof":        "IsSourceOf",
	"describes":         "Describes",
	"isdescribedby":     "IsDescribedBy",
	"ismetadatafor":     "IsMetadataFor",
	"hasmetadata":       "HasMetadata",
	"isannotatedby":     "IsAnnotatedBy",
	"annotates":         "Annotates",
	"iscorrectedby":     "IsCorrectedBy",
	"corrects":          "Corrects",
}

// CMToInvenioRelationTypeMappings maps Commonmeta relation_types to InvenioRDM relation_types
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
	"IsReviewOf":        "reviews",
	"HasReview":         "isreviewedby",
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

// LicenseMappings maps InvenioRDM license types to Commonmeta license types
var LicenseMappings = map[string]string{
	"cc-by-3.0":       "CC-BY-3.0",
	"cc-by-4.0":       "CC-BY-4.0",
	"cc-by-nc-3.0":    "CC-BY-NC-3.0",
	"cc-by-nc-4.0":    "CC-BY-NC-4.0",
	"cc-by-nc-nd-3.0": "CC-BY-NC-ND-3.0",
	"cc-by-nc-nd-4.0": "CC-BY-NC-ND-4.0",
	"cc-by-nc-sa-3.0": "CC-BY-NC-SA-3.0",
	"cc-by-nc-sa-4.0": "CC-BY-NC-SA-4.0",
	"cc-by-nd-3.0":    "CC-BY-ND-3.0",
	"cc-by-nd-4.0":    "CC-BY-ND-4.0",
	"cc-by-sa-3.0":    "CC-BY-SA-3.0",
	"cc-by-sa-4.0":    "CC-BY-SA-4.0",
	"cc0-1.0":         "CC0-1.0",
	"mit":             "MIT",
	"apache-2.0":      "Apache-2.0",
	"gpl-3.0":         "GPL-3.0",
}

// CommunityTypes maps InvenioRDM community types to Commonmeta container types
var CommunityTypes = map[string]string{
	"blog": "Blog",
}

// CommunityTranslations maps Community names in different languages to the community slug
// Also maps synonyms
var CommunityTranslations = map[string]string{
	"ai":                         "artificialintelligence",
	"llms":                       "artificialintelligence",
	"book%20review":              "bookreview",
	"bjps%20review%20of%20books": "bookreview",
	"books":                      "bookreview",
	"nachrichten":                "news",
	"opencitations":              "researchassessment",
	"papers":                     "researchblogging",
	"urheberrecht":               "copyright",
	"workshop":                   "events",
	"veranstaltungen":            "events",
	"veranstaltungshinweise":     "events",
	"asapbio":                    "preprints",
	"biorxiv":                    "preprints",
	"runiverse":                  "r",
	"bericht":                    "report",
}

// Fetch fetches InvenioRDM metadata and returns Commonmeta metadata.
func Fetch(str string, match bool) (commonmeta.Data, error) {
	var data commonmeta.Data
	id, _ := utils.ValidateID(str)
	content, err := Get(id)
	if err != nil {
		return data, err
	}
	data, err = Read(content, match)
	return data, err
}

// FetchAll gets the metadata for a list of records from a InvenioRDM community and returns Commonmeta metadata.
func FetchAll(number int, page int, host string, community string, subject string, type_ string, year string, language string, orcid string, affiliation string, ror string, hasORCID bool, hasROR bool, match bool) ([]commonmeta.Data, error) {
	var data []commonmeta.Data
	content, err := GetAll(number, page, host, community, subject, type_, year, language, orcid, affiliation, ror, hasORCID, hasROR)
	if err != nil {
		return data, err
	}
	data, err = ReadAll(content, match)
	return data, err
}

// Load loads the metadata for a single work from a JSON file
func Load(filename string, match bool) (commonmeta.Data, error) {
	var data commonmeta.Data

	content, err := ReadJSON(filename)
	if err != nil {
		return data, err
	}
	data, err = Read(content, match)
	if err != nil {
		return data, err
	}
	return data, nil
}

// LoadAll loads a list of Inveniordm metadata from a JSON file and returns Commonmeta metadata.
func LoadAll(filename string, match bool) ([]commonmeta.Data, error) {
	var data []commonmeta.Data
	var content []Content
	var err error

	extension := path.Ext(filename)
	switch extension {
	case ".jsonl", ".jsonlines":
		content, err = ReadJSONLines(filename)
		if err != nil {
			return data, err
		}
	case ".json":
		content, err = ReadJSONList(filename)
		if err != nil {
			return data, err
		}
	default:
		return data, errors.New("unsupported file format")
	}
	data, err = ReadAll(content, match)
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
	if resp.StatusCode != http.StatusOK {
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

// GetAll retrieves InvenioRDM metadata for all records in a community.
func GetAll(number int, page int, host string, community string, subject string, type_ string, year string, language string, orcid string, affiliation string, ror string, hasORCID bool, hasROR bool) ([]Content, error) {
	var response Query
	var content []Content

	if number <= 0 {
		number = 10
	} else if number > 500 {
		number = 500
	}
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	url := QueryURL(number, page, host, community, subject, type_, year, language, orcid, affiliation, ror, hasORCID, hasROR)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return content, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return content, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return content, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return content, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("error:", err)
	}
	content = append(content, response.Hits.Hits...)
	return content, err
}

// QueryURL returns the URL for the InvenioRDM API query
func QueryURL(number int, page int, host string, community string, subject string, type_ string, year string, language string, orcid string, affiliation string, ror string, hasORCID bool, hasROR bool) string {
	var requestURL string
	var q string
	if community != "" {
		requestURL = fmt.Sprintf("https://%s/api/communities/%s/records?", host, community)
	} else {
		requestURL = fmt.Sprintf("https://%s/api/records?", host)
	}
	values := url.Values{}
	if subject != "" {
		if q != "" {
			q += " AND "
		}
		values.Set("q", q+"metadata.subjects.subject:"+subject)
	}
	if type_ != "" {
		if q != "" {
			q += " AND "
		}
		values.Set("q", q+"metadata.resource_type.id:"+type_)
	}
	if year != "" {
		q := values.Get("q")
		if q != "" {
			q += " AND "
		}
		values.Set("q", q+"metadata.publication_date:["+year+"-01-01 TO "+year+"-12-31]")
	}
	if orcid != "" {
		o, _ := utils.ValidateORCID(orcid)
		if o != "" {
			q := values.Get("q")
			if q != "" {
				q += " AND "
			}
			values.Set("q", q+"metadata.creators.person_or_org.identifiers.identifier:"+o)
		}
	}
	if ror != "" {
		r, _ := utils.ValidateROR(ror)
		if r != "" {
			q := values.Get("q")
			if q != "" {
				q += " AND "
			}
			values.Set("q", q+"metadata.creators.affiliations.id:"+r)
		}
	}
	if affiliation != "" {
		q := values.Get("q")
		if q != "" {
			q += " AND "
		}
		values.Set("q", q+"metadata.creators.affiliations.name:\""+affiliation+"\"")
	}
	if hasORCID {
		q := values.Get("q")
		if q != "" {
			q += " AND "
		}
		values.Set("q", q+"metadata.creators.person_or_org.identifiers.scheme:orcid")
	}
	if hasROR {
		q := values.Get("q")
		if q != "" {
			q += " AND "
		}
		values.Set("q", q+"metadata.creators.affiliations.id:*")
	}
	if language != "" {
		q := values.Get("q")
		if q != "" {
			q += " AND "
		}
		l := utils.GetLanguage(language, "iso639-3")
		values.Set("q", q+"metadata.languages.id:"+l)
	}
	values.Add("l", "list")
	values.Add("page", strconv.Itoa(page))
	values.Add("size", strconv.Itoa(number))
	values.Add("sort", "newest")

	return requestURL + values.Encode()
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
		return utils.ParseString(query.Hits.Hits[0].ID), nil
	}
}

// SearchBySlug searches InvenioRDM communities by slug.
// Specify type of community (blog or topic) in query, subject area communities are always queried.
func SearchBySlug(slug string, type_ string, client *InvenioRDMClient, cache *cache2go.CacheTable) (string, error) {
	// first check for cached community ID
	res, _ := cache.Value(slug)
	if res != nil {
		id := fmt.Sprintf("%v", res.Data())
		return id, nil
	}

	var query Query
	requestURL := fmt.Sprintf("https://%s/api/communities?q=slug:%s&type=%s&type=subject", client.Host, slug, type_)
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
		id := utils.ParseString(query.Hits.Hits[0].ID)
		cache.Add(slug, 1*time.Hour, id)
		return id, nil
	}
}

// Read reads InvenioRDM JSON API response and converts it into Commonmeta metadata.
func Read(content Content, match bool) (commonmeta.Data, error) {
	var data commonmeta.Data

	if content.DOI != "" {
		data.ID = doiutils.NormalizeDOI(content.DOI)
	} else {
		data.ID = doiutils.NormalizeDOI(content.Pids.DOI.Identifier)
	}

	if content.Metadata.ResourceType.ID != "" {
		data.Type = InvenioToCMMappings[content.Metadata.ResourceType.ID]
	} else if content.Metadata.ResourceType.Subtype != "" {
		data.Type = InvenioToCMMappings[content.Metadata.ResourceType.Subtype]
	} else {
		data.Type = InvenioToCMMappings[content.Metadata.ResourceType.Type]
	}

	if content.Parent.Communities.Default != "" {
		for _, v := range content.Parent.Communities.Entries {
			if v.ID == content.Parent.Communities.Default {
				var identifier, identifierType string
				if content.CustomFields.Journal.ISSN != "" {
					identifier = content.CustomFields.Journal.ISSN
					identifierType = "ISSN"
				} else {
					identifier = utils.CommunitySlugAsURL(v.Slug, "rogue-scholar.org")
					identifierType = "URL"
				}
				type_ := CommunityTypes[v.Metadata.Type.ID]
				if type_ == "" {
					type_ = "Community"
				}
				data.Container = commonmeta.Container{
					Identifier:     identifier,
					IdentifierType: identifierType,
					Type:           type_,
					Title:          v.Metadata.Title,
				}
				if identifierType == "ISSN" {
					identifier = utils.ISSNAsURL(identifier)
					data.Relations = append(data.Relations, commonmeta.Relation{
						ID:   identifier,
						Type: "IsPartOf",
					})
				}
				identifier = utils.CommunitySlugAsURL(v.Slug, "rogue-scholar.org")
				data.Relations = append(data.Relations, commonmeta.Relation{
					ID:   identifier,
					Type: "IsPartOf",
				})
			} else {
				identifier := utils.CommunitySlugAsURL(v.Slug, "rogue-scholar.org")
				data.Relations = append(data.Relations, commonmeta.Relation{
					ID:   identifier,
					Type: "IsPartOf",
				})
			}
		}
	}

	for _, v := range content.Metadata.Creators {
		var contributor commonmeta.Contributor
		if v.PersonOrOrg.Name != "" || v.PersonOrOrg.FamilyName != "" {
			contributor = GetContributor(v, match)
		} else if v.Name != "" {
			contributor = GetZenodoContributor(v)
		}
		containsID := slices.ContainsFunc(data.Contributors, func(e commonmeta.Contributor) bool {
			return e.ID != "" && e.ID == contributor.ID
		})
		if !containsID {
			data.Contributors = append(data.Contributors, contributor)
		}
	}

	for _, v := range content.Metadata.Dates {
		// parse Date as either string or struct
		var tt Type
		var ts, t string
		err := json.Unmarshal(v.Type, &tt)
		if err != nil {
			err = json.Unmarshal(v.Type, &ts)
		}
		if err != nil {
			log.Println(err)
		}
		if ts != "" {
			t = ts
		} else if tt.ID != "" {
			t = tt.ID
		}

		if t == "accepted" {
			data.Date.Accepted = v.Date
		}
		if t == "available" {
			data.Date.Available = v.Date
		}
		if t == "collected" {
			data.Date.Collected = v.Date
		}
		if t == "created" {
			data.Date.Created = v.Date
		}
		if t == "issued" {
			data.Date.Published = v.Date
		}
		if t == "submitted" {
			data.Date.Submitted = v.Date
		}
		if t == "updated" {
			data.Date.Updated = v.Date
		}
		if t == "valid" {
			data.Date.Valid = v.Date
		}
		if t == "withdrawn" {
			data.Date.Withdrawn = v.Date
		}
		if t == "other" {
			data.Date.Other = v.Date
		}
	}
	if data.Date.Published == "" && content.Metadata.PublicationDate != "" {
		data.Date.Published = content.Metadata.PublicationDate
	}

	if content.Metadata.Description != "" {
		description := utils.Sanitize(content.Metadata.Description)
		data.Descriptions = append(data.Descriptions, commonmeta.Description{
			Description: description,
			Type:        "Abstract",
		})
	}

	if content.CustomFields.FeatureImage != "" {
		data.FeatureImage = content.CustomFields.FeatureImage
	}

	if doiutils.IsRogueScholarDOI(data.ID, "") {
		doi, _ := doiutils.ValidateDOI(data.ID)
		data.Files = append(data.Files, commonmeta.File{
			URL:      fmt.Sprintf("https://api.rogue-scholar.org/posts/%s.md", doi),
			MimeType: "text/markdown",
		})
		data.Files = append(data.Files, commonmeta.File{
			URL:      fmt.Sprintf("https://api.rogue-scholar.org/posts/%s.pdf", doi),
			MimeType: "application/pdf",
		})
		data.Files = append(data.Files, commonmeta.File{
			URL:      fmt.Sprintf("https://api.rogue-scholar.org/posts/%s.epub", doi),
			MimeType: "application/epub+zip",
		})
		data.Files = append(data.Files, commonmeta.File{
			URL:      fmt.Sprintf("https://api.rogue-scholar.org/posts/%s.xml", doi),
			MimeType: "application/xml",
		})
	}

	if len(content.Metadata.Funding) > 0 {
		for _, v := range content.Metadata.Funding {
			funderIdentifier, funderIdentifierType := utils.ValidateID(v.Funder.ID)
			if funderIdentifierType == "ROR" {
				funderIdentifier = utils.NormalizeROR(funderIdentifier)
			}
			awardNumber := v.Award.Number
			awardURI, _ := utils.ValidateID(v.Award.ID)
			data.FundingReferences = append(data.FundingReferences, commonmeta.FundingReference{
				FunderIdentifier:     funderIdentifier,
				FunderIdentifierType: funderIdentifierType,
				FunderName:           v.Funder.Name,
				AwardNumber:          awardNumber,
				AwardURI:             awardURI,
			})
		}
	} else if len(content.Metadata.Grants) > 0 {
		for _, v := range content.Metadata.Grants {
			var funderIdentifierType string
			funderIdentifier := doiutils.NormalizeDOI(v.Funder.DOI)
			if funderIdentifier != "" {
				funderIdentifierType = "Crossref Funder ID"
			}
			awardNumber := v.Code
			awardURI, _ := utils.NormalizeURL(v.URL, true, false)
			data.FundingReferences = append(data.FundingReferences, commonmeta.FundingReference{
				FunderIdentifier:     funderIdentifier,
				FunderIdentifierType: funderIdentifierType,
				FunderName:           v.Funder.Name,
				AwardNumber:          awardNumber,
				AwardTitle:           v.Title,
				AwardURI:             awardURI,
			})
		}
	}

	// GeoLocationPoint can be float64 or string
	// for _, v := range content.GeoLocations {
	// 	pointLongitude := ParseGeoCoordinate(v.GeoLocationPointInterface.PointLongitude)
	// 	pointLatitude := ParseGeoCoordinate(v.GeoLocationPointInterface.PointLatitude)
	// 	westBoundLongitude := ParseGeoCoordinate(v.GeoLocationBoxInterface.WestBoundLongitude)
	// 	eastBoundLongitude := ParseGeoCoordinate(v.GeoLocationBoxInterface.EastBoundLongitude)
	// 	southBoundLatitude := ParseGeoCoordinate(v.GeoLocationBoxInterface.SouthBoundLatitude)
	// 	northBoundLatitude := ParseGeoCoordinate(v.GeoLocationBoxInterface.NorthBoundLatitude)
	// 	geoLocation := commonmeta.GeoLocation{
	// 		GeoLocationPlace: v.GeoLocationPlace,
	// 		GeoLocationPoint: commonmeta.GeoLocationPoint{
	// 			PointLongitude: pointLongitude,
	// 			PointLatitude:  pointLatitude,
	// 		},
	// 		GeoLocationBox: commonmeta.GeoLocationBox{
	// 			WestBoundLongitude: westBoundLongitude,
	// 			EastBoundLongitude: eastBoundLongitude,
	// 			SouthBoundLatitude: southBoundLatitude,
	// 			NorthBoundLatitude: northBoundLatitude,
	// 		},
	// 	}
	// 	data.GeoLocations = append(data.GeoLocations, geoLocation)
	//}

	data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
		Identifier:     data.ID,
		IdentifierType: "DOI",
	})
	if len(content.Metadata.Identifiers) > 0 {
		for _, v := range content.Metadata.Identifiers {
			identifier := v.Identifier
			scheme := InvenioToCMIdentifierMappings[v.Scheme]
			if scheme == "URL" {
				data.URL, _ = utils.NormalizeURL(identifier, true, false)
			} else if scheme != "" {
				data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
					Identifier:     identifier,
					IdentifierType: scheme,
				})
			}
		}
	}
	if content.ID != nil {
		switch v := content.ID.(type) {
		case string:
			data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
				Identifier:     v,
				IdentifierType: "RID",
			})
		}
	}

	if len(content.Metadata.Languages) > 0 {
		data.Language = utils.GetLanguage(content.Metadata.Languages[0].ID, "iso639-1")
	} else if content.Metadata.Language != "" {
		data.Language = utils.GetLanguage(content.Metadata.Language, "iso639-1")
	}

	if len(content.Metadata.Rights) > 0 {
		licenseID := LicenseMappings[content.Metadata.Rights[0].ID]
		licenseURL := content.Metadata.Rights[0].Props.URL
		data.License = commonmeta.License{
			ID:  licenseID,
			URL: licenseURL,
		}
	} else if content.Metadata.License.ID != "" {
		var licenseURL string
		licenseID := LicenseMappings[content.Metadata.License.ID]
		license, _ := spdx.Search(licenseID)
		if len(license.SeeAlso) == 0 {
			licenseURL = license.SeeAlso[0]
		}
		data.License = commonmeta.License{
			ID:  licenseID,
			URL: licenseURL,
		}
	}

	if doiutils.IsRogueScholarDOI(data.ID, "") {
		data.Provider = "Crossref"
	}

	if content.Metadata.Publisher != "" {
		data.Publisher = commonmeta.Publisher{
			Name: content.Metadata.Publisher,
		}
	}
	// workaround until InvenioRDM supports BlogPost type
	if data.Type == "Article" && (data.Publisher.Name == "Front Matter" || doiutils.IsRogueScholarDOI(data.ID, "")) {
		data.Type = "BlogPost"
	}

	if len(content.Metadata.Subjects) > 0 {
		for _, v := range content.Metadata.Subjects {
			s := v.Subject
			// if v.Scheme == "FOS" {
			// 	s = "FOS: " + s
			// }
			subject := commonmeta.Subject{
				Subject: s,
			}
			if !slices.Contains(data.Subjects, subject) {
				data.Subjects = append(data.Subjects, subject)
			}
		}
	} else if len(content.Metadata.Keywords) > 0 {
		for _, v := range content.Metadata.Keywords {
			subject := commonmeta.Subject{
				Subject: v,
			}
			if !slices.Contains(data.Subjects, subject) {
				data.Subjects = append(data.Subjects, subject)
			}
		}
	}

	if len(content.Metadata.References) > 0 {
		for _, v := range content.Metadata.References {
			id := utils.NormalizeID(v.Identifier)
			data.References = append(data.References, commonmeta.Reference{
				ID:           id,
				Unstructured: v.Reference,
			})
		}
	}

	if len(content.Metadata.RelatedIdentifiers) > 0 {
		references := []string{
			"cites",
			"references",
		}
		for _, v := range content.Metadata.RelatedIdentifiers {
			id := utils.NormalizeID(v.Identifier)
			if id != "" && slices.Contains(references, v.RelationType.ID) {
				data.References = append(data.References, commonmeta.Reference{
					ID: id,
				})
			} else if id != "" {
				t := InvenioToCMRelationTypeMappings[v.RelationType.ID]
				if t != "" {
					data.Relations = append(data.Relations, commonmeta.Relation{
						ID:   id,
						Type: t,
					})
				}
			}
		}
	}

	if content.Metadata.Title != "" {
		data.Titles = append(data.Titles, commonmeta.Title{
			Title: content.Metadata.Title,
		})
	}

	data.Version = content.Metadata.Version

	// optional full text content
	if content.CustomFields.ContentHTML != "" {
		data.ContentHTML = content.CustomFields.ContentHTML
	}

	return data, nil
}

// ReadAll reads a list of Inveniordm JSON responses and returns a list of works in Commonmeta format
func ReadAll(content []Content, match bool) ([]commonmeta.Data, error) {
	var data []commonmeta.Data
	for _, v := range content {
		d, err := Read(v, match)
		if err != nil {
			log.Println(err)
		}
		data = append(data, d)
	}
	return data, nil
}

// ReadJSON reads JSON from a file and unmarshals it
func ReadJSON(filename string) (Content, error) {
	var content Content

	extension := path.Ext(filename)
	if extension != ".json" {
		return content, errors.New("invalid file extension")
	}
	file, err := os.Open(filename)
	if err != nil {
		return content, errors.New("error reading file")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&content)
	return content, err
}

// ReadJSONList reads JSON list from a file and unmarshals it
func ReadJSONList(filename string) ([]Content, error) {
	var content []Content

	extension := path.Ext(filename)
	if extension != ".json" {
		return content, errors.New("invalid file extension")
	}
	file, err := os.Open(filename)
	if err != nil {
		return content, errors.New("error reading file")
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&content)
	return content, err
}

// ReadJSONLines reads JSON lines from a file and unmarshals them
func ReadJSONLines(filename string) ([]Content, error) {
	var response []Content

	extension := path.Ext(filename)
	if extension != ".jsonl" && extension != ".jsonlines" {
		return nil, errors.New("invalid file extension")
	}
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("error reading file")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		var inveniordm Content
		if err := decoder.Decode(&inveniordm); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		response = append(response, inveniordm)
	}

	return response, nil
}

// GetContributor converts Inveniordm contributor metadata into the Commonmeta format
func GetContributor(v Creator, match bool) commonmeta.Contributor {
	var t string
	switch v.PersonOrOrg.Type {
	case "personal":
		t = "Person"
	case "organizational":
		t = "Organization"
	}
	var id string
	if len(v.PersonOrOrg.Identifiers) > 0 {
		ni := v.PersonOrOrg.Identifiers[0]
		switch ni.Scheme {
		case "orcid":
			id = utils.NormalizeORCID(ni.Identifier)
			t = "Person"
		case "ROR":
			id = utils.NormalizeROR(ni.Identifier)
			t = "Organization"
		default:
			id = ni.Identifier
		}
	}
	name := v.PersonOrOrg.Name
	givenName := v.PersonOrOrg.GivenName
	familyName := v.PersonOrOrg.FamilyName
	if t == "" && (givenName != "" || familyName != "") {
		t = "Person"
	} else if t == "" {
		t = "Organization"
	}
	if t == "Person" && name != "" && familyName == "" {
		// split name for type Person into given/family name if not already provided
		givenName, familyName, name = authorutils.ParseName(name)
	}

	var affiliations []*commonmeta.Affiliation
	for _, a := range v.Affiliations {
		var assertedBy string
		ID, _ := utils.ValidateROR(a.ID)
		if ID != "" {
			assertedBy = "publisher"
		}
		ID, name, assertedBy, err := ror.MapROR(ID, a.Name, assertedBy, match)
		if err != nil {
			fmt.Println("error mapping ROR:", err)
		}
		if ID != "" || name != "" {
			affiliations = append(affiliations, &commonmeta.Affiliation{
				ID:         ID,
				Name:       name,
				AssertedBy: assertedBy,
			})
		}
	}

	var roles []string
	roles = append(roles, "Author")

	return commonmeta.Contributor{
		ID:               id,
		Type:             t,
		Name:             name,
		GivenName:        givenName,
		FamilyName:       familyName,
		Affiliations:     affiliations,
		ContributorRoles: roles,
	}
}

// GetZenodoContributor converts Zenodo contributor metadata into the Commonmeta format
func GetZenodoContributor(v Creator) commonmeta.Contributor {
	var id, t string

	// split name into given/family name
	givenName, familyName, name := authorutils.ParseName(v.Name)

	if v.ORCID != "" {
		id = utils.NormalizeORCID(v.ORCID)
		t = "Person"
	}

	if t == "" && (givenName != "" || familyName != "") {
		t = "Person"
	} else if t == "" {
		t = "Organization"
	}
	if t == "Person" && familyName == "" && name != "" {
		familyName = name
		name = ""
	}

	var affiliations []*commonmeta.Affiliation
	if v.Affiliation != "" {
		affiliations = append(affiliations, &commonmeta.Affiliation{
			Name: v.Affiliation,
		})
	}

	var roles []string
	roles = append(roles, "Author")

	return commonmeta.Contributor{
		ID:               id,
		Type:             t,
		Name:             name,
		GivenName:        givenName,
		FamilyName:       familyName,
		Affiliations:     affiliations,
		ContributorRoles: roles,
	}
}
