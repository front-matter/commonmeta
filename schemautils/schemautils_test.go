package schemautils_test

import (
	"commonmeta/metadata"
	"commonmeta/schemautils"
	"encoding/json"
	"log"
	"os"
	"testing"
)

func TestJSONSchemaErrors(t *testing.T) {
	t.Parallel()
	m := metadata.Metadata{
		ID:   "https://doi.org/10.7554/elife.01567",
		Type: "JournalArticle",
	}
	doc, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	want := 0
	result := schemautils.JSONSchemaErrors(doc)
	got := len(result.Errors())
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
