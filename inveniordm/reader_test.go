package inveniordm_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/front-matter/commonmeta/inveniordm"
	"github.com/muesli/cache2go"
	"golang.org/x/time/rate"
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
	match := true
	for _, tc := range testCases {
		got, err := inveniordm.Fetch(tc.pid, match)
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

func ExampleSearchBySlug() {
	host := "rogue-scholar.org"
	rl := rate.NewLimiter(rate.Every(60*time.Second), 800)
	client := inveniordm.NewClient(rl, host)
	cache := cache2go.Cache("communities")
	s, _ := inveniordm.SearchBySlug("naturalSciences", "topic", client, cache)
	fmt.Println(s)
	// Output:
	// 7d3b25fd-a4a8-4155-8e76-99d6be06706a
}
