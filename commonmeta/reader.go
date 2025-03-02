// Package commonmeta provides functions to read and write commonmeta metadata.
package commonmeta

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path"
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

// ContributorRoles list of contributor roles defined in commonmeta schema.
//
// from commonmeta schema
var ContributorRoles = []string{
	"Author",
	"Editor",
	"Chair",
	"Reviewer",
	"ReviewAssistant",
	"StatsReviewer",
	"ReviewerExternal",
	"Reader",
	"Translator",
	"ContactPerson",
	"DataCollector",
	"DataManager",
	"Distributor",
	"HostingInstitution",
	"Producer",
	"ProjectLeader",
	"ProjectManager",
	"ProjectMember",
	"RegistrationAgency",
	"RegistrationAuthority",
	"RelatedPerson",
	"ResearchGroup",
	"RightsHolder",
	"Researcher",
	"Sponsor",
	"WorkPackageLeader",
	"Conceptualization",
	"DataCuration",
	"FormalAnalysis",
	"FundingAcquisition",
	"Investigation",
	"Methodology",
	"ProjectAdministration",
	"Resources",
	"Software",
	"Supervision",
	"Validation",
	"Visualization",
	"WritingOriginalDraft",
	"WritingReviewEditing",
	"Maintainer",
	"Other",
}

// ContainerTypes maps types to associated container types
var ContainerTypes = map[string]string{
	"BookChapter":        "Book",
	"Dataset":            "Database",
	"JournalArticle":     "Journal",
	"JournalIssue":       "Journal",
	"Book":               "BookSeries",
	"ProceedingsArticle": "Proceedings",
	"Article":            "Periodical",
	"BlogPost":           "Blog",
}

// IdentifierTypes list of identifier types defined in the commonmeta schema.
var IdentifierTypes = []string{
	"ARK",
	"arXiv",
	"Bibcode",
	"DOI",
	"GUID",
	"Handle",
	"ISBN",
	"ISSN",
	"PMID",
	"PMCID",
	"PURL",
	"RID",
	"URL",
	"URN",
	"UUID",
	"Other",
}

// FOSKeyMappings maps OECD FOS keys to OECD FOS strings
var FOSKeyMappings = map[string]string{
	"naturalSciences":                          "Natural sciences",
	"mathematics":                              "Mathematics",
	"computerAndInformationSciences":           "Computer and information sciences",
	"physicalSciences":                         "Physical sciences",
	"chemicalSciences":                         "Chemical sciences",
	"earthAndRelatedEnvironmentalSciences":     "Earth and related environmental sciences",
	"biologicalSciences":                       "Biological sciences",
	"otherNaturalSciences":                     "Other natural sciences",
	"engineeringAndTechnology":                 "Engineering and technology",
	"civilEngineering":                         "Civil engineering",
	"electricalEngineering":                    "Electrical engineering, electronic engineering, information engineering",
	"mechanicalEngineering":                    "Mechanical engineering",
	"chemicalEngineering":                      "Chemical engineering",
	"materialsEngineering":                     "Materials engineering",
	"medicalEngineering":                       "Medical engineering",
	"environmentalEngineering":                 "Environmental engineering",
	"environmentalBiotechnology":               "Environmental biotechnology",
	"industrialBiotechnology":                  "Industrial biotechnology",
	"nanoTechnology":                           "Nano technology",
	"otherEngineeringAndTechnologies":          "Other engineering and technologies",
	"medicalAndHealthSciences":                 "Medical and health sciences",
	"basicMedicine":                            "Basic medicine",
	"clinicalMedicine":                         "Clinical medicine",
	"healthSciences":                           "Health sciences",
	"healthBiotechnology":                      "Health biotechnology",
	"otherMedicalSciences":                     "Other medical sciences",
	"agriculturalSciences":                     "Agricultural sciences",
	"agricultureForestryAndFisheries":          "Agriculture, forestry, and fisheries",
	"animalAndDairyScience":                    "Animal and dairy science",
	"veterinaryScience":                        "Veterinary science",
	"agriculturalBiotechnology":                "Agricultural biotechnology",
	"otherAgriculturalSciences":                "Other agricultural sciences",
	"socialScience":                            "Social science",
	"psychology":                               "Psychology",
	"economicsAndBusiness":                     "Economics and business",
	"educationalSciences":                      "Educational sciences",
	"sociology":                                "Sociology",
	"law":                                      "Law",
	"politicalScience":                         "Political science",
	"socialAndEconomicGeography":               "Social and economic geography",
	"mediaAndCommunications":                   "Media and communications",
	"otherSocialSciences":                      "Other social sciences",
	"humanities":                               "Humanities",
	"historyAndArchaeology":                    "History and archaeology",
	"languagesAndLiterature":                   "Languages and literature",
	"philosophyEthicsAndReligion":              "Philosophy, ethics and religion",
	"artsArtsHistoryOfArtsPerformingArtsMusic": "Arts (arts, history of arts, performing arts, music)",
	"otherHumanities":                          "Other humanities",
}

