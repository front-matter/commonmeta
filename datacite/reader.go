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

// Datacite represents the DataCite metadata.
type Datacite struct {
	ID                   string                `json:"id"`
	DOI                  string                `json:"doi"`
	AlternateIdentifiers []AlternateIdentifier `json:"alternateIdentifiers,omitempty"`
	Creators             []Contributor         `json:"creators"`
	Publisher            Publisher             `json:"publisher"`
	Container            Container             `json:"container,omitempty"`
	PublicationYear      int                   `json:"publicationYear"`
	Titles               []Title               `json:"titles"`
	URL                  string                `json:"url"`
	Subjects             []Subject             `json:"subjects,omitempty"`
	Contributors         []Contributor         `json:"contributors,omitempty"`
	Dates                []Date                `json:"dates,omitempty"`
	Language             string                `json:"language,omitempty"`
	Types                Types                 `json:"types"`
	RelatedIdentifiers   []RelatedIdentifier   `json:"relatedIdentifiers,omitempty"`
	Sizes                []string              `json:"sizes,omitempty"`
	Formats              []string              `json:"formats,omitempty"`
	Version              string                `json:"version,omitempty"`
	RightsList           []Rights              `json:"rightsList,omitempty"`
	Descriptions         []Description         `json:"descriptions,omitempty"`
	GeoLocations         []GeoLocation         `json:"geoLocations,omitempty"`
	FundingReferences    []FundingReference    `json:"fundingReferences,omitempty"`
	SchemaVersion        string                `json:"schemaVersion"`
}

// Content represents the DataCite metadata returned from DataCite. The type is more
// flexible than the Datacite type, allowing for different formats of some metadata.
// Affiliation can be string or struct, PublicationYear can be int or string. Publisher can be string or struct.
type Content struct {
	*Datacite
	Creators        []ContentContributor `json:"creators"`
	Contributors    []ContentContributor `json:"contributors"`
	PublicationYear json.RawMessage      `json:"publicationYear"`
	Publisher       json.RawMessage      `json:"publisher"`
}

// ContentContributor represents a creator or contributor in the DataCite JSONAPI response.
type ContentContributor struct {
	*Contributor
	Affiliation json.RawMessage `json:"affiliation,omitempty"`
}

type Affiliation struct {
	AffiliationIdentifier       string `json:"affiliationIdentifier,omitempty"`
	AffiliationIdentifierScheme string `json:"affiliationIdentifierScheme,omitempty"`
	SchemeURI                   string `json:"schemeUri,omitempty"`
	Name                        string `json:"name"`
}

// AlternateIdentifier represents an alternate identifier in the DataCite metadata.
type AlternateIdentifier struct {
	AlternateIdentifier     string `json:"alternateIdentifier,omitempty"`
	AlternateIdentifierType string `json:"alternateIdentifierType,omitempty"`
}

// Container represents the container of the DataCite JSONAPI response.
type Container struct {
	Type           string `json:"type,omitempty"`
	Identifier     string `json:"identifier,omitempty"`
	IdentifierType string `json:"identifierType,omitempty"`
	Title          string `json:"title,omitempty"`
	Volume         string `json:"volume,omitempty"`
	Issue          string `json:"issue,omitempty"`
	FirstPage      string `json:"firstPage,omitempty"`
	LastPage       string `json:"lastPage,omitempty"`
}

// Contributor represents the contributor of the DataCite JSONAPI response.
type Contributor struct {
	Name            string           `json:"name,omitempty"`
	GivenName       string           `json:"givenName,omitempty"`
	FamilyName      string           `json:"familyName,omitempty"`
	NameType        string           `json:"nameType"`
	Affiliation     []string         `json:"affiliation,omitempty"`
	NameIdentifiers []NameIdentifier `json:"nameIdentifiers,omitempty"`
	ContributorType string           `json:"contributorType,omitempty"`
}

