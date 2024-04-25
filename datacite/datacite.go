// Package datacite provides function to convert DataCite metadata to/from the commonmeta metadata format.
package datacite

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/doiutils"

	"github.com/front-matter/commonmeta/utils"
)

// Content represents the DataCite metadata.
type Content struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

// Attributes represents the attributes of the DataCite JSONAPI response.
type Attributes struct {
	DOI                  string `json:"doi"`
	Prefix               string `json:"prefix"`
	Suffix               string `json:"suffix"`
	AlternateIdentifiers []struct {
		AlternateIdentifier     string `json:"alternateIdentifier"`
		AlternateIdentifierType string `json:"alternateIdentifierType"`
	} `json:"alternateIdentifiers"`
	Creators  []Contributor `json:"creators"`
	Publisher string        `json:"publisher"`
	Container struct {
		Type           string `json:"type"`
		Identifier     string `json:"identifier"`
		IdentifierType string `json:"identifierType"`
		Title          string `json:"title"`
		Volume         string `json:"volume"`
		Issue          string `json:"issue"`
		FirstPage      string `json:"firstPage"`
		LastPage       string `json:"lastPage"`
	} `json:"container"`
	PublicationYear int `json:"publicationYear"`
	Titles          []struct {
		Title     string `json:"title"`
		TitleType string `json:"titleType"`
		Lang      string `json:"lang"`
	} `json:"titles"`
	URL      string `json:"url"`
	Subjects []struct {
		Subject string `json:"subject"`
	} `json:"subjects"`
	Contributors []Contributor `json:"contributors"`
	Dates        []struct {
		Date            string `json:"date"`
		DateType        string `json:"dateType"`
		DateInformation string `json:"dateInformation"`
	} `json:"dates"`
	Language string `json:"language"`
	Types    struct {
		ResourceTypeGeneral string `json:"resourceTypeGeneral"`
		ResourceType        string `json:"resourceType"`
	} `json:"types"`
	RelatedIdentifiers []struct {
		RelatedIdentifier     string `json:"relatedIdentifier"`
		RelatedIdentifierType string `json:"relatedIdentifierType"`
		RelationType          string `json:"relationType"`
	} `json:"relatedIdentifiers"`
	Sizes      []string `json:"sizes"`
	Formats    []string `json:"formats"`
	Version    string   `json:"version"`
	RightsList []struct {
		Rights                 string `json:"rights"`
		RightsURI              string `json:"rightsUri"`
		SchemeURI              string `json:"schemeUri"`
		RightsIdentifier       string `json:"rightsIdentifier"`
		RightsIdentifierScheme string `json:"rightsIdentifierScheme"`
	}
	Descriptions []struct {
		Description     string `json:"description"`
		DescriptionType string `json:"descriptionType"`
		Lang            string `json:"lang"`
	} `json:"descriptions"`
	GeoLocations []struct {
		GeoLocationPoint struct {
			PointLongitude float64 `json:"pointLongitude,string"`
			PointLatitude  float64 `json:"pointLatitude,string"`
		} `json:"geoLocationPoint"`
		GeoLocationBox struct {
			WestBoundLongitude float64 `json:"westBoundLongitude,string"`
			EastBoundLongitude float64 `json:"eastBoundLongitude,string"`
			SouthBoundLatitude float64 `json:"southBoundLatitude,string"`
			NorthBoundLatitude float64 `json:"northBoundLatitude,string"`
		} `json:"geoLocationBox"`
		GeoLocationPlace string `json:"geoLocationPlace"`
	} `json:"geoLocations"`
	FundingReferences []struct {
		FunderName           string `json:"funderName"`
		FunderIdentifier     string `json:"funderIdentifier"`
		FunderIdentifierType string `json:"funderIdentifierType"`
		AwardNumber          string `json:"awardNumber"`
		AwardURI             string `json:"awardUri"`
	} `json:"fundingReferences"`
}

