// Package schemautils provides functions to validate JSON documents against JSON Schema files.
package schemautils

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/xeipuuv/gojsonschema"
)

// JSONSchemas is the embedded JSON Schema files.
//
//go:embed schemas/*.json
var JSONSchemas embed.FS

const schemaVersion = "commonmeta_v0.15"

// JSONSchemaErrors validates a JSON document against a JSON Schema file.
func JSONSchemaErrors(document []byte, schema ...string) *gojsonschema.Result {

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
	data, err := JSONSchemas.ReadFile(filepath.Join(dir, s+".json"))
	if err != nil {
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		fmt.Println(exPath)
		fmt.Print(err)
	}
	schemaLoader := gojsonschema.NewStringLoader(string(data))
	documentLoader := gojsonschema.NewBytesLoader(document)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		fmt.Print(err)
		panic(err.Error())
	}
	return result
}
