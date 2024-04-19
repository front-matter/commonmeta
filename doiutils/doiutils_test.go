package doiutils_test

import (
	"commonmeta/doiutils"
	"testing"
)

func TestValidateDOI(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "10.7554/elife.01567", want: "10.7554/elife.01567"},
		{input: "https://doi.org/10.7554/elife.01567", want: "10.7554/elife.01567"},
		{input: "https://doi.org/10.7554", want: ""},
		{input: "10.7554", want: ""},
		{input: "", want: ""},
	}
	for _, tc := range testCases {
		got, ok := doiutils.ValidateDOI(tc.input)
		if tc.want != got {
			t.Errorf("Validate DOI(%v): want %v, got %v, ok %v",
				tc.input, tc.want, got, ok)
		}
	}
}

func TestValidatePrefix(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "10.7554/elife.01567", want: "10.7554"},
		{input: "https://doi.org/10.7554/elife.01567", want: "10.7554"},
		{input: "https://doi.org/10.7554", want: ""},
		{input: "10.7554", want: ""},
		{input: "", want: ""},
	}
	for _, tc := range testCases {
		got, ok := doiutils.ValidatePrefix(tc.input)
		if tc.want != got {
			t.Errorf("Validate DOI(%v): want %v, got %v, ok %v",
				tc.input, tc.want, got, ok)
		}
	}
}

func TestNormalizeDOI(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "10.7554/elife.01567", want: "https://doi.org/10.7554/elife.01567"},
		{input: "https://doi.org/10.7554/elife.01567", want: "https://doi.org/10.7554/elife.01567"},
		{input: "https://doi.org/10.7554", want: ""},
		{input: "10.7554", want: ""},
		{input: "", want: ""},
	}
	for _, tc := range testCases {
		got := doiutils.NormalizeDOI(tc.input)
		if tc.want != got {
			t.Errorf("Normalize DOI(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}