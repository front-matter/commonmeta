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
	"strconv"
	"strings"
	"time"

	iso639_3 "github.com/barbashov/iso639-3"
	"github.com/front-matter/commonmeta/crockford"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkosilo/iso7064"
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
		return "", nil
	}
	if u.Host == "" {
		return "", nil
	}
	// strip trailing slash if no query
	if u.Path != "" && len(u.RawQuery) == 0 && u.Path[len(u.Path)-1] == '/' {
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

// SPDXToURL provides the SPDX license URL given a Creative Commons SPDX ID
func SPDXToURL(id string) string {
	// appreviated list from https://spdx.org/licenses/
	SPDXLicenses := map[string]string{
		"CC-BY-3.0":       "https://creativecommons.org/licenses/by/3.0/legalcode",
		"CC-BY-4.0":       "https://creativecommons.org/licenses/by/4.0/legalcode",
		"CC-BY-NC-3.0":    "https://creativecommons.org/licenses/by-nc/3.0/legalcode",
		"CC-BY-NC-4.0":    "https://creativecommons.org/licenses/by-nc/4.0/legalcode",
		"CC-BY-NC-ND-3.0": "https://creativecommons.org/licenses/by-nc-nd/3.0/legalcode",
		"CC-BY-NC-ND-4.0": "https://creativecommons.org/licenses/by-nc-nd/4.0/legalcode",
		"CC-BY-NC-SA-3.0": "https://creativecommons.org/licenses/by-nc-sa/3.0/legalcode",
		"CC-BY-NC-SA-4.0": "https://creativecommons.org/licenses/by-nc-sa/4.0/legalcode",
		"CC-BY-ND-3.0":    "https://creativecommons.org/licenses/by-nd/3.0/legalcode",
		"CC-BY-ND-4.0":    "https://creativecommons.org/licenses/by-nd/4.0/legalcode",
		"CC-BY-SA-3.0":    "https://creativecommons.org/licenses/by-sa/3.0/legalcode",
		"CC-BY-SA-4.0":    "https://creativecommons.org/licenses/by-sa/4.0/legalcode",
		"CC0-1.0":         "https://creativecommons.org/publicdomain/zero/1.0/legalcode",
		"MIT":             "https://opensource.org/licenses/MIT",
		"Apache-2.0":      "https://opensource.org/licenses/Apache-2.0",
		"GPL-3.0":         "https://opensource.org/licenses/GPL-3.0",
	}
	url := SPDXLicenses[id]
	return url
}

type Params struct {
	Pid, Str, Ext, Filename string
	Map                     map[string]interface{}
}

// FindFromFormat finds the commonmeta read format
func FindFromFormat(p Params) string {
	// Find reader from format
	if p.Pid != "" {
		return FindFromFormatByID(p.Pid)
	}
	if p.Str != "" && p.Ext != "" {
		return FindFromFormatByExt(p.Ext)
	}
	if p.Map != nil {
		return FindFromFormatByMap(p.Map)
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
	doi, ok := doiutils.ValidateDOI(id)
	if ok {
		registrationAgency, ok := doiutils.GetDOIRA(doi)
		if ok {
			return strings.ToLower(registrationAgency)
		}
		return "datacite"
	}
	if strings.HasSuffix(id, "codemeta.json") {
		return "codemeta"
	}
	if strings.HasSuffix(id, "CITATION.cff") {
		return "cff"
	}
	if strings.Contains(id, "github.com") {
		return "cff"
	}
	if strings.Contains(id, "jsonfeed") {
		return "jsonfeed"
	}
	r := regexp.MustCompile(`^https:/(/)?api\.rogue-scholar\.org/posts/(.+)$`)
	if len(r.FindStringSubmatch(id)) > 0 {
		return "jsonfeed"
	}
	r = regexp.MustCompile(`^https:/(/)(.+)/(api/)?records/(.+)$`)
	if len(r.FindStringSubmatch(id)) > 0 {
		return "inveniordm"
	}
	return "schemaorg"
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

// FindFromFormatByMap finds the commonmeta reader from format by map
func FindFromFormatByMap(m map[string]interface{}) string {
	if m == nil {
		return ""
	}
	if v, ok := m["schema_version"]; ok && strings.HasPrefix(v.(string), "https://commonmeta.org") {
		return "commonmeta"
	}
	if v, ok := m["@context"]; ok && v == "http://schema.org" {
		return "schemaorg"
	}
	if v, ok := m["@context"]; ok && v == "https://raw.githubusercontent.com/codemeta/codemeta/master/codemeta.jsonld" {
		return "codemeta"
	}
	if _, ok := m["guid"]; ok {
		return "jsonfeed"
	}
	if v, ok := m["schemaVersion"]; ok && strings.HasPrefix(v.(string), "http://datacite.org/schema/kernel") {
		return "datacite"
	}
	if v, ok := m["source"]; ok && v == "Crossref" {
		return "crossref"
	}
	if _, ok := m["issued.date-parts"]; ok {
		return "csl"
	}
	if _, ok := m["conceptdoi"]; ok {
		return "inveniordm"
	}
	if _, ok := m["credit_metadata"]; ok {
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
		return "schemaorg"
	}
	if v, ok := data["@context"]; ok && v == "https://raw.githubusercontent.com/codemeta/codemeta/master/codemeta.jsonld" {
		return "codemeta"
	}
	if _, ok := data["guid"]; ok {
		return "jsonfeed"
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

// ValidateISSN validates an ISSN
func ValidateISSN(issn string) (string, bool) {
	r, err := regexp.Compile(`^(?:https://portal\.issn\.org/resource/ISSN/)?(\d{4}\-\d{3}(\d|x|X))$`)
	if err != nil {
		log.Printf("Error compiling regex: %v", err)
		return "", false
	}
	matched := r.FindStringSubmatch(issn)
	if len(matched) == 0 {
		return "", false
	}
	return matched[1], true
}

// CommunitySlugAsURL returns the InvenioRDM community slug expressed as globally unique URL
func CommunitySlugAsURL(slug string, host string) string {
	if host == "" {
		host = "rogue-scholar.org"
	}
	if slug == "" {
		return ""
	}
	return fmt.Sprintf("https://%s/api/communities/%s", host, slug)
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
// ORCID is a 16-character string in blocks of four
// separated by hyphens between
// 0000-0001-5000-0007 and 0000-0003-5000-0001,
// or between 0009-0000-0000-0000 and 0009-0010-0000-0000.
func ValidateORCID(orcid string) (string, bool) {
	r, err := regexp.Compile(`^(?:(?:http|https)://(?:(?:www|sandbox)?\.)?orcid\.org/)?(000[09][ -]00\d{2}[ -]\d{4}[ -]\d{3}[0-9X]+)$`)
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

// ValidateROR validates a ROR ID. The ROR ID starts with 0 followed by a 6-character
// alphanumeric string which is base32-encoded and a 2-digit checksum.
func ValidateROR(ror string) (string, bool) {
	r, err := regexp.Compile(`^(?:(?:http|https)://ror\.org/)?(0[0-9a-z]{6}\d{2})$`)
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
	client := &http.Client{
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

// ValidateID validates an identifier and returns the type
// Can be DOI, UUID, ISSN, ORCID, ROR, URL
func ValidateID(id string) (string, string) {
	doi, ok := doiutils.ValidateDOI(id)
	if ok {
		prefix, _ := doiutils.ValidatePrefix(doi)
		if prefix == "10.13039" {
			return doi, "Crossref Funder ID"
		}
		return doi, "DOI"
	}
	uuid, ok := ValidateUUID(id)
	if ok {
		return uuid, "UUID"
	}
	orcid, ok := ValidateORCID(id)
	if ok {
		return orcid, "ORCID"
	}
	ror, ok := ValidateROR(id)
	if ok {
		return ror, "ROR"
	}
	issn, ok := ValidateISSN(id)
	if ok {
		return issn, "ISSN"
	}
	url := ValidateURL(id)
	if url != "" {
		return id, "URL"
	}
	return "", ""
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
	// don't allow URLs with certain fragments, e.g. from Software Heritage
	// TODO: testing with more URLs
	disallowedFragments := []string{";origin=", ";jsessionid="}
	for _, f := range disallowedFragments {
		if strings.Contains(u.String(), f) {
			return ""
		}
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

// CamelCaseString converts a pascal case string to camel case
func WordsToCamelCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	words := matchFirstCap.ReplaceAllString(str, "${1} ${2}")
	words = matchAllCap.ReplaceAllString(words, "${1} ${2}")
	var s string
	for _, word := range strings.Fields(words) {
		s += strings.ToUpper(string(word)[:1]) + string(word)[1:]
	}
	s = strings.ReplaceAll(s, " ", "")
	return strings.ToLower(s[:1]) + s[1:]
}

// CamelCaseString converts a pascal case string to camel case
func CamelCaseString(str string) string {
	return strings.ToLower(str[:1]) + str[1:]
}

// KebabCaseToCamelCase converts a kebab case string to camel case
func KebabCaseToCamelCase(str string) string {
	var matchCap = regexp.MustCompile("-([a-z])")
	return matchCap.ReplaceAllStringFunc(str, func(s string) string {
		return strings.ToUpper(s[1:])
	})
}

// KebabCaseToPascalCase converts a kebab case string to pascal case
func KebabCaseToPascalCase(str string) string {
	s := KebabCaseToCamelCase(str)
	return strings.ToUpper(s[:1]) + s[1:]
}

func GetLanguage(lang string, format string) string {
	language := iso639_3.FromAnyCode(lang)
	if language == nil {
		return ""
	} else if format == "iso639-3" {
		return language.Part3
	} else if format == "name" {
		return language.Name
	} else {
		return language.Part1
	}
}

func DecodeID(id string) (int64, error) {
	var number int64
	var ok bool
	var err error

	identifier, identifierType := ValidateID(id)
	if identifierType == "DOI" {
		// the format of a DOI is a prefix and a suffix separated by a slash
		// the prefix starts with 10. and is followed by 4-5 digits
		// the suffix is a string of characters and is not case-sensitive
		// suffixes from Rogue Scholar are base32-encoded numbers with checksums
		suffix := strings.Split(identifier, "/")[1]
		number, err = crockford.Decode(suffix, true)
	} else if identifierType == "ROR" {
		// ROR ID is a 9-character string that starts with 0
		// and is a base32-encoded number with a mod 97-1
		number, err = crockford.Decode(identifier, true)
	} else if identifierType == "ORCID" {
		str := identifier
		identifier = strings.ReplaceAll(identifier, "-", "")
		calc := iso7064.NewMod112Calculator()
		ok, _ = calc.Verify(identifier)
		if !ok {
			cs := identifier[len(identifier)-1:]
			return 0, fmt.Errorf("wrong checksum %s for identifier %s", cs, str)
		}
		number, err = strconv.ParseInt(identifier[:len(identifier)-1], 10, 64)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		return 0, fmt.Errorf("identifier %s not recognized", id)
	}
	return number, err
}

// ParseString parses an interface into a string
func ParseString(s interface{}) string {
	var str string
	switch v := s.(type) {
	case string:
		str = v
	case float64:
		str = fmt.Sprintf("%v", v)
	}
	return str
}
