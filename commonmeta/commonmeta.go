package commonmeta

import (
	"commonmeta/schemautils"
	"commonmeta/types"
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

func ReadCommonmeta(content types.Content) (types.Data, error) {
	var data types.Data
	return data, nil
}

func WriteCommonmeta(data types.Data) ([]byte, []gojsonschema.ResultError) {
	output, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	validation := schemautils.JSONSchemaErrors(output)
	if !validation.Valid() {
		return nil, validation.Errors()
	}
	return output, nil
}

func WriteCommonmetaList(list []types.Data) ([]byte, []gojsonschema.ResultError) {
	output, err := json.Marshal(list)
	if err != nil {
		fmt.Println(err)
	}
	validation := schemautils.JSONSchemaErrors(output)
	if !validation.Valid() {
		return nil, validation.Errors()
	}
	return output, nil
}