// FOSStringMappings maps OECD FOS strings to OECD FOS keys
var FOSStringMappings = map[string]string{
	"Natural sciences":                         "naturalSciences",
	"Mathematics":                              "mathematics",
	"Computer and information sciences":        "computerAndInformationSciences",
	"Physical sciences":                        "physicalSciences",
	"Chemical sciences":                        "chemicalSciences",
	"Earth and related environmental sciences": "earthAndRelatedEnvironmentalSciences",
	"Biological sciences":                      "biologicalSciences",
	"Other natural sciences":                   "otherNaturalSciences",
	"Engineering and technology":               "engineeringAndTechnology",
	"Civil engineering":                        "civilEngineering",
	"Electrical engineering, electronic engineering, information engineering": "electricalEngineering",
	"Mechanical engineering":               "mechanicalEngineering",
	"Chemical engineering":                 "chemicalEngineering",
	"Materials engineering":                "materialsEngineering",
	"Medical engineering":                  "medicalEngineering",
	"Environmental engineering":            "environmentalEngineering",
	"Environmental biotechnology":          "environmentalBiotechnology",
	"Industrial biotechnology":             "industrialBiotechnology",
	"Nano technology":                      "nanoTechnology",
	"Other engineering and technologies":   "otherEngineeringAndTechnologies",
	"Medical and health sciences":          "medicalAndHealthSciences",
	"Basic medicine":                       "basicMedicine",
	"Clinical medicine":                    "clinicalMedicine",
	"Health sciences":                      "healthSciences",
	"Health biotechnology":                 "healthBiotechnology",
	"Other medical sciences":               "otherMedicalSciences",
	"Agricultural sciences":                "agriculturalSciences",
	"Agriculture, forestry, and fisheries": "agricultureForestryAndFisheries",
	"Animal and dairy science":             "animalAndDairyScience",
	"Veterinary science":                   "veterinaryScience",
	"Agricultural biotechnology":           "agriculturalBiotechnology",
	"Other agricultural sciences":          "otherAgriculturalSciences",
	"Social science":                       "socialScience",
	"Psychology":                           "psychology",
	"Economics and business":               "economicsAndBusiness",
	"Educational sciences":                 "educationalSciences",
	"Sociology":                            "sociology",
	"Law":                                  "law",
	"Political science":                    "politicalScience",
	"Social and economic geography":        "socialAndEconomicGeography",
	"Media and communications":             "mediaAndCommunications",
	"Other social sciences":                "otherSocialSciences",
	"Humanities":                           "humanities",
	"History and archaeology":              "historyAndArchaeology",
	"Languages and literature":             "languagesAndLiterature",
	"Philosophy, ethics and religion":      "philosophyEthicsAndReligion",
	"Arts (arts, history of arts, performing arts, music)": "artsArtsHistoryOfArtsPerformingArtsMusic",
	"Other humanities": "otherHumanities",
}

// Data represents the commonmeta metadata, defined in the commonmeta JSON Schema.
type Data struct {
	// required fields
	ID   string `db:"id" json:"id"`
	Type string `db:"type" json:"type"`

	// optional fields
	AdditionalType    string             `db:"additional_type" json:"additionalType,omitempty"`
	ArchiveLocations  []string           `db:"archive_locations" json:"archiveLocations,omitempty"`
	Container         Container          `db:"container" json:"container,omitempty"`
	ContentText       string             `db:"content" json:"content_text,omitempty"`
	Contributors      []Contributor      `db:"contributors" json:"contributors,omitempty"`
	Date              Date               `db:"date" json:"date,omitempty"`
	Descriptions      []Description      `db:"descriptions" json:"descriptions,omitempty"`
	FeatureImage      string             `db:"image" json:"feature_image,omitempty"`
	Files             []File             `db:"files" json:"files,omitempty"`
	FundingReferences []FundingReference `db:"funding_references" json:"fundingReferences,omitempty"`
	GeoLocations      []GeoLocation      `db:"geo_locations" json:"geoLocations,omitempty"`
	Identifiers       []Identifier       `db:"identifiers" json:"identifiers,omitempty"`
	Language          string             `db:"language" json:"language,omitempty"`
	License           License            `db:"license" json:"license,omitempty"`
	Provider          string             `db:"provider" json:"provider,omitempty"`
	Publisher         Publisher          `db:"publisher" json:"publisher,omitempty"`
	References        []Reference        `db:"references" json:"references,omitempty"`
	Relations         []Relation         `db:"relations" json:"relations,omitempty"`
	Subjects          []Subject          `db:"subjects" json:"subjects,omitempty"`
	Titles            []Title            `db:"titles" json:"titles,omitempty"`
	URL               string             `db:"url" json:"url,omitempty"`
	Version           string             `db:"version" json:"version,omitempty"`
}

