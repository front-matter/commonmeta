package metadata_test

import (
	"metadata/metadata"
	"testing"
)

func TestMetadata(t *testing.T) {
	t.Parallel()
	_ = metadata.Metadata{
		ID:   "https://doi.org/10.7554/elife.01567",
		Type: "JournalArticle",
	}
}
