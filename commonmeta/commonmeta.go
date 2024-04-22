package commonmeta

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/front-matter/commonmeta/schemautils"
	"github.com/front-matter/commonmeta/types"

	"github.com/xeipuuv/gojsonschema"
)

func ReadCommonmeta(content types.Data) (types.Data, error) {
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
	for _, data := range list {
		o, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
		}
		validation := schemautils.JSONSchemaErrors(o)
		if !validation.Valid() {
			var out bytes.Buffer
			json.Indent(&out, o, "=", "\t")
			fmt.Println(out.String())
			return nil, validation.Errors()
		}
	}

	output, err := json.Marshal(list)
	if err != nil {
		fmt.Println(err)
	}
	return output, nil
}
