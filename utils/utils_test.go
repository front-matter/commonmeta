package utils_test

import (
	"commonmeta/utils"
	"testing"
)

func TestNormalizeUrl(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input  string
		secure bool
		lower  bool
		want   string
	}
	testCases := []testCase{
		{input: "http://elifesciences.org/articles/91729/", secure: true, lower: true, want: "https://elifesciences.org/articles/91729"},
		{input: "https://api.crossref.org/works/10.1101/097196", secure: true, lower: true, want: "https://api.crossref.org/works/10.1101/097196"},
		{input: "http://elifesciences.org/articles/91729/", secure: false, lower: true, want: "http://elifesciences.org/articles/91729"},
		{input: "https://elifesciences.org/Articles/91729/", secure: true, lower: false, want: "https://elifesciences.org/Articles/91729"},
		{input: "http://elifesciences.org/Articles/91729/", secure: false, lower: false, want: "http://elifesciences.org/Articles/91729"},
	}
	for _, tc := range testCases {
		got, err := utils.NormalizeUrl(tc.input, tc.secure, tc.lower)
		if tc.want != got {
			t.Errorf("Normalize URL(%v): want %v, got %v, error %v",
				tc.input, tc.want, got, err)
		}
	}
}

func TestIssnAsUrl(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "2146-8427", want: "https://portal.issn.org/resource/ISSN/2146-8427"},
	}
	for _, tc := range testCases {
		got := utils.IssnAsUrl(tc.input)
		if tc.want != got {
			t.Errorf("ISSN as URL(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}
