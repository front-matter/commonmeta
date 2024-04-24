package commonmeta_test

import (
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
