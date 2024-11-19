package roguescholar

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/front-matter/commonmeta/commonmeta"
)

// UpdateLegacyRecord updates a record in Rogue Scholar legacy database.
func UpdateLegacyRecord(record commonmeta.APIResponse, legacyKey string, field string) (commonmeta.APIResponse, error) {
	if legacyKey == "" {
		return record, fmt.Errorf("no legacy key provided")
	}
	now := strconv.FormatInt(time.Now().Unix(), 10)
	var output []byte
	if field == "rid" && record.ID != "" {
		output = []byte(`{"rid":"` + record.ID + `", "indexed_at":"` + now + `", "indexed":"true"}`)
	} else if record.DOI != "" {
		output = []byte(`{"doi":"` + record.DOI + `", "indexed_at":"` + now + `", "indexed":"true"}`)
	} else {
		return record, fmt.Errorf("no valid field to update")
	}
	requestURL := fmt.Sprintf("https://db.rogue-scholar.org/rest/v1/posts?id=eq.%s", record.UUID)
	req, _ := http.NewRequest(http.MethodPatch, requestURL, bytes.NewReader(output))
	req.Header = http.Header{
		"Content-Type":  {"application/json"},
		"apikey":        {legacyKey},
		"Authorization": {"Bearer " + legacyKey},
		"Prefer":        {"return=minimal"},
	}
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := client.Do(req)
	if resp.StatusCode != 204 {
		return record, err
	}
	record.Status = "updated_legacy"
	return record, nil
}
