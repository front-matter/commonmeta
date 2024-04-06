package inveniordm_test

import (
	"commonmeta/inveniordm"
	"testing"
)

func TestGetInvenioRDM(t *testing.T) {
	t.Parallel()

	type testCase struct {
		id   string
		want string
		err  error
	}

	publication := inveniordm.Record{
		ID:    "5244404",
		DOI:   "10.5281/zenodo.5244404",
		Title: "The Origins of SARS-CoV-2: A Critical Review",
	}
	presentation := inveniordm.Record{
		ID:    "8173303",
		DOI:   "10.5281/zenodo.8173303",
		Title: "11 July 2023 (Day 2) CERN – NASA Open Science Summit Sketch Notes",
	}

	testCases := []testCase{
		{id: presentation.ID, want: presentation.Title, err: nil},
		{id: publication.ID, want: publication.Title, err: nil},
	}
	for _, tc := range testCases {
		got, err := inveniordm.GetInvenioRDM(tc.id)
		if tc.want != got.Title {
			t.Errorf("InvenioRDM ID(%v): want %v, got %v, error %v",
				tc.id, tc.want, got, err)
		}
	}
}
