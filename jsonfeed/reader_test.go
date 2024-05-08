package jsonfeed_test

import (
	"testing"

	"github.com/front-matter/commonmeta/jsonfeed"
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
