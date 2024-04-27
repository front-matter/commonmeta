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

	"github.com/front-matter/commonmeta/bibtex"
	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/csl"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/ris"
	"github.com/front-matter/commonmeta/schemaorg"
	"github.com/front-matter/commonmeta/schemautils"
	"github.com/xeipuuv/gojsonschema"

	"github.com/front-matter/commonmeta/utils"
)

// Content represents the DataCite JSONAPI response.
type Content struct {
	ID         string   `json:"id"`
	Type       string   `json:"type"`
	Attributes Datacite `json:"attributes"`
}

// Datacite represents the DataCite metadata.
type Datacite struct {
	ID                   string                `json:"id"`
	DOI                  string                `json:"doi"`
	AlternateIdentifiers []AlternateIdentifier `json:"alternateIdentifiers,omitempty"`
	Creators             []Contributor         `json:"creators"`
	Publisher            string                `json:"publisher"`
	Container            Container             `json:"container,omitempty"`
	PublicationYear      int                   `json:"publicationYear"`
	Titles               []Title               `json:"titles"`
	URL                  string                `json:"url"`
	Subjects             []Subject             `json:"subjects,omitempty"`
	Contributors         []Contributor         `json:"contributors,omitempty"`
	Dates                []struct {
		Date            string `json:"date,omitempty"`
		DateType        string `json:"dateType,omitempty"`
		DateInformation string `json:"dateInformation,omitempty"`
	} `json:"dates,omitempty"`
	Language           string              `json:"language,omitempty"`
	Types              Types               `json:"types"`
	RelatedIdentifiers []RelatedIdentifier `json:"relatedIdentifiers,omitempty"`
	Sizes              []string            `json:"sizes,omitempty"`
	Formats            []string            `json:"formats,omitempty"`
	Version            string              `json:"version,omitempty"`
	RightsList         []Rights            `json:"rightsList,omitempty"`
	Descriptions       []Description       `json:"descriptions,omitempty"`
	GeoLocations       []GeoLocation       `json:"geoLocations,omitempty"`
	FundingReferences  []FundingReference  `json:"fundingReferences,omitempty"`
	SchemaVersion      string              `json:"schemaVersion"`
}

