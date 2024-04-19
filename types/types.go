package types

type Content struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
	Abstract   string     `json:"abstract"`
	Archive    []string   `json:"archive"`
	Author     []struct {
		Given       string `json:"given"`
		Family      string `json:"family"`
		Name        string `json:"name"`
		ORCID       string `json:"ORCID"`
		Sequence    string `json:"sequence"`
		Affiliation []struct {
			ID []struct {
				ID     string `json:"id"`
				IDType string `json:"id-type"`
			} `json:"id"`
			Name string `json:"name"`
		} `json:"affiliation"`
	} `json:"author"`
	Blog struct {
		ISSN        string `json:"issn"`
		License     string `json:"license"`
		Title       string `json:"title"`
		HomePageUrl string `json:"home_page_url"`
	} `json:"blog"`
	ContainerTitle []string   `json:"container-title"`
	DOI            string     `json:"doi"`
	Files          []struct{} `json:"files"`
	Funder         []struct {
		DOI   string   `json:"DOI"`
		Name  string   `json:"name"`
		Award []string `json:"award"`
	} `json:"funder"`
	GroupTitle string `json:"group-title"`
	GUID       string `json:"guid"`
	Issue      string `json:"issue"`
	Published  struct {
		DateAsParts [][]int `json:"date-parts"`
		DateTime    string  `json:"date-time"`
	} `json:"published"`
	Issued struct {
		DateAsParts [][]int `json:"date-parts"`
		DateTime    string  `json:"date-time"`
	} `json:"issued"`
	Created struct {
		DateAsParts [][]int `json:"date-parts"`
		DateTime    string  `json:"date-time"`
	} `json:"created"`
	ISSN     []string `json:"ISSN"`
	ISBNType []struct {
		Value string `json:"value"`
		Type  string `json:"type"`
	} `json:"isbn-type"`
	Language string `json:"language"`
	License  []struct {
		Url            string `json:"URL"`
		ContentVersion string `json:"content-version"`
	} `json:"license"`
	Link []struct {
		ContentType string `json:"content-type"`
		Url         string `json:"url"`
	} `json:"link"`
	Page        string `json:"page"`
	PublishedAt string `json:"published_at"`
	Publisher   string `json:"publisher"`
	Reference   []struct {
		Key          string `json:"key"`
		DOI          string `json:"DOI"`
		ArticleTitle string `json:"article-title"`
		Year         string `json:"year"`
		Unstructured string `json:"unstructured"`
	} `json:"reference"`
	Relation struct {
		IsNewVersionOf []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-new-version-of"`
		IsPreviousVersionOf []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-previous-version-of"`
		IsVersionOf []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-version-of"`
		HasVersion []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"has-version"`
		IsPartOf []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-part-of"`
		HasPart []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"has-part"`
		IsVariantFormOf []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-variant-form-of"`
		IsOriginalFormOf []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-original-form-of"`
		IsIdenticalTo []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-identical-to"`
		IsTranslationOf []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-translation-of"`
		IsReviewedBy []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-reviewed-by"`
		Reviews []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"reviews"`
		HasReview []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"has-review"`
		IsPreprintOf []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-preprint-of"`
		HasPreprint []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"has-preprint"`
		IsSupplementTo []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-supplement-to"`
		IsSupplementedBy []struct {
			ID     string `json:"id"`
			IDType string `json:"id-type"`
		} `json:"is-supplemented-by"`
	} `json:"relation"`
	Resource struct {
		Primary struct {
			ContentType string `json:"content_type"`
			URL         string `json:"url"`
		} `json:"primary"`
	} `json:"resource"`
	Subject   []string `json:"subject"`
	Subtitle  []string `json:"subtitle"`
	Summary   string   `json:"summary"`
	Tags      []string `json:"tags"`
	Title     []string `json:"title"`
	UpdatedAt string   `json:"updated_at"`
	Url       string   `json:"url"`
	Version   string   `json:"version"`
	Via       string   `json:"via"`
	Volume    string   `json:"volume"`
}

