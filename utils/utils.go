// Package utils provides utility functions for commonmeta.
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/front-matter/commonmeta/doiutils"

	"github.com/microcosm-cc/bluemonday"
)

// ROR represents a Research Organization Registry (ROR) record
type ROR struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Types       []string `json:"types"`
	ExternalIds struct {
		FundRef struct {
			Preferred string   `json:"preferred"`
			All       []string `json:"all"`
		} `json:"FundRef"`
	} `json:"external_ids"`
}

// NormalizeID checks for valid DOI or HTTP(S) URL and normalizes them
func NormalizeID(pid string) string {
	// check for valid DOI
	doi := doiutils.NormalizeDOI(pid)
	if doi != "" {
		return doi
	}

	// check for valid UUID
	uuid, ok := ValidateUUID(pid)
	if ok {
		return uuid
	}

	// check for valid URL
	uri, err := url.Parse(pid)
	if err != nil {
		return ""
	}
	if uri.Scheme == "" {
		return ""
	}

	// check for valid HTTP uri and ensure https
	if uri.Scheme == "http" {
		uri.Scheme = "https"
	}

	// remove trailing slash
	if pid[len(pid)-1] == '/' {
		pid = pid[:len(pid)-1]
	}

	return pid
}

// NormalizeURL normalizes URL
func NormalizeURL(str string, secure bool, lower bool) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	if u.Path != "" && u.Path[len(u.Path)-1] == '/' {
		u.Path = u.Path[:len(u.Path)-1]
	}
	if secure && u.Scheme == "http" {
		u.Scheme = "https"
	}
	if lower {
		return strings.ToLower(u.String()), nil
	}
	return u.String(), nil
}

