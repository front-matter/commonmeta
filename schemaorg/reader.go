package schemaorg

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/crossref"
	"github.com/front-matter/commonmeta/datacite"
	"github.com/front-matter/commonmeta/dateutils"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/utils"
	"github.com/samber/lo"
)

// SchemaOrg represents the Schema.org metadata.
type SchemaOrg struct {
	Context               string        `json:"@context"`
	ID                    string        `json:"@id"`
	Type                  string        `json:"@type"`
	AdditionalType        string        `json:"additionalType,omitempty"`
	Author                []Contributor `json:"author,omitempty"`
	Citation              []Citation    `json:"citation,omitempty"`
	CodeRepository        string        `json:"codeRepository,omitempty"`
	ContentURL            []string      `json:"contentUrl,omitempty"`
	Contributor           []Contributor `json:"contributor,omitempty"`
	Creator               []Contributor `json:"creator,omitempty"`
	DateCreated           string        `json:"dateCreated,omitempty"`
	DatePublished         string        `json:"datePublished,omitempty"`
	DateModified          string        `json:"dateModified,omitempty"`
	Description           string        `json:"description,omitempty"`
	Distribution          []MediaObject `json:"distribution,omitempty"`
	Editor                []Editor      `json:"editor,omitempty"`
	Encoding              []MediaObject `json:"encoding,omitempty"`
	Headline              string        `json:"headline,omitempty"`
	Identifier            []string      `json:"identifier,omitempty"`
	IncludedInDataCatalog DataCatalog   `json:"includedInDataCatalog,omitempty"`
	InLanguage            string        `json:"inLanguage,omitempty"`
	Keywords              string        `json:"keywords,omitempty"`
	License               string        `json:"license,omitempty"`
	Name                  string        `json:"name,omitempty"`
	PageStart             string        `json:"pageStart,omitempty"`
	PageEnd               string        `json:"pageEnd,omitempty"`
	Periodical            Periodical    `json:"periodical,omitempty"`
	Provider              Provider      `json:"provider,omitempty"`
	Publisher             Publisher     `json:"publisher,omitempty"`
	URL                   string        `json:"url,omitempty"`
	Version               string        `json:"version,omitempty"`
}

// Content represents the SchemaOrg metadata returned from SchemaOrg sources. The type is more
// flexible than the SchemaOrg type, allowing for different formats of some metadata.
// Identifier can be string or []string.
type Content struct {
	SchemaOrg
	Author      json.RawMessage `json:"author,omitempty"`
	Creator     json.RawMessage `json:"creator,omitempty"`
	Contributor json.RawMessage `json:"contributor,omitempty"`
	Editor      json.RawMessage `json:"editor,omitempty"`
	Identifier  json.RawMessage `json:"identifier,omitempty"`
	Keywords    json.RawMessage `json:"keywords,omitempty"`
	Version     interface{}     `json:"version,omitempty"`
}

// Citation represents a citation or reference to another creative work, such as another publication, web page, scholarly article, etc.
type Citation struct {
	ID   string `json:"@id,omitempty"`
	Type string `json:"@type,omitempty"`
	Name string `json:"name,omitempty"`
}

// Coderepository
type CodeRepository struct {
}

// Contributor represents the author, creator or contributor of this CreativeWork.
type Contributor struct {
	ID          string       `json:"@id,omitempty"`
	Type        string       `json:"@type,omitempty"`
	GivenName   string       `json:"givenName,omitempty"`
	FamilyName  string       `json:"familyName"`
	Name        string       `json:"name,omitempty"`
	Affiliation Organization `json:"affiliation,omitempty"`
}

// Datacatalog represents a collection of datasets.
type DataCatalog struct {
	ID   string `json:"@id,omitempty"`
	Type string `json:"@type,omitempty"`
	Name string `json:"name,omitempty"`
}

type Editor struct {
	ID          string       `json:"@id,omitempty"`
	Type        string       `json:"@type,omitempty"`
	GivenName   string       `json:"givenName,omitempty"`
	FamilyName  string       `json:"familyName"`
	Name        string       `json:"name,omitempty"`
	Affiliation Organization `json:"affiliation,omitempty"`
}

// MediaObject represents a media object, such as an image, video, audio, or text object
// embedded in a web page or a downloadable dataset i.e. DataDownload.
type MediaObject struct {
	Type           string `json:"@type"`
	ContentURL     string `json:"contentUrl"`
	EncodingFormat string `json:"encodingFormat,omitempty"`
	Name           string `json:"name,omitempty"`
	SHA256         string `json:"sha256,omitempty"`
	Size           string `json:"size,omitempty"`
}