// used by Datacite API
type Attributes struct {
	DOI                  string `json:"doi"`
	Prefix               string `json:"prefix"`
	Suffix               string `json:"suffix"`
	AlternateIdentifiers []struct {
		Identifier     string `json:"identifier"`
		IdentifierType string `json:"identifierType"`
	} `json:"alternateIdentifiers"`
	Creators []struct {
		Name            string `json:"name"`
		GivenName       string `json:"givenName"`
		FamilyName      string `json:"familyName"`
		NameType        string `json:"nameType"`
		NameIdentifiers []struct {
			SchemeURI            string `json:"schemeUri"`
			NameIdentifier       string `json:"nameIdentifier"`
			NameIdentifierScheme string `json:"nameIdentifierScheme"`
		} `json:"nameIdentifiers"`
		Affiliation []string `json:"affiliation"`
	} `json:"creators"`
	Publisher string `json:"publisher"`
	Container struct {
		Type           string `json:"type"`
		Identifier     string `json:"identifier"`
		IdentifierType string `json:"identifierType"`
		Title          string `json:"title"`
		Volume         string `json:"volume"`
		Issue          string `json:"issue"`
		FirstPage      string `json:"firstPage"`
		LastPage       string `json:"lastPage"`
	} `json:"container"`
	PublicationYear int `json:"publicationYear"`
	Titles          []struct {
		Title     string `json:"title"`
		TitleType string `json:"titleType"`
		Lang      string `json:"lang"`
	} `json:"titles"`
	Url      string `json:"url"`
	Subjects []struct {
		Subject string `json:"subject"`
	} `json:"subjects"`
	Contributors []struct {
		Name            string `json:"name"`
		GivenName       string `json:"givenName"`
		FamilyName      string `json:"familyName"`
		NameType        string `json:"nameType"`
		NameIdentifiers []struct {
			SchemeURI            string `json:"schemeUri"`
			NameIdentifier       string `json:"nameIdentifier"`
			NameIdentifierScheme string `json:"nameIdentifierScheme"`
		} `json:"nameIdentifiers"`
		Affiliation     []string `json:"affiliation"`
		ContributorType string   `json:"contributorType"`
	} `json:"contributors"`
	Dates []struct {
		Date            string `json:"date"`
		DateType        string `json:"dateType"`
		DateInformation string `json:"dateInformation"`
	} `json:"dates"`
	Language string `json:"language"`
	Types    struct {
		ResourceTypeGeneral string `json:"resourceTypeGeneral"`
		ResourceType        string `json:"resourceType"`
	} `json:"types"`
	RelatedIdentifiers []struct {
		RelatedIdentifier     string `json:"relatedIdentifier"`
		RelatedIdentifierType string `json:"relatedIdentifierType"`
		RelationType          string `json:"relationType"`
	} `json:"relatedIdentifiers"`
	Sizes      []string `json:"sizes"`
	Formats    []string `json:"formats"`
	Version    string   `json:"version"`
	RightsList []struct {
		Rights                 string `json:"rights"`
		RightsURI              string `json:"rightsUri"`
		SchemeURI              string `json:"schemeUri"`
		RightsIdentifier       string `json:"rightsIdentifier"`
		RightsIdentifierScheme string `json:"rightsIdentifierScheme"`
	}
	Descriptions []struct {
		Description     string `json:"description"`
		DescriptionType string `json:"descriptionType"`
		Lang            string `json:"lang"`
	} `json:"descriptions"`
	GeoLocations []struct {
		GeoLocationPoint struct {
			PointLongitude float64 `json:"pointLongitude,string"`
			PointLatitude  float64 `json:"pointLatitude,string"`
		} `json:"geoLocationPoint"`
		GeoLocationBox struct {
			WestBoundLongitude float64 `json:"westBoundLongitude,string"`
			EastBoundLongitude float64 `json:"eastBoundLongitude,string"`
			SouthBoundLatitude float64 `json:"southBoundLatitude,string"`
			NorthBoundLatitude float64 `json:"northBoundLatitude,string"`
		} `json:"geoLocationBox"`
		GeoLocationPlace string `json:"geoLocationPlace"`
	} `json:"geoLocations"`
	FundingReferences []struct {
		FunderName           string `json:"funderName"`
		FunderIdentifier     string `json:"funderIdentifier"`
		FunderIdentifierType string `json:"funderIdentifierType"`
		AwardNumber          string `json:"awardNumber"`
		AwardURI             string `json:"awardUri"`
	} `json:"fundingReferences"`
}

type Data struct {
	// required fields
	ID   string `db:"id" json:"id"`
	Type string `db:"type" json:"type"`

	// optional fields
	AdditionalType       string                `db:"additional_type" json:"additional_type,omitempty"`
	Url                  string                `db:"url" json:"url,omitempty"`
	Contributors         []Contributor         `db:"contributors" json:"contributors,omitempty"`
	Publisher            Publisher             `db:"publisher" json:"publisher,omitempty"`
	Date                 Date                  `db:"date" json:"date,omitempty"`
	Titles               []Title               `db:"titles" json:"titles,omitempty"`
	Container            Container             `db:"container" json:"container,omitempty"`
	Subjects             []Subject             `db:"subjects" json:"subjects,omitempty"`
	Sizes                []string              `db:"sizes" json:"sizes,omitempty"`
	Formats              []string              `db:"formats" json:"formats,omitempty"`
	Language             string                `db:"language" json:"language,omitempty"`
	License              License               `db:"license" json:"license,omitempty"`
	Version              string                `db:"version" json:"version,omitempty"`
	References           []Reference           `db:"references" json:"references,omitempty"`
	Relations            []Relation            `db:"relations" json:"relations,omitempty"`
	FundingReferences    []FundingReference    `db:"funding_references" json:"funding_references,omitempty"`
	Descriptions         []Description         `db:"descriptions" json:"descriptions,omitempty"`
	GeoLocations         []GeoLocation         `db:"geo_locations" json:"geo_locations,omitempty"`
	Provider             string                `db:"provider" json:"provider,omitempty"`
	AlternateIdentifiers []AlternateIdentifier `db:"alternate_identifiers" json:"alternate_identifiers,omitempty"`
	Files                []File                `db:"files" json:"files,omitempty"`
	ArchiveLocations     []string              `db:"archive_locations" json:"archive_locations,omitempty"`
}

