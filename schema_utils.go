package schema_utils

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func main() {

	schemaLoader := gojsonschema.NewReferenceLoader("file:///home/me/schema.json")
	documentLoader := gojsonschema.NewReferenceLoader("file:///home/me/document.json")

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
}
