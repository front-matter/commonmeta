package utils

import (
	"commonmeta/doiutils"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
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
	if u.Path[len(u.Path)-1] == '/' {
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
	doi, err := doiutils.ValidateDOI(pid)
	if err != nil {
		return ""
	}
	if doi != "" {
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
	return fmt.Sprintf("https://portal.issn.org/resource/ISSN/%s", issn)
}