type Date struct {
	Date            string `json:"date,omitempty"`
	DateType        string `json:"dateType,omitempty"`
	DateInformation string `json:"dateInformation,omitempty"`
}

type Description struct {
	Description     string `json:"description,omitempty"`
	DescriptionType string `json:"descriptionType,omitempty"`
	Lang            string `json:"lang,omitempty"`
}

type FundingReference struct {
	FunderName           string `json:"funderName,omitempty"`
	FunderIdentifier     string `json:"funderIdentifier,omitempty"`
	FunderIdentifierType string `json:"funderIdentifierType,omitempty"`
	AwardNumber          string `json:"awardNumber,omitempty"`
	AwardURI             string `json:"awardUri,omitempty"`
}

type GeoLocation struct {
	GeoLocationPoint `json:"geoLocationPoint,omitempty"`
	GeoLocationBox   `json:"geoLocationBox,omitempty"`
	GeoLocationPlace string `json:"geoLocationPlace,omitempty"`
}

type GeoLocationBox struct {
	WestBoundLongitude float64 `json:"westBoundLongitude,string,omitempty"`
	EastBoundLongitude float64 `json:"eastBoundLongitude,string,omitempty"`
	SouthBoundLatitude float64 `json:"southBoundLatitude,string,omitempty"`
	NorthBoundLatitude float64 `json:"northBoundLatitude,string,omitempty"`
}

type GeoLocationPoint struct {
	PointLongitude float64 `json:"pointLongitude,string,omitempty"`
	PointLatitude  float64 `json:"pointLatitude,string,omitempty"`
}

type NameIdentifier struct {
	NameIdentifier       string `json:"nameIdentifier,omitempty"`
	NameIdentifierScheme string `json:"nameIdentifierScheme,omitempty"`
	SchemeURI            string `json:"schemeUri,omitempty"`
}

type RelatedIdentifier struct {
	RelatedIdentifier     string `json:"relatedIdentifier,omitempty"`
	RelatedIdentifierType string `json:"relatedIdentifierType,omitempty"`
	RelationType          string `json:"relationType,omitempty"`
}

type Rights struct {
	Rights                 string `json:"rights,omitempty"`
	RightsURI              string `json:"rightsUri,omitempty"`
	SchemeURI              string `json:"schemeUri,omitempty"`
	RightsIdentifier       string `json:"rightsIdentifier,omitempty"`
	RightsIdentifierScheme string `json:"rightsIdentifierScheme,omitempty"`
}

type Subject struct {
	Subject string `json:"subject,omitempty"`
}

type Title struct {
	Title     string `json:"title"`
	TitleType string `json:"titleType,omitempty"`
	Lang      string `json:"lang,omitempty"`
}

type Publisher struct {
	Name                      string `json:"name"`
	PublisherIdentifier       string `json:"publisherIdentifier,omitempty"`
	PublisherIdentifierScheme string `json:"publisherIdentifierScheme,omitempty"`
	SchemeURI                 string `json:"schemeUri,omitempty"`
	Lang                      string `json:"lang,omitempty"`
}

type Types struct {
	ResourceTypeGeneral string `json:"resourceTypeGeneral"`
	ResourceType        string `json:"resourceType,omitempty"`
	Ris                 string `json:"ris,omitempty"`
	Bibtex              string `json:"bibtex,omitempty"`
	Citeproc            string `json:"citeproc,omitempty"`
	SchemaOrg           string `json:"schemaOrg,omitempty"`
}

