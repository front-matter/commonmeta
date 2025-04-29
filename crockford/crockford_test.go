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
		{input: 736381604818, want: "9ed5m-ytn30", splitEvery: 5, length: 10, checksum: true},
		{input: 258706475165200172, want: "75rw5cg-n1bsc64", splitEvery: 7, length: 14, checksum: true},
		{input: 161006169, want: "4shg-js75", splitEvery: 4, length: 8, checksum: true},
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
		{input: "twwjw-1ww98", want: 924377286556, checksum: true},
		{input: "9ed5m-ytn", want: 324712168277, checksum: false},
		{input: "9ed5m-ytn30", want: 324712168277, checksum: true},
		{input: "elife.01567", want: 0, checksum: false},
	}
	for _, tc := range testCases {
		got, err := crockford.Decode(tc.input, tc.checksum)
		if tc.want != got {
			t.Errorf("Decode(%v): want %v, got %v, error %v",
				tc.input, tc.want, got, err)
		}
	}
}

func TestNormalize(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input string
		want  string
	}
	testCases := []testCase{
		{input: "f9ZQNSF065", want: "f9zqnsf065"},
		{input: "f9zqn-sf065", want: "f9zqnsf065"},
		{input: "f9Llio", want: "f91110"},
	}
	for _, tc := range testCases {
		got := crockford.Normalize(tc.input)
		if tc.want != got {
			t.Errorf("Normalize(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}

func TestGenerateChecksum(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input int64
		want  int64
	}
	testCases := []testCase{
		{input: 450320459383, want: 85},
		{input: 123456789012, want: 44},
	}
	for _, tc := range testCases {
		got := crockford.GenerateChecksum(tc.input)
		if tc.want != got {
			t.Errorf("GenerateChecksum(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}

func TestValidate(t *testing.T) {
	t.Parallel()
	type testCase struct {
		input    int64
		checksum int64
		want     bool
	}
	testCases := []testCase{
		{input: 375301249367, checksum: 92, want: true},
		{input: 930412369850, checksum: 36, want: true},
	}
	for _, tc := range testCases {
		got := crockford.Validate(tc.input, tc.checksum)
		if tc.want != got {
			t.Errorf("Validate(%v): want %v, got %v",
				tc.input, tc.want, got)
		}
	}
}
