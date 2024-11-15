package inveniordm

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/crossrefxml"
	"github.com/front-matter/commonmeta/dateutils"
	"github.com/front-matter/commonmeta/doiutils"
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

	resourceType := CMToInvenioMappings[data.Type]
	if resourceType == "" {
		resourceType = "other"
	}
	inveniordm.Metadata.ResourceType = ResourceType{
		ID: resourceType,
	}
	if len(data.Titles) > 0 {
		inveniordm.Metadata.Title = data.Titles[0].Title
	}
	if len(data.Date.Published) >= 4 {
		inveniordm.Metadata.PublicationDate = dateutils.ParseDate(data.Date.Published)
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
					Type:        strings.ToLower(v.Type + "al"),
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
			scheme := CMToInvenioIdentifierMappings[v.IdentifierType]
			if v.Identifier != data.ID && scheme != "" {
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

	// using JSON to iterate over data.Date struct
	var dates map[string]interface{}
	d, _ := json.Marshal(data.Date)
	json.Unmarshal(d, &dates)
	for t, d := range dates {
		if d != "" {
			date := fmt.Sprintf("%v", d)
			id := strings.ToLower(t)
			if id == "published" {
				id = "issued"
			} else if id == "accessed" {
				id = "other"
			}
			inveniordm.Metadata.Dates = append(inveniordm.Metadata.Dates, Date{
				Date: date,
				Type: Type{
					ID: id,
				},
			})
		}
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
					if a.Title.En != "" {
						award = Award{
							Number: v.AwardNumber,
							// Title: AwardTitle{
							// 	En: a.Title.En,
							// },
						}
					} else {
						award = Award{
							Number: v.AwardNumber,
						}
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
				// if v.AwardTitle != "" {
				// 	award.Title = AwardTitle{
				// 		En: v.AwardTitle,
				// 	}
				// }
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
			// ID := FOSMappings[v.Subject]
			ID := ""
			var scheme string
			if ID != "" {
				scheme = "FOS"
			}
			subject := Subject{Subject: v.Subject, ID: ID, Scheme: scheme}
			inveniordm.Metadata.Subjects = append(inveniordm.Metadata.Subjects, subject)
		}
	}
	var right Right
	if data.License.ID != "" {
		right = Right{
			ID: strings.ToLower(data.License.ID),
		}
	} else if data.License.URL != "" {
		right = Right{
			ID: utils.URLToSPDX(data.License.URL),
		}
	}
	if right != (Right{}) {
		inveniordm.Metadata.Rights = append(inveniordm.Metadata.Rights, right)
	}

	if len(data.Relations) > 0 {
		for _, v := range data.Relations {
			// skip IsPartOf relation with InvendioRDM community identifier
			if v.Type == "IsPartOf" && strings.HasPrefix(v.ID, "https://rogue-scholar.org/api/communities/") {
				continue
			}
			id, identifierType := utils.ValidateID(v.ID)
			scheme := CMToInvenioIdentifierMappings[identifierType]
			relationType := CMToInvenioRelationTypeMappings[v.Type]
			if id != "" && scheme != "" && relationType != "" {
				RelatedIdentifier := RelatedIdentifier{
					Identifier:   id,
					Scheme:       scheme,
					RelationType: Type{ID: relationType},
				}
				inveniordm.Metadata.RelatedIdentifiers = append(inveniordm.Metadata.RelatedIdentifiers, RelatedIdentifier)
			}
		}
	}
	if len(data.References) > 0 {
		for _, v := range data.References {
			id, identifierType := utils.ValidateID(v.ID)
			scheme := CMToInvenioIdentifierMappings[identifierType]
			relationType := Type{ID: "references"}
			if id != "" && scheme != "" {
				RelatedIdentifier := RelatedIdentifier{
					Identifier:   id,
					Scheme:       scheme,
					RelationType: relationType,
				}
				inveniordm.Metadata.RelatedIdentifiers = append(inveniordm.Metadata.RelatedIdentifiers, RelatedIdentifier)
			}
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
	// validation := schemautils.JSONSchemaErrors(output, "datacite-v4.5")
	// if !validation.Valid() {
	// 	return nil, validation.Errors()
	// }

	return output, nil
}

// Post a list of inveniordm metadata to an InvenioRDM instance.
func PostAll(list []commonmeta.Data, host string, apiKey string) ([]byte, error) {
	type PostResponse struct {
		ID        string `json:"id"`
		DOI       string `json:"doi"`
		Community string `json:"community,omitempty"`
	}
	var postList []PostResponse
	for _, data := range list {
		inveniordm, err := Convert(data)
		if err != nil {
			return nil, err
		}

		// remove IsPartOf relation with InvendioRDM community identifier after storing it
		var communitySlug string
		var communityIndex int
		for i, v := range data.Relations {
			if v.Type == "IsPartOf" && strings.HasPrefix(v.ID, "https://rogue-scholar.org/api/communities/") {
				communitySlug = strings.Split(v.ID, "/")[5]
				communityIndex = i
			}
			if communityIndex != 0 {
				data.Relations = slices.Delete(data.Relations, i, i)
			}
		}

		// workaround until JSON schema validation is implemented
		// check for required fields
		if inveniordm.Metadata.Title == "" {
			// fmt.Println("Title is required: ", data.ID)
			continue
		}
		if inveniordm.Metadata.ResourceType.ID == "" {
			// fmt.Println("ResourceType is required: ", data.ID)
			continue
		}
		if inveniordm.Metadata.PublicationDate == "" {
			// fmt.Println("PublicationDate is required: ", data.ID)
			continue
		}
		if len(inveniordm.Metadata.Creators) == 0 {
			// fmt.Println("Creators is required: ", data.ID)
			continue
		}

		output, err := json.Marshal(inveniordm)
		if err != nil {
			return nil, err
		}

		// create draft record
		requestURL := fmt.Sprintf("https://%s/api/records", host)
		req, _ := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(output))
		req.Header = http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer " + apiKey},
		}
		client := &http.Client{
			Timeout: time.Second * 10,
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != 201 {
			fmt.Println(data.ID)
			return body, err
		}

		// publish draft record
		type Draft struct {
			ID string `json:"id"`
		}
		var draft Draft
		err = json.Unmarshal(body, &draft)
		if err != nil {
			return nil, err
		}
		requestURL = fmt.Sprintf("https://%s/api/records/%s/draft/actions/publish", host, draft.ID)
		req, _ = http.NewRequest(http.MethodPost, requestURL, nil)
		req.Header = http.Header{
			"Content-Type":  {"application/json"},
			"Authorization": {"Bearer " + apiKey},
		}
		client = &http.Client{
			Timeout: time.Second * 10,
		}
		resp, err = client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		record, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != 202 {
			fmt.Println(data.ID)
			return record, err
		}

		// optionally add record to community
		if communitySlug != "" {
			//get community ID from community slug
			requestURL := fmt.Sprintf("https://%s/api/communities/%s", host, communitySlug)
			req, _ = http.NewRequest(http.MethodGet, requestURL, nil)
			req.Header = http.Header{
				"Content-Type":  {"application/json"},
				"Authorization": {"Bearer " + apiKey},
			}
			client = &http.Client{
				Timeout: time.Second * 10,
			}
			resp, err = client.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			body, _ = io.ReadAll(resp.Body)
			if resp.StatusCode == 404 {
				continue // skip if community does not exist
			} else if resp.StatusCode >= 400 {
				return body, err
			}

			type Community struct {
				ID   string `json:"id"`
				Slug string `json:"slug,omitempty"`
			}
			var community Community
			err = json.Unmarshal(body, &community)
			if err != nil {
				return nil, err
			}
			type Communities struct {
				Communities []Community `json:"communities"`
			}
			com := Community{ID: community.ID}
			var communities Communities
			communities.Communities = append(communities.Communities, com)
			c, _ := json.Marshal(communities)
			requestURL = fmt.Sprintf("https://%s/api/records/%s/communities", host, draft.ID)
			req, _ = http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(c))
			req.Header = http.Header{
				"Content-Type":  {"application/json"},
				"Authorization": {"Bearer " + apiKey},
			}
			client = &http.Client{
				Timeout: time.Second * 10,
			}
			resp, err = client.Do(req)
			defer resp.Body.Close()
			body, _ = io.ReadAll(resp.Body)

			if resp.StatusCode >= 400 {
				return body, err
			}
		}
		post := PostResponse{
			ID:        draft.ID,
			DOI:       data.ID,
			Community: communitySlug,
		}
		postList = append(postList, post)
	}
	output, err := json.Marshal(postList)
	if err != nil {
		fmt.Println(err)
	}
	return output, nil
}
