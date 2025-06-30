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

const Version = "v0.31.0"

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
	"OpenAlex",
	"PMID",
	"PMCID",
	"PURL",
	"RID",
	"URL",
	"URN",
	"UUID",
	"Other",
}

// WorkTypes are the types Commonmeta supports for works
var WorkTypes = []string{"DOI", "Wikidata", "Openalex", "PMID", "PMCID", "UUID"}

// PersonTypes are the types Commonmeta supports for people
var PersonTypes = []string{"ORCID", "ISNI", "Openalex", "Wikidata"}

// OrganizationTypes are the types Commonmeta supports for organizations
var OrganizationTypes = []string{"ROR", "Wikidata", "Openalex", "Crossref Funder ID", "GRID", "ISNI"}

// FOSKeyMappings maps OECD FOS keys to OECD FOS strings
var FOSKeyMappings = map[string]string{
	"naturalSciences":                      "Natural sciences",
	"mathematics":                          "Mathematics",
	"computerAndInformationSciences":       "Computer and information sciences",
	"physicalSciences":                     "Physical sciences",
	"chemicalSciences":                     "Chemical sciences",
	"earthAndRelatedEnvironmentalSciences": "Earth and related environmental sciences",
	"biologicalSciences":                   "Biological sciences",
	"otherNaturalSciences":                 "Other natural sciences",
	"engineeringAndTechnology":             "Engineering and technology",
	"civilEngineering":                     "Civil engineering",
	"electricalEngineeringElectronicEngineeringInformationEngineering": "Electrical engineering, electronic engineering, information engineering",
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
	ID   string `avro:"id" json:"id"`
	Type string `avro:"type" json:"type"`

	// optional fields
	AdditionalType    string             `avro:"additional_type,omitempty" json:"additionalType,omitempty"`
	ArchiveLocations  []string           `avro:"archive_locations,omitempty" json:"archiveLocations,omitempty"`
	Container         Container          `avro:"container" json:"container"`
	ContentHTML       string             `avro:"content_html,omitempty" json:"content_html,omitempty"`
	Contributors      []Contributor      `avro:"contributors,omitempty" json:"contributors,omitempty"`
	Date              Date               `avro:"date" json:"date"`
	Descriptions      []Description      `avro:"descriptions,omitempty" json:"descriptions,omitempty"`
	FeatureImage      string             `avro:"image,omitempty" json:"feature_image,omitempty"`
	Files             []File             `avro:"files,omitempty" json:"files,omitempty"`
	FundingReferences []FundingReference `avro:"funding_references,omitempty" json:"fundingReferences,omitempty"`
	GeoLocations      []*GeoLocation     `avro:"geo_locations,omitempty" json:"geoLocations,omitempty"`
	Identifiers       []Identifier       `avro:"identifiers,omitempty" json:"identifiers,omitempty"`
	Language          string             `avro:"language,omitempty" json:"language,omitempty"`
	License           License            `avro:"license" json:"license"`
	Provider          string             `avro:"provider,omitempty" json:"provider,omitempty"`
	Publisher         Publisher          `avro:"publisher" json:"publisher"`
	References        []Reference        `avro:"references,omitempty" json:"references,omitempty"`
	Relations         []Relation         `avro:"relations,omitempty" json:"relations,omitempty"`
	Subjects          []Subject          `avro:"subjects,omitempty" json:"subjects,omitempty"`
	Titles            []Title            `avro:"titles,omitempty" json:"titles,omitempty"`
	URL               string             `avro:"url,omitempty" json:"url,omitempty"`
	Version           string             `avro:"version,omitempty" json:"version,omitempty"`
}

// Affiliation represents the affiliation of a contributor, defined in the commonmeta JSON Schema.
type Affiliation struct {
	ID         string `avro:"id,omitempty" json:"id,omitempty"`
	Name       string `avro:"name,omitempty" json:"name,omitempty"`
	AssertedBy string `avro:"assertedBy,omitempty" json:"assertedBy,omitempty"`
}

// Container represents the container of a publication, defined in the commonmeta JSON Schema.
type Container struct {
	Identifier     string   `avro:"identifier,omitempty" json:"identifier,omitempty"`
	IdentifierType string   `avro:"identifierType,omitempty" json:"identifierType,omitempty"`
	Type           string   `avro:"type,omitempty" json:"type,omitempty"`
	Title          string   `avro:"title,omitempty" json:"title,omitempty"`
	Description    string   `avro:"description,omitempty" json:"description,omitempty"`
	Language       string   `avro:"language,omitempty" json:"language,omitempty"`
	License        *License `avro:"license,omitempty" json:"license,omitempty"`
	Platform       string   `avro:"platform,omitempty" json:"platform,omitempty"`
	Favicon        string   `avro:"favicon,omitempty" json:"favicon,omitempty"`
	FirstPage      string   `avro:"firstPage,omitempty" json:"firstPage,omitempty"`
	LastPage       string   `avro:"lastPage,omitempty" json:"lastPage,omitempty"`
	Volume         string   `avro:"volume,omitempty" json:"volume,omitempty"`
	Issue          string   `avro:"issue,omitempty" json:"issue,omitempty"`
}

