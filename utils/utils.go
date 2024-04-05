package utils

// ISSN as URL
func IssnAsUrl(issn string) *string {
	if issn == nil {
		return nil
	}
	return "https://portal.issn.org/resource/ISSN/" + issn
}
