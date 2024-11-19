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
	"github.com/front-matter/commonmeta/roguescholar"
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
			ID := FOSMappings[v.Subject]
			// ID := ""
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
			id, identifierType := utils.ValidateID(v.ID)
			// skip IsPartOf relation with InvendioRDM community identifier
			// skip IsPartOf relation with ISSN as that is already captured in the container
			if v.Type == "IsPartOf" && (strings.HasPrefix(v.ID, "https://rogue-scholar.org/api/communities/") || identifierType == "ISSN") {
				continue
			}
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
	validation := schemautils.JSONSchemaErrors(output, "invenio-rdm-v0.1")
	if !validation.Valid() {
		return nil, validation.Errors()
	}

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
	validation := schemautils.JSONSchemaErrors(output, "invenio-rdm-v0.1")
	if !validation.Valid() {
		return nil, validation.Errors()
	}

	return output, nil
}

// Upsert updates or creates a record in InvenioRDM.
func Upsert(record commonmeta.APIResponse, host string, apiKey string, legacyKey string, data commonmeta.Data) (commonmeta.APIResponse, error) {
	inveniordm, err := Convert(data)
	if err != nil {
		return record, err
	}

	record.DOI = data.ID

	// remove IsPartOf relation with InvendioRDM community identifier after storing it
	var communityIndex int
	for i, v := range data.Relations {
		if v.Type == "IsPartOf" && strings.HasPrefix(v.ID, "https://rogue-scholar.org/api/communities/") {
			record.Community = strings.Split(v.ID, "/")[5]
			communityIndex = i
		}
		if communityIndex != 0 {
			data.Relations = slices.Delete(data.Relations, i, i)
		}
	}

	// remove InvenioRDM rid after storing it
	var RIDIndex int
	for i, v := range data.Identifiers {
		if v.IdentifierType == "RID" && v.Identifier != "" {
			record.ID = v.Identifier
			RIDIndex = i
		} else if v.IdentifierType == "UUID" && v.Identifier != "" {
			record.UUID = v.Identifier
		}
		if RIDIndex != 0 {
			data.Identifiers = slices.Delete(data.Identifiers, i, i)
		}
	}

	// check if record already exists in InvenioRDM
	record.ID, _ = SearchByDOI(data.ID, host)

	if record.ID != "" {
		// create draft record from published record
		record, err = EditPublishedRecord(record, host, apiKey)
		if err != nil {
			return record, err
		}

		// update draft record
		record, err = UpdateDraftRecord(record, host, apiKey, inveniordm)
		if err != nil {
			return record, err
		}
	} else {
		// create draft record
		record, err = CreateDraftRecord(record, host, apiKey, inveniordm)
		if err != nil {
			return record, err
		}
	}

	// publish draft record
	record, err = PublishDraftRecord(record, host, apiKey)
	if err != nil {
		return record, err
	}

	// add record to community if community is specified and exists
	if record.Community != "" {
		communityID, err := SearchBySlug(record.Community, host)
		if err != nil {
			return record, err
		}
		record, err = AddRecordToCommunity(record, host, apiKey, communityID)
		if err != nil {
			return record, err
		}
	}

	// update rogue-scholar legacy record with Invenio rid if host is rogue-scholar.org
	if host == "rogue-scholar.org" && legacyKey != "" {
		record, err = roguescholar.UpdateLegacyRecord(record, legacyKey, "rid")
		if err != nil {
			return record, err
		}
	}

	return record, nil
}

// UpsertAll updates or creates a list of records in InvenioRDM.
func UpsertAll(list []commonmeta.Data, host string, apiKey string, legacyKey string) ([]commonmeta.APIResponse, error) {
	var records []commonmeta.APIResponse

	// iterate over list of records
	for _, data := range list {
		record := commonmeta.APIResponse{DOI: data.ID}

		record, err := Upsert(record, host, apiKey, legacyKey, data)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}

	return records, nil
}

