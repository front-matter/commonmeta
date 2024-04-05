package schemautils_test

import (
	"commonmeta/metadata"
	"commonmeta/schemautils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func TestJSONSchemaErrors(t *testing.T) {
	t.Parallel()
	type testCase struct {
		meta metadata.Metadata
		want int
	}
	m := metadata.Metadata{
		ID:   "https://doi.org/10.7554/elife.01567",
		Type: "JournalArticle",
		Url:  "https://elifesciences.org/articles/01567",
	}

	// missing required ID, defaults to empty string
	n := metadata.Metadata{
		Type: "JournalArticle",
	}

	// Type is not supported
	o := metadata.Metadata{
		ID:   "https://doi.org/10.7554/elife.01567",
		Type: "Umbrella",
	}

	testCases := []testCase{
		{meta: m, want: 0},
		{meta: n, want: 2},
		{meta: o, want: 1},
	}
	for _, tc := range testCases {
		documentJSON, err := json.Marshal(tc.meta)
		if err != nil {
			log.Fatal(err)
		}
		result := schemautils.JSONSchemaErrors(documentJSON)
		got := len(result.Errors())
		if tc.want != got {
			t.Errorf("want %d, got %d", tc.want, got)
		}
	}
}

func TestJSONSchemaErrorsTestdata(t *testing.T) {
	t.Parallel()
	type testCase struct {
		meta   string
		schema string
		want   int
	}

	testCases := []testCase{
		{meta: "citeproc.json", schema: "csl-data", want: 0},
		{meta: "datacite.json", schema: "datacite-v4.5", want: 3},
	}
	for _, tc := range testCases {
		data, err := os.ReadFile("../testdata/" + tc.meta)
		if err != nil {
			fmt.Print(err)
		}
		result := schemautils.JSONSchemaErrors(data, tc.schema)
		got := len(result.Errors())
		if tc.want != got {
			t.Errorf("want %d, got %d", tc.want, got)
		}
	}
}
