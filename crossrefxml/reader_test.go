package crossrefxml_test

import (
	"testing"

	"github.com/front-matter/commonmeta/crossrefxml"
)

func TestResource(t *testing.T) {
	t.Parallel()

	type testCase struct {
		id   string
		want string
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
		{id: journalArticle.DOI, want: journalArticle.Resource},
		{id: postedContent.DOI, want: postedContent.Resource},
	}
	for _, tc := range testCases {
		content, err := crossrefxml.Get(tc.id)
		if err != nil {
			t.Errorf("Get (%v): error %v", tc.id, err)
		}
		var got string
		if content.DOI.Type == "journal_article" {
			got = content.DOIRecord.Crossref.Journal.JournalArticle.DOIData.Resource
		} else if content.DOI.Type == "posted_content" {
			got = content.DOIRecord.Crossref.PostedContent.DOIData.Resource
		}
		if tc.want != got {
			t.Errorf("Get (%v): want %v, got %v",
				tc.id, tc.want, got)
		}
	}
}

func TestType(t *testing.T) {
	t.Parallel()

	type testCase struct {
		id   string
		want string
	}

	testCases := []testCase{
		{id: "10.5040/9781718219342", want: "Book"},
		{id: "10.4324/9780429432057-2", want: "BookChapter"},
		{id: "10.1017/9781009458351.012", want: "BookPart"},
		{id: "10.1093/oso/9780190948146.002.0004", want: "BookSection"},
		{id: "10.1039/2050-9251", want: "BookSeries"},
		{id: "10.7139/2017.978-1-56900-592-7", want: "BookSet"},
		{id: "10.11647/obp.0193.29", want: "BookTrack"},
		{id: "10.1021/acs.analchem.2c04701.s002", want: "Component"},
		{id: "10.54985/peeref.2404p5092873", want: "Database"},
		{id: "10.1037/e542602008-001", want: "Dataset"},
		{id: "10.11606/t.8.2017.tde-08052017-100442", want: "Dissertation"},
		{id: "10.51566/ceper2117_55", want: "Book"}, // edited-book
		{id: "10.46936/cpcy.proj.2019.50733/60006578", want: "Grant"},
		{id: "10.35841/immunology-case-reports", want: "Journal"},
		{id: "10.1016/0016-2361(85)90041-9", want: "JournalArticle"},
		{id: "10.24114/jnc.v2i2", want: "JournalIssue"},
		{id: "10.53738/revmed.2017.13.567", want: "JournalVolume"},
		{id: "10.3917/puf.gcosl.1996.01", want: "Book"}, // monograph
		{id: "10.1039/9781847550378-fx001", want: "Other"},
		{id: "10.1111/pim.13023/v2/decision1", want: "PeerReview"},
		{id: "10.5194/egusphere-egu24-18880", want: "Article"}, // posted-content
		{id: "10.1109/isscs39599.2017", want: "Proceedings"},
		{id: "10.1145/3448016.3452841", want: "ProceedingsArticle"},
		{id: "10.15405/epsbs(2357-1330).2021.6.1", want: "ProceedingsSeries"},
		{id: "10.4135/9781529708455", want: "Book"}, // reference-book
		{id: "10.32388/nbim0r", want: "Entry"},      // reference-entry
		{id: "10.3133/tei371", want: "Report"},
		{id: "10.32468/evo-end-ext-ban-col.excel1.tr1-2024", want: "ReportComponent"},
		{id: "10.1787/ece3a6d6-en", want: "ReportSeries"},
		{id: "10.3403/30323840u", want: "Standard"},
	}

	for _, tc := range testCases {
		content, err := crossrefxml.Get(tc.id)
		if err != nil {
			t.Errorf("Get (%v): error %v", tc.id, err)
		}
		got := content.Type()
		if tc.want != got {
			t.Errorf("Get (%v): want %v, got %v",
				tc.id, tc.want, got)
		}
	}
}
