// Package doiutils provides a set of functions to work with DOIs
package doiutils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/front-matter/commonmeta/crockford"
)

// PrefixFromUrl extracts DOI prefix from URL
func PrefixFromUrl(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	if u.Host == "" || u.Host != "doi.org" || !strings.HasPrefix(u.Path, "/10.") {
		return "", nil
	}
	path := strings.Split(u.Path, "/")
	return path[1], nil
}

// NormalizeDOI normalizes a DOI
func NormalizeDOI(doi string) string {
	doistr, ok := ValidateDOI(doi)
	if !ok {
		return ""
	}
	resolver := DOIResolver(doi, false)
	return resolver + strings.ToLower(doistr)
}

// ValidateDOI validates a DOI
func ValidateDOI(doi string) (string, bool) {
	r, err := regexp.Compile(`^(?:(http|https):/(/)?(dx\.)?(doi\.org|handle\.stage\.datacite\.org|handle\.test\.datacite\.org)/)?(doi:)?(10\.\d{4,5}/[^\s]+)$`)
	if err != nil {
		log.Printf("Error compiling regex: %v", err)
		return "", false
	}
	matched := r.FindStringSubmatch(doi)
	if len(matched) == 0 {
		return "", false
	}
	return matched[6], true
}

// EscapeDOI escapes a DOI, i.e. replaces '/' with '%2F'
func EscapeDOI(doi string) string {
	doistr, ok := ValidateDOI(doi)
	if !ok {
		return ""
	}
	return strings.ReplaceAll(doistr, "/", "%2F")
}

func EncodeDOI(prefix string) string {
	suffix := crockford.Generate(10, 5, true)
	doi := fmt.Sprintf("https://doi.org/%s/%s", prefix, suffix)
	if IsRegisteredDOI(doi) {
		return EncodeDOI(prefix)
	}
	return doi
}

func DecodeDOI(doi string) int64 {
	d, ok := ValidateDOI(doi)
	if !ok {
		return 0
	}
	suffix := strings.Split(d, "/")[1]
	number, err := crockford.Decode(suffix, true)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return number
}

// IsRegisteredDOI checks if a DOI resolves (i.e. redirects) via the DOI handle servers
func IsRegisteredDOI(doi string) bool {
	url := NormalizeDOI(doi)
	if url == "" {
		return false
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode <= 308
}

// ValidatePrefix validates a DOI prefix for a given DOI
func ValidatePrefix(doi string) (string, bool) {
	r, err := regexp.Compile(`^(?:(http|https):/(/)?(dx\.)?(doi\.org|handle\.stage\.datacite\.org|handle\.test\.datacite\.org)/)?(doi:)?(10\.\d{4,5})`)
	if err != nil {
		log.Printf("Error compiling regex: %v", err)
		return "", false
	}
	matched := r.FindStringSubmatch(doi)
	if len(matched) == 0 {
		return "", false
	}
	return matched[6], true
}

// DOIResolver returns a DOI resolver for a given DOI
func DOIResolver(doi string, sandbox bool) string {
	d, err := url.Parse(doi)
	if err != nil {
		return ""
	}
	if d.Host == "stage.datacite.org" || sandbox {
		return "https://handle.stage.datacite.org/"
	}
	return "https://doi.org/"
}

// GetDOIRA returns the DOI registration agency for a given DOI or prefix
func GetDOIRA(doi string) (string, bool) {
	var knownCrossrefPrefixes = []string{
		"10.53731",
		"10.54900",
		"10.59347",
		"10.59348",
		"10.59349",
		"10.59350",
		"10.59351",
		"10.64000",
	}
	var knownDatacitePrefixes = []string{
		"10.34732",
		"10.57689",
		"10.83132",
	}
	prefix, ok := ValidatePrefix(doi)
	if !ok {
		return "", false
	}

	// check for known prefixes, e.g. regularly used by Rogue Scholar
	if slices.Contains(knownCrossrefPrefixes, prefix) {
		return "Crossref", true
	} else if slices.Contains(knownDatacitePrefixes, prefix) {
		return "DataCite", true
	}

	type Response []struct {
		DOI string `json:"DOI"`
		RA  string `json:"RA"`
	}
	var result Response
	resp, err := http.Get(fmt.Sprintf("https://doi.org/ra/%s", prefix))
	if err != nil {
		return "", false
	}
	if resp.StatusCode == 404 {
		return "", false
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false
	}
	defer resp.Body.Close()
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", false
	}
	return result[0].RA, true
}

// IsRogueScholarDOI checks if a DOI is from Rogue Scholar
func IsRogueScholarDOI(doi string, ra string) bool {
	var rogueScholarCrossrefPrefixes = []string{
		"10.53731",
		"10.54900",
		"10.57689",
		"10.59347",
		"10.59348",
		"10.59349",
		"10.59350",
		"10.63485",
		"10.64000",
	}
	var rogueScholarDatacitePrefixes = []string{
		"10.5438",
		"10.34732",
		"10.57689",
		"10.58079",
		"10.60804",
		// "10.83132",
	}
	prefix, ok := ValidatePrefix(doi)
	if !ok {
		return false
	}
	isCrossref := slices.Contains(rogueScholarCrossrefPrefixes, prefix)
	isDatacite := slices.Contains(rogueScholarDatacitePrefixes, prefix)
	if ra == "crossref" {
		return isCrossref
	} else if ra == "datacite" {
		return isDatacite
	}
	return isCrossref || isDatacite
}