type Affiliation struct {
	Name string `json:"name"`
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
	NameIdentifiers []NameIdentifier `json:"nameIdentifiers,omitempty"`
	Affiliation     []Affiliation    `json:"affiliation,omitempty"`
	ContributorType string           `json:"contributorType,omitempty"`
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
	log.Println(data)
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
func Get(id string) (Datacite, error) {
	// the envelope for the JSON response from the DataCite API
	type Response struct {
		Data Datacite `json:"data"`
	}

	var response Response
	doi, ok := doiutils.ValidateDOI(id)
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
func Read(datacite Datacite) (commonmeta.Data, error) {
	var data = commonmeta.Data{}
	var err error

	data.ID = doiutils.NormalizeDOI(datacite.DOI)
	data.Type = DCToCMMappings[datacite.Types.ResourceTypeGeneral]

	// ArchiveLocations not yet supported

	// Support the additional types added in schema 4.4
	AdditionalType := DCToCMMappings[datacite.Types.ResourceType]
	if AdditionalType != "" {
		data.Type = AdditionalType
	} else if datacite.Types.ResourceType != "" && !strings.EqualFold(datacite.Types.ResourceType, data.Type) {
		data.AdditionalType = datacite.Types.ResourceType
	}

	data.Container = commonmeta.Container{
		Identifier:     datacite.Container.Identifier,
		IdentifierType: datacite.Container.IdentifierType,
		Type:           datacite.Container.Type,
		Title:          datacite.Container.Title,
		Volume:         datacite.Container.Volume,
		Issue:          datacite.Container.Issue,
		FirstPage:      datacite.Container.FirstPage,
		LastPage:       datacite.Container.LastPage,
	}

	for _, v := range datacite.Creators {
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
	for _, v := range datacite.Contributors {
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

	for _, v := range datacite.Dates {
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
		data.Date.Published = strconv.Itoa(datacite.PublicationYear)
	}

	for _, v := range datacite.Descriptions {
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

	for _, v := range datacite.FundingReferences {
		data.FundingReferences = append(data.FundingReferences, commonmeta.FundingReference{
			FunderIdentifier:     v.FunderIdentifier,
			FunderIdentifierType: v.FunderIdentifierType,
			FunderName:           v.FunderName,
			AwardNumber:          v.AwardNumber,
			AwardURI:             v.AwardURI,
		})
	}
	for _, v := range datacite.GeoLocations {
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

	if len(datacite.AlternateIdentifiers) > 0 {
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
		for _, v := range datacite.AlternateIdentifiers {
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

	if datacite.Publisher != "" {
		data.Publisher = commonmeta.Publisher{
			Name: datacite.Publisher,
		}
	}

	for _, v := range datacite.Subjects {
		subject := commonmeta.Subject{
			Subject: v.Subject,
		}
		if !slices.Contains(data.Subjects, subject) {
			data.Subjects = append(data.Subjects, subject)
		}
	}

	data.Language = datacite.Language

	if len(datacite.RightsList) > 0 {
		url, _ := utils.NormalizeCCUrl(datacite.RightsList[0].RightsURI)
		id := utils.URLToSPDX(url)
		data.License = commonmeta.License{
			ID:  id,
			URL: url,
		}
	}

	data.Provider = "DataCite"

	if len(datacite.RelatedIdentifiers) > 0 {
		supportedRelations := []string{
			"Cites",
			"References",
		}
		for i, v := range datacite.RelatedIdentifiers {
			id := utils.NormalizeID(v.RelatedIdentifier)
			if id != "" && slices.Contains(supportedRelations, v.RelationType) {
				data.References = append(data.References, commonmeta.Reference{
					Key: "ref" + strconv.Itoa(i+1),
					ID:  id,
				})
			}
		}
	}

	if len(datacite.RelatedIdentifiers) > 0 {
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
		for _, v := range datacite.RelatedIdentifiers {
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

	for _, v := range datacite.Titles {
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

	data.URL, err = utils.NormalizeURL(datacite.URL, true, false)
	if err != nil {
		log.Println(err)
	}

	data.Version = datacite.Version

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
		affiliations = append(affiliations, commonmeta.Affiliation{
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
func GetList(number int, sample bool) ([]Datacite, error) {
	// the envelope for the JSON response from the DataCite API
	type Response struct {
		Data []Datacite `json:"data"`
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
func ReadList(content []Datacite) ([]commonmeta.Data, error) {
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
func readJSON(filename string) (Datacite, error) {
	var content Datacite

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

// readJSONLines reads JSON lines from a file and unmarshals them
func readJSONLines(filename string) ([]Datacite, error) {
	var response []Datacite

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
		var datacite Datacite
		if err := decoder.Decode(&datacite); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		response = append(response, datacite)
	}

	return response, nil
}

// Convert converts Commonmeta metadata to DataCite metadata
func Convert(data commonmeta.Data) (Datacite, error) {
	var datacite Datacite

	// required properties
	datacite.ID = data.ID
	datacite.DOI, _ = doiutils.ValidateDOI(data.ID)
	datacite.Types.ResourceTypeGeneral = CMToDCMappings[data.Type]
	datacite.Types.SchemaOrg = schemaorg.CMToSOMappings[data.Type]
	datacite.Types.Citeproc = csl.CMToCSLMappings[data.Type]
	datacite.Types.Bibtex = bibtex.CMToBibMappings[data.Type]
	datacite.Types.Ris = ris.CMToRISMappings[data.Type]
	if data.AdditionalType != "" {
		datacite.Types.ResourceType = data.AdditionalType
	}
	if datacite.Types.ResourceTypeGeneral == "" {
		datacite.Types.ResourceTypeGeneral = "Other"
	}

	if len(data.Date.Published) >= 4 {
		datacite.PublicationYear, _ = strconv.Atoi(data.Date.Published[:4])
	}

	if len(data.Titles) > 0 {
		for _, v := range data.Titles {
			title := Title{
				Title:     v.Title,
				TitleType: v.Type,
				Lang:      v.Language,
			}
			datacite.Titles = append(datacite.Titles, title)
		}
	}

	if len(data.Contributors) > 0 {
		for _, v := range data.Contributors {
			if slices.Contains(commonmeta.ContributorRoles, "Author") {
				var nameIdentifiers []NameIdentifier
				if v.ID != "" {
					nameIdentifier := NameIdentifier{
						NameIdentifier:       v.ID,
						NameIdentifierScheme: "ORCID",
						SchemeURI:            "https://orcid.org",
					}
					nameIdentifiers = append(nameIdentifiers, nameIdentifier)
				}
				var affiliations []Affiliation
				for _, a := range v.Affiliations {
					affiliation := Affiliation{
						Name: a.Name,
					}
					affiliations = append(affiliations, affiliation)
				}
				contributor := Contributor{
					Name:            v.Name,
					GivenName:       v.GivenName,
					FamilyName:      v.FamilyName,
					NameType:        v.Type + "al",
					NameIdentifiers: nameIdentifiers,
					Affiliation:     affiliations,
				}
				datacite.Creators = append(datacite.Creators, contributor)
			}
		}
	}

	datacite.Publisher = data.Publisher.Name
	datacite.URL = data.URL
	datacite.SchemaVersion = "http://datacite.org/schema/kernel-4"

	// optional properties

	datacite.Container = Container{
		Type:           data.Container.Type,
		Identifier:     data.Container.Identifier,
		IdentifierType: data.Container.IdentifierType,
		Title:          data.Container.Title,
		Volume:         data.Container.Volume,
		Issue:          data.Container.Issue,
		FirstPage:      data.Container.FirstPage,
		LastPage:       data.Container.LastPage,
	}

	if len(data.Identifiers) > 0 {
		for _, v := range data.Identifiers {
			if v.Identifier != data.ID {
				AlternateIdentifier := AlternateIdentifier{
					AlternateIdentifier:     v.Identifier,
					AlternateIdentifierType: v.IdentifierType,
				}
				datacite.AlternateIdentifiers = append(datacite.AlternateIdentifiers, AlternateIdentifier)
			}
		}
	}

	if len(data.Descriptions) > 0 {
		for _, v := range data.Descriptions {
			description := Description{
				Description:     v.Description,
				DescriptionType: v.Type,
				Lang:            v.Language,
			}
			datacite.Descriptions = append(datacite.Descriptions, description)
		}
	}

	if len(data.FundingReferences) > 0 {
		for _, v := range data.FundingReferences {
			fundingReference := FundingReference{
				FunderName:           v.FunderName,
				FunderIdentifier:     v.FunderIdentifier,
				FunderIdentifierType: v.FunderIdentifierType,
				AwardNumber:          v.AwardNumber,
				AwardURI:             v.AwardURI,
			}
			datacite.FundingReferences = append(datacite.FundingReferences, fundingReference)
		}
	}
	if len(data.GeoLocations) > 0 {
		for _, v := range data.GeoLocations {
			geoLocation := GeoLocation{
				GeoLocationPlace: v.GeoLocationPlace,
				GeoLocationPoint: GeoLocationPoint{
					PointLongitude: v.GeoLocationPoint.PointLongitude,
					PointLatitude:  v.GeoLocationPoint.PointLatitude,
				},
				GeoLocationBox: GeoLocationBox{
					WestBoundLongitude: v.GeoLocationBox.WestBoundLongitude,
					EastBoundLongitude: v.GeoLocationBox.EastBoundLongitude,
					SouthBoundLatitude: v.GeoLocationBox.SouthBoundLatitude,
					NorthBoundLatitude: v.GeoLocationBox.NorthBoundLatitude,
				},
			}
			datacite.GeoLocations = append(datacite.GeoLocations, geoLocation)
		}
	}
	datacite.Language = data.Language
	if len(data.Subjects) > 0 {
		for _, v := range data.Subjects {
			subject := Subject{Subject: v.Subject}
			datacite.Subjects = append(datacite.Subjects, subject)
		}
	}
	if data.License.URL != "" {
		rights := Rights{
			RightsURI:              data.License.URL,
			RightsIdentifier:       data.License.ID,
			RightsIdentifierScheme: "SPDX",
			SchemeURI:              "https://spdx.org/licenses/",
		}
		datacite.RightsList = append(datacite.RightsList, rights)
	}
	if len(data.Relations) > 0 {
		for _, v := range data.Relations {
			id := doiutils.NormalizeDOI(v.ID)
			relatedIdentifierType := "DOI"
			if id == "" {
				relatedIdentifierType = "URL"
			}
			RelatedIdentifier := RelatedIdentifier{
				RelatedIdentifier:     id,
				RelatedIdentifierType: relatedIdentifierType,
				RelationType:          v.Type,
			}
			datacite.RelatedIdentifiers = append(datacite.RelatedIdentifiers, RelatedIdentifier)
		}
	}

	if len(data.References) > 0 {
		for _, v := range data.References {
			id := doiutils.NormalizeDOI(v.ID)
			relatedIdentifierType := "DOI"
			if id == "" {
				relatedIdentifierType = "URL"
			}
			RelatedIdentifier := RelatedIdentifier{
				RelatedIdentifier:     id,
				RelatedIdentifierType: relatedIdentifierType,
				RelationType:          "References",
			}
			datacite.RelatedIdentifiers = append(datacite.RelatedIdentifiers, RelatedIdentifier)
		}
	}

	if len(data.Relations) > 0 {
		for _, v := range data.Relations {
			RelatedIdentifier := RelatedIdentifier{
				RelatedIdentifier:     v.ID,
				RelatedIdentifierType: "DOI",
				RelationType:          v.Type,
			}
			datacite.RelatedIdentifiers = append(datacite.RelatedIdentifiers, RelatedIdentifier)
		}
	}

	datacite.Version = data.Version

	// "creators": creators,
	// "titles": metadata.titles,
	// "contributors": contributors,
	// "dates": dates,

	// "relatedIdentifiers": related_identifiers,

	return datacite, nil
}

// Write writes commonmeta metadata.
func Write(data commonmeta.Data) ([]byte, []gojsonschema.ResultError) {
	datacite, err := Convert(data)
	if err != nil {
		fmt.Println(err)
	}
	output, err := json.Marshal(datacite)
	if err != nil {
		fmt.Println(err)
	}
	validation := schemautils.JSONSchemaErrors(output, "datacite-v4.5")
	if !validation.Valid() {
		return nil, validation.Errors()
	}

	return output, nil
}

// WriteList writes a list of commonmeta metadata.
func WriteList(list []commonmeta.Data) ([]byte, []gojsonschema.ResultError) {
	var dataciteList []Datacite
	for _, data := range list {
		datacite, err := Convert(data)
		if err != nil {
			fmt.Println(err)
		}
		dataciteList = append(dataciteList, datacite)
	}
	output, err := json.Marshal(dataciteList)
	if err != nil {
		fmt.Println(err)
	}
	validation := schemautils.JSONSchemaErrors(output, "datacite-v4.5")
	if !validation.Valid() {
		return nil, validation.Errors()
	}

	return output, nil
}
