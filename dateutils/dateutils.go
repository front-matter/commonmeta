package dateutils

import (
	"fmt"
	"strings"
	"time"
)

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

// Get date parts
// func GetDateParts(iso8601Time string) map[string][]int {
// 	if iso8601Time == nil {
// 		return map[string][]int{"date-parts": []int{}}
// 	}
// 	if len(iso8601Time) < 10 {
// 		iso8601Time = strings.Repeat(iso8601Time, 10, "0")
// 	}
// 	year, _ := strconv.Atoi(iso8601Time[0:4])
// 	month, _ := strconv.Atoi(iso8601Time[5:7])
// 	day, _ := strconv.Atoi(iso8601Time[8:10])
// 	dateParts := []int{year, month, day}
// 	return map[string][]int{"date-parts": dateParts}
// }

func GetDateFromUnixTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).Format(Iso8601DateFormat)
}

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
