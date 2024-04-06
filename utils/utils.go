package utils

import (
	"fmt"
	"net/url"
	"strings"
)

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

// ISSN as URL
func IssnAsUrl(issn string) string {
	return fmt.Sprintf("https://portal.issn.org/resource/ISSN/%s", issn)
}
