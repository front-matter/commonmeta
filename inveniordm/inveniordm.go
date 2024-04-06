package inveniordm

import (
	"commonmeta/metadata"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// the JSON response containing the InvenioRDM metadata
type Record struct {
	ID    string `json:"id"`
	DOI   string `json:"doi"`
	Title string `json:"title"`
}

var result Record

func GetInvenioRDM(pid string) (Record, error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	url := "https://zenodo.org/api/records/" + pid
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
	return result, err
}

func ReadInvenioRDM(record Record) (metadata.Metadata, error) {
	var m metadata.Metadata
	return m, nil
}
