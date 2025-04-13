package ror

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path"
	"slices"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/fileutils"
	"github.com/front-matter/commonmeta/utils"
	"github.com/hamba/avro/v2"
)

// ROR represents the minimal ROR metadata record.
type ROR struct {
	ID        string     `avro:"id" json:"id"`
	Locations []Location `avro:"locations" json:"locations"`
	Names     []Name     `avro:"names" json:"names"`
	Types     []string   `avro:"types" json:"types"`
}

// Content represents the full ROR metadata record.
type Content struct {
	*ROR
	Established   int            `json:"established"`
	ExternalIDs   []ExternalID   `json:"external_ids"`
	Links         []Link         `json:"links"`
	Relationships []Relationship `json:"relationships"`
	Status        string         `json:"status"`
	Admin         struct {
		Created struct {
			Date          string `json:"date"`
			SchemaVersion string `json:"schema_version"`
		} `json:"created"`
		LastModified struct {
			Date          string `json:"date"`
			SchemaVersion string `json:"schema_version"`
		} `json:"last_modified"`
	}
}

type ExternalID struct {
	Type      string   `json:"type"`
	All       []string `json:"all"`
	Preferred string   `json:"preferred"`
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
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Location struct {
	GeonamesID      int             `avro:"geonames_id" json:"geonames_id"`
	GeonamesDetails GeonamesDetails `avro:"geonames_details" json:"geonames_details"`
}

type Name struct {
	Value string   `avro:"value" json:"value"`
	Types []string `avro:"types" json:"types"`
	Lang  string   `avro:"lang,omitempty" json:"lang,omitempty" yaml:"lang,omitempty"`
}

type Relationship struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	ID    string `json:"id"`
}

// RORSchema is the Avro schema for the minimal ROR metadata.
var RORSchema = `{
  "type": "array",
  "items": {
    "name": "ROR",
    "type": "record",
    "fields": [
      { "name": "id", "type": "string" },
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
              {
                "name": "lang",
                "type": ["null", "string"], "default": null
              }
            ]
          }
        }
      },
      {
        "name": "types",
        "type": {
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
      }
    ]
  }
}`

// ContentSchema is the Avro schema for the complete ROR metadata.
var ContentSchema = `{
  "type": "array",
  "items": {
    "name": "ROR",
    "type": "record",
    "fields": [
      { "name": "established", "type": "int" },
      {
        "name": "external_ids",
        "type": {
          "type": "array",
          "items": {
            "type": [
              {
                "name": "external_id",
                "type": {
                  "name": "external_id",
                  "type": "record",
                  "fields": [
                    {
                      "type": "enum",
                      "name": "type",
                      "symbols": ["fundref", "grid", "isni", "wikidata"]
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
                    { "name": "preferred", "type": "string" }
                  ]
                }
              }
            ]
          }
        }
      },
      { "name": "id", "type": "string" },
      {
        "name": "links",
        "type": {
          "type": "array",
          "items": {
            "name": "link",
            "type": {
              "name": "link",
              "type": "record",
              "fields": [
                {
                  "type": "enum",
                  "name": "type",
                  "symbols": ["website", "wikipedia"]
                },
                { "name": "value", "type": "string" }
              ]
            }
          }
        }
      },
      {
        "name": "locations",
        "type": {
          "type": "array",
          "items": {
            "name": "location",
            "type": {
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
                        "type": ["null", "string"]
                      },
                      {
                        "name": "country_subdivision_name",
                        "type": ["null", "string"]
                      },
                      { "name": "lat", "type": ["null", "double"] },
                      { "name": "lng", "type": ["null", "double"] },
                      { "name": "name", "type": "string" }
                    ]
                  }
                }
              ]
            }
          }
        }
      },
      {
        "name": "names",
        "type": {
          "type": "array",
          "items": {
            "name": "name",
            "type": {
              "name": "name",
              "type": "record",
              "fields": [
                { "name": "value", "type": "string" },
                {
                  "type": "enum",
                  "name": "type",
                  "symbols": ["acronym", "alias", "label", "ror_display"]
                },
                { "name": "lang", "type": ["null", "string"], "default": null }
              ]
            }
          }
        }
      },
      {
        "name": "relationships",
        "type": {
          "type": "array",
          "items": {
            "name": "relationship",
            "type": {
              "name": "relationship",
              "type": "record",
              "fields": [
                {
                  "type": "enum",
                  "name": "type",
                  "symbols": ["child", "parent", "related"]
                },
                { "name": "label", "type": "string" },
                { "name": "id", "type": "string" }
              ]
            }
          }
        }
      },
      { "name": "status", "type": "enum", "symbols": ["active"] },
      {
        "name": "types",
        "type": {
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

var RORTypes = []string{"archive", "company", "education", "facility", "funder", "government", "healthcare", "nonprofit", "other"}

// LoadAll loads the metadata for a list of organizations from a ROR JSON file
func LoadAll(filename string, type_ string, country string) ([]ROR, error) {
	var data []ROR
	var content []Content
	var err error

	extension := path.Ext(filename)
	if extension == ".json" {
		file, err := os.Open(filename)
		if err != nil {
			return data, errors.New("error reading file")
		}
		defer file.Close()

		decoder := json.NewDecoder(file)
		err = decoder.Decode(&content)
		if err != nil {
			return data, err
		}
	} else if extension != ".json" {
		return data, errors.New("invalid file extension")
	} else {
		return data, errors.New("unsupported file format")
	}

	data, err = ReadAll(content)
	if err != nil {
		return data, err
	}
	return data, nil
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

// Read reads ROR full metadata and converts it into InvenioRDM metadata.
func Read(content Content) (ROR, error) {
	var data ROR

	data.ID = content.ID
	data.Locations = content.Locations
	data.Names = content.Names
	data.Types = content.Types

	return data, nil
}

// ReadAll reads a list of ROR JSON organizations
func ReadAll(content []Content) ([]ROR, error) {
	var data []ROR

	for _, v := range content {
		d, err := Read(v)
		if err != nil {
			log.Println(err)
		}
		data = append(data, d)
	}
	return data, nil
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
