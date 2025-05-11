// Package csl converts citation-style language (CSL) metadata to/from the commonmeta metadata format.
package csl

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/dateutils"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/utils"
)

// CSL represents the CSL metadata.
type CSL struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Abstract string `json:"abstract,omitempty"`
	Accessed struct {
		DateAsParts []dateutils.DateSlice `json:"date-parts"`
		DateTime    string                `json:"date-time"`
	} `json:"accessed"`
	Author              []Author `json:"author,omitempty"`
	Categories          []string `json:"categories,omitempty"`
	ContainerTitle      string   `json:"container-title,omitempty"`
	ContainerTitleShort string   `json:"container-title-short,omitempty"`
	DOI                 string   `json:"DOI,omitempty"`
	Editor              []Author `json:"editor,omitempty"`
	ISSN                string   `json:"ISSN,omitempty"`
	Issue               string   `json:"issue,omitempty"`
	Issued              struct {
		DateAsParts []dateutils.DateSlice `json:"date-parts"`
		DateTime    string                `json:"date-time"`
	} `json:"issued"`
	Keyword   string `json:"keyword,omitempty"`
	Language  string `json:"language,omitempty"`
	License   string `json:"license,omitempty"`
	Note      string `json:"note,omitempty"`
	Page      string `json:"page,omitempty"`
	PMID      string `json:"PMID,omitempty"`
	Publisher string `json:"publisher,omitempty"`
	Source    string `json:"source,omitempty"`
	Submitted struct {
		DateAsParts []dateutils.DateSlice `json:"date-parts"`
		DateTime    string                `json:"date-time"`
	} `json:"submitted"`
	Title   string `json:"title,omitempty"`
	URL     string `json:"URL,omitempty"`
	Version string `json:"version,omitempty"`
	Volume  string `json:"volume,omitempty"`
}

// Content represents the CSL metadata. The type is more flexible than the CSL type,
// allowing for different formats of some metadata. Subjects can be a string named keyword,
// or an array of strings named categories. Date parts can be int or string.
// Publisher can be string or struct.

type Content struct {
	*CSL
	Publisher json.RawMessage `json:"publisher"`
}

// Authors represents the author in the CSL item.
type Author struct {
	Family              string `json:"family,omitempty"`
	Given               string `json:"given,omitempty"`
	NonDroppingParticle string `json:"non-dropping-particle,omitempty"`
	Literal             string `json:"literal,omitempty"`
}

type Publisher struct {
	Name string `json:"name,omitempty"`
}

// source: https://docs.citationstyles.org/en/stable/specification.html?highlight=book#appendix-iii-types
var CSLToCMMappings = map[string]string{
	"article":                "Article",
	"article-journal":        "JournalArticle",
	"article-magazine":       "Article",
	"article-newspaper":      "Article",
	"bill":                   "LegalDocument",
	"book":                   "Book",
	"broadcast":              "Audiovisual",
	"chapter":                "BookChapter",
	"classic":                "Book",
	"collection":             "Collection",
	"dataset":                "Dataset",
	"document":               "Document",
	"entry":                  "Entry",
	"entry-dictionary":       "Entry",
	"entry-encyclopedia":     "Entry",
	"event":                  "Event",
	"figure":                 "Figure",
	"graphic":                "Image",
	"hearing":                "LegalDocument",
	"interview":              "Document",
	"legal_case":             "LegalDocument",
	"legislation":            "LegalDocument",
	"manuscript":             "Manuscript",
	"map":                    "Map",
	"motion_picture":         "Audiovisual",
	"musical_score":          "Document",
	"pamphlet":               "Document",
	"paper-conference":       "ProceedingsArticle",
	"patent":                 "Patent",
	"performance":            "Performance",
	"periodical":             "Journal",
	"personal_communication": "PersonalCommunication",
	"post":                   "Post",
	"post-weblog":            "BlogPost",
	"regulation":             "LegalDocument",
	"report":                 "Report",
	"review":                 "Review",
	"review-book":            "Review",
	"software":               "Software",
	"song":                   "Audiovisual",
	"speech":                 "Presentation",
	"standard":               "Standard",
	"thesis":                 "Dissertation",
	"treaty":                 "LegalDocument",
	"webpage":                "WebPage",
}

// Load loads the metadata for a single work from a CSL file
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

// LoadAll loads the metadata for a list of works from a CSL file and converts it to the Commonmeta format
func LoadAll(filename string) ([]commonmeta.Data, error) {
	var data []commonmeta.Data
	var content []Content
	var err error

	extension := path.Ext(filename)
	if extension == ".json" {
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
	} else {
		return data, errors.New("unsupported file format")
	}

	data, err = ReadAll(content)
	if err != nil {
		return data, err
	}
	return data, nil
}

