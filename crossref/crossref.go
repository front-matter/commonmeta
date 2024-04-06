package crossref

import (
	"commonmeta/doiutils"
	"commonmeta/metadata"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// the envelope for the JSON response from the Crossref API
type Result struct {
	Status         string `json:"status"`
	MessageType    string `json:"message-type"`
	MessageVersion string `json:"message-version"`
	Message        Record `json:"message"`
}

// the JSON response containing the metadata for the DOI
type Record struct {
	URL       string   `json:"URL"`
	DOI       string   `json:"DOI"`
	Type      string   `json:"type"`
	Title     []string `json:"title"`
	Publisher string   `json:"publisher"`
	Volume    string   `json:"volume"`
	Issue     string   `json:"issue"`
	Page      string   `json:"page"`
}

var result Result

func GetCrossref(pid string) (Record, error) {
	doi, err := doiutils.DOIFromUrl(pid)
	if err != nil {
		return Record{}, err
	}
	url := "https://api.crossref.org/works/" + doi
	client := http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(url)
	if err != nil {
		return Record{}, err
	}
	if resp.StatusCode >= 400 {
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
	return result.Message, err
}

func ReadCrossref(record Record) (metadata.Metadata, error) {
	var m metadata.Metadata
	return m, nil
}