type Affiliation struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

type AlternateIdentifier struct {
	Identifier     string `json:"identifier"`
	IdentifierType string `json:"identifierType"`
}

type Container struct {
	Identifier     string `json:"identifier,omitempty"`
	IdentifierType string `json:"identifierType,omitempty"`
	Type           string `json:"type,omitempty"`
	Title          string `json:"title,omitempty"`
	FirstPage      string `json:"firstPage,omitempty"`
	LastPage       string `json:"lastPage,omitempty"`
	Volume         string `json:"volume,omitempty"`
	Issue          string `json:"issue,omitempty"`
}

type Contributor struct {
	ID               string        `json:"id,omitempty"`
	Type             string        `json:"type"`
	Name             string        `json:"name,omitempty"`
	GivenName        string        `json:"givenName,omitempty"`
	FamilyName       string        `json:"familyName,omitempty"`
	Affiliations     []Affiliation `json:"affiliations,omitempty"`
	ContributorRoles []string      `json:"contributorRoles"`
}

type Date struct {
	Created     string `json:"created,omitempty"`
	Submitted   string `json:"submitted,omitempty"`
	Accepted    string `json:"accepted,omitempty"`
	Published   string `json:"published,omitempty"`
	Updated     string `json:"updated,omitempty"`
	Accessed    string `json:"accessed,omitempty"`
	Available   string `json:"available,omitempty"`
	Copyrighted string `json:"copyrighted,omitempty"`
	Collected   string `json:"collected,omitempty"`
	Valid       string `json:"valid,omitempty"`
	Withdrawn   string `json:"withdrawn,omitempty"`
	Other       string `json:"other,omitempty"`
}

type Description struct {
	Description     string `json:"description"`
	DescriptionType string `json:"descriptionType,omitempty"`
	Language        string `json:"language,omitempty"`
}

type File struct {
	Bucket   string `json:"bucket,omitempty"`
	Key      string `json:"key,omitempty"`
	Checksum string `json:"checksum,omitempty"`
	Url      string `json:"url"`
	Size     int    `json:"size,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

type FundingReference struct {
	FunderIdentifier     string `json:"funderIdentifier,omitempty"`
	FunderIdentifierType string `json:"funderIdentifierType,omitempty"`
	FunderName           string `json:"funderName"`
	AwardNumber          string `json:"awardNumber,omitempty"`
	AwardURI             string `json:"award_uri,omitempty"`
}

type GeoLocation struct {
	GeoLocationPlace string           `json:"geoLocationPlace,omitempty"`
	GeoLocationPoint GeoLocationPoint `json:"geoLocationPoint,omitempty"`
	GeoLocationBox   GeoLocationBox   `json:"geoLocationBox,omitempty"`
}

type GeoLocationPoint struct {
	PointLongitude float64 `json:"pointLongitude,omitempty"`
	PointLatitude  float64 `json:"pointLatitude,omitempty"`
}

type GeoLocationBox struct {
	EastBoundLongitude float64 `json:"eastBoundLongitude,omitempty"`
	WestBoundLongitude float64 `json:"westBoundLongitude,omitempty"`
	SouthBoundLatitude float64 `json:"southBoundLatitude,omitempty"`
	NorthBoundLatitude float64 `json:"northBoundLatitude,omitempty"`
}

type GeoLocationPolygon struct {
	PolygonPoints  []GeoLocationPoint `json:"polygon_points,omitempty"`
	InPolygonPoint GeoLocationPoint   `json:"in_polygon_point,omitempty"`
}

type License struct {
	ID  string `json:"id,omitempty"`
	Url string `json:"url,omitempty"`
}

type Publisher struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

type Reference struct {
	Key             string `json:"key"`
	Doi             string `json:"doi,omitempty"`
	Title           string `json:"title,omitempty"`
	PublicationYear string `json:"publicationYear,omitempty"`
	Unstructured    string `json:"unstructured,omitempty"`
}

type Relation struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type Subject struct {
	Subject string `json:"subject"`
}

type Title struct {
	Title     string `json:"title"`
	TitleType string `json:"titleType,omitempty"`
	Language  string `json:"language,omitempty"`
}
