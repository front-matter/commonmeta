package ror

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/fileutils"
	"github.com/front-matter/commonmeta/utils"
	"github.com/hamba/avro/v2"
	"gopkg.in/yaml.v3"
)

// ROR represents the ROR metadata record.
type ROR struct {
	ID            string        `avro:"id" json:"id" csv:"id"`
	Domains       Strings       `avro:"domains,omitempty" json:"domains,omitempty" yaml:"domains,omitempty"`
	Established   int           `avro:"established,omitempty" json:"established,omitempty" yaml:"established,omitempty"`
	ExternalIDs   ExternalIDS   `avro:"external_ids" json:"external_ids" yaml:"external_ids,omitempty"`
	Links         Links         `avro:"links" json:"links" yaml:"links,omitempty"`
	Locations     Locations     `avro:"locations" json:"locations"`
	Names         Names         `avro:"names" json:"names"`
	Relationships Relationships `avro:"relationships" json:"relationships" yaml:"relationships,omitempty"`
	Status        string        `avro:"status" json:"status"`
	Types         Strings       `avro:"types" json:"types"`
	Admin         Admin         `avro:"admin" json:"admin"`
}

type Strings []string
type ExternalIDS []ExternalID
type Links []Link
type Locations []Location
type Names []Name
type Relationships []Relationship

type Admin struct {
	Created      Date `avro:"created" json:"created"`
	LastModified Date `avro:"last_modified" json:"last_modified"`
}

type Date struct {
	Date          string `avro:"date" json:"date"`
	SchemaVersion string `avro:"schema_version" json:"schema_version"`
}

type ExternalID struct {
	Type      string  `avro:"type" json:"type"`
	All       Strings `avro:"all" json:"all"`
	Preferred string  `avro:"preferred,omitempty" json:"preferred,omitempty" yaml:"preferred,omitempty"`
}

type GeonamesDetails struct {
	ContinentCode          string  `avro:"continent_code" json:"continent_code" yaml:"continent_code"`
	ContinentName          string  `avro:"continent_name" json:"continent_name" yaml:"continent_name"`
	CountryCode            string  `avro:"country_code" json:"country_code" yaml:"country_code"`
	CountryName            string  `avro:"country_name" json:"country_name" yaml:"country_name"`
	CountrySubdivisionCode string  `avro:"country_subdivision_code,omitempty" json:"country_subdivision_code,omitempty" yaml:"country_subdivision_code,omitempty"`
	CountrySubdivisionName string  `avro:"country_subdivision_name,omitempty" json:"country_subdivision_name,omitempty" yaml:"country_subdivision_name,omitempty"`
	Lat                    float64 `avro:"lat" json:"lat"`
	Lng                    float64 `avro:"lng" json:"lng"`
	Name                   string  `avro:"name" json:"name"`
}

type Link struct {
	Type  string `avro:"type" json:"type"`
	Value string `avro:"value" json:"value"`
}

type Location struct {
	GeonamesID      int             `avro:"geonames_id" json:"geonames_id" yaml:"geonames_id"`
	GeonamesDetails GeonamesDetails `avro:"geonames_details" json:"geonames_details" yaml:"geonames_details"`
}

type Name struct {
	Value string  `avro:"value" json:"value"`
	Types Strings `avro:"types" json:"types"`
	Lang  string  `avro:"lang,omitempty" json:"lang,omitempty" yaml:"lang,omitempty"`
}

type Relationship struct {
	Type  string `avro:"type" json:"type"`
	Label string `avro:"label" json:"label"`
	ID    string `avro:"id" json:"id"`
}

