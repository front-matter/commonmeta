package inveniordm_test

import (
	"fmt"
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

	publication := inveniordm.Content{
		ID:    "https://zenodo.org/api/records/5244404",
		Title: "The Origins of SARS-CoV-2: A Critical Review",
	}
	presentation := inveniordm.Content{
		ID:    "https://zenodo.org/api/records/8173303",
		Title: "11 July 2023 (Day 2) CERN – NASA Open Science Summit Sketch Notes",
	}
	preprint := inveniordm.Content{
		ID:    "https://rogue-scholar.org/api/records/42jxf-4yd62",
		Title: "The Origins of SARS-CoV-2: A Critical Review",
	}

	testCases := []testCase{
		{pid: presentation.ID, want: presentation.Title, err: nil},
		{pid: publication.ID, want: publication.Title, err: nil},
		{pid: preprint.ID, want: preprint.Title, err: nil},
	}
	for _, tc := range testCases {
		got, err := inveniordm.Get(tc.pid)
		if tc.want != got.Title {
			t.Errorf("InvenioRDM ID(%v): want %v, got %v, error %v",
				tc.pid, tc.want, got, err)
		}
	}
}

func ExampleSearchByDOI() {
	s, _ := inveniordm.SearchByDOI("https://doi.org/10.59350/k0746-rsc44", "rogue-scholar.org")
	fmt.Println(s)
	// Output:
	// [xm2mv-r7378]
}
