// Package openalex provides functions to convert OpenAlex metadata to the commonmeta metadata format.
package openalex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/front-matter/commonmeta/authorutils"
	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/utils"
)

const openAlexBaseURL = "https://api.openalex.org"

// Reader struct to hold any configuration for the OpenAlex reader
type Reader struct {
	Email string // Email for polite pool of OpenAlex API
}

// NewReader creates a new OpenAlex reader
func NewReader(email string) *Reader {
	return &Reader{
		Email: email,
	}
}

// OAToCMMappings maps OpenAlex types to Commonmeta types
var OAToCMMappings = map[string]string{
	"article":                 "Article",
	"book":                    "Book",
	"book-chapter":            "BookChapter",
	"dataset":                 "Dataset",
	"dissertation":            "Dissertation",
	"editorial":               "Document",
	"erratum":                 "Other",
	"grant":                   "Grant",
	"letter":                  "Article",
	"libguides":               "InteractiveResource",
	"other":                   "Other",
	"paratext":                "Component",
	"peer-review":             "PeerReview",
	"preprint":                "Article",
	"reference-entry":         "Other",
	"report":                  "Report",
	"retraction":              "Other",
	"review":                  "Article",
	"standard":                "Standard",
	"supplementary-materials": "Component",
}

// OpenAlexContainerTypes maps OpenAlex container types to Commonmeta container types
var OpenAlexContainerTypes = map[string]string{
	"journal":       "Journal",
	"proceedings":   "Proceedings",
	"reference":     "Collection",
	"repository":    "Repository",
	"book-series":   "BookSeries",
	"book":          "Book",
	"report-series": "ReportSeries",
}

// OpenAlexIdentifierTypes maps OpenAlex identifier types to Commonmeta identifier types
var OpenAlexIdentifierTypes = map[string]string{
	"openalex": "OpenAlex",
	"doi":      "DOI",
	"mag":      "MAG",
	"pmid":     "PMID",
	"pmcid":    "PMCID",
}

// OpenAlexLicenses maps OpenAlex license strings to SPDX licenseId
var OpenAlexLicenses = map[string]string{
	"cc-by": "CC-BY-4.0",
	"cc0":   "CC0-1.0",
}

// Work represents an OpenAlex work
type Work struct {
	ID                    string            `json:"id"`
	DOI                   string            `json:"doi"`
	Title                 string            `json:"title"`
	DisplayName           string            `json:"display_name"`
	Type                  string            `json:"type"`
	TypeCrossref          string            `json:"type_crossref"`
	PublicationDate       string            `json:"publication_date"`
	CreatedDate           string            `json:"created_date"`
	Language              string            `json:"language"`
	Version               string            `json:"version"`
	AbstractInvertedIndex map[string][]int  `json:"abstract_inverted_index"`
	AuthorShips           []AuthorShip      `json:"authorships"`
	Ids                   map[string]string `json:"ids"`
	PrimaryLocation       Location          `json:"primary_location"`
	BestOALocation        Location          `json:"best_oa_location"`
	Topics                []Topic           `json:"topics"`
	Biblio                Biblio            `json:"biblio"`
	ReferencedWorks       []string          `json:"referenced_works"`
	RelatedWorks          []string          `json:"related_works"`
	Grants                []Grant           `json:"grants"`
}

// AuthorShip represents author information in OpenAlex
type AuthorShip struct {
	Author       Author        `json:"author"`
	Institutions []Institution `json:"institutions"`
}

// Author represents an author in OpenAlex
type Author struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	ORCID       string `json:"orcid"`
}

// Biblio represents bibliographic information in OpenAlex
type Biblio struct {
	Volume    string `json:"volume"`
	Issue     string `json:"issue"`
	FirstPage string `json:"first_page"`
	LastPage  string `json:"last_page"`
}

// Funder represents funder information in OpenAlex
type Funder struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Ids         struct {
		ROR string `json:"ror"`
	} `json:"ids"`
}

