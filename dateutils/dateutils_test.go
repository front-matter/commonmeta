package dateutils_test

import (
	"fmt"
	"testing"

	"github.com/front-matter/commonmeta/dateutils"
	"github.com/google/go-cmp/cmp"
)

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

func TestGetDateParts(t *testing.T) {
	t.Parallel()
	type testCase struct {
		date string
		want map[string][][]int
	}
	testCases := []testCase{
		{date: "2021-01-22", want: map[string][][]int{"date-parts": {{2021, 1, 22}}}},
		{date: "2021-01", want: map[string][][]int{"date-parts": {{2021, 1, 0}}}},
		{date: "2021", want: map[string][][]int{"date-parts": {{2021, 0, 0}}}},
		{date: "", want: map[string][][]int{"date-parts": {}}},
	}
	for _, tc := range testCases {
		got := dateutils.GetDateParts(tc.date)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("GetDateParts mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestGetDateStruct(t *testing.T) {
	t.Parallel()
	type testCase struct {
		date string
		want dateutils.DateStruct
	}
	testCases := []testCase{
		{date: "2021-01-22", want: dateutils.DateStruct{Year: "2021", Month: "01", Day: "22"}},
		{date: "2021-01", want: dateutils.DateStruct{Year: "2021", Month: "01", Day: "00"}},
		{date: "2021", want: dateutils.DateStruct{Year: "2021", Month: "00", Day: "00"}},
		{date: "", want: dateutils.DateStruct{Year: "", Month: "", Day: ""}},
	}
	for _, tc := range testCases {
		got := dateutils.GetDateStruct(tc.date)
		if diff := cmp.Diff(tc.want, got); diff != "" {
			t.Errorf("GetDateStruct mismatch (-want +got):\n%s", diff)
		}
	}
}

func ExampleGetDateParts() {
	m := dateutils.GetDateParts("2023-12-06")
	fmt.Println(m)
	// Output:
	// map[date-parts:[[2023 12 6]]]
}

func ExampleGetDateFromUnixTimestamp() {
	s := dateutils.GetDateFromUnixTimestamp(1611312000)
	fmt.Println(s)
	// Output:
	// 2021-01-22
}

func ExampleStripMilliseconds() {
	s := dateutils.StripMilliseconds("2021-01-22T10:00:00.000Z")
	fmt.Println(s)
	// Output:
	// 2021-01-22T10:00:00Z
}

func ExampleParseDate() {
	s := dateutils.ParseDate("2021-01-22T10:00:00")
	fmt.Println(s)
	// Output:
	// 2021-01-22
}
