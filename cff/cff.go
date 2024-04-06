package cff

import "commonmeta/metadata"

type Record struct {
	URL string `json:"URL"`
}

var result Record

func GetCFF(pid string) (Record, error) {
	var r Record
	return r, nil
}

func ReadCFF(record Record) (metadata.Metadata, error) {
	var m metadata.Metadata
	return m, nil
}
