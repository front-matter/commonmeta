// Package jsonfeed converts JSON Feed metadata to/from the commonmeta metadata format.
package jsonfeed

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
	"strings"
	"time"

	"github.com/front-matter/commonmeta/authorutils"
	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/dateutils"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/utils"
)

// Query represents the InvenioRDM JSON API query.
type Query struct {
	Items        []Content `json:"items"`
	TotalResults int       `json:"total-results"`
}

// Content represents the JSON Feed metadata.
type Content struct {
	ID                string             `json:"id"`
	DOI               string             `json:"doi"`
	GUID              string             `json:"guid"`
	RID               string             `json:"rid"`
	Abstract          string             `json:"abstract"`
	ArchiveURL        string             `json:"archive_url"`
	Authors           Authors            `json:"authors"`
	Blog              Blog               `json:"blog"`
	BlogName          string             `json:"blog_name"`
	BlogSlug          string             `json:"blog_slug"`
	ContentText       string             `json:"content_text"`
	FeatureImage      string             `json:"image"`
	IndexedAt         int64              `json:"indexed_at"`
	Language          string             `json:"language"`
	PublishedAt       int64              `json:"published_at"`
	Relationships     []Relation         `json:"relationships"`
	Reference         []Reference        `json:"reference"`
	FundingReferences []FundingReference `json:"funding_references"`
	Summary           string             `json:"summary"`
	Tags              []string           `json:"tags"`
	Title             string             `json:"title"`
	UpdatedAt         int64              `json:"updated_at"`
	URL               string             `json:"url"`
}

// Affiliation represents an affiliation in the JSON Feed item.
type Affiliation struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Authors represents the authors in the JSON Feed item.
type Authors []struct {
	Name        string        `json:"name"`
	URL         string        `json:"url"`
	Affiliation []Affiliation `json:"affiliation"`
}

