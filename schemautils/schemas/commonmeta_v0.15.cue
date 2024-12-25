package compose

import (
	"net"
	"list"
)

#commonmeta: {
	// Commonmeta v0.15
	//
	// JSON representation of the Commonmeta schema.
	@jsonschema(schema="http://json-schema.org/draft-07/schema#")
	@jsonschema(id="https://commonmeta.org/commonmeta_v0.15.json")
	_

	#affiliations: [...{
		organization?: #organization
		...
	}]

	#schema: close({
		id!:   #id
		type!: #type

		// The additional type of the resource.
		additionalType?: string

		// The location where content is archived.
		archiveLocations?: [..."CLOCKSS" | "LOCKSS" | "Portico" | "KB" | "Internet Archive" | "DWT"]

		// The container of the resource.
		container?: {
			// The identifier for the container.
			identifier?: string

			// The identifierType for the container.
			identifierType?: string

			// The type of the container.
			type?: "Book" | "BookSeries" | "Journal" | "Proceedings" | "ProceedingsSeries" | "Repository" | "DataRepository" | "Periodical" | "Series"

			// The title of the container.
			title?: string

			// The first page of the resource.
			firstPage?: string

			// The last page of the resource.
			lastPage?: string

			// The volume of the resource.
			volume?: string

			// The issue of the resource.
			issue?: string
			...
		}

		// The contributors to the resource.
		contributors?: [...{
			organization?: #organization
			person?:       #person

			// List of roles assumed by the contributor when working on the
			// resource.
			contributorRoles?: [...#contributorRole]
			...
		}] & [_, ...]

		// The dates for the resource.
		date?: {
			// The date the resource was created.
			created?: string

			// The date the resource was submitted.
			submitted?: string

			// The date the resource was accepted.
			accepted?: string

			// The date the resource was published.
			published?: string

			// The date the resource was updated.
			updated?: string

			// The date the resource was accessed.
			accessed?: string

			// The date the resource was made available.
			available?: string

			// The date the resource was withdrawn.
			withdrawn?: string
			...
		}

		// The descriptions of the resource.
		descriptions?: [...{
			// The description of the resource.
			description!: string

			// The type of the description.
			type?: "Abstract" | "Summary" | "Methods" | "TechnicalInfo" | "Other"

			// The language of the title. Use one of the language codes from
			// the IETF BCP 47 standard.
			language?: string
			...
		}]

		// The downloadable files for the resource.
		files?: [...{
			bucket?:   string
			key?:      string
			checksum?: string
			url!:      net.AbsURL
			size?:     int
			mimeType?: string
			...
		}] & [_, ...]

		// The funding references for the resource.
		fundingReferences?: [...{
			funderIdentifier?:     string
			funderIdentifierType?: "Crossref Funder ID" | "ROR" | "GRID" | "ISNI" | "Ringgold" | "Other"
			funderName!:           string
			awardNumber?:          string
			awardTitle?:           string
			awardUri?:             net.AbsURL
			...
		}]
		geoLocations?: list.UniqueItems() & [...{
			geoLocationPlace?: string
			geoLocationPoint?: #geoLocationPoint
			geoLocationBox?:   #geoLocationBox
			geoLocationPolygons?: list.UniqueItems() & [...{
				polygonPoints!: [...#geoLocationPoint] & [_, _, _, _, ...]
				inPolygonPoint?: #geoLocationPoint
				...
			}]
			...
		}]

		// Identifiers for the resource, including the id.
		identifiers?: [...{
			identifier!:     string
			identifierType!: "ARK" | "arXiv" | "Bibcode" | "DOI" | "GUID" | "Handle" | "ISBN" | "ISSN" | "PMID" | "PMCID" | "PURL" | "RID" | "URL" | "URN" | "UUID" | "Other"
			...
		}]

		// The language of the resource. Use one of the language codes
		// from the IETF BCP 47 standard.
		language?: string

		// The license for the resource. Use one of the SPDX license
		// identifiers.
		license?: {
			id?:  string
			url?: net.AbsURL
			...
		}

		// The provider of the resource. This can be a DOI registration
		// agency or a repository.
		provider?: "Crossref" | "DataCite" | "GitHub" | "JaLC" | "KISTI" | "mEDRA" | "OP"

		// The publisher of the resource.
		publisher?: {
			organization?: #organization
			...
		}

		// Other resolvable persistent unique IDs related to the resource.
		relations?: [...{
			id!:   net.AbsURL
			type!: "IsNewVersionOf" | "IsPreviousVersionOf" | "IsVersionOf" | "HasVersion" | "IsPartOf" | "HasPart" | "IsVariantFormOf" | "IsOriginalFormOf" | "IsIdenticalTo" | "IsTranslationOf" | "HasTranslation" | "IsReviewedBy" | "Reviews" | "HasReview" | "IsPreprintOf" | "HasPreprint" | "IsSupplementTo" | "IsSupplementedBy"
			...
		}] & [_, ...]
		references?: [...{
			id?:              #id
			type?:            #type
			key!:             string
			contributor?:     string
			title?:           string
			publisher?:       string
			publicationYear?: string
			volume?:          string
			issue?:           string
			firstPage?:       string
			lastPage?:        string
			containerTitle?:  string
			edition?:         string
			unstructured?:    string
			...
		}]
		subjects?: [...{
			subject!: string

			// The language of the subject. Use one of the language codes from
			// the IETF BCP 47 standard.
			language?: string
			...
		}]

		// The titles of the resource.
		titles?: [...{
			// The title of the resource.
			title!: string

			// The type of the title.
			type?: "AlternativeTitle" | "Subtitle" | "TranslatedTitle"

			// The language of the title. Use one of the language codes from
			// the IETF BCP 47 standard.
			language?: string
			...
		}]

		// The URL of the resource.
		url?: net.AbsURL

		// The version of the resource.
		version?: string
	})

	#contributorRole: "Author" | "Editor" | "Chair" | "Reviewer" | "ReviewAssistant" | "StatsReviewer" | "ReviewerExternal" | "Reader" | "Translator" | "ContactPerson" | "DataCollector" | "DataManager" | "Distributor" | "HostingInstitution" | "Producer" | "ProjectLeader" | "ProjectManager" | "ProjectMember" | "RegistrationAgency" | "RegistrationAuthority" | "RelatedPerson" | "ResearchGroup" | "RightsHolder" | "Researcher" | "Sponsor" | "WorkPackageLeader" | "Conceptualization" | "DataCuration" | "FormalAnalysis" | "FundingAcquisition" | "Investigation" | "Methodology" | "ProjectAdministration" | "Resources" | "Software" | "Supervision" | "Validation" | "Visualization" | "WritingOriginalDraft" | "WritingReviewEditing" | "Maintainer" | "Other"

	#geoLocationBox: {
		westBoundLongitude?: #longitude
		eastBoundLongitude?: #longitude
		southBoundLatitude?: #latitude
		northBoundLatitude?: #latitude
		...
	}

	#geoLocationPoint: {
		pointLongitude?: #longitude
		pointLatitude?:  #latitude
		...
	}

	#id: net.AbsURL

	#latitude: >=-90 & <=90

	#longitude: >=-180 & <=180

	#organization: {
		// The unique identifier for the organization.
		id?:   net.AbsURL
		type!: "Organization"

		// The name of the organization.
		name!: string
		...
	}

	#person: {
		id?:   net.AbsURL
		type!: "Person"

		// The given name of the person.
		givenName?: string

		// The family name of the person.
		familyName!:  string
		affiliation?: #affiliations
		...
	}

	#type: "Article" | "Audiovisual" | "BookChapter" | "BookPart" | "BookSection" | "BookSeries" | "BookSet" | "Book" | "Collection" | "Component" | "Database" | "Dataset" | "Dissertation" | "Document" | "Entry" | "Event" | "Grant" | "Image" | "Instrument" | "InteractiveResource" | "JournalArticle" | "JournalIssue" | "JournalVolume" | "Journal" | "PeerReview" | "PhysicalObject" | "Poster" | "Presentation" | "ProceedingsArticle" | "ProceedingsSeries" | "Proceedings" | "ReportComponent" | "ReportSeries" | "Report" | "Software" | "Standard" | "StudyRegistration" | "WebPage" | "Other"
}
