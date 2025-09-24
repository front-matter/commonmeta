package vocabularies

import (
	"embed"
	"fmt"
)

//go:embed *.json *.parquet
var Files embed.FS

var RORFile = "ror_v1.71.parquet"
var SPDXFile = "licenses.json"

func LoadVocabulary(name string) ([]byte, error) {
	var output []byte

	switch name {
	case "ROR.Organizations":
		output, _ = Files.ReadFile(RORFile)
	case "SPDX.Licenses":
		output, _ = Files.ReadFile(SPDXFile)
	default:
		return output, fmt.Errorf("unsupported vocabulary: %s", name)
	}
	return output, nil
}
