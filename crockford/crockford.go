package crockford

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strconv"
	"strings"
)

// Generate, encode and decode random base32 identifiers.
// This encoder/decoder:
// - uses Douglas Crockford Base32 encoding: https://www.crockford.com/base32.html
// - allows for ISO 7064 checksum
// - encodes the checksum using only characters in the base32 set
// - produces string that are URI-friendly (no '=' or '/' for instance)
// This is based on: https://github.com/front-matter/base32-url

// NO i, l, o or u
const ENCODING_CHARS = "0123456789abcdefghjkmnpqrstvwxyz"

// Encode a number to a URI-friendly Douglas Crockford base32 string.
// optionally split with ' -' every n characters, pad with zeros to a minimum length,
// and append a checksum using modulo 97-10 (ISO 7064).
func Encode(number int64, splitEvery int, length int, checksum bool) string {
	encoded := ""
	originalNumber := number
	if number == 0 {
		encoded = "0"
	} else {
		for number > 0 {
			remainder := number % 32
			number /= 32
			encoded = string(ENCODING_CHARS[remainder]) + encoded
		}
	}

	if checksum && length > 2 {
		length -= 2
	}
	if length > 0 && len(encoded) < length {
		encoded = strings.Repeat("0", length-len(encoded)) + encoded
	}

	if checksum {
		computedChecksum := GenerateChecksum(originalNumber)
		encoded += fmt.Sprintf("%02d", computedChecksum)
	}

	if splitEvery > 0 {
		var splits []string
		for i := 0; i < len(encoded); i += splitEvery {
			nn := i + splitEvery
			if nn > len(encoded) {
				nn = len(encoded)
			}
			splits = append(splits, string(encoded[i:nn]))
		}
		encoded = strings.Join(splits, "-")
	}

	return encoded
}

// Generate a random Cockroft base32 string.
// optionally split with ' -' every n characters, pad with zeros to a minimum length,
// and append a checksum using modulo 97-10 (ISO 7064).
func Generate(length int, splitEvery int, checksum bool) string {
	if checksum && length < 3 {
		panic("Invalid 'length'. Must be >= 3 if checksum enabled.")
	}
	// fixes number size, otherwise decoding checksum check will fail
	if checksum {
		length -= 2
	}
	// generate a random number between 0 and 32^length
	n := math.Pow(float64(32), float64(length))
	number := int64(rand.IntN(int(n)))
	str := Encode(number, splitEvery, length, checksum)
	return str
}

// Decode a URI-friendly Douglas Crockford base32 string to a number.
func Decode(str string, checksum bool) (int64, error) {
	var encoded string
	var number, cs int64
	var ok bool
	var err error

	encoded = Normalize(str)
	if checksum {
		// checksum is the last two characters
		cs, err = strconv.ParseInt(encoded[len(encoded)-2:], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid checksum: %s", encoded[len(encoded)-2:])
		}
		encoded = encoded[:len(encoded)-2]
	}
	for _, c := range encoded {
		number *= 32
		pos := strings.Index(ENCODING_CHARS, string(c))
		// invalid character, stop decoding and return 0
		if pos == -1 {
			return 0, fmt.Errorf("invalid character: %s", string(c))
		}
		number += int64(pos)
	}
	if checksum {
		ok = Validate(number, cs)
		if !ok {
			return 0, fmt.Errorf("wrong checksum %02d for identifier %s", cs, str)
		}
	}
	return number, err
}

// Normalize returns a normalized encoded string for base32 encoding.
func Normalize(str string) string {
	normalized := strings.ToLower(str)
	normalized = strings.ReplaceAll(normalized, "-", "")
	normalized = strings.ReplaceAll(normalized, "i", "1")
	normalized = strings.ReplaceAll(normalized, "l", "1")
	normalized = strings.ReplaceAll(normalized, "o", "0")
	return normalized
}

// Validate returns true if the encoded string is a valid base32 string with checksum.
func Validate(number int64, checksum int64) bool {
	return checksum == GenerateChecksum(number)
}

// GenerateChecksum returns the checksum for a number using ISO 7064 (mod 97-10).
func GenerateChecksum(number int64) int64 {
	checksum := 97 - ((100 * number) % 97) + 1
	return checksum
}
