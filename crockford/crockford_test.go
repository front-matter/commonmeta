package crockford_test

import (
	"testing"

	"github.com/front-matter/commonmeta/crockford"
)

func TestEncode(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input      int64
		splitEvery int
		length     int
		checksum   bool
		want       string
	}
	testCases := []testCase{
		{input: 0, want: "0", splitEvery: 0, length: 0, checksum: false},
		{input: 1234, want: "16j", splitEvery: 0, length: 0, checksum: false},
		{input: 1234, want: "16-j", splitEvery: 2, length: 0, checksum: false},
		{input: 1234, want: "01-6j", splitEvery: 2, length: 4, checksum: false},
		{input: 538751765283013, want: "f9zqn-sf065", splitEvery: 5, length: 10, checksum: false},
		{input: 712266282077, want: "mqb61-x2x15", splitEvery: 5, length: 10, checksum: true},
	}
	for _, tc := range testCases {
		got := crockford.Encode(tc.input, tc.splitEvery, tc.length, tc.checksum)
		if tc.want != got {
			t.Errorf("Encode(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}

func TestGenerate(t *testing.T) {
	t.Parallel()
	type testCase struct {
		length     int
		splitEvery int
		checksum   bool
	}
	testCases := []testCase{
		{length: 4, splitEvery: 0, checksum: false},
		{length: 10, splitEvery: 5, checksum: false},
		{length: 10, splitEvery: 5, checksum: true},
	}
	for _, tc := range testCases {
		got := crockford.Generate(tc.length, tc.splitEvery, tc.checksum)
		length := tc.length
		if tc.splitEvery > 0 {
			length += tc.length/tc.splitEvery - 1
		}
		if len(got) != length {
			t.Errorf("Generate(%v): want length %v, got %v",
				tc.length, tc.length, len(got))
		}
	}
}

func TestDecode(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input    string
		checksum bool
		want     int64
	}
	testCases := []testCase{
		{input: "0", want: 0, checksum: false},
		{input: "16j", want: 1234, checksum: false},
		{input: "16-j", want: 1234, checksum: false},
		{input: "01-6j", want: 1234, checksum: false},
		{input: "f9zqn-sf065", want: 538751765283013, checksum: false},
		{input: "mqb61-x2x15", want: 712266282077, checksum: true},
		{input: "axgv5-6aq97", want: 375301249367, checksum: true},
	}
	for _, tc := range testCases {
		got, err := crockford.Decode(tc.input, tc.checksum)
		if tc.want != got {
			t.Errorf("Decode(%v): want %v, got %v, error %v",
				tc.input, tc.want, got, err)
		}
	}
}
