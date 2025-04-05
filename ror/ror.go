// Package ror converts ROR (Research Organization Registry) metadata.
package ror

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"slices"

	"gopkg.in/yaml.v3"

	"github.com/front-matter/commonmeta/utils"
)

// ROR represents the minimal ROR metadata record.
type ROR struct {
	ID    string `json:"id"`
	Names []Name `json:"names"`
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
	Admin         struct {
		Created struct {
			Date          string `json:"date"`
			SchemaVersion string `json:"schema_version"`
		}
		LastModified struct {
			Date          string `json:"date"`
			SchemaVersion string `json:"schema_version"`
		}
	}
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

// LoadAll loads the metadata for a list of organizations from a ROR JSON file
func LoadAll(filename string) ([]ROR, error) {
	var data []ROR
	var content []Content
	var err error

	extension := path.Ext(filename)
	if extension == ".json" {
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
		err = decoder.Decode(&content)
		if err != nil {
			return data, err
		}
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

// WriteAll writes a list of ROR metadata to InvenioRDM YAML format.
func WriteAll(list []ROR, input string) (string, error) {
	var inveniordmList []InvenioRDM

	for _, data := range list {
		inveniordm, err := Convert(data)
		if err != nil {
			fmt.Println(err)
		}
		inveniordmList = append(inveniordmList, inveniordm)
	}
	output, err := yaml.Marshal(inveniordmList)
	if err != nil {
		fmt.Println(err)
	}

	filename := path.Base("affiliations_ror.yaml")
	output = append([]byte("# file generated from "+input+"\n\n"), output...)
	err = os.WriteFile(filename, output, 0644)
	return filename, err
}