// Organization represents an organization such as a school, NGO, corporation, club, etc.
type Organization struct {
	ID     string `json:"@id,omitempty"`
	SameAs string `json:"sameAs,omitempty"`
	Name   string `json:"name"`
}

// Periodical represents a publication in any medium issued in successive parts bearing numerical or chronological designations and intended to continue indefinitely, such as a magazine, scholarly journal, or newspaper.
type Periodical struct {
	ID   string `json:"@id,omitempty"`
	Type string `json:"@type"`
	Name string `json:"name,omitempty"`
	ISSN string `json:"issn,omitempty"`
}

// Person represents a person (alive, dead, undead, or fictional).
type Person struct {
	ID         string `json:"@id,omitempty"`
	GivenName  string `json:"givenName,omitempty"`
	FamilyName string `json:"familyName"`
}

// Provider represents the provider of the metadata.
type Provider struct {
	Type string `json:"@type"`
	Name string `json:"name"`
}

// Publisher represents the publisher of the metadata.
type Publisher struct {
	Type string `json:"@type"`
	Name string `json:"name"`
}

// SOToCMMappings maps Schema.org types to Commonmeta types.
var SOToCMMappings = map[string]string{
	"Article":            "Article",
	"BlogPosting":        "Article",
	"Book":               "Book",
	"BookChapter":        "BookChapter",
	"CreativeWork":       "Other",
	"Dataset":            "Dataset",
	"DigitalDocument":    "Document",
	"Dissertation":       "Dissertation",
	"Instrument":         "Instrument",
	"NewsArticle":        "Article",
	"Legislation":        "LegalDocument",
	"Report":             "Report",
	"ScholarlyArticle":   "JournalArticle",
	"SoftwareSourceCode": "Software",
}

// Fetch fetches Schemaorg metadata for a given URL and returns Commonmeta metadata.
func Fetch(url string) (commonmeta.Data, error) {
	var data commonmeta.Data

	content, err := Get(url)
	if err != nil {
		return data, err
	}
	// if url represents (Crossref or DataCite) DOI, fetch metadata from Crossref or DataCite API
	if content.Provider.Name == "Crossref" {
		data, err = crossref.Fetch(content.ID)
	} else if content.Provider.Name == "DataCite" {
		data, err = datacite.Fetch(content.ID)
	} else {
		data, err = Read(content)
	}
	return data, err
}

