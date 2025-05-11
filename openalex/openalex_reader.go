// Package openalex provides functions to convert Openalex metadata to the commonmeta metadata format.
package openalex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/spdx"
	"github.com/front-matter/commonmeta/utils"
)

const openAlexBaseURL = "https://api.openalex.org"

// Reader struct to hold any configuration for the Openalex reader
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
	AuthorShips           []Authorship      `json:"authorships"`
	Ids                   map[string]string `json:"ids"`
	PrimaryLocation       Location          `json:"primary_location"`
	BestOALocation        Location          `json:"best_oa_location"`
	Topics                []Topic           `json:"topics"`
	Biblio                Biblio            `json:"biblio"`
	ReferencedWorks       []string          `json:"referenced_works"`
	RelatedWorks          []string          `json:"related_works"`
	Grants                []Grant           `json:"grants"`
}

// Authorship represents author information in OpenAlex
type Authorship struct {
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

// Source represents an OpenAlex source
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
func (r *Reader) QueryURL(number int, page int, publisher string, type_ string, sample bool, ids string, year string, orcid string, ror string, hasORCID bool, hasROR bool, hasReferences bool, hasRelation bool, hasAbstract bool, hasAward bool, hasLicense bool, hasArchive bool) string {
	types := []string{
		"article",
		"book-chapter",
		"dataset",
		"preprint",
		"dissertation",
		"book",
		"review",
		"paratext",
		"libguides",
		"letter",
		"other",
		"reference-entry",
		"report",
		"editorial",
		"peer-review",
		"erratum",
		"standard",
		"grant",
		"supplementary-materials",
		"retraction",
	}

	u, _ := url.Parse("https://api.openalex.org")
	u.Path = path.Join(u.Path, "works")
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
		values.Add("per-page", strconv.Itoa(number))
		values.Add("page", strconv.Itoa(page))

		// sort results by published date in descending order
		values.Add("sort", "publication_date:desc")
	}

	var filters []string
	if ids != "" {
		filters = append(filters, "member:"+ids)
	}
	if type_ != "" && slices.Contains(types, type_) {
		filters = append(filters, "type:"+type_)
	}
	if ror != "" {
		r, _ := utils.ValidateROR(ror)
		if r != "" {
			filters = append(filters, "authorships.institutions.ror:"+r)
		}
	}
	if orcid != "" {
		o, _ := utils.ValidateORCID(orcid)
		if o != "" {
			filters = append(filters, "authorships.author.orcid:"+o)
		}
	}
	if year != "" {
		filters = append(filters, "publication_year:"+year)
	}
	if hasORCID {
		filters = append(filters, "has-orcid:true")
	}
	// if hasROR {
	// 	filters = append(filters, "has-ror-id:true")
	// }
	if hasReferences {
		filters = append(filters, "has-references:true")
	}
	if hasAbstract {
		filters = append(filters, "has-abstract:true")
	}
	// if hasAward {
	// 	filters = append(filters, "has-award:true")
	// }
	// if hasLicense {
	// 	filters = append(filters, "has-license:true")
	// }
	// if hasArchive {
	// 	filters = append(filters, "has-archive:true")
	// }
	if len(filters) > 0 {
		values.Add("filter", strings.Join(filters[:], ","))
	}
	u.RawQuery = values.Encode()
	return u.String()
}

// APIURL constructs a URL for accessing the OpenAlex API with a specific ID
func (r *Reader) APIURL(id string, idType string) string {
	u, _ := url.Parse(openAlexBaseURL)
	var query = url.Values{}
	if r.Email != "" {
		query.Add("mailto", r.Email)
	}

	// Different paths based on ID type
	if idType == "OpenAlex" {
		u.Path = path.Join(u.Path, "works", id)
	} else if idType == "DOI" {
		u.Path = path.Join(u.Path, "works", doiutils.NormalizeDOI(id))
	} else {
		u.Path = path.Join(u.Path, "works")
		filter := fmt.Sprintf("ids.%s:%s", strings.ToLower(idType), id)
		query.Add("filter", filter)
	}

	u.RawQuery = query.Encode()
	return u.String()
}

