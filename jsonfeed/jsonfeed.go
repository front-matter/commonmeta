package jsonfeed

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/front-matter/commonmeta-go/types"
)

type Content struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

func GetJsonFeedItem(id string) (Content, error) {
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

func ReadJsonFeedItem(content Content) (types.Data, error) {
	var data types.Data
	return data, nil
}
