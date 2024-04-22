package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"

	"github.com/front-matter/commonmeta-go/doiutils"

	"github.com/microcosm-cc/bluemonday"
)

// Check for valid DOI or HTTP(S) Url
func NormalizeID(pid string) string {
	// check for valid DOI
	doi := doiutils.NormalizeDOI(pid)
	if doi != "" {
		return doi
	}

	// check for valid HTTP uri and ensure https
	uri, err := url.Parse(pid)
	if err != nil {
		return ""
	}
	if uri.Scheme == "http" {
		uri.Scheme = "https"
	}

	// remove trailing slash
	if pid[len(pid)-1] == '/' {
		pid = pid[:len(pid)-1]
	}

	return pid
}

// Normalize URL
func NormalizeUrl(str string, secure bool, lower bool) (string, error) {
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

// return true if the URL is a Creative Commons License URL
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
	url, err = NormalizeUrl(url, true, false)
	if err != nil {
		return "", false
	}
	normalizedUrl, ok := NormalizedLicenses[url]
	if !ok {
		return url, false
	}
	return normalizedUrl, true
}

func UrlToSPDX(url string) string {
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

type Params struct {
	Pid, Str, Ext, Filename string
	Dct                     map[string]interface{}
}

func FindFromFormat(p Params) string {
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

// Find reader from format by id
func FindFromFormatByID(pid string) string {
	_, ok := doiutils.ValidateDOI(pid)
	if ok {
		return "datacite"
	}
	if strings.Contains(pid, "github.com") {
		return "cff"
	}
	if strings.Contains(pid, "codemeta.json") {
		return "codemeta"
	}
	if strings.Contains(pid, "json_feed_item") {
		return "json_feed_item"
	}
	if strings.Contains(pid, "zenodo.org") {
		return "inveniordm"
	}
	return "schema_org"
}

func FindFromFormatByExt(ext string) string {
	if ext == ".bib" {
		return "bibtex"
	}
	if ext == ".ris" {
		return "ris"
	}
	return ""
}

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

func FindFromFormatByFilename(filename string) string {
	if filename == "CITATION.cff" {
		return "cff"
	}
	return ""
}

// ISSN as URL
func IssnAsUrl(issn string) string {
	if issn == "" {
		return ""
	}
	return fmt.Sprintf("https://portal.issn.org/resource/ISSN/%s", issn)
}

func Sanitize(html string) string {
	policy := bluemonday.StrictPolicy()
	policy.AllowElements("b", "br", "code", "em", "i", "sub", "sup", "strong")
	policy.AllowElements("i")
	sanitizedHTML := policy.Sanitize(html)
	str := strings.Trim(sanitizedHTML, "\n")
	return str
}

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

func NormalizeORCID(orcid string) string {
	orcid_str, ok := ValidateORCID(orcid)
	if !ok {
		return ""
	}
	return "https://orcid.org/" + orcid_str
}

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

func NormalizeROR(ror string) string {
	ror_str, ok := ValidateROR(ror)
	if !ok {
		return ""
	}
	return "https://ror.org/" + ror_str
}

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