// NormalizeCCUrl returns the normalized Creative Commons License URL
func NormalizeCCUrl(url string) (string, bool) {
	NormalizedLicenses := map[string]string{
		"https://creativecommons.org/licenses/by/1.0":          "https://creativecommons.org/licenses/by/1.0/legalcode",
		"https://creativecommons.org/licenses/by/2.0":          "https://creativecommons.org/licenses/by/2.0/legalcode",
		"https://creativecommons.org/licenses/by/2.5":          "https://creativecommons.org/licenses/by/2.5/legalcode",
		"https://creativecommons.org/licenses/by/3.0":          "https://creativecommons.org/licenses/by/3.0/legalcode",
		"https://creativecommons.org/licenses/by/3.0/us":       "https://creativecommons.org/licenses/by/3.0/legalcode",
		"https://creativecommons.org/licenses/by/4.0":          "https://creativecommons.org/licenses/by/4.0/legalcode",
		"https://creativecommons.org/licenses/by-nc/1.0":       "https://creativecommons.org/licenses/by-nc/1.0/legalcode",
		"https://creativecommons.org/licenses/by-nc/2.0":       "https://creativecommons.org/licenses/by-nc/2.0/legalcode",
		"https://creativecommons.org/licenses/by-nc/2.5":       "https://creativecommons.org/licenses/by-nc/2.5/legalcode",
		"https://creativecommons.org/licenses/by-nc/3.0":       "https://creativecommons.org/licenses/by-nc/3.0/legalcode",
		"https://creativecommons.org/licenses/by-nc/4.0":       "https://creativecommons.org/licenses/by-nc/4.0/legalcode",
		"https://creativecommons.org/licenses/by-nd-nc/1.0":    "https://creativecommons.org/licenses/by-nd-nc/1.0/legalcode",
		"https://creativecommons.org/licenses/by-nd-nc/2.0":    "https://creativecommons.org/licenses/by-nd-nc/2.0/legalcode",
		"https://creativecommons.org/licenses/by-nd-nc/2.5":    "https://creativecommons.org/licenses/by-nd-nc/2.5/legalcode",
		"https://creativecommons.org/licenses/by-nd-nc/3.0":    "https://creativecommons.org/licenses/by-nd-nc/3.0/legalcode",
		"https://creativecommons.org/licenses/by-nd-nc/4.0":    "https://creativecommons.org/licenses/by-nd-nc/4.0/legalcode",
		"https://creativecommons.org/licenses/by-nc-sa/1.0":    "https://creativecommons.org/licenses/by-nc-sa/1.0/legalcode",
		"https://creativecommons.org/licenses/by-nc-sa/2.0":    "https://creativecommons.org/licenses/by-nc-sa/2.0/legalcode",
		"https://creativecommons.org/licenses/by-nc-sa/2.5":    "https://creativecommons.org/licenses/by-nc-sa/2.5/legalcode",
		"https://creativecommons.org/licenses/by-nc-sa/3.0":    "https://creativecommons.org/licenses/by-nc-sa/3.0/legalcode",
		"https://creativecommons.org/licenses/by-nc-sa/3.0/us": "https://creativecommons.org/licenses/by-nc-sa/3.0/legalcode",
		"https://creativecommons.org/licenses/by-nc-sa/4.0":    "https://creativecommons.org/licenses/by-nc-sa/4.0/legalcode",
		"https://creativecommons.org/licenses/by-nd/1.0":       "https://creativecommons.org/licenses/by-nd/1.0/legalcode",
		"https://creativecommons.org/licenses/by-nd/2.0":       "https://creativecommons.org/licenses/by-nd/2.0/legalcode",
		"https://creativecommons.org/licenses/by-nd/2.5":       "https://creativecommons.org/licenses/by-nd/2.5/legalcode",
		"https://creativecommons.org/licenses/by-nd/3.0":       "https://creativecommons.org/licenses/by-nd/3.0/legalcode",
		"https://creativecommons.org/licenses/by-nd/4.0":       "https://creativecommons.org/licenses/by-nd/2.0/legalcode",
		"https://creativecommons.org/licenses/by-sa/1.0":       "https://creativecommons.org/licenses/by-sa/1.0/legalcode",
		"https://creativecommons.org/licenses/by-sa/2.0":       "https://creativecommons.org/licenses/by-sa/2.0/legalcode",
		"https://creativecommons.org/licenses/by-sa/2.5":       "https://creativecommons.org/licenses/by-sa/2.5/legalcode",
		"https://creativecommons.org/licenses/by-sa/3.0":       "https://creativecommons.org/licenses/by-sa/3.0/legalcode",
		"https://creativecommons.org/licenses/by-sa/4.0":       "https://creativecommons.org/licenses/by-sa/4.0/legalcode",
		"https://creativecommons.org/licenses/by-nc-nd/1.0":    "https://creativecommons.org/licenses/by-nc-nd/1.0/legalcode",
		"https://creativecommons.org/licenses/by-nc-nd/2.0":    "https://creativecommons.org/licenses/by-nc-nd/2.0/legalcode",
		"https://creativecommons.org/licenses/by-nc-nd/2.5":    "https://creativecommons.org/licenses/by-nc-nd/2.5/legalcode",
		"https://creativecommons.org/licenses/by-nc-nd/3.0":    "https://creativecommons.org/licenses/by-nc-nd/3.0/legalcode",
		"https://creativecommons.org/licenses/by-nc-nd/4.0":    "https://creativecommons.org/licenses/by-nc-nd/4.0/legalcode",
		"https://creativecommons.org/licenses/publicdomain":    "https://creativecommons.org/licenses/publicdomain/",
		"https://creativecommons.org/publicdomain/zero/1.0":    "https://creativecommons.org/publicdomain/zero/1.0/legalcode",
	}

	if url == "" {
		return "", false
	}
	var err error
	url, err = NormalizeURL(url, true, false)
	if err != nil {
		return "", false
	}
	normalizedURL, ok := NormalizedLicenses[url]
	if !ok {
		return url, false
	}
	return normalizedURL, true
}

// URLToSPDX provides the SPDX license ID given a Creative Commons URL
func URLToSPDX(url string) string {
	// appreviated list from https://spdx.org/licenses/
	SPDXLicenses := map[string]string{
		"https://creativecommons.org/licenses/by/3.0/legalcode":       "CC-BY-3.0",
		"https://creativecommons.org/licenses/by/4.0/legalcode":       "CC-BY-4.0",
		"https://creativecommons.org/licenses/by-nc/3.0/legalcode":    "CC-BY-NC-3.0",
		"https://creativecommons.org/licenses/by-nc/4.0/legalcode":    "CC-BY-NC-4.0",
		"https://creativecommons.org/licenses/by-nc-nd/3.0/legalcode": "CC-BY-NC-ND-3.0",
		"https://creativecommons.org/licenses/by-nc-nd/4.0/legalcode": "CC-BY-NC-ND-4.0",
		"https://creativecommons.org/licenses/by-nc-sa/3.0/legalcode": "CC-BY-NC-SA-3.0",
		"https://creativecommons.org/licenses/by-nc-sa/4.0/legalcode": "CC-BY-NC-SA-4.0",
		"https://creativecommons.org/licenses/by-nd/3.0/legalcode":    "CC-BY-ND-3.0",
		"https://creativecommons.org/licenses/by-nd/4.0/legalcode":    "CC-BY-ND-4.0",
		"https://creativecommons.org/licenses/by-sa/3.0/legalcode":    "CC-BY-SA-3.0",
		"https://creativecommons.org/licenses/by-sa/4.0/legalcode":    "CC-BY-SA-4.0",
		"https://creativecommons.org/publicdomain/zero/1.0/legalcode": "CC0-1.0",
		"https://creativecommons.org/licenses/publicdomain/":          "CC0-1.0",
		"https://opensource.org/licenses/MIT":                         "MIT",
		"https://opensource.org/licenses/Apache-2.0":                  "Apache-2.0",
		"https://opensource.org/licenses/GPL-3.0":                     "GPL-3.0",
	}
	id := SPDXLicenses[url]
	return id
}

