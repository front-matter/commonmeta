package crossref_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/front-matter/commonmeta-go/crossref"
	"github.com/front-matter/commonmeta-go/doiutils"
	"github.com/front-matter/commonmeta-go/types"

	"github.com/google/go-cmp/cmp"
)

func TestGetCrossref(t *testing.T) {
	t.Parallel()

	type testCase struct {
		id   string
		want string
		err  error
	}

	journalArticle := crossref.Content{
		ID:        "https://doi.org/10.7554/elife.01567",
		Publisher: "eLife Sciences Publications, Ltd",
	}
	postedContent := crossref.Content{
		ID:        "https://doi.org/10.1101/097196",
		Publisher: "Cold Spring Harbor Laboratory",
	}

	testCases := []testCase{
		{id: journalArticle.ID, want: journalArticle.Publisher, err: nil},
		{id: postedContent.ID, want: postedContent.Publisher, err: nil},
	}
	for _, tc := range testCases {
		got, err := crossref.GetCrossref(tc.id)
		if tc.want != got.Publisher {
			t.Errorf("Get Crossref(%v): want %v, got %v, error %v",
				tc.id, tc.want, got, err)
		}
	}
}

func TestFetchCrossref(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name string
		id   string
	}

	testCases := []testCase{
		{name: "test doi", id: "https://doi.org/10.5555/12345678"},
		{name: "journal article with data citation", id: "https://doi.org/10.7554/elife.01567"},
		{name: "posted content", id: "https://doi.org/10.1101/097196"},
		{name: "book", id: "https://doi.org/10.1017/9781108348843"},
		{name: "book chapter", id: "https://doi.org/10.1007/978-3-662-46370-3_13"},
		{name: "proceedings article", id: "https://doi.org/10.1145/3448016.3452841"},
		{name: "dataset", id: "https://doi.org/10.2210/pdb4hhb/pdb"},
		{name: "component", id: "https://doi.org/10.1371/journal.pmed.0030277.g001"},
		{name: "peer review", id: "https://doi.org/10.7554/elife.55167.sa2"},
		{name: "blog post", id: "https://doi.org/10.59350/2shz7-ehx26"},
		{name: "dissertation", id: "https://doi.org/10.14264/uql.2020.791"},
		{name: "with ror id", id: "https://doi.org/10.1364/oe.490112"},
		{name: "archived", id: "10.5694/j.1326-5377.1943.tb44329.x"},
	}
	for _, tc := range testCases {
		got, err := crossref.FetchCrossref(tc.id)
		if err != nil {
			t.Errorf("Crossref Metadata(%v): error %v", tc.id, err)
		}
		// read json file from testdata folder and convert to Data struct
		doi, ok := doiutils.ValidateDOI(tc.id)
		if !ok {
			t.Fatal("invalid doi")
		}
		filename := strings.ReplaceAll(doi, "/", "_") + ".json"
		filepath := filepath.Join("testdata", filename)
		content, err := os.ReadFile(filepath)
		if err != nil {
			t.Fatal(err)
		}
		want := types.Data{
			// FundingReferences: []types.FundingReference{},
		}
		err = json.Unmarshal(content, &want)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("FetchCrossref(%s) mismatch (-want +got):\n%s", tc.id, diff)
		}
	}
}

func TestCrossrefApiSampleUrl(t *testing.T) {
	t.Parallel()

	type testCase struct {
		number int
		member string
		_type  string
		want   string
	}

	testCases := []testCase{
		{number: 10, member: "340", _type: "journal-article", want: "https://api.crossref.org/works?filter=member%3A340%2Ctype%3Ajournal-article&sample=10"},
		{number: 1, member: "", _type: "posted-content", want: "https://api.crossref.org/works?filter=type%3Aposted-content&sample=1"},
		{number: 20, member: "78", _type: "", want: "https://api.crossref.org/works?filter=member%3A78&sample=20"},
		{number: 120, member: "", _type: "", want: "https://api.crossref.org/works?sample=120"},
	}
	for _, tc := range testCases {
		got := crossref.CrossrefApiSampleUrl(tc.number, tc.member, tc._type, false, false, false, false, false, false, false, false)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("CrossrefApiSampleUrl mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestGetCrossrefSample(t *testing.T) {
	t.Parallel()

	type testCase struct {
		number int
		member string
		_type  string
	}

	testCases := []testCase{
		{number: 3, member: "340", _type: "journal-article"},
		{number: 1, member: "", _type: "posted-content"},
		{number: 2, member: "", _type: ""},
	}
	for _, tc := range testCases {
		got, err := crossref.GetCrossrefSample(tc.number, tc.member, tc._type, false, false, false, false, false, false, false, false)
		if err != nil {
			t.Errorf("GetCrossrefSample(%v): error %v", tc.number, err)
		}
		if diff := cmp.Diff(tc.number, len(got)); diff != "" {
			t.Errorf("GetCrossrefSample mismatch (-want +got):\n%s", diff)
		}
	}
}
