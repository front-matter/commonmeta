package metadata

type Metadata struct {
	// required properties
	ID           string        `json:"id"`
	Type         string        `json:"type"`
	Url          string        `json:"url,omitempty"`
	Contributors []Contributor `json:"contributors,omitempty"`
	Titles       []Title       `json:"titles,omitempty"`
	Publisher    struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name"`
	}
	Date struct {
		Created   string `json:"created,omitempty"`
		Submitted string `json:"submitted,omitempty"`
		Accepted  string `json:"accepted,omitempty"`
		Published string `json:"published,omitempty"`
		Updated   string `json:"updated,omitempty"`
		Accessed  string `json:"accessed,omitempty"`
		Available string `json:"available,omitempty"`
		Withdrawn string `json:"withdrawn,omitempty"`
	}

	// recommended and optional properties
	AdditionalType string    `json:"additional_name,omitempty"`
	Subjects       []Subject `json:"subjects,omitempty"`
	// The language of the resource. Use one of the language codes from the IETF BCP 47 standard.
	Language             string `json:"language,omitempty"`
	AlternateIdentifiers struct {
		AlternateIdentifier     string `json:"alternate_identifier"`
		AlternateIdentifierType string `json:"alternate_identifier_type"`
	}
	Relations []Relation `json:"relations,omitempty"`
	Sizes     []string   `json:"sizes,omitempty"`
	Formats   []string   `json:"formats,omitempty"`
	Version   string     `json:"version,omitempty"`
	// The license for the resource. Use one of the SPDX license identifiers.
	License struct {
		ID  string `json:"id,omitempty"`
		Url string `json:"url,omitempty"`
	}
	Descriptions      []Description      `json:"descriptions,omitempty"`
	GeoLocations      []GeoLocation      `json:"geo_locations,omitempty"`
	FundingReferences []FundingReference `json:"funding_references,omitempty"`
	References        []Reference        `json:"references,omitempty"`

	// other properties
	DateCreated    string `json:"date_created,omitempty"`
	DateRegistered string `json:"date_registered,omitempty"`
	DatePublished  string `json:"date_published,omitempty"`
	DateUpdated    string `json:"date_updated,omitempty"`
	Files          []File `json:"files,omitempty"`
	Container      struct {
		Identifier     string `json:"identifier,omitempty"`
		IdentifierType string `json:"identifier_type,omitempty"`
		Type           string `json:"type,omitempty"`
		Title          string `json:"title,omitempty"`
		FirstPage      string `json:"first_page,omitempty"`
		LastPage       string `json:"last_page,omitempty"`
		Volume         string `json:"volume,omitempty"`
		Issue          string `json:"issue,omitempty"`
	}
	Provider struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}
	SchemaVersion string `json:"schema_version,omitempty"`
	// The location where content is archived.
	ArchiveLocations []string `json:"archive_locations,omitempty"`
	State            string   `json:"state,omitempty"`
}

type Affiliation struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Contributor struct {
	ID               string        `json:"id,omitempty"`
	Type             string        `json:"type"`
	ContributorRoles []string      `json:"contributor_roles"`
	Name             string        `json:"name,omitempty"`
	GivenName        string        `json:"given_name,omitempty"`
	FamilyName       string        `json:"family_name,omitempty"`
	Affiliations     []Affiliation `json:"affiliations,omitempty"`
}

type Description struct {
	Description string `json:"Description"`
	Type        string `json:"Type,omitempty"`
	Language    string `json:"Language,omitempty"`
}

type File struct {
	Bucket   string `json:"bucket,omitempty"`
	Key      string `json:"key,omitempty"`
	Checksum string `json:"checksum,omitempty"`
	Url      string `json:"url"`
	Size     int    `json:"size,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
}

type FundingReference struct {
	FunderIdentifier     string `json:"funder_identifier,omitempty"`
	FunderIdentifierType string `json:"funder_identifier_type,omitempty"`
	FunderName           string `json:"funder_name"`
	AwardNumber          string `json:"award_number,omitempty"`
	AwardUri             string `json:"award_uri,omitempty"`
}

type GeoLocation struct {
	GeoLocationPlace string           `json:"geo_location_place,omitempty"`
	GeoLocationPoint GeoLocationPoint `json:"geo_location_point,omitempty"`
	GeoLocationBox   struct {
		WestBoundLongitude float64 `json:"west_bound_longitude"`
		EastBoundLongitude float64 `json:"east_bound_longitude"`
		SouthBoundLatitude float64 `json:"south_bound_latitude"`
		NorthBoundLatitude float64 `json:"north_bound_latitude"`
	}
	GeoLocationPolygons []GeoLocationPolygon `json:"geo_location_polygon,omitempty"`
}

type GeoLocationPoint struct {
	PointLongitude float64 `json:"point_longitude,omitempty"`
	PointLatitude  float64 `json:"point_latitude,omitempty"`
}

type GeoLocationPolygon struct {
	PolygonPoints  []GeoLocationPoint `json:"polygon_points"`
	InPolygonPoint GeoLocationPoint   `json:"in_polygon_point,omitempty"`
}

type Relation struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Reference struct {
	Key             string `json:"key"`
	Doi             string `json:"doi,omitempty"`
	Contributor     string `json:"contributor,omitempty"`
	Title           string `json:"title,omitempty"`
	Publisher       string `json:"publisher,omitempty"`
	PublicationYear string `json:"publication_year,omitempty"`
	Volume          string `json:"volume,omitempty"`
	Issue           string `json:"issue,omitempty"`
	FirstPage       string `json:"first_page,omitempty"`
	LastPage        string `json:"last_page,omitempty"`
	ContainerTitle  string `json:"container_title,omitempty"`
	Edition         string `json:"edition,omitempty"`
	Unstructured    string `json:"unstructured,omitempty"`
}

type Subject struct {
	Subject string `json:"subject"`
}

type Title struct {
	Title    string `json:"title"`
	Type     string `json:"type,omitempty"`
	Language string `json:"language,omitempty"`
}