type params struct {
	Pid, Str, Ext, Filename string
	Dct                     map[string]interface{}
}

// FindFromFormat finds the commonmeta read format
func FindFromFormat(p params) string {
	// Find reader from format
	if p.Pid != "" {
		return FindFromFormatByID(p.Pid)
	}
	if p.Str != "" && p.Ext != "" {
		return FindFromFormatByExt(p.Ext)
	}
	if p.Dct != nil {
		return FindFromFormatByDict(p.Dct)
	}
	if p.Str != "" {
		return FindFromFormatByString(p.Str)
	}
	if p.Filename != "" {
		return FindFromFormatByFilename(p.Filename)
	}
	return "datacite"
}

// FindFromFormatByID finds the commonmeta reader from format by id
func FindFromFormatByID(id string) string {
	_, ok := doiutils.ValidateDOI(id)
	if ok {
		return "datacite"
	}
	if strings.Contains(id, "github.com") {
		return "cff"
	}
	if strings.Contains(id, "codemeta.json") {
		return "codemeta"
	}
	if strings.Contains(id, "json_feed_item") {
		return "json_feed_item"
	}
	if strings.Contains(id, "zenodo.org") {
		return "inveniordm"
	}
	return "schema_org"
}

// FindFromFormatByExt finds the commonmeta reader from format by file extension
func FindFromFormatByExt(ext string) string {
	if ext == ".bib" {
		return "bibtex"
	}
	if ext == ".ris" {
		return "ris"
	}
	return ""
}

// FindFromFormatByDict finds the commonmeta reader from format by dictionary
func FindFromFormatByDict(dct map[string]interface{}) string {
	if dct == nil {
		return ""
	}
	if v, ok := dct["schema_version"]; ok && strings.HasPrefix(v.(string), "https://commonmeta.org") {
		return "commonmeta"
	}
	if v, ok := dct["@context"]; ok && v == "http://schema.org" {
		return "schema_org"
	}
	if v, ok := dct["@context"]; ok && v == "https://raw.githubusercontent.com/codemeta/codemeta/master/codemeta.jsonld" {
		return "codemeta"
	}
	if _, ok := dct["guid"]; ok {
		return "json_feed_item"
	}
	if v, ok := dct["schemaVersion"]; ok && strings.HasPrefix(v.(string), "http://datacite.org/schema/kernel") {
		return "datacite"
	}
	if v, ok := dct["source"]; ok && v == "Crossref" {
		return "crossref"
	}
	if _, ok := dct["issued.date-parts"]; ok {
		return "csl"
	}
	if _, ok := dct["conceptdoi"]; ok {
		return "inveniordm"
	}
	if _, ok := dct["credit_metadata"]; ok {
		return "kbase"
	}
	return ""
}

// FindFromFormatByString finds the commonmeta reader from format by string
func FindFromFormatByString(str string) string {
	if str == "" {
		return ""
	}
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(str), &data); err != nil {
		return ""
	}
	if v, ok := data["schema_version"]; ok && strings.HasPrefix(v.(string), "https://commonmeta.org") {
		return "commonmeta"
	}
	if v, ok := data["@context"]; ok && v == "http://schema.org" {
		return "schema_org"
	}
	if v, ok := data["@context"]; ok && v == "https://raw.githubusercontent.com/codemeta/codemeta/master/codemeta.jsonld" {
		return "codemeta"
	}
	if _, ok := data["guid"]; ok {
		return "json_feed_item"
	}
	if v, ok := data["schemaVersion"]; ok && strings.HasPrefix(v.(string), "http://datacite.org/schema/kernel") {
		return "datacite"
	}
	if v, ok := data["source"]; ok && v == "Crossref" {
		return "crossref"
	}
	if _, ok := data["issued.date-parts"]; ok {
		return "csl"
	}
	if _, ok := data["conceptdoi"]; ok {
		return "inveniordm"
	}
	if _, ok := data["credit_metadata"]; ok {
		return "kbase"
	}
	return ""
}