// Grant represents grant information in OpenAlex
type Grant struct {
	Funder  string `json:"funder"`
	AwardID string `json:"award_id"`
}

// Institution represents an institution in OpenAlex
type Institution struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	ROR         string `json:"ror"`
}

// Location represents a location in OpenAlex
type Location struct {
	LandingPageURL string `json:"landing_page_url"`
	PDFURL         string `json:"pdf_url"`
	Source         Source `json:"source"`
	License        string `json:"license"`
}

// Source represents a source in OpenAlex
type Source struct {
	ID                   string `json:"id"`
	DisplayName          string `json:"display_name"`
	ISSN                 string `json:"issn_l"`
	Type                 string `json:"type"`
	HostOrganizationName string `json:"host_organization_name"`
	HomepageURL          string `json:"homepage_url"`
}

// SubfieldTopic represents a subfield topic in OpenAlex
type SubfieldTopic struct {
	DisplayName string `json:"display_name"`
}

// Topic represents a topic in OpenAlex
type Topic struct {
	ID       string        `json:"id"`
	Subfield SubfieldTopic `json:"subfield"`
}

// QueryURL constructs a URL for querying the OpenAlex API
func (r *Reader) QueryURL(query url.Values) string {
	u, err := url.Parse(openAlexBaseURL)
	if err != nil {
		return ""
	}
	u.Path = path.Join(u.Path, "works")
	if r.Email != "" {
		query.Add("mailto", r.Email)
	}
	u.RawQuery = query.Encode()
	return u.String()
}

// SampleURL constructs a URL for getting a random sample from OpenAlex API
func (r *Reader) SampleURL(count int) string {
	u, err := url.Parse(openAlexBaseURL)
	if err != nil {
		return ""
	}
	u.Path = path.Join(u.Path, "works", "random")
	query := url.Values{}
	query.Add("per-page", fmt.Sprintf("%d", count))
	if r.Email != "" {
		query.Add("mailto", r.Email)
	}
	u.RawQuery = query.Encode()
	return u.String()
}

// APIURL constructs a URL for accessing the OpenAlex API with a specific ID
func (r *Reader) APIURL(id string, idType string) string {
	u, err := url.Parse(openAlexBaseURL)
	if err != nil {
		return ""
	}

	var query = url.Values{}
	if r.Email != "" {
		query.Add("mailto", r.Email)
	}

	// Different paths based on ID type
	if idType == "OpenAlex" {
		u.Path = path.Join(u.Path, "works", id)
	} else {
		u.Path = path.Join(u.Path, "works")
		filter := fmt.Sprintf("ids.%s:%s", strings.ToLower(idType), id)
		query.Add("filter", filter)
	}

	u.RawQuery = query.Encode()
	return u.String()
}

// GetAll gets the metadata for a list of works from the OpenAlex API
func (r *Reader) GetAll(query url.Values) ([]Work, error) {
	url := r.QueryURL(query)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAlex API returned status %d", resp.StatusCode)
	}

	var response struct {
		Results []Work `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// Get fetches a single work from OpenAlex based on ID
func (r *Reader) Get(pid string) (*Work, error) {
	id, idType := utils.ValidateID(pid)
	if idType == "" || (idType != "DOI" && idType != "MAG" && idType != "OpenAlex" && idType != "PMID" && idType != "PMCID") {
		return nil, fmt.Errorf("invalid identifier: %s", pid)
	}

	url := r.APIURL(id, idType)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAlex API returned status %d", resp.StatusCode)
	}

	// For MAG, PMID, PMCID we get a list response
	if idType == "MAG" || idType == "PMID" || idType == "PMCID" {
		var response struct {
			Results []Work `json:"results"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, err
		}
		if len(response.Results) == 0 {
			return nil, fmt.Errorf("no results found for %s", pid)
		}
		work := response.Results[0]
		return &work, nil
	}

	// For DOI and OpenAlex we get a single object response
	var work Work
	if err := json.NewDecoder(resp.Body).Decode(&work); err != nil {
		return nil, err
	}

	return &work, nil
}

