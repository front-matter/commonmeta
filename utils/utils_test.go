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

func ExampleNormalizeID() {
	s := utils.NormalizeID("10.1101/097196")
	fmt.Println(s)
	// Output:
	// https://doi.org/10.1101/097196
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

func ExampleNormalizeURL() {
	s, _ := utils.NormalizeURL("http://elifesciences.org/Articles/91729/", true, true)
	fmt.Println(s)
	// Output:
	// https://elifesciences.org/articles/91729
}

func ExampleNormalizeCCUrl() {
	s, ok := utils.NormalizeCCUrl("https://creativecommons.org/licenses/by/4.0")
	fmt.Println(s, ok)
	// Output:
	// https://creativecommons.org/licenses/by/4.0/legalcode true
}

func ExampleURLToSPDX() {
	s := utils.URLToSPDX("https://creativecommons.org/licenses/by/4.0/legalcode")
	fmt.Println(s)
	// Output:
	// CC-BY-4.0
}

func TestFindFromFormatByID(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "10.1371/journal.pone.0042793", want: "crossref"},
		{input: "https://doi.org/10.5061/dryad.8515", want: "datacite"},
		{input: "10.1392/roma081203", want: "medra"},
		{input: "https://doi.org/10.5012/bkcs.2013.34.10.2889", want: "kisti"},
		{input: "https://doi.org/10.11367/grsj1979.12.283", want: "jalc"},
		{input: "https://doi.org/10.2903/j.efsa.2018.5239", want: "op"},
		{input: "https://github.com/citation-file-format/ruby-cff/blob/main/CITATION.cff", want: "cff"},
		{input: "https://github.com/datacite/maremma/blob/master/codemeta.json", want: "codemeta"},
		{input: "https://dataverse.harvard.edu/dataset.xhtml?persistentId=doi:10.7910/DVN/GAOC03", want: "schemaorg"},
		{input: "https://api.rogue-scholar.org/posts/c3095752-2af0-40a4-a229-3ceb7424bce2", want: "jsonfeed"},
		{input: "https://rogue-scholar.org/records/h9ckh-tke61", want: "inveniordm"},
	}
	for _, tc := range testCases {
		got := utils.FindFromFormatByID(tc.input)
		if tc.want != got {
			t.Errorf("Find from format by ID (%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}

func ExampleFindFromFormatByID() {
	s := utils.FindFromFormatByID("https://api.rogue-scholar.org/posts/4e4bf150-751f-4245-b4ca-fe69e3c3bb24")
	fmt.Println(s)
	// Output:
	// jsonfeed
}

func ExampleFindFromFormatByExt() {
	s := utils.FindFromFormatByExt(".bib")
	fmt.Println(s)
	// Output:
	// bibtex
}

func ExampleFindFromFormatByMap() {
	f := map[string]any{
		"schemaVersion": "http://datacite.org/schema/kernel-4",
	}
	s := utils.FindFromFormatByMap(f)
	fmt.Println(s)
	// Output:
	// datacite
}

// func ExampleFindFromFormatByString() {
// 	s := utils.FindFromFormatByString(".bib")
// 	fmt.Println(s)
// 	// Output:
// 	// bibtex
// }

func ExampleFindFromFormatByFilename() {
	s := utils.FindFromFormatByFilename("CITATION.cff")
	fmt.Println(s)
	// Output:
	// cff
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

func ExampleNormalizeORCID() {
	s := utils.NormalizeORCID("0000-0002-1825-0097")
	fmt.Println(s)
	// Output:
	// https://orcid.org/0000-0002-1825-0097
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

func ExampleValidateORCID() {
	s, _ := utils.ValidateORCID("https://orcid.org/0000-0002-1825-0097")
	fmt.Println(s)
	// Output:
	// 0000-0002-1825-0097
}

func TestORCIDNumberRange(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  bool
	}
	testCases := []testCase{
		{input: "0000-0002-1825-0097", want: true},
		{input: "0000-0002-2590-225X", want: true},
		{input: "0000 0001 2112 2291", want: false}, // ISNI
		{input: "0000000121122291", want: false},    // ISNI
	}
	for _, tc := range testCases {
		got := utils.CheckORCIDNumberRange(tc.input)
		if tc.want != got {
			t.Errorf("ORCID Number Range(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}

func TestValidateISNI(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
		ok    bool
	}
	testCases := []testCase{
		{input: "https://isni.org/isni/0000000121122291", want: "0000000121122291", ok: true},
		{input: "https://isni.org/isni/0000 0001 2112 2291", want: "0000000121122291", ok: true},
		{input: "https://isni.org/isni/0000-0001-2112-2291", want: "0000000121122291", ok: true},
		{input: "0000 0001 2112 2291", want: "0000000121122291", ok: true},
		{input: "0000-0001-2112-2291", want: "0000000121122291", ok: true},
		{input: "https://isni.org/isni/000000021825009", want: "", ok: false}, // invalid ISNI
	}
	for _, tc := range testCases {
		got, ok := utils.ValidateISNI(tc.input)
		if tc.want != got || tc.ok != ok {
			t.Errorf("Validate ISNI(%v): want %v, got %v, ok %v",
				tc.input, tc.want, got, ok)
		}
	}
}

func ExampleValidateISNI() {
	s, _ := utils.ValidateISNI("https://isni.org/isni/000000012146438X")
	fmt.Println(s)
	// Output:
	// 000000012146438X
}

func TestValidateWikidata(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "https://www.wikidata.org/wiki/Q7186", want: "Q7186"},     // Marie Curie (person)
		{input: "https://www.wikidata.org/wiki/Q251061", want: "Q251061"}, // Potsdam Institute for Climate Impact Research (organization)
		{input: "Q251061", want: "Q251061"},
		{input: "https://www.wikidata.org/wiki/Property:P610", want: ""}, // Wikidata property not item
	}
	for _, tc := range testCases {
		got, ok := utils.ValidateWikidata(tc.input)
		if tc.want != got {
			t.Errorf("Validate Wikidata(%v): want %v, got %v, ok %v",
				tc.input, tc.want, got, ok)
		}
	}
}

func ExampleValidateWikidata() {
	s, _ := utils.ValidateWikidata("https://www.wikidata.org/wiki/Q7186")
	fmt.Println(s)
	// Output:
	// Q7186
}

func TestValidateGRID(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "https://www.grid.ac/institutes/grid.1017.7", want: "grid.1017.7"}, // Royal Melbourne Institute of Technology University
		{input: "https://grid.ac/institutes/grid.1017.7", want: "grid.1017.7"},
		{input: "grid.1017.7", want: "grid.1017.7"},
		{input: "1017.7", want: ""}, // invalid GRID
	}
	for _, tc := range testCases {
		got, ok := utils.ValidateGRID(tc.input)
		if tc.want != got {
			t.Errorf("Validate GRID(%v): want %v, got %v, ok %v",
				tc.input, tc.want, got, ok)
		}
	}
}

func ExampleValidateGRID() {
	s, _ := utils.ValidateGRID("https://www.grid.ac/institutes/grid.1017.7")
	fmt.Println(s)
	// Output:
	// grid.1017.7
}

func ExampleValidateIDCategory() {
	_, _, s := utils.ValidateIDCategory("https://ror.org/0342dzm54")
	fmt.Println(s)
	// Output:
	// Organization
}

func TestNormalizeOrganizationID(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "https://ror.org/0342dzm54", want: "https://ror.org/0342dzm54"},
		{input: "0342dzm54", want: "https://ror.org/0342dzm54"},
		{input: "grid.1017.7", want: "https://grid.ac/institutes/grid.1017.7"},
		{input: "https://grid.ac/institutes/grid.1017.7", want: "https://grid.ac/institutes/grid.1017.7"},
		{input: "Q7186", want: "https://www.wikidata.org/wiki/Q7186"},
		{input: "https://www.wikidata.org/wiki/Q7186", want: "https://www.wikidata.org/wiki/Q7186"},
		{input: "0000000121122291", want: "https://isni.org/isni/0000000121122291"},
		{input: "0000-0001-2112-2291", want: "https://isni.org/isni/0000000121122291"},
		{input: "https://isni.org/isni/0000000121122291", want: "https://isni.org/isni/0000000121122291"},
	}
	for _, tc := range testCases {
		got := utils.NormalizeOrganizationID(tc.input)
		if tc.want != got {
			t.Errorf("Normalize Organization ID(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}

func ExampleNormalizeOrganizationID() {
	s := utils.NormalizeOrganizationID("https://ror.org/0342dzm54")
	fmt.Println(s)
	// Output:
	// https://ror.org/0342dzm54
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

func ExampleNormalizeROR() {
	s := utils.NormalizeROR("http://ror.org/0342dzm54")
	fmt.Println(s)
	// Output:
	// https://ror.org/0342dzm54
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

func ExampleGetROR() {
	ror, _ := utils.GetROR("https://ror.org/0342dzm54")
	fmt.Println(ror.Name)
	// Output:
	// Liberate Science
}

func TestValidateROR(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "https://ror.org/0342dzm54", want: "0342dzm54"},
		//TODO: {input: "ror.org/0342dzm54", want: "0342dzm54"},
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

func TestValidateCrossrefFunderID(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "https://doi.org/10.13039/501100000155", want: "501100000155"},
		{input: "10.13039/501100000155/", want: "501100000155"},
		{input: "100010540/", want: "100010540"},
	}
	for _, tc := range testCases {
		got, ok := utils.ValidateCrossrefFunderID(tc.input)
		if tc.want != got {
			t.Errorf("Validate Crossref Funder ID(%v): want %v, got %v, ok %v",
				tc.input, tc.want, got, ok)
		}
	}
}

func ExampleValidateCrossrefFunderID() {
	s, _ := utils.ValidateCrossrefFunderID("https://doi.org/10.13039/501100001659")
	fmt.Println(s)
	// Output:
	// 501100001659
}

func ExampleValidateUUID() {
	s, _ := utils.ValidateUUID("2491b2d5-7daf-486b-b78b-e5aab48064c1")
	fmt.Println(s)
	// Output:
	// 2491b2d5-7daf-486b-b78b-e5aab48064c1
}

func ExampleValidateRID() {
	s, _ := utils.ValidateRID("nryd8-14284")
	fmt.Println(s)
	// Output:
	// nryd8-14284
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
		{input: "https://api.rogue-scholar.org/posts/f5dd4c59-47ac-44de-8aac-0c0ea6583b5a", want: "JSONFEEDID"},
		{input: "https://api.rogue-scholar.org/posts/10.59350/at255-j1j24", want: "JSONFEEDID"},
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
		{input: "nryd8-14284", want: "RID"},
		{input: "https://ror.org/0342dzm54", want: "ROR"},
		{input: "https://grid.ac/institutes/grid.1017.7", want: "GRID"},
		{input: "grid.1017.7", want: "GRID"},
		{input: "https://www.wikidata.org/wiki/Q7186", want: "Wikidata"},
		{input: "Q7186", want: "Wikidata"},
		{input: "https://isni.org/isni/0000000121122291", want: "ISNI"},
		{input: "0000-0001-2112-2291", want: "ISNI"},
		{input: "https://orcid.org/0000-0002-1825-0097", want: "ORCID"},
		{input: "https://datadryad.org/stash/dataset/doi:10.5061/dryad.8515", want: "URL"},
		// {input: "https://archive.softwareheritage.org/swh:1:dir:44641d8369477d44432fdf50b2eae38e5d079742;origin=https://github.com/murrayds/sci-text-disagreement;visit=swh:1:snp:5695398f6bd0811d67792e16a2684052abe9dc37;anchor=swh:1:rev:b361157a9cfeb536ca255422280e154855b4e9a3", want: "URL"},
		{input: "https://portal.issn.org/resource/ISSN/1094-4087", want: "ISSN"},
		{input: "2749-9952", want: "ISSN"},
		{input: "https://api.rogue-scholar.org/posts/f5dd4c59-47ac-44de-8aac-0c0ea6583b5a", want: "JSONFEEDID"},
		{input: "https://api.rogue-scholar.org/posts/10.59350/at255-j1j24", want: "JSONFEEDID"},
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

func ExampleValidateID() {
	s, t := utils.ValidateID("https://ror.org/0342dzm54")
	fmt.Println(s, t)
	// Output:
	// 0342dzm54 ROR
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

func ExampleDecodeID() {
	s, _ := utils.DecodeID("https://doi.org/10.59350/b8pcg-q9k70")
	fmt.Println(s)
	// Output:
	// 387298385203
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

func ExampleWordsToCamelCase() {
	s := utils.WordsToCamelCase("Earth and related environmental sciences")
	fmt.Println(s)
	// Output:
	// earthAndRelatedEnvironmentalSciences
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

func ExampleStringToSlug() {
	s := utils.StringToSlug("Investigación-Digital💿")
	fmt.Println(s)
	// Output:
	// investigaciondigital
}

func ExampleNormalizeString() {
	s, _ := utils.NormalizeString("Investigación-Digital💿")
	fmt.Println(s)
	// Output:
	// Investigacion-Digital💿
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

func ExampleSplitString() {
	s := utils.SplitString("0000000121122291", 4, " ")
	fmt.Println(s)
	// Output:
	// 0000 0001 2112 2291
}
