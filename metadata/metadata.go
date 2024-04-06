package metadata

import (
	"commonmeta/cff"
	"commonmeta/codemeta"
	"commonmeta/crossref"
	"commonmeta/crossrefxml"
	"commonmeta/datacite"
	"commonmeta/inveniordm"
	"commonmeta/jsonfeed"
	"commonmeta/schemaorg"
	"commonmeta/utils"
)

type Metadata struct {
	// required properties
	ID           string        `json:"id"`
	Type         string        `json:"type"`
	Url          string        `json:"url,omitempty"`
	Contributors []Contributor `json:"contributors,omitempty"`
	Titles       []Title       `json:"titles,omitempty"`
	Publisher    Publisher     `json:"publisher,omitempty"`
	Date         Date          `json:"date,omitempty"`
	// recommended and optional properties
	AdditionalType string    `json:"additional_name,omitempty"`
	Subjects       []Subject `json:"subjects,omitempty"`
	// The language of the resource. Use one of the language codes from the IETF BCP 47 standard.
	Language             string                `json:"language,omitempty"`
	AlternateIdentifiers []AlternateIdentifier `json:"alternate_identifiers,omitempty"`
	Relations            []Relation            `json:"relations,omitempty"`
	Sizes                []string              `json:"sizes,omitempty"`
	Formats              []string              `json:"formats,omitempty"`
	Version              string                `json:"version,omitempty"`
	// The license for the resource. Use one of the SPDX license identifiers.
	License           License            `json:"license,omitempty"`
	Descriptions      []Description      `json:"descriptions,omitempty"`
	GeoLocations      []GeoLocation      `json:"geo_locations,omitempty"`
	FundingReferences []FundingReference `json:"funding_references,omitempty"`
	References        []Reference        `json:"references,omitempty"`

	// other properties
	DateCreated    string    `json:"date_created,omitempty"`
	DateRegistered string    `json:"date_registered,omitempty"`
	DatePublished  string    `json:"date_published,omitempty"`
	DateUpdated    string    `json:"date_updated,omitempty"`
	Files          []File    `json:"files,omitempty"`
	Container      Container `json:"container,omitempty"`
	Provider       string    `json:"provider,omitempty"`
	SchemaVersion  string    `json:"schema_version,omitempty"`
	// The location where content is archived.
	ArchiveLocations []string `json:"archive_locations,omitempty"`
	State            string   `json:"state,omitempty"`
	Via              string   `json:"via,omitempty"`
}

type Affiliation struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type AlternateIdentifier struct {
	AlternateIdentifier     string `json:"alternate_identifier"`
	AlternateIdentifierType string `json:"alternate_identifier_type"`
}

