// Package crossref provides function to convert Crossref metadata to/from the commonmeta metadata format.
package crossref

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"mvdan.cc/xurls/v2"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/dateutils"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/ror"
	"github.com/front-matter/commonmeta/utils"
)

type Reader struct {
	r *bufio.Reader
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: bufio.NewReader(r),
	}
}

// Content is the struct for the message in the JSON response from the Crossref API
type Content struct {
	ID       string   `json:"id"`
	Type     string   `json:"type"`
	Abstract string   `json:"abstract"`
	Archive  []string `json:"archive"`
	Author   []struct {
		Given       string `json:"given"`
		Family      string `json:"family"`
		Name        string `json:"name"`
		ORCID       string `json:"ORCID"`
		Sequence    string `json:"sequence"`
		Affiliation []struct {
			ID []struct {
				ID         string `json:"id"`
				IDType     string `json:"id-type"`
				AssertedBy string `json:"asserted-by"`
			} `json:"id"`
			Name string `json:"name"`
		} `json:"affiliation"`
	} `json:"author"`
	Member         string     `json:"member"`
	ContainerTitle []string   `json:"container-title"`
	DOI            string     `json:"doi"`
	Files          []struct{} `json:"files"`
	Funder         []struct {
		DOI   string   `json:"DOI"`
		Name  string   `json:"name"`
		Award []string `json:"award"`
	} `json:"funder"`
	GroupTitle string `json:"group-title"`
	Issue      string `json:"issue"`
	Published  struct {
		DateAsParts []dateutils.DateSlice `json:"date-parts"`
		DateTime    string                `json:"date-time"`
	} `json:"published"`
	Issued struct {
		DateAsParts []dateutils.DateSlice `json:"date-parts"`
		DateTime    string                `json:"date-time"`
	} `json:"issued"`
	Created struct {
		DateAsParts []dateutils.DateSlice `json:"date-parts"`
		DateTime    string                `json:"date-time"`
	} `json:"created"`
	Institution []struct {
		Name string `json:"name"`
	} `json:"institution"`
	ISSNType []struct {
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"issn-type"`
	ISBNType []struct {
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"isbn-type"`
	Language string `json:"language"`
	License  []struct {
		URL            string `json:"URL"`
		ContentVersion string `json:"content-version"`
	} `json:"license"`
	Link []struct {
		ContentType string `json:"content-type"`
		URL         string `json:"url"`
	} `json:"link"`
	OriginalTitle []string `json:"original-title"`
	Page          string   `json:"page"`
	PublishedAt   string   `json:"published_at"`
	Publisher     string   `json:"publisher"`
	Reference     []struct {
		Key          string `json:"key"`
		Type         string `json:"type"`
		DOI          string `json:"DOI"`
		ArticleTitle string `json:"article-title"`
		Year         string `json:"year"`
		Unstructured string `json:"unstructured"`
		AssertedBy   string `json:"doi-asserted-by"`
	} `json:"reference"`
	Relation struct {
		IsNewVersionOf []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-new-version-of"`
		IsPreviousVersionOf []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-previous-version-of"`
		IsVersionOf []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-version-of"`
		HasVersion []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"has-version"`
		IsPartOf []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-part-of"`
		HasPart []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"has-part"`
		IsVariantFormOf []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-variant-form-of"`
		IsOriginalFormOf []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-original-form-of"`
		IsIdenticalTo []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-identical-to"`
		IsTranslationOf []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-translation-of"`
		IsReviewOf []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"reviews"`
		HasReview []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"has-review"`
		IsPreprintOf []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-preprint-of"`
		HasPreprint []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"has-preprint"`
		IsSupplementTo []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-supplement-to"`
		IsSupplementedBy []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-supplemented-by"`
	} `json:"relation"`
	Resource struct {
		Primary struct {
			ContentType string `json:"content_type"`
			URL         string `json:"url"`
		} `json:"primary"`
	} `json:"resource"`
	Subject  []string `json:"subject"`
	Subtitle []string `json:"subtitle,omitempty"`
	Title    []string `json:"title"`
	URL      string   `json:"url"`
	Version  string   `json:"version,omitempty"`
	Volume   string   `json:"volume,omitempty"`
}

// CRToCMMappings maps Crossref types to Commonmeta types
// source: http://api.crossref.org/types
var CRToCMMappings = map[string]string{
	"book-chapter":        "BookChapter",
	"book-part":           "BookPart",
	"book-section":        "BookSection",
	"book-series":         "BookSeries",
	"book-set":            "BookSet",
	"book-track":          "BookTrack",
	"book":                "Book",
	"component":           "Component",
	"database":            "Database",
	"dataset":             "Dataset",
	"dissertation":        "Dissertation",
	"edited-book":         "Book",
	"grant":               "Grant",
	"journal-article":     "JournalArticle",
	"journal-issue":       "JournalIssue",
	"journal-volume":      "JournalVolume",
	"journal":             "Journal",
	"monograph":           "Book",
	"other":               "Other",
	"peer-review":         "PeerReview",
	"posted-content":      "Article",
	"proceedings-article": "ProceedingsArticle",
	"proceedings-series":  "ProceedingsSeries",
	"proceedings":         "Proceedings",
	"reference-book":      "Book",
	"reference-entry":     "Entry",
	"report-component":    "ReportComponent",
	"report-series":       "ReportSeries",
	"report":              "Report",
	"standard":            "Standard",
}

// CRCitationToCMMappings maps Crossref citation types to Commonmeta types
var CRCitationToCMMappings = map[string]string{
	"blog-post":              "BlogPost",
	"book":                   "Book",
	"book-chapter":           "BookChapter",
	"dataset":                "Dataset",
	"dissertation":           "Dissertation",
	"journal":                "Journal",
	"journal-article":        "JournalArticle",
	"patent":                 "Patent",
	"peer-review":            "PeerReview",
	"preprint":               "Preprint",
	"conference-proceedings": "Proceedings",
	"conference-paper":       "ProceedingsArticle",
	"protocol":               "Protocol",
	"report":                 "Report",
	"software":               "Software",
	"standard":               "Standard",
	"web-resource":           "Webpage",
	"other":                  "Other",
}

// CrossrefContainerTypes maps Crossref types to Crossref container types
var CrossrefContainerTypes = map[string]string{
	"book-chapter":        "book",
	"dataset":             "database",
	"journal-article":     "journal",
	"journal-issue":       "journal",
	"monograph":           "book-series",
	"proceedings-article": "proceedings",
	"posted-content":      "periodical",
}

// CRToCMContainerTranslations maps Crossref container types to Commonmeta container types
var CRToCMContainerTranslations = map[string]string{
	"book":        "Book",
	"book-series": "BookSeries",
	"database":    "DataRepository",
	"journal":     "Journal",
	"proceedings": "Proceedings",
	"periodical":  "Periodical",
}

// relation types to include
var relationTypes = []string{"IsPartOf", "HasPart", "IsVariantFormOf", "IsOriginalFormOf", "IsIdenticalTo", "IsTranslationOf", "IsReviewOf", "HasReview", "IsPreprintOf", "HasPreprint", "IsSupplementTo", "IsSupplementedBy"}

// Fetch gets the metadata for a single work from the Crossref API and converts it to the Commonmeta format
func Fetch(str string, match bool) (commonmeta.Data, error) {
	var data commonmeta.Data
	id, ok := doiutils.ValidateDOI(str)
	if !ok {
		return data, errors.New("invalid DOI")
	}
	content, err := Get(id)
	if err != nil {
		return data, err
	}
	data, err = Read(content, match)
	if err != nil {
		return data, err
	}
	return data, nil
}

// FetchAll gets the metadata for a list of works from the Crossref API and converts it to the Commonmeta format
func FetchAll(number int, page int, member string, type_ string, sample bool, year string, ror string, orcid string, hasORCID bool, hasROR bool, hasReferences bool, hasRelation bool, hasAbstract bool, hasAward bool, hasLicense bool, hasArchive bool, match bool) ([]commonmeta.Data, error) {
	var data []commonmeta.Data
	content, err := GetAll(number, page, member, type_, sample, year, orcid, ror, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense, hasArchive)
	if err != nil {
		return data, err
	}

	data, err = ReadAll(content, match)
	if err != nil {
		return data, err
	}
	return data, nil
}

// Get gets the metadata for a single work from the Crossref API
func Get(pid string) (Content, error) {
	// the envelope for the JSON response from the Crossref API
	type Response struct {
		Status         string  `json:"status"`
		MessageType    string  `json:"message-type"`
		MessageVersion string  `json:"message-version"`
		Message        Content `json:"message"`
	}

	var response Response
	doi, ok := doiutils.ValidateDOI(pid)
	if !ok {
		return response.Message, errors.New("invalid DOI")
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	url := "https://api.crossref.org/works/" + doi
	req, err := http.NewRequest(http.MethodGet, url, nil)
	u := "info@front-matter.io"
	userAgent := fmt.Sprintf("commonmeta/%s (https://commonmeta.org/; mailto: %s)", commonmeta.Version, u)
	req.Header.Set("User-Agent", userAgent)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return response.Message, err
	}
	if resp.StatusCode >= 400 {
		return response.Message, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.Message, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("error:", err)
	}
	return response.Message, err
}

// GetAll gets the metadata for a list of works from the Crossref API
func GetAll(number int, page int, member string, type_ string, sample bool, year string, ror string, orcid string, hasORCID bool, hasROR bool, hasReferences bool, hasRelation bool, hasAbstract bool, hasAward bool, hasLicense bool, hasArchive bool) ([]Content, error) {
	// the envelope for the JSON response from the Crossref API
	type Response struct {
		Status         string `json:"status"`
		MessageType    string `json:"message-type"`
		MessageVersion string `json:"message-version"`
		Message        struct {
			TotalResults int       `json:"total-results"`
			Items        []Content `json:"items"`
		}
	}
	var response Response
	if number > 100 {
		number = 1000
	}
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	url := QueryURL(number, page, member, type_, sample, year, orcid, ror, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense, hasArchive)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	u := "info@front-matter.io"
	userAgent := fmt.Sprintf("commonmeta/%s (https://commonmeta.org; mailto: %s)", commonmeta.Version, u)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Cache-Control", "private")
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("error:", err)
	}
	return response.Message.Items, nil
}

