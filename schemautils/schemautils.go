// Package schemautils provides functions to validate JSON documents against JSON Schema files.
package schemautils

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/encoding/json"
)

// JSONSchemas is the embedded JSON Schema files.
//
//go:embed schemas/*.cue
var Schemas embed.FS

const schemaVersion = "commonmeta_v0.15"

// SchemaErrors validates a JSON document against a Cue Schema file.
func SchemaErrors(document []byte, schema ...string) error {

	// If no schema is provided, default to const schema_version
	if len(schema) == 0 {
		schema = append(schema, schemaVersion)
	}
	s := schema[len(schema)-1]

	// Cue Schema files stored locally to validate against
	schemata := []string{schemaVersion, "datacite-v4.5", "crossref-v0.2", "csl-data", "cff_v1.2.0", "invenio-rdm-v0.1"}
	if !slices.Contains(schemata, s) {
		log.Fatalf("Schema %s not found", s)
	}
	dir := "schemas"
	data, err := Schemas.ReadFile(filepath.Join(dir, s+".cue"))
	if err != nil {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		fmt.Println(exPath)
		fmt.Print(err)
	}
	cueSource := string(data)
	ctx := cuecontext.New()

	// Build the schema
	cueSchema := ctx.CompileString(cueSource).LookupPath(cue.ParsePath("schema"))
	fmt.Println(cueSchema)
	// Load the JSON file specified (the program's sole argument) as a CUE value
	dataExpr, err := json.Extract("_", document)
	if err != nil {
		log.Fatal(err)
	}
	dataAsCUE := ctx.BuildExpr(dataExpr)

	// Validate the JSON data using the schema
	unified := cueSchema.Unify(dataAsCUE)
	return unified.Validate()
}
