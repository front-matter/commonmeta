package ror_test

import (
	"fmt"

	"github.com/front-matter/commonmeta/ror"
)

func ExampleGetCountryCodes() {
	s := ror.GetCountryCodes("dk")

	fmt.Println(s)
	// Output:
	// DK
}
