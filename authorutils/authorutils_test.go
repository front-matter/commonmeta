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
		{input: "John Doe", want: true},
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

func TestParseName(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input      string
		givenName  string
		familyName string
		name       string
	}
	testCases := []testCase{
		{input: "John Doe", givenName: "John", familyName: "Doe", name: ""},
		{input: "Rainer Maria Rilke", givenName: "Rainer Maria", familyName: "Rilke", name: ""},
		{input: "Harvard University", givenName: "", familyName: "", name: "Harvard University"},
		{input: "LiberateScience", givenName: "", familyName: "", name: "LiberateScience"},
		{input: "Jane Smith, MD", givenName: "Jane", familyName: "Smith", name: ""},
		{input: "John", givenName: "", familyName: "", name: "John"},
	}
	for _, tc := range testCases {
		givenName, familyName, name := authorutils.ParseName(tc.input)
		if tc.givenName != givenName || tc.familyName != familyName || tc.name != name {
			t.Errorf("Parse Name(%v): want (%v %v) - %v, got (%v %v) - %v",
				tc.input, tc.givenName, tc.familyName, tc.name, givenName, familyName, name)
		}
	}
}
