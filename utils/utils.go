// Package utils provides utility functions for commonmeta.
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	iso639_3 "github.com/barbashov/iso639-3"
	"github.com/front-matter/commonmeta/crockford"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/spdx"
	"github.com/microcosm-cc/bluemonday"
	"github.com/pkosilo/iso7064"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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

	// check for valid Wikidata item ID
	wikidata, ok := ValidateWikidata(pid)
	if ok {
		return "https://www.wikidata.org/wiki/" + wikidata
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

// NormalizeWorkID normalizes work ID
func NormalizeWorkID(id string) string {
	pid, type_, category := ValidateIDCategory(id)
	var allowedCategories = []string{"Work", "All"}
	if !slices.Contains(allowedCategories, category) {
		return ""
	}
	if type_ == "DOI" {
		return doiutils.NormalizeDOI(pid)
	} else if type_ == "UUID" {
		return pid
	} else if type_ == "URL" {
		return pid
	} else if type_ == "Wikidata" {
		return "https://www.wikidata.org/wiki/" + pid
	}
	return ""
}

// NormalizeOrganizationID normalizes organization ID
func NormalizeOrganizationID(id string) string {
	pid, type_, category := ValidateIDCategory(id)
	var allowedCategories = []string{"Organization", "Contributor", "All"}
	if !slices.Contains(allowedCategories, category) {
		return ""
	}
	if type_ == "ROR" {
		return "https://ror.org/" + pid
	} else if type_ == "Crossref Funder ID" {
		return "https://doi.org/" + pid
	} else if type_ == "GRID" {
		return "https://grid.ac/institutes/" + pid
	} else if type_ == "Wikidata" {
		return "https://www.wikidata.org/wiki/" + pid
	} else if type_ == "ISNI" {
		return "https://isni.org/isni/" + pid
	}
	return ""
}

// NormalizePersonID normalizes person ID
func NormalizePersonID(id string) string {
	pid, type_, category := ValidateIDCategory(id)
	var allowedCategories = []string{"Person", "Contributor", "All"}
	if !slices.Contains(allowedCategories, category) {
		return ""
	}
	if type_ == "ORCID" {
		return "https://orcid.org/" + pid
	} else if type_ == "ISNI" {
		return "https://isni.org/isni/" + pid
	} else if type_ == "Wikidata" {
		return "https://www.wikidata.org/wiki/" + pid
	}
	return ""
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
	if secure && u.Scheme == "http" {
		u.Scheme = "https"
	}
	if lower {
		return strings.ToLower(u.String()), nil
	}
	return u.String(), nil
}

// NormalizeCCUrl returns the normalized Creative Commons License URL
func NormalizeCCUrl(url_ string) (string, bool) {
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

	if url_ == "" {
		return "", false
	}
	var err error
	url_, err = NormalizeURL(url_, true, false)
	if err != nil {
		return "", false
	}
	u, err := url.Parse(url_)
	if err != nil {
		return "", false
	}
	// strip trailing slash if no query
	if u.Path != "" && len(u.RawQuery) == 0 && u.Path[len(u.Path)-1] == '/' {
		u.Path = u.Path[:len(u.Path)-1]
	}
	normalizedURL, ok := NormalizedLicenses[u.String()]
	return normalizedURL, ok
}

// URLToSPDX provides the SPDX license ID given a Creative Commons URL
func URLToSPDX(url string) string {
	license, err := spdx.Search(url)
	if err != nil {
		return ""
	}
	return license.LicenseID
}

// SPDXToURL provides the SPDX license URL given a Creative Commons SPDX ID
func SPDXToURL(id string) string {
	license, err := spdx.Search(id)
	if err != nil {
		return ""
	}
	if license.SeeAlso == nil || len(license.SeeAlso) == 0 {
		return ""
	}
	return license.SeeAlso[0]
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
	r := regexp.MustCompile(`^(?:https://portal\.issn\.org/resource/ISSN/)?(\d{4}\-\d{3}(\d|x|X))$`)
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
	r := regexp.MustCompile(`^(?:(?:http|https)://(?:(?:www|sandbox)?\.)?orcid\.org/)?(000[09][ -]000[123][ -]\d{4}[ -]\d{3}[0-9X]+)$`)
	matched := r.FindStringSubmatch(orcid)
	if len(matched) == 0 {
		return "", false
	}
	return matched[1], CheckORCIDNumberRange(matched[1])
}

// ValidateISNI validates an ISNI
// ISNI is a 16-character string in blocks of four
// optionally separated by hyphens or spaces and NOT
// between 0000-0001-5000-0007 and 0000-0003-5000-0001,
// or between 0009-0000-0000-0000 and 0009-0010-0000-0000
// (the ranged reserved for ORCID).
func ValidateISNI(isni string) (string, bool) {
	r := regexp.MustCompile(`^(?:(?:http|https)://(?:(?:www)?\.)?isni\.org/)?(?:isni/)?(0000[ -]?00\d{2}[ -]?\d{4}[ -]?\d{3}[0-9X]+)$`)
	matched := r.FindStringSubmatch(isni)
	if len(matched) == 0 {
		return "", false
	}
	// workaround until regex capture group is fixed
	match := strings.ReplaceAll(matched[1], " ", "")
	match = strings.ReplaceAll(match, "-", "")
	return match, !CheckORCIDNumberRange(match)
}

// check if ORCID is in the range 0000-0001-5000-0007 and 0000-0003-5000-0001
// or between 0009-0000-0000-0000 and 0009-0010-0000-0000
func CheckORCIDNumberRange(orcid string) bool {
	number := strings.ReplaceAll(orcid, "-", "")
	if number >= "0000000150000007" && number <= "0000000350000001" {
		return true
	}
	if number >= "0009000000000000" && number <= "0009001000000000" {
		return true
	}
	return false
}

// ValidateWikidata validates a Wikidata item ID
// Wikidata item ID is a string prefixed with Q followed by a number
func ValidateWikidata(wikidata string) (string, bool) {
	r := regexp.MustCompile(`^(?:(?:http|https)://(?:(?:www)?\.)?wikidata\.org/wiki/)?(Q\d+)$`)
	matched := r.FindStringSubmatch(wikidata)
	if len(matched) == 0 {
		return "", false
	}
	return matched[1], true
}

// ValidateGRID validates a GRID ID
// GRID ID is a string prefixed with grid followed by dot number dot string
func ValidateGRID(grid string) (string, bool) {
	r := regexp.MustCompile(`^(?:(?:http|https)://(?:(?:www)?\.)?grid\.ac/)?(?:institutes/)?(grid\.[0-9]+\.[a-f0-9]{1,2})$`)
	matched := r.FindStringSubmatch(grid)
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
	r := regexp.MustCompile(`^(?:(?:http|https)://ror\.org/)?(0[0-9a-z]{6}\d{2})$`)
	matched := r.FindStringSubmatch(ror)
	if len(matched) == 0 {
		return "", false
	}
	return matched[1], true
}

// ValidateOpenalex validates an OpenAlex ID. The first letter indicates the type of resource
// (A author, F funder, I institution, P publisher, S source W work), followed by 8-10 digits.
func ValidateOpenalex(openalex string) (string, bool) {
	r := regexp.MustCompile(`^(?:(?:http|https)://openalex\.org/)?([AFIPSW]\d{8,10})$`)
	matched := r.FindStringSubmatch(openalex)
	if len(matched) == 0 {
		return "", false
	}
	return matched[1], true
}

// ValidatePMID validates a PubdMed ID
func ValidatePMID(pmid string) (string, bool) {
	r := regexp.MustCompile(`^(?:(?:http|https)://pubmed\.ncbi\.nlm\.nih\.gov/)?(\d{4,8})$`)
	matched := r.FindStringSubmatch(pmid)
	if len(matched) == 0 {
		return "", false
	}
	return matched[1], true
}

// ValidatePMCID validates a PubMed Central ID
func ValidatePMCID(pmcid string) (string, bool) {
	r := regexp.MustCompile(`^(?:(?:http|https)://www\.ncbi\.nlm\.nih\.gov/pmc/articles/)?(\d{4,8})$`)
	matched := r.FindStringSubmatch(pmcid)
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
	if resp.StatusCode != http.StatusOK {
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

func ValidateCrossrefFunderID(fundref string) (string, bool) {
	r, err := regexp.Compile(`^(?:https?://doi\.org/)?(?:10\.13039/)?((501)?1000[0-9]{5})$`)
	if err != nil {
		fmt.Println("error:", err)
		return "", false
	}
	// r := regexp.MustCompile(`^(?:(http|https):/(/)?(dx\.)?(doi\.org/))?(?:10\.13039/)?((501)?1000[0-1]+[0-9]{4})$`)
	matched := r.FindStringSubmatch(fundref)
	if len(matched) == 0 {
		return "", false
	}
	return matched[1], true
}

// ValidateID validates an identifier and returns the type
// Can be DOI, UUID, ISSN, ORCID, ROR, URL, RID, Wikidata, ISNI, OpenAlex,
// PMID, PMCID or GRID
func ValidateID(id string) (string, string) {
	fundref, ok := ValidateCrossrefFunderID(id)
	if ok {
		return fundref, "Crossref Funder ID"
	}
	doi, ok := doiutils.ValidateDOI(id)
	if ok {
		return doi, "DOI"
	}
	uuid, ok := ValidateUUID(id)
	if ok {
		return uuid, "UUID"
	}
	pmid, ok := ValidatePMID(id)
	if ok {
		return pmid, "PMID"
	}
	pmcid, ok := ValidatePMCID(id)
	if ok {
		return pmcid, "OpenAlex"
	}
	openalex, ok := ValidateOpenalex(id)
	if ok {
		return openalex, "OpenAlex"
	}
	orcid, ok := ValidateORCID(id)
	if ok {
		return orcid, "ORCID"
	}
	ror, ok := ValidateROR(id)
	if ok {
		return ror, "ROR"
	}
	grid, ok := ValidateGRID(id)
	if ok {
		return grid, "GRID"
	}
	rid, ok := ValidateRID(id)
	if ok {
		return rid, "RID"
	}
	wikidata, ok := ValidateWikidata(id)
	if ok {
		return wikidata, "Wikidata"
	}
	isni, ok := ValidateISNI(id)
	if ok {
		return isni, "ISNI"
	}
	issn, ok := ValidateISSN(id)
	if ok {
		return issn, "ISSN"
	}
	url := ValidateURL(id)
	if url != "" {
		return id, url
	}
	return "", ""
}

// ValidateIDCategory validates an identifier and returns the identifier,
// type, and category
// Category can be work, person, organization, contributor (person or organization),
// or all
func ValidateIDCategory(id string) (string, string, string) {
	id, type_ := ValidateID(id)
	switch type_ {
	case "ROR", "Crossref Funder ID", "GRID":
		return id, type_, "Organization"
	case "ORCID":
		return id, type_, "Person"
	case "ISNI":
		return id, type_, "Contributor"
	case "DOI", "PMID", "PMCID":
		return id, type_, "Work"
	case "Wikidata", "OpenAlex", "URL", "UUID":
		return id, type_, "All"
	default:
		return id, type_, ""
	}
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
	if u.Scheme == "https" && u.Host == "api.rogue-scholar.org" {
		path := strings.Split(u.Path, "/")
		if len(path) == 3 && path[1] == "posts" {
			_, ok = ValidateUUID(path[2])
			if ok {
				return "JSONFEEDID"
			}
		} else if len(path) == 4 && path[1] == "posts" {
			_, ok = doiutils.ValidateDOI(path[2] + "/" + path[3])
			if ok {
				return "JSONFEEDID"
			}
		}
	} else if u.Scheme == "http" || u.Scheme == "https" {
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

// ValidateRID validates a RID
// RID is the unique identifier used by the InvenioRDM platform
func ValidateRID(rid string) (string, bool) {
	r := regexp.MustCompile("^[" + crockford.ENCODING_CHARS + "]{5}-[" + crockford.ENCODING_CHARS + "]{3}[0-9]{2}$")
	if !r.MatchString(rid) {
		return "", false
	}
	return r.FindString(rid), true
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
	s = strings.ReplaceAll(s, "-", "")
	if s == "" {
		return s
	}
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

// StringToSlug makes a string lowercase and removes non-alphanumeric characters
func StringToSlug(str string) string {
	s, _ := NormalizeString(str)
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, s)
}

func NormalizeString(s string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, s)
	if err != nil {
		return "", err
	}

	return result, nil
}

// func SplitString adds a hyphen every n characters
func SplitString(str string, n int, s string) string {
	if n <= 0 {
		return str
	}

	var splits []string
	for i := 0; i < len(str); i += n {
		nn := i + n
		if nn > len(str) {
			nn = len(str)
		}
		splits = append(splits, string(str[i:nn]))
	}
	return strings.Join(splits, s)
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
	} else if identifierType == "RID" {
		// RID is a 10-character string with a hyphen after five digits.
		// It is a base32-encoded numbers with checksum.
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
