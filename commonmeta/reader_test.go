package commonmeta_test

import (
	"fmt"
	"testing"

	"github.com/front-matter/commonmeta/commonmeta"
)

func TestData(t *testing.T) {
	t.Parallel()
	_ = commonmeta.Data{
		ID:   "https://doi.org/10.7554/elife.01567",
		Type: "JournalArticle",
	}
}

func ExamplePages() {
	book := commonmeta.Container{
		Type:           "Book",
		Identifier:     "9783662463703",
		IdentifierType: "ISBN",
		Title:          "Shoulder Stiffness",
		FirstPage:      "155",
		LastPage:       "158",
	}
	s := book.Pages()
	fmt.Println(s)
	// Output:
	// 155-158
}
