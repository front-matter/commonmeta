package ghost_test

import (
	"strings"
	"testing"

	"github.com/front-matter/commonmeta/ghost"
)

func TestGenerateGhostToken(t *testing.T) {
	t.Parallel()
	type testCase struct {
		key  string
		want string
	}
	testCases := []testCase{
		{key: "abc12345678:0018256bcfdcde6ef155481e96ab45b907e147ca562369d823fbc4a22755fa18", want: "eyJhbGciOiJIUzI1NiIsImtpZCI6ImFiYzEyMzQ1Njc4IiwidHlwIjoiSldUIn0"},
	}
	for _, tc := range testCases {
		got, err := ghost.GenerateGhostToken(tc.key)
		if err != nil {
			t.Errorf("GenerateGhostToken(%v): error %v", tc.key, err)
		}
		header := strings.Split(got, ".")[0]
		if tc.want != header {
			t.Errorf("Generate Ghost Token(%v): want %v, got %v",
				tc.key, tc.want, header)
		}
	}
}

// func TestUpdateGhostPost(t *testing.T) {
// 	t.Parallel()
// 	type testCase struct {
// 		postID string
// 		post   ghost.Post
// 		want   ghost.Post
// 	}
// 	testCases := []testCase{
// 		{postID: "5f4d3c1f5d1e2b001f4d3c1f", post: ghost.Post{Title: "Test Post", Slug: "test-post", HTML: "<p>Test post</p>", PublishedAt: "2020-09-01T00:00:00Z", UpdatedAt: "2020-09-01T00:00:00Z"}, want: ghost.Post{Title: "Test Post", Slug: "test-post", HTML: "<p>Test post</p>", PublishedAt: "2020-09-01T00:00:00Z", UpdatedAt: "2020-09-01T00:00:00Z"}},
// 	}
// 	for _, tc := range testCases {
// 		got, err := ghost.UpdateGhostPost(tc.postID, tc.post)
// 		if fmt.Sprintf("%v", tc.want) != fmt.Sprintf("%v", got) {
// 			t.Errorf("Update Ghost Post(%v): want %v, got %v, error %v",
// 				tc.postID, tc.want, got, err)
// 		}
// 	}
// }
