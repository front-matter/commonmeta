package doiutils_test

import (
	"fmt"
	"testing"

	"github.com/front-matter/commonmeta/doiutils"
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
		{input: "10.3201/eid1503.081203 10.1083/jcb.1843iti1", want: ""},
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
		{input: "https://doi.org/10.7554", want: "10.7554"},
		{input: "10.7554", want: "10.7554"},
		{input: "", want: ""},
	}
	for _, tc := range testCases {
		got, ok := doiutils.ValidatePrefix(tc.input)
		if tc.want != got {
			t.Errorf("Validate Prefix (%v): want %v, got %v, ok %v",
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

func TestGetDOIRA(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "10.7554/elife.01567", want: "Crossref"},
		{input: "https://doi.org/10.5061/dryad.8515", want: "DataCite"},
		{input: "10.9999", want: ""},
		{input: "10.7554", want: "Crossref"},
		{input: "10.83132", want: "DataCite"},
		{input: "", want: ""},
	}
	for _, tc := range testCases {
		got, _ := doiutils.GetDOIRA(tc.input)
		if tc.want != got {
			t.Errorf("Get DOI RA(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}

func TestIsRogueScholarDOI(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		ra    string
		want  bool
	}
	testCases := []testCase{
		{input: "https://doi.org/10.59350/dybzk-cs537", ra: "crossref", want: true},
		{input: "10.59350/dybzk-cs537", ra: "crossref", want: true},
		{input: "https://doi.org/10.1101/097196", want: false},
		{input: "123456", want: false},
	}
	for _, tc := range testCases {
		got := doiutils.IsRogueScholarDOI(tc.input, tc.ra)
		if tc.want != got {
			t.Errorf("Is Rogue Scholar Prefix(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}

func ExampleEscapeDOI() {
	s := doiutils.EscapeDOI("10.59350/k0746-rsc44")
	fmt.Println(s)
	// Output:
	// 10.59350%2Fk0746-rsc44
}
