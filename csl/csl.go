package csl

import "github.com/front-matter/commonmeta-go/types"

type Content struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func ReadCsl(content Content) (types.Data, error) {
	var data types.Data
	return data, nil
}