// GetAbstract extracts the abstract from OpenAlex's inverted index format
func GetAbstract(invertedIndex map[string][]int) string {
	if len(invertedIndex) == 0 {
		return ""
	}

	// Find the maximum position
	var maxPos int
	for _, positions := range invertedIndex {
		for _, pos := range positions {
			if pos > maxPos {
				maxPos = pos
			}
		}
	}

	// Create an array to hold all words in their positions
	abstractWords := make([]string, maxPos+1)

	// Fill in the words at their positions
	for word, positions := range invertedIndex {
		for _, pos := range positions {
			abstractWords[pos] = word
		}
	}

	// Join all words to form the abstract
	return strings.Join(abstractWords, " ")
}

// GetContributors extracts contributor information from authorships
func GetContributors(authorships []AuthorShip) []commonmeta.Contributor {
	var contributors []commonmeta.Contributor

	for _, authorship := range authorships {
		var affiliations []commonmeta.Affiliation
		for _, inst := range authorship.Institutions {
			if inst.DisplayName != "" || inst.ROR != "" {
				affiliations = append(affiliations, commonmeta.Affiliation{
					ID:   inst.ROR,
					Name: inst.DisplayName,
				})
			}
		}

		if authorship.Author.DisplayName != "" || authorship.Author.ORCID != "" {
			contributors = append(contributors, commonmeta.Contributor{
				ID:           authorship.Author.ORCID,
				Name:         authorship.Author.DisplayName,
				Affiliations: affiliations,
			})
		}
	}

	return contributors
}

// GetWorks fetches multiple works from OpenAlex based on IDs
func (r *Reader) GetWorks(ids []string) ([]Work, error) {
	var works []Work

	// Process in batches of 49 to respect API limits
	for i := 0; i < len(ids); i += 49 {
		end := i + 49
		if end > len(ids) {
			end = len(ids)
		}

		batch := ids[i:end]
		idsString := strings.Join(batch, "|")

		query := url.Values{}
		query.Add("filter", fmt.Sprintf("ids.openalex:%s", idsString))

		batchWorks, err := r.GetList(query)
		if err != nil {
			return nil, err
		}

		works = append(works, batchWorks...)
	}

	return works, nil
}

