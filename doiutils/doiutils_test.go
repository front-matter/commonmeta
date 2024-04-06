package doiutils_test

import (
	"commonmeta/doiutils"
	"testing"
)

func TestDOIFromUrl(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "https://doi.org/10.7554/elife.01567", want: "10.7554/elife.01567"},
		{input: "10.1101/097196", want: "10.1101/097196"},
	}
	for _, tc := range testCases {
		got, err := doiutils.DOIFromUrl(tc.input)
		if tc.want != got {
			t.Errorf("DOI from Url(%v): want %v, got %v, error %v",
				tc.input, tc.want, got, err)
		}
	}
}
