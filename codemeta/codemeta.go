package codemeta

import "commonmeta/metadata"

type Record struct {
	URL string `json:"URL"`
}

var result Record

func GetCodemeta(pid string) (Record, error) {
	var r Record
	return r, nil
}

func ReadCodemeta(record Record) (metadata.Metadata, error) {
	var m metadata.Metadata
	return m, nil
}
