package datacite_test

import (
	"commonmeta/datacite"
	"testing"
)

func TestGetDataCite(t *testing.T) {
	t.Parallel()

	type testCase struct {
		id   string
		want string
		err  error
	}

	publication := datacite.Record{
		Attributes: datacite.Attributes{
			DOI: "https://doi.org/10.5281/zenodo.5244404",
			Url: "https://zenodo.org/record/5244404",
		},
	}
	presentation := datacite.Record{
		Attributes: datacite.Attributes{
			DOI: "10.5281/zenodo.8173303",
			Url: "https://zenodo.org/record/8173303",
		},
	}

	testCases := []testCase{
		{id: presentation.Attributes.DOI, want: presentation.Attributes.Url, err: nil},
		{id: publication.Attributes.DOI, want: publication.Attributes.Url, err: nil},
	}
	for _, tc := range testCases {
		got, err := datacite.GetDatacite(tc.id)
		if tc.want != got.Attributes.Url {
			t.Errorf("InvenioRDM ID(%v): want %v, got %v, error %v",
				tc.id, tc.want, got, err)
		}
	}
}