// GetAll gets the metadata for a list of works from the OpenAlex API
func (r *Reader) GetAll(number int, page int, publisher string, type_ string, sample bool, ids string, year string, ror string, orcid string, hasORCID bool, hasROR bool, hasReferences bool, hasRelation bool, hasAbstract bool, hasAward bool, hasLicense bool, hasArchive bool) ([]Work, error) {
	var response struct {
		Results []Work `json:"results"`
	}
	url := r.QueryURL(number, page, publisher, type_, sample, ids, year, orcid, ror, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense, hasArchive)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAlex API returned status %d", resp.StatusCode)
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
		fmt.Println(err, url)
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
func GetContributors(authorships []Authorship) []commonmeta.Contributor {
	var contributors []commonmeta.Contributor

	for _, authorship := range authorships {
		var affiliations []*commonmeta.Affiliation
		for _, inst := range authorship.Institutions {
			if inst.DisplayName != "" || inst.ROR != "" {
				affiliations = append(affiliations, &commonmeta.Affiliation{
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

	// Parse in batches of 49 to respect API limits
	for i := 0; i < len(ids); i += 49 {
		end := min(i+49, len(ids))

		batch := ids[i:end]
		idsString := strings.Join(batch, "|")

		batchWorks, err := r.GetAll(end, 1, "", "", false, idsString, "", "", "", false, false, false, false, false, false, false, false)
		if err != nil {
			return nil, err
		}
		works = append(works, batchWorks...)
	}
	return works, nil
}

// GetFunders fetches multiple funders from OpenAlex based on IDs
func (r *Reader) GetFunders(ids []string) ([]Funder, error) {
	var funders []Funder

	// Process in batches of 49 to respect API limits
	for i := 0; i < len(ids); i += 49 {
		end := min(i+49, len(ids))

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
func (r *Reader) GetContainer(work *Work) *commonmeta.Container {
	container := commonmeta.Container{}

	if work.PrimaryLocation.Source.ID == "" {
		return &container
	}

	// Try to get extended source information
	source, err := r.GetSource(work.PrimaryLocation.Source.ID)
	if err != nil {
		// Fall back to basic information in the work
		container.Type = OpenAlexContainerTypes[work.PrimaryLocation.Source.Type]
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
		container.Type = OpenAlexContainerTypes[source.Type]
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

	return &container
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
		ref.Title = work.PrimaryLocation.Source.DisplayName
	}

	if work.PrimaryLocation.Source.HostOrganizationName != "" {
		ref.Publisher = work.PrimaryLocation.Source.HostOrganizationName
	}

	return ref
}

// ParseReferences fetches and processes references for a work
func (r *Reader) ParseReferences(referencedWorks []string) ([]commonmeta.Reference, error) {
	if len(referencedWorks) == 0 {
		return nil, nil
	}

	// Extract OpenAlex IDs from the reference URLs
	var openAlexIDs []string
	for _, refWork := range referencedWorks {
		if id, _ := utils.ValidateOpenalex(refWork); id != "" {
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

// ParseFunding processes funding information from a work
func (r *Reader) ParseFunding(grants []Grant) ([]commonmeta.FundingReference, error) {
	if len(grants) == 0 {
		return nil, nil
	}

	// Extract funder IDs
	var funderIDs []string
	for _, grant := range grants {
		if id, _ := utils.ValidateOpenalex(grant.Funder); id != "" {
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

// Read OpenAlex response and return work struct in Commonmeta format
func (r *Reader) Read(work *Work) (commonmeta.Data, error) {
	var data commonmeta.Data
	if work == nil {
		return data, fmt.Errorf("no work data provided")
	}

	if work.DOI != "" {
		data.ID = doiutils.NormalizeDOI(work.DOI)
	} else if work.ID != "" {
		data.ID = work.ID
	}

	data.Type = "Other"
	if work.TypeCrossref != "" && OAToCMMappings[work.TypeCrossref] != "" {
		data.Type = OAToCMMappings[work.TypeCrossref]
	} else if work.Type != "" && OAToCMMappings[work.Type] != "" {
		data.Type = OAToCMMappings[work.Type]
	}

	// Parse additional type if different from type
	if work.Type != "" && OAToCMMappings[work.Type] != "" &&
		OAToCMMappings[work.Type] != data.Type {
		data.AdditionalType = OAToCMMappings[work.Type]
	}

	title := work.Title
	if title != "" {
		data.Titles = []commonmeta.Title{{Title: utils.Sanitize(title)}}
	}

	if work.PrimaryLocation.LandingPageURL != "" {
		url, err := utils.NormalizeURL(work.PrimaryLocation.LandingPageURL, true, true)
		if err != nil {
			fmt.Println(err)
		}
		data.URL = url
	} else if work.ID != "" {
		url, err := utils.NormalizeURL(work.ID, true, true)
		if err != nil {
			fmt.Println(err)
		}
		data.URL = url
	}

	if work.PrimaryLocation.Source.HostOrganizationName != "" {
		data.Publisher = &commonmeta.Publisher{
			Name: work.PrimaryLocation.Source.HostOrganizationName,
		}
	}

	if work.PublicationDate != "" || work.CreatedDate != "" {
		data.Date = commonmeta.Date{
			Published: work.PublicationDate,
		}
		if data.Date.Published == "" {
			data.Date.Published = work.CreatedDate
		}
	}

	for idType, idValue := range work.Ids {
		if standardType, ok := OpenAlexIdentifierTypes[idType]; ok {
			data.Identifiers = append(data.Identifiers, commonmeta.Identifier{
				Identifier:     idValue,
				IdentifierType: standardType,
			})
		}
	}

	if work.BestOALocation.License != "" {
		var ID, URL string
		licenseID := work.BestOALocation.License
		if spdxID, ok := OpenAlexLicenses[licenseID]; ok {
			licenseID = spdxID
		}
		license, err := spdx.Search(licenseID)
		if err != nil {
			fmt.Println(err)
		}
		ID = license.LicenseID
		if len(license.SeeAlso) > 0 {
			URL = license.SeeAlso[0]
		}
		data.License = &commonmeta.License{
			ID:  ID,
			URL: URL,
		}
	}
	data.Container = r.GetContainer(work)

	contributors := GetContributors(work.AuthorShips)
	if len(contributors) > 0 {
		data.Contributors = contributors
	}

	if len(work.AbstractInvertedIndex) > 0 {
		abstract := GetAbstract(work.AbstractInvertedIndex)
		if abstract != "" {
			data.Descriptions = []commonmeta.Description{{
				Description: utils.Sanitize(abstract),
				Type:        "Abstract",
			}}
		}
	}

	for _, topic := range work.Topics {
		subject := commonmeta.Subject{
			Subject: topic.Subfield.DisplayName,
		}
		if topic.Subfield.DisplayName != "" && !slices.Contains(data.Subjects, subject) {
			data.Subjects = append(data.Subjects, subject)
		}
	}

	data.Files = GetFiles(work)

	// parse references and funding asynchronously
	var wg sync.WaitGroup
	var references []commonmeta.Reference
	var fundingRefs []commonmeta.FundingReference
	var refErr, fundErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		references, refErr = r.ParseReferences(work.ReferencedWorks)
	}()

	go func() {
		defer wg.Done()
		fundingRefs, fundErr = r.ParseFunding(work.Grants)
	}()

	wg.Wait()
	if refErr == nil && len(references) > 0 {
		data.References = references
	}

	if fundErr == nil && len(fundingRefs) > 0 {
		data.FundingReferences = fundingRefs
	}

	data.Language = work.Language
	data.Version = work.Version
	data.Provider = "OpenAlex"

	return data, nil
}

// ReadAll reads a list of OpenAlex JSON responses and returns a list of works in Commonmeta format
func (r *Reader) ReadAll(works []Work) ([]commonmeta.Data, error) {
	var data []commonmeta.Data
	for _, v := range works {
		d, err := r.Read(&v)
		if err != nil {
			fmt.Println(err)
		}
		data = append(data, d)
	}
	return data, nil
}

// Fetch retrieves and parses metadata from OpenAlex by ID
func (r *Reader) Fetch(pid string) (commonmeta.Data, error) {
	var data commonmeta.Data
	work, err := r.Get(pid)
	if err != nil {
		return data, err
	}

	return r.Read(work)
}

// FetchAll retrieves and parses metadata from OpenAlex by query
func (r *Reader) FetchAll(number int, page int, publisher string, type_ string, sample bool, ids string, year string, ror string, orcid string, hasORCID bool, hasROR bool, hasReferences bool, hasRelation bool, hasAbstract bool, hasAward bool, hasLicense bool, hasArchive bool) ([]commonmeta.Data, error) {
	var data []commonmeta.Data
	content, err := r.GetAll(number, page, publisher, type_, sample, ids, year, orcid, ror, hasORCID, hasROR, hasReferences, hasRelation, hasAbstract, hasAward, hasLicense, hasArchive)
	if err != nil {
		return data, err
	}

	data, err = r.ReadAll(content)
	if err != nil {
		return data, err
	}
	return data, nil
}

// FetchRandom retrieves and parses random metadata from OpenAlex