// Contributor represents the contributor of the DataCite JSONAPI response.
type Contributor struct {
	Name            string `json:"name"`
	GivenName       string `json:"givenName"`
	FamilyName      string `json:"familyName"`
	NameType        string `json:"nameType"`
	NameIdentifiers []struct {
		SchemeURI            string `json:"schemeUri"`
		NameIdentifier       string `json:"nameIdentifier"`
		NameIdentifierScheme string `json:"nameIdentifierScheme"`
	} `json:"nameIdentifiers"`
	Affiliation []struct {
		AffiliationIdentifier       string `json:"affiliationIdentifier"`
		AffiliationIdentifierScheme string `json:"affiliationIdentifierScheme"`
		Name                        string `json:"name"`
	} `json:"affiliation"`
	ContributorType string `json:"contributorType"`
}

// DCToCMTranslations maps DataCite resource types to Commonmeta types
// source: https://github.com/datacite/schema/blob/master/source/meta/kernel-4/include/datacite-resourceType-v4.xsd
var DCToCMTranslations = map[string]string{
	"Audiovisual":           "Audiovisual",
	"BlogPosting":           "Article",
	"Book":                  "Book",
	"BookChapter":           "BookChapter",
	"Collection":            "Collection",
	"ComputationalNotebook": "ComputationalNotebook",
	"ConferencePaper":       "ProceedingsArticle",
	"ConferenceProceeding":  "Proceedings",
	"DataPaper":             "JournalArticle",
	"Dataset":               "Dataset",
	"Dissertation":          "Dissertation",
	"Event":                 "Event",
	"Image":                 "Image",
	"Instrument":            "Instrument",
	"InteractiveResource":   "InteractiveResource",
	"Journal":               "Journal",
	"JournalArticle":        "JournalArticle",
	"Model":                 "Model",
	"OutputManagementPlan":  "OutputManagementPlan",
	"PeerReview":            "PeerReview",
	"PhysicalObject":        "PhysicalObject",
	"Poster":                "Presentation",
	"Preprint":              "Article",
	"Report":                "Report",
	"Service":               "Service",
	"Software":              "Software",
	"Sound":                 "Sound",
	"Standard":              "Standard",
	"StudyRegistration":     "StudyRegistration",
	"Text":                  "Document",
	"Thesis":                "Dissertation",
	"Workflow":              "Workflow",
	"Other":                 "Other",
}