// Contributor represents a contributor of a publication, defined in the commonmeta JSON Schema.
type Contributor struct {
	ID               string         `avro:"id,omitempty" json:"id,omitempty"`
	Type             string         `avro:"type,omitempty" json:"type,omitempty"`
	Name             string         `avro:"name,omitempty" json:"name,omitempty"`
	GivenName        string         `avro:"givenName,omitempty" json:"givenName,omitempty"`
	FamilyName       string         `avro:"familyName,omitempty" json:"familyName,omitempty"`
	Affiliations     []*Affiliation `avro:"affiliations,omitempty" json:"affiliations,omitempty"`
	ContributorRoles []string       `avro:"contributorRoles,omitempty" json:"contributorRoles,omitempty"`
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
	AwardURI             string `json:"awardUri,omitempty"`
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
	Publisher       string `json:"publisher,omitempty"`
	PublicationYear string `json:"publicationYear,omitempty"`
	Volume          string `json:"volume,omitempty"`
	Issue           string `json:"issue,omitempty"`
	FirstPage       string `json:"first_page,omitempty"`
	LastPage        string `json:"last_page,omitempty"`
	Unstructured    string `json:"unstructured,omitempty"`
	AssertedBy      string `json:"assertedBy,omitempty"`
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

// FOSMappings maps OECD FOS strings to OECD FOS identifiers
var FOSMappings = map[string]string{
	"Natural sciences":                         "http://www.oecd.org/science/inno/38235147.pdf?1",
	"Mathematics":                              "http://www.oecd.org/science/inno/38235147.pdf?1.1",
	"Computer and information sciences":        "http://www.oecd.org/science/inno/38235147.pdf?1.2",
	"Physical sciences":                        "http://www.oecd.org/science/inno/38235147.pdf?1.3",
	"Chemical sciences":                        "http://www.oecd.org/science/inno/38235147.pdf?1.4",
	"Earth and related environmental sciences": "http://www.oecd.org/science/inno/38235147.pdf?1.5",
	"Biological sciences":                      "http://www.oecd.org/science/inno/38235147.pdf?1.6",
	"Other natural sciences":                   "http://www.oecd.org/science/inno/38235147.pdf?1.7",
	"Engineering and technology":               "http://www.oecd.org/science/inno/38235147.pdf?2",
	"Civil engineering":                        "http://www.oecd.org/science/inno/38235147.pdf?2.1",
	"Electrical engineering, electronic engineering, information engineering": "http://www.oecd.org/science/inno/38235147.pdf?2.2",
	"Mechanical engineering":               "http://www.oecd.org/science/inno/38235147.pdf?2.3",
	"Chemical engineering":                 "http://www.oecd.org/science/inno/38235147.pdf?2.4",
	"Materials engineering":                "http://www.oecd.org/science/inno/38235147.pdf?2.5",
	"Medical engineering":                  "http://www.oecd.org/science/inno/38235147.pdf?2.6",
	"Environmental engineering":            "http://www.oecd.org/science/inno/38235147.pdf?2.7",
	"Environmental biotechnology":          "http://www.oecd.org/science/inno/38235147.pdf?2.8",
	"Industrial biotechnology":             "http://www.oecd.org/science/inno/38235147.pdf?2.9",
	"Nano technology":                      "http://www.oecd.org/science/inno/38235147.pdf?2.10",
	"Other engineering and technologies":   "http://www.oecd.org/science/inno/38235147.pdf?2.11",
	"Medical and health sciences":          "http://www.oecd.org/science/inno/38235147.pdf?3",
	"Basic medicine":                       "http://www.oecd.org/science/inno/38235147.pdf?3.1",
	"Clinical medicine":                    "http://www.oecd.org/science/inno/38235147.pdf?3.2",
	"Health sciences":                      "http://www.oecd.org/science/inno/38235147.pdf?3.3",
	"Health biotechnology":                 "http://www.oecd.org/science/inno/38235147.pdf?3.4",
	"Other medical sciences":               "http://www.oecd.org/science/inno/38235147.pdf?3.5",
	"Agricultural sciences":                "http://www.oecd.org/science/inno/38235147.pdf?4",
	"Agriculture, forestry, and fisheries": "http://www.oecd.org/science/inno/38235147.pdf?4.1",
	"Animal and dairy science":             "http://www.oecd.org/science/inno/38235147",
	"Veterinary science":                   "http://www.oecd.org/science/inno/38235147",
	"Agricultural biotechnology":           "http://www.oecd.org/science/inno/38235147",
	"Other agricultural sciences":          "http://www.oecd.org/science/inno/38235147",
	"Social science":                       "http://www.oecd.org/science/inno/38235147.pdf?5",
	"Psychology":                           "http://www.oecd.org/science/inno/38235147.pdf?5.1",
	"Economics and business":               "http://www.oecd.org/science/inno/38235147.pdf?5.2",
	"Educational sciences":                 "http://www.oecd.org/science/inno/38235147.pdf?5.3",
	"Sociology":                            "http://www.oecd.org/science/inno/38235147.pdf?5.4",
	"Law":                                  "http://www.oecd.org/science/inno/38235147.pdf?5.5",
	"Political science":                    "http://www.oecd.org/science/inno/38235147.pdf?5.6",
	"Social and economic geography":        "http://www.oecd.org/science/inno/38235147.pdf?5.7",
	"Media and communications":             "http://www.oecd.org/science/inno/38235147.pdf?5.8",
	"Other social sciences":                "http://www.oecd.org/science/inno/38235147.pdf?5.9",
	"Humanities":                           "http://www.oecd.org/science/inno/38235147.pdf?6",
	"History and archaeology":              "http://www.oecd.org/science/inno/38235147.pdf?6.1",
	"Languages and literature":             "http://www.oecd.org/science/inno/38235147.pdf?6.2",
	"Philosophy, ethics and religion":      "http://www.oecd.org/science/inno/38235147.pdf?6.3",
	"Arts (arts, history of arts, performing arts, music)": "http://www.oecd.org/science/inno/38235147.pdf?6.4",
	"Other humanities": "http://www.oecd.org/science/inno/38235147.pdf?6.5",
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
