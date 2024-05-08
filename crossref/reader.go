// Package crossref provides function to convert Crossref metadata to/from the commonmeta metadata format.
package crossref

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
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/dateutils"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/utils"
)

// Content is the struct for the message in tge JSON response from the Crossref API
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
				ID     string `json:"id"`
				IDType string `json:"id-type"`
			} `json:"id"`
			Name string `json:"name"`
		} `json:"affiliation"`
	} `json:"author"`
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
		DateAsParts [][]int `json:"date-parts"`
		DateTime    string  `json:"date-time"`
	} `json:"published"`
	Issued struct {
		DateAsParts [][]int `json:"date-parts"`
		DateTime    string  `json:"date-time"`
	} `json:"issued"`
	Created struct {
		DateAsParts [][]int `json:"date-parts"`
		DateTime    string  `json:"date-time"`
	} `json:"created"`
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
		DOI          string `json:"DOI"`
		ArticleTitle string `json:"article-title"`
		Year         string `json:"year"`
		Unstructured string `json:"unstructured"`
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
		IsReviewedBy []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-reviewed-by"`
		Reviews []struct {
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
	Subtitle []string `json:"subtitle"`
	Title    []string `json:"title"`
	URL      string   `json:"url"`
	Version  string   `json:"version"`
	Volume   string   `json:"volume"`
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
var relationTypes = []string{"IsPartOf", "HasPart", "IsVariantFormOf", "IsOriginalFormOf", "IsIdenticalTo", "IsTranslationOf", "IsReviewedBy", "Reviews", "HasReview", "IsPreprintOf", "HasPreprint", "IsSupplementTo", "IsSupplementedBy"}

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

// FetchList gets the metadata for a list of works from the Crossref API and converts it to the Commonmeta format
func FetchList(number int, member string, _type string, sample bool, hasORCID bool, hasROR bool, hasReferences bool, hasRelation bool, hasAbstract bool, hasAward bool, hasLicense bool, hasArchive bool) ([]commonmeta.Data, error) {

	var data []commonmeta.Data
	content, err := GetList(number, member, _type, sample, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense, hasArchive)
	if err != nil {
		return data, err
	}

	data, err = ReadAll(content)
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
	url := "https://api.crossref.org/works/" + doi
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

// GetList gets the metadata for a list of works from the Crossref API
func GetList(number int, member string, _type string, sample bool, hasORCID bool, hasROR bool, hasReferences bool, hasRelation bool, hasAbstract bool, hasAward bool, hasLicense bool, hasArchive bool) ([]Content, error) {
	// the envelope for the JSON response from the Crossref API
	type Response struct {
		Status         string `json:"status"`
		MessageType    string `json:"message-type"`
		MessageVersion string `json:"message-version"`
		Message        struct {
			TotalResults int       `json:"total-results"`
			Items        []Content `json:"items"`
		} `json:"message`
	}
	var response Response
	if number > 100 {
		number = 100
	}
	url := QueryURL(number, member, _type, sample, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense, hasArchive)
	req, err := http.NewRequest("GET", url, nil)
	v := "0.1"
	u := "info@front-matter.io"
	userAgent := fmt.Sprintf("commonmeta/%s (https://commonmeta.org; mailto: %s)", v, u)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Cache-Control", "private")
	if err != nil {
		log.Fatalln(err)
	}
	client := http.Client{
		Timeout: 20 * time.Second,
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

// LoadList loads the metadata for a list of works from a JSON file and converts it to the Commonmeta format
func LoadList(filename string) ([]commonmeta.Data, error) {
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

// Read Crossref JSON response and return work struct in Commonmeta format
func Read(content Content) (commonmeta.Data, error) {
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
			var affiliations []commonmeta.Affiliation
			if len(v.Affiliation) > 0 {
				for _, a := range v.Affiliation {
					var ID string
					if len(a.ID) > 0 && a.ID[0].IDType == "ROR" {
						ID = utils.NormalizeROR(a.ID[0].ID)
					}
					if a.Name != "" {
						affiliations = append(affiliations, commonmeta.Affiliation{
							ID:   ID,
							Name: a.Name,
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
	if content.License != nil && len(content.License) > 0 {
		url, _ := utils.NormalizeCCUrl(content.License[0].URL)
		id := utils.URLToSPDX(url)
		data.License = commonmeta.License{
			ID:  id,
			URL: url,
		}
	}

	data.Provider = "Crossref"

	if content.Publisher != "" {
		data.Publisher = commonmeta.Publisher{
			Name: content.Publisher,
		}
	}

	for _, v := range content.Reference {
		reference := commonmeta.Reference{
			Key:             v.Key,
			ID:              doiutils.NormalizeDOI(v.DOI),
			Title:           v.ArticleTitle,
			PublicationYear: v.Year,
			Unstructured:    v.Unstructured,
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
				} else if v.IDType == "issn" {
					id = utils.ISSNAsURL(v.ID)
				} else if utils.ValidateURL(v.ID) == "URL" {
					id = v.ID
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
		data.Titles = append(data.Titles, commonmeta.Title{
			Title: content.Title[0],
		})
	}
	if len(content.Subtitle) > 0 && content.Subtitle[0] != "" {
		data.Titles = append(data.Titles, commonmeta.Title{
			Title: content.Subtitle[0],
			Type:  "Subtitle",
		})
	}
	if len(content.OriginalTitle) > 0 && content.OriginalTitle[0] != "" {
		data.Titles = append(data.Titles, commonmeta.Title{
			Title: content.OriginalTitle[0],
			Type:  "TranslatedTitle",
		})
	}

	data.URL = content.Resource.Primary.URL

	return data, nil
}

// ReadAll reads a list of Crossref JSON responses and returns a list of works in Commonmeta format
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

// QueryURL returns the URL for the Crossref API query
func QueryURL(number int, member string, _type string, sample bool, hasORCID bool, hasROR bool, hasReferences bool, hasRelation bool, hasAbstract bool, hasAward bool, hasLicense bool, hasArchive bool) string {
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
	if sample {
		values.Add("sample", strconv.Itoa(number))
	} else {
		values.Add("rows", strconv.Itoa(number))
	}

	// sort results by published date in descending order
	values.Add("sort", "published")
	values.Add("order", "desc")
	var filters []string
	if member != "" {
		filters = append(filters, "member:"+member)
	}
	if _type != "" && slices.Contains(types, _type) {
		filters = append(filters, "type:"+_type)
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
	return result.Message.PrimaryName, true
}
