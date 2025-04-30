package spdx

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"slices"
	"strings"

	"github.com/front-matter/commonmeta/fileutils"
	"github.com/front-matter/commonmeta/vocabularies"
)

type SPDX struct {
	LicenseListVersion string    `json:"licenseListVersion"`
	LicenseList        []License `json:"licenses"`
	ReleaseDate        string    `json:"releaseDate"`
}

type License struct {
	Reference             string   `json:"reference"`
	IsDeprecatedLicenseID bool     `json:"isDeprecatedLicenseId"`
	DetailsURL            string   `json:"detailsUrl"`
	ReferenceNumber       int      `json:"referenceNumber"`
	Name                  string   `json:"name"`
	LicenseID             string   `json:"licenseId"`
	SeeAlso               []string `json:"seeAlso"`
	IsOsiApproved         bool     `json:"isOsiApproved"`
	IsFsfLibre            bool     `json:"isFsfLibre"`
}

// SPDXDownloadURL is the URL to download the SPDX licenses JSON file.
const SPDXDownloadURL = "https://raw.githubusercontent.com/spdx/license-list-data/refs/heads/main/json/licenses.json"

// SPDXFilename is the name of the SPDX licenses JSON file.
const SPDXFilename = "licenses.json"

// FetchAll fetches the SPDX licenses from GitHub.
func FetchAll() ([]License, error) {
	var output []byte
	var content SPDX
	var list []License
	var err error

	// download the SPDX license file
	output, err = fileutils.DownloadFile(SPDXDownloadURL, false)
	if err != nil {
		fmt.Println(err)
		return list, fmt.Errorf("error downloading SPDX file")
	}

	// write the json file to disk
	err = fileutils.WriteFile(SPDXFilename, output)
	if err != nil {
		return list, fmt.Errorf("error writing json file: %w", err)
	}

	err = json.Unmarshal(output, &content)
	if err != nil {
		return list, fmt.Errorf("error parsing json data: %w", err)
	}
	list = content.LicenseList
	return list, nil
}

// LoadBuiltin loads the SPDX metadata from the embedded SPDX catalog.
func LoadBuiltin() ([]License, error) {
	var content SPDX
	var list []License

	bytes, err := vocabularies.LoadVocabulary("SPDX.Licenses")
	if err != nil {
		return nil, fmt.Errorf("error loading spdx file: %w", err)
	}
	err = json.Unmarshal(bytes, &content)
	if err != nil {
		return nil, fmt.Errorf("error parsing spdx data: %w", err)
	}
	list = content.LicenseList
	return list, nil
}

// LoadJSON loads a local licenses json file.
func LoadJSON(version string) ([]SPDX, error) {
	var bytes []byte
	var list []SPDX

	// check that local SPDX licenses file exists.
	info, err := os.Stat(SPDXFilename)
	if err == nil && (info.Size()/1024) > 50 {
		// read json file
		bytes, err = fileutils.ReadFile(SPDXFilename)
		if err != nil {
			return nil, fmt.Errorf("error reading json file: %w", err)
		}
	} else {
		bytes, err = fileutils.DownloadFile(SPDXDownloadURL, false)
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

// Search searches local SPDX metadata for a given spdx licenseId or url.
func Search(id string) (License, error) {
	var url_ string
	var idx int
	var license License

	// check if the id is a URL
	u, err := url.Parse(id)
	if err == nil && u.Scheme != "" && u.Host != "" {
		url_ = u.String()
	}

	list, err := LoadBuiltin()
	if err != nil {
		return license, err
	}
	if url_ != "" {
		// check if the id is a URL found in the seeAlso field
		idx = slices.IndexFunc(list, func(d License) bool {
			return slices.Contains(d.SeeAlso, url_)
		})
	} else {
		// check if the id is a licenseId
		idx = slices.IndexFunc(list, func(l License) bool {
			return strings.EqualFold(l.LicenseID, id)
		})
	}
	if idx == -1 {
		return license, nil
	}
	license = list[idx]
	return license, nil
}
