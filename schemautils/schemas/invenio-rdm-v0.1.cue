package compose

import "time"

#schema: {
	// InvenioRDM v0.1
	//
	// JSON schema representation of the InvenioRDM v12 schema.
	@jsonschema(schema="http://json-schema.org/draft-07/schema#")
	_

	#resource: close({
		// The unique identifier of the record.
		id?: string

		// The persistent identifiers of the record.
		pids?: {
			// The digital object identifier (DOI) of the record.
			doi!: {
				// The digital object identifier (DOI).
				identifier?: string

				// The provider of the DOI.
				provider?: "external"
				...
			}
			...
		}

		// The access of the record.
		access?: {
			// The access of the record.
			record!: "public" | "restricted"

			// The access of the files.
			files!: "public" | "restricted"
			...
		}

		// The files of the record.
		files?: {
			// Whether the files are enabled.
			enabled!: bool
			...
		}

		// The metadata of the record.
		metadata?: {
			// The type of the resource.
			resource_type!: {
				// The unique identifier of the resource type.
				id!: "dataset" | "publication-preprint"
				...
			}

			// The creators of the resource.
			creators!: [...{
				// The person or organization.
				person_or_org?: {
					// The type of the person or organization.
					type?: "personal" | "organizational"

					// The given name of the person.
					given_name?: string

					// The family name of the person.
					family_name?: string

					// The name of the organization.
					name?: string

					// The identifiers of the person or organization.
					identifiers?: [...{
						// The identifier of the person or organization.
						identifier?: string

						// The scheme of the identifier.
						scheme?: "orcid" | "ror"
						...
					}]
					...
				}
				...
			}]

			// The title of the resource.
			title!: string

			// The publication date of the resource.
			publication_date!: time.Format("2006-01-02")

			// The subjects of the resource.
			subjects?: [...{
				// The unique identifier of the subject.
				id?: string

				// The title of the subject.
				subject?: string

				// The scheme of the subject.
				scheme?: "FOS"
				...
			}]

			// The dates of the resource.
			dates?: [...{
				// The date of the resource.
				date?: time.Time

				// The type of the date.
				type?: {
					// The unique identifier of the date type.
					id?: "accepted" | "available" | "collected" | "copyrighted" | "created" | "issued" | "other" | "submitted" | "updated" | "valid" | "withdrawn"
					...
				}
				...
			}]

			// The languages of the resource.
			languages?: [...{
				// The ISO-639-3 language code.
				id?: "chi" | "dan" | "dut" | "eng" | "fre" | "ger" | "ita" | "jpn" | "pol" | "por" | "rus" | "spa" | "swe" | "tur"
				...
			}]

			// The identifiers of the resource.
			identifiers?: [...{
				// The identifier of the resource.
				identifier?: string

				// The scheme of the identifier.
				scheme?: "ark" | "arxiv" | "bibcode" | "doi" | "ean13" | "eissn" | "handle" | "igsn" | "isbn" | "issn" | "istc" | "lissn" | "lsid" | "pmid" | "purl" | "upc" | "url" | "urn" | "w3id"
				...
			}]

			// The related identifiers of the resource.
			related_identifiers?: [...{
				// The identifier of the related resource.
				identifier?: string

				// The scheme of the related identifier.
				scheme?: "doi" | "url" | "issn"

				// The type of the relation.
				relation_type?: {
					// The relation type.
					id?: "isnewversionof" | "ispreviousversionof" | "isversionof" | "hasversion" | "ispartof" | "haspart" | "isvariantformof" | "isoriginalformof" | "isidenticalto" | "istranslationof" | "isreviewedby" | "reviews" | "ispreprintof" | "haspreprint" | "issupplementto" | "references"
					...
				}
				...
			}]

			// The rights of the resource.
			rights?: [...{
				// The unique identifier of the rights.
				id?: "cc-by-4.0"
				...
			}]

			// The description of the resource.
			description?: string

			// The funding of the resource.
			funding?: [...{
				// The funder of the resource.
				funder?: {
					// The identifier of the funder.
					id?: string

					// The name of the funder.
					name?: string
					...
				}

				// The award of the resource.
				award?: {
					// The identifier of the award.
					id?: string

					// The number of the award.
					number?: string

					// The title of the award.
					title?: string

					// The identifiers of the award.
					identifiers?: [...{
						// The identifier of the award.
						identifier?: string

						// The scheme of the identifier.
						scheme?: "grid" | "ror" | "doi"
						...
					}]
					...
				}
				...
			}]
			...
		}

		// The custom fields of the record.
		custom_fields?: {
			// The journal of the record.
			"journal:journal"?: {
				// The title of the journal.
				title?: string

				// The volume of the journal.
				volume?: string

				// The issue of the journal.
				issue?: string

				// The pages of the journal.
				pages?: string

				// The International Standard Serial Number (ISSN) of the journal.
				issn?: string
				...
			}
			...
		}
	})
}