// DCToCMMappings maps DataCite resource types to Commonmeta types
// source: https://github.com/datacite/schema/blob/master/source/meta/kernel-4/include/datacite-resourceType-v4.xsd
var DCToCMMappings = map[string]string{
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

var CMToDCMappings = map[string]string{
	"Article":               "Preprint",
	"Audiovisual":           "Audiovisual",
	"Book":                  "Book",
	"BookChapter":           "BookChapter",
	"Collection":            "Collection",
	"Dataset":               "Dataset",
	"Document":              "Text",
	"Entry":                 "Text",
	"Event":                 "Event",
	"Figure":                "Image",
	"Image":                 "Image",
	"Instrument":            "Instrument",
	"JournalArticle":        "JournalArticle",
	"LegalDocument":         "Text",
	"Manuscript":            "Text",
	"Map":                   "Image",
	"Patent":                "Text",
	"Performance":           "Audiovisual",
	"PersonalCommunication": "Text",
	"Post":                  "Text",
	"ProceedingsArticle":    "ConferencePaper",
	"Proceedings":           "ConferenceProceeding",
	"Report":                "Report",
	"Review":                "PeerReview",
	"Software":              "Software",
	"Sound":                 "Sound",
	"Standard":              "Standard",
	"StudyRegistration":     "StudyRegistration",
	"WebPage":               "Text",
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

	content, err := ReadJSON(filename)
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

	response, err := ReadJSONLines(filename)
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
func Get(id string) (Content, error) {
	// the envelope for the JSON response from the DataCite API
	type Response struct {
		Data struct {
			ID         string  `json:"id"`
			Attributes Content `json:"attributes"`
		} `json:"data"`
	}

	var response Response
	doi, ok := doiutils.ValidateDOI(id)
	if !ok {
		return response.Data.Attributes, errors.New("invalid DOI")
	}
	url := "https://api.datacite.org/dois/" + doi
	client := http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(url)
	if err != nil {
		return response.Data.Attributes, err
	}
	if resp.StatusCode >= 400 {
		return response.Data.Attributes, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.Data.Attributes, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("error:", err)
	}
	return response.Data.Attributes, err
}

// Read reads DataCite JSON response and return work struct in Commonmeta format
func Read(content Content) (commonmeta.Data, error) {
	var data = commonmeta.Data{}
	var err error

	data.ID = doiutils.NormalizeDOI(content.DOI)
	data.Type = DCToCMMappings[content.Types.ResourceTypeGeneral]

	// ArchiveLocations not yet supported

	// Support the additional types added in schema 4.4
	AdditionalType := DCToCMMappings[content.Types.ResourceType]
	if AdditionalType != "" {
		data.Type = AdditionalType
	} else if content.Types.ResourceType != "" && !strings.EqualFold(content.Types.ResourceType, data.Type) {
		data.AdditionalType = content.Types.ResourceType
	}

	data.Container = commonmeta.Container{
		Identifier:     content.Container.Identifier,
		IdentifierType: content.Container.IdentifierType,
		Type:           content.Container.Type,
		Title:          content.Container.Title,
		Volume:         content.Container.Volume,
		Issue:          content.Container.Issue,
		FirstPage:      content.Container.FirstPage,
		LastPage:       content.Container.LastPage,
	}

	for _, v := range content.Creators {
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

	// merge creators and contributors
	for _, v := range content.Contributors {
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

	for _, v := range content.Dates {
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
		data.Date.Published = string(content.PublicationYear)
	}

	for _, v := range content.Descriptions {
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

	for _, v := range content.FundingReferences {
		data.FundingReferences = append(data.FundingReferences, commonmeta.FundingReference{
			FunderIdentifier:     v.FunderIdentifier,
			FunderIdentifierType: v.FunderIdentifierType,
			FunderName:           v.FunderName,
			AwardNumber:          v.AwardNumber,
			AwardURI:             v.AwardURI,
		})
	}
	for _, v := range content.GeoLocations {
		geoLocation := commonmeta.GeoLocation{
			GeoLocationPlace: v.GeoLocationPlace,
			GeoLocationPoint: commonmeta.GeoLocationPoint{
				PointLongitude: v.GeoLocationPoint.PointLongitude,
				PointLatitude:  v.GeoLocationPoint.PointLatitude,
			},
			GeoLocationBox: commonmeta.GeoLocationBox{
				WestBoundLongitude: v.GeoLocationBox.WestBoundLongitude,
				EastBoundLongitude: v.GeoLocationBox.EastBoundLongitude,
				SouthBoundLatitude: v.GeoLocationBox.SouthBoundLatitude,
				NorthBoundLatitude: v.GeoLocationBox.NorthBoundLatitude,
			},
		}
		data.GeoLocations = append(data.GeoLocations, geoLocation)
	}

	if len(content.AlternateIdentifiers) > 0 {
		for _, v := range content.AlternateIdentifiers {
			identifierType := "Other"
			if slices.Contains(commonmeta.IdentifierTypes, v.AlternateIdentifierType) {
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

	// parse Publisher as either string (up to schema 4.4) or struct (schema 4.5)
	var publisher Publisher
	var publisherName string
	err = json.Unmarshal(content.Publisher, &publisher)
	if err != nil {
		err = json.Unmarshal(content.Publisher, &publisherName)
	}
	if err != nil {
		log.Println(err)
	}
	if publisher.Name != "" {
		id := utils.NormalizeROR(publisher.PublisherIdentifier)
		data.Publisher = commonmeta.Publisher{
			ID:   id,
			Name: publisher.Name,
		}
	} else if publisherName != "" {
		data.Publisher = commonmeta.Publisher{
			Name: publisherName,
		}
	}

	for _, v := range content.Subjects {
		subject := commonmeta.Subject{
			Subject: v.Subject,
		}
		if !slices.Contains(data.Subjects, subject) {
			data.Subjects = append(data.Subjects, subject)
		}
	}

	data.Language = content.Language

	if len(content.RightsList) > 0 {
		url, _ := utils.NormalizeCCUrl(content.RightsList[0].RightsURI)
		id := utils.URLToSPDX(url)
		data.License = commonmeta.License{
			ID:  id,
			URL: url,
		}
	}

	data.Provider = "DataCite"

	if len(content.RelatedIdentifiers) > 0 {
		supportedRelations := []string{
			"Cites",
			"References",
		}
		for i, v := range content.RelatedIdentifiers {
			id := utils.NormalizeID(v.RelatedIdentifier)
			if id != "" && slices.Contains(supportedRelations, v.RelationType) {
				data.References = append(data.References, commonmeta.Reference{
					Key: "ref" + strconv.Itoa(i+1),
					ID:  id,
				})
			}
		}
	}

	if len(content.RelatedIdentifiers) > 0 {
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
		for _, v := range content.RelatedIdentifiers {
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

	for _, v := range content.Titles {
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

	data.URL, err = utils.NormalizeURL(content.URL, true, false)
	if err != nil {
		log.Println(err)
	}

	data.Version = content.Version

	return data, nil
}

// GetContributor converts DataCite contributor metadata into the Commonmeta format
func GetContributor(v ContentContributor) commonmeta.Contributor {
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

	//parse Affiliation as either slice of string or slice of struct
	var affiliationStructs []Affiliation
	var affiliationNames []string
	var affiliations []commonmeta.Affiliation
	var err error
	err = json.Unmarshal(v.Affiliation, &affiliationNames)
	if err != nil {
		err = json.Unmarshal(v.Affiliation, &affiliationStructs)
	}
	if err != nil {
		log.Println(err)
	}
	if len(affiliationStructs) > 0 {
		for _, v := range affiliationStructs {
			id := utils.NormalizeROR(v.AffiliationIdentifier)
			af := commonmeta.Affiliation{
				ID:   id,
				Name: v.Name,
			}
			affiliations = append(affiliations, af)
		}
	} else if len(affiliationNames) > 0 {
		af := commonmeta.Affiliation{
			Name: affiliationNames[0],
		}
		affiliations = append(affiliations, af)
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
		Affiliations:     affiliations,
		ContributorRoles: roles,
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
	if err != nil {
		return content, err
	}
	return content, nil
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
		var datacite Content
		if err := decoder.Decode(&datacite); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		response = append(response, datacite)
	}

	return response, nil
}
