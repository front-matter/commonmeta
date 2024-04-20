package schemautils_test

import (
	"commonmeta/schemautils"
	"commonmeta/types"
	"encoding/json"
	"path/filepath"

	"fmt"
	"log"
	"os"
	"testing"

	"sigs.k8s.io/yaml"
)

func TestJSONSchemaErrors(t *testing.T) {
	t.Parallel()
	type testCase struct {
		meta types.Data
		want int
	}
	m := types.Data{
		ID:   "https://doi.org/10.7554/elife.01567",
		Type: "JournalArticle",
		Url:  "https://elifesciences.org/articles/01567",
	}

	// missing required ID, defaults to empty string
	n := types.Data{
		Type: "JournalArticle",
	}

	// Type is not supported
	o := types.Data{
		ID:   "https://doi.org/10.1515/9789048535248-011",
		Type: "Umbrella",
	}

	testCases := []testCase{
		{meta: m, want: 1},
		{meta: n, want: 3},
		{meta: o, want: 2},
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
			fmt.Printf("The document %s is not valid. see errors :\n", tc.meta.ID)
			for _, desc := range result.Errors() {
				fmt.Printf("- %s\n", desc)
			}
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
		{meta: "journal_article.commonmeta.json", schema: "commonmeta_v0.13", want: 0},
		{meta: "citeproc.json", schema: "csl-data", want: 0},
		{meta: "datacite.json", schema: "datacite-v4.5", want: 3},
		{meta: "datacite-instrument.json", schema: "datacite-v4.5", want: 27},
		{meta: "datacite_software_version.json", schema: "datacite-v4.5", want: 7},
	}
	for _, tc := range testCases {
		filepath := filepath.Join("testdata", tc.meta)
		data, err := os.ReadFile(filepath)
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

func TestJSONSchemaErrorsTestdataYAML(t *testing.T) {
	t.Parallel()
	type testCase struct {
		meta   string
		schema string
		want   int
	}

	testCases := []testCase{
		{meta: "CITATION.cff", schema: "cff_v1.2.0", want: 0},
	}
	for _, tc := range testCases {
		filepath := filepath.Join("testdata", tc.meta)
		data, err := os.ReadFile(filepath)
		if err != nil {
			fmt.Print(err)
		}
		YAMLdata, err := yaml.YAMLToJSON(data)
		if err != nil {
			fmt.Print(err)
		}
		result := schemautils.JSONSchemaErrors(YAMLdata, tc.schema)
		got := len(result.Errors())
		if tc.want != got {
			t.Errorf("want %d, got %d", tc.want, got)
		}
	}
}
