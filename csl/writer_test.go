package csl_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/crossref"
	"github.com/front-matter/commonmeta/csl"
	"github.com/front-matter/commonmeta/datacite"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/google/go-cmp/cmp"
)

func TestWrite(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name string
		id   string
		from string
	}

	testCases := []testCase{
		{name: "journal article", id: "https://doi.org/10.7554/elife.01567", from: "crossref"},
		{name: "preprint", id: "https://doi.org/10.1101/097196", from: "crossref"},
		{name: "dataset", id: "https://doi.org/10.5061/dryad.8515", from: "datacite"},
	}
	match := true
	for _, tc := range testCases {
		var data commonmeta.Data
		var err error
		if tc.from == "crossref" {
			data, err = crossref.Fetch(tc.id, match)
		} else if tc.from == "datacite" {
			data, err = datacite.Fetch(tc.id, match)
		}
		if err != nil {
			t.Errorf("Crossref Fetch (%v): error %v", tc.id, err)
		}
		got, err := csl.Write(data)
		if err != nil {
			t.Errorf("CSL Write (%v): error %v", tc.id, err)
		}
		// read json file from testdata folder and convert to CSL struct
		doi, ok := doiutils.ValidateDOI(tc.id)
		if !ok {
			t.Fatal("invalid doi")
		}
		filename := strings.ReplaceAll(doi, "/", "_") + ".json"
		filepath := filepath.Join("testdata", filename)
		bytes, err := os.ReadFile(filepath)
		if err != nil {
			t.Fatal(err)
		}

		var want csl.CSL
		err = json.Unmarshal(bytes, &want)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Fetch (%s) mismatch (-want +got):\n%s", tc.id, diff)
		}
	}
}