// Load loads the metadata for a single work from a JSON file
func Load(filename string, match bool) (commonmeta.Data, error) {
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
	data, err = Read(content, match)
	if err != nil {
		return data, err
	}
	return data, nil
}

// LoadAll loads the metadata for a list of works from a JSON file and converts it to the Commonmeta format
func LoadAll(filename string, match bool) ([]commonmeta.Data, error) {
	var data []commonmeta.Data
	var content []Content
	var err error

	extension := path.Ext(filename)
	if extension == ".jsonl" || extension == ".jsonlines" {
		var response []Content
		file, err := os.Open(filename)
		if err != nil {
			return nil, errors.New("error reading file")
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		for {
			var c Content
			if err := decoder.Decode(&c); err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			response = append(response, c)
		}
		content = response
	} else if extension == ".json" {
		type Response struct {
			Items []Content `json:"items"`
		}
		var response Response
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

	data, err = ReadAll(content, match)
	if err != nil {
		return data, err
	}
	return data, nil
}

// Read Crossref JSON response and return work struct in Commonmeta format
func Read(content Content, match bool) (commonmeta.Data, error) {
	var data = commonmeta.Data{}

	data.ID = doiutils.NormalizeDOI(content.DOI)
	data.Type = CRToCMMappings[content.Type]
	if data.Type == "" {
		data.Type = "Other"
	}
	containerType := CrossrefContainerTypes[content.Type]
	containerType = CRToCMContainerTranslations[containerType]

	for _, v := range content.Archive {
		if !slices.Contains(data.ArchiveLocations, v) {
			data.ArchiveLocations = append(data.ArchiveLocations, v)
		}
	}

	var identifier, identifierType string
	if len(content.ISSNType) > 0 {
		i := make(map[string]string)
		for _, issn := range content.ISSNType {
			i[issn.Type] = issn.Value
		}
		if i["electronic"] != "" {
			identifier = i["electronic"]
			identifierType = "ISSN"
		} else if i["print"] != "" {
			identifier = i["print"]
			identifierType = "ISSN"
		}
	} else if len(content.ISBNType) > 0 {
		i := make(map[string]string)
		for _, isbn := range content.ISBNType {
			i[isbn.Type] = isbn.Value
		}
		if i["electronic"] != "" {
			identifier = i["electronic"]
			identifierType = "ISBN"
		} else if i["print"] != "" {
			identifier = i["print"]
			identifierType = "ISBN"
		}
	}
	var containerTitle string
	if len(content.ContainerTitle) > 0 {
		containerTitle = content.ContainerTitle[0]
	} else if len(content.Institution) > 0 {
		containerTitle = content.Institution[0].Name
	}
	var lastPage string
	pages := strings.Split(content.Page, "-")
	firstPage := pages[0]
	if len(pages) > 1 {
		lastPage = pages[1]
	}
	data.Container = commonmeta.Container{
		Identifier:     identifier,
		IdentifierType: identifierType,
		Type:           containerType,
		Title:          containerTitle,
		Volume:         content.Volume,
		Issue:          content.Issue,
		FirstPage:      firstPage,
		LastPage:       lastPage,
	}

	for _, v := range content.Author {
		if v.Name != "" || v.Given != "" || v.Family != "" {
			var ID, Type string
			if v.ORCID != "" {
				// enforce HTTPS
				ID, _ = utils.NormalizeURL(v.ORCID, true, false)
			}
			if v.Name != "" {
				Type = "Organization"
			} else {
				Type = "Person"
			}
			var affiliations []*commonmeta.Affiliation
			if len(v.Affiliation) > 0 {
				for _, a := range v.Affiliation {
					var ID, assertedBy string
					if len(a.ID) > 0 && a.ID[0].IDType == "ROR" {
						ID = utils.NormalizeROR(a.ID[0].ID)
						assertedBy = a.ID[0].AssertedBy
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
			}

			contributor := commonmeta.Contributor{
				ID:               ID,
				Type:             Type,
				GivenName:        v.Given,
				FamilyName:       v.Family,
				Name:             v.Name,
				ContributorRoles: []string{"Author"},
				Affiliations:     affiliations,
			}
			containsName := slices.ContainsFunc(data.Contributors, func(e commonmeta.Contributor) bool {
				return e.Name != "" && e.Name == contributor.Name || e.GivenName == contributor.GivenName && e.FamilyName != "" && e.FamilyName == contributor.FamilyName
			})
			if !containsName {
				data.Contributors = append(data.Contributors, contributor)
			}
		}
	}

	if content.Published.DateTime != "" {
		data.Date.Published = content.Published.DateTime
	} else if len(content.Published.DateAsParts) > 0 {
		data.Date.Published = dateutils.GetDateFromDateParts(content.Published.DateAsParts)
	} else if content.Issued.DateTime != "" {
		data.Date.Published = content.Issued.DateTime
	} else if len(content.Issued.DateAsParts) > 0 {
		data.Date.Published = dateutils.GetDateFromDateParts(content.Issued.DateAsParts)
	}
	if data.Date.Published == "" {
		if content.Created.DateTime != "" {
			data.Date.Published = content.Created.DateTime
		} else if len(content.Created.DateAsParts) > 0 {
			data.Date.Published = dateutils.GetDateFromDateParts(content.Created.DateAsParts)
		}
	}

	if content.Abstract != "" {
		abstract := utils.Sanitize(content.Abstract)
		data.Descriptions = append(data.Descriptions, commonmeta.Description{
			Description: abstract,
			Type:        "Abstract",
		})
	}

	for _, v := range content.Link {
		if v.ContentType != "unspecified" {
			data.Files = append(data.Files, commonmeta.File{
				URL:      v.URL,
				MimeType: v.ContentType,
			})
		}
	}
	if len(content.Link) > 1 {
		data.Files = utils.DedupeSlice(data.Files)
	}

	if len(content.Funder) > 1 {
		for _, v := range content.Funder {
			funderIdentifier := doiutils.NormalizeDOI(v.DOI)
			var funderIdentifierType string
			if strings.HasPrefix(v.DOI, "10.13039") {
				funderIdentifierType = "Crossref Funder ID"
			}
			if len(v.Award) > 0 {
				for _, award := range v.Award {
					data.FundingReferences = append(data.FundingReferences, commonmeta.FundingReference{
						FunderIdentifier:     funderIdentifier,
						FunderIdentifierType: funderIdentifierType,
						FunderName:           v.Name,
						AwardNumber:          award,
					})
				}
			} else {
				data.FundingReferences = append(data.FundingReferences, commonmeta.FundingReference{
					FunderIdentifier:     funderIdentifier,
					FunderIdentifierType: funderIdentifierType,
					FunderName:           v.Name,
				})
			}
		}
		// if len(content.Funder) > 1 {
		data.FundingReferences = utils.DedupeSlice(data.FundingReferences)
		// }
	}

	data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
		Identifier:     data.ID,
		IdentifierType: "DOI",
	})

	data.Language = content.Language
	if len(content.License) > 0 {
		url, _ := utils.NormalizeCCUrl(content.License[0].URL)
		id := utils.URLToSPDX(url)
		data.License = commonmeta.License{
			ID:  id,
			URL: url,
		}
	}

	data.Provider = "Crossref"

	if content.Publisher != "" || content.Member != "" {
		var id string
		if content.Member != "" {
			id = fmt.Sprintf("https://api.crossref.org/members/%s", content.Member)
		}
		data.Publisher = commonmeta.Publisher{
			ID:   id,
			Name: content.Publisher,
		}
	}
	// workaround until Crossref supports BlogPost as posted-content type
	if data.Type == "Article" && data.Publisher.Name == "Front Matter" {
		data.Type = "BlogPost"
	}

	rxStrict := xurls.Strict()
	for _, v := range content.Reference {
		ID := doiutils.NormalizeDOI(v.DOI)
		if ID == "" && v.Unstructured != "" {
			ID = rxStrict.FindString(v.Unstructured)
		}
		type_ := CRCitationToCMMappings[v.Type]
		if type_ == "" {
			type_ = "Other"
		}
		reference := commonmeta.Reference{
			Key:             v.Key,
			Type:            type_,
			ID:              ID,
			Title:           v.ArticleTitle,
			PublicationYear: v.Year,
			Unstructured:    v.Unstructured,
			AssertedBy:      v.AssertedBy,
		}
		containsKey := slices.ContainsFunc(data.References, func(e commonmeta.Reference) bool {
			return e.Key != "" && e.Key == reference.Key
		})
		if !containsKey {
			data.References = append(data.References, reference)
		}
	}

	fields := reflect.VisibleFields(reflect.TypeOf(content.Relation))
	for _, field := range fields {
		if slices.Contains(relationTypes, field.Name) {
			relationByType := reflect.ValueOf(content.Relation).FieldByName(field.Name)
			for _, v := range relationByType.Interface().([]struct {
				ID     string `json:"id"`
				IDType string `json:"id-type"`
			}) {
				var id string
				if v.IDType == "doi" {
					id = doiutils.NormalizeDOI(v.ID)
				} else if utils.ValidateURL(v.ID) == "URL" {
					id = v.ID
				} else if v.IDType == "issn" {
					data.Container.IdentifierType = "ISSN"
					data.Container.Identifier = v.ID
				}
				relation := commonmeta.Relation{
					ID:   id,
					Type: field.Name,
				}
				if id != "" && !slices.Contains(data.Relations, relation) {
					data.Relations = append(data.Relations, relation)
				}
			}
			sort.Slice(data.Relations, func(i, j int) bool {
				return data.Relations[i].Type < data.Relations[j].Type
			})
		}
	}

	// add relation to subject area community
	if content.GroupTitle != "" && data.Type == "BlogPost" {
		groupTitle := utils.WordsToCamelCase(content.GroupTitle)
		data.Relations = append(data.Relations, commonmeta.Relation{
			ID:   utils.CommunitySlugAsURL(groupTitle, "rogue-scholar.org"),
			Type: "IsPartOf",
		})
	}
	if data.Container.IdentifierType == "ISSN" {
		data.Relations = append(data.Relations, commonmeta.Relation{
			ID:   utils.ISSNAsURL(data.Container.Identifier),
			Type: "IsPartOf",
		})
	}

	for _, v := range content.Subject {
		subject := commonmeta.Subject{
			Subject: v,
		}
		if !slices.Contains(data.Subjects, subject) {
			data.Subjects = append(data.Subjects, subject)
		}
	}

	if content.GroupTitle != "" {
		data.Subjects = append(data.Subjects, commonmeta.Subject{
			Subject: content.GroupTitle,
		})
	}

	if len(content.Title) > 0 && content.Title[0] != "" {
		for _, v := range content.Title {
			data.Titles = append(data.Titles, commonmeta.Title{
				Title: v,
			})
		}
	}
	if len(content.Subtitle) > 0 && content.Subtitle[0] != "" {
		for _, v := range content.Subtitle {
			data.Titles = append(data.Titles, commonmeta.Title{
				Title: v,
				Type:  "Subtitle",
			})
		}
	}
	if len(content.OriginalTitle) > 0 && content.OriginalTitle[0] != "" {
		for _, v := range content.OriginalTitle {
			data.Titles = append(data.Titles, commonmeta.Title{
				Title: v,
				Type:  "OriginalTitle",
			})
		}
	}

	data.URL = content.Resource.Primary.URL

	return data, nil
}

// ReadAll reads a list of Crossref JSON responses and returns a list of works in Commonmeta format
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

// QueryURL returns the URL for the Crossref API query
func QueryURL(number int, page int, member string, type_ string, sample bool, year string, orcid string, ror string, hasORCID bool, hasROR bool, hasReferences bool, hasRelation bool, hasAbstract bool, hasAward bool, hasLicense bool, hasArchive bool) string {
	types := []string{
		"book",
		"book-chapter",
		"book-part",
		"book-section",
		"book-series",
		"book-set",
		"book-track",
		"component",
		"database",
		"dataset",
		"dissertation",
		"edited-book",
		"grant",
		"journal",
		"journal-article",
		"journal-issue",
		"journal-volume",
		"monograph",
		"other",
		"peer-review",
		"posted-content",
		"proceedings",
		"proceedings-article",
		"proceedings-series",
		"reference-book",
		"reference-entry",
		"report",
		"report-component",
		"report-series",
		"standard",
	}

	u, _ := url.Parse("https://api.crossref.org/works")
	values := u.Query()
	if number <= 0 {
		number = 10
	}
	if number > 1000 {
		number = 1000
	}
	if page <= 0 {
		page = 1
	}
	if sample {
		values.Add("sample", strconv.Itoa(number))
	} else {
		values.Add("rows", strconv.Itoa(number))
		values.Add("offset", strconv.Itoa((page-1)*number))
	}

	// sort results by published date in descending order
	values.Add("sort", "published")
	values.Add("order", "desc")
	var filters []string
	if member != "" {
		filters = append(filters, "member:"+member)
	}
	if type_ != "" && slices.Contains(types, type_) {
		filters = append(filters, "type:"+type_)
	}
	if ror != "" {
		r, _ := utils.ValidateROR(ror)
		if r != "" {
			filters = append(filters, "ror-id:"+r)
		}
	}
	if orcid != "" {
		o, _ := utils.ValidateORCID(orcid)
		if o != "" {
			filters = append(filters, "orcid:"+o)
		}
	}
	if year != "" {
		filters = append(filters, "from-pub-date:"+year+"-01-01")
		filters = append(filters, "until-pub-date:"+year+"-12-31")
	}
	if hasORCID {
		filters = append(filters, "has-orcid:true")
	}
	if hasROR {
		filters = append(filters, "has-ror-id:true")
	}
	if hasReferences {
		filters = append(filters, "has-references:true")
	}
	if hasRelation {
		filters = append(filters, "has-relation:true")
	}
	if hasAbstract {
		filters = append(filters, "has-abstract:true")
	}
	if hasAward {
		filters = append(filters, "has-award:true")
	}
	if hasLicense {
		filters = append(filters, "has-license:true")
	}
	if hasArchive {
		filters = append(filters, "has-archive:true")
	}
	if len(filters) > 0 {
		values.Add("filter", strings.Join(filters[:], ","))
	}
	u.RawQuery = values.Encode()
	return u.String()
}

// Get the Crossref member name for a given memberId
func GetMember(memberId string) (string, bool) {
	type Response struct {
		Message struct {
			PrimaryName string `json:"primary-name"`
		} `json:"message"`
	}
	var result Response
	if memberId == "" {
		return "", false
	}
	resp, err := http.Get(fmt.Sprintf("https://api.crossref.org/members/%s", memberId))
	if err != nil {
		return "", false
	}
	if resp.StatusCode == 404 {
		return "", false
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", false
	}
	return string(result.Message.PrimaryName), true
}
