package schemautils_test

import (
	"commonmeta/metadata"
	"commonmeta/schemautils"
	"testing"
)

func TestJSONSchemaErrors(t *testing.T) {
	t.Parallel()
	m := metadata.Metadata{
		ID:   "https://doi.org/10.7554/elife.01567",
		Type: "JournalArticle",
		Url:  "https://elifesciences.org/articles/01567",
	}
	want := 0
	result := schemautils.JSONSchemaErrors(m)
	got := len(result.Errors())
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
