package types

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
	Url               string             `db:"url" json:"url,omitempty"`
	Version           string             `db:"version" json:"version,omitempty"`
}

type Affiliation struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

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

type Contributor struct {
	ID               string        `json:"id,omitempty"`
	Type             string        `json:"type,omitempty"`
	Name             string        `json:"name,omitempty"`
	GivenName        string        `json:"givenName,omitempty"`
	FamilyName       string        `json:"familyName,omitempty"`
	Affiliations     []Affiliation `json:"affiliations,omitempty"`
	ContributorRoles []string      `json:"contributorRoles,omitempty"`
}

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

type Description struct {
	Description string `json:"description"`
	Type        string `json:"type,omitempty"`
	Language    string `json:"language,omitempty"`
}

type File struct {
	Bucket   string `json:"bucket,omitempty"`
	Key      string `json:"key,omitempty"`
	Checksum string `json:"checksum,omitempty"`
	Url      string `json:"url"`
	Size     int    `json:"size,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

type FundingReference struct {
	FunderIdentifier     string `json:"funderIdentifier,omitempty"`
	FunderIdentifierType string `json:"funderIdentifierType,omitempty"`
	FunderName           string `json:"funderName"`
	AwardNumber          string `json:"awardNumber,omitempty"`
	AwardURI             string `json:"award_uri,omitempty"`
}

type GeoLocation struct {
	GeoLocationPlace string           `json:"geoLocationPlace,omitempty"`
	GeoLocationPoint GeoLocationPoint `json:"geoLocationPoint,omitempty"`
	GeoLocationBox   GeoLocationBox   `json:"geoLocationBox,omitempty"`
}

type GeoLocationPoint struct {
	PointLongitude float64 `json:"pointLongitude,omitempty"`
	PointLatitude  float64 `json:"pointLatitude,omitempty"`
}

type GeoLocationBox struct {
	EastBoundLongitude float64 `json:"eastBoundLongitude,omitempty"`
	WestBoundLongitude float64 `json:"westBoundLongitude,omitempty"`
	SouthBoundLatitude float64 `json:"southBoundLatitude,omitempty"`
	NorthBoundLatitude float64 `json:"northBoundLatitude,omitempty"`
}

type GeoLocationPolygon struct {
	PolygonPoints  []GeoLocationPoint `json:"polygon_points,omitempty"`
	InPolygonPoint GeoLocationPoint   `json:"in_polygon_point,omitempty"`
}

type Identifier struct {
	Identifier     string `json:"identifier"`
	IdentifierType string `json:"identifierType"`
}

type License struct {
	ID  string `json:"id,omitempty"`
	Url string `json:"url,omitempty"`
}

type Publisher struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

type Reference struct {
	Key             string `json:"key"`
	ID              string `json:"id,omitempty"`
	Type            string `json:"type,omitempty"`
	Title           string `json:"title,omitempty"`
	PublicationYear string `json:"publicationYear,omitempty"`
	Unstructured    string `json:"unstructured,omitempty"`
}

type Relation struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Subject struct {
	Subject string `json:"subject"`
}

type Title struct {
	Title    string `json:"title,omitempty"`
	Type     string `json:"type,omitempty"`
	Language string `json:"language,omitempty"`
}
