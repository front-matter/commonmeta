package schemautils

import (
	"fmt"
	"log"
	"slices"

	"github.com/xeipuuv/gojsonschema"
)

func JSONSchemaErrors(document []byte, schema ...string) *gojsonschema.Result {
	// If no schema is provided, default to commonmeta_v0.13
	if len(schema) == 0 {
		schema = append(schema, "commonmeta_v0.12")
	}
	s := schema[len(schema)-1]
	// JSON Schema files stored locally to validate against
	schemata := []string{"commonmeta_v0.12", "datacite-v4.5", "crossref-v0.2", "csl-data", "cff_v1.2.0"}
	if !slices.Contains(schemata, s) {
		log.Fatalf("Schema %s not found", s)
	}
	schemaPath := "file://../resources/" + s + ".json"
	schemaLoader := gojsonschema.NewReferenceLoader(schemaPath)
	documentLoader := gojsonschema.NewBytesLoader(document)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("Input: %v\n", string(document))
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
	return result
}
