package doiutils

import (
	"net/url"
	"regexp"
	"strings"
)

// extract DOI from URL
func DOIFromUrl(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	if u.Host == "" {
		return str, nil
	}
	if u.Host != "doi.org" || !strings.HasPrefix(u.Path, "/10.") {
		return "", nil
	}
	return strings.TrimLeft(u.Path, "/"), nil
}

// def doi_from_url(url: Optional[str]) -> Optional[str]:
//     """Return a DOI from a URL"""
//     if url is None:
//         return None

//     f = furl(url)
//     # check for allowed scheme if string is a URL
//     if f.host is not None and f.scheme not in ["http", "https", "ftp"]:
//         return None

//     # url is for a short DOI
//     if f.host == "doi.org" and not f.path.segments[0].startswith("10."):
//         return short_doi_as_doi(url)

//     # special rules for specific hosts
//     if f.host == "onlinelibrary.wiley.com":
//         if f.path.segments[-1] in ["epdf"]:
//             f.path.segments.pop()
//     elif f.host == "www.plosone.org":
//         if (
//             f.path.segments[-1] in ["fetchobject.action"]
//             and f.args.get("uri", None) is not None
//         ):
//             f.path = f.args.get("uri")
//     path = str(f.path)
//     match = re.search(
//         r"(10\.\d{4,5}/.+)\Z",
//         path,
//     )
//     if match is None:
//         return None
//     return match.group(0).lower()

// Normalize a DOI
func NormalizeDOI(doi string) string {
	doistr, err := ValidateDOI(doi)
	if err != nil {
		return ""
	}
	resolver := DOIResolver(doi, false)
	return resolver + strings.ToLower(doistr)
}

// Validate a DOI
func ValidateDOI(doi string) (string, error) {
	matched, err := regexp.MatchString(`^(?:(http|https):/(/)?(dx\.)?(doi\.org|handle\.stage\.datacite\.org|handle\.test\.datacite\.org)/)?(doi:)?(10\.\d{4,5}/.+)$`, doi)
	if err != nil {
		return "", err
	}
	if !matched {
		return "", nil
	}
	return doi, nil
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