// RORSchema is the Avro schema for the minimal ROR metadata.
var RORSchema = `{
  "type": "array",
  "items": {
    "name": "ROR",
    "type": "record",
    "fields": [
      { "name": "id", "type": "string" },
      { "name": "established", "type": ["null", "int"], "default": null },
      {
        "name": "external_ids",
        "type": {
          "type": "array",
          "items": {
            "name": "external_id",
            "type": "record",
            "fields": [
              {
                "name": "type",
                "type": {
                  "name": "external_id_type",
                  "type": "enum",
                  "symbols": ["fundref", "grid", "isni", "wikidata"]
                }
              },
              {
                "name": "all",
                "type": {
                  "type": "array",
                  "items": {
                    "name": "external_id",
                    "type": "string"
                  }
                }
              },
              {
                "name": "preferred",
                "type": ["null", "string"],
                "default": null
              }
            ]
          }
        }
      },
      {
        "name": "links",
        "type": {
          "type": "array",
          "items": {
            "name": "link",
            "type": "record",
            "fields": [
              {
                "name": "type",
                "type": {
                  "name": "link_type",
                  "type": "enum",
                  "symbols": ["website", "wikipedia"]
                }
              },
              { "name": "value", "type": "string" }
            ]
          }
        }
      },
      {
        "name": "locations",
        "type": {
          "type": "array",
          "items": {
            "name": "location",
            "type": "record",
            "fields": [
              { "name": "geonames_id", "type": "long" },
              {
                "name": "geonames_details",
                "type": {
                  "name": "geonames_details",
                  "type": "record",
                  "fields": [
                    { "name": "continent_code", "type": "string" },
                    { "name": "continent_name", "type": "string" },
                    { "name": "country_code", "type": "string" },
                    { "name": "country_name", "type": "string" },
                    {
                      "name": "country_subdivision_code",
                      "type": ["null", "string"],
                      "default": null
                    },
                    {
                      "name": "country_subdivision_name",
                      "type": ["null", "string"],
                      "default": null
                    },
                    { "name": "lat", "type": "double" },
                    { "name": "lng", "type": "double" },
                    { "name": "name", "type": "string" }
                  ]
                }
              }
            ]
          }
        }
      },
      {
        "name": "names",
        "type": {
          "type": "array",
          "items": {
            "name": "name",
            "type": "record",
            "fields": [
              { "name": "value", "type": "string" },
              {
                "name": "types",
                "type": {
                  "type": "array",
                  "items": {
                    "name": "name_type",
                    "type": "enum",
                    "symbols": ["acronym", "alias", "label", "ror_display"]
                  }
                }
              },
              { "name": "lang", "type": ["null", "string"], "default": null }
            ]
          }
        }
      },
      {
        "name": "relationships",
        "type": {
          "type": "array",
          "items": {
            "name": "relationship",
            "type": "record",
            "fields": [
              {
                "name": "type",
                "type": {
                  "name": "relationship_type",
                  "type": "enum",
                  "symbols": [
                    "child",
                    "parent",
                    "related",
                    "predecessor",
                    "successor"
                  ]
                }
              },
              { "name": "label", "type": "string" },
              { "name": "id", "type": "string" }
            ]
          }
        }
      },
      {
        "name": "status",
        "type": "string"
      },
      {
        "name": "types",
        "type": {
          "name": "type",
          "type": "array",
          "items": {
            "name": "type",
            "type": "enum",
            "symbols": [
              "archive",
              "company",
              "education",
              "facility",
              "funder",
              "government",
              "healthcare",
              "nonprofit",
              "other"
            ]
          }
        }
      },
      {
        "name": "admin",
        "type": {
          "name": "admin",
          "type": "record",
          "fields": [
            {
              "name": "created",
              "type": {
                "name": "created",
                "type": "record",
                "fields": [
                  { "name": "date", "type": "string" },
                  {
                    "name": "schema_version",
                    "type": "string"
                  }
                ]
              }
            },
            {
              "name": "last_modified",
              "type": {
                "name": "last_modified",
                "type": "record",
                "fields": [
                  { "name": "date", "type": "string" },
                  {
                    "name": "schema_version",
                    "type": "string"
                  }
                ]
              }
            }
          ]
        }
      }
    ]
  }
}`

// RORVersions contains the ROR versions and their release dates, published on Zenodo.
// The ROR version is the first part of the filename, e.g., v1.63-2025-04-03-ror-data_schema_v2.json
// Beginning with release v1.45 on 11 April 2024, data releases contain JSON and CSV files formatted
// according to both schema v1 and schema v2. Version 2 files have _schema_v2 appended to the end of
// the filename, e.g., v1.45-2024-04-11-ror-data_schema_v2.json.
var RORVersions = map[string]string{
	"v1.50": "2024-07-29",
	"v1.51": "2024-08-21",
	"v1.52": "2024-09-16",
	"v1.53": "2023-10-14",
	"v1.54": "2024-10-21",
	"v1.55": "2024-10-31",
	"v1.56": "2024-11-19",
	"v1.58": "2024-12-11",
	"v1.59": "2025-01-23",
	"v1.60": "2025-02-27",
	"v1.61": "2025-03-18",
	"v1.62": "2025-03-27",
	"v1.63": "2025-04-03",
}

var SupportedTypes = []string{"ROR", "Wikidata", "Crossref Funder ID", "GRID", "ISNI"}
var RORTypes = []string{"archive", "company", "education", "facility", "funder", "government", "healthcare", "nonprofit", "other"}
var Extensions = []string{".avro", ".yaml", ".json"}

// Fetch fetches ROR metadata for a given ror id.
func Fetch(str string) (ROR, error) {
	var data ROR
	id, ok := utils.ValidateROR(str)
	if !ok {
		return data, errors.New("invalid ror id")
	}
	data, err := Get(id)
	return data, err
}

// Get gets ROR metadata for a given ror id.
func Get(id string) (ROR, error) {
	var data ROR

	id, ok := utils.ValidateROR(id)
	if !ok {
		return data, errors.New("invalid ror id")
	}
	url := "https://api.ror.org/v2/organizations/" + id
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(url)
	if err != nil {
		return data, err
	}
	if resp.StatusCode >= 400 {
		return data, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &data)
	return data, err
}

