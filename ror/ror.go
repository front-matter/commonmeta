// Package ror converts ROR (Research Organization Registry) metadata.
package ror

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"

	"gopkg.in/yaml.v3"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/fileutils"
	"github.com/front-matter/commonmeta/utils"
)

// ROR represents the minimal ROR metadata record.
type ROR struct {
	ID    string `json:"id"`
	Names []Name `json:"names"`
	Admin struct {
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

// Content represents the full ROR metadata record.
type Content struct {
	*ROR
	Locations     []Location     `json:"locations"`
	Established   int            `json:"established"`
	ExternalIDs   []ExternalID   `json:"external_ids"`
	Links         []Link         `json:"links"`
	Relationships []Relationship `json:"relationships"`
	Status        string         `json:"status"`
	Types         []string       `json:"types"`
}

// InvenioRDM represents the ROR metadata record in InvenioRDM format.
type InvenioRDM struct {
	ID          string       `json:"id"`
	Identifiers []Identifier `json:"identifiers"`
	Name        string       `json:"name"`
}

type ExternalID struct {
	Type      string   `json:"type"`
	All       []string `json:"all"`
	Preferred string   `json:"preferred"`
}

type GeonamesDetails struct {
	ContinentCode          string  `json:"continent_code"`
	ContinentName          string  `json:"continent_name"`
	CountryCode            string  `json:"country_code"`
	CountryName            string  `json:"country_name"`
	CountrySubdivisionCode string  `json:"country_subdivision_code"`
	CountrySubdivisionName string  `json:"country_subdivision_name"`
	Lat                    float64 `json:"lat"`
	Lng                    float64 `json:"lng"`
	Name                   string  `json:"name"`
}

type Identifier struct {
	Identifier string `json:"identifier"`
	Scheme     string `json:"scheme"`
}

type Location struct {
	GeonamesID      int             `json:"geonames_id"`
	GeonamesDetails GeonamesDetails `json:"geonames_details"`
}

type Link struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Name struct {
	Value string   `json:"value"`
	Types []string `json:"types"`
	Lang  string   `json:"lang"`
}

type Relationship struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	ID    string `json:"id"`
}

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

// LoadAll loads the metadata for a list of organizations from a ROR JSON file
func LoadAll(filename string) ([]ROR, error) {
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

// Read reads ROR full metadata and converts it into ROR minimal metadata.
func Read(content Content) (ROR, error) {
	var data ROR

	data.ID = content.ID
	data.Names = content.Names
	data.Admin.LastModified.Date = content.Admin.LastModified.Date

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
	var data []InvenioRDM
	var extracted []InvenioRDM
	var ids []string
	var err error

	// Load the ROR metadata from the embedded ZIP file with all ROR records
	filename := filepath.Join("ror", "affiliations_ror.yaml.zip")
	out, err := fileutils.ReadZIPFile(filename)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(out, &data)
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
							idx := slices.IndexFunc(data, func(d InvenioRDM) bool { return d.ID == id })
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

	output, err := yaml.Marshal(extracted)
	return output, err
}

// Convert converts ROR metadata into InvenioRDM format.
func Convert(data ROR) (InvenioRDM, error) {
	var inveniordm InvenioRDM

	id, _ := utils.ValidateROR(data.ID)
	inveniordm.ID = id
	inveniordm.Identifiers = []Identifier{
		{
			Identifier: id,
			Scheme:     "ror",
		},
	}
	for _, name := range data.Names {
		if slices.Contains(name.Types, "ror_display") {
			inveniordm.Name = name.Value
		}
	}
	return inveniordm, nil
}

// Write writes ROR metadata to InvenioRDM YAML format.
func Write(data ROR) ([]byte, error) {
	inveniordm, err := Convert(data)
	if err != nil {
		fmt.Println(err)
	}
	output, err := yaml.Marshal(inveniordm)
	return output, err
}

// WriteAll writes a list of ROR metadata in InvenioRDM YAML format.
func WriteAll(list []ROR, to string) ([]byte, error) {
	var inveniordmList []InvenioRDM
	var err error
	var output []byte

	if to != "inveniordm" {
		return output, errors.New("unsupported output format")
	}

	for _, data := range list {
		inveniordm, err := Convert(data)
		if err != nil {
			fmt.Println(err)
		}
		if inveniordm.ID != "" {
			inveniordmList = append(inveniordmList, inveniordm)
		}
	}

	output, err = yaml.Marshal(inveniordmList)
	return output, err
}