// FindFromFormatByFilename finds the commonmeta reader from format by filename
func FindFromFormatByFilename(filename string) string {
	if filename == "CITATION.cff" {
		return "cff"
	}
	return ""
}

// ISSNAsURL returns the ISSN expressed as URL
func ISSNAsURL(issn string) string {
	if issn == "" {
		return ""
	}
	return fmt.Sprintf("https://portal.issn.org/resource/ISSN/%s", issn)
}

// Sanitize removes all HTML tags except for a whitelist of allowed tags. Used for
// title and description fields.
func Sanitize(html string) string {
	policy := bluemonday.StrictPolicy()
	policy.AllowElements("b", "br", "code", "em", "i", "sub", "sup", "strong")
	sanitizedHTML := policy.Sanitize(html)
	str := strings.Trim(sanitizedHTML, "\n")
	return str
}

// TitleCase capitalizes the first letter of a string without changing the rest
func TitleCase(str string) string {
	return strings.ToUpper(string(str[0])) + str[1:]
}

// UnescapeUTF8 unescapes UTF-8 characters
func UnescapeUTF8(inStr string) (outStr string, err error) {
	jsonStr := `"` + strings.ReplaceAll(inStr, `"`, `\"`) + `"`
	err = json.Unmarshal([]byte(jsonStr), &outStr)
	return
}

// DedupeSlice removes duplicates from a slice
// https://stackoverflow.com/questions/66643946/how-to-remove-duplicates-strings-or-int-from-slice-in-go/76948712#76948712
func DedupeSlice[T comparable](sliceList []T) []T {
	dedupeMap := make(map[T]struct{})
	list := []T{}

	for _, slice := range sliceList {
		if _, exists := dedupeMap[slice]; !exists {
			dedupeMap[slice] = struct{}{}
			list = append(list, slice)
		}
	}

	return list
}

// NormalizeORCID returns a normalized ORCID URL
func NormalizeORCID(orcid string) string {
	orcidStr, ok := ValidateORCID(orcid)
	if !ok {
		return ""
	}
	return "https://orcid.org/" + orcidStr
}

// ValidateORCID validates an ORCID
func ValidateORCID(orcid string) (string, bool) {
	r, err := regexp.Compile(`^(?:(?:http|https)://(?:(?:www|sandbox)?\.)?orcid\.org/)?(\d{4}[ -]\d{4}[ -]\d{4}[ -]\d{3}[0-9X]+)$`)
	if err != nil {
		log.Printf("Error compiling regex: %v", err)
		return "", false
	}
	matched := r.FindStringSubmatch(orcid)
	if len(matched) == 0 {
		return "", false
	}
	return matched[1], true
}

// NormalizeROR returns a normalized ROR URL
func NormalizeROR(ror string) string {
	rorStr, ok := ValidateROR(ror)
	if !ok {
		return ""
	}
	return "https://ror.org/" + rorStr
}

// ValidateROR validates a ROR
func ValidateROR(ror string) (string, bool) {
	r, err := regexp.Compile(`^(?:(?:http|https)://ror\.org/)?([0-9a-z]{7}\d{2})$`)
	if err != nil {
		log.Printf("Error compiling regex: %v", err)
		return "", false
	}
	matched := r.FindStringSubmatch(ror)
	if len(matched) == 0 {
		return "", false
	}
	return matched[1], true
}

// GetROR
func GetROR(ror string) (ROR, error) {
	var content ROR
	client := http.Client{
		Timeout: time.Second * 10,
	}
	url := "https://api.ror.org/organizations/" + ror
	resp, err := client.Get(url)
	if err != nil {
		return content, err
	}
	if resp.StatusCode != 200 {
		return content, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return content, err
	}
	err = json.Unmarshal(body, &content)
	if err != nil {
		fmt.Println("error:", err)
	}
	return content, err
}

// ValidateURL validates a URL and checks if it is a DOI
func ValidateURL(str string) string {
	_, ok := doiutils.ValidateDOI(str)
	if ok {
		return "DOI"
	}
	u, err := url.Parse(str)
	if err != nil {
		return ""
	}
	if u.Scheme == "http" || u.Scheme == "https" {
		return "URL"
	}
	return ""
}

// ValidateUUID validates a UUID
func ValidateUUID(uuid string) (string, bool) {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	if !r.MatchString(uuid) {
		return "", false
	}
	return r.FindString(uuid), true
}

func CamelCaseToWords(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	words := matchFirstCap.ReplaceAllString(str, "${1} ${2}")
	words = matchAllCap.ReplaceAllString(words, "${1} ${2}")
	return strings.ToUpper(words[:1]) + strings.ToLower(words[1:])
}
