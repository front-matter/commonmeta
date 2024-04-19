package authorutils

import "strings"

// IsPersonalName checks if a name is for a Person
func IsPersonalName(name string) bool {
	// personal names are not allowed to contain semicolons
	if strings.Contains(name, ";") {
		return false
	}

	// check if a name has only one word, e.g. "FamousOrganization", not including commas
	if len(strings.Split(name, " ")) == 1 && !strings.Contains(name, ",") {
		return false
	}

	// check if name contains words known to be used in organization names
	organizationWords := []string{
		"University",
		"College",
		"Institute",
		"School",
		"Center",
		"Department",
		"Laboratory",
		"Library",
		"Museum",
		"Foundation",
		"Society",
		"Association",
		"Company",
		"Corporation",
		"Collaboration",
		"Consortium",
		"Incorporated",
		"Inc.",
		"Institut",
		"Research",
		"Science",
		"Team",
		"Ministry",
		"Government",
	}

	for _, word := range organizationWords {
		if strings.Contains(name, word) {
			return false
		}
	}

	// check for suffixes, e.g. "John Smith, MD"
	suffix := strings.Split(name, ", ")[1]
	suffixes := []string{"MD", "PhD", "BS"}
	for _, s := range suffixes {
		if suffix == s {
			return true
		}
	}

	//default to false
	return false
}