// GetFunders fetches funders from OpenAlex based on IDs
func (r *Reader) GetFunders(ids []string) ([]Funder, error) {
	var funders []Funder

	// Process in batches of 49 to respect API limits
	for i := 0; i < len(ids); i += 49 {
		end := i + 49
		if end > len(ids) {
			end = len(ids)
		}

		batch := ids[i:end]
		idsString := strings.Join(batch, "|")

		u, err := url.Parse(openAlexBaseURL)
		if err != nil {
			return nil, err
		}
		u.Path = path.Join(u.Path, "funders")

		query := url.Values{}
		query.Add("filter", fmt.Sprintf("ids.openalex:%s", idsString))
		if r.Email != "" {
			query.Add("mailto", r.Email)
		}
		u.RawQuery = query.Encode()

		resp, err := http.Get(u.String())
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("OpenAlex API returned status %d", resp.StatusCode)
		}

		var response struct {
			Results []Funder `json:"results"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, err
		}

		funders = append(funders, response.Results...)
	}

	return funders, nil
}

// GetSource fetches source information from OpenAlex
func (r *Reader) GetSource(sourceID string) (*Source, error) {
	if sourceID == "" || !strings.HasPrefix(sourceID, "https://openalex.org/") {
		return nil, fmt.Errorf("invalid OpenAlex source ID: %s", sourceID)
	}

	u, err := url.Parse(openAlexBaseURL)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "sources", strings.TrimPrefix(sourceID, "https://openalex.org/"))

	query := url.Values{}
	if r.Email != "" {
		query.Add("mailto", r.Email)
	}
	u.RawQuery = query.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAlex API returned status %d", resp.StatusCode)
	}

	var source Source
	if err := json.NewDecoder(resp.Body).Decode(&source); err != nil {
		return nil, err
	}

	return &source, nil
}

// GetContainer extracts container information from a work
func (r *Reader) GetContainer(work *Work) commonmeta.Container {
	container := commonmeta.Container{}

	if work.PrimaryLocation.Source.ID == "" {
		return container
	}

	// Try to get extended source information
	source, err := r.GetSource(work.PrimaryLocation.Source.ID)
	if err != nil {
		// Fall back to basic information in the work
		container.Type = OpenAlexToCommonmetaContainerTypes[work.PrimaryLocation.Source.Type]
		container.Title = work.PrimaryLocation.Source.DisplayName

		if work.PrimaryLocation.Source.ISSN != "" {
			container.Identifier = work.PrimaryLocation.Source.ISSN
			container.IdentifierType = "ISSN"
		} else if work.PrimaryLocation.Source.HomepageURL != "" {
			container.Identifier = work.PrimaryLocation.Source.HomepageURL
			container.IdentifierType = "URL"
		}
	} else {
		// Use extended source information
		container.Type = OpenAlexToCommonmetaContainerTypes[source.Type]
		container.Title = source.DisplayName

		if source.ISSN != "" {
			container.Identifier = source.ISSN
			container.IdentifierType = "ISSN"
		} else if source.HomepageURL != "" {
			container.Identifier = source.HomepageURL
			container.IdentifierType = "URL"
		}
	}

	// Add bibliographic information
	container.Volume = work.Biblio.Volume
	container.Issue = work.Biblio.Issue
	container.FirstPage = work.Biblio.FirstPage
	container.LastPage = work.Biblio.LastPage

	return container
}

// GetFiles extracts file information from a work
func GetFiles(work *Work) []commonmeta.File {
	if work.BestOALocation.PDFURL == "" {
		return nil
	}

	return []commonmeta.File{
		{
			URL:      work.BestOALocation.PDFURL,
			MimeType: "application/pdf",
		},
	}
}

// ValidateOpenAlex checks if a string is a valid OpenAlex ID
func ValidateOpenAlex(id string) string {
	if strings.HasPrefix(id, "https://openalex.org/") {
		return strings.TrimPrefix(id, "https://openalex.org/")
	}
	return ""
}

// GetRelated converts an OpenAlex work to a reference
func GetRelated(work *Work) commonmeta.Reference {
	if work == nil {
		return commonmeta.Reference{}
	}

	ref := commonmeta.Reference{
		ID:              doiutils.NormalizeDOI(work.DOI),
		Title:           work.Title,
		PublicationYear: strings.Split(work.PublicationDate, "-")[0],
		Volume:          work.Biblio.Volume,
		Issue:           work.Biblio.Issue,
		FirstPage:       work.Biblio.FirstPage,
		LastPage:        work.Biblio.LastPage,
	}

	if work.PrimaryLocation.Source.DisplayName != "" {
		ref.ContainerTitle = work.PrimaryLocation.Source.DisplayName
	}

	if work.PrimaryLocation.Source.HostOrganizationName != "" {
		ref.Publisher = work.PrimaryLocation.Source.HostOrganizationName
	}

	return ref
}

// ProcessReferences fetches and processes references for a work
func (r *Reader) ProcessReferences(referencedWorks []string) ([]commonmeta.Reference, error) {
	if len(referencedWorks) == 0 {
		return nil, nil
	}

	// Extract OpenAlex IDs from the reference URLs
	var openAlexIDs []string
	for _, refWork := range referencedWorks {
		if id := ValidateOpenAlex(refWork); id != "" {
			openAlexIDs = append(openAlexIDs, id)
		}
	}

	if len(openAlexIDs) == 0 {
		return nil, nil
	}

	works, err := r.GetWorks(openAlexIDs)
	if err != nil {
		return nil, err
	}

	var references []commonmeta.Reference
	for _, work := range works {
		ref := GetRelated(&work)
		if ref.ID != "" || ref.Title != "" {
			references = append(references, ref)
		}
	}

	return references, nil
}

// ProcessFunding processes funding information from a work
func (r *Reader) ProcessFunding(grants []Grant) ([]commonmeta.FundingReference, error) {
	if len(grants) == 0 {
		return nil, nil
	}

	// Extract funder IDs
	var funderIDs []string
	for _, grant := range grants {
		if id := ValidateOpenAlex(grant.Funder); id != "" {
			funderIDs = append(funderIDs, id)
		}
	}

	if len(funderIDs) == 0 {
		return nil, nil
	}

	// Get funders
	funders, err := r.GetFunders(funderIDs)
	if err != nil {
		return nil, err
	}

	// Map funder IDs to funder info
	funderMap := make(map[string]Funder)
	for _, funder := range funders {
		funderMap[funder.ID] = funder
	}

	// Create funding references
	var fundingRefs []commonmeta.FundingReference
	for _, grant := range grants {
		if funder, ok := funderMap[grant.Funder]; ok {
			ref := commonmeta.FundingReference{
				FunderName:  funder.DisplayName,
				AwardNumber: grant.AwardID,
			}

			if funder.Ids.ROR != "" {
				ref.FunderIdentifier = funder.Ids.ROR
				ref.FunderIdentifierType = "ROR"
			}

			fundingRefs = append(fundingRefs, ref)
		}
	}

	return fundingRefs, nil
}

// GetRandomIDs fetches random IDs from OpenAlex
func (r *Reader) GetRandomIDs(count int) ([]Work, error) {
	if count > 20 {
		count = 20 // Limit to 20 to be reasonable
	}

	url := r.SampleURL(count)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAlex API returned status %d", resp.StatusCode)
	}

	var response struct {
		Results []Work `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

// Read processes an OpenAlex work into CommonMeta format
func (r *Reader) Read(work *Work, options map[string]interface{}) (commonmeta.Data, error) {
	if work == nil {
		return commonmeta.Data{}, fmt.Errorf("no work data provided")
	}

	// Create a CommonMeta data object
	var data commonmeta.Data

	// Process ID - prefer DOI if available
	if work.DOI != "" {
		data.ID = doiutils.NormalizeDOI(work.DOI)
	} else if work.ID != "" {
		data.ID = work.ID
	}

	// Process type
	data.Type = "Other"
	if work.TypeCrossref != "" && OpenAlexToCommonmetaTypes[work.TypeCrossref] != "" {
		data.Type = OpenAlexToCommonmetaTypes[work.TypeCrossref]
	} else if work.Type != "" && OpenAlexToCommonmetaTypes[work.Type] != "" {
		data.Type = OpenAlexToCommonmetaTypes[work.Type]
	}

	// Process additional type if different from main type
	if work.Type != "" && OpenAlexToCommonmetaTypes[work.Type] != "" &&
		OpenAlexToCommonmetaTypes[work.Type] != data.Type {
		data.AdditionalType = OpenAlexToCommonmetaTypes[work.Type]
	}

	// Process title
	title := work.Title
	if title != "" {
		data.Titles = []commonmeta.Title{{Title: utils.Sanitize(title)}}
	}

	// Process URL
	if work.PrimaryLocation.LandingPageURL != "" {
		data.URL = utils.NormalizeURL(work.PrimaryLocation.LandingPageURL)
	} else if work.ID != "" {
		data.URL = utils.NormalizeURL(work.ID)
	}

	// Process publisher
	if work.PrimaryLocation.Source.HostOrganizationName != "" {
		data.Publisher = commonmeta.Publisher{
			Name: work.PrimaryLocation.Source.HostOrganizationName,
		}
	}

	// Process date
	if work.PublicationDate != "" || work.CreatedDate != "" {
		data.Date = commonmeta.Date{
			Published: work.PublicationDate,
		}
		if data.Date.Published == "" {
			data.Date.Published = work.CreatedDate
		}
	}

	// Process identifiers
	for idType, idValue := range work.Ids {
		if standardType, ok := OpenAlexIdentifierTypes[idType]; ok {
			data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
				Identifier:     idValue,
				IdentifierType: standardType,
			})
		}
	}

	// Process license
	if work.BestOALocation.License != "" {
		licenseID := work.BestOALocation.License
		if spdxID, ok := OpenAlexLicenses[licenseID]; ok {
			licenseID = spdxID
		}
		data.License = utils.URLToSPDX(licenseID)
	}

	// Process container
	data.Container = r.GetContainer(work)

	// Process contributors
	contributors := GetContributors(work.AuthorShips)
	if len(contributors) > 0 {
		processedContributors, _ := authorutils.ProcessContributors(contributors)
		data.Contributors = processedContributors
	}

	// Process abstract
	if len(work.AbstractInvertedIndex) > 0 {
		abstract := GetAbstract(work.AbstractInvertedIndex)
		if abstract != "" {
			data.Descriptions = []commonmeta.Description{{
				Description: utils.Sanitize(abstract),
				Type:        "Abstract",
			}}
		}
	}

	// Process subjects
	for _, topic := range work.Topics {
		if topic.Subfield.DisplayName != "" {
			data.Subjects = append(data.Subjects, commonmeta.Subject{
				Subject: topic.Subfield.DisplayName,
			})
		}
	}

	// Remove duplicate subjects
	if len(data.Subjects) > 0 {
		subjectMap := make(map[string]bool)
		var uniqueSubjects []commonmeta.Subject

		for _, subj := range data.Subjects {
			if !subjectMap[subj.Subject] {
				subjectMap[subj.Subject] = true
				uniqueSubjects = append(uniqueSubjects, subj)
			}
		}

		data.Subjects = uniqueSubjects
	}

	// Process files
	data.Files = GetFiles(work)

	// Process references and funding asynchronously
	var wg sync.WaitGroup
	var references []commonmeta.Reference
	var fundingRefs []commonmeta.FundingReference
	var refErr, fundErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		references, refErr = r.ProcessReferences(work.ReferencedWorks)
	}()

	go func() {
		defer wg.Done()
		fundingRefs, fundErr = r.ProcessFunding(work.Grants)
	}()

	wg.Wait()

	if refErr == nil && len(references) > 0 {
		data.References = references
	}

	if fundErr == nil && len(fundingRefs) > 0 {
		data.FundingReferences = fundingRefs
	}

	// Set language and version
	data.Language = work.Language
	data.Version = work.Version

	// Set provider
	data.Provider = "OpenAlex"

	// Apply any options
	for k, v := range options {
		switch k {
		case "validate":
			// Validation would be applied here
		}
	}

	return data, nil
}

// Fetch retrieves and processes metadata from OpenAlex by ID
func (r *Reader) Fetch(pid string, options map[string]interface{}) (commonmeta.Data, error) {
	work, err := r.Get(pid)
	if err != nil {
		return commonmeta.Data{}, err
	}

	return r.Read(work, options)
}

// FetchRandom retrieves and processes random metadata from OpenAlex
func (r *Reader) FetchRandom(count int, options map[string]interface{}) ([]commonmeta.Data, error) {
	works, err := r.GetRandomIDs(count)
	if err != nil {
		return nil, err
	}

	results := make([]commonmeta.Data, 0, len(works))
	for _, work := range works {
		data, err := r.Read(&work, options)
		if err == nil {
			results = append(results, data)
		}
	}

	return results, nil
}
