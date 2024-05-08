package datacite

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"

	"github.com/front-matter/commonmeta/bibtex"
	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/csl"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/ris"
	"github.com/front-matter/commonmeta/schemaorg"
	"github.com/front-matter/commonmeta/schemautils"
	"github.com/xeipuuv/gojsonschema"
)

// Convert converts Commonmeta metadata to DataCite metadata
func Convert(data commonmeta.Data) (Datacite, error) {
	var datacite Datacite

	// required properties
	datacite.ID = data.ID
	datacite.DOI, _ = doiutils.ValidateDOI(data.ID)
	datacite.Types.ResourceTypeGeneral = CMToDCMappings[data.Type]
	datacite.Types.SchemaOrg = schemaorg.CMToSOMappings[data.Type]
	datacite.Types.Citeproc = csl.CMToCSLMappings[data.Type]
	datacite.Types.Bibtex = bibtex.CMToBibMappings[data.Type]
	datacite.Types.Ris = ris.CMToRISMappings[data.Type]
	if data.AdditionalType != "" {
		datacite.Types.ResourceType = data.AdditionalType
	}
	if datacite.Types.ResourceTypeGeneral == "" {
		datacite.Types.ResourceTypeGeneral = "Other"
	}

	if len(data.Date.Published) >= 4 {
		datacite.PublicationYear, _ = strconv.Atoi(data.Date.Published[:4])
	}

	if len(data.Titles) > 0 {
		for _, v := range data.Titles {
			title := Title{
				Title:     v.Title,
				TitleType: v.Type,
				Lang:      v.Language,
			}
			datacite.Titles = append(datacite.Titles, title)
		}
	}

	if len(data.Contributors) > 0 {
		for _, v := range data.Contributors {
			var nameIdentifiers []NameIdentifier
			if v.ID != "" {
				nameIdentifier := NameIdentifier{
					NameIdentifier:       v.ID,
					NameIdentifierScheme: "ORCID",
					SchemeURI:            "https://orcid.org",
				}
				nameIdentifiers = append(nameIdentifiers, nameIdentifier)
			}
			var affiliations []string
			for _, a := range v.Affiliations {
				affiliation := a.Name
				affiliations = append(affiliations, affiliation)
			}
			if slices.Contains(v.ContributorRoles, "Author") {
				contributor := Contributor{
					Name:            v.Name,
					GivenName:       v.GivenName,
					FamilyName:      v.FamilyName,
					NameType:        v.Type + "al",
					NameIdentifiers: nameIdentifiers,
					Affiliation:     affiliations,
				}
				datacite.Creators = append(datacite.Creators, contributor)
			} else {
				contributorType := v.ContributorRoles[0]
				contributor := Contributor{
					Name:            v.Name,
					GivenName:       v.GivenName,
					FamilyName:      v.FamilyName,
					NameType:        v.Type + "al",
					NameIdentifiers: nameIdentifiers,
					Affiliation:     affiliations,
					ContributorType: contributorType,
				}
				datacite.Contributors = append(datacite.Contributors, contributor)
			}
		}
	}

	datacite.Publisher = Publisher{
		Name: data.Publisher.Name,
	}
	datacite.URL = data.URL
	datacite.SchemaVersion = "http://datacite.org/schema/kernel-4"

	// optional properties

	datacite.Container = Container{
		Type:           data.Container.Type,
		Identifier:     data.Container.Identifier,
		IdentifierType: data.Container.IdentifierType,
		Title:          data.Container.Title,
		Volume:         data.Container.Volume,
		Issue:          data.Container.Issue,
		FirstPage:      data.Container.FirstPage,
		LastPage:       data.Container.LastPage,
	}

	if len(data.Identifiers) > 0 {
		for _, v := range data.Identifiers {
			if v.Identifier != data.ID {
				AlternateIdentifier := AlternateIdentifier{
					AlternateIdentifier:     v.Identifier,
					AlternateIdentifierType: v.IdentifierType,
				}
				datacite.AlternateIdentifiers = append(datacite.AlternateIdentifiers, AlternateIdentifier)
			}
		}
	}

	if data.Date.Created != "" {
		datacite.Dates = append(datacite.Dates, Date{
			Date:     data.Date.Created,
			DateType: "Created",
		})
	} else if data.Date.Submitted != "" {
		datacite.Dates = append(datacite.Dates, Date{
			Date:     data.Date.Submitted,
			DateType: "Submitted",
		})
	} else if data.Date.Accepted != "" {
		datacite.Dates = append(datacite.Dates, Date{
			Date:     data.Date.Accepted,
			DateType: "Accepted",
		})
	} else if data.Date.Published != "" {
		datacite.Dates = append(datacite.Dates, Date{
			Date:     data.Date.Published,
			DateType: "Issued",
		})
	} else if data.Date.Updated != "" {
		datacite.Dates = append(datacite.Dates, Date{
			Date:     data.Date.Updated,
			DateType: "Updated",
		})
	} else if data.Date.Accessed != "" {
		datacite.Dates = append(datacite.Dates, Date{
			Date:     data.Date.Accessed,
			DateType: "Accessed",
		})
	} else if data.Date.Available != "" {
		datacite.Dates = append(datacite.Dates, Date{
			Date:     data.Date.Available,
			DateType: "Available",
		})
	} else if data.Date.Collected != "" {
		datacite.Dates = append(datacite.Dates, Date{
			Date:     data.Date.Collected,
			DateType: "Collected",
		})
	} else if data.Date.Valid != "" {
		datacite.Dates = append(datacite.Dates, Date{
			Date:     data.Date.Valid,
			DateType: "Valid",
		})
	} else if data.Date.Withdrawn != "" {
		datacite.Dates = append(datacite.Dates, Date{
			Date:     data.Date.Withdrawn,
			DateType: "Withdrawn",
		})
	} else if data.Date.Other != "" {
		datacite.Dates = append(datacite.Dates, Date{
			Date:     data.Date.Other,
			DateType: "Other",
		})
	}

	if len(data.Descriptions) > 0 {
		for _, v := range data.Descriptions {
			description := Description{
				Description:     v.Description,
				DescriptionType: v.Type,
				Lang:            v.Language,
			}
			datacite.Descriptions = append(datacite.Descriptions, description)
		}
	}

	if len(data.FundingReferences) > 0 {
		for _, v := range data.FundingReferences {
			fundingReference := FundingReference{
				FunderName:           v.FunderName,
				FunderIdentifier:     v.FunderIdentifier,
				FunderIdentifierType: v.FunderIdentifierType,
				AwardNumber:          v.AwardNumber,
				AwardURI:             v.AwardURI,
			}
			datacite.FundingReferences = append(datacite.FundingReferences, fundingReference)
		}
	}
	if len(data.GeoLocations) > 0 {
		for _, v := range data.GeoLocations {
			geoLocation := GeoLocation{
				GeoLocationPlace: v.GeoLocationPlace,
				GeoLocationPoint: GeoLocationPoint{
					PointLongitude: v.GeoLocationPoint.PointLongitude,
					PointLatitude:  v.GeoLocationPoint.PointLatitude,
				},
				GeoLocationBox: GeoLocationBox{
					WestBoundLongitude: v.GeoLocationBox.WestBoundLongitude,
					EastBoundLongitude: v.GeoLocationBox.EastBoundLongitude,
					SouthBoundLatitude: v.GeoLocationBox.SouthBoundLatitude,
					NorthBoundLatitude: v.GeoLocationBox.NorthBoundLatitude,
				},
			}
			datacite.GeoLocations = append(datacite.GeoLocations, geoLocation)
		}
	}
	datacite.Language = data.Language
	if len(data.Subjects) > 0 {
		for _, v := range data.Subjects {
			subject := Subject{Subject: v.Subject}
			datacite.Subjects = append(datacite.Subjects, subject)
		}
	}
	if data.License.URL != "" {
		rights := Rights{
			RightsURI:              data.License.URL,
			RightsIdentifier:       data.License.ID,
			RightsIdentifierScheme: "SPDX",
			SchemeURI:              "https://spdx.org/licenses/",
		}
		datacite.RightsList = append(datacite.RightsList, rights)
	}
	if len(data.Relations) > 0 {
		for _, v := range data.Relations {
			id := doiutils.NormalizeDOI(v.ID)
			relatedIdentifierType := "DOI"
			if id == "" {
				relatedIdentifierType = "URL"
			}
			RelatedIdentifier := RelatedIdentifier{
				RelatedIdentifier:     id,
				RelatedIdentifierType: relatedIdentifierType,
				RelationType:          v.Type,
			}
			datacite.RelatedIdentifiers = append(datacite.RelatedIdentifiers, RelatedIdentifier)
		}
	}

	if len(data.References) > 0 {
		for _, v := range data.References {
			id := doiutils.NormalizeDOI(v.ID)
			relatedIdentifierType := "DOI"
			if id == "" {
				relatedIdentifierType = "URL"
			}
			RelatedIdentifier := RelatedIdentifier{
				RelatedIdentifier:     id,
				RelatedIdentifierType: relatedIdentifierType,
				RelationType:          "References",
			}
			datacite.RelatedIdentifiers = append(datacite.RelatedIdentifiers, RelatedIdentifier)
		}
	}

	if len(data.Relations) > 0 {
		for _, v := range data.Relations {
			RelatedIdentifier := RelatedIdentifier{
				RelatedIdentifier:     v.ID,
				RelatedIdentifierType: "DOI",
				RelationType:          v.Type,
			}
			datacite.RelatedIdentifiers = append(datacite.RelatedIdentifiers, RelatedIdentifier)
		}
	}

	datacite.Version = data.Version

	return datacite, nil
}

// Write writes commonmeta metadata.
func Write(data commonmeta.Data) ([]byte, []gojsonschema.ResultError) {
	datacite, err := Convert(data)
	if err != nil {
		fmt.Println(err)
	}
	output, err := json.Marshal(datacite)
	if err != nil {
		fmt.Println(err)
	}
	validation := schemautils.JSONSchemaErrors(output, "datacite-v4.5")
	if !validation.Valid() {
		return nil, validation.Errors()
	}

	return output, nil
}

// WriteAll writes a list of commonmeta metadata.
func WriteAll(list []commonmeta.Data) ([]byte, []gojsonschema.ResultError) {
	var dataciteList []Datacite
	for _, data := range list {
		datacite, err := Convert(data)
		if err != nil {
			fmt.Println(err)
		}
		dataciteList = append(dataciteList, datacite)
	}
	output, err := json.Marshal(dataciteList)
	if err != nil {
		fmt.Println(err)
	}
	validation := schemautils.JSONSchemaErrors(output, "datacite-v4.5")
	if !validation.Valid() {
		return nil, validation.Errors()
	}

	return output, nil
}
