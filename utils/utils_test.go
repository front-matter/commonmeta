package utils_test

import (
	"testing"
	"utils"
)

func TestIssnAsUrl(t *testing.T) {
	t.Parallel()
	type testCase struct {
		issn string
		want string
	}
	testCases := []testCase{
		{issn: "2146-8427", want: "https://portal.issn.org/resource/ISSN/2146-8427"},
		{issn: nil, want: nil},
	}
	for _, tc := range testCases {
		got := utils.IssnAsUrl(tc.issn)
		if tc.want != got {
			t.Errorf("ISSN as URL(%f): want %f, got %f",
				tc.issn, tc.want, got)
		}
	}
}
