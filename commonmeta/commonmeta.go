// Package commonmeta provides functions to read and write commonmeta metadata.
package commonmeta

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/front-matter/commonmeta/schemautils"

	"github.com/xeipuuv/gojsonschema"
)

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

// Data represents the commonmeta metadata, defined in the commonmeta JSON Schema.
type Data struct {
	// required fields
	ID   string `db:"id" json:"id"`
	Type string `db:"type" json:"type"`

	// optional fields
	AdditionalType    string             `db:"additional_type" json:"additionalType,omitempty"`
	ArchiveLocations  []string           `db:"archive_locations" json:"archiveLocations,omitempty"`
	Container         Container          `db:"container" json:"container,omitempty"`
	Contributors      []Contributor      `db:"contributors" json:"contributors,omitempty"`
	Date              Date               `db:"date" json:"date,omitempty"`
	Descriptions      []Description      `db:"descriptions" json:"descriptions,omitempty"`
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
	Identifier     string `json:"identifier,omitempty"`
	IdentifierType string `json:"identifierType,omitempty"`
	Type           string `json:"type,omitempty"`
	Title          string `json:"title,omitempty"`
	FirstPage      string `json:"firstPage,omitempty"`
	LastPage       string `json:"lastPage,omitempty"`
	Volume         string `json:"volume,omitempty"`
	Issue          string `json:"issue,omitempty"`
}

// Contributor represents a contributor of a publication, defined in the commonmeta JSON Schema.
type Contributor struct {
	ID               string        `json:"id,omitempty"`
	Type             string        `json:"type,omitempty"`
	Name             string        `json:"name,omitempty"`
	GivenName        string        `json:"givenName,omitempty"`
	FamilyName       string        `json:"familyName,omitempty"`
	Affiliations     []Affiliation `json:"affiliations,omitempty"`
	ContributorRoles []string      `json:"contributorRoles,omitempty"`
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
	Name string `json:"name"`
}

// Reference represents the reference of a publication, defined in the commonmeta JSON Schema.
type Reference struct {
	Key             string `json:"key"`
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

// Read reads commonmeta metadata.
func Read(content Data) (Data, error) {
	data := content
	return data, nil
}

// Write writes commonmeta metadata.
func Write(data Data) ([]byte, []gojsonschema.ResultError) {
	output, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	validation := schemautils.JSONSchemaErrors(output)
	if !validation.Valid() {
		return nil, validation.Errors()
	}
	return output, nil
}

// WriteList writes commonmeta metadata in slice format.
func WriteList(list []Data) ([]byte, []gojsonschema.ResultError) {
	output, err := json.Marshal(list)
	if err != nil {
		fmt.Println(err)
	}
	validation := schemautils.JSONSchemaErrors(output)
	if !validation.Valid() {
		return nil, validation.Errors()
	}
	return output, nil
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
