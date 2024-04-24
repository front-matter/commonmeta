// Package jsonfeed converts JSON Feed metadata to/from the commonmeta metadata format.
package jsonfeed

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
)

// Content represents the JSON Feed metadata.
type Content struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// Get retrieves JSON Feed metadata.
func Get(id string) (Content, error) {
	var content Content
	client := http.Client{
		Timeout: time.Second * 10,
	}
	url := "https://api.rogue-scholar.org/posts/" + id
	resp, err := client.Get(url)
	if err != nil {
		return content, err
	}
	if resp.StatusCode != 200 {
		return content, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return content, err
	}
	err = json.Unmarshal(body, &content)
	if err != nil {
		fmt.Println("error:", err)
	}
	return content, err
}

// Read reads JSON Feed metadata and converts it into Commonmeta metadata.
func Read(content Content) (commonmeta.Data, error) {
	var data commonmeta.Data
	data.ID = content.ID
	return data, nil
}
