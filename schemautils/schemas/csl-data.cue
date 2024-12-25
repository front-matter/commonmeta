package compose

import "list"

#schema: {
	// JSON schema for CSL input data. Modified for commonmeta
	@jsonschema(schema="http://json-schema.org/draft-07/schema#")
	@jsonschema(id="https://commonmeta.org/schema/v1.0.1/input/json/csl-data.json")
	_

	#citation: null | bool | number | string | [...] | close({
		type!:           "article" | "article-journal" | "article-magazine" | "article-newspaper" | "bill" | "book" | "broadcast" | "chapter" | "classic" | "collection" | "dataset" | "document" | "entry" | "entry-dictionary" | "entry-encyclopedia" | "event" | "figure" | "graphic" | "hearing" | "interview" | "legal_case" | "legislation" | "manuscript" | "map" | "motion_picture" | "musical_score" | "pamphlet" | "paper-conference" | "patent" | "performance" | "periodical" | "personal_communication" | "post" | "post-weblog" | "regulation" | "report" | "review" | "review-book" | "software" | "song" | "speech" | "standard" | "thesis" | "treaty" | "webpage"
		id!:             number | string
		"citation-key"?: string
		categories?: [...string]
		language?:            string
		journalAbbreviation?: string
		shortTitle?:          string
		author?: [...#["name-variable"]]
		chair?: [...#["name-variable"]]
		"collection-editor"?: [...#["name-variable"]]
		compiler?: [...#["name-variable"]]
		composer?: [...#["name-variable"]]
		"container-author"?: [...#["name-variable"]]
		contributor?: [...#["name-variable"]]
		curator?: [...#["name-variable"]]
		director?: [...#["name-variable"]]
		editor?: [...#["name-variable"]]
		"editorial-director"?: [...#["name-variable"]]
		"executive-producer"?: [...#["name-variable"]]
		guest?: [...#["name-variable"]]
		host?: [...#["name-variable"]]
		interviewer?: [...#["name-variable"]]
		illustrator?: [...#["name-variable"]]
		narrator?: [...#["name-variable"]]
		organizer?: [...#["name-variable"]]
		"original-author"?: [...#["name-variable"]]
		performer?: [...#["name-variable"]]
		producer?: [...#["name-variable"]]
		recipient?: [...#["name-variable"]]
		"reviewed-author"?: [...#["name-variable"]]
		"script-writer"?: [...#["name-variable"]]
		"series-creator"?: [...#["name-variable"]]
		translator?: [...#["name-variable"]]
		accessed?:                #["date-variable"]
		"available-date"?:        #["date-variable"]
		"event-date"?:            #["date-variable"]
		issued?:                  #["date-variable"]
		"original-date"?:         #["date-variable"]
		submitted?:               #["date-variable"]
		abstract?:                string
		annote?:                  string
		archive?:                 string
		archive_collection?:      string
		archive_location?:        string
		"archive-place"?:         string
		authority?:               string
		"call-number"?:           string
		"chapter-number"?:        number | string
		"citation-number"?:       number | string
		"citation-label"?:        string
		"collection-number"?:     number | string
		"collection-title"?:      string
		"container-title"?:       string
		"container-title-short"?: string
		dimensions?:              string
		division?:                string
		DOI?:                     string
		edition?:                 number | string

		// [Deprecated - use 'event-title' instead. Will be removed in
		// 1.1]
		event?:                         string
		"event-title"?:                 string
		"event-place"?:                 string
		"first-reference-note-number"?: number | string
		genre?:                         string
		ISBN?:                          string
		ISSN?:                          string
		issue?:                         number | string
		jurisdiction?:                  string
		keyword?:                       string
		locator?:                       number | string
		medium?:                        string
		note?:                          string
		"number"?:                      number | string
		"number-of-pages"?:             number | string
		"number-of-volumes"?:           number | string
		"original-publisher"?:          string
		"original-publisher-place"?:    string
		"original-title"?:              string
		page?:                          number | string
		"page-first"?:                  number | string
		part?:                          number | string
		"part-title"?:                  string
		PMCID?:                         string
		PMID?:                          string
		printing?:                      number | string
		publisher?:                     string
		"publisher-place"?:             string
		references?:                    string
		"reviewed-genre"?:              string
		"reviewed-title"?:              string
		scale?:                         string
		section?:                       string
		source?:                        string
		status?:                        string
		supplement?:                    number | string
		title?:                         string
		"title-short"?:                 string
		URL?:                           string
		version?:                       string
		volume?:                        number | string
		"volume-title"?:                string
		"volume-title-short"?:          string
		"year-suffix"?:                 string

		// Custom key-value pairs.
		//
		// Used to store additional information that does not have a
		// designated CSL JSON field. The custom field is preferred over
		// the note field for storing custom data, particularly for
		// storing key-value pairs, as the note field is used for user
		// annotations in annotated bibliography styles.
		custom?: {
			...
		}
	})

	#: "name-variable": close({
		family?:                  string
		given?:                   string
		"dropping-particle"?:     string
		"non-dropping-particle"?: string
		suffix?:                  string
		"comma-suffix"?:          bool | number | string
		"static-ordering"?:       bool | number | string
		literal?:                 string
		"parse-names"?:           bool | number | string
	})

	#: "date-variable": close({
		"date-parts"?: list.MaxItems(2) & [...list.MaxItems(3) & [...number | string] & [_, ...]] & [_, ...]
		season?:  number | string
		circa?:   bool | number | string
		literal?: string
		raw?:     string
	})
}
