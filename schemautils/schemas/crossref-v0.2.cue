package compose

import "net"

#schema: {
	// Crossref v0.2
	//
	// Unofficial JSON representation of the Crossref content registration schema.
	@jsonschema(schema="http://json-schema.org/draft-07/schema#")
	@jsonschema(id="https://data.crossref.org/schemas/crossref_v0.2.json")

	// The digital object identifier (DOI) of the content.
	doi!: string

	// Content Type describes the type of content registered with
	// Crossref
	type!: "BookChapter" | "BookPart" | "BookSection" | "BookSeries" | "BookSet" | "BookTrack" | "Book" | "Component" | "Database" | "Dataset" | "Dissertation" | "EditedBook" | "Entry" | "Grant" | "JournalArticle" | "JournalIssue" | "JournalVolume" | "Journal" | "Monograph" | "Other" | "PeerReview" | "PostedContent" | "ProceedingsArticle" | "ProceedingsSeries" | "Proceedings" | "ReferenceBook" | "ReferenceEntry" | "ReportComponent" | "ReportSeries" | "Report" | "Standard"

	// The URL for the content.
	url!: net.AbsURL
	...
}
