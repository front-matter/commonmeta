package inveniordm_test

import (
	"testing"

	"github.com/front-matter/commonmeta/inveniordm"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type testCase struct {
		pid  string
		want string
		err  error
	}

	var metadata inveniordm.Metadata
	metadata.Title = "The Origins of SARS-CoV-2: A Critical Review"
	publication := inveniordm.Inveniordm{
		ID:       "https://zenodo.org/api/records/524440",
		Metadata: metadata,
	}

	metadata.Title = "11 July 2023 (Day 2) CERN – NASA Open Science Summit Sketch Notes"
	presentation := inveniordm.Inveniordm{
		ID:       "https://zenodo.org/api/records/8173303",
		Metadata: metadata,
	}

	metadata.Title = "The Origins of SARS-CoV-2: A Critical Review"
	preprint := inveniordm.Inveniordm{
		ID:       "https://rogue-scholar.org/api/records/42jxf-4yd62",
		Metadata: metadata,
	}

	testCases := []testCase{
		{pid: presentation.ID, want: presentation.Metadata.Title, err: nil},
		{pid: publication.ID, want: publication.Metadata.Title, err: nil},
		{pid: preprint.ID, want: preprint.Metadata.Title, err: nil},
	}
	for _, tc := range testCases {
		got, err := inveniordm.Get(tc.pid)
		if tc.want != got.Metadata.Title {
			t.Errorf("InvenioRDM ID(%v): want %v, got %v, error %v",
				tc.pid, tc.want, got, err)
		}
	}
}

func TestFetch(t *testing.T) {
	t.Parallel()

	type testCase struct {
		pid  string
		want string
		err  error
	}

	var metadata inveniordm.Metadata
	metadata.Title = "The Origins of SARS-CoV-2: A Critical Review"
	publication := inveniordm.Inveniordm{
		ID:       "https://zenodo.org/api/records/524440",
		Metadata: metadata,
	}

	metadata.Title = "11 July 2023 (Day 2) CERN – NASA Open Science Summit Sketch Notes"
	presentation := inveniordm.Inveniordm{
		ID:       "https://zenodo.org/api/records/8173303",
		Metadata: metadata,
	}

	metadata.Title = "The Origins of SARS-CoV-2: A Critical Review"
	preprint := inveniordm.Inveniordm{
		ID:       "https://rogue-scholar.org/api/records/42jxf-4yd62",
		Metadata: metadata,
	}

	testCases := []testCase{
		{pid: presentation.ID, want: presentation.Metadata.Title, err: nil},
		{pid: publication.ID, want: publication.Metadata.Title, err: nil},
		{pid: preprint.ID, want: preprint.Metadata.Title, err: nil},
	}
	for _, tc := range testCases {
		got, err := inveniordm.Fetch(tc.pid)
		if tc.want != got.Titles[0].Title {
			t.Errorf("InvenioRDM ID(%v): want %v, got %v, error %v",
				tc.pid, tc.want, got, err)
		}
	}
}

// func ExampleSearchByDOI() {
// 	s, _ := inveniordm.SearchByDOI("https://doi.org/10.59350/k0746-rsc44", "rogue-scholar.org")
// 	fmt.Println(s)
// 	// Output:
// 	// [xm2mv-r7378]
// }
