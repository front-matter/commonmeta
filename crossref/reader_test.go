package crossref_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/crossref"
	"github.com/front-matter/commonmeta/doiutils"

	"github.com/google/go-cmp/cmp"
)

func TestGet(t *testing.T) {
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
		got, err := crossref.Get(tc.id)
		if tc.want != got.Publisher {
			t.Errorf("Get Crossref(%v): want %v, got %v, error %v",
				tc.id, tc.want, got, err)
		}
	}
}

func TestFetch(t *testing.T) {
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
		got, err := crossref.Fetch(tc.id)
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
		bytes, err := os.ReadFile(filepath)
		if err != nil {
			t.Fatal(err)
		}

		var want commonmeta.Data
		err = json.Unmarshal(bytes, &want)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Fetch (%s) mismatch (-want +got):\n%s", tc.id, diff)
		}
	}
}

func TestQueryURL(t *testing.T) {
	t.Parallel()

	type testCase struct {
		number        int
		page          int
		member        string
		type_         string
		sample        bool
		year          string
		orcid         string
		ror           string
		hasORCID      bool
		hasROR        bool
		hasReferences bool
		hasRelation   bool
		hasAbstract   bool
		hasAward      bool
		hasLicense    bool
		hasArchive    bool
		want          string
	}

	testCases := []testCase{
		{want: "https://api.crossref.org/works?offset=0&order=desc&rows=10&sort=published"},
		{number: 100, page: 3, want: "https://api.crossref.org/works?offset=200&order=desc&rows=100&sort=published"},
		{sample: true, want: "https://api.crossref.org/works?order=desc&sample=10&sort=published"},
		{sample: true, number: 120, want: "https://api.crossref.org/works?order=desc&sample=100&sort=published"},
		{sample: true, member: "340", want: "https://api.crossref.org/works?filter=member%3A340&order=desc&sample=10&sort=published"},
		{sample: true, year: "2022", want: "https://api.crossref.org/works?filter=from-pub-date%3A2022-01-01%2Cuntil-pub-date%3A2022-12-31&order=desc&sample=10&sort=published"},
		{sample: true, orcid: "0000-0002-8635-8390", want: "https://api.crossref.org/works?filter=orcid%3A0000-0002-8635-8390&order=desc&sample=10&sort=published"},
		{sample: true, ror: "041kmwe10", want: "https://api.crossref.org/works?filter=ror-id%3A041kmwe10&order=desc&sample=10&sort=published"},
		{sample: true, hasORCID: true, want: "https://api.crossref.org/works?filter=has-orcid%3Atrue&order=desc&sample=10&sort=published"},
		{sample: true, hasROR: true, want: "https://api.crossref.org/works?filter=has-ror-id%3Atrue&order=desc&sample=10&sort=published"},
		{sample: true, hasReferences: true, want: "https://api.crossref.org/works?filter=has-references%3Atrue&order=desc&sample=10&sort=published"},
		{sample: true, hasRelation: true, want: "https://api.crossref.org/works?filter=has-relation%3Atrue&order=desc&sample=10&sort=published"},
		{sample: true, hasAbstract: true, want: "https://api.crossref.org/works?filter=has-abstract%3Atrue&order=desc&sample=10&sort=published"},
		{sample: true, hasAward: true, want: "https://api.crossref.org/works?filter=has-award%3Atrue&order=desc&sample=10&sort=published"},
		{sample: true, hasLicense: true, want: "https://api.crossref.org/works?filter=has-license%3Atrue&order=desc&sample=10&sort=published"},
	}
	for _, tc := range testCases {
		got := crossref.QueryURL(tc.number, tc.page, tc.member, tc.type_, tc.sample, tc.year, tc.orcid, tc.ror, tc.hasORCID, tc.hasROR, tc.hasReferences, tc.hasRelation, tc.hasAbstract, tc.hasAward, tc.hasLicense, tc.hasArchive)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("CrossrefQueryUrl mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestGetAll(t *testing.T) {
	t.Parallel()

	type testCase struct {
		number        int
		page          int
		member        string
		type_         string
		sample        bool
		year          string
		ror           string
		orcid         string
		hasORCID      bool
		hasROR        bool
		hasReferences bool
		hasRelation   bool
		hasAbstract   bool
		hasAward      bool
		hasLicense    bool
		hasArchive    bool
	}

	testCases := []testCase{
		{number: 3, member: "340", type_: "journal-article"},
		{number: 1, type_: "posted-content", sample: true},
		{number: 2, type_: "", sample: true},
	}
	for _, tc := range testCases {
		got, err := crossref.GetAll(tc.number, tc.page, tc.member, tc.type_, tc.sample, tc.year, tc.ror, tc.orcid, tc.hasORCID, tc.hasROR, tc.hasReferences, tc.hasRelation, tc.hasAbstract, tc.hasAward, tc.hasLicense, tc.hasArchive)
		if err != nil {
			t.Errorf("GetAll (%v): error %v", tc.number, err)
		}
		if diff := cmp.Diff(tc.number, len(got)); diff != "" {
			t.Errorf("GetAll mismatch (-want +got):\n%s", diff)
		}
	}
}
func TestGetMember(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "340", want: "Public Library of Science (PLoS)"},
		{input: "1313", want: ""},
		{input: "", want: ""},
	}
	for _, tc := range testCases {
		got, _ := crossref.GetMember(tc.input)
		if tc.want != got {
			t.Errorf("Get Crossref Member(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}
