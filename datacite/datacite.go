package datacite

import (
	"commonmeta/doiutils"
	"commonmeta/types"
	"commonmeta/utils"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
)

// source: https://github.com/datacite/schema/blob/master/source/meta/kernel-4/include/datacite-resourceType-v4.xsd
var DCToCMTranslations = map[string]string{
	"Audiovisual":           "Audiovisual",
	"BlogPosting":           "Article",
	"Book":                  "Book",
	"BookChapter":           "BookChapter",
	"Collection":            "Collection",
	"ComputationalNotebook": "ComputationalNotebook",
	"ConferencePaper":       "ProceedingsArticle",
	"ConferenceProceeding":  "Proceedings",
	"DataPaper":             "JournalArticle",
	"Dataset":               "Dataset",
	"Dissertation":          "Dissertation",
	"Event":                 "Event",
	"Image":                 "Image",
	"Instrument":            "Instrument",
	"InteractiveResource":   "InteractiveResource",
	"Journal":               "Journal",
	"JournalArticle":        "JournalArticle",
	"Model":                 "Model",
	"OutputManagementPlan":  "OutputManagementPlan",
	"PeerReview":            "PeerReview",
	"PhysicalObject":        "PhysicalObject",
	"Poster":                "Presentation",
	"Preprint":              "Article",
	"Report":                "Report",
	"Service":               "Service",
	"Software":              "Software",
	"Sound":                 "Sound",
	"Standard":              "Standard",
	"StudyRegistration":     "StudyRegistration",
	"Text":                  "Document",
	"Thesis":                "Dissertation",
	"Workflow":              "Workflow",
	"Other":                 "Other",
}

// from commonmeta schema
var CommonmetaContributorRoles = []string{
	"Author",
	"Editor",
	"Chair",
	"Reviewer",
	"ReviewAssistant",
	"StatsReviewer",
	"ReviewerExternal",
	"Reader",
	"Translator",
	"ContactPerson",
	"DataCollector",
	"DataManager",
	"Distributor",
	"HostingInstitution",
	"Producer",
	"ProjectLeader",
	"ProjectManager",
	"ProjectMember",
	"RegistrationAgency",
	"RegistrationAuthority",
	"RelatedPerson",
	"ResearchGroup",
	"RightsHolder",
	"Researcher",
	"Sponsor",
	"WorkPackageLeader",
	"Conceptualization",
	"DataCuration",
	"FormalAnalysis",
	"FundingAcquisition",
	"Investigation",
	"Methodology",
	"ProjectAdministration",
	"Resources",
	"Software",
	"Supervision",
	"Validation",
	"Visualization",
	"WritingOriginalDraft",
	"WritingReviewEditing",
	"Maintainer",
	"Other",
}

func FetchDatacite(str string) (types.Data, error) {
	var data types.Data
	id, ok := doiutils.ValidateDOI(str)
	if !ok {
		return data, errors.New("invalid doi")
	}
	content, err := GetDatacite(id)
	if err != nil {
		return data, err
	}
	data, err = ReadDatacite(content)
	if err != nil {
		return data, err
	}
	return data, nil
}

