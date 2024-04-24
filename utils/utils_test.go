package utils_test

import (
	"fmt"
	"testing"

	"github.com/front-matter/commonmeta/utils"
)

func TestNormalizeURL(t *testing.T) {
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
		got, err := utils.NormalizeURL(tc.input, tc.secure, tc.lower)
		if tc.want != got {
			t.Errorf("Normalize URL(%v): want %v, got %v, error %v",
				tc.input, tc.want, got, err)
		}
	}
}

func TestISSNAsURL(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "2146-8427", want: "https://portal.issn.org/resource/ISSN/2146-8427"},
	}
	for _, tc := range testCases {
		got := utils.ISSNAsURL(tc.input)
		if tc.want != got {
			t.Errorf("ISSN as URL(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}

func ExampleISSNAsUrl() {
	s := utils.ISSNAsURL("2146-8427")
	fmt.Println(s)
	// Output:
	// https://portal.issn.org/resource/ISSN/2146-8427
}

func TestNormalizeORCID(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "0000-0002-1825-0097", want: "https://orcid.org/0000-0002-1825-0097"},
		{input: "https://orcid.org/0000-0002-1825-0097", want: "https://orcid.org/0000-0002-1825-0097"},
	}
	for _, tc := range testCases {
		got := utils.NormalizeORCID(tc.input)
		if tc.want != got {
			t.Errorf("Normalize ORCID(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}

func TestValidateORCID(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "http://orcid.org/0000-0002-2590-225X", want: "0000-0002-2590-225X"},
		{input: "https://orcid.org/0000-0002-1825-0097", want: "0000-0002-1825-0097"},
		{input: "0000-0002-1825-0097", want: "0000-0002-1825-0097"},
		{input: "https://sandbox.orcid.org/0000-0002-1825-0097", want: "0000-0002-1825-0097"},
		{input: "0000-0002-1825-009", want: ""}, // invalid ORCID
	}
	for _, tc := range testCases {
		got, ok := utils.ValidateORCID(tc.input)
		if tc.want != got {
			t.Errorf("Validate ORCID(%v): want %v, got %v, ok %v",
				tc.input, tc.want, got, ok)
		}
	}
}

func TestNormalizeROR(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "https://ror.org/0342dzm54", want: "https://ror.org/0342dzm54"},
		{input: "http://ror.org/0342dzm54", want: "https://ror.org/0342dzm54"},
	}
	for _, tc := range testCases {
		got := utils.NormalizeROR(tc.input)
		if tc.want != got {
			t.Errorf("Normalize ORCID(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}

func TestValidateROR(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "https://ror.org/0342dzm54", want: "0342dzm54"},
	}
	for _, tc := range testCases {
		got, ok := utils.ValidateROR(tc.input)
		if tc.want != got {
			t.Errorf("Validate ROR(%v): want %v, got %v, ok %v",
				tc.input, tc.want, got, ok)
		}
	}
}

func ExampleValidateROR() {
	s, _ := utils.ValidateROR("https://ror.org/0342dzm54")
	fmt.Println(s)
	// Output:
	// 0342dzm54
}

func TestSanitize(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "<p>The Origins of SARS-CoV-2: A Critical <a href=\"index.html\">Review</a></p>", want: "The Origins of SARS-CoV-2: A Critical Review"},
		{input: "11 July 2023 (Day 2) CERN – NASA Open Science Summit Sketch Notes", want: "11 July 2023 (Day 2) CERN – NASA Open Science Summit Sketch Notes"},
	}
	for _, tc := range testCases {
		got := utils.Sanitize(tc.input)
		if tc.want != got {
			t.Errorf("Sanitize String(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}

func ExampleSanitize() {
	s := utils.Sanitize("<p>The Origins of SARS-CoV-2: A <i>Critical</i> <a href=\"index.html\">Review</a></p>")
	fmt.Println(s)
	// Output:
	// The Origins of SARS-CoV-2: A <i>Critical</i> Review
}

func ExampleUnescapeUTF8() {
	s, err := utils.UnescapeUTF8("capable of signi\"cance.")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s)
	// Output:
	// capable of signi"cance.
}
