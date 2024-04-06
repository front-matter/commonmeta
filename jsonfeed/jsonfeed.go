package jsonfeed

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// the JSON response containing the JSON Feed metadata
type Record struct {
	ID    string `json:"id"`
	DOI   string `json:"doi"`
	Title string `json:"title"`
	Url   string `json:"url"`
}

var result Record

func GetJsonfeedItem(pid string) (Record, error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	url := "https://api.rogue-scholar.org/posts/" + pid
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
