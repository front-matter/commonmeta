package ror

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"os"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/fileutils"
	"github.com/front-matter/commonmeta/utils"
	"github.com/front-matter/commonmeta/vocabularies"
	"github.com/hamba/avro/v2"
	"gopkg.in/yaml.v3"
)

// ROR represents the ROR metadata record.
type ROR struct {
	ID            string        `avro:"id" json:"id" csv:"id"`
	Domains       Strings       `avro:"domains,omitempty" json:"domains,omitempty" yaml:"domains,omitempty"`
	Established   int           `avro:"established" json:"established,omitempty" yaml:"established,omitempty"`
	ExternalIDs   ExternalIDS   `avro:"external_ids" json:"external_ids,omitempty" yaml:"external_ids,omitempty"`
	Links         Links         `avro:"links" json:"links" yaml:"links,omitempty"`
	Locations     Locations     `avro:"locations" json:"locations"`
	Names         Names         `avro:"names" json:"names"`
	Relationships Relationships `avro:"relationships" json:"relationships,omitempty" yaml:"relationships,omitempty"`
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
	Preferred string  `avro:"preferred" json:"preferred,omitempty" yaml:"preferred,omitempty"`
}

type GeonamesDetails struct {
	ContinentCode          string  `avro:"continent_code" json:"continent_code" yaml:"continent_code"`
	ContinentName          string  `avro:"continent_name" json:"continent_name" yaml:"continent_name"`
	CountryCode            string  `avro:"country_code" json:"country_code" yaml:"country_code"`
	CountryName            string  `avro:"country_name" json:"country_name" yaml:"country_name"`
	CountrySubdivisionCode string  `avro:"country_subdivision_code" json:"country_subdivision_code,omitempty" yaml:"country_subdivision_code,omitempty"`
	CountrySubdivisionName string  `avro:"country_subdivision_name" json:"country_subdivision_name,omitempty" yaml:"country_subdivision_name,omitempty"`
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
	Lang  string  `avro:"lang" json:"lang,omitempty" yaml:"lang,omitempty"`
}

type Relationship struct {
	Type  string `avro:"type" json:"type"`
	Label string `avro:"label" json:"label"`
	ID    string `avro:"id" json:"id"`
}

