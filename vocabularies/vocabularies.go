package vocabularies

import (
	"embed"
	"fmt"
	"path/filepath"

	"github.com/front-matter/commonmeta/fileutils"
)

//go:embed *.zip
var Files embed.FS

var RORFile = "v1.63-2025-04-03-ror-data.avro"

func LoadVocabulary(name string) ([]byte, error) {
	var output []byte

	switch name {
	case "ROR.Organizations":
		input, err := Files.ReadFile(filepath.Join(RORFile + ".zip"))
		output, err = fileutils.UnzipContent(input, RORFile)
		if err != nil {
			return nil, err
		}
	default:
		return output, fmt.Errorf("unsupported vocabulary: %s", name)
	}
	return output, nil
}
