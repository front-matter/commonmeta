package schemaorg

import "commonmeta/metadata"

type Record struct {
	ID string `json:"id"`
}

func GetSchemaOrg(pid string) string {
	return "Schema.org"
}

func ReadSchemaorg(record Record) (metadata.Metadata, error) {
	var m metadata.Metadata
	return m, nil
}