// Fetch fetches DataCite metadata for a given DOI and returns Commonmeta metadata.
func Fetch(str string) (commonmeta.Data, error) {
	var data commonmeta.Data
	id, ok := doiutils.ValidateDOI(str)
	if !ok {
		return data, errors.New("invalid doi")
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

// FetchList gets the metadata for a list of works from the DataCite API and returns Commonmeta metadata.
func FetchList(number int, sample bool) ([]commonmeta.Data, error) {

	var data []commonmeta.Data
	content, err := GetList(number, sample)
	if err != nil {
		return data, err
	}
	for _, v := range content {
		d, err := Read(v)
		if err != nil {
			log.Println(err)
		}
		data = append(data, d)
	}
	return data, nil
}

// Load loads the metadata for a single work from a JSON file
func Load(filename string) (commonmeta.Data, error) {
	var data commonmeta.Data

	content, err := readJSON(filename)
	if err != nil {
		return data, err
	}
	data, err = Read(content)
	if err != nil {
		return data, err
	}
	return data, nil
}

// LoadList loads a list of DataCite metadata from a JSON string and returns Commonmeta metadata.
func LoadList(filename string) ([]commonmeta.Data, error) {
	var data []commonmeta.Data

	response, err := readJSONLines(filename)
	if err != nil {
		return data, err
	}

	data, err = ReadList(response)
	if err != nil {
		return data, err
	}
	return data, nil
}

// Get gets DataCite metadata for a given DOI
func Get(pid string) (Content, error) {
	// the envelope for the JSON response from the DataCite API
	type Response struct {
		Data Content `json:"data"`
	}

	var response Response
	doi, ok := doiutils.ValidateDOI(pid)
	if !ok {
		return response.Data, errors.New("invalid DOI")
	}
	url := "https://api.datacite.org/dois/" + doi
	client := http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(url)
	if err != nil {
		return response.Data, err
	}
	if resp.StatusCode >= 400 {
		return response.Data, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.Data, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("error:", err)
	}
	return response.Data, err
}

// Read reads DataCite JSON response and return work struct in Commonmeta format
func Read(content Content) (commonmeta.Data, error) {
	var data = commonmeta.Data{}
	var err error

	data.ID = doiutils.NormalizeDOI(content.Attributes.DOI)
	data.Type = DCToCMTranslations[content.Attributes.Types.ResourceTypeGeneral]

	// ArchiveLocations not yet supported

	// Support the additional types added in schema 4.4
	AdditionalType := DCToCMTranslations[content.Attributes.Types.ResourceType]
	if AdditionalType != "" {
		data.Type = AdditionalType
	} else if content.Attributes.Types.ResourceType != "" && !strings.EqualFold(content.Attributes.Types.ResourceType, data.Type) {
		data.AdditionalType = content.Attributes.Types.ResourceType
	}

	data.Container = commonmeta.Container{
		Identifier:     content.Attributes.Container.Identifier,
		IdentifierType: content.Attributes.Container.IdentifierType,
		Type:           content.Attributes.Container.Type,
		Title:          content.Attributes.Container.Title,
		Volume:         content.Attributes.Container.Volume,
		Issue:          content.Attributes.Container.Issue,
		FirstPage:      content.Attributes.Container.FirstPage,
		LastPage:       content.Attributes.Container.LastPage,
	}

	for _, v := range content.Attributes.Creators {
		if v.Name != "" || v.GivenName != "" || v.FamilyName != "" {
			contributor := GetContributor(v)
			containsID := slices.ContainsFunc(data.Contributors, func(e commonmeta.Contributor) bool {
				return e.ID != "" && e.ID == contributor.ID
			})
			containsName := slices.ContainsFunc(data.Contributors, func(e commonmeta.Contributor) bool {
				return e.Name != "" && e.Name == contributor.Name || e.GivenName != "" && e.GivenName == contributor.GivenName && e.FamilyName != "" && e.FamilyName == contributor.FamilyName
			})
			if !containsID && !containsName {
				data.Contributors = append(data.Contributors, contributor)
			}
		}
	}

	// merge creators and contributors
	for _, v := range content.Attributes.Contributors {
		if v.Name != "" || v.GivenName != "" || v.FamilyName != "" {
			contributor := GetContributor(v)
			containsID := slices.ContainsFunc(data.Contributors, func(e commonmeta.Contributor) bool {
				return e.ID != "" && e.ID == contributor.ID
			})
			if containsID {
				log.Printf("Contributor with ID %s already exists", contributor.ID)
			} else {
				data.Contributors = append(data.Contributors, contributor)

			}
		}
	}

	for _, v := range content.Attributes.Dates {
		if v.DateType == "Accepted" {
			data.Date.Accepted = v.Date
		}
		if v.DateType == "Available" {
			data.Date.Available = v.Date
		}
		if v.DateType == "Collected" {
			data.Date.Collected = v.Date
		}
		if v.DateType == "Created" {
			data.Date.Created = v.Date
		}
		if v.DateType == "Issued" {
			data.Date.Published = v.Date
		} else if v.DateType == "Published" {
			data.Date.Published = v.Date
		}
		if v.DateType == "Submitted" {
			data.Date.Submitted = v.Date
		}
		if v.DateType == "Updated" {
			data.Date.Updated = v.Date
		}
		if v.DateType == "Valid" {
			data.Date.Valid = v.Date
		}
		if v.DateType == "Withdrawn" {
			data.Date.Withdrawn = v.Date
		}
		if v.DateType == "Other" {
			data.Date.Other = v.Date
		}
	}
	if data.Date.Published == "" {
		data.Date.Published = strconv.Itoa(content.Attributes.PublicationYear)
	}

	for _, v := range content.Attributes.Descriptions {
		var t string
		if slices.Contains([]string{"Abstract", "Summary", "Methods", "TechnicalInfo", "Other"}, v.DescriptionType) {
			t = v.DescriptionType
		} else {
			t = "Other"
		}
		description := utils.Sanitize(v.Description)
		data.Descriptions = append(data.Descriptions, commonmeta.Description{
			Description: description,
			Type:        t,
			Language:    v.Lang,
		})
	}

	// Files not yet supported. Sizes and formats are part of the file object,
	// but can't be mapped directly

	for _, v := range content.Attributes.FundingReferences {
		data.FundingReferences = append(data.FundingReferences, commonmeta.FundingReference{
			FunderIdentifier:     v.FunderIdentifier,
			FunderIdentifierType: v.FunderIdentifierType,
			FunderName:           v.FunderName,
			AwardNumber:          v.AwardNumber,
			AwardURI:             v.AwardURI,
		})
	}

	for _, v := range content.Attributes.GeoLocations {
		if v.GeoLocationPoint.PointLongitude != 0 && v.GeoLocationPoint.PointLatitude != 0 && v.GeoLocationBox.WestBoundLongitude != 0 && v.GeoLocationBox.EastBoundLongitude != 0 && v.GeoLocationBox.SouthBoundLatitude != 0 && v.GeoLocationBox.NorthBoundLatitude != 0 {
			geoLocationPoint := commonmeta.GeoLocationPoint{
				PointLongitude: v.GeoLocationPoint.PointLongitude,
				PointLatitude:  v.GeoLocationPoint.PointLatitude,
			}
			geoLocationBox := commonmeta.GeoLocationBox{
				EastBoundLongitude: v.GeoLocationBox.EastBoundLongitude,
				WestBoundLongitude: v.GeoLocationBox.WestBoundLongitude,
				SouthBoundLatitude: v.GeoLocationBox.SouthBoundLatitude,
				NorthBoundLatitude: v.GeoLocationBox.NorthBoundLatitude,
			}
			geoLocation := commonmeta.GeoLocation{
				GeoLocationPoint: geoLocationPoint,
				GeoLocationPlace: v.GeoLocationPlace,
				GeoLocationBox:   geoLocationBox}
			data.GeoLocations = append(data.GeoLocations, geoLocation)
		}
	}

	if len(content.Attributes.AlternateIdentifiers) > 0 {
		supportedIdentifiers := []string{
			"ARK",
			"arXiv",
			"Bibcode",
			"DOI",
			"Handle",
			"ISBN",
			"ISSN",
			"PMID",
			"PMCID",
			"PURL",
			"URL",
			"URN",
			"Other",
		}
		for _, v := range content.Attributes.AlternateIdentifiers {
			identifierType := "Other"
			if slices.Contains(supportedIdentifiers, v.AlternateIdentifierType) {
				identifierType = v.AlternateIdentifierType
			}
			if v.AlternateIdentifier != "" {
				data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
					Identifier:     v.AlternateIdentifier,
					IdentifierType: identifierType,
				})
			}
		}
	}
	data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
		Identifier:     data.ID,
		IdentifierType: "DOI",
	})
	if len(data.Identifiers) > 1 {
		data.Identifiers = utils.DedupeSlice(data.Identifiers)
	}

	if content.Attributes.Publisher != "" {
		data.Publisher = commonmeta.Publisher{
			Name: content.Attributes.Publisher,
		}
	}

	for _, v := range content.Attributes.Subjects {
		subject := commonmeta.Subject{
			Subject: v.Subject,
		}
		if !slices.Contains(data.Subjects, subject) {
			data.Subjects = append(data.Subjects, subject)
		}
	}

	data.Language = content.Attributes.Language

	if len(content.Attributes.RightsList) > 0 {
		url, _ := utils.NormalizeCCUrl(content.Attributes.RightsList[0].RightsURI)
		id := utils.URLToSPDX(url)
		data.License = commonmeta.License{
			ID:  id,
			URL: url,
		}
	}

	data.Provider = "DataCite"

	if len(content.Attributes.RelatedIdentifiers) > 0 {
		supportedRelations := []string{
			"Cites",
			"References",
		}
		for i, v := range content.Attributes.RelatedIdentifiers {
			id := utils.NormalizeID(v.RelatedIdentifier)
			if id != "" && slices.Contains(supportedRelations, v.RelationType) {
				data.References = append(data.References, commonmeta.Reference{
					Key: "ref" + strconv.Itoa(i+1),
					ID:  id,
				})
			}
		}
	}

	if len(content.Attributes.RelatedIdentifiers) > 0 {
		supportedRelations := []string{
			"IsNewVersionOf",
			"IsPreviousVersionOf",
			"IsVersionOf",
			"HasVersion",
			"IsPartOf",
			"HasPart",
			"IsVariantFormOf",
			"IsOriginalFormOf",
			"IsIdenticalTo",
			"IsTranslationOf",
			"IsReviewedBy",
			"Reviews",
			"IsPreprintOf",
			"HasPreprint",
			"IsSupplementTo",
		}
		for _, v := range content.Attributes.RelatedIdentifiers {
			id := utils.NormalizeID(v.RelatedIdentifier)
			if id != "" && slices.Contains(supportedRelations, v.RelationType) {
				relation := commonmeta.Relation{
					ID:   id,
					Type: v.RelationType,
				}
				if !slices.Contains(data.Relations, relation) {
					data.Relations = append(data.Relations, relation)
				}
			}
		}
	}

	for _, v := range content.Attributes.Titles {
		var t string
		if slices.Contains([]string{"MainTitle", "Subtitle", "TranslatedTitle"}, v.TitleType) {
			t = v.TitleType
		}
		data.Titles = append(data.Titles, commonmeta.Title{
			Title:    v.Title,
			Type:     t,
			Language: v.Lang,
		})
	}

	data.URL, err = utils.NormalizeURL(content.Attributes.URL, true, false)
	if err != nil {
		log.Println(err)
	}

	data.Version = content.Attributes.Version

	return data, nil
}

