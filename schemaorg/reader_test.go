package schemaorg_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/schemaorg"
	"github.com/google/go-cmp/cmp"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type testCase struct {
		id   string
		want string
	}

	testCases := []testCase{
		{id: "https://blog.front-matter.io/posts/eating-your-own-dog-food", want: "https://doi.org/10.53731/r79vxn1-97aq74v-ag58n"},
		{id: "https://www.zenodo.org/records/1196821", want: "https://doi.org/10.5281/zenodo.1196821"},
		{id: "https://doi.pangaea.de/10.1594/PANGAEA.836178", want: "https://doi.org/10.1594/PANGAEA.836178"},
		// {id: "https://dataverse.harvard.edu/dataset.xhtml?persistentId=doi:10.7910/DVN/NJ7XSO", want: "Harvard Dataverse"},
		{id: "https://datadryad.org/stash/dataset/doi:10.5061/dryad.8515", want: "https://doi.org/10.5061/dryad.8515"},
	}
	for _, tc := range testCases {
		content, err := schemaorg.Get(tc.id)
		if err != nil {
			t.Errorf("Schemaorg Get (%v): error %v", tc.id, err)
		}
		got := content.ID
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("Schemaorg Get (%v): -want +got %s", tc.id, diff)
		}
	}
}

func TestFetch(t *testing.T) {
	t.Parallel()

	type testCase struct {
		id  string
		url string
	}

	testCases := []testCase{
		{id: "10.53731/r79vxn1-97aq74v-ag58n", url: "https://blog.front-matter.io/posts/eating-your-own-dog-food"},
		{id: "10.5281/zenodo.1196821", url: "https://www.zenodo.org/records/1196821"},
		{id: "10.1594/pangaea.836178", url: "https://doi.pangaea.de/10.1594/PANGAEA.836178"},
		// {id: "10.7910/dvn/nj7xso", url: "https://dataverse.harvard.edu/dataset.xhtml?persistentId=doi:10.7910/DVN/NJ7XSO"},
		{id: "10.5061/dryad.8515", url: "https://datadryad.org/stash/dataset/doi:10.5061/dryad.8515"},
		{id: "10.7554/eLife.93170.2", url: "https://elifesciences.org/reviewed-preprints/93170"},
	}

	for _, tc := range testCases {
		got, err := schemaorg.Fetch(tc.url)
		if err != nil {
			t.Errorf("Schemaorg Fetch (%v): error %v", tc.url, err)
			got = commonmeta.Data{}
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
		_ = json.Unmarshal(content, &want)
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Schemaorg Fetch (%v): -want +got %s", tc.id, diff)
		}
	}
}