type Container struct {
	Identifier     string `json:"identifier,omitempty"`
	IdentifierType string `json:"identifier_type,omitempty"`
	Type           string `json:"type,omitempty"`
	Title          string `json:"title,omitempty"`
	FirstPage      string `json:"first_page,omitempty"`
	LastPage       string `json:"last_page,omitempty"`
	Volume         string `json:"volume,omitempty"`
	Issue          string `json:"issue,omitempty"`
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

type Date struct {
	Created   string `json:"created,omitempty"`
	Submitted string `json:"submitted,omitempty"`
	Accepted  string `json:"accepted,omitempty"`
	Published string `json:"published,omitempty"`
	Updated   string `json:"updated,omitempty"`
	Accessed  string `json:"accessed,omitempty"`
	Available string `json:"available,omitempty"`
	Withdrawn string `json:"withdrawn,omitempty"`
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

type License struct {
	ID  string `json:"id,omitempty"`
	Url string `json:"url,omitempty"`
}

type Publisher struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
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

func NewMetadata(str string, via string) *Metadata {
	pid := utils.NormalizeID(str)
	// p := utils.Params{
	// 	Pid: pid,
	// }
	// via := utils.FindFromFormat(p)
	return &Metadata{ID: pid, Type: "", Via: via}
}

func (m *Metadata) GetMetadata(pid string, str string) map[string]interface{} {
	via := m.Via
	if pid != "" {
		if via == "schema_org" {
			return schemaorg.GetSchemaOrg(pid)
		} else if via == "datacite" {
			return datacite.GetDatacite(pid)
		} else if via == "crossref" || via == "op" {
			return crossref.GetCrossref(pid)
		} else if via == "crossref_xml" {
			return crossrefxml.GetCrossrefXML(pid)
		} else if via == "codemeta" {
			return codemeta.GetCodemeta(pid)
		} else if via == "cff" {
			return cff.GetCFF(pid)
		} else if via == "json_feed_item" {
			return jsonfeed.GetJsonFeedItem(pid)
		} else if via == "inveniordm" {
			return inveniordm.GetInvenioRDM(pid)
		}
	} else if str != "" {
		// if via == "datacite_xml" {
		// 	return ParseXML(str)
		// } else if via == "crossref_xml" {
		// 	return ParseXML(str, "crossref")
		// } else if via == "cff" {
		// 	return ParseYAML(str)
		// } else if via == "bibtex" {
		// 	panic("Bibtex not supported")
		// } else if via == "ris" {
		// 	return ParseRIS(str)
		// } else if via == "commonmeta" || via == "crossref" || via == "datacite" || via == "schema_org" || via == "csl" || via == "json_feed_item" || via == "codemeta" || via == "kbase" || via == "inveniordm" {
		// 	return ParseJSON(str)
		// } else {
		// 	panic("No input format found")
		// }
	} else {
		panic("No metadata found")
	}
	return nil
}

// def get_metadata(self, pid, string) -> dict:
// via = self.via
// if pid is not None:
// 		if via == "schema_org":
// 				return get_schema_org(pid)
// 		elif via == "datacite":
// 				return get_datacite(pid)
// 		elif via in ["crossref", "op"]:
// 				return get_crossref(pid)
// 		elif via == "crossref_xml":
// 				return get_crossref_xml(pid)
// 		elif via == "codemeta":
// 				return get_codemeta(pid)
// 		elif via == "cff":
// 				return get_cff(pid)
// 		elif via == "json_feed_item":
// 				return get_json_feed_item(pid)
// 		elif via == "inveniordm":
// 				return get_inveniordm(pid)
// elif string is not None:
// 		if via == "datacite_xml":
// 				return parse_xml(string)
// 		elif via == "crossref_xml":
// 				return parse_xml(string, dialect="crossref")
// 		elif via == "cff":
// 				return yaml.safe_load(string)
// 		elif via == "bibtex":
// 				raise ValueError("Bibtex not supported")
// 		elif via == "ris":
// 				return string
// 		elif via in [
// 				"commonmeta",
// 				"crossref",
// 				"datacite",
// 				"schema_org",
// 				"csl",
// 				"json_feed_item",
// 				"codemeta",
// 				"kbase",
// 				"inveniordm",
// 		]:
// 				return json.loads(string)
// 		else:
// 				raise ValueError("No input format found")
// else:
// 		raise ValueError("No metadata found")

// def read_metadata(self, data: dict, **kwargs) -> dict:
// """get_metadata"""
// via = isinstance(data, dict) and data.get("via", None) or self.via
// if via == "commonmeta":
// 		return read_commonmeta(data, **kwargs)
// elif via == "schema_org":
// 		return read_schema_org(data)
// elif via == "datacite":
// 		return read_datacite(data)
// elif via == "datacite_xml":
// 		return read_datacite_xml(data)
// elif via in ["crossref", "op"]:
// 		return read_crossref(data)
// elif via == "crossref_xml":
// 		return read_crossref_xml(data)
// elif via == "csl":
// 		return read_csl(data, **kwargs)
// elif via == "codemeta":
// 		return read_codemeta(data)
// elif via == "cff":
// 		return read_cff(data)
// elif via == "json_feed_item":
// 		return read_json_feed_item(data, **kwargs)
// elif via == "inveniordm":
// 		return read_inveniordm(data)
// elif via == "kbase":
// 		return read_kbase(data)
// elif via == "ris":
// 		return read_ris(data)
// else:
// 		raise ValueError("No input format found")