// CreateDraftRecord creates a draft record in InvenioRDM.
func CreateDraftRecord(record commonmeta.APIResponse, host string, apiKey string, inveniordm Inveniordm) (commonmeta.APIResponse, error) {
	output, err := json.Marshal(inveniordm)
	if err != nil {
		return record, err
	}

	type Response struct {
		*Inveniordm
		Created string `json:"created,omitempty"`
		Updated string `json:"updated,omitempty"`
	}
	var response Response

	var requestURL string
	var req *http.Request
	var resp *http.Response
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	requestURL = fmt.Sprintf("https://%s/api/records", host)
	req, _ = http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(output))
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + apiKey},
	}
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != 201 {
		return record, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return record, err
	}
	if response != (Response{}) {
		record.ID = response.ID
		record.Created = response.Created
		record.Updated = response.Updated
		record.Status = "draft"
	}
	return record, nil
}

// EditPublishedRecord creates a draft record from a published record in InvenioRDM.
func EditPublishedRecord(record commonmeta.APIResponse, host string, apiKey string) (commonmeta.APIResponse, error) {
	type Response struct {
		*Inveniordm
		Created string `json:"created,omitempty"`
		Updated string `json:"updated,omitempty"`
	}
	var response Response
	var requestURL string
	var req *http.Request
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	requestURL = fmt.Sprintf("https://%s/api/records/%s/draft", host, record.ID)
	req, _ = http.NewRequest(http.MethodPost, requestURL, nil)
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + apiKey},
	}
	resp, err := client.Do(req)
	if err != nil {
		return record, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return record, err
	}
	record.Updated = response.Updated
	return record, nil
}

// UpdateDraftRecord updates a draft record in InvenioRDM.
func UpdateDraftRecord(record commonmeta.APIResponse, host string, apiKey string, inveniordm Inveniordm) (commonmeta.APIResponse, error) {
	output, err := json.Marshal(inveniordm)
	if err != nil {
		return record, err
	}

	type Response struct {
		*Inveniordm
		Created string `json:"created,omitempty"`
		Updated string `json:"updated,omitempty"`
	}
	var response Response
	requestURL := fmt.Sprintf("https://%s/api/records/%s/draft", host, record.ID)
	req, _ := http.NewRequest(http.MethodPut, requestURL, bytes.NewReader(output))
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + apiKey},
	}
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := client.Do(req)
	if err != nil {
		return record, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return record, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return record, err
	}
	record.Updated = response.Updated
	return record, nil
}

// PublishDraftRecord publishes a draft record in InvenioRDM.
func PublishDraftRecord(record commonmeta.APIResponse, host string, apiKey string) (commonmeta.APIResponse, error) {
	type Response struct {
		*Inveniordm
		Created string `json:"created,omitempty"`
		Updated string `json:"updated,omitempty"`
		Status  string `json:"status,omitempty"`
	}
	var response Response
	requestURL := fmt.Sprintf("https://%s/api/records/%s/draft/actions/publish", host, record.ID)
	req, _ := http.NewRequest(http.MethodPost, requestURL, nil)
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + apiKey},
	}
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := client.Do(req)
	if err != nil {
		return record, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 202 {
		return record, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return record, err
	}
	record.Created = response.Created
	record.Updated = response.Updated
	record.Status = response.Status
	return record, nil
}

// CreateCommunity creates a community in InvenioRDM.
func CreateCommunity(community string, host string, apiKey string) (string, error) {
	type Response struct {
		ID string `json:"id"`
	}
	var response Response
	var requestURL string
	var req *http.Request
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	requestURL = fmt.Sprintf("https://%s/api/communities", host)
	req, _ = http.NewRequest(http.MethodPost, requestURL, nil)
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + apiKey},
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	return response.ID, nil
}

// AddRecordToCommunity adds record to InvenioRDM community.
func AddRecordToCommunity(record commonmeta.APIResponse, host string, apiKey string, communityID string) (commonmeta.APIResponse, error) {
	type Response struct {
		ID string `json:"id"`
	}
	var response Response
	client := &http.Client{
		Timeout: time.Second * 30,
	}

	var output = []byte(`{"communities":[{"id":"` + communityID + `"}]}`)
	requestURL := fmt.Sprintf("https://%s/api/records/%s/communities", host, record.ID)
	req, _ := http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(output))
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + apiKey},
	}
	resp, err := client.Do(req)
	if err != nil {
		return record, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &response)
	record.Status = "added_to_community"
	return record, err
}
