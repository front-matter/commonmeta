package jsonfeed_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/front-matter/commonmeta/commonmeta"
	"github.com/front-matter/commonmeta/jsonfeed"
	"github.com/front-matter/commonmeta/utils"
	"github.com/google/go-cmp/cmp"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type testCase struct {
		pid  string
		want string
		err  error
	}

	ghostPost := jsonfeed.Content{
		ID:    "5adbb6d4-1fe2-4da2-8cf4-c897f88a02d9",
		Title: "INFORMATE: Where Are the Data?",
	}
	wordpressPost := jsonfeed.Content{
		ID:    "4e4bf150-751f-4245-b4ca-fe69e3c3bb24",
		Title: "New paper: Curtice et al. (2023) on the first <i>Haplocanthosaurus</i> from Dry Mesa",
	}

	testCases := []testCase{
		{pid: ghostPost.ID, want: ghostPost.Title, err: nil},
		{pid: wordpressPost.ID, want: wordpressPost.Title, err: nil},
	}
	for _, tc := range testCases {
		got, err := jsonfeed.Get(tc.pid)
		if tc.want != got.Title {
			t.Errorf("JSON Feed ID(%v): want %v, got %v, error %v",
				tc.pid, tc.want, got, err)
		}
	}
}

func TestFetch(t *testing.T) {
	t.Parallel()
	type testCase struct {
		name string
		id   string
	}

	testCases := []testCase{
		{name: "blog post with funding", id: "8a4de443-3347-4b82-b57d-e3c82b6485fc"},
		{name: "project blog", id: "4d51f3c9-151d-4030-9893-ddbca37f54bb"},
		{name: "url with uppercase characters", id: "3d02cf64-c600-4eb1-91b4-02f5bade5691"},
	}
	for _, tc := range testCases {
		got, err := jsonfeed.Fetch(tc.id)
		if err != nil {
			t.Errorf("JSON Feed Metadata(%v): error %v", tc.id, err)
		}
		// read json file from testdata folder and convert to Data struct
		id, ok := utils.ValidateUUID(tc.id)
		if !ok {
			t.Fatal("invalid uuid")
		}
		filename := id + ".json"
		filepath := filepath.Join("testdata", filename)
		bytes, err := os.ReadFile(filepath)
		if err != nil {
			t.Fatal(err)
		}

		want := commonmeta.Data{}
		err = json.Unmarshal(bytes, &want)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("Fetch (%s) mismatch (-want +got):\n%s", tc.id, diff)
		}
	}
}
