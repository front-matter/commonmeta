package inveniordm

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"slices"
	"strings"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/dateutils"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/roguescholar"
	"github.com/front-matter/commonmeta/ror"
	"github.com/front-matter/commonmeta/schemautils"
	"github.com/front-matter/commonmeta/utils"
	"gopkg.in/yaml.v3"
)

// Vocabularies is the embedded vocabulary yaml files.
//
//go:embed vocabularies/*.yaml
var Vocabularies embed.FS

// Convert converts Commonmeta metadata to InvenioRDM metadata
func Convert(data commonmeta.Data, fromHost string) (Inveniordm, error) {
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
	if inveniordm.Metadata.Title == "" {
		inveniordm.Metadata.Title = "No title"
	}
	if len(data.Date.Published) >= 4 {
		inveniordm.Metadata.PublicationDate = dateutils.ParseDate(data.Date.Published)
	} else if len(data.Date.Available) >= 4 {
		inveniordm.Metadata.PublicationDate = dateutils.ParseDate(data.Date.Available)
	} else if len(data.Date.Created) >= 4 {
		inveniordm.Metadata.PublicationDate = dateutils.ParseDate(data.Date.Created)
	}

	if len(data.Contributors) > 0 {
		for _, v := range data.Contributors {
			var identifiers []Identifier
			if v.ID != "" {
				id, ok := utils.ValidateORCID(v.ID)
				if ok {
					Identifier := Identifier{
						Identifier: id,
						Scheme:     "orcid",
					}
					identifiers = append(identifiers, Identifier)
				}
			}
			var affiliations []Affiliation
			for _, a := range v.Affiliations {
				id, _ := utils.ValidateROR(a.ID)
				affiliation := Affiliation{
					ID:   id,
					Name: a.Name,
				}
				// avoid duplicate affiliations
				if !slices.ContainsFunc(affiliations, func(e Affiliation) bool {
					return e.ID != "" && e.ID == affiliation.ID
				}) {
					affiliations = append(affiliations, affiliation)
				}
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
	} else {
		// add placeholder author if no authors are provided
		personOrOrg := PersonOrOrg{
			Name: "No author",
			Type: "organizational",
		}
		contributor := Creator{
			PersonOrOrg: personOrOrg,
		}
		inveniordm.Metadata.Creators = append(inveniordm.Metadata.Creators, contributor)
	}

	inveniordm.Metadata.Publisher = data.Publisher.Name

	// optional properties

	if data.Container.Title != "" {
		inveniordm.CustomFields.Journal.Title = data.Container.Title
	}
	if data.Container.Platform != "" {
		inveniordm.CustomFields.Generator = data.Container.Platform
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
	inveniordm.CustomFields.ContentHTML = data.ContentHTML
	inveniordm.CustomFields.FeatureImage = data.FeatureImage

	if len(data.Identifiers) > 0 {
		for _, v := range data.Identifiers {
			scheme := CMToInvenioIdentifierMappings[v.IdentifierType]
			if scheme == "" || (v.Identifier == data.ID && scheme == "doi") {
				continue
			}
			identifier := Identifier{
				Identifier: v.Identifier,
				Scheme:     scheme,
			}
			inveniordm.Metadata.Identifiers = append(inveniordm.Metadata.Identifiers, identifier)
		}
	}

	// add URL as identifier
	if data.URL != "" {
		identifier := Identifier{
			Identifier: data.URL,
			Scheme:     "url",
		}
		inveniordm.Metadata.Identifiers = append(inveniordm.Metadata.Identifiers, identifier)
	}

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
			// convert Open Funder Registry DOI to ROR
			if identifierType == "Crossref Funder ID" {
				r, _ := ror.Fetch(v.FunderIdentifier)
				id = r.ID
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
				} else {
					award.Title = AwardTitle{
						En: "No title",
					}
				}
				if v.AwardURI != "" {
					id, identifierType := utils.ValidateID(v.AwardURI)
					if id != "" && identifierType != "" {
						identifier := Identifier{
							Identifier: id,
							Scheme:     strings.ToLower(identifierType),
						}
						award.Identifiers = append(award.Identifiers, identifier)
					}
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
			ID := commonmeta.FOSMappings[v.Subject]
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
			unstructured := v.Unstructured
			if unstructured == "" {
				// use title as unstructured reference
				if v.Title != "" {
					unstructured = v.Title
				} else {
					unstructured = "Unknown title"
				}
				if v.PublicationYear != "" {
					unstructured += " (" + v.PublicationYear + ")."
				}
			} else {
				if v.ID != "" {
					// remove duplicate ID from unstructured reference
					unstructured = strings.Replace(unstructured, v.ID, "", 1)
				}
				// remove optional trailing period
				unstructured = strings.TrimSuffix(unstructured, " .")
			}
			reference := Reference{
				Reference:  unstructured,
				Scheme:     scheme,
				Identifier: id,
			}
			inveniordm.Metadata.References = append(inveniordm.Metadata.References, reference)
		}
	}

	if len(data.Relations) > 0 {
		for _, v := range data.Relations {
			id, identifierType := utils.ValidateID(v.ID)
			if id != "" && identifierType == "URL" {
				u, _ := url.Parse(id)
				if u.Host == fromHost {
					u.Host = fromHost
					id = u.String()
				}
			}
			// skip IsPartOf relation with ISSN as that is already captured in the container
			if identifierType == "ISSN" {
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
func Write(data commonmeta.Data, fromHost string) ([]byte, error) {
	inveniordm, err := Convert(data, fromHost)
	if err != nil {
		fmt.Println(err)
	}
	output, err := json.Marshal(inveniordm)
	if err != nil {
		fmt.Println(err)
	}
	err = schemautils.JSONSchemaErrors(output, "invenio-rdm-v0.1")
	return output, err
}

// WriteAll writes a list of inveniordm metadata.
func WriteAll(list []commonmeta.Data, fromHost string) ([]byte, error) {
	var inveniordmList []Inveniordm
	for _, data := range list {
		inveniordm, err := Convert(data, fromHost)
		if err != nil {
			fmt.Println(err)
		}
		inveniordmList = append(inveniordmList, inveniordm)
	}
	output, err := json.Marshal(inveniordmList)
	if err != nil {
		fmt.Println(err)
	}
	err = schemautils.JSONSchemaErrors(output, "invenio-rdm-v0.1")
	return output, err
}

// Upsert updates or creates a record in InvenioRDM.
func Upsert(record commonmeta.APIResponse, fromHost string, apiKey string, legacyKey string, data commonmeta.Data, client *InveniordmClient) (commonmeta.APIResponse, error) {
	if client.Host == "rogue-scholar.org" && !doiutils.IsRogueScholarDOI(data.ID, "") {
		record.Status = "failed_not_rogue_scholar_doi"
		return record, nil
	}

	inveniordm, err := Convert(data, fromHost)
	if err != nil {
		return record, err
	}

	doi, ok := doiutils.ValidateDOI(data.ID)
	if !ok {
		record.Status = "failed"
		return record, fmt.Errorf("missing or invalid DOI")
	}

	record.DOI = doi

	// check if required metadata are present
	if inveniordm.Metadata.PublicationDate == "" {
		record.Status = "failed"
		return record, fmt.Errorf("missing publication date: %s", record.DOI)
	}

	// remove IsPartOf relation with InvendioRDM community identifier after storing it
	var communityIndex int
	for i, v := range data.Relations {
		if v.Type == "IsPartOf" && strings.HasPrefix(v.ID, fmt.Sprintf("https://%s/api/communities/", client.Host)) {
			slug := strings.Split(v.ID, "/")[5]
			communityID, _ := SearchBySlug(slug, "blog", client)
			if communityID != "" {
				record.Community = slug
				record.CommunityID = communityID
				communityIndex = i
			}
		}
		if communityIndex != 0 {
			data.Relations = slices.Delete(data.Relations, i, i)
		}
	}

	// remove InvenioRDM rid after storing it
	for i := len(data.Identifiers) - 1; i >= 0; i-- {
		v := data.Identifiers[i]
		if v.IdentifierType == "RID" && v.Identifier != "" {
			record.ID = v.Identifier
			data.Identifiers = slices.Delete(data.Identifiers, i, i+1)
		} else if v.IdentifierType == "UUID" && v.Identifier != "" {
			record.UUID = v.Identifier
		}
	}
	// check if record already exists in InvenioRDM
	record.ID, _ = SearchByDOI(data.ID, client)
	if record.ID == "" {
		// create draft record
		record, err = CreateDraftRecord(record, apiKey, inveniordm, client)
		if err != nil {
			return record, err
		}
	} else {
		// create draft record from published record
		record, err = EditPublishedRecord(record, apiKey, client)
		if err != nil {
			return record, err
		}
		// update draft record
		record, err = UpdateDraftRecord(record, apiKey, inveniordm, client)
		if err != nil {
			return record, err
		}
	}

	// publish draft record
	record, err = PublishDraftRecord(record, apiKey, client)
	if err != nil {
		return record, err
	}

	// add record to blog community if blog community is specified and exists
	if record.CommunityID != "" {
		record, err = AddRecordToCommunity(record, client, apiKey, record.CommunityID)
		if err != nil {
			return record, err
		}
	}

	// add record to subject area community if subject area community is specified and exists
	// subject area communities should exist for all subjects in the FOSMappings
	if len(data.Subjects) > 0 {
		var slug, communityID string
		for _, v := range data.Subjects {
			slug = utils.StringToSlug(v.Subject)
			if synonym := CommunityTranslations[slug]; synonym != "" {
				slug = synonym
			}
			communityID, err = SearchBySlug(slug, "subject", client)
			if err != nil {
				fmt.Println(err)
			}
			if communityID != "" {
				record, err = AddRecordToCommunity(record, client, apiKey, communityID)
				if err != nil {
					return record, err
				}
			}
		}
	}

	// add record to communities defined as IsPartOf relation in inveniordm.Metadata.RelatedIdentifiers
	if len(inveniordm.Metadata.RelatedIdentifiers) > 0 {
		for i := len(inveniordm.Metadata.RelatedIdentifiers) - 1; i >= 0; i-- {
			if inveniordm.Metadata.RelatedIdentifiers[i].RelationType.ID == "ispartof" {
				u, _ := url.Parse(inveniordm.Metadata.RelatedIdentifiers[i].Identifier)
				c := strings.Split(u.Path, "/")
				if u.Host == fromHost && len(c) == 4 && c[2] == "communities" {
					record, err = AddRecordToCommunity(record, client, apiKey, c[3])
					if err != nil {
						return record, err
					}
					// remove related identifier
					fmt.Println("Removing related identifier:", inveniordm.Metadata.RelatedIdentifiers[i].Identifier)
					inveniordm.Metadata.RelatedIdentifiers = slices.Delete(inveniordm.Metadata.RelatedIdentifiers, i, i+1)
				}
			}
		}
	}
	fmt.Println(len(inveniordm.Metadata.RelatedIdentifiers), "related identifiers after removing IsPartOf relations")

	// update rogue-scholar legacy record with Invenio rid if host is rogue-scholar.org
	if client.Host == "rogue-scholar.org" && legacyKey != "" {
		record, err = roguescholar.UpdateLegacyRecord(record, legacyKey, "rid")
		if err != nil {
			return record, err
		}
	}
	return record, nil
}

// UpsertAll updates or creates a list of records in InvenioRDM.
func UpsertAll(list []commonmeta.Data, fromHost string, apiKey string, legacyKey string, client *InveniordmClient) ([]commonmeta.APIResponse, error) {
	var records []commonmeta.APIResponse

	for _, data := range list {
		record := commonmeta.APIResponse{ID: data.ID}
		doi, ok := doiutils.ValidateDOI(data.ID)
		if !ok && doi == "" {
			record.Status = "failed_missing_doi"
		} else if client.Host == "rogue-scholar.org" && !doiutils.IsRogueScholarDOI(data.ID, "") {
			record.Status = "failed_not_rogue_scholar_doi"
		} else {
			record, _ = Upsert(record, fromHost, apiKey, legacyKey, data, client)
		}

		records = append(records, record)
	}

	return records, nil
}

// CreateDraftRecord creates a draft record in InvenioRDM.
func CreateDraftRecord(record commonmeta.APIResponse, apiKey string, inveniordm Inveniordm, client *InveniordmClient) (commonmeta.APIResponse, error) {
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
	requestURL = fmt.Sprintf("https://%s/api/records", client.Host)
	req, _ = http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(output))
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"Authorization": {"Bearer " + apiKey},
	}
	resp, err = client.Do(req)
	if err != nil {
		return record, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 429 {
		record.Status = "failed_rate_limited"
		return record, fmt.Errorf("rate limited")
	}
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 201 {
		fmt.Println(string(body), record)
		record.Status = "failed_create_draft"
		return record, errors.New("failed to create draft record:" + string(body))
	}
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
func EditPublishedRecord(record commonmeta.APIResponse, apiKey string, client *InveniordmClient) (commonmeta.APIResponse, error) {
	type Response struct {
		*Inveniordm
		Created string `json:"created,omitempty"`
		Updated string `json:"updated,omitempty"`
	}
	var response Response
	var requestURL string
	var req *http.Request
	requestURL = fmt.Sprintf("https://%s/api/records/%s/draft", client.Host, record.ID)
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
func UpdateDraftRecord(record commonmeta.APIResponse, apiKey string, inveniordm Inveniordm, client *InveniordmClient) (commonmeta.APIResponse, error) {
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
	requestURL := fmt.Sprintf("https://%s/api/records/%s/draft", client.Host, record.ID)
	req, _ := http.NewRequest(http.MethodPut, requestURL, bytes.NewReader(output))
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
	err = json.Unmarshal(body, &response)
	if err != nil {
		return record, err
	}
	record.Updated = response.Updated
	return record, nil
}

// PublishDraftRecord publishes a draft record in InvenioRDM.
func PublishDraftRecord(record commonmeta.APIResponse, apiKey string, client *InveniordmClient) (commonmeta.APIResponse, error) {
	type Response struct {
		*Inveniordm
		Created string `json:"created,omitempty"`
		Updated string `json:"updated,omitempty"`
		Status  string `json:"status,omitempty"`
	}
	var response Response
	requestURL := fmt.Sprintf("https://%s/api/records/%s/draft/actions/publish", client.Host, record.ID)
	req, _ := http.NewRequest(http.MethodPost, requestURL, nil)
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
	if resp.StatusCode != 202 {
		fmt.Println(string(body), record)
		return record, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return record, err
	}
	record.Created = response.Created
	record.Updated = response.Updated
	record.Status = "published"
	return record, nil
}

// DeleteDraftRecord publishes a draft record in InvenioRDM.
func DeleteDraftRecord(record commonmeta.APIResponse, apiKey string, client *InveniordmClient) (commonmeta.APIResponse, error) {
	requestURL := fmt.Sprintf("https://%s/api/records/%s/draft", client.Host, record.ID)
	req, _ := http.NewRequest(http.MethodDelete, requestURL, nil)
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
	if resp.StatusCode != 204 {
		fmt.Println(string(body), record)
		return record, err
	}
	record.Status = "deleted"
	return record, nil
}

// CreateSubjectCommunities creates communities for each subject in FOSKeyMappings
func CreateSubjectCommunities(apiKey string, client *InveniordmClient) ([]byte, error) {
	var communities []Community
	var id string
	var err error
	var output []byte

	for slug := range commonmeta.FOSKeyMappings {
		title := commonmeta.FOSKeyMappings[slug]
		community := Community{
			Access: &CommunityAccess{
				Visibility:   "public",
				MemberPolicy: "open",
				RecordPolicy: "open",
				ReviewPolicy: "open",
			},
			Slug: utils.StringToSlug(slug),
			Metadata: CommunityMetadata{
				Title:       title,
				Description: title + " subject area.",
				Type: Type{
					ID: "subject",
				},
			},
		}
		id, err = UpsertCommunity(community, apiKey, client)
		if err != nil {
			return output, err
		}
		community.ID = id
		communities = append(communities, community)
	}

	output, err = json.Marshal(communities)
	return output, err
}

// TransferCommunities transfers communities between InvenioRDM instances
// Transfer is my community type, e.g. blog, topic or subject
func TransferCommunities(type_ string, apiKey string, oldApiKey string, oldClient *InveniordmClient, client *InveniordmClient) ([]byte, error) {
	var oldCommunities, communities []Community
	var id string
	var err error
	var output []byte

	// get all communities by type from old InvenioRDM instance
	oldCommunities, err = SearchByType(type_, oldApiKey, oldClient)
	if err != nil {
		return output, err
	}

	for _, community := range oldCommunities {
		id, err = UpsertCommunity(community, apiKey, client)
		if err != nil {
			return output, err
		}
		community.ID = id
		communities = append(communities, community)

		// transfer optional logo if it exists
		logo, err := GetCommunityLogo(community.Slug, oldClient)
		if err != nil {
			return output, err
		}
		if len(logo) > 0 {
			_, err = UpdateCommunityLogo(community.Slug, logo, apiKey, client)
			if err != nil {
				return output, err
			}
		}
	}

	output, err = json.Marshal(communities)
	return output, err
}

// UpsertCommunity updates or creates a community in InvenioRDM.
func UpsertCommunity(community Community, apiKey string, client *InveniordmClient) (string, error) {
	var err error
	var communityID string

	// check if community already exists
	communityID, _ = SearchBySlug(community.Slug, community.Metadata.Type.ID, client)
	if communityID != "" {
		community.ID = communityID
		communityID, err = UpdateCommunity(community, client, apiKey)
		if err != nil {
			return communityID, err
		}
	} else {
		// Create the community if it doesn't exist
		communityID, err = CreateCommunity(community, client, apiKey)
		if err != nil {
			return communityID, err
		}
	}
	return communityID, nil
}

// UpdateCommunityLogo updates a community logo in InvenioRDM.
func UpdateCommunityLogo(slug string, logo []byte, apiKey string, client *InveniordmClient) (string, error) {
	if len(logo) == 0 {
		return "", fmt.Errorf("empty logo data")
	}

	requestURL := fmt.Sprintf("https://%s/api/communities/%s/logo", client.Host, slug)
	req, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewReader(logo))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header = http.Header{
		"Content-Type":  {"application/octet-stream"},
		"Authorization": {"Bearer " + apiKey},
	}

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to upload logo (status %d): %s", resp.StatusCode, string(body))
	}

	return slug, nil
}

// CreateCommunity creates a community in InvenioRDM.
func CreateCommunity(community Community, client *InveniordmClient, apiKey string) (string, error) {
	type Response struct {
		*Community
		Created string `json:"created,omitempty"`
		Updated string `json:"updated,omitempty"`
	}
	var response Response
	var requestURL string
	var req *http.Request

	jsonData, err := json.Marshal(community)
	if err != nil {
		return response.ID, err
	}
	requestURL = fmt.Sprintf("https://%s/api/communities", client.Host)
	req, _ = http.NewRequest(http.MethodPost, requestURL, bytes.NewReader(jsonData))
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
	fmt.Println(string(body))
	if resp.StatusCode != 201 {
		return "", fmt.Errorf("failed to create community: %s", string(body))
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	return response.ID, nil
}

// UpdateCommunity updates a community in InvenioRDM.
func UpdateCommunity(community Community, client *InveniordmClient, apiKey string) (string, error) {
	type Response struct {
		*Community
		Created string `json:"created,omitempty"`
		Updated string `json:"updated,omitempty"`
	}
	var response Response
	var requestURL string
	var req *http.Request

	jsonData, err := json.Marshal(community)
	if err != nil {
		return response.ID, err
	}
	requestURL = fmt.Sprintf("https://%s/api/communities/%s", client.Host, community.Slug)
	req, _ = http.NewRequest(http.MethodPut, requestURL, bytes.NewReader(jsonData))
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
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to update community: %s", string(body))
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}
	return response.ID, nil
}

// AddRecordToCommunity adds record to InvenioRDM community.
func AddRecordToCommunity(record commonmeta.APIResponse, client *InveniordmClient, apiKey string, communityID string) (commonmeta.APIResponse, error) {
	type Response struct {
		ID string `json:"id"`
	}
	var response Response
	var output = []byte(`{"communities":[{"id":"` + communityID + `"}]}`)
	requestURL := fmt.Sprintf("https://%s/api/records/%s/communities", client.Host, record.ID)
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
	record.Status = "added_to_community"
	return record, err
}