type Blog struct {
	ID          string           `json:"id"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Favicon     string           `json:"favicon"`
	Funding     FundingReference `json:"funding"`
	Generator   string           `json:"generator"`
	HomePageURL string           `json:"home_page_url"`
	ISSN        string           `json:"issn"`
	Language    string           `json:"language"`
	License     string           `json:"license"`
	Prefix      string           `json:"prefix"`
	Slug        string           `json:"slug"`
	Status      string           `json:"status"`
	Title       string           `json:"title"`
}

// FundingReference represents the funding reference of a publication, defined in the commonmeta JSON Schema.
type FundingReference struct {
	FunderIdentifier     string `json:"funderIdentifier,omitempty"`
	FunderIdentifierType string `json:"funderIdentifierType,omitempty"`
	FunderName           string `json:"funderName,omitempty"`
	AwardNumber          string `json:"awardNumber,omitempty"`
	AwardTitle           string `json:"awardTitle,omitempty"`
	AwardURI             string `json:"awardUri,omitempty"`
}

// Relation represents a relation in the JSON Feed item.
type Relation struct {
	Type string   `json:"type"`
	Urls []string `json:"urls"`
}

// Reference represents a reference in the JSON Feed item.
type Reference struct {
	Key             string `json:"key,omitempty"`
	ID              string `json:"id,omitempty"`
	PublicationYear string `json:"publicationYear,omitempty"`
	Title           string `json:"title,omitempty"`
	Unstructured    string `json:"unstructured,omitempty"`
}

// relation types to include
var relationTypes = []string{"IsPartOf", "HasPart", "IsVariantFormOf", "IsOriginalFormOf", "IsIdenticalTo", "IsTranslationOf", "IsReviewedBy", "Reviews", "HasReview", "IsPreprintOf", "HasPreprint", "IsSupplementTo", "IsSupplementedBy"}

// Fetch fetches JSON Feed metadata and returns Commonmeta metadata.
func Fetch(str string) (commonmeta.Data, error) {
	var data commonmeta.Data
	var id string

	_, identifierType := utils.ValidateID(str)
	if identifierType == "JSONFEEDID" {
		id = str
	} else if identifierType == "DOI" {
		doi, _ := doiutils.ValidateDOI(str)
		id = "https://api.rogue-scholar.org/posts/" + doi
	} else if identifierType == "UUID" {
		id = "https://api.rogue-scholar.org/posts/" + str
	} else {
		return data, errors.New("invalid ID")
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

// FetchAll fetches a list of JSON Feed metadata and returns Commonmeta metadata.
func FetchAll(number int, page int, community string, archived bool) ([]commonmeta.Data, error) {
	var data []commonmeta.Data
	content, err := GetAll(number, page, community, archived)
	if err != nil {
		return data, err
	}
	data, err = ReadAll(content)
	return data, err
}

// Get retrieves JSON Feed metadata.
func Get(id string) (Content, error) {
	var content Content
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(id)
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

// GetAll retrieves JSON Feed metadata for all records in a community.
func GetAll(number int, page int, community string, archived bool) ([]Content, error) {
	var response Query
	var content []Content

	client := &http.Client{
		Timeout: time.Second * 30,
	}
	url := QueryURL(number, page, community, archived)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return content, err
	}
	resp, err := client.Do(req)
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
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("error:", err)
	}
	content = append(content, response.Items...)
	return content, err
}

// QueryURL returns the URL for the Rogue Scholar API query
func QueryURL(number int, page int, community string, archived bool) string {
	requestURL := "https://api.rogue-scholar.org/posts?"
	values := url.Values{}
	if !archived {
		values.Set("flag", "needs_update")
	}
	if community != "" {
		values.Set("blog_slug", community)
	}
	values.Add("per_page", strconv.Itoa(number))
	values.Add("page", strconv.Itoa(page))

	return requestURL + values.Encode()
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

// LoadAll loads the metadata for a list of works from a JSON file and converts it to the Commonmeta format
func LoadAll(filename string) ([]commonmeta.Data, error) {
	var data []commonmeta.Data
	var content []Content
	var err error

	extension := path.Ext(filename)
	if extension == ".json" {
		type Response struct {
			Items []Content `json:"items"`
		}
		var response Response

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
		err = decoder.Decode(&response)
		if err != nil {
			return data, err
		}
		content = response.Items
	} else {
		return data, errors.New("unsupported file format")
	}

	data, err = ReadAll(content)
	if err != nil {
		return data, err
	}
	return data, nil
}

// Read reads JSON Feed metadata and converts it into Commonmeta metadata.
func Read(content Content) (commonmeta.Data, error) {
	var data commonmeta.Data

	if content.DOI != "" {
		data.ID = doiutils.NormalizeDOI(content.DOI)
	} else if content.GUID != "" {
		// use GUID as DOI string if prefix matches blog prefix
		// and suffix is base32 encoded number with checksum
		doi := doiutils.NormalizeDOI(content.GUID)
		prefix, ok := doiutils.ValidatePrefix(doi)
		number, err := utils.DecodeID(doi)
		if ok && prefix == content.Blog.Prefix && number != 0 && err == nil {
			data.ID = doi
		}
	}
	if data.ID == "" && content.Blog.Prefix != "" {
		// optionally generate a DOI string if missing but a DOI prefix is provided
		data.ID = doiutils.EncodeDOI(content.Blog.Prefix)
	} else if data.ID == "" {
		data.ID = content.URL
	}
	data.Type = "BlogPost"

	identifier := content.Blog.HomePageURL
	identifierType := "URL"
	if content.Blog.ISSN != "" {
		identifier = content.Blog.ISSN
		identifierType = "ISSN"
		data.Relations = append(data.Relations, commonmeta.Relation{
			ID:   utils.ISSNAsURL(identifier),
			Type: "IsPartOf",
		})
	}
	if content.Blog.Slug != "" {
		data.Relations = append(data.Relations, commonmeta.Relation{
			ID:   utils.CommunitySlugAsURL(content.Blog.Slug, "rogue-scholar.org"),
			Type: "IsPartOf",
		})
	}
	data.Container = commonmeta.Container{
		Type:           "Blog",
		Title:          content.Blog.Title,
		Description:    content.Blog.Description,
		Language:       content.Blog.Language,
		License:        commonmeta.License{URL: content.Blog.License},
		Favicon:        content.Blog.Favicon,
		Platform:       content.Blog.Generator,
		Identifier:     identifier,
		IdentifierType: identifierType,
	}

	if len(content.Authors) > 0 {
		contrib, err := GetContributors(content.Authors)
		if err != nil {
			return data, err
		}
		data.Contributors = append(data.Contributors, contrib...)
	}

	data.Date.Published = dateutils.GetDateTimeFromUnixTimestamp(content.PublishedAt)
	data.Date.Updated = dateutils.GetDateTimeFromUnixTimestamp(content.UpdatedAt)

	description := content.Summary
	if content.Abstract != "" {
		description = content.Abstract
	}
	data.Descriptions = []commonmeta.Description{
		{Description: utils.Sanitize(description), Type: "Abstract"},
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

		data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
			Identifier:     data.ID,
			IdentifierType: "DOI",
		})

		data.Provider = "Crossref"
	}

	data.FundingReferences = GetFundingReferences(content)

	data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
		Identifier:     content.ID,
		IdentifierType: "UUID",
	})
	if content.GUID != "" {
		data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
			Identifier:     content.GUID,
			IdentifierType: "GUID",
		})
	}
	if content.RID != "" {
		data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
			Identifier:     content.RID,
			IdentifierType: "RID",
		})
	}

	data.Language = content.Language

	licenseURL, err := utils.NormalizeURL(content.Blog.License, true, true)
	if err != nil {
		return data, err
	}
	licenseID := utils.URLToSPDX(licenseURL)
	data.License = commonmeta.License{
		ID:  licenseID,
		URL: licenseURL,
	}

	data.Publisher = commonmeta.Publisher{
		Name: "Front Matter",
	}
	for _, v := range content.Relationships {
		if slices.Contains(relationTypes, v.Type) {
			for _, u := range v.Urls {
				url, err := utils.NormalizeURL(u, true, true)
				if err != nil {
					return data, err
				}
				data.Relations = append(data.Relations, commonmeta.Relation{
					ID:   url,
					Type: v.Type,
				})
			}
		}
	}

	for _, v := range content.Reference {
		reference := commonmeta.Reference{
			Key:             v.Key,
			ID:              v.ID,
			Title:           v.Title,
			PublicationYear: v.PublicationYear,
			Unstructured:    v.Unstructured,
		}
		containsKey := slices.ContainsFunc(data.References, func(e commonmeta.Reference) bool {
			return e.Key != "" && e.Key == reference.Key
		})
		containsID := slices.ContainsFunc(data.References, func(e commonmeta.Reference) bool {
			return e.ID != "" && e.ID == reference.ID
		})
		if !containsKey && !containsID {
			data.References = append(data.References, reference)
		}
	}

	if content.Blog.Category != "" {
		subject := commonmeta.FOSKeyMappings[content.Blog.Category]
		data.Subjects = append(data.Subjects, commonmeta.Subject{
			Subject: subject,
		})
		data.Relations = append(data.Relations, commonmeta.Relation{
			ID:   utils.CommunitySlugAsURL(content.Blog.Category, "rogue-scholar.org"),
			Type: "IsPartOf",
		})
	}

	for _, v := range content.Tags {
		data.Subjects = append(data.Subjects, commonmeta.Subject{
			Subject: v,
		})
	}

	data.Titles = []commonmeta.Title{
		{Title: utils.Sanitize(content.Title)},
	}

	url, err := utils.NormalizeURL(content.URL, true, false)
	if content.Blog.Status == "archived" && content.ArchiveURL != "" {
		url, err = utils.NormalizeURL(content.ArchiveURL, true, false)
	}
	if err != nil {
		return data, err
	}

	data.URL = url
	data.ContentText = content.ContentText
	data.FeatureImage = content.FeatureImage

	return data, nil
}

// ReadAll reads a list of JSON Feed responses and returns a list of works in Commonmeta format
func ReadAll(content []Content) ([]commonmeta.Data, error) {
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

func GetContributors(contrib Authors) ([]commonmeta.Contributor, error) {
	var contributors []commonmeta.Contributor

	if len(contrib) > 0 {
		for _, v := range contrib {
			ID := utils.NormalizeORCID(v.URL)
			GivenName, FamilyName, Name := authorutils.ParseName(v.Name)
			var Type string
			if Name == "" {
				Type = "Person"
			} else {
				Type = "Organization"
			}

			var affiliations []*commonmeta.Affiliation
			if len(v.Affiliation) > 0 {
				for _, a := range v.Affiliation {
					if a.Name != "" {
						affiliations = append(affiliations, &commonmeta.Affiliation{
							ID:   a.ID,
							Name: a.Name,
						})
					}
				}
			}

			contributor := commonmeta.Contributor{
				ID:               ID,
				Type:             Type,
				GivenName:        GivenName,
				FamilyName:       FamilyName,
				Name:             Name,
				ContributorRoles: []string{"Author"},
				Affiliations:     affiliations,
			}
			contributors = append(contributors, contributor)
		}
	}
	return contributors, nil
}

// GetFundingReferences returns the funding references from the JSON Feed metadata.
// Either provided by the blog metadata or via HasAward relationships
func GetFundingReferences(content Content) []commonmeta.FundingReference {
	var fundingReferences []commonmeta.FundingReference

	// Funding references from blog metadata
	if content.Blog.Funding.FunderName != "" {
		fundingReferences = append(fundingReferences, commonmeta.FundingReference{
			FunderName:           content.Blog.Funding.FunderName,
			FunderIdentifier:     content.Blog.Funding.FunderIdentifier,
			FunderIdentifierType: content.Blog.Funding.FunderIdentifierType,
			AwardTitle:           content.Blog.Funding.AwardTitle,
			AwardNumber:          content.Blog.Funding.AwardNumber,
			AwardURI:             content.Blog.Funding.AwardURI,
		})
	}
	if len(content.FundingReferences) > 0 {
		for _, v := range content.FundingReferences {
			fundingReferences = append(fundingReferences, commonmeta.FundingReference{
				FunderName:           v.FunderName,
				FunderIdentifier:     v.FunderIdentifier,
				FunderIdentifierType: v.FunderIdentifierType,
				AwardTitle:           v.AwardTitle,
				AwardNumber:          v.AwardNumber,
				AwardURI:             v.AwardURI,
			})
		}
	} else {
		// Funding references from relationships
		for _, v := range content.Relationships {
			if "HasAward" == v.Type {
				// Urls can either be a list of grant IDs or a funder identifier
				// (Open Funder Registry ID or ROR), followed by a grant URL
				if len(v.Urls) == 1 {
					prefix, _ := doiutils.ValidatePrefix(v.Urls[0])
					u, _ := url.Parse(v.Urls[0])
					if prefix == "10.3030" || u.Host == "cordis.europa.eu" {
						// Prefix 10.3030 means grant ID from funder is European Commission.
						// CORDIS is the grants portal of the European Commission.
						paths := strings.Split(u.Path, "/")
						awardNumber := paths[len(paths)-1]
						fundingReferences = append(fundingReferences, commonmeta.FundingReference{
							FunderName:           "European Commission",
							FunderIdentifier:     "https://doi.org/10.13039/501100000780",
							FunderIdentifierType: "Crossref Funder ID",
							AwardNumber:          awardNumber,
							AwardURI:             v.Urls[0],
						})

					}
				} else if len(v.Urls) == 2 {
					var funderName string
					prefix, _ := doiutils.ValidatePrefix(v.Urls[0])
					u, _ := url.Parse(v.Urls[1])
					if prefix == "10.13039" {
						// Prefix 10.13039 means funder ID from Open Funder registry.
						if v.Urls[0] == "https://doi.org/10.13039/100000001" {
							funderName = "National Science Foundation"
						}
						var awardNumber string
						if q := u.Query(); q != nil {
							awardNumber = q["awd_id"][0]
						} else {
							awardNumber = u.Path
						}
						fundingReferences = append(fundingReferences, commonmeta.FundingReference{
							FunderName:           funderName,
							FunderIdentifier:     v.Urls[0],
							FunderIdentifierType: "Crossref Funder ID",
							AwardNumber:          awardNumber,
							AwardURI:             v.Urls[1],
						})
					} else if _, ok := utils.ValidateROR(v.Urls[0]); ok {
						// URL is ROR ID for funder. Need to transform to Crossref Funder ID
						// until Crossref production service supports ROR IDs.
						ror, _ := utils.GetROR(v.Urls[0])
						funderIdentifier := ror.ExternalIds.FundRef.All[0]
						if funderIdentifier != "" {
							funderIdentifier = "https://doi.org/" + funderIdentifier
							var awardNumber string
							if q := u.Query(); q != nil {
								awardNumber = q["awd_id"][0]
							} else {
								paths := strings.Split(u.Path, "/")
								awardNumber = paths[len(paths)-1]
							}
							fundingReferences = append(fundingReferences, commonmeta.FundingReference{
								FunderName:           ror.Name,
								FunderIdentifier:     funderIdentifier,
								FunderIdentifierType: "Crossref Funder ID",
								AwardNumber:          awardNumber,
								AwardURI:             v.Urls[1],
							})
						}
					}
				}
			}
		}
	}
	return fundingReferences
}
