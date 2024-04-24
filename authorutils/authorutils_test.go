package authorutils_test

import (
	"testing"

	"github.com/front-matter/commonmeta/authorutils"
)

func TestIsPersonalName(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  bool
	}
	testCases := []testCase{
		{input: "John Doe", want: false},
		{input: "Harvard University", want: false},
		{input: "LiberateScience", want: false},
		{input: "Jane Smith, MD", want: true},
		{input: "John", want: false},
	}
	for _, tc := range testCases {
		got := authorutils.IsPersonalName(tc.input)
		if tc.want != got {
			t.Errorf("Is Personal Name(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}
