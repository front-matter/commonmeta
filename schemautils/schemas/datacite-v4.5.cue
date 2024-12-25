package compose

import (
	"net"
	"list"
	"time"
)

#schema: {
	// DataCite v4.5
	//
	// JSON representation of the DataCite v4.5 schema.
	@jsonschema(schema="http://json-schema.org/draft-07/schema#")
	_

	#resource: close({
		id?:  net.AbsURL
		doi?: =~"^10.\\d{4,9}/[-._;()/:a-z0-9A-Z]+$"
		url?: net.AbsURL
		types!: {
			resourceType?:        string
			resourceTypeGeneral!: #resourceTypeGeneral
			...
		}
		creators!: [...#creator & {
			...
		} & {
			name!: _
			...
		}] & [_, ...]
		titles!: list.UniqueItems() & [...{
			title!:     string
			titleType?: #titleType
			lang?:      string
			...
		}] & [_, ...]
		publisher!: {
			name!:                      string
			publisherIdentifier?:       string
			publisherIdentifierScheme?: string
			schemeURI?:                 net.AbsURL
			lang?:                      string
			...
		}
		publicationYear!: =~"^[0-9]{4}$"
		subjects?: list.UniqueItems() & [...{
			subject!:            string
			subjectScheme?:      string
			schemeUri?:          net.AbsURL
			valueUri?:           net.AbsURL
			classificationCode?: string
			lang?:               string
			...
		}]
		contributors?: [...#contributor & {
			...
		} & {
			contributorType!: _
			name!:            _
			...
		}]
		dates?: list.UniqueItems() & [...{
			date!:            #date
			dateType!:        #dateType
			dateInformation?: string
			...
		}]
		language?: string
		alternateIdentifiers?: list.UniqueItems() & [...{
			alternateIdentifier!:     string
			alternateIdentifierType!: string
			...
		}]
		relatedIdentifiers?: [...#relatedObject & {
			...
		} & matchIf(#relatedObjectIf & {
			...
		}, _, #relatedObjectElse & {
			...
		}) & {
			relatedIdentifier!:     string
			relatedIdentifierType!: #relatedIdentifierType
			relationType!:          _
			...
		}]
		relatedItems?: list.UniqueItems() & [...#relatedObject & {
			...
		} & matchIf(#relatedObjectIf & {
			...
		}, _, #relatedObjectElse & {
			...
		}) & {
			relatedItemIdentifier?: {
				relatedItemIdentifier!:     string
				relatedItemIdentifierType!: #relatedIdentifierType
				...
			}
			relatedItemType!: #resourceTypeGeneral
			creators?: [...#creator & {
				...
			} & {
				name!: _
				...
			}]
			contributors?: [...#contributor & {
				...
			} & {
				contributorType!: _
				name!:            _
				...
			}]
			titles!: list.UniqueItems() & [...{
				title!:     string
				titleType?: #titleType
				lang?:      string
				...
			}] & [_, ...]
			publicationYear?: =~"^[0-9]{4}$"
			volume?:          string
			issue?:           string
			firstPage?:       string
			lastPage?:        string
			edition?:         string
			publisher?:       string
			number?:          string
			numberType?:      "Article" | "Chapter" | "Report" | "Other"
			relationType!:    _
			...
		}]
		sizes?: list.UniqueItems() & [...string]
		formats?: list.UniqueItems() & [...string]
		version?: string
		rightsList?: list.UniqueItems() & [...{
			rights?:                 string
			rightsUri?:              net.AbsURL
			rightsIdentifier?:       string
			rightsIdentifierScheme?: string
			schemeUri?:              net.AbsURL
			lang?:                   string
			...
		}]
		descriptions?: list.UniqueItems() & [...{
			description!:     string
			descriptionType!: #descriptionType
			lang?:            string
			...
		}]
		geoLocations?: list.UniqueItems() & [...{
			geoLocationPlace?: string
			geoLocationPoint?: #geoLocationPoint
			geoLocationBox?: {
				westBoundLongitude!: #longitude
				eastBoundLongitude!: #longitude
				southBoundLatitude!: #latitude
				northBoundLatitude!: #latitude
				...
			}
			geoLocationPolygons?: list.UniqueItems() & [...{
				polygonPoints!: [...#geoLocationPoint] & [_, _, _, _, ...]
				inPolygonPoint?: #geoLocationPoint
				...
			}]
			...
		}]
		fundingReferences?: [...{
			funderName!:           string
			funderIdentifier?:     string
			funderIdentifierType?: #funderIdentifierType
			awardNumber?:          string
			awardUri?:             net.AbsURL
			awardTitle?:           string
			...
		}]
		schemaVersion!: "http://datacite.org/schema/kernel-4"
		container?: {
			type?:      string
			title?:     string
			firstPage?: string
			...
		}
	})

	#nameType: "Organizational" | "Personal"

	#nameIdentifiers: list.UniqueItems() & [...{
		nameIdentifier!:       string
		nameIdentifierScheme!: string
		schemeUri?:            net.AbsURL
		...
	}]

	#affiliation: list.UniqueItems() & [...{
		name!:                        string
		affiliationIdentifier?:       string
		affiliationIdentifierScheme?: string
		schemeUri?:                   net.AbsURL
		...
	}]

	#creator: {
		name!:            string
		nameType?:        #nameType
		givenName?:       string
		familyName?:      string
		nameIdentifiers?: #nameIdentifiers
		affiliation?:     #affiliation
		lang?:            string
		...
	}

	#contributor: #creator & {
		...
	} & {
		contributorType!: #contributorType
		name!:            _
		...
	}

	#contributorType: "ContactPerson" | "DataCollector" | "DataCurator" | "DataManager" | "Distributor" | "Editor" | "HostingInstitution" | "Producer" | "ProjectLeader" | "ProjectManager" | "ProjectMember" | "RegistrationAgency" | "RegistrationAuthority" | "RelatedPerson" | "Researcher" | "ResearchGroup" | "RightsHolder" | "Sponsor" | "Supervisor" | "WorkPackageLeader" | "Other"

	#titleType: "AlternativeTitle" | "Subtitle" | "TranslatedTitle" | "Other"

	#longitude: <=180 & >=-180

	#latitude: <=90 & >=-90

	#date: matchN(>=1, [string, string, time.Format("2006-01-02"), string, string, string, string, string])

	#dateType: "Accepted" | "Available" | "Copyrighted" | "Collected" | "Created" | "Issued" | "Submitted" | "Updated" | "Valid" | "Withdrawn" | "Other"

	#resourceTypeGeneral: "Audiovisual" | "Book" | "BookChapter" | "Collection" | "ComputationalNotebook" | "ConferencePaper" | "ConferenceProceeding" | "DataPaper" | "Dataset" | "Dissertation" | "Event" | "Image" | "Instrument" | "InteractiveResource" | "Journal" | "JournalArticle" | "Model" | "OutputManagementPlan" | "PeerReview" | "PhysicalObject" | "Preprint" | "Report" | "Service" | "Software" | "Sound" | "Standard" | "StudyRegistration" | "Text" | "Workflow" | "Other"

	#relatedIdentifierType: "ARK" | "arXiv" | "bibcode" | "DOI" | "EAN13" | "EISSN" | "Handle" | "IGSN" | "ISBN" | "ISSN" | "ISTC" | "LISSN" | "LSID" | "PMID" | "PURL" | "UPC" | "URL" | "URN" | "w3id"

	#relationType: "IsCitedBy" | "Cites" | "IsCollectedBy" | "Collects" | "IsSupplementTo" | "IsSupplementedBy" | "IsContinuedBy" | "Continues" | "IsDescribedBy" | "Describes" | "HasMetadata" | "IsMetadataFor" | "HasVersion" | "IsVersionOf" | "IsNewVersionOf" | "IsPartOf" | "IsPreviousVersionOf" | "IsPublishedIn" | "HasPart" | "IsReferencedBy" | "References" | "IsDocumentedBy" | "Documents" | "IsCompiledBy" | "Compiles" | "IsVariantFormOf" | "IsOriginalFormOf" | "IsIdenticalTo" | "IsReviewedBy" | "Reviews" | "IsDerivedFrom" | "IsSourceOf" | "IsRequiredBy" | "Requires" | "IsObsoletedBy" | "Obsoletes"

	#relatedObject: {
		relationType!:          #relationType
		relatedMetadataScheme?: string
		schemeUri?:             net.AbsURL
		schemeType?:            string
		resourceTypeGeneral?:   #resourceTypeGeneral
		...
	}

	#relatedObjectIf: null | bool | number | string | [...] | {
		relationType?: "HasMetadata" | "IsMetadataFor"
		...
	}

	#relatedObjectElse: null | bool | number | string | [...] | {
		relatedMetadataScheme?: _|_
		schemeUri?:             _|_
		schemeType?:            _|_
		...
	}

	#descriptionType: "Abstract" | "Methods" | "SeriesInformation" | "TableOfContents" | "TechnicalInfo" | "Other"

	#geoLocationPoint: {
		pointLongitude!: #longitude
		pointLatitude!:  #latitude
		...
	}

	#funderIdentifierType: "ISNI" | "GRID" | "Crossref Funder ID" | "ROR" | "Other"
}
