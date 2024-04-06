package jsonfeed

import (
	"commonmeta/metadata"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// the JSON response containing the JSON Feed metadata
type Record struct {
	ID          string   `json:"id"`
	Abstract    string   `json:"abstract"`
	Authors     []Author `json:"authors"`
	Blog        Blog     `json:"blog"`
	DOI         string   `json:"doi"`
	Files       []File   `json:"files"`
	GUID        string   `json:"guid"`
	Language    string   `json:"language"`
	PublishedAt string   `json:"published_at"`
	References  string   `json:"reference"`
	Relations   string   `json:"relations"`
	Summary     string   `json:"summary"`
	Tags        []string `json:"tags"`
	Title       string   `json:"title"`
	UpdatedAt   string   `json:"updated_at"`
	Url         string   `json:"url"`
}

type Author struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Blog struct {
	ISSN        string `json:"issn"`
	License     string `json:"license"`
	Title       string `json:"title"`
	HomePageUrl string `json:"home_page_url"`
}

type File struct {
	Url string `json:"url"`
}

var result Record

func GetJsonFeedItem(pid string) (Record, error) {
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

func ReadJsonFeedItem(record Record) (metadata.Metadata, error) {
	var m metadata.Metadata
	return m, nil
}