// GetContributor converts DataCite contributor metadata into the Commonmeta format
func GetContributor(v Contributor) commonmeta.Contributor {
	var t string
	if len(v.NameType) > 2 {
		t = v.NameType[:len(v.NameType)-2]
	}
	var id string
	if len(v.NameIdentifiers) > 0 {
		ni := v.NameIdentifiers[0]
		if ni.NameIdentifierScheme == "ORCID" || ni.NameIdentifierScheme == "https://orcid.org/" {
			id = utils.NormalizeORCID(ni.NameIdentifier)
			t = "Person"
		} else if ni.NameIdentifierScheme == "ROR" {
			id = ni.NameIdentifier
			t = "Organization"
		} else {
			id = ni.NameIdentifier
		}
	}
	name := v.Name
	GivenName := v.GivenName
	FamilyName := v.FamilyName
	if t == "" && (v.GivenName != "" || v.FamilyName != "") {
		t = "Person"
	} else if t == "" {
		t = "Organization"
	}
	if t == "Person" && name != "" {
		// split name for type Person into given/family name if not already provided
		names := strings.Split(name, ",")
		if len(names) == 2 {
			GivenName = strings.TrimSpace(names[1])
			FamilyName = names[0]
			name = ""
		}
	}
	var affiliations []commonmeta.Affiliation
	for _, a := range v.Affiliation {
		id := utils.NormalizeROR(a.AffiliationIdentifier)
		affiliations = append(affiliations, commonmeta.Affiliation{
			ID:   id,
			Name: a.Name,
		})
	}
	var roles []string
	if slices.Contains(commonmeta.ContributorRoles, v.ContributorType) {
		roles = append(roles, v.ContributorType)
	} else {
		roles = append(roles, "Author")
	}
	return commonmeta.Contributor{
		ID:               id,
		Type:             t,
		Name:             name,
		GivenName:        GivenName,
		FamilyName:       FamilyName,
		ContributorRoles: roles,
		Affiliations:     affiliations,
	}
}