// Affiliation represents the affiliation of a contributor, defined in the commonmeta JSON Schema.
type Affiliation struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Container represents the container of a publication, defined in the commonmeta JSON Schema.
type Container struct {
	Identifier     string  `json:"identifier,omitempty"`
	IdentifierType string  `json:"identifierType,omitempty"`
	Type           string  `json:"type,omitempty"`
	Title          string  `json:"title,omitempty"`
	Description    string  `json:"description,omitempty"`
	Language       string  `json:"language,omitempty"`
	License        License `json:"license,omitempty"`
	Platform       string  `json:"platform,omitempty"`
	Favicon        string  `json:"favicon,omitempty"`
	FirstPage      string  `json:"firstPage,omitempty"`
	LastPage       string  `json:"lastPage,omitempty"`
	Volume         string  `json:"volume,omitempty"`
	Issue          string  `json:"issue,omitempty"`
}

// Contributor represents a contributor of a publication, defined in the commonmeta JSON Schema.
type Contributor struct {
	ID               string         `json:"id,omitempty"`
	Type             string         `json:"type,omitempty"`
	Name             string         `json:"name,omitempty"`
	GivenName        string         `json:"givenName,omitempty"`
	FamilyName       string         `json:"familyName,omitempty"`
	Affiliations     []*Affiliation `json:"affiliations,omitempty"`
	ContributorRoles []string       `json:"contributorRoles,omitempty"`
}

// Date represents the date of a publication, defined in the commonmeta JSON Schema.
type Date struct {
	Created     string `json:"created,omitempty"`
	Submitted   string `json:"submitted,omitempty"`
	Accepted    string `json:"accepted,omitempty"`
	Published   string `json:"published,omitempty"`
	Updated     string `json:"updated,omitempty"`
	Accessed    string `json:"accessed,omitempty"`
	Available   string `json:"available,omitempty"`
	Copyrighted string `json:"copyrighted,omitempty"`
	Collected   string `json:"collected,omitempty"`
	Valid       string `json:"valid,omitempty"`
	Withdrawn   string `json:"withdrawn,omitempty"`
	Other       string `json:"other,omitempty"`
}

// Description represents the description of a publication, defined in the commonmeta JSON Schema.
type Description struct {
	Description string `json:"description"`
	Type        string `json:"type,omitempty"`
	Language    string `json:"language,omitempty"`
}

