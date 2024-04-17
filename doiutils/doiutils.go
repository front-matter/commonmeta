package doiutils

import (
	"log"
	"net/url"
	"regexp"
	"strings"
)

// extract DOI prefix from URL
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

func DOIAsUrl(str string) string {
	if str == "" {
		return ""
	}
	return "https://doi.org/" + strings.ToLower(str)
}

// Normalize a DOI
func NormalizeDOI(doi string) string {
	doistr, ok := ValidateDOI(doi)
	if !ok {
		return ""
	}
	resolver := DOIResolver(doi, false)
	return resolver + strings.ToLower(doistr)
}

// Validate a DOI
func ValidateDOI(doi string) (string, bool) {
	r, err := regexp.Compile(`^(?:(http|https):/(/)?(dx\.)?(doi\.org|handle\.stage\.datacite\.org|handle\.test\.datacite\.org)/)?(doi:)?(10\.\d{4,5}/.+)$`)
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

// Return a DOI resolver for a given DOI
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
