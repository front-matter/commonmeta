package schemautils

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func JSONSchemaErrors(document []byte) *gojsonschema.Result {
	schemaLoader := gojsonschema.NewReferenceLoader("file://../resources/commonmeta_v0.12.json")
	documentLoader := gojsonschema.NewGoLoader(document)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		panic(err.Error())
	}

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
	return result
}
