package datacite_test

import (
	"commonmeta/datacite"
	"commonmeta/doiutils"
	"commonmeta/types"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetDatacite(t *testing.T) {
	t.Parallel()

	type testCase struct {
		id   string
		want string
		err  error
	}

	// PID as DOI Url
	publication := types.Content{
		ID: "https://doi.org/10.5281/zenodo.5244404",
		Attributes: types.Attributes{
			Url: "https://zenodo.org/record/5244404",
		},
	}
	// PID as DOI string
	presentation := types.Content{
		ID: "10.5281/zenodo.8173303",
		Attributes: types.Attributes{
			Url: "https://zenodo.org/record/8173303",
		},
	}

	testCases := []testCase{
		{id: presentation.ID, want: presentation.Attributes.Url, err: nil},
		{id: publication.ID, want: publication.Attributes.Url, err: nil},
	}

	for _, tc := range testCases {
		got, err := datacite.GetDatacite(tc.id)
		if tc.want != got.Attributes.Url {
			t.Errorf("Get DataCite(%v): want %v, got %v, error %v",
				tc.id, tc.want, got, err)
		}
	}
}

func TestFetchDatacite(t *testing.T) {
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
		{name: "geolocation point", id: "10.4121/UUID:7B900822-4EFE-42F1-9B6E-A099EDA4BA02"},
	}
	for _, tc := range testCases {
		got, err := datacite.FetchDatacite(tc.id)
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
		want := types.Data{}
		err = json.Unmarshal(content, &want)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("FetchDatacite(%s) mismatch (-want +got):\n%s", tc.id, diff)
		}
	}
}
