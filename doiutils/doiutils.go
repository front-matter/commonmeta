package doiutils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

// Validate a DOI prefix for a given DOI
func ValidatePrefix(doi string) (string, bool) {
	r, err := regexp.Compile(`^(?:(http|https):/(/)?(dx\.)?(doi\.org|handle\.stage\.datacite\.org|handle\.test\.datacite\.org)/)?(doi:)?(10\.\d{4,5})/.+$`)
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

// return the DOI registration agency for a given DOI
func GetDOIRA(doi string) (string, bool) {
	prefix, ok := ValidatePrefix(doi)
	if !ok {
		return "", false
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

// Get the Crossref member name for a given member_id
func GetCrossrefMember(memberId string) (string, bool) {
	type Response struct {
		Message struct {
			PrimaryName string `json:"primary-name"`
		} `json:"message"`
	}
	var result Response
	if memberId == "" {
		return "", false
	}
	resp, err := http.Get(fmt.Sprintf("https://api.crossref.org/members/%s", memberId))
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
	return result.Message.PrimaryName, true
}
