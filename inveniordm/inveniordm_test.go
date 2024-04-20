package inveniordm_test

import (
	"commonmeta/inveniordm"
	"testing"
)

func TestGetInvenioRDM(t *testing.T) {
	t.Parallel()

	type testCase struct {
		pid  string
		want string
		err  error
	}

	publication := inveniordm.Content{
		ID:    "5244404",
		Title: "The Origins of SARS-CoV-2: A Critical Review",
	}
	presentation := inveniordm.Content{
		ID:    "8173303",
		Title: "11 July 2023 (Day 2) CERN – NASA Open Science Summit Sketch Notes",
	}

	testCases := []testCase{
		{pid: presentation.ID, want: presentation.Title, err: nil},
		{pid: publication.ID, want: publication.Title, err: nil},
	}
	for _, tc := range testCases {
		got, err := inveniordm.GetInvenioRDM(tc.pid)
		if tc.want != got.Title {
			t.Errorf("InvenioRDM ID(%v): want %v, got %v, error %v",
				tc.pid, tc.want, got, err)
		}
	}
}
