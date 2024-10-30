package inveniordm

import (
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/crossrefxml"
	"github.com/front-matter/commonmeta/dateutils"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/schemautils"
	"github.com/front-matter/commonmeta/utils"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

// Vocabularies is the embedded vocabulary yaml files.
//
//go:embed vocabularies/*.yaml
var Vocabularies embed.FS

// Convert converts Commonmeta metadata to InvenioRDM metadata
func Convert(data commonmeta.Data) (Inveniordm, error) {
	var inveniordm Inveniordm

	// load awards vocabulary
	var awardsVocabulary []AwardVocabulary
	yamlAwardsVocabulary, _ := Vocabularies.ReadFile(filepath.Join("vocabularies", "awards.yaml"))
	err := yaml.Unmarshal(yamlAwardsVocabulary, &awardsVocabulary)
	if err != nil {
		panic(err)
	}

	// required properties
	doi, _ := doiutils.ValidateDOI(data.ID)
	inveniordm.Pids.DOI = DOI{
		Identifier: doi,
		Provider:   "external",
	}
	inveniordm.Access.Record = "public"
	inveniordm.Access.Files = "public"
	inveniordm.Files.Enabled = false

	inveniordm.Metadata.ResourceType = ResourceType{
		ID: CMToInvenioMappings[data.Type],
	}
	if len(data.Titles) > 0 {
		inveniordm.Metadata.Title = data.Titles[0].Title
	}
	if len(data.Date.Published) >= 4 {
		inveniordm.Metadata.PublicationDate = dateutils.GetDateFromDatetime(data.Date.Published)
	}

	if len(data.Contributors) > 0 {
		for _, v := range data.Contributors {
			var identifiers []Identifier
			if v.ID != "" {
				id, _ := utils.ValidateORCID(v.ID)
				Identifier := Identifier{
					Identifier: id,
					Scheme:     "orcid",
				}
				identifiers = append(identifiers, Identifier)
			}
			var affiliations []Affiliation
			for _, a := range v.Affiliations {
				id, _ := utils.ValidateROR(a.ID)
				affiliation := Affiliation{
					ID:   id,
					Name: a.Name,
				}
				affiliations = append(affiliations, affiliation)
			}
			if slices.Contains(v.ContributorRoles, "Author") {
				personOrOrg := PersonOrOrg{
					Name:        v.Name,
					GivenName:   v.GivenName,
					FamilyName:  v.FamilyName,
					Type:        v.Type + "al",
					Identifiers: identifiers,
				}
				contributor := Creator{
					PersonOrOrg:  personOrOrg,
					Affiliations: affiliations,
				}
				inveniordm.Metadata.Creators = append(inveniordm.Metadata.Creators, contributor)
			}
		}
	}

	// currently not using publisher
	// inveniordm.Metadata.Publisher = data.Publisher.Name
	// inveniordm.URL = data.URL

	// optional properties

	if data.Container.Title != "" {
		inveniordm.CustomFields.Journal.Title = data.Container.Title
	}
	if data.Container.Volume != "" {
		inveniordm.CustomFields.Journal.Volume = data.Container.Volume
	}
	if data.Container.Issue != "" {
		inveniordm.CustomFields.Journal.Issue = data.Container.Issue
	}
	if data.Container.FirstPage != "" {
		inveniordm.CustomFields.Journal.Pages = data.Container.Pages()
	}
	if data.Container.Identifier != "" && data.Container.IdentifierType == "ISSN" {
		inveniordm.CustomFields.Journal.ISSN = data.Container.Identifier
	}

	// optional custom fields
	inveniordm.CustomFields.ContentText = data.ContentText

	inveniordm.CustomFields.FeatureImage = data.FeatureImage

	if len(data.Identifiers) > 0 {
		for _, v := range data.Identifiers {
			if v.Identifier != data.ID {
				scheme := CMToInvenioIdentifierMappings[v.IdentifierType]
				identifier := Identifier{
					Identifier: v.Identifier,
					Scheme:     scheme,
				}
				inveniordm.Metadata.Identifiers = append(inveniordm.Metadata.Identifiers, identifier)
			}
		}
	}

	// add URL as identifier
	identifier := Identifier{
		Identifier: data.URL,
		Scheme:     "url",
	}
	inveniordm.Metadata.Identifiers = append(inveniordm.Metadata.Identifiers, identifier)

	if data.Date.Accepted != "" {
		inveniordm.Metadata.Dates = append(inveniordm.Metadata.Dates, Date{
			Date: data.Date.Accepted,
			Type: Type{
				ID: "accepted",
			},
		})
	}
	if data.Date.Available != "" {
		inveniordm.Metadata.Dates = append(inveniordm.Metadata.Dates, Date{
			Date: data.Date.Available,
			Type: Type{
				ID: "available",
			},
		})
	}
	if data.Date.Collected != "" {
		inveniordm.Metadata.Dates = append(inveniordm.Metadata.Dates, Date{
			Date: data.Date.Collected,
			Type: Type{
				ID: "collected",
			},
		})
	}
	if data.Date.Created != "" {
		inveniordm.Metadata.Dates = append(inveniordm.Metadata.Dates, Date{
			Date: data.Date.Created,
			Type: Type{
				ID: "created",
			},
		})
	}
	if data.Date.Published != "" {
		inveniordm.Metadata.Dates = append(inveniordm.Metadata.Dates, Date{
			Date: data.Date.Published,
			Type: Type{
				ID: "issued",
			},
		})
	}
	if data.Date.Accessed != "" {
		inveniordm.Metadata.Dates = append(inveniordm.Metadata.Dates, Date{
			Date: data.Date.Accessed,
			Type: Type{
				ID: "other",
			},
		})
	}
	if data.Date.Other != "" {
		inveniordm.Metadata.Dates = append(inveniordm.Metadata.Dates, Date{
			Date: data.Date.Other,
			Type: Type{
				ID: "other",
			},
		})
	}
	if data.Date.Submitted != "" {
		inveniordm.Metadata.Dates = append(inveniordm.Metadata.Dates, Date{
			Date: data.Date.Submitted,
			Type: Type{
				ID: "submitted",
			},
		})
	}
	if data.Date.Updated != "" {
		inveniordm.Metadata.Dates = append(inveniordm.Metadata.Dates, Date{
			Date: data.Date.Updated,
			Type: Type{
				ID: "updated",
			},
		})
	}
	if data.Date.Valid != "" {
		inveniordm.Metadata.Dates = append(inveniordm.Metadata.Dates, Date{
			Date: data.Date.Valid,
			Type: Type{
				ID: "valid",
			},
		})
	}
	if data.Date.Withdrawn != "" {
		inveniordm.Metadata.Dates = append(inveniordm.Metadata.Dates, Date{
			Date: data.Date.Withdrawn,
			Type: Type{
				ID: "withdrawn",
			},
		})
	}

	if len(data.Descriptions) > 0 {
		inveniordm.Metadata.Description = data.Descriptions[0].Description
	}

	if len(data.FundingReferences) > 0 {
		for _, v := range data.FundingReferences {
			id, identifierType := utils.ValidateID(v.FunderIdentifier)

			// convert Open Funder Registry DOI to ROR using mapping file
			if identifierType == "Crossref Funder ID" {
				id = crossrefxml.OFRToRORMappings[v.FunderIdentifier]
			}
			if id != "" {
				id, _ = utils.ValidateROR(id)
			}
			funder := Funder{
				ID:   id,
				Name: v.FunderName,
			}
			var award Award
			// first check if award number is in the awards vocabulary
			for _, a := range awardsVocabulary {
				if a.ID == v.AwardNumber {
					award = Award{
						Number: v.AwardNumber,
						Title: AwardTitle{
							En: a.Title.En,
						},
					}
					// if a.Identifiers.Identifier != "" {
					// 	id, identifierType := utils.ValidateID(a.AwardURI)
					// 	if id == "" {
					// 		id = a.AwardURI
					// 	}
					// 	identifier := Identifier{
					// 		Identifier: id,
					// 		Scheme:     strings.ToLower(identifierType),
					// 	}
					// 	award.Identifiers = append(award.Identifiers, identifier)
					// }
				}
			}
			if award.Number == "" && v.AwardNumber != "" {
				award = Award{
					Number: v.AwardNumber,
				}
				if v.AwardTitle != "" {
					award.Title = AwardTitle{
						En: v.AwardTitle,
					}
				}
			}
			if v.AwardURI != "" {
				id, identifierType := utils.ValidateID(v.AwardURI)
				if id == "" {
					id = v.AwardURI
				}
				identifier := Identifier{
					Identifier: id,
					Scheme:     strings.ToLower(identifierType),
				}
				award.Identifiers = append(award.Identifiers, identifier)
			}
			funding := Funding{
				Funder: funder,
				Award:  award,
			}
			inveniordm.Metadata.Funding = append(inveniordm.Metadata.Funding, funding)
		}
	}

	// if len(data.GeoLocations) > 0 {
	// 	for _, v := range data.GeoLocations {
	// 		geoLocation := GeoLocation{
	// 			GeoLocationPlace: v.GeoLocationPlace,
	// 			GeoLocationPoint: GeoLocationPoint{
	// 				PointLongitude: v.GeoLocationPoint.PointLongitude,
	// 				PointLatitude:  v.GeoLocationPoint.PointLatitude,
	// 			},
	// 			GeoLocationBox: GeoLocationBox{
	// 				WestBoundLongitude: v.GeoLocationBox.WestBoundLongitude,
	// 				EastBoundLongitude: v.GeoLocationBox.EastBoundLongitude,
	// 				SouthBoundLatitude: v.GeoLocationBox.SouthBoundLatitude,
	// 				NorthBoundLatitude: v.GeoLocationBox.NorthBoundLatitude,
	// 			},
	// 		}
	// 		inveniordm.GeoLocations = append(inveniordm.GeoLocations, geoLocation)
	// 	}
	// }

	if data.Language != "" {
		language := Language{
			ID: utils.GetLanguage(data.Language, "iso639-3"),
		}
		inveniordm.Metadata.Languages = append(inveniordm.Metadata.Languages, language)
	}
	if len(data.Subjects) > 0 {
		for _, v := range data.Subjects {
			ID := FOSMappings[v.Subject]
			var scheme string
			if ID != "" {
				scheme = "FOS"
			}
			subject := Subject{Subject: v.Subject, ID: ID, Scheme: scheme}
			inveniordm.Metadata.Subjects = append(inveniordm.Metadata.Subjects, subject)
		}
	}
	if data.License.URL != "" {
		right := Right{
			ID: data.License.ID,
		}
		inveniordm.Metadata.Rights = append(inveniordm.Metadata.Rights, right)
	}
	if len(data.Relations) > 0 {
		for _, v := range data.Relations {
			id, identifierType := utils.ValidateID(v.ID)
			scheme := CMToInvenioIdentifierMappings[identifierType]
			relationType := Type{ID: strings.ToLower(v.Type)}
			RelatedIdentifier := RelatedIdentifier{
				Identifier:   id,
				Scheme:       scheme,
				RelationType: relationType,
			}
			inveniordm.Metadata.RelatedIdentifiers = append(inveniordm.Metadata.RelatedIdentifiers, RelatedIdentifier)
		}
	}
	if len(data.References) > 0 {
		for _, v := range data.References {
			id, identifierType := utils.ValidateID(v.ID)
			scheme := CMToInvenioIdentifierMappings[identifierType]
			relationType := Type{ID: "references"}
			RelatedIdentifier := RelatedIdentifier{
				Identifier:   id,
				Scheme:       scheme,
				RelationType: relationType,
			}
			inveniordm.Metadata.RelatedIdentifiers = append(inveniordm.Metadata.RelatedIdentifiers, RelatedIdentifier)
		}
	}

	if len(data.Relations) > 0 {
		for _, v := range data.Relations {
			// skip ISSN IsPartOf relations as they are already included in the Journal metadata
			if v.Type == "IsPartOf" {
				continue
			}
			id, _ := doiutils.ValidateDOI(v.ID)
			RelatedIdentifier := RelatedIdentifier{
				Identifier:   id,
				Scheme:       "doi",
				RelationType: Type{ID: strings.ToLower(v.Type)},
			}
			inveniordm.Metadata.RelatedIdentifiers = append(inveniordm.Metadata.RelatedIdentifiers, RelatedIdentifier)
		}
	}

	inveniordm.Metadata.Version = data.Version

	return inveniordm, nil
}

// Write writes inveniordm metadata.
func Write(data commonmeta.Data) ([]byte, []gojsonschema.ResultError) {
	inveniordm, err := Convert(data)
	if err != nil {
		fmt.Println(err)
	}
	output, err := json.Marshal(inveniordm)
	if err != nil {
		fmt.Println(err)
	}
	// validation := schemautils.JSONSchemaErrors(output, "datacite-v4.5")
	// if !validation.Valid() {
	// 	return nil, validation.Errors()
	// }

	return output, nil
}

// WriteAll writes a list of inveniordm metadata.
func WriteAll(list []commonmeta.Data) ([]byte, []gojsonschema.ResultError) {
	var inveniordmList []Inveniordm
	for _, data := range list {
		inveniordm, err := Convert(data)
		if err != nil {
			fmt.Println(err)
		}
		inveniordmList = append(inveniordmList, inveniordm)
	}
	output, err := json.Marshal(inveniordmList)
	if err != nil {
		fmt.Println(err)
	}
	validation := schemautils.JSONSchemaErrors(output, "datacite-v4.5")
	if !validation.Valid() {
		return nil, validation.Errors()
	}

	return output, nil
}