func GetDatacite(pid string) (types.Content, error) {
	// the envelope for the JSON response from the DataCite API
	type Response struct {
		Data types.Content `json:"data"`
	}

	var response Response
	doi, ok := doiutils.ValidateDOI(pid)
	if !ok {
		return response.Data, errors.New("Invalid DOI")
	}
	url := "https://api.datacite.org/dois/" + doi
	client := http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(url)
	if err != nil {
		return response.Data, err
	}
	if resp.StatusCode >= 400 {
		return response.Data, errors.New(resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return response.Data, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("error:", err)
	}
	return response.Data, err
}

// read DataCite JSON response and return work struct in Commonmeta format
func ReadDatacite(content types.Content) (types.Data, error) {
	var data = types.Data{}

	data.ID = doiutils.NormalizeDOI(content.Attributes.DOI)
	data.Type = DCToCMTranslations[content.Attributes.Types.ResourceTypeGeneral]
	var err error

	data.Identifiers = append(data.Identifiers, types.Identifier{
		Identifier:     data.ID,
		IdentifierType: "DOI",
	})
	if len(content.Attributes.AlternateIdentifiers) > 0 {
		for _, v := range content.Attributes.AlternateIdentifiers {
			if content.Attributes.AlternateIdentifiers[0].Identifier != "" {
				data.Identifiers = append(data.Identifiers, types.Identifier{
					Identifier:     v.Identifier,
					IdentifierType: v.IdentifierType,
				})
			}
		}
	}

	data.AdditionalType = DCToCMTranslations[content.Attributes.Types.ResourceType]
	if data.AdditionalType != "" {
		data.Type = data.AdditionalType
		data.AdditionalType = ""
	} else {
		data.AdditionalType = content.Attributes.Types.ResourceType
	}
	data.Url, err = utils.NormalizeUrl(content.Attributes.Url, true, false)
	if err != nil {
		log.Println(err)
	}

	if len(content.Attributes.Creators) > 0 {
		for _, v := range content.Attributes.Creators {
			if v.Name != "" || v.GivenName != "" || v.FamilyName != "" {
				contributor := GetContributor(v)
				containsID := slices.ContainsFunc(data.Contributors, func(e types.Contributor) bool {
					return e.ID != "" && e.ID == contributor.ID
				})
				if containsID {
					log.Printf("Contributor with ID %s already exists", contributor.ID)
				} else {
					data.Contributors = append(data.Contributors, contributor)

				}
			}
		}

		// merge creators and contributors
		for _, v := range content.Attributes.Contributors {
			if v.Name != "" || v.GivenName != "" || v.FamilyName != "" {
				contributor := GetContributor(v)
				containsID := slices.ContainsFunc(data.Contributors, func(e types.Contributor) bool {
					return e.ID != "" && e.ID == contributor.ID
				})
				if containsID {
					log.Printf("Contributor with ID %s already exists", contributor.ID)
				} else {
					data.Contributors = append(data.Contributors, contributor)

				}
			}
		}
	}

	if len(content.Attributes.Titles) > 0 {
		for _, v := range content.Attributes.Titles {
			var t string
			if slices.Contains([]string{"MainTitle", "Subtitle", "TranslatedTitle"}, v.TitleType) {
				t = v.TitleType
			}
			data.Titles = append(data.Titles, types.Title{
				Title:    v.Title,
				Type:     t,
				Language: v.Lang,
			})
		}
	}

	if content.Attributes.Publisher != "" {
		data.Publisher = types.Publisher{
			Name: content.Attributes.Publisher,
		}
	}

	for _, v := range content.Attributes.Dates {
		if v.DateType == "Accepted" {
			data.Date.Accepted = v.Date
		}
		if v.DateType == "Available" {
			data.Date.Available = v.Date
		}
		if v.DateType == "Collected" {
			data.Date.Collected = v.Date
		}
		if v.DateType == "Copyrighted" {
			data.Date.Copyrighted = v.Date
		}
		if v.DateType == "Created" {
			data.Date.Created = v.Date
		}
		if v.DateType == "Issued" {
			data.Date.Published = v.Date
		} else if v.DateType == "Published" {
			data.Date.Published = v.Date
		}
		if v.DateType == "Submitted" {
			data.Date.Submitted = v.Date
		}
		if v.DateType == "Updated" {
			data.Date.Updated = v.Date
		}
		if v.DateType == "Valid" {
			data.Date.Valid = v.Date
		}
		if v.DateType == "Withdrawn" {
			data.Date.Withdrawn = v.Date
		}
		if v.DateType == "Other" {
			data.Date.Other = v.Date
		}
	}

	data.Container = types.Container{
		Identifier:     content.Attributes.Container.Identifier,
		IdentifierType: content.Attributes.Container.IdentifierType,
		Type:           content.Attributes.Container.Type,
		Title:          content.Attributes.Container.Title,
		Volume:         content.Attributes.Container.Volume,
		Issue:          content.Attributes.Container.Issue,
		FirstPage:      content.Attributes.Container.FirstPage,
		LastPage:       content.Attributes.Container.LastPage,
	}

	if len(content.Attributes.Descriptions) > 0 {
		for _, v := range content.Attributes.Descriptions {
			var t string
			if slices.Contains([]string{"Abstract", "Summary", "Methods", "TechnicalInfo", "Other"}, v.DescriptionType) {
				t = v.DescriptionType
			} else {
				t = "Other"
			}
			description := utils.Sanitize(v.Description)
			data.Descriptions = append(data.Descriptions, types.Description{
				Description: description,
				Type:        t,
				Language:    v.Lang,
			})
		}
	}

	if len(content.Attributes.Subjects) > 0 {
		for _, v := range content.Attributes.Subjects {
			data.Subjects = append(data.Subjects, types.Subject{
				Subject: v.Subject,
			})
		}
	}

	data.Language = content.Attributes.Language

	if len(content.Attributes.RightsList) > 0 {
		url := content.Attributes.RightsList[0].RightsURI
		id := utils.UrlToSPDX(url)
		if id == "" {
			log.Printf("License URL %s not found in SPDX", url)
		}
		data.License = types.License{
			ID:  id,
			Url: url,
		}
	}

	data.Version = content.Attributes.Version

	if len(content.Attributes.RelatedIdentifiers) > 0 {
		supportedRelations := []string{
			"Cites",
			"References",
		}
		for i, v := range content.Attributes.RelatedIdentifiers {
			if slices.Contains(supportedRelations, v.RelationType) {
				isDoi, _ := regexp.MatchString(`^10\.\d{4,9}/.+$`, v.RelatedIdentifier)
				var doi, unstructured string
				if isDoi {
					doi = doiutils.NormalizeDOI(v.RelatedIdentifier)
				} else {
					unstructured = v.RelatedIdentifier
				}
				data.References = append(data.References, types.Reference{
					Key:          "ref" + strconv.Itoa(i+1),
					Doi:          doi,
					Unstructured: unstructured,
				})
			}
		}
	}

	if len(content.Attributes.RelatedIdentifiers) > 0 {
		supportedRelations := []string{
			"IsNewVersionOf",
			"IsPreviousVersionOf",
			"IsVersionOf",
			"HasVersion",
			"IsPartOf",
			"HasPart",
			"IsVariantFormOf",
			"IsOriginalFormOf",
			"IsIdenticalTo",
			"IsTranslationOf",
			"IsReviewedBy",
			"Reviews",
			"IsPreprintOf",
			"HasPreprint",
			"IsSupplementTo",
		}
		for _, v := range content.Attributes.RelatedIdentifiers {
			if slices.Contains(supportedRelations, v.RelationType) {
				isDoi, _ := regexp.MatchString(`^10\.\d{4,9}/.+$`, v.RelatedIdentifier)
				identifier := v.RelatedIdentifier
				if isDoi {
					identifier = doiutils.NormalizeDOI(v.RelatedIdentifier)
				} else if {
					identifier = v.RelatedIdentifier
				}
				data.Relations = append(data.Relations, types.Relation{
					ID:   identifier,
					Type: v.RelationType,
				})
			}
		}
	}

	if len(content.Attributes.FundingReferences) > 0 {
		for _, v := range content.Attributes.FundingReferences {
			data.FundingReferences = append(data.FundingReferences, types.FundingReference{
				FunderIdentifier:     v.FunderIdentifier,
				FunderIdentifierType: v.FunderIdentifierType,
				FunderName:           v.FunderName,
				AwardNumber:          v.AwardNumber,
				AwardURI:             v.AwardURI,
			})
		}
	} else {
		data.FundingReferences = []types.FundingReference{}
	}

	if len(content.Attributes.GeoLocations) > 0 {
		for _, v := range content.Attributes.GeoLocations {
			data.GeoLocations = append(data.GeoLocations, types.GeoLocation{
				GeoLocationPoint: types.GeoLocationPoint{
					PointLongitude: v.GeoLocationPoint.PointLongitude,
					PointLatitude:  v.GeoLocationPoint.PointLatitude,
				},
				GeoLocationPlace: v.GeoLocationPlace,
				GeoLocationBox: types.GeoLocationBox{
					EastBoundLongitude: v.GeoLocationBox.EastBoundLongitude,
					WestBoundLongitude: v.GeoLocationBox.WestBoundLongitude,
					SouthBoundLatitude: v.GeoLocationBox.SouthBoundLatitude,
					NorthBoundLatitude: v.GeoLocationBox.NorthBoundLatitude,
				},
			})
		}
	}

	data.Files = []types.File{}
	// sizes and formats are part of the file object, but can't be mapped directly

	data.ArchiveLocations = []string{}

	data.Provider = "DataCite"

	return data, nil
}

func GetContributor(v types.DCContributor) types.Contributor {
	var t string
	if len(v.NameType) > 2 {
		t = v.NameType[:len(v.NameType)-2]
	}
	var id string
	if len(v.NameIdentifiers) > 0 {
		ni := v.NameIdentifiers[0]
		if ni.NameIdentifierScheme == "ORCID" || ni.NameIdentifierScheme == "https://orcid.org/" {
			id = utils.NormalizeORCID(ni.NameIdentifier)
			t = "Person"
		} else if ni.NameIdentifierScheme == "ROR" {
			id = ni.NameIdentifier
			t = "Organization"
		} else {
			id = ni.NameIdentifier
		}
	}
	name := v.Name
	GivenName := v.GivenName
	FamilyName := v.FamilyName
	if t == "" && (v.GivenName != "" || v.FamilyName != "") {
		t = "Person"
	} else if t == "" {
		t = "Organization"
	}
	if t == "Person" && name != "" {
		// split name for type Person into given/family name if not already provided
		names := strings.Split(name, ",")
		if len(names) == 2 {
			GivenName = strings.TrimSpace(names[1])
			FamilyName = names[0]
			name = ""
		}
	}
	var affiliations []types.Affiliation
	for _, a := range v.Affiliation {
		affiliations = append(affiliations, types.Affiliation{
			ID:   "",
			Name: a,
		})
	}
	var roles []string
	if slices.Contains(CommonmetaContributorRoles, v.ContributorType) {
		roles = append(roles, v.ContributorType)
	} else {
		roles = append(roles, "Author")
	}
	return types.Contributor{
		ID:               id,
		Type:             t,
		Name:             name,
		GivenName:        GivenName,
		FamilyName:       FamilyName,
		ContributorRoles: roles,
		Affiliations:     affiliations,
	}
}