// Get gets Schemaorg metadata for a given URL
func Get(url string) (Content, error) {
	var content Content
	var err error

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(url)
	if err != nil {
		return content, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return content, fmt.Errorf("HTTP status code: %d", resp.StatusCode)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return content, err
	}

	// Find the Schema.org JSON-LD script
	jsonLD := doc.FindMatcher(goquery.Single("script[type='application/ld+json']"))
	if len(jsonLD.Nodes) > 0 {
		json.Unmarshal([]byte(jsonLD.Text()), &content)
	}

	// Find ID
	if content.ID == "" {
		var ids []*goquery.Selection
		ids = append(ids, doc.FindMatcher(goquery.Single("meta[name='citation_doi']")))
		ids = append(ids, doc.FindMatcher(goquery.Single("meta[name='dc.identifier']")))
		ids = append(ids, doc.FindMatcher(goquery.Single("meta[name='DC.identifier']")))
		ids = append(ids, doc.FindMatcher(goquery.Single("meta[name='bepress_citation_doi']")))
		ids = lo.Compact(ids)
		if len(ids) > 0 {
			content.ID, _ = ids[0].Attr("content")
		}
	}

	// if id represents a DOI, get metadata from Crossref or DataCite
	doi, ok := doiutils.ValidateDOI(content.ID)
	if ok {
		ra, ok := doiutils.GetDOIRA(doi)
		if ok {
			if ra == "Crossref" {
				content.Provider = Provider{
					Type: "Organization",
					Name: "Crossref",
				}
				return content, nil
			} else if ra == "DataCite" {
				content.Provider = Provider{
					Type: "Organization",
					Name: "DataCite",
				}
				return content, nil
			}
		}
	}

	if content.Type == "" {
		var types []*goquery.Selection
		types = append(types, doc.FindMatcher(goquery.Single("meta[property='og:type']")))
		types = append(types, doc.FindMatcher(goquery.Single("meta[name='dc.type']")))
		types = append(types, doc.FindMatcher(goquery.Single("meta[name='DC.type']")))
		types = lo.Compact(types)
		if len(types) > 0 {
			content.Type, _ = types[0].Attr("content")
		}
	}

	// Find name
	if content.Name == "" {
		var names []*goquery.Selection
		names = append(names, doc.FindMatcher(goquery.Single("meta[name='citation_title']")))
		names = append(names, doc.FindMatcher(goquery.Single("meta[name='dc.title']")))
		names = append(names, doc.FindMatcher(goquery.Single("meta[name='DC.title']")))
		names = append(names, doc.FindMatcher(goquery.Single("meta[property='og:title']")))
		names = append(names, doc.FindMatcher(goquery.Single("meta[name='twitter:title']")))
		if len(names) > 0 {
			content.Name = lo.Compact(names)[0].Text()
		}
	}

	// Find description// 	// Find description
	if len(content.Description) == 0 {
		var descriptions []*goquery.Selection
		descriptions = append(descriptions, doc.FindMatcher(goquery.Single("meta[name='citation_abstract']")))
		descriptions = append(descriptions, doc.FindMatcher(goquery.Single("meta[name='dc.description']")))
		descriptions = append(descriptions, doc.FindMatcher(goquery.Single("meta[property='og:description']")))
		descriptions = append(descriptions, doc.FindMatcher(goquery.Single("meta[name='twitter:description']")))
		if len(descriptions) == 0 {
			content.Description = lo.Compact(descriptions)[0].Text()
		}
	}

	// Find date published
	if content.DatePublished == "" {
		var datePublished []*goquery.Selection
		datePublished = append(datePublished, doc.FindMatcher(goquery.Single("meta[name='citation_publication_date']")))
		datePublished = append(datePublished, doc.FindMatcher(goquery.Single("meta[name='citation_date']")))
		datePublished = append(datePublished, doc.FindMatcher(goquery.Single("meta[name='dc.date']")))
		datePublished = append(datePublished, doc.FindMatcher(goquery.Single("meta[property='article:published_time']")))
		datePublished = lo.Compact(datePublished)
		if len(datePublished) == 0 {
			dp, _ := datePublished[0].Attr("content")
			content.DatePublished = dateutils.StripMilliseconds(dp)
		}
	}

	// Find date modified
	if content.DateModified == "" {
		var dateModified []*goquery.Selection
		dateModified = append(dateModified, doc.FindMatcher(goquery.Single("meta[name='og:updated_time']")))
		dateModified = append(dateModified, doc.FindMatcher(goquery.Single("meta[name='article:modified_time']")))
		dateModified = lo.Compact(dateModified)
		if len(dateModified) == 0 {
			dm, _ := dateModified[0].Attr("content")
			content.DateModified = dateutils.StripMilliseconds(dm)
		}
	}

	// 	Find language
	if content.InLanguage == "" {
		content.InLanguage = doc.FindMatcher(goquery.Single("html")).AttrOr("lang", "")
	}

	// 	Find license
	if content.License == "" {
		license := doc.FindMatcher(goquery.Single("link[rel='license']"))
		if len(license.Nodes) > 0 {
			content.License = license.AttrOr("href", "")
		}
	}

	// author and creator are synonyms
	if len(content.Author) == 0 && len(content.Creator) > 0 {
		content.Author = content.Creator
	}
	return content, err
}

