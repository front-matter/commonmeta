// Package dateutils provides functions to work with dates.
package dateutils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Iso8601DateFormat is the ISO 8601 date format without time.
const Iso8601DateFormat = "2006-01-02"

// func MonthNames() struct {
// 	return {
// 		"01": "jan",
// 		"02": "feb",
// 		"03": "mar",
// 		"04": "apr",
// 		"05": "may",
// 		"06": "jun",
// 		"07": "jul",
// 		"08": "aug",
// 		"09": "sep",
// 		"10": "oct",
// 		"11": "nov",
// 		"12": "dec",
// 	}
// }

// func MonthShortNames() slice {
// 	return [
// 		"jan",
// 		"feb",
// 		"mar",
// 		"apr",
// 		"may",
// 		"jun",
// 		"jul",
// 		"aug",
// 		"sep",
// 		"oct",
// 		"nov",
// 		"dec",
// 	]
// }

// GetDateParts return date parts from an ISO 8601 date string
func GetDateParts(iso8601Time string) map[string][][]int {
	if iso8601Time == "" {
		return map[string][][]int{"date-parts": {}}
	}

	// optionally add missing zeros to the date string
	if len(iso8601Time) < 10 {
		iso8601Time = iso8601Time + strings.Repeat("0", 10-len(iso8601Time))
	}
	year, _ := strconv.Atoi(iso8601Time[0:4])
	month, _ := strconv.Atoi(iso8601Time[5:7])
	day, _ := strconv.Atoi(iso8601Time[8:10])
	dateParts := [][]int{{year, month, day}}
	return map[string][][]int{"date-parts": dateParts}
}

// GetDateFromUnixTimestamp returns a date string from a Unix timestamp
func GetDateFromUnixTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(Iso8601DateFormat)
}

// GetDateFromDateParts returns a date string from date parts
func GetDateFromDateParts(dateAsParts [][]int) string {
	dateParts := dateAsParts[0]
	switch len(dateParts) {
	case 0:
		return ""
	case 1:
		year := dateParts[0]
		if year == 0 {
			return ""
		}
		return GetDateFromParts(year)
	case 2:
		year, month := dateParts[0], dateParts[1]
		return GetDateFromParts(year, month)
	case 3:
		year, month, day := dateParts[0], dateParts[1], dateParts[2]
		return GetDateFromParts(year, month, day)
	}
	return ""
}

// GetDateFromParts returns a date string from parts
func GetDateFromParts(parts ...int) string {
	var arr []string
	switch len(parts) {
	case 0:
		return ""
	case 1:
		year := fmt.Sprintf("%04d", parts[0])
		arr = []string{year}
	case 2:
		year, month := fmt.Sprintf("%04d", parts[0]), fmt.Sprintf("%02d", parts[1])
		arr = []string{year, month}
	case 3:
		year, month, day := fmt.Sprintf("%04d", parts[0]), fmt.Sprintf("%02d", parts[1]), fmt.Sprintf("%02d", parts[2])
		arr = []string{year, month, day}
	}
	return strings.Join(arr, "-")
}

// GetDateFromCrossrefParts returns a date string from Crossref XML date parts
func GetDateFromCrossrefParts(strParts ...string) string {
	parts := make([]int, 0)
	for _, s := range strParts {
		if s != "" {
			v, _ := strconv.Atoi(s)
			parts = append(parts, v)
		}
	}
	return GetDateFromParts(parts...)
}
