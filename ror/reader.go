package ror

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	"gopkg.in/yaml.v3"
)

// ROR represents the ROR metadata record.
type ROR struct {
	ID            string        `json:"id" csv:"id"`
	Domains       Strings       `json:"domains,omitempty" yaml:"domains,omitempty"`
	Established   int           `json:"established,omitempty" yaml:"established,omitempty"`
	ExternalIDs   ExternalIDS   `json:"external_ids,omitempty" yaml:"external_ids,omitempty"`
	Links         Links         `json:"links" yaml:"links,omitempty"`
	Locations     Locations     `json:"locations"`
	Names         Names         `json:"names"`
	Relationships Relationships `json:"relationships,omitempty" yaml:"relationships,omitempty"`
	Status        string        `json:"status"`
	Types         Strings       `json:"types"`
	Admin         Admin         `json:"admin"`
}

type Strings []string
type ExternalIDS []ExternalID
type Links []Link
type Locations []Location
type Names []Name
type Relationships []Relationship

type Admin struct {
	Created      Date `json:"created"`
	LastModified Date `json:"last_modified"`
}

type Date struct {
	Date          string `json:"date"`
	SchemaVersion string `json:"schema_version"`
}

type ExternalID struct {
	Type      string  `json:"type"`
	All       Strings `json:"all"`
	Preferred string  `json:"preferred,omitempty" yaml:"preferred,omitempty"`
}

type GeonamesDetails struct {
	ContinentCode          string  `json:"continent_code" yaml:"continent_code"`
	ContinentName          string  `json:"continent_name" yaml:"continent_name"`
	CountryCode            string  `json:"country_code" yaml:"country_code"`
	CountryName            string  `json:"country_name" yaml:"country_name"`
	CountrySubdivisionCode string  `json:"country_subdivision_code,omitempty" yaml:"country_subdivision_code,omitempty"`
	CountrySubdivisionName string  `json:"country_subdivision_name,omitempty" yaml:"country_subdivision_name,omitempty"`
	Lat                    float64 `json:"lat"`
	Lng                    float64 `json:"lng"`
	Name                   string  `json:"name"`
}

type Link struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Location struct {
	GeonamesID      int             `json:"geonames_id" yaml:"geonames_id"`
	GeonamesDetails GeonamesDetails `json:"geonames_details" yaml:"geonames_details"`
}

type Name struct {
	Value string  `json:"value"`
	Types Strings `json:"types"`
	Lang  string  `json:"lang,omitempty" yaml:"lang,omitempty"`
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
	"v1.59": "2025-01-23",
	"v1.60": "2025-02-27",
	"v1.61": "2025-03-18",
	"v1.62": "2025-03-27",
	"v1.63": "2025-04-03",
	"v1.64": "2025-04-28",
	"v1.65": "2025-05-05",
	"v1.66": "2025-05-20",
	"v1.67": "2025-06-24",
}
var DefaultVersion = "v1.67"

// RORZenodoIDs contains the Zenodo IDs for the ROR data releases.
var RORZenodoIDs = map[string]string{
	"v1.59": "14728473",
	"v1.60": "14797924",
	"v1.61": "15047759",
	"v1.62": "15098078",
	"v1.63": "15132361",
	"v1.64": "15298417",
	"v1.65": "15343380",
	"v1.66": "15475023",
	"v1.67": "15731450",
}

var ArchivedFilename = fmt.Sprintf("%s-%s-ror-data_schema_v2.json", DefaultVersion, RORVersions[DefaultVersion])
var RORFilename = fmt.Sprintf("%s-%s-ror-data.zip", DefaultVersion, RORVersions[DefaultVersion])
var RORDownloadURL = fmt.Sprintf("https://github.com/ror-community/ror-data/blob/main/%s", RORFilename)

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
	if !slices.Contains(commonmeta.OrganizationTypes, type_) {
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
	if !slices.Contains(commonmeta.OrganizationTypes, type_) {
		return ror, errors.New("not a supported organization id")
	}

	list, err := LoadBuiltin()
	if err != nil {
		return ror, err
	}
	if type_ == "ISNI" {
		// ROR expects ISNI IDs to be in the form of 0000 0002 1234 5678
		pid = utils.SplitString(pid, 4, " ")
		pid = strings.ReplaceAll(pid, "-", " ")
	}
	if type_ == "ROR" {
		idx = slices.IndexFunc(list, func(d ROR) bool { return d.ID == utils.NormalizeROR(pid) })
	} else {
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
	}
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
func FetchAll(version string) ([]ROR, error) {
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
	input, err = fileutils.DownloadFile(url, true)
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
	return list, nil
}

// LoadAll loads the metadata for a list of organizations from a ROR file.
func LoadAll(filename string) ([]ROR, error) {
	var list []ROR
	var output []byte
	var err error

	filename, extension, compress := fileutils.GetExtension(filename, ".json")

	switch compress {
	case "gz":
		output, err = fileutils.ReadGZFile(filename + ".gz")
		if err != nil {
			return list, errors.New("error reading gz file")
		}
	case "zip":
		output, err = fileutils.ReadZIPFile(filename+".zip", path.Base(filename))
		if err != nil {
			return list, errors.New("error reading zip file")
		}
	default:
		output, err = fileutils.ReadFile(filename)
		if err != nil {
			return list, errors.New("error reading file")
		}
	}

	switch extension {
	case ".json":
		err = json.Unmarshal(output, &list)
		if err != nil {
			return list, errors.New("error unmarshalling json file")
		}
	case ".yaml":
		err = yaml.Unmarshal(output, &list)
		if err != nil {
			return list, errors.New("error unmarshalling yaml file")
		}
	case ".jsonl":
		decoder := json.NewDecoder(bytes.NewReader(output))
		for {
			var item ROR
			if err := decoder.Decode(&item); err == io.EOF {
				break
			} else if err != nil {
				return list, errors.New("error unmarshalling jsonl file")
			}
			list = append(list, item)
		}
	case ".sql":

	default:
		return list, errors.New("unsupported file format")
	}
	return list, err
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

// LoadBuiltin loads the ROR metadata from the embedded ROR catalog in zipped json format.
func LoadBuiltin() ([]ROR, error) {
	var list []ROR

	bytes, err := vocabularies.LoadVocabulary("ROR.Organizations")
	err = json.Unmarshal(bytes, &list)
	if err != nil {
		return nil, fmt.Errorf("error parsing json data: %w", err)
	}

	return list, nil
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
		bytes, err = fileutils.DownloadFile(RORDownloadURL, true)
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
	var extracted []ROR
	var ids []string

	// Load the ROR metadata from the embedded ZIP file with all ROR records
	list, err := LoadBuiltin()
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
							idx := slices.IndexFunc(list, func(d ROR) bool { return d.ID == id })
							if idx != -1 {
								ids = append(ids, a.ID)
								extracted = append(extracted, list[idx])
							}
						}
					}
				}
			}
		}
	}

	output, err := json.Marshal(extracted)
	return output, err
}
