package schemautils_test

import (
	"encoding/json"
	"path/filepath"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/schemautils"

	"fmt"
	"log"
	"os"
	"testing"

	"sigs.k8s.io/yaml"
)

func TestSchemaErrors(t *testing.T) {
	t.Parallel()
	type testCase struct {
		meta commonmeta.Data
		want error
	}
	m := commonmeta.Data{
		ID:   "https://doi.org/10.7554/elife.01567",
		Type: "JournalArticle",
		URL:  "https://elifesciences.org/articles/01567",
	}

	// missing required ID, defaults to empty string
	n := commonmeta.Data{
		Type: "JournalArticle",
	}

	// Type is not supported
	o := commonmeta.Data{
		ID:   "https://doi.org/10.1515/9789048535248-011",
		Type: "Umbrella",
	}

	testCases := []testCase{
		{meta: m, want: nil},
		{meta: n, want: nil},
		{meta: o, want: nil},
	}
	for _, tc := range testCases {
		documentJSON, err := json.Marshal(tc.meta)
		if err != nil {
			log.Fatal(err)
		}
		err = schemautils.SchemaErrors(documentJSON)
		if tc.want != err {
			t.Errorf("want %d, got %d", tc.want, err)
			fmt.Printf("The document %s is not valid. see errors :\n", tc.meta.ID)
			fmt.Printf("%v", err)
		}
	}
}

func TestSchemaErrorsTestdata(t *testing.T) {
	t.Parallel()
	type testCase struct {
		meta   string
		schema string
		err    error
	}

	testCases := []testCase{
		{meta: "journal_article.commonmeta.json", schema: "commonmeta_v0.15", err: nil},
		{meta: "citeproc.json", schema: "csl-data", err: nil},
		{meta: "datacite.json", schema: "datacite-v4.5", err: nil},
		{meta: "datacite-instrument.json", schema: "datacite-v4.5", err: nil},
		{meta: "datacite_software_version.json", schema: "datacite-v4.5", err: nil},
		{meta: "inveniordm.json", schema: "invenio-rdm-v0.1", err: nil},
	}
	for _, tc := range testCases {
		filepath := filepath.Join("testdata", tc.meta)
		data, err := os.ReadFile(filepath)
		if err != nil {
			fmt.Print(err)
		}
		err = schemautils.SchemaErrors(data, tc.schema)
		if tc.err != err {
			t.Errorf("want %v, got %d", tc.meta, err)
		}
	}
}

func TestSchemaErrorsTestdataYAML(t *testing.T) {
	t.Parallel()
	type testCase struct {
		meta   string
		schema string
		err    error
	}

	testCases := []testCase{
		{meta: "CITATION.cff", schema: "cff_v1.2.0", err: nil},
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
		got := schemautils.SchemaErrors(YAMLdata, tc.schema)
		if tc.err != got {
			t.Errorf("want %d, got %d", tc.err, got)
		}
	}
}
