package codemeta

import "github.com/front-matter/commonmeta/types"

type Content struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func GetCodemeta(id string) (Content, error) {
	var content Content
	return content, nil
}

func ReadCodemeta(content Content) (types.Data, error) {
	var data types.Data
	return data, nil
}
