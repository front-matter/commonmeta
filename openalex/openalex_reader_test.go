package openalex_test

import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/doiutils"
	"github.com/front-matter/commonmeta/openalex"

	"github.com/google/go-cmp/cmp"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type testCase struct {
		id   string
		want string
		err  error
	}

	journalArticle := openalex.Work{
		ID:    "https://doi.org/10.7554/elife.01567",
		Title: "",
	}
	postedContent := openalex.Work{
		ID:    "https://doi.org/10.1101/097196",
		Title: "Cold Spring Harbor Laboratory",
	}

	testCases := []testCase{
		{id: journalArticle.ID, want: journalArticle.Title, err: nil},
		{id: postedContent.ID, want: postedContent.Title, err: nil},
	}
	r := openalex.NewReader("info@front-matter.io")
	for _, tc := range testCases {
		got, err := r.Get(tc.id)
		if tc.want != got.Title {
			t.Errorf("Get OpenAlex(%v): want %v, got %v, error %v",
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
		// {name: "book", id: "https://doi.org/10.1017/9781108348843"},
		// {name: "book chapter", id: "https://doi.org/10.1007/978-3-662-46370-3_13"},
		// {name: "proceedings article", id: "https://doi.org/10.1145/3448016.3452841"},
		{name: "peer review", id: "https://doi.org/10.7554/elife.55167.sa2"},
		{name: "blog post", id: "https://doi.org/10.59350/2shz7-ehx26"},
		{name: "dissertation", id: "https://doi.org/10.14264/uql.2020.791"},
		// {name: "with ror id", id: "https://doi.org/10.1364/oe.490112"},
		// {name: "archived", id: "10.5694/j.1326-5377.1943.tb44329.x"},
	}
	r := openalex.NewReader("info@front-matter.io")
	for _, tc := range testCases {
		got, err := r.Fetch(tc.id)
		if err != nil {
			t.Errorf("OpenAlex Metadata(%v): error %v", tc.id, err)
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
		publisher     string
		type_         string
		sample        bool
		ids           string
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
		{want: "https://api.openalex.org/works?page=1&per-page=10&sort=publication_date%3Adesc"},
		{number: 100, page: 3, want: "https://api.openalex.org/works?page=3&per-page=100&sort=publication_date%3Adesc"},
		{sample: true, want: "https://api.openalex.org/works?sample=10"},
		{sample: true, number: 120, want: "https://api.openalex.org/works?sample=120"},
		// {sample: true, year: "2022", want: "https://api.openalex.org/works?filter=from-pub-date%3A2022-01-01%2Cuntil-pub-date%3A2022-12-31&order=desc&sample=10"},
		{sample: true, orcid: "0000-0002-8635-8390", want: "https://api.openalex.org/works?filter=authorships.author.id%3A0000-0002-8635-8390&sample=10"},
		{sample: true, ror: "041kmwe10", want: "https://api.openalex.org/works?filter=authorships.institutions.ror%3A041kmwe10&sample=10"},
		// {sample: true, hasORCID: true, want: "https://api.openalex.org/works?filter=has-orcid%3Atrue&order=desc&sample=10"},
		// {sample: true, hasROR: true, want: "https://api.openalex.org/works?filter=has-ror-id%3Atrue&order=desc&sample=10"},
		// {sample: true, hasReferences: true, want: "https://api.openalex.org/works?filter=has-references%3Atrue&order=desc&sample=10"},
		// {sample: true, hasRelation: true, want: "https://api.openalex.org/works?filter=has-relation%3Atrue&order=desc&sample=10"},
		// {sample: true, hasAbstract: true, want: "https://api.openalex.org/works?filter=has-abstract%3Atrue&order=desc&sample=10"},
		// {sample: true, hasAward: true, want: "https://api.openalex.org/works?filter=has-award%3Atrue&order=desc&sample=10"},
		// {sample: true, hasLicense: true, want: "https://api.openalex.org/works?filter=has-license%3Atrue&order=desc&sample=10"},
	}
	r := openalex.NewReader("info@front-matter.io")
	query := url.Values{}
	for _, tc := range testCases {
		if tc.sample {
			query.Set("sample", "10")
		} else if tc.number > 0 {
			query.Set("number", strconv.Itoa(tc.number))
		}
		got := r.QueryURL(tc.number, tc.page, tc.publisher, tc.type_, tc.sample, tc.ids, tc.year, tc.orcid, tc.ror, tc.hasORCID, tc.hasROR, tc.hasReferences, tc.hasRelation, tc.hasAbstract, tc.hasAward, tc.hasLicense, tc.hasArchive)
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
		publisher     string
		type_         string
		sample        bool
		ids           string
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
		{number: 3, type_: "journal-article"},
		{number: 1, type_: "posted-content", sample: true},
		{number: 2, type_: "", sample: true},
	}
	r := openalex.NewReader("info@front-matter.io")
	query := url.Values{}
	for _, tc := range testCases {
		if tc.sample {
			query.Set("sample", "10")
		} else if tc.number > 0 {
			query.Set("number", strconv.Itoa(tc.number))
		}
		got, err := r.GetAll(tc.number, tc.page, tc.publisher, tc.type_, tc.sample, tc.ids, tc.year, tc.orcid, tc.ror, tc.hasORCID, tc.hasROR, tc.hasReferences, tc.hasRelation, tc.hasAbstract, tc.hasAward, tc.hasLicense, tc.hasArchive)
		if err != nil {
			t.Errorf("GetAll (%v): error %v", tc.sample, err)
		}
		if diff := cmp.Diff(tc.number, len(got)); diff != "" {
			t.Errorf("GetAll mismatch (-want +got):\n%s", diff)
		}
	}
}
