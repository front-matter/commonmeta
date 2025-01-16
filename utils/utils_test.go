package utils_test

import (
	"fmt"
	"testing"

	"github.com/front-matter/commonmeta/utils"
)

func TestNormalizeID(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "https://doi.org/10.7554/eLife.01567", want: "https://doi.org/10.7554/elife.01567"},
		{input: "10.1101/097196", want: "https://doi.org/10.1101/097196"},
		{input: "https://datadryad.org/stash/dataset/doi:10.5061/dryad.8515", want: "https://datadryad.org/stash/dataset/doi:10.5061/dryad.8515"},
		{input: "2491b2d5-7daf-486b-b78b-e5aab48064c1", want: "2491b2d5-7daf-486b-b78b-e5aab48064c1"},
		{input: "dryad.8515", want: ""},
	}
	for _, tc := range testCases {
		got := utils.NormalizeID(tc.input)
		if tc.want != got {
			t.Errorf("Normalize ID(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}

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
		{input: "https://www.ch.ic.ac.uk/rzepa/blog/?p=27133", secure: true, lower: true, want: "https://www.ch.ic.ac.uk/rzepa/blog/?p=27133"},
		{input: "https://infomgnt.org/posts/2024-07-01-Vorstellung-OA-Datenpraxis/", secure: true, lower: false, want: "https://infomgnt.org/posts/2024-07-01-Vorstellung-OA-Datenpraxis"},
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

func ExampleISSNAsURL() {
	s := utils.ISSNAsURL("2146-8427")
	fmt.Println(s)
	// Output:
	// https://portal.issn.org/resource/ISSN/2146-8427
}

func ExampleValidateISSN() {
	s, _ := utils.ValidateISSN("https://portal.issn.org/resource/ISSN/2146-8427")
	fmt.Println(s)
	// Output:
	// 2146-8427
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
		{input: "0009-0002-1825-0097", want: "0009-0002-1825-0097"},
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

func TestGetROR(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input   string
		name    string
		fundref string
	}
	testCases := []testCase{
		{input: "https://ror.org/021nxhr62", name: "U.S. National Science Foundation", fundref: "100000001"},
		{input: "https://ror.org/018mejw64", name: "Deutsche Forschungsgemeinschaft", fundref: "501100001659"},
	}
	for _, tc := range testCases {
		got, _ := utils.GetROR(tc.input)
		if tc.name != got.Name || tc.fundref != got.ExternalIds.FundRef.All[0] {
			t.Errorf("Get ROR (%v): want %v %v, got %v %v",
				tc.input, tc.name, tc.fundref, got.Name, got.ExternalIds.FundRef.All[0])
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
		//{input: "ror.org/0342dzm54", want: "0342dzm54"},
		{input: "0342dzm54", want: "0342dzm54"},
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

func ExampleValidateUUID() {
	s, _ := utils.ValidateUUID("2491b2d5-7daf-486b-b78b-e5aab48064c1")
	fmt.Println(s)
	// Output:
	// 2491b2d5-7daf-486b-b78b-e5aab48064c1
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

func TestValidateURL(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "https://elifesciences.org/articles/91729", want: "URL"},
		{input: "https://doi.org/10.7554/eLife.91729.3", want: "DOI"},
		{input: "10.7554/eLife.91729.3", want: "DOI"},
		{input: "https://doi.org/10.1101", want: "URL"},
		{input: "10.1101", want: ""},
	}
	for _, tc := range testCases {
		got := utils.ValidateURL(tc.input)
		if tc.want != got {
			t.Errorf("Validate URL(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}

func TestValidateID(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "https://doi.org/10.7554/eLife.01567", want: "DOI"},
		{input: "10.1101/097196", want: "DOI"},
		{input: "2491b2d5-7daf-486b-b78b-e5aab48064c1", want: "UUID"},
		{input: "https://ror.org/0342dzm54", want: "ROR"},
		{input: "https://orcid.org/0000-0002-1825-0097", want: "ORCID"},
		{input: "https://datadryad.org/stash/dataset/doi:10.5061/dryad.8515", want: "URL"},
		// {input: "https://archive.softwareheritage.org/swh:1:dir:44641d8369477d44432fdf50b2eae38e5d079742;origin=https://github.com/murrayds/sci-text-disagreement;visit=swh:1:snp:5695398f6bd0811d67792e16a2684052abe9dc37;anchor=swh:1:rev:b361157a9cfeb536ca255422280e154855b4e9a3", want: "URL"},
		{input: "https://portal.issn.org/resource/ISSN/1094-4087", want: "ISSN"},
		{input: "2749-9952", want: "ISSN"},
		{input: "dryad.8515", want: ""},
	}
	for _, tc := range testCases {
		got, type_ := utils.ValidateID(tc.input)
		if tc.want != type_ {
			t.Errorf("Validate ID(%v): want %v, got %v (%v)",
				tc.input, tc.want, got, type_)
		}
	}
}

func TestDecodeID(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  int64
	}
	testCases := []testCase{
		{input: "https://doi.org/10.59350/b8pcg-q9k70", want: 387298385203},
		{input: "10.5555/ka4bq-90315", want: 663718962179},
		{input: "https://doi.org/10.7554/elife.01567", want: 0},
		{input: "10.1101/097196", want: 0},
		{input: "https://ror.org/0342dzm54", want: 104937460},
		{input: "https://orcid.org/0000-0003-1419-2405", want: 31419240},
	}
	for _, tc := range testCases {
		got, _ := utils.DecodeID(tc.input)
		if tc.want != got {
			t.Errorf("Decode ID(%v): want %v, got %v",
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

func ExampleCamelCaseToWords() {
	s := utils.CamelCaseToWords("earthAndRelatedEnvironmentalSciences")
	fmt.Println(s)
	// Output:
	// Earth and related environmental sciences
}

func ExampleKebabCaseToCamelCase() {
	s := utils.KebabCaseToCamelCase("earth-and-related-environmental-sciences")
	fmt.Println(s)
	// Output:
	// earthAndRelatedEnvironmentalSciences
}

func ExampleKebabCaseToPascalCase() {
	s := utils.KebabCaseToPascalCase("earth-and-related-environmental-sciences")
	fmt.Println(s)
	// Output:
	// EarthAndRelatedEnvironmentalSciences
}

func ExampleGetLanguage() {
	i := utils.GetLanguage("de", "iso639-3")
	fmt.Println(i)
	// Output:
	// deu
}

func ExampleCommunitySlugAsURL() {
	s := utils.CommunitySlugAsURL("irights", "rogue-scholar.org")
	fmt.Println(s)
	// Output:
	// https://rogue-scholar.org/api/communities/irights
}
