package types_test

import (
	"testing"

	"commonmeta/types"
)

func TestMetadata(t *testing.T) {
	t.Parallel()
	_ = types.Data{
		ID:   "https://doi.org/10.7554/elife.01567",
		Type: "JournalArticle",
	}
}