// MatchAffiliation searches ROR metadata for a given affiliation name, using their
// matching strategies.
func MatchAffiliation(name string) (ROR, error) {
	type Response struct {
		Substring    string  `json:"substring"`
		Score        float32 `json:"score"`
		MatchingType string  `json:"matching_type"`
		Chosen       bool    `json:"chosen"`
		Organization ROR     `json:"organization"`
	}

	// Content is the wrapper around the response from the ROR API
	type Content struct {
		NumberOfResults int        `json:"number_of_results"`
		Items           []Response `json:"items"`
	}

	var content Content
	var data ROR

	url := "https://api.ror.org/v2/organizations?affiliation=" + url.QueryEscape(name)
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(url)
	if err != nil {
		return data, err
	}
	if resp.StatusCode >= 400 {
		return data, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(body, &content)
	if err != nil {
		fmt.Println(err)
		return data, errors.New("error unmarshalling response")
	}

	// Check if there is a chosen organization in the response
	chosen := slices.IndexFunc(content.Items, func(d Response) bool { return d.Chosen })
	if chosen != -1 {
		// fmt.Println("Chosen:", content.Items[chosen].MatchingType, content.Items[chosen].Score, content.Items[chosen].Organization.Names[0].Value)
		data = content.Items[chosen].Organization
		return data, nil
	}
	return data, err
}

// Search searches local ROR metadata for a given ror id,
// Crossref Funder ID, grid ID, or Wikidata ID.
func Search(id string) (ROR, error) {
	var idx int
	var ror ROR

	pid, type_ := utils.ValidateID(id)
	if !slices.Contains(SupportedTypes, type_) {
		return ror, errors.New("not a supported organization id")
	}

	data, err := LoadBuiltin()
	if err != nil {
		return ror, err
	}
	if type_ == "ISNI" {
		// ROR expects ISNI IDs to be in the form of 0000 0002 1234 5678
		pid = utils.SplitString(pid, 4, " ")
		pid = strings.ReplaceAll(pid, "-", " ")
	}
	if type_ == "ROR" {
		idx = slices.IndexFunc(data, func(d ROR) bool { return d.ID == utils.NormalizeROR(pid) })
	} else {
		idx = slices.IndexFunc(data, func(d ROR) bool {
			for _, e := range d.ExternalIDs {
				for _, all := range e.All {
					if all == pid {
						return true
					}
				}
			}
			return false
		})
	}
	if idx == -1 {
		return ror, errors.New("no organization found")
	}

	ror = data[idx]
	return ror, err
}

// LoadAll loads the metadata for a list of organizations from a ROR JSON file
func LoadAll(filename string) ([]ROR, error) {
	var data []ROR

	extension := path.Ext(filename)
	if !slices.Contains(Extensions, extension) {
		return data, errors.New("invalid file extension")
	}
	output, err := fileutils.ReadFile(filename)
	if err != nil {
		return data, errors.New("error reading file")
	}

	if extension == ".avro" {
		schema, err := avro.Parse(RORSchema)
		if err != nil {
			return nil, err
		}
		err = avro.Unmarshal(schema, output, &data)
		if err != nil {
			fmt.Println(err)
			return data, errors.New("error unmarshalling avro file")
		}
	} else if extension == ".json" {
		err = json.Unmarshal(output, &data)
		if err != nil {
			return data, errors.New("error unmarshalling json file")
		}
	} else if extension == ".yaml" {
		err = yaml.Unmarshal(output, &data)
		if err != nil {
			return data, errors.New("error unmarshalling yaml file")
		}
	}
	return data, err
}

// LoadBuiltin loads the embedded ROR metadata from the ZIP file with all ROR records.
func LoadBuiltin() ([]ROR, error) {
	var data []ROR
	var err error
	schema, err := avro.Parse(RORSchema)
	if err != nil {
		return nil, err
	}

	output, err := fileutils.ReadZIPFile("v1.63-2025-04-03-ror-data.avro.zip")
	if err != nil {
		return nil, err
	}
	err = avro.Unmarshal(schema, output, &data)
	if err != nil {
		return nil, err
	}
	return data, err
}

// ExtractAll extracts ROR metadata from a JSON file in commonmeta format.
func ExtractAll(content []commonmeta.Data) ([]byte, error) {
	var data []ROR
	var extracted []ROR
	var ids []string
	var err error
	schema, err := avro.Parse(RORSchema)
	if err != nil {
		return nil, err
	}

	// Load the ROR metadata from the embedded ZIP file with all ROR records
	out, err := fileutils.ReadZIPFile("ror.avro.zip")
	if err != nil {
		return nil, err
	}
	err = avro.Unmarshal(schema, out, &data)
	if err != nil {
		return nil, err
	}

	// Extract ROR IDs from the content
	for _, v := range content {
		if len(v.Contributors) > 0 {
			for _, c := range v.Contributors {
				if len(c.Affiliations) > 0 {
					for _, a := range c.Affiliations {
						if a.ID != "" && !slices.Contains(ids, a.ID) {
							id, _ := utils.ValidateROR(a.ID)
							idx := slices.IndexFunc(data, func(d ROR) bool { return d.ID == id })
							if idx != -1 {
								ids = append(ids, a.ID)
								extracted = append(extracted, data[idx])
							}
						}
					}
				}
			}
		}
	}

	output, err := avro.Marshal(schema, extracted)
	return output, err
}
