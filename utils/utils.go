package utils

import "fmt"

// ISSN as URL
func IssnAsUrl(issn string) string {
	return fmt.Sprintf("https://portal.issn.org/resource/ISSN/%s", issn)
}
