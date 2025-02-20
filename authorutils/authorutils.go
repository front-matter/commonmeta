// Package authorutils provides utility functions to work with authors
package authorutils

import (
	"strings"
)

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
		"Redaktion",
		"Count",
	}

	for _, word := range organizationWords {
		if strings.Contains(name, word) {
			return false
		}
	}

	// check for suffixes, e.g. "John Smith, MD"
	suffix := strings.Split(name, ", ")
	if len(suffix) > 1 {
		suffixes := []string{"MD", "PhD", "BS"}
		for _, s := range suffixes {
			if suffix[1] == s {
				return true
			}
		}
	}

	//default to true
	return true
}

func ParseName(name string) (string, string, string) {
	var givenName, familyName string

	if !IsPersonalName(name) {
		return givenName, familyName, name
	}

	// check for suffixes, e.g. "John Smith, MD"
	suffix := strings.Split(name, ", ")
	if len(suffix) > 1 {
		suffixes := []string{"MD", "PhD", "BS"}
		for _, s := range suffixes {
			if suffix[1] == s {
				name = suffix[0]
				break
			}
		}
	}

	// check for comma separated names, e.g. "Doe, John"
	comma := strings.Split(name, ", ")
	if len(comma) > 1 {
		givenName = comma[1]
		familyName = comma[0]
		return givenName, familyName, ""
	}

	// default to the last word as family name
	words := strings.Split(name, " ")
	if len(words) == 1 {
		familyName = name
		return givenName, familyName, ""
	} else if len(words) > 1 {
		familyName = words[len(words)-1]
		givenName = strings.Join(words[:len(words)-1], " ")
		name = ""
	}
	return givenName, familyName, name
}
