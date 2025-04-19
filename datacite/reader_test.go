package datacite_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/datacite"

	"github.com/front-matter/commonmeta/doiutils"

	"github.com/google/go-cmp/cmp"
)

// func TestGet(t *testing.T) {
// 	t.Parallel()

// 	type testCase struct {
// 		id   string
// 		want string
// 		err  error
// 	}

// 	// PID as DOI Url
// 	publication := datacite.Content{
// 		ID:  "https://doi.org/10.5281/zenodo.5244404",
// 		URL: "https://zenodo.org/record/5244404",
// 	}
// 	// PID as DOI string
// 	presentation := datacite.Content{
// 		ID:  "10.5281/zenodo.8173303",
// 		URL: "https://zenodo.org/record/8173303",
// 	}

// 	testCases := []testCase{
// 		{id: presentation.ID, want: presentation.URL, err: nil},
// 		{id: publication.ID, want: publication.URL, err: nil},
// 	}

// 	for _, tc := range testCases {
// 		got, err := datacite.Get(tc.id)
// 		if tc.want != got.URL {
// 			t.Errorf("Get DataCite(%v): want %v, got %v, error %v",
// 				tc.id, tc.want, got.URL, err)
// 		}
// 	}
// }

func TestFetch(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name string
		id   string
	}

	testCases := []testCase{
		{name: "dataset", id: "https://doi.org/10.5061/dryad.8515"},
		{name: "blog posting", id: "https://doi.org/10.5438/zhyx-n122"},
		{name: "proceedings article", id: "https://doi.org/10.4230/lipics.tqc.2013.93"},
		{name: "subject scheme FOR", id: "https://doi.org/10.6084/m9.figshare.1449060"},
		{name: "geolocation box", id: "https://doi.org/10.6071/z7wc73"},
	}
	match := true
	for _, tc := range testCases {
		got, err := datacite.Fetch(tc.id, match)
		if err != nil {
			t.Errorf("DataCite Metadata(%v): error %v", tc.id, err)
		}
		// read json file from testdata folder and convert to Data struct
		doi, ok := doiutils.ValidateDOI(tc.id)
		if !ok {
			t.Fatal(errors.New("invalid doi"))
		}
		filename := strings.ReplaceAll(doi, "/", "_") + ".json"
		filepath := filepath.Join("testdata", filename)
		content, err := os.ReadFile(filepath)
		if err != nil {
			t.Fatal(err)
		}
		var want commonmeta.Data
		err = json.Unmarshal(content, &want)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("FetchDatacite(%s) mismatch (-want +got):\n%s", tc.id, diff)
		}
	}
}

func TestQueryURL(t *testing.T) {
	t.Parallel()

	type testCase struct {
		number        int
		page          int
		client        string
		type_         string
		sample        bool
		year          string
		language      string
		orcid         string
		ror           string
		hasORCID      bool
		hasROR        bool
		hasReferences bool
		hasRelation   bool
		hasAbstract   bool
		hasAward      bool
		hasLicense    bool
		want          string
	}

	testCases := []testCase{
		{want: "https://api.datacite.org/dois?page[size]=10&page[number]=1&sort=-published"},
		{sample: true, want: "https://api.datacite.org/dois?page[size]=10&random=true"},
		{number: 100, want: "https://api.datacite.org/dois?page[size]=100&page[number]=1&sort=-published"},
		{number: 10, page: 3, want: "https://api.datacite.org/dois?page[size]=10&page[number]=3&sort=-published"},
		{sample: true, client: "cern.zenodo", want: "https://api.datacite.org/dois?page[size]=10&random=true&client-id=cern.zenodo"},
		{sample: true, year: "2022", want: "https://api.datacite.org/dois?page[size]=10&random=true&query=publicationYear:2022"},
		{sample: true, language: "es", want: "https://api.datacite.org/dois?page[size]=10&random=true&query=language:es"},
		{sample: true, orcid: "0000-0002-8635-8390", want: "https://api.datacite.org/dois?page[size]=10&random=true&query=creators.nameIdentifiers.nameIdentifier:0000-0002-8635-8390"},
		{sample: true, ror: "039kas258", want: "https://api.datacite.org/dois?page[size]=10&random=true&affiliation-id=039kas258"},
		{sample: true, hasORCID: true, want: "https://api.datacite.org/dois?page[size]=10&random=true&query=creators.nameIdentifiers.nameIdentifierScheme:ORCID"},
		{sample: true, hasROR: true, want: "https://api.datacite.org/dois?page[size]=10&random=true&query=creators.affiliation.affiliationIdentifierScheme:ROR&affiliation=true"},
		{sample: true, hasReferences: true, want: "https://api.datacite.org/dois?page[size]=10&random=true&query=relatedIdentifiers.relationType:Cites"},
		{sample: true, hasRelation: true, want: "https://api.datacite.org/dois?page[size]=10&random=true&query=relatedIdentifiers.relationType:*"},
		{sample: true, hasAbstract: true, want: "https://api.datacite.org/dois?page[size]=10&random=true&query=descriptions.descriptionType:Abstract"},
		{sample: true, hasAward: true, want: "https://api.datacite.org/dois?page[size]=10&random=true&query=fundingReferences.funderIdentifier:*"},
		{sample: true, hasLicense: true, want: "https://api.datacite.org/dois?page[size]=10&random=true&query=rightsList.rightsIdentifierScheme:SPDX"},
	}
	for _, tc := range testCases {
		got := datacite.QueryURL(tc.number, tc.page, tc.client, tc.type_, tc.sample, tc.year, tc.language, tc.orcid, tc.ror, tc.hasORCID, tc.hasROR, tc.hasReferences, tc.hasRelation, tc.hasAbstract, tc.hasAward, tc.hasLicense)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("DataciteApiQueryUrl mismatch (-want +got):\n%s", diff)
		}
	}
}

// func TestGetDataciteSample(t *testing.T) {
// 	t.Parallel()

// 	type testCase struct {
// 		number int
// 		want   string
// 	}

// 	testCases := []testCase{
// 		{number: 10, want: "https://api.datacite.org/works?query=member:340,type:journal-article&rows=10"},
// 	}
// 	for _, tc := range testCases {
// 		got, err := datacite.GetDataciteSample(tc.number)
// 		if err != nil {
// 			t.Errorf("Datacite Sample(%v): error %v", tc.number, err)
// 		}
// 		if diff := cmp.Diff(tc.want, got); diff != "" {
// 			t.Errorf("DataciteApiSampleUrl mismatch (-want +got):\n%s", diff)
// 		}
// 	}
// }
