// Package dateutils provides functions to work with dates.
package dateutils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type DateStruct struct {
	Year  int
	Month int
	Day   int
}

type FormattedDateStruct struct {
	Year  string
	Month string
	Day   string
}

type DateSlice []interface{}

// Iso8601DateFormat is the ISO 8601 date format without time.
const Iso8601DateFormat = "2006-01-02"

// Iso8601DateMonthFormat is the ISO 8601 date format without time and day.
const Iso8601DateMonthFormat = "2006-01"

// Iso8601DateYearFormat is the ISO 8601 date format without time, month and day.
const Iso8601DateYearFormat = "2006"

// Iso8601DateTimeFormat is the ISO 8601 date format with time.
const Iso8601DateTimeFormat = "2006-01-02T15:04:05Z"

// CrossrefDateTimeFormat is the Crossref date format with time, used in XML for content registration.
const CrossrefDateTimeFormat = "20060102150405"

// DateTimeWithSpacesFormat is the date time format with spaces and no timezone information.
const DateTimeWithSpacesFormat = "2006-01-02 15:04:05"

// ParseDate parses date strings in various formats and returns a date string in ISO 8601 format.
func ParseDate(iso8601Time string) string {
	date := GetDateStruct(iso8601Time)
	if date.Year == 0 {
		return ""
	}
	dateStr := fmt.Sprintf("%04d", date.Year)
	if date.Month != 0 {
		dateStr += "-" + fmt.Sprintf("%02d", date.Month)
	}
	if date.Day != 0 {
		dateStr += "-" + fmt.Sprintf("%02d", date.Day)
	}
	return dateStr
}

// ParseDateTime parses datetime strings in various formats and returns a datetime string in ISO 8601 format.
func ParseDateTime(input string) string {
	time, err := ParseTime(input)
	if err != nil {
		return ""
	}
	if time.Hour() == 0 && time.Minute() == 0 && time.Second() == 0 {
		return ParseDate(input)
	}
	return time.Format(Iso8601DateTimeFormat)
}

// ParseTime parses datetime strings in various formats and returns a go time, which can then be formatted to an iso8601 string using ParseDateTime.
func ParseTime(input string) (time.Time, error) {
	formats := []string{Iso8601DateTimeFormat, Iso8601DateFormat, Iso8601DateMonthFormat, Iso8601DateYearFormat, DateTimeWithSpacesFormat, CrossrefDateTimeFormat}
	for _, format := range formats {
		datetime, err := time.Parse(format, input)
		if err == nil {
			return datetime, nil
		}
	}
	return time.Time{}, errors.New("unrecognized time format")
}

// GetDateParts return date parts from an ISO 8601 date string
func GetDateParts(iso8601Time string) []DateSlice {
	var dateParts []DateSlice
	if iso8601Time == "" {
		return dateParts
	}

	// optionally add missing zeros to the date string
	if len(iso8601Time) < 10 {
		iso8601Time = iso8601Time + strings.Repeat("0", 10-len(iso8601Time))
	}
	year, _ := strconv.Atoi(iso8601Time[0:4])
	month, _ := strconv.Atoi(iso8601Time[5:7])
	day, _ := strconv.Atoi(iso8601Time[8:10])
	dateParts = append(dateParts, DateSlice{year, month, day})
	return dateParts
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
	year, _ := strconv.Atoi(iso8601Time[0:4])
	month, _ := strconv.Atoi(iso8601Time[5:7])
	day, _ := strconv.Atoi(iso8601Time[8:10])

	return DateStruct{
		Year:  year,
		Month: month,
		Day:   day,
	}
}

// GetFormattedDateStruct returns a formatted date struct from an ISO 8601 date string
func GetFormattedDateStruct(iso8601Time string) FormattedDateStruct {
	var year, month, day string
	dateStruct := GetDateStruct(iso8601Time)
	if dateStruct.Year == 0 {
		return FormattedDateStruct{}
	}
	year = fmt.Sprintf("%04d", dateStruct.Year)
	if dateStruct.Month != 0 {
		month = fmt.Sprintf("%02d", dateStruct.Month)
	}
	if dateStruct.Day != 0 {
		day = fmt.Sprintf("%02d", dateStruct.Day)
	}
	return FormattedDateStruct{
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
// uses interface{} to allow for float64 and string types
func GetDateFromDateParts(dateAsParts []DateSlice) string {
	dateParts := dateAsParts[0]
	length := len(dateParts)
	parts := make([]float64, length)
	if length == 0 || dateParts[0] == nil {
		return ""
	}
	for i, part := range dateParts {
		switch v := part.(type) {
		case string:
			parts[i], _ = strconv.ParseFloat(part.(string), 64)
		case int:
			parts[i] = float64(part.(int))
		case float64:
			parts[i] = part.(float64)
		default:
			fmt.Printf("I don't know about type %T!\n", v)
		}
	}
	if parts[0] == 0 {
		return ""
	} else if length == 1 {
		return GetDateFromParts(int(parts[0]))
	} else if length == 2 {
		return GetDateFromParts(int(parts[0]), int(parts[1]))
	} else if length == 3 {
		return GetDateFromParts(int(parts[0]), int(parts[1]), int(parts[2]))
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

// GetUnixTimestampFromDatetime returns a Unix timestamp from a datetime
func GetUnixTimestampFromDatetime(iso8601Time string) int64 {
	time, _ := time.Parse(Iso8601DateTimeFormat, iso8601Time)
	return time.Unix()
}

// StripMilliseconds removes milliseconds from an ISO 8601 datetime string
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
