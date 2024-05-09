// Package jsonfeed converts JSON Feed metadata to/from the commonmeta metadata format.
package jsonfeed

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/front-matter/commonmeta/authorutils"
	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/dateutils"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/utils"
)

// Content represents the JSON Feed metadata.
type Content struct {
	ID            string      `json:"id"`
	DOI           string      `json:"doi"`
	GUID          string      `json:"guid"`
	Abstract      string      `json:"abstract"`
	ArchiveURL    string      `json:"archive_url"`
	Authors       Authors     `json:"authors"`
	Blog          Blog        `json:"blog"`
	BlogName      string      `json:"blog_name"`
	BlogSlug      string      `json:"blog_slug"`
	Image         string      `json:"image"`
	IndexedAt     int64       `json:"indexed_at"`
	Language      string      `json:"language"`
	PublishedAt   int64       `json:"published_at"`
	Relationships []Relation  `json:"relationships"`
	Reference     []Reference `json:"reference"`
	Summary       string      `json:"summary"`
	Title         string      `json:"title"`
	UpdatedAt     int64       `json:"updated_at"`
	URL           string      `json:"url"`
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
	ID          string  `json:"id"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Favicon     string  `json:"favicon"`
	Funding     Funding `json:"funding"`
	HomePageURL string  `json:"home_page_url"`
	ISSN        string  `json:"issn"`
	Language    string  `json:"language"`
	License     string  `json:"license"`
	Prefix      string  `json:"prefix"`
	Slug        string  `json:"slug"`
	Title       string  `json:"title"`
}

type Funding struct {
	AwardNumber string `json:"award_number"`
	AwardURI    string `json:"award_uri"`
	FunderID    string `json:"funder_id"`
	FunderName  string `json:"funder_name"`
}

// Relation represents a relation in the JSON Feed item.
type Relation struct {
	Type string   `json:"type"`
	Urls []string `json:"urls"`
}

// Reference represents a reference in the JSON Feed item.
type Reference struct {
	Key             string `json:"key"`
	DOI             string `json:"doi"`
	URL             string `json:"url"`
	PublicationYear string `json:"publicationYear"`
	Title           string `json:"title"`
}

// relation types to include
var relationTypes = []string{"IsPartOf", "HasPart", "IsVariantFormOf", "IsOriginalFormOf", "IsIdenticalTo", "IsTranslationOf", "IsReviewedBy", "Reviews", "HasReview", "IsPreprintOf", "HasPreprint", "IsSupplementTo", "IsSupplementedBy"}

// Fetch fetches JSON Feed metadata and returns Commonmeta metadata.
func Fetch(str string) (commonmeta.Data, error) {
	var data commonmeta.Data
	id, ok := utils.ValidateUUID(str)
	if !ok {
		return data, errors.New("invalid uuid")
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

// Get retrieves JSON Feed metadata.
func Get(id string) (Content, error) {
	var content Content
	client := http.Client{
		Timeout: time.Second * 10,
	}
	url := "https://api.rogue-scholar.org/posts/" + id
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

// Read reads JSON Feed metadata and converts it into Commonmeta metadata.
func Read(content Content) (commonmeta.Data, error) {
	var data commonmeta.Data

	if content.DOI != "" {
		data.ID = doiutils.NormalizeDOI(content.DOI)
	} else {
		data.ID = content.ID
	}
	data.Type = "Article"

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
	data.Container = commonmeta.Container{
		Type:           "Periodical",
		Title:          content.Blog.Title,
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

	if doiutils.IsRogueScholarDOI(data.ID) {
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
			Identifier:     content.DOI,
			IdentifierType: "DOI",
		})

		data.Provider = "Crossref"
	}

	data.FundingReferences = GetFundingReferences(content)

	data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
		Identifier:     content.ID,
		IdentifierType: "UUID",
	})

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
		Name: content.Blog.Title,
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
			ID:              doiutils.NormalizeDOI(v.DOI),
			Title:           v.Title,
			PublicationYear: v.PublicationYear,
		}
		containsKey := slices.ContainsFunc(data.References, func(e commonmeta.Reference) bool {
			return e.Key != "" && e.Key == reference.Key
		})
		if !containsKey {
			data.References = append(data.References, reference)
		}
	}

	if content.Blog.Category != "" {
		data.Subjects = []commonmeta.Subject{
			{Subject: content.Blog.Category},
		}
	}

	data.Titles = []commonmeta.Title{
		{Title: utils.Sanitize(content.Title)},
	}

	url, err := utils.NormalizeURL(content.URL, true, true)
	if err != nil {
		return data, err
	}
	data.URL = url
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

			var affiliations []commonmeta.Affiliation
			if len(v.Affiliation) > 0 {
				for _, a := range v.Affiliation {
					if a.Name != "" {
						affiliations = append(affiliations, commonmeta.Affiliation{
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
			FunderIdentifier:     content.Blog.Funding.FunderID,
			FunderIdentifierType: "Crossref Funder ID",
			AwardNumber:          content.Blog.Funding.AwardNumber,
			AwardURI:             content.Blog.Funding.AwardURI,
		})
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
					} else if _, ok := utils.ValidateROR(v.Urls[0]); ok == true {
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