// GetList gets the metadata for a list of works from the DataCite API
func GetList(number int, sample bool) ([]Content, error) {
	// the envelope for the JSON response from the DataCite API
	type Response struct {
		Data []Content `json:"data"`
	}
	if number > 100 {
		number = 100
	}
	var response Response
	url := QueryURL(number, sample)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	client := http.Client{
		Timeout: 30 * time.Second,
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
	return response.Data, nil
}

// ReadList reads a list of DataCite JSON responses and returns a list of works in Commonmeta format
func ReadList(content []Content) ([]commonmeta.Data, error) {
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

// QueryURL returns the URL for the DataCite API query
func QueryURL(number int, sample bool) string {
	if sample {
		number = 10
	}
	url := "https://api.datacite.org/dois?random=true&page[size]=" + strconv.Itoa(number)
	return url
}

// readJSON reads JSON from a file and unmarshals it
func readJSON(filename string) (Content, error) {
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
	err = decoder.Decode(&content.Attributes)
	if err != nil {
		return content, err
	}
	return content, nil
}

// readJSONLines reads JSON lines from a file and unmarshals them
func readJSONLines(filename string) ([]Content, error) {
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
		var attributes Attributes
		if err := decoder.Decode(&attributes); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		response = append(response, Content{Attributes: attributes})
	}

	return response, nil
}
