package crossrefxml

import (
	"github.com/front-matter/commonmeta-go/types"
)

type Content struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func GetCrossrefXML(pid string) (Content, error) {
	var result Content
	return result, nil
}

func ReadCrossrefXML(content Content) (types.Data, error) {
	var data types.Data
	return data, nil
}
