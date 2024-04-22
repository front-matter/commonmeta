package cff

import "commonmeta-go/types"

type Content struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func GetCFF(id string) (Content, error) {
	var content Content
	return content, nil
}

func ReadCFF(content Content) (types.Data, error) {
	var data types.Data
	return data, nil
}
