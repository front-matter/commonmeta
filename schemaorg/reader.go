package schemaorg

import (
	"encoding/json"

	"github.com/front-matter/commonmeta/commonmeta"
)

// Content represents the SchemaOrg metadata returned from SchemaOrg sources. The type is more
// flexible than the SchemaOrg type, allowing for different formats of some metadata.
// Identifier can be string or []string.
type Content struct {
	*SchemaOrg
	Identifier json.RawMessage `json:"identifier"`
}

// Read reads Schema.org metadata and converts it to commonmeta.
func Read(content Content) (commonmeta.Data, error) {
	var data commonmeta.Data
	data.ID = content.ID
	return data, nil
}
