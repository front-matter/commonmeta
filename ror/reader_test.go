package ror_test

import (
	"fmt"

	"github.com/front-matter/commonmeta/ror"
)

func ExampleFetch() {
	ror, _ := ror.Fetch("https://doi.org/10.13039/501100000780")
	s := ror.ID
	fmt.Println(s)
	// Output:
	// https://ror.org/00k4n6c32
}
