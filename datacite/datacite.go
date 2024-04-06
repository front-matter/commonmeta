package datacite

import (
	"commonmeta/doiutils"
	"commonmeta/metadata"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// the envelope for the JSON response from the DataCite API
type Result struct {
	Data Record `json:"data"`
}

// the JSON response containing the DataCite metadata
type Record struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	DOI             string    `json:"doi"`
	Prefix          string    `json:"prefix"`
	Suffix          string    `json:"suffix"`
	Creators        []Creator `json:"creators"`
	Publisher       string    `json:"publisher"`
	Container       Container `json:"container"`
	PublicationYear int       `json:"publicationYear"`
	Titles          []Title   `json:"titles"`
	Url             string    `json:"url"`
}

type Container struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Creator struct {
	Type           string `json:"type"`
	Identifier     string `json:"identifier"`
	IdentifierType string `json:"identifierType"`
	Name           string `json:"name"`
}

type Title struct {
	Title    string `json:"title"`
	Language string `json:"language"`
}

var result Result

func GetDatacite(pid string) (Record, error) {
	doi, err := doiutils.DOIFromUrl(pid)
	if err != nil {
		return Record{}, err
	}
	url := "https://api.datacite.org/dois/" + doi
	client := http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(url)
	if err != nil {
		return Record{}, err
	}
	if resp.StatusCode != 200 {
		return Record{}, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Record{}, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("error:", err)
	}
	return result.Data, err
}

func ReadDatacite(record Record) (metadata.Metadata, error) {
	var m metadata.Metadata
	return m, nil
}
