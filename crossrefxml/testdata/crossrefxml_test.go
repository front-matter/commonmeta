package crossrefxml_test

import (
	"testing"

	"github.com/front-matter/commonmeta/crossrefxml"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type testCase struct {
		id   string
		want string
		err  error
	}

	journalArticle := crossrefxml.DOIData{
		DOI:      "10.7554/elife.01567",
		Resource: "https://elifesciences.org/articles/01567",
	}
	postedContent := crossrefxml.DOIData{
		DOI:      "10.1101/097196",
		Resource: "http://biorxiv.org/lookup/doi/10.1101/097196",
	}

	testCases := []testCase{
		{id: journalArticle.DOI, want: journalArticle.Resource, err: nil},
		{id: postedContent.DOI, want: postedContent.Resource, err: nil},
	}
	for _, tc := range testCases {
		got, err := crossrefxml.Get(tc.id)
		if err != nil {
			t.Errorf("Get (%v): error %v", tc.id, err)
		}
		var resource string
		if got.Query.DOI.Type == "journal-article" {
			resource = got.Query.DOIRecord.Crossref.Journal.JournalArticle.DOIData.Resource
		} else if got.Query.DOI.Type == "posted-content" {
			resource = got.Query.DOIRecord.Crossref.PostedContent.DOIData.Resource
		}
		if tc.want != resource {
			t.Errorf("Get (%v): want %v, got %v, error %v",
				tc.id, tc.want, got, err)
		}
	}
}