// Load loads the metadata for a single work from a JSON file
func Load(filename string) (commonmeta.Data, error) {
	var data commonmeta.Data
	var content Content

	extension := path.Ext(filename)
	if extension != ".json" {
		return data, errors.New("invalid file extension")
	}
	file, err := os.Open(filename)
	if err != nil {
		return data, errors.New("error reading file")
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
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

// Read reads Schema.org metadata and converts it to commonmeta.
func Read(content Content) (commonmeta.Data, error) {
	var data commonmeta.Data

	data.ID = utils.NormalizeID(content.ID)
	data.Type = SOToCMMappings[content.Type]
	if data.Type == "" {
		data.Type = "WebPage"
	}
	data.AdditionalType = content.AdditionalType

	var contributor Contributor
	var contributors []Contributor
	err := json.Unmarshal(content.Author, &contributor)
	if err != nil {
		_ = json.Unmarshal(content.Author, &contributors)
	}
	if len(contributors) == 0 {
		contributors = append(contributors, contributor)
	}
	for _, v := range contributors {
		if v.Name != "" || v.GivenName != "" || v.FamilyName != "" {
			contributor := GetContributor(v)
			containsID := slices.ContainsFunc(data.Contributors, func(e commonmeta.Contributor) bool {
				return e.ID != "" && e.ID == contributor.ID
			})
			if !containsID {
				data.Contributors = append(data.Contributors, contributor)
			}
		}
	}

	if content.DatePublished != "" {
		data.Date.Published = content.DatePublished
	}
	if content.DateModified != "" {
		data.Date.Updated = content.DateModified
	}
	if content.DateCreated != "" {
		data.Date.Created = content.DateCreated
	}

	if len(content.Description) > 0 {
		data.Descriptions = append(data.Descriptions, commonmeta.Description{
			Description: utils.Sanitize(content.Description),
			Type:        "Abstract",
		})
	}

	var identifier string
	var identifiers []string
	err = json.Unmarshal(content.Identifier, &identifier)
	if err != nil {
		_ = json.Unmarshal(content.Identifier, &identifiers)
	}
	if identifier != "" {
		identifiers = append(identifiers, identifier)
	}
	if len(identifiers) > 0 {
		for _, id := range identifiers {
			if id != data.ID {
				identifier, identifierType := utils.ValidateID(id)
				if identifierType == "DOI" {
					identifier = doiutils.NormalizeDOI(identifier)
				}
				data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
					Identifier:     identifier,
					IdentifierType: identifierType,
				})
			}
		}
	}

	data.Language = content.InLanguage

	if content.License != "" {
		licenseURL, ok := utils.NormalizeCCUrl(content.License)
		if ok {
			licenseID := utils.URLToSPDX(licenseURL)
			data.License = commonmeta.License{
				ID:  licenseID,
				URL: licenseURL,
			}
		}
	}

	doi, ok := doiutils.ValidateDOI(data.ID)
	if ok {
		data.Provider, _ = doiutils.GetDOIRA(doi)
	}

	if content.Publisher.Name != "" {
		data.Publisher = commonmeta.Publisher{
			Name: content.Publisher.Name,
		}
	}

	var keyword string
	var keywords []string
	err = json.Unmarshal(content.Keywords, &keyword)
	if err != nil {
		_ = json.Unmarshal(content.Keywords, &keywords)
	}
	if keyword != "" {
		keywords = strings.Split(keyword, ",")
	}
	if len(keywords) > 0 {
		for _, subject := range keywords {
			data.Subjects = append(data.Subjects, commonmeta.Subject{
				Subject: subject,
			})
		}
	}

	if content.Name != "" {
		data.Titles = append(data.Titles, commonmeta.Title{
			Title: content.Name,
		})
	} else if content.Headline != "" {
		data.Titles = append(data.Titles, commonmeta.Title{
			Title: content.Headline,
		})
	}

	data.URL, _ = utils.NormalizeURL(content.URL, true, false)

	// version can be a string or a number
	if content.Version != nil {
		switch v := content.Version.(type) {
		case string:
			data.Version = v
		case float64:
			data.Version = fmt.Sprintf("%v", v)
		}
	}

	return data, nil
}

// GetContributor converts Schemaorg contributor metadata into the Commonmeta format
func GetContributor(v Contributor) commonmeta.Contributor {
	var t string
	if v.Type != "" {
		t = v.Type
	}
	var id string
	if v.ID != "" {
		id = utils.NormalizeORCID(v.ID)
		if id != "" {
			t = "Person"
		} else {
			id = utils.NormalizeROR(v.ID)
			t = "Organization"
		}
	}
	name := v.Name
	givenName := v.GivenName
	familyName := v.FamilyName
	if t == "" && (v.GivenName != "" || v.FamilyName != "") {
		t = "Person"
	} else if t == "" {
		t = "Organization"
	}
	if t == "Person" && name != "" {
		// split name for type Person into given/family name if not already provided
		names := strings.Split(name, ",")
		l := len(names)
		if l == 2 {
			givenName = strings.TrimSpace(names[1])
			familyName = names[0]
		} else if l == 1 {
			names = strings.Split(name, " ")
			l = len(names)
			givenName = names[0]
			if l > 1 {
				familyName = strings.Join(names[1:l], " ")
			}
		}
		name = ""
	}

	var affiliations []*commonmeta.Affiliation
	if v.Affiliation.Name != "" {
		id := utils.NormalizeROR(v.Affiliation.ID)
		if id == "" {
			id = utils.NormalizeROR(v.Affiliation.SameAs)
		}
		af := commonmeta.Affiliation{
			ID:   id,
			Name: v.Affiliation.Name,
		}
		affiliations = append(affiliations, &af)
	}

	var roles []string
	if slices.Contains(commonmeta.ContributorRoles, v.Type) {
		roles = append(roles, v.Type)
	} else {
		roles = append(roles, "Author")
	}

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