// File represents a file of a publication, defined in the commonmeta JSON Schema.
type File struct {
	Bucket   string `json:"bucket,omitempty"`
	Key      string `json:"key,omitempty"`
	Checksum string `json:"checksum,omitempty"`
	URL      string `json:"url"`
	Size     int    `json:"size,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

// FundingReference represents the funding reference of a publication, defined in the commonmeta JSON Schema.
type FundingReference struct {
	FunderIdentifier     string `json:"funderIdentifier,omitempty"`
	FunderIdentifierType string `json:"funderIdentifierType,omitempty"`
	FunderName           string `json:"funderName,omitempty"`
	AwardNumber          string `json:"awardNumber,omitempty"`
	AwardTitle           string `json:"awardTitle,omitempty"`
	AwardURI             string `json:"award_uri,omitempty"`
}

// GeoLocation represents the geographical location of a publication, defined in the commonmeta JSON Schema.
type GeoLocation struct {
	GeoLocationPlace string           `json:"geoLocationPlace,omitempty"`
	GeoLocationPoint GeoLocationPoint `json:"geoLocationPoint,omitempty"`
	GeoLocationBox   GeoLocationBox   `json:"geoLocationBox,omitempty"`
}

// GeoLocationPoint represents a point in a geographical location, defined in the commonmeta JSON Schema.
type GeoLocationPoint struct {
	PointLongitude float64 `json:"pointLongitude,omitempty"`
	PointLatitude  float64 `json:"pointLatitude,omitempty"`
}

// GeoLocationBox represents a box in a geographical location, defined in the commonmeta JSON Schema.
type GeoLocationBox struct {
	EastBoundLongitude float64 `json:"eastBoundLongitude,omitempty"`
	WestBoundLongitude float64 `json:"westBoundLongitude,omitempty"`
	SouthBoundLatitude float64 `json:"southBoundLatitude,omitempty"`
	NorthBoundLatitude float64 `json:"northBoundLatitude,omitempty"`
}

// GeoLocationPolygon represents a polygon in a geographical location, defined in the commonmeta JSON Schema.
type GeoLocationPolygon struct {
	PolygonPoints  []GeoLocationPoint `json:"polygon_points,omitempty"`
	InPolygonPoint GeoLocationPoint   `json:"in_polygon_point,omitempty"`
}

// Identifier represents the identifier of a publication, defined in the commonmeta JSON Schema.
type Identifier struct {
	Identifier     string `json:"identifier"`
	IdentifierType string `json:"identifierType"`
}

// License represents the license of a publication, defined in the commonmeta JSON Schema.
type License struct {
	ID  string `json:"id,omitempty"`
	URL string `json:"url,omitempty"`
}

// Publisher represents the publisher of a publication, defined in the commonmeta JSON Schema.
type Publisher struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Reference represents the reference of a publication, defined in the commonmeta JSON Schema.
type Reference struct {
	Key             string `json:"key,omitempty"`
	ID              string `json:"id,omitempty"`
	Type            string `json:"type,omitempty"`
	Title           string `json:"title,omitempty"`
	PublicationYear string `json:"publicationYear,omitempty"`
	Unstructured    string `json:"unstructured,omitempty"`
}

// Relation represents the relation of a publication, defined in the commonmeta JSON Schema.
type Relation struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// Subject represents the subject of a publication, defined in the commonmeta JSON Schema.
type Subject struct {
	Subject string `json:"subject"`
}

// Title represents the title of a publication, defined in the commonmeta JSON Schema.
type Title struct {
	Title    string `json:"title,omitempty"`
	Type     string `json:"type,omitempty"`
	Language string `json:"language,omitempty"`
}

type APIResponse struct {
	DOI         string `json:"doi"`
	ID          string `json:"id,omitempty"`
	DOIBatchID  string `json:"doi_batch_id,omitempty"`
	UUID        string `json:"uuid,omitempty"`
	Community   string `json:"community,omitempty"`
	CommunityID string `json:"community_id,omitempty"`
	Created     string `json:"created,omitempty"`
	Updated     string `json:"updated,omitempty"`
	Timestamp   string `json:"timestamp,omitempty"`
	Status      string `json:"status,omitempty"`
}

// CMToSOMappings maps Commonmeta types to Schema.org types.
var CMToSOMappings = map[string]string{
	"Article":        "Article",
	"Audiovisual":    "CreativeWork",
	"BlogPost":       "BlogPosting",
	"Book":           "Book",
	"BookChapter":    "BookChapter",
	"Collection":     "CreativeWork",
	"Dataset":        "Dataset",
	"Dissertation":   "Dissertation",
	"Document":       "CreativeWork",
	"Entry":          "CreativeWork",
	"Event":          "CreativeWork",
	"Figure":         "CreativeWork",
	"Image":          "CreativeWork",
	"Instrument":     "Instrument",
	"JournalArticle": "ScholarlyArticle",
	"LegalDocument":  "Legislation",
	"Software":       "SoftwareSourceCode",
	"Presentation":   "PresentationDigitalDocument",
}

// Load loads the metadata for a single work from a JSON file
func Load(filename string) (Data, error) {
	var data Data

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
	err = decoder.Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

// LoadAll loads a list of commonmeta metadata from a JSON string and returns Commonmeta metadata.
func LoadAll(filename string) ([]Data, error) {
	var data []Data

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
	err = decoder.Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}

// Read reads commonmeta metadata.
func Read(content Data) (Data, error) {
	data := content
	return data, nil
}

// ReadAll reads commonmeta metadata in slice format.
func ReadAll(content []Data) ([]Data, error) {
	data := content
	return data, nil
}

// Pages returns the first and last page of a work as a string.
func (c *Container) Pages() string {
	if c.FirstPage == "" {
		return c.LastPage
	}
	if c.LastPage == "" {
		return c.FirstPage
	}
	return c.FirstPage + "-" + c.LastPage
}
