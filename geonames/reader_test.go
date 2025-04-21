package geonames_test

import (
	"fmt"

	"github.com/front-matter/commonmeta/geonames"
)

func ExampleLoadGeonamesCountries() {
	countries, _ := geonames.LoadGeonamesCountries()
	s := countries["DK"]
	fmt.Println(s.Name)
	// Output: Denmark
}

func ExampleLoadGeonamesCities() {
	cities, _ := geonames.LoadGeonamesCities()
	s := cities[2950159]
	fmt.Println(s.Name)
	// Output: Berlin
}
