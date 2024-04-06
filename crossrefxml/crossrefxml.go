package crossrefxml

import "commonmeta/metadata"

type Record struct {
	URL string `json:"URL"`
}

func GetCrossrefXML(pid string) (Record, error) {
	var r Record
	return r, nil
}

func ReadCrossrefXML(record Record) (metadata.Metadata, error) {
	var m metadata.Metadata
	return m, nil
}