// RORSchema is the Avro schema for the minimal ROR metadata.
var RORSchema = `{
  "type": "map",
  "values": {
    "name": "ROR",
    "type": "record",
    "fields": [
      { "name": "id", "type": "string" },
      { "name": "established", "type": "int", "default": 0 },
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
                "type": "string",
                "default": ""
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
                      "type": "string",
                      "default": ""
                    },
                    {
                      "name": "country_subdivision_name",
                      "type": "string",
                      "default": ""
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
              { "name": "lang", "type": "string", "default": "" }
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
	"v1.59": "2025-01-23",
	"v1.60": "2025-02-27",
	"v1.61": "2025-03-18",
	"v1.62": "2025-03-27",
	"v1.63": "2025-04-03",
}

// RORZenodoIDs contains the Zenodo IDs for the ROR data releases.
var RORZenodoIDs = map[string]string{
	"v1.59": "14728473",
	"v1.60": "14797924",
	"v1.61": "15047759",
	"v1.62": "15098078",
	"v1.63": "15132361",
}

var ArchivedFilename = "v1.63-2025-04-03-ror-data_schema_v2.json"
var RORFilename = "v1.63-2025-04-03-ror-data.zip"
var RORDownloadURL = fmt.Sprintf("https://zenodo.org/records/15132361/files/%s?download=1", RORFilename)
var RORAvroFilename = "v1.63-2025-04-03-ror-data.avro"

var SupportedTypes = []string{"ROR", "Wikidata", "Crossref Funder ID", "GRID", "ISNI"}
var RORTypes = []string{"archive", "company", "education", "facility", "funder", "government", "healthcare", "nonprofit", "other"}
var Extensions = []string{".avro", ".yaml", ".json", ".jsonl", ".csv"}

// Fetch fetches ROR metadata for a given ror id.
func Fetch(str string) (ROR, error) {
	var ror ROR

	ror, err := Get(str)
	return ror, err
}

// Get gets ROR metadata for a given ror id.
func Get(str string) (ROR, error) {
	// Content is the wrapper around the response from the ROR API
	type Content struct {
		NumberOfResults int   `json:"number_of_results"`
		Items           []ROR `json:"items"`
	}

	var ror ROR
	var content Content
	var url_ string

	id, type_ := utils.ValidateID(str)
	if !slices.Contains(SupportedTypes, type_) {
		return ror, errors.New("not a supported organization id")
	}
	if type_ == "ROR" {
		url_ = "https://api.ror.org/v2/organizations/" + id
	} else {
		url_ = "https://api.ror.org/v2/organizations?query=" + url.QueryEscape(id)
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(url_)
	if err != nil {
		return ror, err
	}
	if resp.StatusCode >= 400 {
		return ror, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ror, err
	}
	if type_ == "ROR" {
		err = json.Unmarshal(body, &ror)
	} else {
		err = json.Unmarshal(body, &content)
		if content.NumberOfResults == 1 {
			ror = content.Items[0]
		}
	}
	if err != nil {
		return ror, errors.New("error unmarshalling response")
	}
	return ror, err
}

// MatchOrganization searches ROR metadata for a given affiliation name, using their
// matching strategies.
func MatchOrganization(name string) (ROR, error) {
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

	// fmt.Println("Number of results:", content.NumberOfResults)

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

	catalog, err := LoadBuiltin()
	if err != nil {
		return ror, err
	}
	if type_ == "ISNI" {
		// ROR expects ISNI IDs to be in the form of 0000 0002 1234 5678
		pid = utils.SplitString(pid, 4, " ")
		pid = strings.ReplaceAll(pid, "-", " ")
	}
	if type_ == "ROR" {
		return catalog[utils.NormalizeROR(pid)], nil
	}
	// convert map into a slice to search for value
	list := slices.Collect(maps.Values(catalog))
	idx = slices.IndexFunc(list, func(d ROR) bool {
		for _, e := range d.ExternalIDs {
			for _, all := range e.All {
				if all == pid {
					return true
				}
			}
		}
		return false
	})
	if idx == -1 {
		return ror, errors.New("no organization found")
	}

	ror = list[idx]
	return ror, err
}

// Basename returns the basename of the ROR data dump file for a given version.
func Basename(version string) string {
	date := RORVersions[version]
	if date == "" {
		return ""
	}
	basename := fmt.Sprintf("%s-%s-ror-data", version, date)
	return basename
}

// FetchAll fetches the ROR Data dump from Zenodo.
func FetchAll(version string) (map[string]ROR, error) {
	var input, output []byte
	var list []ROR
	var err error

	// constuct the URL for the ROR data dump
	basename := Basename(version)
	if basename == "" {
		return nil, fmt.Errorf("invalid version: %s", version)
	}
	zenodoID := RORZenodoIDs[version]
	zipname := fmt.Sprintf("%s.zip", basename)
	jsonname := fmt.Sprintf("%s_schema_v2.json", basename)
	url := fmt.Sprintf("https://zenodo.org/records/%s/files/%s?download=1", zenodoID, zipname)

	// download the ROR data zip file
	input, err = fileutils.DownloadFile(url)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error downloading zip file")
	}

	// unzip the json version 2 in the ROR data zip file
	output, err = fileutils.UnzipContent(input, jsonname)
	if err != nil {
		return nil, fmt.Errorf("error unzipping file: %w", err)
	}

	// write the unzipped json file to disk
	err = fileutils.WriteFile(jsonname, output)
	if err != nil {
		return nil, fmt.Errorf("error writing json file: %w", err)
	}

	err = json.Unmarshal(output, &list)
	if err != nil {
		return nil, fmt.Errorf("error parsing json data: %w", err)
	}

	// convert to map
	catalog := make(map[string]ROR)
	for _, v := range list {
		catalog[v.ID] = v
	}
	return catalog, nil
}

// LoadAll loads the metadata for a list of organizations from a ROR file.
func LoadAll(filename string) (map[string]ROR, error) {
	var list []ROR
	var catalog = make(map[string]ROR)
	var output []byte
	var err error

	filename, extension, compress := fileutils.GetExtension(filename, ".json")

	if !slices.Contains(Extensions, extension) {
		return catalog, errors.New("invalid file extension")
	}
	if compress {
		output, err = fileutils.ReadZIPFile(filename+".zip", path.Base(filename))
		if err != nil {
			return catalog, errors.New("error reading zip file")
		}
	} else {
		output, err = fileutils.ReadFile(filename)
		if err != nil {
			return catalog, errors.New("error reading file")
		}
	}
	if extension == ".avro" {
		schema, err := avro.Parse(RORSchema)
		if err != nil {
			return nil, err
		}
		err = avro.Unmarshal(schema, output, &catalog)
		if err != nil {
			fmt.Println(err)
			return catalog, errors.New("error unmarshalling avro file")
		}
	} else if extension == ".json" {
		err = json.Unmarshal(output, &list)
		if err != nil {
			return catalog, errors.New("error unmarshalling json file")
		}
		for _, v := range list {
			catalog[v.ID] = v
		}
	} else if extension == ".yaml" {
		err = yaml.Unmarshal(output, &list)
		if err != nil {
			return catalog, errors.New("error unmarshalling yaml file")
		}
		for _, v := range list {
			catalog[v.ID] = v
		}
	} else if extension == ".jsonl" {
		decoder := json.NewDecoder(bytes.NewReader(output))
		for {
			var item ROR
			if err := decoder.Decode(&item); err == io.EOF {
				break
			} else if err != nil {
				return catalog, errors.New("error unmarshalling jsonl file")
			}
			catalog[item.ID] = item
		}
	}
	return catalog, err
}

// MapROR maps between a ROR ID and organization name
//
// The function accepts a ROR ID and/or name and returns both values if possible:
// - If both ID and name are provided, they are returned unchanged
// - If only ID is provided, the name is fetched from ROR API
// - If only name is provided and match=true, attempts to find a matching ROR ID
func MapROR(id string, name string, assertedBy string, match bool) (string, string, string, error) {
	// Both ID and name provided, nothing to do
	if id != "" && name != "" {
		return id, name, assertedBy, nil
	}

	// Only ID provided, fetch the name
	if id != "" {
		ror, err := Fetch(id)
		if err != nil {
			return id, "", assertedBy, err
		}
		return id, GetDisplayName(ror), assertedBy, nil
	}

	// Only name provided and matching requested
	// Do not match against ROR for names that are known to be missing
	exceptions := []string{"Front Matter"}
	if name != "" && !slices.Contains(exceptions, name) && match {
		ror, err := MatchOrganization(name)
		if err != nil {
			return "", name, assertedBy, err
		}
		if ror.ID == "" {
			return "", name, assertedBy, nil
		}
		return ror.ID, name, "ror", nil
	}

	// No useful input or matching not requested
	return id, name, assertedBy, nil
}

// LoadBuiltin loads the ROR metadata from the embedded ROR catalog in zipped avro format.
func LoadBuiltin() (map[string]ROR, error) {
	var catalog map[string]ROR

	bytes, err := vocabularies.LoadVocabulary("ROR.Organizations")
	schema, err := avro.Parse(RORSchema)
	if err != nil {
		return nil, fmt.Errorf("error parsing avro schema: %w", err)
	}
	err = avro.Unmarshal(schema, bytes, &catalog)
	if err != nil {
		return nil, fmt.Errorf("error parsing avro data: %w", err)
	}

	return catalog, nil
}

// LoadJSON loads a local ROR json file with the specified version.
// If the file does not exist or is smaller than 50 MB, it will be downloaded from Zenodo.
func LoadJSON(version string) ([]ROR, error) {
	var bytes []byte
	var list []ROR

	date := RORVersions[version]
	if date == "" {
		return nil, fmt.Errorf("invalid version: %s", version)
	}
	filename := fmt.Sprintf("%s-%s-ror-data_schema_v2.json", version, date)

	// check that local ROR data json exists and is large enough
	info, err := os.Stat(filename)
	if err == nil && (info.Size()/1048576) > 50 {
		// read json ZIP file
		bytes, err = fileutils.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("error reading json file: %w", err)
		}
	} else {
		bytes, err = fileutils.DownloadFile(RORDownloadURL)
		if err != nil {
			return nil, fmt.Errorf("error downloading json file")
		}
	}

	err = json.Unmarshal(bytes, &list)
	if err != nil {
		return nil, fmt.Errorf("error parsing json data: %w", err)
	}

	return list, nil
}

// GetDisplayName returns the display name of the organization
func GetDisplayName(ror ROR) string {
	for _, name := range ror.Names {
		if slices.Contains(name.Types, "ror_display") {
			return name.Value
		}
	}
	return ""
}

// ExtractAll extracts ROR metadata from a JSON file in commonmeta format.
func ExtractAll(content []commonmeta.Data) ([]byte, error) {
	var extracted map[string]ROR
	var ids []string

	schema, err := avro.Parse(RORSchema)
	if err != nil {
		return nil, err
	}

	// Load the ROR metadata from the embedded ZIP file with all ROR records
	catalog, err := LoadBuiltin()
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
							item, ok := catalog[id]
							if ok {
								extracted[id] = item
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
