package schemautils

import (
	"commonmeta/metadata"
	"encoding/json"
	"fmt"
	"log"

	"github.com/xeipuuv/gojsonschema"
)

func JSONSchemaErrors(document metadata.Metadata) *gojsonschema.Result {
	documentJSON, err := json.Marshal(document)
	if err != nil {
		log.Fatal(err)
	}
	schemaLoader := gojsonschema.NewReferenceLoader("file://../resources/commonmeta_v0.12.json")
	documentLoader := gojsonschema.NewBytesLoader(documentJSON)
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("Input: %v\n", string(documentJSON))
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
	return result
}
