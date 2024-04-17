package data_test

import (
	"commonmeta/data"
	"testing"
)

// func TestNewMetadata(t *testing.T) {
// 	t.Parallel()
// 	type testCase struct {
// 		str  string
// 		want string
// 	}
// 	testCases := []testCase{
// 		{str: "https://doi.org/10.5281/zenodo.8173303", want: "https://portal.issn.org/resource/ISSN/2146-8427"},
// 	}
// 	for _, tc := range testCases {
// 		got := data.NewMetadata(tc.str)
// 		if tc.want != got {
// 			t.Errorf("Metadata(%v): want %v, got %v",
// 				tc.str, tc.want, got)
// 		}
// 	}
// }

func TestGetMetadata(t *testing.T) {
	t.Parallel()
	type testCase struct {
		str  string
		want string
	}
	testCases := []testCase{
		{str: "https://doi.org/10.5281/zenodo.8173303", want: "https://portal.issn.org/resource/ISSN/2146-8427"},
	}
	for _, tc := range testCases {
		got := data.GetMetadata(tc.str)
		if tc.want != got {
			t.Errorf("Get Metadata(%v): want %v, got %v",
				tc.str, tc.want, got)
		}
	}
}
