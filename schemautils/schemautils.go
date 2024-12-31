// Package schemautils provides functions to validate JSON documents against JSON Schema files.
package schemautils

import (
	"bytes"
	"embed"
	"log"
	"path/filepath"
	"slices"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// Schemas is the embedded JSON Schema files.
//
//go:embed schemas/*.json
var Schemas embed.FS

const schemaVersion = "commonmeta_v0.15"

// SchemaErrors validates a JSON document against a JSON Schema.
func SchemaErrors(document []byte, schema ...string) error {

	// If no schema is provided, default to const schema_version
	if len(schema) == 0 {
		schema = append(schema, schemaVersion)
	}
	s := schema[len(schema)-1]

	// JSON Schema files stored locally to validate against
	schemata := []string{schemaVersion, "datacite-v4.5", "crossref-v0.2", "csl-data", "cff_v1.2.0", "invenio-rdm-v0.1"}
	if !slices.Contains(schemata, s) {
		log.Fatalf("Schema %s not found", s)
	}
	dir := "schemas"
	schemaPath := filepath.Join(dir, s+".json")

	c := jsonschema.NewCompiler()
	sch, err := c.Compile(schemaPath)
	if err != nil {
		return err
	}

	// Load the document to validate
	inst, err := jsonschema.UnmarshalJSON(bytes.NewReader(document))
	if err != nil {
		return err
	}
	return sch.Validate(inst)
}
