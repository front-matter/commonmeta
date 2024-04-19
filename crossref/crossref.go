package crossref

import (
	"commonmeta/dateutils"
	"commonmeta/doiutils"
	"commonmeta/types"
	"commonmeta/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

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

var CrossrefContainerTypes = map[string]string{
	"book-chapter":        "book",
	"dataset":             "database",
	"journal-article":     "journal",
	"journal-issue":       "journal",
	"monograph":           "book-series",
	"proceedings-article": "proceedings",
	"posted-content":      "periodical",
}

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

func FetchCrossref(str string) (types.Data, error) {
	var data types.Data
	id, ok := doiutils.ValidateDOI(str)
	if !ok {
		return data, errors.New("invalid DOI")
	}
	content, err := GetCrossref(id)
	if err != nil {
		return data, err
	}
	data, err = ReadCrossref(content)
	if err != nil {
		return data, err
	}
	return data, nil
}

func FetchCrossrefSample(number int, member string, _type string) ([]types.Data, error) {
	var data []types.Data
	content, err := GetCrossrefSample(number, member, _type)
	if err != nil {
		return data, err
	}
	for _, v := range content {
		d, err := ReadCrossref(v)
		if err != nil {
			log.Println(err)
		}
		data = append(data, d)
	}
	return data, nil
}

func GetCrossref(pid string) (types.Content, error) {
	// the envelope for the JSON response from the Crossref API
	type Response struct {
		Status         string        `json:"status"`
		MessageType    string        `json:"message-type"`
		MessageVersion string        `json:"message-version"`
		Message        types.Content `json:"message"`
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
	userAgent := fmt.Sprintf("commonmeta-go/%s (https://commonmeta.org/commonmeta-go/; mailto: %s)", v, u)
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

// read Crossref JSON response and return work struct in Commonmeta format
func ReadCrossref(content types.Content) (types.Data, error) {
	var data = types.Data{}

	data.ID = doiutils.DOIAsUrl(content.DOI)
	data.Type = CRToCMMappings[content.Type]
	data.Url = content.Resource.Primary.URL

	for _, v := range content.Author {
		if v.Name != "" || v.Given != "" || v.Family != "" {
			var ID, Type string
			if v.ORCID != "" {
				// enforce HTTPS
				ID, _ = utils.NormalizeUrl(v.ORCID, true, false)
			}
			if v.Name != "" {
				Type = "Organization"
			} else {
				Type = "Person"
			}
			var affiliations []types.Affiliation
			if len(v.Affiliation) > 0 {
				for _, a := range v.Affiliation {
					affiliations = append(affiliations, types.Affiliation{
						Name: a.Name,
					})
				}
			}
			data.Contributors = append(data.Contributors, types.Contributor{
				ID:               ID,
				Type:             Type,
				GivenName:        v.Given,
				FamilyName:       v.Family,
				Name:             v.Name,
				ContributorRoles: []string{"Author"},
				Affiliations:     affiliations,
			})

		}
	}

	if content.Publisher != "" {
		data.Publisher = types.Publisher{
			Name: content.Publisher,
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

	if len(content.Title) > 0 {
		data.Titles = append(data.Titles, types.Title{
			Title: content.Title[0],
		})
	}
	if len(content.Subtitle) > 0 {
		data.Titles = append(data.Titles, types.Title{
			Title:     content.Subtitle[0],
			TitleType: "Subtitle",
		})
	}

	if content.Abstract != "" {
		abstract := utils.Sanitize(content.Abstract)
		data.Descriptions = append(data.Descriptions, types.Description{
			Description:     abstract,
			DescriptionType: "Abstract",
		})
	}

	containerType := CrossrefContainerTypes[content.Type]
	containerType = CRToCMContainerTranslations[containerType]
	var identifier, identifierType string
	if content.ISSN != nil {
		identifier = content.ISSN[0]
		identifierType = "ISSN"
	}
	if len(content.ISBNType) > 0 {
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
	data.Container = types.Container{
		Identifier:     identifier,
		IdentifierType: identifierType,
		Type:           containerType,
		Title:          containerTitle,
		Volume:         content.Volume,
		Issue:          content.Issue,
		FirstPage:      firstPage,
		LastPage:       lastPage,
	}

	for _, v := range content.Subject {
		subject := types.Subject{
			Subject: v,
		}
		if !slices.Contains(data.Subjects, subject) {
			data.Subjects = append(data.Subjects, subject)
		}
	}

	if content.GroupTitle != "" {
		data.Subjects = append(data.Subjects, types.Subject{
			Subject: content.GroupTitle,
		})
	}

	for _, v := range content.Reference {
		data.References = append(data.References, types.Reference{
			Key:             v.Key,
			Doi:             doiutils.DOIAsUrl(v.DOI),
			Title:           v.ArticleTitle,
			PublicationYear: v.Year,
			Unstructured:    v.Unstructured,
		})
	}

	fields := reflect.VisibleFields(reflect.TypeOf(content.Relation))
	for _, field := range fields {
		if slices.Contains(relationTypes, field.Name) {
			relationByType := reflect.ValueOf(content.Relation).FieldByName(field.Name)
			for _, v := range relationByType.Interface().([]struct {
				ID     string `json:"id"`
				IDType string `json:"id-type"`
			}) {
				data.Relations = append(data.Relations, types.Relation{
					ID:   doiutils.DOIAsUrl(v.ID),
					Type: field.Name,
				})
			}
			sort.Slice(data.Relations, func(i, j int) bool {
				return data.Relations[i].Type < data.Relations[j].Type
			})
		}
	}
	if content.ISSN != nil {
		data.Relations = append(data.Relations, types.Relation{
			ID:   utils.IssnAsUrl(content.ISSN[0]),
			Type: "IsPartOf",
		})
	}

	for _, v := range content.Funder {
		funderIdentifier := doiutils.DOIAsUrl(v.DOI)
		var funderIdentifierType string
		if strings.HasPrefix(v.DOI, "10.13039") {
			funderIdentifierType = "Crossref Funder ID"
		}
		if len(v.Award) > 0 {
			for _, award := range v.Award {
				data.FundingReferences = append(data.FundingReferences, types.FundingReference{
					FunderIdentifier:     funderIdentifier,
					FunderIdentifierType: funderIdentifierType,
					FunderName:           v.Name,
					AwardNumber:          award,
				})
			}
		} else {
			data.FundingReferences = append(data.FundingReferences, types.FundingReference{
				FunderIdentifier:     funderIdentifier,
				FunderIdentifierType: funderIdentifierType,
				FunderName:           v.Name,
			})
		}
	}
	data.FundingReferences = utils.DedupeSlice(data.FundingReferences)

	data.Language = content.Language
	if content.License != nil && len(content.License) > 0 {
		url, _ := utils.NormalizeCCUrl(content.License[0].Url)
		id := utils.UrlToSPDX(url)
		data.License = types.License{
			ID:  id,
			Url: url,
		}
	}
	data.Provider = "Crossref"
	for _, v := range content.Link {
		if v.ContentType != "unspecified" {
			data.Files = append(data.Files, types.File{
				Url:      v.Url,
				MimeType: v.ContentType,
			})
		}
	}
	data.Files = utils.DedupeSlice(data.Files)

	copy(data.ArchiveLocations, content.Archive)

	return data, nil
}

func GetCrossrefSample(number int, member string, _type string) ([]types.Content, error) {
	// the envelope for the JSON response from the Crossref API
	type Response struct {
		Status         string `json:"status"`
		MessageType    string `json:"message-type"`
		MessageVersion string `json:"message-version"`
		Message        struct {
			TotalResults int             `json:"total-results"`
			Items        []types.Content `json:"items"`
		} `json:"message`
	}
	var response Response
	if number > 100 {
		number = 100
	}
	url := CrossrefApiSampleUrl(number, member, _type)
	req, err := http.NewRequest("GET", url, nil)
	v := "0.1"
	u := "info@front-matter.io"
	userAgent := fmt.Sprintf("commonmeta-go/%s (https://commonmeta.org/commonmeta-go/; mailto: %s)", v, u)
	req.Header.Set("User-Agent", userAgent)
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
	log.Println("Total results:", response.Message.TotalResults)
	return response.Message.Items, nil
}

func CrossrefApiSampleUrl(number int, member string, _type string) string {
	types := []string{
		"book-section",
		"monograph",
		"report-component",
		"report",
		"peer-review",
		"book-track",
		"journal-article",
		"book-part",
		"other",
		"book",
		"journal-volume",
		"book-set",
		"reference-entry",
		"proceedings-article",
		"journal",
		"component",
		"book-chapter",
		"proceedings-series",
		"report-series",
		"proceedings",
		"database",
		"standard",
		"reference-book",
		"posted-content",
		"journal-issue",
		"dissertation",
		"grant",
		"dataset",
		"book-series",
		"edited-book",
		"journal-section",
		"monograph-series",
		"journal-meta",
		"book-series-meta",
		"component-list",
		"journal-issue-meta",
		"journal-meta",
		"book-part-meta",
		"book-meta",
		"proceedings-meta",
		"book-series-meta",
		"book-set",
	}
	u, _ := url.Parse("https://api.crossref.org/works")
	values := u.Query()
	values.Add("sample", strconv.Itoa(number))
	var filters []string
	if member != "" {
		filters = append(filters, "member:"+member)
	}
	if _type != "" && slices.Contains(types, _type) {
		filters = append(filters, "type:"+_type)
	}
	if len(filters) > 0 {
		values.Add("filter", strings.Join(filters[:], ","))
	}
	u.RawQuery = values.Encode()
	return u.String()
}
