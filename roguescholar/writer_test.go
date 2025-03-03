package roguescholar_test

import (
	"fmt"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/roguescholar"
)

func ExampleUpdateLegacyRecord() {
	record := commonmeta.APIResponse{
		ID: "https://doi.org/10.7554/elife.01567",
	}
	_, err := roguescholar.UpdateLegacyRecord(record, "", "doi")
	fmt.Println(err)
	// Output:
	// no legacy key provided
}
