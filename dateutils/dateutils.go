// Package dateutils provides functions to work with dates.
package dateutils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type DateStruct struct {
	Year  string
	Month string
	Day   string
}

// Iso8601DateFormat is the ISO 8601 date format without time.
const Iso8601DateFormat = "2006-01-02"

// Iso8601DateTimeFormat is the ISO 8601 date format with time.
const Iso8601DateTimeFormat = "2006-01-02T15:04:05Z"

// CrossrefDateTimeFormat is the Crossref date format with time, used in XML for content registration.
const CrossrefDateTimeFormat = "20060102150405"

// ParseDate parses date strings in various formats and returns a date string in ISO 8601 format
func ParseDate(date string) string {
	t, _ := time.Parse(Iso8601DateFormat, date)
	if t.Year() == 0001 {
		t, _ = time.Parse("02 January 2006", date)
	}
	if t.Year() == 0001 {
		t, _ = time.Parse("2006-02", date)
	}
	if t.Year() == 0001 {
		t, _ = time.Parse("2006", date)
	}
	if t.Year() == 0001 {
		return ""
	}
	return t.Format(Iso8601DateFormat)
}

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

// GetDateStruct returns struct with date (year, month, day) from an ISO 8601 date string
func GetDateStruct(iso8601Time string) DateStruct {
	if iso8601Time == "" {
		return DateStruct{}
	}

	// optionally add missing zeros to the date string
	if len(iso8601Time) < 10 {
		iso8601Time = iso8601Time + strings.Repeat("0", 10-len(iso8601Time))
	}
	year := iso8601Time[0:4]
	month := iso8601Time[5:7]
	day := iso8601Time[8:10]

	return DateStruct{
		Year:  year,
		Month: month,
		Day:   day,
	}
}

// GetDateFromUnixTimestamp returns a date string from a Unix timestamp
func GetDateFromUnixTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(Iso8601DateFormat)
}

// GetDateTimeFromUnixTimestamp returns a datetime string from a Unix timestamp
func GetDateTimeFromUnixTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(Iso8601DateTimeFormat)
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

// GetDateFromDatetime returns a datetime string from a Unix timestamp
func GetDateFromDatetime(iso8601Time string) string {
	date, _ := time.Parse(Iso8601DateTimeFormat, iso8601Time)
	return date.Format(Iso8601DateFormat)
}

func StripMilliseconds(iso8601Time string) string {
	if iso8601Time == "" {
		return ""
	}
	if strings.Contains(iso8601Time, "T00:00:00") {
		return strings.Split(iso8601Time, "T")[0]
	}
	if strings.Contains(iso8601Time, ".") {
		return strings.Split(iso8601Time, ".")[0] + "Z"
	}
	if strings.Contains(iso8601Time, "+00:00") {
		return strings.Split(iso8601Time, "+")[0] + "Z"
	}
	return iso8601Time
}

// ValidateEdtf validates an EDTF date string.
// Workaround for a bug in InvenioRDM (edtf 4.0.1)
func ValidateEdtf(iso8601Time string) string {
	if iso8601Time == "" {
		return ""
	}
	if strings.Contains(iso8601Time, "T23") {
		return ""
	}
	return iso8601Time
}
