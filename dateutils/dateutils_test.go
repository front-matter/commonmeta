package dateutils_test

import (
	"testing"

	"github.com/front-matter/commonmeta/dateutils"
)

// func TestGetDateParts(t *testing.T) {
// 	t.Parallel()
// 	type testCase struct {
// 		date string
// 		want map[string][]int
// 	}
// 	testCases := []testCase{
// 		{date: "2021-01-22", want: {"date_parts": [[2021, 1, 22]]}},
// 		{date: "2021-01", want: {"date_parts": [[2021, 1]]}},
// 		{date: "2021", want: {"date_parts": [[2021]]}},
// 		{date: nil, want: {"date_parts": nil}},
// 	}
// 	for _, tc := range testCases {
// 		got := date_utils.GetDateParts(tc.date)
// 		for i := 0; i < 3; i++ {
// 			if tc.want[i] != got[i] {
// 				t.Errorf("Get date parts(%s) from date: want %v, got %v",
// 					tc.date, tc.want, got)
// 			}
// 		}
// 	}
// }

func TestGetDateFromUnixTimestamp(t *testing.T) {
	t.Parallel()
	type testCase struct {
		timestamp int64
		want      string
	}
	testCases := []testCase{
		{timestamp: 0, want: "1970-01-01"},
		{timestamp: 1611312000, want: "2021-01-22"},
	}
	for _, tc := range testCases {
		got := dateutils.GetDateFromUnixTimestamp(tc.timestamp)
		if tc.want != got {
			t.Errorf("Get Date From Unix Timestamp(%d): want %s, got %s",
				tc.timestamp, tc.want, got)
		}
	}
}
