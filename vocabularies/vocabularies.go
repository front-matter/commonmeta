package vocabularies

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/front-matter/commonmeta/fileutils"
)

//go:embed *.zip *.json
var Files embed.FS

var RORFile = "v1.68-2025-07-15-ror-data_schema_v2.json"
var SPDXFile = "licenses.json"

func LoadVocabulary(name string) ([]byte, error) {
	var input, output []byte
	var err error

	switch name {
	case "ROR.Organizations":
		input, _ = Files.ReadFile(filepath.Join(RORFile + ".zip"))
		output, err = fileutils.UnzipContent(input, RORFile)
		if err != nil {
			return nil, err
		}
	case "SPDX.Licenses":
		output, _ = Files.ReadFile(SPDXFile)
	default:
		return output, fmt.Errorf("unsupported vocabulary: %s", name)
	}
	return output, nil
}
