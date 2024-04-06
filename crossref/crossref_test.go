package crossref_test

import (
	"commonmeta/crossref"
	"testing"
)

func TestGetCrossref(t *testing.T) {
	t.Parallel()

	type testCase struct {
		doi  string
		want string
		err  error
	}

	journalArticle := crossref.Record{
		URL:       "https://api.crossref.org/works/10.7554/elife.01567",
		DOI:       "10.7554/elife.01567",
		Publisher: "eLife Sciences Publications, Ltd",
	}
	postedContent := crossref.Record{
		URL:       "https://api.crossref.org/works/10.1101/097196",
		DOI:       "10.1101/097196",
		Publisher: "Cold Spring Harbor Laboratory",
	}

	testCases := []testCase{
		{doi: journalArticle.DOI, want: journalArticle.Publisher, err: nil},
		{doi: postedContent.DOI, want: postedContent.Publisher, err: nil},
	}
	for _, tc := range testCases {
		got, err := crossref.GetCrossref(tc.doi)
		if tc.want != got.Publisher {
			t.Errorf("Get Crossref(%v): want %v, got %v, error %v",
				tc.doi, tc.want, got, err)
		}
	}
}
