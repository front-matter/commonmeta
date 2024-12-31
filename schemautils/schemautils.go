// Package schemautils provides functions to validate JSON documents against JSON Schema files.
package schemautils

import (
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"slices"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/encoding/json"
	"cuelang.org/go/encoding/jsonschema"
)

// Schemas is the embedded JSON Schema files.
//
//go:embed schemas/*.json
var Schemas embed.FS

const schemaVersion = "commonmeta_v0.15"

// SchemaErrors validates a JSON document against a JSON Schema using Cue.
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
	schemaFile, err := Schemas.ReadFile(schemaPath)
	if err != nil {
		fmt.Print(err)
	}

	ctx := cuecontext.New()
	schemaJsonAst, err := json.Extract(schemaPath, schemaFile)
	if err != nil {
		log.Fatal(err)
	}
	schemaJson := ctx.BuildExpr(schemaJsonAst)

	// Extract JSON Schema from the JSON
	schemaAst, err := jsonschema.Extract(schemaJson, &jsonschema.Config{
		Strict: false,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Build a cue.Value of the schema
	cueSchema := ctx.BuildFile(schemaAst)

	// Load the data JSON
	dataAst, err := json.Extract(".", document)
	if err != nil {
		log.Fatal(err)
	}

	// Build a cue.Value of the data
	cueData := ctx.BuildExpr(dataAst)

	// Unify the schema and data
	res := cueSchema.Unify(cueData)

	// Validate whether the combined (unified) result has errors or not.
	err = res.Validate(cue.Concrete(true))
	return err
}
