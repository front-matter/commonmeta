package spdx_test

import (
	"fmt"
	"testing"

	"github.com/front-matter/commonmeta/spdx"
)

func ExampleLoadBuiltin() {
	licenses, _ := spdx.LoadBuiltin()
	s := licenses[65].LicenseID
	fmt.Println(s)
	// Output:
	// BSD-1-Clause
}

func TestSearch(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "CC-BY-4.0", want: "Creative Commons Attribution 4.0 International"},
		{input: "MIT", want: "MIT License"},
		{input: "https://creativecommons.org/licenses/by/4.0/legalcode", want: "Creative Commons Attribution 4.0 International"},
		{input: "BSL", want: ""}, // BSL is not a valid SPDX license ID
	}
	for _, tc := range testCases {
		license, err := spdx.Search(tc.input)
		if err != nil {
			t.Errorf("Error searching for license: %v", err)
		}
		if license.Name != tc.want {
			t.Errorf("Expected '%s', got '%s'", tc.want, license.Name)
		}
	}
}

func ExampleSearch() {
	license, _ := spdx.Search("CC-BY-4.0")
	s := license.Name
	fmt.Println(s)
	// Output:
	// Creative Commons Attribution 4.0 International
}