// Read reads CSL metadata and converts it into Commonmeta metadata.
func Read(content Content) (commonmeta.Data, error) {
	var data commonmeta.Data

	if content.DOI != "" {
		data.ID = doiutils.NormalizeDOI(content.DOI)
	} else if strings.HasPrefix(content.Note, "DOI: ") {
		doi := strings.TrimPrefix(content.Note, "DOI: ")
		data.ID = doiutils.NormalizeDOI(doi)
	} else if content.URL != "" {
		data.ID = content.URL
	}
	data.Type = CSLToCMMappings[content.Type]
	if data.Type == "" {
		data.Type = "Other"
	}

	var identifier, identifierType string
	if len(content.ISSN) >= 9 {
		identifier = content.ISSN[:9] // remove extra info, e.g. (Electronic)
		identifierType = "ISSN"
		data.Relations = append(data.Relations, commonmeta.Relation{
			ID:   utils.ISSNAsURL(identifier),
			Type: "IsPartOf",
		})
	}

	var firstPage, lastPage string
	if content.Page != "" {
		pages := strings.Split(content.Page, "-")
		firstPage = pages[0]
		if len(pages) > 1 && pages[1] > firstPage {
			lastPage = pages[1]
		}
	}

	data.Container = commonmeta.Container{
		Type:           "Periodical",
		Title:          content.ContainerTitle,
		Identifier:     identifier,
		IdentifierType: identifierType,
		Volume:         content.Volume,
		Issue:          content.Issue,
		FirstPage:      firstPage,
		LastPage:       lastPage,
	}

	if len(content.Author) > 0 {
		for _, author := range content.Author {
			t := "Person"
			if author.Literal != "" {
				t = "Organization"
			}
			author := commonmeta.Contributor{
				Type:             t,
				GivenName:        author.Given,
				FamilyName:       author.Family,
				Name:             author.Literal,
				ContributorRoles: []string{"Author"},
			}
			data.Contributors = append(data.Contributors, author)
		}
	}
	if len(content.Editor) > 0 {
		for _, editor := range content.Editor {
			t := "Person"
			if editor.Literal != "" {
				t = "Organization"
			}
			editor := commonmeta.Contributor{
				Type:             t,
				GivenName:        editor.Given,
				FamilyName:       editor.Family,
				Name:             editor.Literal,
				ContributorRoles: []string{"Editor"},
			}
			data.Contributors = append(data.Contributors, editor)
		}
	}

	// parse date parts, which can be int or string
	if len(content.Issued.DateAsParts) > 0 {
		data.Date.Published = dateutils.GetDateFromDateParts(content.Issued.DateAsParts)
	}
	if len(content.Submitted.DateAsParts) > 0 {
		data.Date.Submitted = dateutils.GetDateFromDateParts(content.Submitted.DateAsParts)
	}
	if len(content.Accessed.DateAsParts) > 0 {
		data.Date.Accessed = dateutils.GetDateFromDateParts(content.Accessed.DateAsParts)
	}

	if content.Abstract != "" {
		data.Descriptions = []commonmeta.Description{
			{Description: utils.Sanitize(content.Abstract), Type: "Abstract"},
		}
	}

	if content.ID != "" {
		id, identifierType := utils.ValidateID(content.ID)
		if id == "" {
			id = content.ID
		}
		if identifierType == "" {
			identifierType = "Other"
		}
		if identifierType != "DOI" {
			data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
				Identifier:     id,
				IdentifierType: identifierType,
			})
		}
	}

	if content.Language != "" {
		data.Language = utils.GetLanguage(content.Language, "iso639-1")
	}
	if content.License != "" {
		licenseURL, err := utils.NormalizeURL(content.License, true, true)
		if err != nil {
			return data, err
		}
		licenseID := utils.URLToSPDX(licenseURL)
		data.License = &commonmeta.License{
			ID:  licenseID,
			URL: licenseURL,
		}
	}

	// parse Publisher as either string or struct
	if len(content.Publisher) != 0 {
		var publisher Publisher
		var publisherName string
		var err error
		err = json.Unmarshal(content.Publisher, &publisher)
		if err != nil {
			err = json.Unmarshal(content.Publisher, &publisherName)
		}
		if err != nil {
			fmt.Println(2, err)
		}
		if publisher.Name != "" {
			data.Publisher = &commonmeta.Publisher{
				Name: publisher.Name,
			}
		} else if publisherName != "" {
			data.Publisher = &commonmeta.Publisher{
				Name: publisherName,
			}
		}
	}

	if content.Keyword != "" {
		keywords := strings.Split(content.Keyword, ",")
		for _, subject := range keywords {
			data.Subjects = []commonmeta.Subject{
				{Subject: subject},
			}
		}
	} else if len(content.Categories) > 0 {
		for _, category := range content.Categories {
			data.Subjects = append(data.Subjects, commonmeta.Subject{
				Subject: category,
			})
		}
	}

	data.Titles = []commonmeta.Title{
		{Title: utils.Sanitize(content.Title)},
	}

	url, err := utils.NormalizeURL(content.URL, true, false)
	if err != nil {
		return data, err
	}
	data.URL = url

	data.Version = content.Version

	return data, nil
}

// ReadAll reads a list of CSL metadata and returns a list of works in Commonmeta format
func ReadAll(content []Content) ([]commonmeta.Data, error) {
	var data []commonmeta.Data
	for _, v := range content {
		d, err := Read(v)
		if err != nil {
			fmt.Println(v.ID, err)
		}
		data = append(data, d)
	}
	return data, nil
}
