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

//go:embed schemas/*.json
var JsonSchemas embed.FS

func JSONSchemaErrors(document []byte, schema ...string) *gojsonschema.Result {
	// If no schema is provided, default to commonmeta_v0.13
	if len(schema) == 0 {
		schema = append(schema, "commonmeta_v0.13")
	}
	s := schema[len(schema)-1]
	// JSON Schema files stored locally to validate against
	schemata := []string{"commonmeta_v0.13", "datacite-v4.5", "crossref-v0.2", "csl-data", "cff_v1.2.0"}
	if !slices.Contains(schemata, s) {
		log.Fatalf("Schema %s not found", s)
	}
	dir := "schemas"
	data, err := JsonSchemas.ReadFile(filepath.Join(dir, s+".json"))
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
