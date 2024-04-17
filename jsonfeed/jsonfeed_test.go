package jsonfeed_test

import (
	"commonmeta/jsonfeed"
	"commonmeta/types"
	"testing"
)

func TestGetJsonFeedItem(t *testing.T) {
	t.Parallel()

	type testCase struct {
		pid  string
		want string
		err  error
	}

	ghostPost := types.Content{
		PID:   "5adbb6d4-1fe2-4da2-8cf4-c897f88a02d9",
		Title: "INFORMATE: Where Are the Data?",
	}
	wordpressPost := types.Content{
		PID:   "4e4bf150-751f-4245-b4ca-fe69e3c3bb24",
		Title: "New paper: Curtice et al. (2023) on the first <i>Haplocanthosaurus</i> from Dry Mesa",
	}

	testCases := []testCase{
		{pid: ghostPost.PID, want: ghostPost.Title, err: nil},
		{pid: wordpressPost.PID, want: wordpressPost.Title, err: nil},
	}
	for _, tc := range testCases {
		got, err := jsonfeed.GetJsonFeedItem(tc.pid)
		if tc.want != got.Title {
			t.Errorf("JSON Feed ID(%v): want %v, got %v, error %v",
				tc.pid, tc.want, got, err)
		}
	}
}
