package metadata

type Metadata struct {
	// required properties
	id           string
	_type        string
	url          string
	contributors []Contributor
	titles       []Title
	publisher    struct {
		id   string
		name string
	}
	date struct {
		created   string
		submitted string
		accepted  string
		published string
		updated   string
		accessed  string
		available string
		withdrawn string
	}

	// recommended and optional properties
	additional_type       string
	subjects              []string
	language              string
	alternate_identifiers []Identifier
	relations             []Relation
	sizes                 []string
	formats               []string
	version               string
	license               []string
	descriptions          []Description
	geo_locations         []GeoLocation
	funding_references    []FundingReference
	references            []Reference

	// other properties
	date_created    string
	date_registered string
	date_published  string
	date_updated    string
	files           []File
	container       struct {
		id   string
		name string
	}
	provider struct {
		id   string
		name string
	}
	schema_version    string
	archive_locations []string
	state             string
}
